package baseclient

import (
	"testing"
	"time"

	"github.com/get-convex/convex-go/internal/testutil"
)

func TestBackoffGrowthAndReset(t *testing.T) {
	rng := testutil.NewDeterministicRNG(1.0, 1.0, 1.0)
	b := NewBackoff(100*time.Millisecond, 2*time.Second, rng)
	first := b.Next()
	second := b.Next()
	third := b.Next()
	if first != 100*time.Millisecond || second != 200*time.Millisecond || third != 400*time.Millisecond {
		t.Fatalf("unexpected deterministic backoff sequence: %s %s %s", first, second, third)
	}
	b.Reset()
	if b.Failures() != 0 {
		t.Fatalf("expected reset failures")
	}
}

func TestBackoffCapAndJitterParity(t *testing.T) {
	rng := testutil.NewDeterministicRNG(0.25, 0.75)
	b := NewBackoff(100*time.Millisecond, 250*time.Millisecond, rng)
	b.SetFailures(4)

	first := b.Next()
	if first != 62500*time.Microsecond {
		t.Fatalf("expected capped jittered delay of 62.5ms, got %s", first)
	}
	second := b.Next()
	if second != 187500*time.Microsecond {
		t.Fatalf("expected capped jittered delay of 187.5ms, got %s", second)
	}
}
