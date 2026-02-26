package protocol

type QueryID uint64
type SessionID string
type IdentityVersion uint64
type QuerySetVersion uint64
type Timestamp uint64
type RequestSequenceNumber uint64

func (t Timestamp) Succ() Timestamp {
    return t + 1
}

func MinTimestamp() Timestamp {
    return 0
}

type AuthenticationToken struct {
    Token    *string           `json:"token,omitempty"`
    Admin    bool              `json:"admin,omitempty"`
    ActingAs map[string]string `json:"actingAs,omitempty"`
}

type Query struct {
    UDFPath string         `json:"udfPath"`
    Args    map[string]any `json:"args"`
}

type QuerySetModification struct {
    Type    string  `json:"type"`
    QueryID QueryID `json:"queryId"`
    Query   *Query  `json:"query,omitempty"`
}

type StateModification struct {
    Type    string `json:"type"`
    QueryID QueryID `json:"queryId"`
    Value   any    `json:"value,omitempty"`
}

type ClientMessage struct {
    Type      string         `json:"type"`
    RequestID uint64         `json:"requestId,omitempty"`
    Payload   map[string]any `json:"payload,omitempty"`
}

type ServerMessage struct {
    Type      string         `json:"type"`
    RequestID uint64         `json:"requestId,omitempty"`
    Payload   map[string]any `json:"payload,omitempty"`
}
