package protocol

import (
	"encoding/json"
	"fmt"
)

type QueryID uint32
type SessionID string
type IdentityVersion uint32
type QuerySetVersion uint32
type Timestamp uint64
type RequestSequenceNumber uint32

type StateVersion struct {
	QuerySet QuerySetVersion `json:"querySet"`
	Identity IdentityVersion `json:"identity"`
	TS       Timestamp       `json:"-"`
}

type stateVersionJSON struct {
	QuerySet QuerySetVersion `json:"querySet"`
	Identity IdentityVersion `json:"identity"`
	TS       string          `json:"ts"`
}

func (version StateVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(stateVersionJSON{
		QuerySet: version.QuerySet,
		Identity: version.Identity,
		TS:       EncodeTimestamp(version.TS.Uint64()),
	})
}

func (version *StateVersion) UnmarshalJSON(data []byte) error {
	var wire stateVersionJSON
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	ts, err := DecodeTimestamp(wire.TS)
	if err != nil {
		return fmt.Errorf("invalid state version timestamp: %w", err)
	}
	version.QuerySet = wire.QuerySet
	version.Identity = wire.Identity
	version.TS = NewTimestamp(ts)
	return nil
}

type Query struct {
	QueryID       QueryID         `json:"queryId"`
	UDFPath       string          `json:"udfPath"`
	Args          json.RawMessage `json:"args,omitempty"`
	Journal       *string         `json:"journal,omitempty"`
	ComponentPath *string         `json:"componentPath,omitempty"`
}

type QuerySetModification struct {
	add    *Query
	remove *QueryID
}

func NewQuerySetAdd(query Query) QuerySetModification {
	copy := query
	return QuerySetModification{add: &copy}
}

func NewQuerySetRemove(queryID QueryID) QuerySetModification {
	copy := queryID
	return QuerySetModification{remove: &copy}
}

func (mod QuerySetModification) IsAdd() bool {
	return mod.add != nil
}

func (mod QuerySetModification) IsRemove() bool {
	return mod.remove != nil
}

func (mod QuerySetModification) QueryID() QueryID {
	if mod.add != nil {
		return mod.add.QueryID
	}
	if mod.remove != nil {
		return *mod.remove
	}
	return 0
}

func (mod QuerySetModification) Query() (Query, bool) {
	if mod.add == nil {
		return Query{}, false
	}
	return *mod.add, true
}

func (mod QuerySetModification) MarshalJSON() ([]byte, error) {
	if mod.add != nil && mod.remove != nil {
		return nil, fmt.Errorf("query set modification must be either add or remove")
	}
	if mod.add != nil {
		if mod.add.UDFPath == "" {
			return nil, fmt.Errorf("query set add modification requires udfPath")
		}
		return json.Marshal(struct {
			Type          string          `json:"type"`
			QueryID       QueryID         `json:"queryId"`
			UDFPath       string          `json:"udfPath"`
			Args          json.RawMessage `json:"args,omitempty"`
			Journal       *string         `json:"journal,omitempty"`
			ComponentPath *string         `json:"componentPath,omitempty"`
		}{
			Type:          "Add",
			QueryID:       mod.add.QueryID,
			UDFPath:       mod.add.UDFPath,
			Args:          mod.add.Args,
			Journal:       mod.add.Journal,
			ComponentPath: mod.add.ComponentPath,
		})
	}
	if mod.remove != nil {
		return json.Marshal(struct {
			Type    string  `json:"type"`
			QueryID QueryID `json:"queryId"`
		}{
			Type:    "Remove",
			QueryID: *mod.remove,
		})
	}
	return nil, fmt.Errorf("query set modification variant not set")
}

func (mod *QuerySetModification) UnmarshalJSON(data []byte) error {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	switch base.Type {
	case "Add":
		var payload struct {
			QueryID       QueryID         `json:"queryId"`
			UDFPath       string          `json:"udfPath"`
			Args          json.RawMessage `json:"args,omitempty"`
			Journal       *string         `json:"journal,omitempty"`
			ComponentPath *string         `json:"componentPath,omitempty"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		if payload.UDFPath == "" {
			return fmt.Errorf("query set add modification missing udfPath")
		}
		query := Query{
			QueryID:       payload.QueryID,
			UDFPath:       payload.UDFPath,
			Args:          payload.Args,
			Journal:       payload.Journal,
			ComponentPath: payload.ComponentPath,
		}
		mod.add = &query
		mod.remove = nil
		return nil
	case "Remove":
		var payload struct {
			QueryID *QueryID `json:"queryId"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		if payload.QueryID == nil {
			return fmt.Errorf("query set remove modification missing queryId")
		}
		copy := *payload.QueryID
		mod.remove = &copy
		mod.add = nil
		return nil
	default:
		return fmt.Errorf("unknown query set modification type %q", base.Type)
	}
}

type StateModification struct {
	updated *StateQueryUpdated
	failed  *StateQueryFailed
	removed *QueryID
}

type StateQueryUpdated struct {
	QueryID QueryID
	Value   json.RawMessage
	Journal *string
}

type StateQueryFailed struct {
	QueryID      QueryID
	ErrorMessage string
	ErrorData    json.RawMessage
	Journal      *string
}

func NewStateModificationQueryUpdated(queryID QueryID, value json.RawMessage, journal *string) StateModification {
	copy := StateQueryUpdated{QueryID: queryID, Value: value, Journal: journal}
	return StateModification{updated: &copy}
}

func NewStateModificationQueryFailed(queryID QueryID, message string, errorData json.RawMessage, journal *string) StateModification {
	copy := StateQueryFailed{QueryID: queryID, ErrorMessage: message, ErrorData: errorData, Journal: journal}
	return StateModification{failed: &copy}
}

func NewStateModificationQueryRemoved(queryID QueryID) StateModification {
	copy := queryID
	return StateModification{removed: &copy}
}

func (mod StateModification) Kind() string {
	if mod.updated != nil {
		return "QueryUpdated"
	}
	if mod.failed != nil {
		return "QueryFailed"
	}
	if mod.removed != nil {
		return "QueryRemoved"
	}
	return ""
}

func (mod StateModification) QueryUpdated() (StateQueryUpdated, bool) {
	if mod.updated == nil {
		return StateQueryUpdated{}, false
	}
	return *mod.updated, true
}

func (mod StateModification) QueryFailed() (StateQueryFailed, bool) {
	if mod.failed == nil {
		return StateQueryFailed{}, false
	}
	return *mod.failed, true
}

func (mod StateModification) QueryRemoved() (QueryID, bool) {
	if mod.removed == nil {
		return 0, false
	}
	return *mod.removed, true
}

func (mod StateModification) MarshalJSON() ([]byte, error) {
	variants := 0
	if mod.updated != nil {
		variants++
	}
	if mod.failed != nil {
		variants++
	}
	if mod.removed != nil {
		variants++
	}
	if variants != 1 {
		return nil, fmt.Errorf("state modification must contain exactly one variant")
	}

	if mod.updated != nil {
		if len(mod.updated.Value) == 0 {
			return nil, fmt.Errorf("query updated modification requires value")
		}
		return json.Marshal(struct {
			Type    string          `json:"type"`
			QueryID QueryID         `json:"queryId"`
			Value   json.RawMessage `json:"value"`
			Journal *string         `json:"journal,omitempty"`
		}{
			Type:    "QueryUpdated",
			QueryID: mod.updated.QueryID,
			Value:   mod.updated.Value,
			Journal: mod.updated.Journal,
		})
	}

	if mod.failed != nil {
		if mod.failed.ErrorMessage == "" {
			return nil, fmt.Errorf("query failed modification requires errorMessage")
		}
		return json.Marshal(struct {
			Type         string          `json:"type"`
			QueryID      QueryID         `json:"queryId"`
			ErrorMessage string          `json:"errorMessage"`
			ErrorData    json.RawMessage `json:"errorData,omitempty"`
			Journal      *string         `json:"journal,omitempty"`
		}{
			Type:         "QueryFailed",
			QueryID:      mod.failed.QueryID,
			ErrorMessage: mod.failed.ErrorMessage,
			ErrorData:    mod.failed.ErrorData,
			Journal:      mod.failed.Journal,
		})
	}

	return json.Marshal(struct {
		Type    string  `json:"type"`
		QueryID QueryID `json:"queryId"`
	}{
		Type:    "QueryRemoved",
		QueryID: *mod.removed,
	})
}

func (mod *StateModification) UnmarshalJSON(data []byte) error {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	switch base.Type {
	case "QueryUpdated":
		var payload struct {
			QueryID QueryID          `json:"queryId"`
			Value   *json.RawMessage `json:"value"`
			Journal *string          `json:"journal,omitempty"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		if payload.Value == nil {
			return fmt.Errorf("query updated modification missing value")
		}
		mod.updated = &StateQueryUpdated{QueryID: payload.QueryID, Value: *payload.Value, Journal: payload.Journal}
		mod.failed = nil
		mod.removed = nil
		return nil
	case "QueryFailed":
		var payload struct {
			QueryID      QueryID         `json:"queryId"`
			ErrorMessage *string         `json:"errorMessage"`
			ErrorData    json.RawMessage `json:"errorData,omitempty"`
			Journal      *string         `json:"journal,omitempty"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		if payload.ErrorMessage == nil {
			return fmt.Errorf("query failed modification missing errorMessage")
		}
		mod.failed = &StateQueryFailed{QueryID: payload.QueryID, ErrorMessage: *payload.ErrorMessage, ErrorData: payload.ErrorData, Journal: payload.Journal}
		mod.updated = nil
		mod.removed = nil
		return nil
	case "QueryRemoved":
		var payload struct {
			QueryID *QueryID `json:"queryId"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		if payload.QueryID == nil {
			return fmt.Errorf("query removed modification missing queryId")
		}
		copy := *payload.QueryID
		mod.removed = &copy
		mod.updated = nil
		mod.failed = nil
		return nil
	default:
		return fmt.Errorf("unknown state modification type %q", base.Type)
	}
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
