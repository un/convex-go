package sync

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/get-convex/convex-go/internal/protocol"
	"github.com/gorilla/websocket"
)

func TestWebSocketOpenHandshakeFailureIncludesDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "handshake denied", http.StatusUnauthorized)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")

	_, err := manager.Open(context.Background(), ReconnectRequest{})
	if err == nil {
		t.Fatalf("expected websocket handshake failure")
	}
	errText := err.Error()
	if !strings.Contains(errText, "401") {
		t.Fatalf("expected handshake error to include status code, got: %s", errText)
	}
	if !strings.Contains(errText, "handshake denied") {
		t.Fatalf("expected handshake error to include response body, got: %s", errText)
	}
}

func TestWebSocketOpenHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	manager := NewWebSocketManager("ws://127.0.0.1:1", "test-client")
	_, err := manager.Open(ctx, ReconnectRequest{})
	if err == nil {
		t.Fatalf("expected context cancellation error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestWebSocketConnectMessageOpenAndReconnect(t *testing.T) {
	upgrader := websocket.Upgrader{}
	connectMessages := make(chan protocol.ClientMessage, 4)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			_, payload, err := conn.ReadMessage()
			if err != nil {
				return
			}
			message, err := protocol.DecodeClientMessage(payload)
			if err == nil {
				connectMessages <- message
			}
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	defer manager.Close()

	if _, err := manager.Open(context.Background(), ReconnectRequest{}); err != nil {
		t.Fatalf("open failed: %v", err)
	}
	first := awaitConnectMessage(t, connectMessages)
	if first.Type != "Connect" {
		t.Fatalf("unexpected first message type: %s", first.Type)
	}
	if first.ConnectionCount != 0 {
		t.Fatalf("expected first connectionCount=0, got %d", first.ConnectionCount)
	}
	if first.LastCloseReason != "InitialConnect" {
		t.Fatalf("expected initial close reason, got %q", first.LastCloseReason)
	}
	if first.MaxObservedTimestamp != "" {
		t.Fatalf("expected empty max observed timestamp, got %q", first.MaxObservedTimestamp)
	}
	if first.ClientTS == nil || *first.ClientTS != 0 {
		t.Fatalf("expected clientTs=0, got %v", first.ClientTS)
	}

	reconnectRequest := ReconnectRequest{Reason: "InactiveServer", MaxObservedTimestamp: 42}
	if err := manager.Reconnect(context.Background(), reconnectRequest); err != nil {
		t.Fatalf("reconnect failed: %v", err)
	}
	second := awaitConnectMessage(t, connectMessages)
	if second.ConnectionCount != 1 {
		t.Fatalf("expected reconnect connectionCount=1, got %d", second.ConnectionCount)
	}
	if second.LastCloseReason != reconnectRequest.Reason {
		t.Fatalf("expected reconnect reason %q, got %q", reconnectRequest.Reason, second.LastCloseReason)
	}
	expectedTS := protocol.EncodeTimestamp(reconnectRequest.MaxObservedTimestamp)
	if second.MaxObservedTimestamp != expectedTS {
		t.Fatalf("expected reconnect maxObservedTimestamp %q, got %q", expectedTS, second.MaxObservedTimestamp)
	}
}

func awaitConnectMessage(t *testing.T, messages <-chan protocol.ClientMessage) protocol.ClientMessage {
	t.Helper()
	select {
	case message := <-messages:
		return message
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for connect message")
		return protocol.ClientMessage{}
	}
}

func TestWebSocketWriteQueuePreservesOrdering(t *testing.T) {
	upgrader := websocket.Upgrader{}
	received := make(chan string, 8)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			for {
				_, payload, err := conn.ReadMessage()
				if err != nil {
					return
				}
				message, err := protocol.DecodeClientMessage(payload)
				if err != nil {
					continue
				}
				received <- message.EventType
			}
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	defer manager.Close()

	if _, err := manager.Open(context.Background(), ReconnectRequest{}); err != nil {
		t.Fatalf("open failed: %v", err)
	}

	messageTypes := []string{"first", "second", "third"}
	for _, eventType := range messageTypes {
		eventPayload := []byte(`{"ok":true}`)
		if err := manager.Send(context.Background(), protocol.ClientMessage{
			Type:      "Event",
			EventType: eventType,
			Event:     eventPayload,
		}); err != nil {
			t.Fatalf("send failed for %s: %v", eventType, err)
		}
	}

	for _, expected := range messageTypes {
		select {
		case actual := <-received:
			if actual != expected {
				t.Fatalf("event ordering mismatch: got %s want %s", actual, expected)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for event %s", expected)
		}
	}
}

func TestWebSocketReadLoopClassifiesDecodeFailures(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`not-json`))
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	defer manager.Close()

	responses, err := manager.Open(context.Background(), ReconnectRequest{})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	select {
	case response := <-responses:
		if response.Err == nil || !strings.Contains(response.Err.Error(), "protocol decode failure") {
			t.Fatalf("expected protocol decode failure, got %+v", response)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for decode failure response")
	}
}

func TestWebSocketReadLoopClassifiesCloseFrames(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second))
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	defer manager.Close()

	responses, err := manager.Open(context.Background(), ReconnectRequest{})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	select {
	case response := <-responses:
		if response.Err == nil || !strings.Contains(response.Err.Error(), "close frame") {
			t.Fatalf("expected close frame classification, got %+v", response)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for close frame classification")
	}
}

func TestWebSocketHeartbeatSendsPingFrames(t *testing.T) {
	upgrader := websocket.Upgrader{}
	pingSeen := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.SetPingHandler(func(appData string) error {
			select {
			case pingSeen <- struct{}{}:
			default:
			}
			return nil
		})
		go func() {
			defer conn.Close()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	manager.heartbeatInterval = 20 * time.Millisecond
	defer manager.Close()

	if _, err := manager.Open(context.Background(), ReconnectRequest{}); err != nil {
		t.Fatalf("open failed: %v", err)
	}

	select {
	case <-pingSeen:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for heartbeat ping frame")
	}
}

func TestWebSocketInactivityWatchdogTriggersFailure(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go func() {
			defer conn.Close()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	manager := NewWebSocketManager(wsURL, "test-client")
	manager.heartbeatInterval = 20 * time.Millisecond
	manager.inactivityThreshold = 40 * time.Millisecond
	defer manager.Close()

	responses, err := manager.Open(context.Background(), ReconnectRequest{})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	select {
	case response := <-responses:
		if response.Err == nil || !strings.Contains(response.Err.Error(), "InactiveServer") {
			t.Fatalf("expected inactivity failure, got %+v", response)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for inactivity failure")
	}
}
