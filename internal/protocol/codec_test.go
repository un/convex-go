package protocol

import (
	"encoding/json"
	"testing"
)

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

func TestStateVersionTimestampWireShape(t *testing.T) {
	version := StateVersion{
		QuerySet: NewQuerySetVersion(3),
		Identity: NewIdentityVersion(5),
		TS:       NewTimestamp(11),
	}
	raw, err := json.Marshal(version)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var wire map[string]any
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatalf("unmarshal wire map failed: %v", err)
	}
	if wire["ts"] != EncodeTimestamp(11) {
		t.Fatalf("unexpected ts wire value: %v", wire["ts"])
	}

	var decoded StateVersion
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("state version decode failed: %v", err)
	}
	if decoded.QuerySet != version.QuerySet || decoded.Identity != version.Identity || decoded.TS != version.TS {
		t.Fatalf("decoded state version mismatch: got %+v want %+v", decoded, version)
	}
}

func TestStateVersionRejectsInvalidTimestamp(t *testing.T) {
	raw := []byte(`{"querySet":1,"identity":2,"ts":"bad"}`)
	var version StateVersion
	err := json.Unmarshal(raw, &version)
	if err == nil {
		t.Fatalf("expected invalid timestamp error")
	}
}

func TestQuerySetModificationAddRoundTrip(t *testing.T) {
	input := NewQuerySetAdd(Query{
		QueryID: NewQueryID(9),
		UDFPath: "messages:list",
		Args:    json.RawMessage(`[{"limit":10}]`),
	})

	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var output QuerySetModification
	if err := json.Unmarshal(raw, &output); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	query, ok := output.Query()
	if !ok {
		t.Fatalf("expected add variant after roundtrip")
	}
	if query.QueryID != NewQueryID(9) || query.UDFPath != "messages:list" {
		t.Fatalf("unexpected query payload: %+v", query)
	}
}

func TestQuerySetModificationRemoveRoundTrip(t *testing.T) {
	input := NewQuerySetRemove(NewQueryID(3))
	raw, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var output QuerySetModification
	if err := json.Unmarshal(raw, &output); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !output.IsRemove() {
		t.Fatalf("expected remove variant after roundtrip")
	}
	if output.QueryID() != NewQueryID(3) {
		t.Fatalf("unexpected remove query id: %d", output.QueryID())
	}
}

func TestQuerySetModificationRejectsMalformedAdd(t *testing.T) {
	raw := []byte(`{"type":"Add","queryId":1}`)
	var output QuerySetModification
	if err := json.Unmarshal(raw, &output); err == nil {
		t.Fatalf("expected malformed add decode error")
	}
}
