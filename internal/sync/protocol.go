package sync

import (
    "context"

    "github.com/get-convex/convex-go/internal/protocol"
)

type ReconnectRequest struct {
    Reason               string
    MaxObservedTimestamp uint64
}

type ProtocolResponse struct {
    Message *protocol.ServerMessage
    Err     error
}

type SyncProtocol interface {
    Open(ctx context.Context, request ReconnectRequest) (<-chan ProtocolResponse, error)
    Send(ctx context.Context, message protocol.ClientMessage) error
    Reconnect(ctx context.Context, request ReconnectRequest) error
    Close() error
}
