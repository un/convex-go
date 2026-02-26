package protocol

import (
	"encoding/json"
	"os"
	"testing"
)

func FuzzDecodeClientMessage(f *testing.F) {
	for _, seed := range loadClientFuzzSeeds(f) {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, payload []byte) {
		if len(payload) > 1<<20 {
			t.Skip()
		}
		msg, err := DecodeClientMessage(payload)
		if err != nil {
			return
		}
		_, _ = EncodeClientMessage(msg)
	})
}

func FuzzDecodeServerMessage(f *testing.F) {
	for _, seed := range loadServerFuzzSeeds(f) {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, payload []byte) {
		if len(payload) > 1<<20 {
			t.Skip()
		}
		msg, err := DecodeServerMessage(payload)
		if err != nil {
			return
		}
		_, _ = EncodeServerMessage(msg)
	})
}

func loadClientFuzzSeeds(t *testing.F) [][]byte {
	seeds := [][]byte{}

	validClientMessages := []ClientMessage{
		{Type: "Connect", SessionID: MustSessionID("f47ac10b-58cc-4372-a567-0e02b2c3d479"), ConnectionCount: 0},
		{Type: "Authenticate", BaseVersion: 0, Token: NewNoAuthenticationToken()},
		{Type: "Event", EventType: "ping", Event: json.RawMessage(`{"ok":true}`)},
	}
	for _, message := range validClientMessages {
		encoded, err := EncodeClientMessage(message)
		if err != nil {
			t.Fatalf("seed encode failed: %v", err)
		}
		seeds = append(seeds, encoded)
	}

	seeds = append(seeds, loadMalformedCorpusByKind(t, "client")...)
	seeds = append(seeds, loadRustClientFixtureSeeds(t)...)
	return seeds
}

func loadServerFuzzSeeds(t *testing.F) [][]byte {
	seeds := [][]byte{}

	success := true
	validServerMessages := []ServerMessage{
		{Type: "Ping"},
		{Type: "ActionResponse", RequestID: NewRequestSequenceNumber(1), Success: &success, Result: json.RawMessage(`{"ok":true}`)},
		{
			Type:         "Transition",
			StartVersion: &StateVersion{QuerySet: NewQuerySetVersion(0), Identity: NewIdentityVersion(0), TS: NewTimestamp(1)},
			EndVersion:   &StateVersion{QuerySet: NewQuerySetVersion(1), Identity: NewIdentityVersion(0), TS: NewTimestamp(2)},
			Modifications: []StateModification{
				NewStateModificationQueryRemoved(NewQueryID(1)),
			},
		},
	}
	for _, message := range validServerMessages {
		encoded, err := EncodeServerMessage(message)
		if err != nil {
			t.Fatalf("seed encode failed: %v", err)
		}
		seeds = append(seeds, encoded)
	}

	seeds = append(seeds, loadMalformedCorpusByKind(t, "server")...)
	return seeds
}

func loadMalformedCorpusByKind(t *testing.F, kind string) [][]byte {
	t.Helper()
	type malformedCase struct {
		Kind    string `json:"kind"`
		Payload string `json:"payload"`
	}
	raw, err := os.ReadFile("testdata/malformed_corpus.json")
	if err != nil {
		t.Fatalf("read malformed corpus failed: %v", err)
	}
	var cases []malformedCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("decode malformed corpus failed: %v", err)
	}
	var seeds [][]byte
	for _, tc := range cases {
		if tc.Kind == kind {
			seeds = append(seeds, []byte(tc.Payload))
		}
	}
	return seeds
}

func loadRustClientFixtureSeeds(t *testing.F) [][]byte {
	t.Helper()
	type fixtureVector struct {
		Kind    string          `json:"kind"`
		Payload json.RawMessage `json:"payload"`
	}
	type fixtureFile struct {
		Vectors []fixtureVector `json:"vectors"`
	}
	raw, err := os.ReadFile("testdata/rust_fixture_vectors.json")
	if err != nil {
		t.Fatalf("read rust fixtures failed: %v", err)
	}
	var fixtures fixtureFile
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		t.Fatalf("decode rust fixtures failed: %v", err)
	}
	var seeds [][]byte
	for _, vector := range fixtures.Vectors {
		if vector.Kind == "client_decode" {
			seeds = append(seeds, vector.Payload)
		}
	}
	return seeds
}
