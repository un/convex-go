package baseclient

import (
    "testing"
    "time"

    "github.com/get-convex/convex-go/internal/testutil"
)

func TestBackoffGrowthAndReset(t *testing.T) {
    rng := testutil.NewDeterministicRNG(0.5, 0.5, 0.5)
    b := NewBackoff(100*time.Millisecond, 2*time.Second, rng)
    a := b.Next()
    c := b.Next()
    if c <= a {
        t.Fatalf("expected growth, got %s then %s", a, c)
    }
    b.Reset()
    if b.Failures() != 0 {
        t.Fatalf("expected reset failures")
    }
}
