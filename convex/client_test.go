package convex

import (
    "context"
    "testing"
)

func TestQueryMatchesSubscribeFirstValue(t *testing.T) {
    c := NewClient()
    got, err := c.Query(context.Background(), "test:query", map[string]any{"x": 1})
    if err != nil {
        t.Fatalf("query failed: %v", err)
    }
    _, err = got.Unwrap()
    if err != nil {
        t.Fatalf("unwrap failed: %v", err)
    }
}

func TestWatchAllSnapshot(t *testing.T) {
    c := NewClient()
    sub, err := c.Subscribe(context.Background(), "test:query", map[string]any{})
    if err != nil {
        t.Fatalf("subscribe failed: %v", err)
    }
    defer sub.Close()

    watch := c.WatchAll()
    defer watch.Close()

    snapshot := <-watch.Updates()
    if len(snapshot) == 0 {
        t.Fatalf("expected non-empty snapshot")
    }
}

func TestCloneSharesInstance(t *testing.T) {
    c := NewClient()
    clone := c.Clone()
    if c != clone {
        t.Fatalf("expected clone to share connection instance")
    }
}
