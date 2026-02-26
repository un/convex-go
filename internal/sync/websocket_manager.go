package sync

import (
    "context"
    "sync"

    "github.com/get-convex/convex-go/internal/protocol"
)

type WebSocketManager struct {
    mu        sync.Mutex
    responses chan ProtocolResponse
    sent      []protocol.ClientMessage
}

func NewWebSocketManager() *WebSocketManager {
    return &WebSocketManager{responses: make(chan ProtocolResponse, 64)}
}

func (m *WebSocketManager) Open(ctx context.Context, request ReconnectRequest) (<-chan ProtocolResponse, error) {
    _ = ctx
    _ = request
    return m.responses, nil
}

func (m *WebSocketManager) Send(ctx context.Context, message protocol.ClientMessage) error {
    _ = ctx
    m.mu.Lock()
    defer m.mu.Unlock()
    m.sent = append(m.sent, message)
    return nil
}

func (m *WebSocketManager) Reconnect(ctx context.Context, request ReconnectRequest) error {
    _ = ctx
    _ = request
    return nil
}

func (m *WebSocketManager) Close() error {
    return nil
}

func (m *WebSocketManager) Inject(message protocol.ServerMessage) {
    m.responses <- ProtocolResponse{Message: &message}
}
