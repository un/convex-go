package protocol

import "testing"

func TestQueryIDFromUint64(t *testing.T) {
	id, err := QueryIDFromUint64(42)
	if err != nil {
		t.Fatalf("expected conversion success, got %v", err)
	}
	if id.Uint64() != 42 {
		t.Fatalf("unexpected query id value: got %d", id.Uint64())
	}
}

func TestQueryIDFromUint64Overflow(t *testing.T) {
	_, err := QueryIDFromUint64(1 << 40)
	if err == nil {
		t.Fatalf("expected overflow error")
	}
}

func TestVersionConversionHelpers(t *testing.T) {
	querySet, err := QuerySetVersionFromUint64(7)
	if err != nil {
		t.Fatalf("unexpected query set version conversion error: %v", err)
	}
	if querySet.Uint32() != 7 {
		t.Fatalf("unexpected query set version: %d", querySet.Uint32())
	}

	identity, err := IdentityVersionFromUint64(11)
	if err != nil {
		t.Fatalf("unexpected identity version conversion error: %v", err)
	}
	if identity.Uint32() != 11 {
		t.Fatalf("unexpected identity version: %d", identity.Uint32())
	}

	seq, err := RequestSequenceNumberFromUint64(13)
	if err != nil {
		t.Fatalf("unexpected request sequence number conversion error: %v", err)
	}
	if seq.Uint32() != 13 {
		t.Fatalf("unexpected request sequence number: %d", seq.Uint32())
	}
}

func TestSessionIDValidation(t *testing.T) {
	session, err := NewSessionID("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	if err != nil {
		t.Fatalf("expected session id validation success, got %v", err)
	}
	if session.String() != "f47ac10b-58cc-4372-a567-0e02b2c3d479" {
		t.Fatalf("unexpected session id value: %s", session)
	}

	if _, err := NewSessionID("not-a-uuid"); err == nil {
		t.Fatalf("expected invalid session id error")
	}
}
