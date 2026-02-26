package convex

import "testing"

func TestTypeShapes(t *testing.T) {
    _ = ConvexError{Message: "x", Data: map[string]any{"k": "v"}}
    _ = FunctionResult{}
    _ = AuthenticationToken{}
    _ = WebSocketStateConnected
}
