package baseclient

import "testing"

func TestQueryTokenDeterminism(t *testing.T) {
    a, err := CanonicalQueryToken("f", map[string]any{"a": 1, "b": 2})
    if err != nil {
        t.Fatalf("token failed: %v", err)
    }
    b, err := CanonicalQueryToken("f", map[string]any{"b": 2, "a": 1})
    if err != nil {
        t.Fatalf("token failed: %v", err)
    }
    if a != b {
        t.Fatalf("expected deterministic token")
    }
}

func TestSubscribeUnsubscribe(t *testing.T) {
    s := NewLocalSyncState()
    q, subA, added, err := s.Subscribe("f", map[string]any{"x": 1})
    if err != nil || !added {
        t.Fatalf("subscribe failed")
    }
    _, subB, added, err := s.Subscribe("f", map[string]any{"x": 1})
    if err != nil || added {
        t.Fatalf("dedupe failed")
    }
    if _, removed, err := s.Unsubscribe(subA); err != nil || removed {
        t.Fatalf("first unsubscribe should not remove")
    }
    got, removed, err := s.Unsubscribe(subB)
    if err != nil || !removed || got != q {
        t.Fatalf("final unsubscribe should remove")
    }
}

func TestAuthCallbackBehavior(t *testing.T) {
    s := NewLocalSyncState()
    calls := []bool{}
    err := s.SetAuthCallback(func(forceRefresh bool) (*string, error) {
        calls = append(calls, forceRefresh)
        token := "t"
        return &token, nil
    })
    if err != nil {
        t.Fatalf("set callback failed: %v", err)
    }
    if len(calls) != 1 || calls[0] {
        t.Fatalf("expected immediate non-refresh call")
    }
    if err := s.RefreshAuthOnReconnect(); err != nil {
        t.Fatalf("refresh failed: %v", err)
    }
    if len(calls) != 2 || !calls[1] {
        t.Fatalf("expected reconnect refresh call")
    }
}
