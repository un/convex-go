package protocol

import (
	"fmt"
	"math"

	"github.com/google/uuid"
)

func NewQueryID(value uint32) QueryID {
	return QueryID(value)
}

func QueryIDFromUint64(value uint64) (QueryID, error) {
	if value > math.MaxUint32 {
		return 0, fmt.Errorf("query id %d exceeds uint32", value)
	}
	return QueryID(value), nil
}

func (id QueryID) Uint32() uint32 {
	return uint32(id)
}

func (id QueryID) Uint64() uint64 {
	return uint64(id)
}

func NewIdentityVersion(value uint32) IdentityVersion {
	return IdentityVersion(value)
}

func IdentityVersionFromUint64(value uint64) (IdentityVersion, error) {
	if value > math.MaxUint32 {
		return 0, fmt.Errorf("identity version %d exceeds uint32", value)
	}
	return IdentityVersion(value), nil
}

func (version IdentityVersion) Uint32() uint32 {
	return uint32(version)
}

func NewQuerySetVersion(value uint32) QuerySetVersion {
	return QuerySetVersion(value)
}

func QuerySetVersionFromUint64(value uint64) (QuerySetVersion, error) {
	if value > math.MaxUint32 {
		return 0, fmt.Errorf("query set version %d exceeds uint32", value)
	}
	return QuerySetVersion(value), nil
}

func (version QuerySetVersion) Uint32() uint32 {
	return uint32(version)
}

func NewRequestSequenceNumber(value uint32) RequestSequenceNumber {
	return RequestSequenceNumber(value)
}

func RequestSequenceNumberFromUint64(value uint64) (RequestSequenceNumber, error) {
	if value > math.MaxUint32 {
		return 0, fmt.Errorf("request sequence number %d exceeds uint32", value)
	}
	return RequestSequenceNumber(value), nil
}

func (number RequestSequenceNumber) Uint32() uint32 {
	return uint32(number)
}

func NewTimestamp(value uint64) Timestamp {
	return Timestamp(value)
}

func (ts Timestamp) Uint64() uint64 {
	return uint64(ts)
}

func NewSessionID(value string) (SessionID, error) {
	if _, err := uuid.Parse(value); err != nil {
		return "", fmt.Errorf("invalid session id %q: %w", value, err)
	}
	return SessionID(value), nil
}

func MustSessionID(value string) SessionID {
	id, err := NewSessionID(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (id SessionID) String() string {
	return string(id)
}
