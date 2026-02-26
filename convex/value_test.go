package convex

import (
    "encoding/json"
    "math"
    "testing"
)

func TestValueIntegerRoundTrip(t *testing.T) {
    in := NewValue(int64(42))
    raw, err := json.Marshal(in)
    if err != nil {
        t.Fatalf("marshal failed: %v", err)
    }
    var out Value
    if err := json.Unmarshal(raw, &out); err != nil {
        t.Fatalf("unmarshal failed: %v", err)
    }
    if got := out.Raw().(int64); got != 42 {
        t.Fatalf("expected 42, got %d", got)
    }
}

func TestValueSpecialFloatEncoding(t *testing.T) {
    in := NewValue(math.Inf(1))
    raw, err := json.Marshal(in)
    if err != nil {
        t.Fatalf("marshal failed: %v", err)
    }
    if string(raw) != `{"$float":"Infinity"}` {
        t.Fatalf("unexpected float encoding: %s", string(raw))
    }
}

func TestUnsupportedSetMap(t *testing.T) {
    var out Value
    if err := json.Unmarshal([]byte(`{"$set":[]}`), &out); err == nil {
        t.Fatalf("expected $set error")
    }
    if err := json.Unmarshal([]byte(`{"$map":[]}`), &out); err == nil {
        t.Fatalf("expected $map error")
    }
}
