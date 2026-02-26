package baseclient

import "time"

type RNG interface {
    Float64() float64
}

type Backoff struct {
    min      time.Duration
    max      time.Duration
    failures int
    rng      RNG
}

type defaultRNG struct{}

func (defaultRNG) Float64() float64 {
    return 0.5
}

func NewBackoff(min, max time.Duration, rng RNG) *Backoff {
    if rng == nil {
        rng = defaultRNG{}
    }
    return &Backoff{min: min, max: max, rng: rng}
}

func (b *Backoff) Next() time.Duration {
    factor := 1 << minInt(b.failures, 10)
    delay := time.Duration(int64(b.min) * int64(factor))
    if delay > b.max {
        delay = b.max
    }
    jitter := 0.5 + (b.rng.Float64() * 0.5)
    b.failures++
    return time.Duration(float64(delay) * jitter)
}

func (b *Backoff) Reset() {
    b.failures = 0
}

func (b *Backoff) Failures() int {
    return b.failures
}

func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}
