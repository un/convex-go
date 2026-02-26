package testutil

import "sync"

type DeterministicRNG struct {
    mu     sync.Mutex
    values []float64
    index  int
}

func NewDeterministicRNG(values ...float64) *DeterministicRNG {
    return &DeterministicRNG{values: values}
}

func (r *DeterministicRNG) Float64() float64 {
    r.mu.Lock()
    defer r.mu.Unlock()
    if len(r.values) == 0 {
        return 0.5
    }
    value := r.values[r.index%len(r.values)]
    r.index++
    return value
}
