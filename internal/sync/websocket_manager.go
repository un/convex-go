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
	mu              sync.Mutex
	wsURL           string
	clientID        string
	conn            *websocket.Conn
	responses       chan ProtocolResponse
	closed          bool
	connectionCount uint32
	lastCloseReason string
}

func NewWebSocketManager(wsURL, clientID string) *WebSocketManager {
	return &WebSocketManager{
		wsURL:           wsURL,
		clientID:        clientID,
		responses:       make(chan ProtocolResponse, 256),
		lastCloseReason: "InitialConnect",
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
	conn := m.conn
	closed := m.closed
	m.mu.Unlock()

	if closed {
		return fmt.Errorf("websocket manager closed")
	}
	if conn == nil {
		return fmt.Errorf("websocket not connected")
	}

	deadline, hasDeadline := ctx.Deadline()
	if hasDeadline {
		_ = conn.SetWriteDeadline(deadline)
	} else {
		_ = conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	}
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		m.responses <- ProtocolResponse{Err: err}
		return err
	}
	return nil
}

func (m *WebSocketManager) Reconnect(ctx context.Context, request ReconnectRequest) error {
	m.mu.Lock()
	conn := m.conn
	m.conn = nil
	m.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}
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
	m.conn = nil
	close(m.responses)
	m.mu.Unlock()

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
	m.conn = conn
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

	go m.readLoop(conn)
	return nil
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

func (m *WebSocketManager) readLoop(conn *websocket.Conn) {
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			m.mu.Lock()
			if m.conn == conn {
				m.conn = nil
			}
			m.mu.Unlock()
			m.emitResponse(ProtocolResponse{Err: err})
			return
		}

		message, err := protocol.DecodeServerMessage(payload)
		if err != nil {
			m.emitResponse(ProtocolResponse{Err: err})
			continue
		}

		m.emitResponse(ProtocolResponse{Message: &message})
	}
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
