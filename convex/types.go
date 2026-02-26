package convex

import "fmt"

type ConvexError struct {
    Message string         `json:"message"`
    Data    map[string]any `json:"data,omitempty"`
}

func (e ConvexError) Error() string {
    if e.Message == "" {
        return "convex error"
    }
    return e.Message
}

type FunctionResult struct {
    Value *Value
    Err   error
}

func Success(v Value) FunctionResult {
    vv := v
    return FunctionResult{Value: &vv}
}

func Failure(err error) FunctionResult {
    return FunctionResult{Err: err}
}

func (r FunctionResult) Unwrap() (Value, error) {
    if r.Err != nil {
        return Value{}, r.Err
    }
    if r.Value == nil {
        return Value{}, fmt.Errorf("missing value")
    }
    return *r.Value, nil
}

type AuthenticationToken struct {
    Token    *string           `json:"token,omitempty"`
    Admin    bool              `json:"admin,omitempty"`
    ActingAs map[string]string `json:"actingAs,omitempty"`
}

type WebSocketState string

const (
    WebSocketStateDisconnected WebSocketState = "disconnected"
    WebSocketStateConnecting   WebSocketState = "connecting"
    WebSocketStateConnected    WebSocketState = "connected"
    WebSocketStateReconnecting WebSocketState = "reconnecting"
)
