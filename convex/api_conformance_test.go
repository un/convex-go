package convex

import (
    "context"
    "testing"
)

func TestAPISignaturesCompile(t *testing.T) {
    c := NewClient()
    _, _ = c.Subscribe(context.Background(), "test:query", map[string]any{})
    _, _ = c.Query(context.Background(), "test:query", map[string]any{})
    _, _ = c.Mutation(context.Background(), "test:mutation", map[string]any{})
    _, _ = c.Action(context.Background(), "test:action", map[string]any{})
    _ = c.WatchAll()
    c.SetAuth(nil)
}
