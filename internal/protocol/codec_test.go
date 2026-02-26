package protocol

import "testing"

func TestClientMessageCodecRoundTrip(t *testing.T) {
    in := ClientMessage{Type: "Connect", RequestID: 1, Payload: map[string]any{"a": "b"}}
    raw, err := EncodeClientMessage(in)
    if err != nil {
        t.Fatalf("encode failed: %v", err)
    }
    out, err := DecodeClientMessage(raw)
    if err != nil {
        t.Fatalf("decode failed: %v", err)
    }
    if out.Type != in.Type {
        t.Fatalf("type mismatch: %s != %s", out.Type, in.Type)
    }
}
