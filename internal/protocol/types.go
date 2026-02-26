package protocol

import "encoding/json"

type QueryID uint32
type SessionID string
type IdentityVersion uint32
type QuerySetVersion uint32
type Timestamp uint64
type RequestSequenceNumber uint32

type StateVersion struct {
	QuerySet QuerySetVersion `json:"querySet"`
	Identity IdentityVersion `json:"identity"`
	TS       string          `json:"ts"`
}

type Query struct {
	QueryID       QueryID         `json:"queryId"`
	UDFPath       string          `json:"udfPath"`
	Args          json.RawMessage `json:"args,omitempty"`
	Journal       *string         `json:"journal,omitempty"`
	ComponentPath *string         `json:"componentPath,omitempty"`
}

type QuerySetModification struct {
	Type          string          `json:"type"`
	QueryID       QueryID         `json:"queryId"`
	UDFPath       string          `json:"udfPath,omitempty"`
	Args          json.RawMessage `json:"args,omitempty"`
	Journal       *string         `json:"journal,omitempty"`
	ComponentPath *string         `json:"componentPath,omitempty"`
}

type StateModification struct {
	Type         string          `json:"type"`
	QueryID      QueryID         `json:"queryId"`
	Value        json.RawMessage `json:"value,omitempty"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	ErrorData    json.RawMessage `json:"errorData,omitempty"`
	Journal      *string         `json:"journal,omitempty"`
}

type AuthenticationToken struct {
	TokenType string          `json:"tokenType"`
	Value     string          `json:"value,omitempty"`
	ActingAs  json.RawMessage `json:"actingAs,omitempty"`
}

type ClientMessage struct {
	Type                 string                 `json:"type"`
	SessionID            SessionID              `json:"sessionId,omitempty"`
	ConnectionCount      uint32                 `json:"connectionCount,omitempty"`
	LastCloseReason      string                 `json:"lastCloseReason,omitempty"`
	MaxObservedTimestamp string                 `json:"maxObservedTimestamp,omitempty"`
	ClientTS             *int64                 `json:"clientTs,omitempty"`
	BaseVersion          uint32                 `json:"baseVersion,omitempty"`
	NewVersion           uint32                 `json:"newVersion,omitempty"`
	Modifications        []QuerySetModification `json:"modifications,omitempty"`
	RequestID            RequestSequenceNumber  `json:"requestId,omitempty"`
	UDFPath              string                 `json:"udfPath,omitempty"`
	Args                 json.RawMessage        `json:"args,omitempty"`
	ComponentPath        *string                `json:"componentPath,omitempty"`
	TokenType            string                 `json:"tokenType,omitempty"`
	Value                string                 `json:"value,omitempty"`
	ActingAs             json.RawMessage        `json:"actingAs,omitempty"`
	EventType            string                 `json:"eventType,omitempty"`
	Event                json.RawMessage        `json:"event,omitempty"`
}

type ServerMessage struct {
	Type                string                `json:"type"`
	StartVersion        *StateVersion         `json:"startVersion,omitempty"`
	EndVersion          *StateVersion         `json:"endVersion,omitempty"`
	Modifications       []StateModification   `json:"modifications,omitempty"`
	RequestID           RequestSequenceNumber `json:"requestId,omitempty"`
	Success             *bool                 `json:"success,omitempty"`
	Result              json.RawMessage       `json:"result,omitempty"`
	TS                  string                `json:"ts,omitempty"`
	ErrorData           json.RawMessage       `json:"errorData,omitempty"`
	Error               string                `json:"error,omitempty"`
	BaseVersion         *IdentityVersion      `json:"baseVersion,omitempty"`
	AuthUpdateAttempted *bool                 `json:"authUpdateAttempted,omitempty"`
	Chunk               string                `json:"chunk,omitempty"`
	PartNumber          uint32                `json:"partNumber,omitempty"`
	TotalParts          uint32                `json:"totalParts,omitempty"`
	TransitionID        string                `json:"transitionId,omitempty"`
}
