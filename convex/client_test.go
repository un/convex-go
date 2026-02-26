package convex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/get-convex/convex-go/internal/protocol"
)

func TestQueryMatchesSubscribeFirstValue(t *testing.T) {
	server := newSyncTestServer(t)
	defer server.Close()

	client := NewClientBuilder().WithDeploymentURL(server.URL).Build()
	defer client.Close()

	result, err := client.Query(context.Background(), "test:query", map[string]any{"x": int64(1)})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	value, err := result.Unwrap()
	if err != nil {
		t.Fatalf("unwrap failed: %v", err)
	}

	raw, ok := value.Raw().(map[string]any)
	if !ok {
		t.Fatalf("expected object result, got %T", value.Raw())
	}
	if raw["source"] != "server" {
		t.Fatalf("expected server-sourced value, got %#v", raw)
	}
}

func TestWatchAllSnapshot(t *testing.T) {
	server := newSyncTestServer(t)
	defer server.Close()

	client := NewClientBuilder().WithDeploymentURL(server.URL).Build()
	defer client.Close()

	sub, err := client.Subscribe(context.Background(), "test:query", map[string]any{})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	defer sub.Close()

	watch := client.WatchAll()
	defer watch.Close()

	select {
	case snapshot := <-watch.Updates():
		if len(snapshot) == 0 {
			t.Fatalf("expected non-empty snapshot")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for watch_all snapshot")
	}
}

func TestCloneSharesInstance(t *testing.T) {
	c := NewClient()
	clone := c.Clone()
	if c != clone {
		t.Fatalf("expected clone to share connection instance")
	}
}

func TestWebSocketStateCallbackOrderingInitialConnect(t *testing.T) {
	server := newSyncTestServer(t)
	defer server.Close()

	states := make(chan WebSocketState, 8)
	client := NewClientBuilder().
		WithDeploymentURL(server.URL).
		WithWebSocketStateCallback(func(state WebSocketState) {
			states <- state
		}).
		Build()
	defer client.Close()

	_, err := client.Query(context.Background(), "test:query", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	first := awaitState(t, states)
	second := awaitState(t, states)
	if first != WebSocketStateConnecting || second != WebSocketStateConnected {
		t.Fatalf("unexpected initial state ordering: %s -> %s", first, second)
	}
}

func TestWebSocketStateCallbackOrderingReconnect(t *testing.T) {
	server := newReconnectingSyncTestServer(t)
	defer server.Close()

	var (
		mu     sync.Mutex
		states []WebSocketState
	)
	client := NewClientBuilder().
		WithDeploymentURL(server.URL).
		WithWebSocketStateCallback(func(state WebSocketState) {
			mu.Lock()
			states = append(states, state)
			mu.Unlock()
		}).
		Build()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Query(ctx, "test:query", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("query with reconnect failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(states) < 4 {
		t.Fatalf("expected reconnect state sequence, got %v", states)
	}
	expectedPrefix := []WebSocketState{WebSocketStateConnecting, WebSocketStateConnected, WebSocketStateReconnecting, WebSocketStateConnected}
	for i, expected := range expectedPrefix {
		if states[i] != expected {
			t.Fatalf("unexpected state at %d: got %s want %s (full: %v)", i, states[i], expected, states)
		}
	}
}

func awaitState(t *testing.T, states <-chan WebSocketState) WebSocketState {
	t.Helper()
	select {
	case state := <-states:
		return state
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for websocket state callback")
		return ""
	}
}

func newReconnectingSyncTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	upgrader := websocket.Upgrader{}
	var mu sync.Mutex
	connections := 0

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sync" {
			http.NotFound(w, r)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		mu.Lock()
		connections++
		connectionIndex := connections
		mu.Unlock()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			message, err := protocol.DecodeClientMessage(data)
			if err != nil {
				continue
			}

			if message.Type == "Connect" && connectionIndex == 1 {
				_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "force reconnect"), time.Now().Add(time.Second))
				return
			}

			if message.Type != "ModifyQuerySet" {
				continue
			}
			for _, mod := range message.Modifications {
				query, ok := mod.Query()
				if !ok {
					continue
				}
				payload, err := json.Marshal(NewValue(map[string]any{"source": "server", "query": query.UDFPath}))
				if err != nil {
					t.Fatalf("marshal failed: %v", err)
				}
				transition := protocol.ServerMessage{
					Type: "Transition",
					StartVersion: &protocol.StateVersion{
						QuerySet: 0,
						Identity: 0,
						TS:       protocol.NewTimestamp(0),
					},
					EndVersion: &protocol.StateVersion{
						QuerySet: 1,
						Identity: 0,
						TS:       protocol.NewTimestamp(1),
					},
					Modifications: []protocol.StateModification{
						protocol.NewStateModificationQueryUpdated(query.QueryID, payload, nil),
					},
				}
				bytes, err := protocol.EncodeServerMessage(transition)
				if err != nil {
					t.Fatalf("encode failed: %v", err)
				}
				if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
					return
				}
			}
		}
	}))
}

func newSyncTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	upgrader := websocket.Upgrader{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/sync" {
			http.NotFound(w, r)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			message, err := protocol.DecodeClientMessage(data)
			if err != nil {
				continue
			}

			switch message.Type {
			case "Connect":
			case "Authenticate":
			case "ModifyQuerySet":
				for _, mod := range message.Modifications {
					query, ok := mod.Query()
					if !ok {
						continue
					}

					payload, err := json.Marshal(NewValue(map[string]any{
						"source": "server",
						"query":  query.UDFPath,
					}))
					if err != nil {
						t.Fatalf("marshal failed: %v", err)
					}
					transition := protocol.ServerMessage{
						Type: "Transition",
						StartVersion: &protocol.StateVersion{
							QuerySet: 0,
							Identity: 0,
							TS:       protocol.NewTimestamp(0),
						},
						EndVersion: &protocol.StateVersion{
							QuerySet: 1,
							Identity: 0,
							TS:       protocol.NewTimestamp(1),
						},
						Modifications: []protocol.StateModification{
							protocol.NewStateModificationQueryUpdated(query.QueryID, payload, nil),
						},
					}
					bytes, err := protocol.EncodeServerMessage(transition)
					if err != nil {
						t.Fatalf("encode failed: %v", err)
					}
					if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
						return
					}
				}
			case "Mutation":
				success := true
				result, err := json.Marshal(NewValue(map[string]any{"ok": true}))
				if err != nil {
					t.Fatalf("marshal failed: %v", err)
				}

				response := protocol.ServerMessage{
					Type:      "MutationResponse",
					RequestID: message.RequestID,
					Success:   &success,
					Result:    result,
					TS:        protocol.EncodeTimestamp(2),
				}
				respBytes, err := protocol.EncodeServerMessage(response)
				if err != nil {
					t.Fatalf("encode failed: %v", err)
				}
				if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
					return
				}

				transition := protocol.ServerMessage{
					Type:         "Transition",
					StartVersion: &protocol.StateVersion{TS: protocol.NewTimestamp(1)},
					EndVersion:   &protocol.StateVersion{TS: protocol.NewTimestamp(2)},
				}
				transitionBytes, err := protocol.EncodeServerMessage(transition)
				if err != nil {
					t.Fatalf("encode failed: %v", err)
				}
				if err := conn.WriteMessage(websocket.TextMessage, transitionBytes); err != nil {
					return
				}
			case "Action":
				success := true
				result, err := json.Marshal(NewValue(map[string]any{"ok": true}))
				if err != nil {
					t.Fatalf("marshal failed: %v", err)
				}

				response := protocol.ServerMessage{
					Type:      "ActionResponse",
					RequestID: message.RequestID,
					Success:   &success,
					Result:    result,
				}
				respBytes, err := protocol.EncodeServerMessage(response)
				if err != nil {
					t.Fatalf("encode failed: %v", err)
				}
				if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
					return
				}
			}
		}
	}))
}
