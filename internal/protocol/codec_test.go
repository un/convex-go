package protocol

import "testing"

func TestClientMessageCodecRoundTrip(t *testing.T) {
	in := ClientMessage{Type: "Connect", SessionID: SessionID("session"), ConnectionCount: 1}
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

func TestTimestampBase64RoundTrip(t *testing.T) {
	values := []uint64{0, 1, 42, 1<<53 + 7, ^uint64(0)}
	for _, value := range values {
		encoded := EncodeTimestamp(value)
		decoded, err := DecodeTimestamp(encoded)
		if err != nil {
			t.Fatalf("decode failed for %d: %v", value, err)
		}
		if decoded != value {
			t.Fatalf("roundtrip mismatch: got %d want %d", decoded, value)
		}
	}
}
