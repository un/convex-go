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

func TestVersionAndSubscriberLifecycleSemantics(t *testing.T) {
	s := NewLocalSyncState()
	if s.QuerySetVersion() != 0 || s.IdentityVersion() != 0 {
		t.Fatalf("expected initial versions to be zero")
	}

	queryID, subA, added, err := s.Subscribe("f", map[string]any{"x": 1})
	if err != nil || !added {
		t.Fatalf("expected first subscribe add, got added=%v err=%v", added, err)
	}
	if s.QuerySetVersion() != 1 {
		t.Fatalf("expected query set version to increment on first add")
	}

	_, subB, added, err := s.Subscribe("f", map[string]any{"x": 1})
	if err != nil || added {
		t.Fatalf("expected deduped subscribe, got added=%v err=%v", added, err)
	}
	if s.QuerySetVersion() != 1 {
		t.Fatalf("expected query set version to stay stable on dedupe")
	}

	s.SetQueryValue(queryID, "value")
	results := s.ResultsBySubscriber()
	if results[subA] != "value" || results[subB] != "value" {
		t.Fatalf("expected both subscribers to observe same query value")
	}

	if _, removed, err := s.Unsubscribe(subA); err != nil || removed {
		t.Fatalf("expected first unsubscribe not to remove query")
	}
	if s.QuerySetVersion() != 1 {
		t.Fatalf("expected query set version unchanged after partial unsubscribe")
	}

	if _, removed, err := s.Unsubscribe(subB); err != nil || !removed {
		t.Fatalf("expected final unsubscribe to remove query")
	}
	if s.QuerySetVersion() != 2 {
		t.Fatalf("expected query set version increment on query removal")
	}
}

func TestObservedTimestampMonotonicity(t *testing.T) {
	s := NewLocalSyncState()
	s.UpdateObservedTimestamp(5)
	s.UpdateObservedTimestamp(3)
	s.UpdateObservedTimestamp(9)
	if s.ObservedTimestamp() != 9 {
		t.Fatalf("expected monotonic observed timestamp, got %d", s.ObservedTimestamp())
	}
}

func TestIdentityVersionIncrementsAcrossAuthUpdates(t *testing.T) {
	s := NewLocalSyncState()
	if s.IdentityVersion() != 0 {
		t.Fatalf("expected initial identity version zero")
	}

	token := "token-1"
	s.SetAuthToken(&token)
	if s.IdentityVersion() != 1 {
		t.Fatalf("expected identity version increment after set auth token")
	}

	if err := s.SetAuthCallback(func(forceRefresh bool) (*string, error) {
		t := "token-2"
		return &t, nil
	}); err != nil {
		t.Fatalf("set auth callback failed: %v", err)
	}
	if s.IdentityVersion() != 2 {
		t.Fatalf("expected identity version increment after auth callback registration")
	}

	if err := s.RefreshAuthOnReconnect(); err != nil {
		t.Fatalf("refresh auth on reconnect failed: %v", err)
	}
	if s.IdentityVersion() != 3 {
		t.Fatalf("expected identity version increment after reconnect auth refresh")
	}
}
