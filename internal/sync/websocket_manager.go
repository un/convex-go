package sync

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/get-convex/convex-go/internal/protocol"
)

type WebSocketManager struct {
	mu                  sync.Mutex
	wsURL               string
	clientID            string
	conn                *websocket.Conn
	writeQueue          chan []byte
	writeStop           chan struct{}
	writeQueueCapacity  int
	heartbeatStop       chan struct{}
	heartbeatInterval   time.Duration
	inactivityThreshold time.Duration
	lastServerResponse  time.Time
	responses           chan ProtocolResponse
	closed              bool
	connectionCount     uint32
	lastCloseReason     string
}

func NewWebSocketManager(wsURL, clientID string) *WebSocketManager {
	return &WebSocketManager{
		wsURL:               wsURL,
		clientID:            clientID,
		writeQueueCapacity:  256,
		heartbeatInterval:   5 * time.Second,
		inactivityThreshold: 30 * time.Second,
		responses:           make(chan ProtocolResponse, 256),
		lastCloseReason:     "InitialConnect",
	}
}

func (m *WebSocketManager) Open(ctx context.Context, request ReconnectRequest) (<-chan ProtocolResponse, error) {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil, fmt.Errorf("websocket manager closed")
	}
	m.mu.Unlock()
	if err := m.openConn(ctx, request); err != nil {
		return nil, err
	}
	return m.responses, nil
}

func (m *WebSocketManager) Send(ctx context.Context, message protocol.ClientMessage) error {
	payload, err := protocol.EncodeClientMessage(message)
	if err != nil {
		return err
	}

	m.mu.Lock()
	queue := m.writeQueue
	closed := m.closed
	m.mu.Unlock()

	if closed {
		return fmt.Errorf("websocket manager closed")
	}
	if queue == nil {
		return fmt.Errorf("websocket not connected")
	}
	select {
	case queue <- payload:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("websocket write queue full")
	}
}

func (m *WebSocketManager) Reconnect(ctx context.Context, request ReconnectRequest) error {
	m.closeActiveConnection()
	return m.openConn(ctx, request)
}

func (m *WebSocketManager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	conn := m.conn
	stop := m.writeStop
	heartbeatStop := m.heartbeatStop
	m.conn = nil
	m.writeQueue = nil
	m.writeStop = nil
	m.heartbeatStop = nil
	m.lastServerResponse = time.Time{}
	m.mu.Unlock()

	safeClose(stop)
	safeClose(heartbeatStop)
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (m *WebSocketManager) openConn(ctx context.Context, request ReconnectRequest) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("websocket dial cancelled: %w", err)
	}

	headers := http.Header{}
	if m.clientID != "" {
		headers.Set("Convex-Client", m.clientID)
	}

	d := websocket.Dialer{HandshakeTimeout: 20 * time.Second}
	conn, response, err := d.DialContext(ctx, m.wsURL, headers)
	if err != nil {
		return formatDialError(ctx, m.wsURL, response, err)
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		_ = conn.Close()
		return fmt.Errorf("websocket manager closed")
	}
	if m.writeStop != nil {
		safeClose(m.writeStop)
	}
	if m.heartbeatStop != nil {
		safeClose(m.heartbeatStop)
	}
	m.conn = conn
	m.writeQueue = make(chan []byte, m.writeQueueCapacity)
	m.writeStop = make(chan struct{})
	m.heartbeatStop = make(chan struct{})
	m.lastServerResponse = time.Now()
	writeQueue := m.writeQueue
	writeStop := m.writeStop
	heartbeatStop := m.heartbeatStop
	connCount := m.connectionCount
	m.connectionCount++
	if request.Reason != "" {
		m.lastCloseReason = request.Reason
	}
	lastCloseReason := m.lastCloseReason
	m.mu.Unlock()

	sessionID := protocol.MustSessionID(uuid.NewString())
	maxObserved := ""
	if request.MaxObservedTimestamp > 0 {
		maxObserved = protocol.EncodeTimestamp(request.MaxObservedTimestamp)
	}
	clientTS := int64(0)
	connect := protocol.ClientMessage{
		Type:                 "Connect",
		SessionID:            sessionID,
		ConnectionCount:      connCount,
		LastCloseReason:      lastCloseReason,
		MaxObservedTimestamp: maxObserved,
		ClientTS:             &clientTS,
	}
	payload, err := protocol.EncodeClientMessage(connect)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		_ = conn.Close()
		return err
	}

	go m.writeLoop(conn, writeQueue, writeStop)
	go m.heartbeatLoop(conn, heartbeatStop)
	go m.readLoop(conn)
	return nil
}

func (m *WebSocketManager) closeActiveConnection() {
	m.mu.Lock()
	conn := m.conn
	stop := m.writeStop
	heartbeatStop := m.heartbeatStop
	m.conn = nil
	m.writeQueue = nil
	m.writeStop = nil
	m.heartbeatStop = nil
	m.lastServerResponse = time.Time{}
	m.mu.Unlock()
	safeClose(stop)
	safeClose(heartbeatStop)
	if conn != nil {
		_ = conn.Close()
	}
}

func formatDialError(ctx context.Context, wsURL string, response *http.Response, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("websocket dial cancelled: %w", ctxErr)
	}
	if response == nil {
		return fmt.Errorf("websocket dial to %s failed: %w", wsURL, err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
	if len(body) == 0 {
		return fmt.Errorf("websocket handshake to %s failed with %s: %w", wsURL, response.Status, err)
	}
	return fmt.Errorf("websocket handshake to %s failed with %s: %w: %s", wsURL, response.Status, err, string(body))
}

func (m *WebSocketManager) writeLoop(conn *websocket.Conn, queue <-chan []byte, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case payload := <-queue:
			if err := conn.SetWriteDeadline(time.Now().Add(15 * time.Second)); err != nil {
				m.emitResponse(ProtocolResponse{Err: err})
				_ = conn.Close()
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				m.emitResponse(ProtocolResponse{Err: err})
				_ = conn.Close()
				return
			}
		}
	}
}

func (m *WebSocketManager) heartbeatLoop(conn *websocket.Conn, stop <-chan struct{}) {
	ticker := time.NewTicker(m.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if m.serverInactive(conn) {
				m.emitResponse(ProtocolResponse{Err: fmt.Errorf("InactiveServer")})
				_ = conn.Close()
				return
			}
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				m.emitResponse(ProtocolResponse{Err: fmt.Errorf("heartbeat ping failed: %w", err)})
				_ = conn.Close()
				return
			}
		}
	}
}

func (m *WebSocketManager) serverInactive(conn *websocket.Conn) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn != conn {
		return false
	}
	if m.lastServerResponse.IsZero() {
		return false
	}
	return time.Since(m.lastServerResponse) > m.inactivityThreshold
}

func (m *WebSocketManager) markServerResponse(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn == conn {
		m.lastServerResponse = time.Now()
	}
}

func (m *WebSocketManager) readLoop(conn *websocket.Conn) {
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			var stop chan struct{}
			m.mu.Lock()
			if m.conn == conn {
				m.conn = nil
				m.writeQueue = nil
				stop = m.writeStop
				m.writeStop = nil
				safeClose(m.heartbeatStop)
				m.heartbeatStop = nil
				m.lastServerResponse = time.Time{}
			}
			m.mu.Unlock()
			safeClose(stop)
			m.emitResponse(ProtocolResponse{Err: classifyReadError(err)})
			return
		}
		m.markServerResponse(conn)

		if messageType != websocket.TextMessage {
			m.emitResponse(ProtocolResponse{Err: fmt.Errorf("unsupported websocket frame type %d", messageType)})
			continue
		}

		message, err := protocol.DecodeServerMessage(payload)
		if err != nil {
			m.emitResponse(ProtocolResponse{Err: fmt.Errorf("protocol decode failure: %w", err)})
			continue
		}

		m.emitResponse(ProtocolResponse{Message: &message})
	}
}

func classifyReadError(err error) error {
	if closeErr, ok := err.(*websocket.CloseError); ok {
		return fmt.Errorf("websocket close frame: code=%d reason=%s", closeErr.Code, closeErr.Text)
	}
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return fmt.Errorf("websocket closed: %w", err)
	}
	return fmt.Errorf("websocket read failed: %w", err)
}

func (m *WebSocketManager) emitResponse(response ProtocolResponse) {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	responses := m.responses
	m.mu.Unlock()

	defer func() {
		_ = recover()
	}()

	select {
	case responses <- response:
	default:
	}
}

func safeClose(ch chan struct{}) {
	if ch == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	close(ch)
}
