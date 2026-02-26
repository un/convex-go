package sync

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
