package baseclient

import (
	"math"
	"math/rand"
	"time"
)

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
	return rand.Float64()
}

func NewBackoff(min, max time.Duration, rng RNG) *Backoff {
	if rng == nil {
		rng = defaultRNG{}
	}
	return &Backoff{min: min, max: max, rng: rng}
}

func (b *Backoff) Next() time.Duration {
	base := b.min
	if b.failures > 0 {
		if b.failures >= 30 {
			base = b.max
		} else {
			multiplier := 1 << b.failures
			if int64(base) > math.MaxInt64/int64(multiplier) {
				base = b.max
			} else {
				base = time.Duration(int64(base) * int64(multiplier))
				if base > b.max {
					base = b.max
				}
			}
		}
	}
	jitter := b.rng.Float64()
	if jitter < 0 {
		jitter = 0
	}
	if jitter > 1 {
		jitter = 1
	}
	b.failures++
	return time.Duration(float64(base) * jitter)
}

func (b *Backoff) Reset() {
	b.failures = 0
}

func (b *Backoff) Failures() int {
	return b.failures
}

func (b *Backoff) SetFailures(failures int) {
	if failures < 0 {
		failures = 0
	}
	b.failures = failures
}
