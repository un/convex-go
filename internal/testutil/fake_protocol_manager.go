package testutil

import (
    "context"
    "sync"

    "github.com/get-convex/convex-go/internal/protocol"
    syncproto "github.com/get-convex/convex-go/internal/sync"
)

type FakeProtocolManager struct {
    mu        sync.Mutex
    opened    bool
    requests  []syncproto.ReconnectRequest
    sent      []protocol.ClientMessage
    responses chan syncproto.ProtocolResponse
}

func NewFakeProtocolManager() *FakeProtocolManager {
    return &FakeProtocolManager{responses: make(chan syncproto.ProtocolResponse, 64)}
}

func (m *FakeProtocolManager) Open(ctx context.Context, request syncproto.ReconnectRequest) (<-chan syncproto.ProtocolResponse, error) {
    _ = ctx
    m.mu.Lock()
    m.opened = true
    m.requests = append(m.requests, request)
    m.mu.Unlock()
    return m.responses, nil
}

func (m *FakeProtocolManager) Send(ctx context.Context, message protocol.ClientMessage) error {
    _ = ctx
    m.mu.Lock()
    m.sent = append(m.sent, message)
    m.mu.Unlock()
    return nil
}

func (m *FakeProtocolManager) Reconnect(ctx context.Context, request syncproto.ReconnectRequest) error {
    _ = ctx
    m.mu.Lock()
    m.requests = append(m.requests, request)
    m.mu.Unlock()
    return nil
}

func (m *FakeProtocolManager) Close() error {
    return nil
}

func (m *FakeProtocolManager) Inject(msg protocol.ServerMessage) {
    m.responses <- syncproto.ProtocolResponse{Message: &msg}
}

func (m *FakeProtocolManager) SentMessages() []protocol.ClientMessage {
    m.mu.Lock()
    defer m.mu.Unlock()
    out := make([]protocol.ClientMessage, len(m.sent))
    copy(out, m.sent)
    return out
}
