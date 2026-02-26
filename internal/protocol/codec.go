package protocol

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func EncodeClientMessage(msg ClientMessage) ([]byte, error) {
	if err := validateClientMessageForEncode(msg); err != nil {
		return nil, err
	}
	return json.Marshal(msg)
}

func validateClientMessageForEncode(msg ClientMessage) error {
	switch msg.Type {
	case "Connect":
		if msg.SessionID == "" {
			return fmt.Errorf("connect message missing sessionId")
		}
		if _, err := NewSessionID(msg.SessionID.String()); err != nil {
			return err
		}
		if msg.MaxObservedTimestamp != "" {
			if _, err := DecodeTimestamp(msg.MaxObservedTimestamp); err != nil {
				return fmt.Errorf("connect message invalid maxObservedTimestamp: %w", err)
			}
		}
	case "ModifyQuerySet":
		if msg.Modifications == nil {
			return fmt.Errorf("modify query set message missing modifications")
		}
	case "Mutation", "Action":
		if msg.UDFPath == "" {
			return fmt.Errorf("%s message missing udfPath", msg.Type)
		}
		if len(msg.Args) == 0 {
			return fmt.Errorf("%s message missing args", msg.Type)
		}
	case "Authenticate":
		if msg.Token.Kind() == "" {
			return fmt.Errorf("authenticate message missing token")
		}
	case "Event":
		if msg.EventType == "" {
			return fmt.Errorf("event message missing eventType")
		}
		if len(msg.Event) == 0 {
			return fmt.Errorf("event message missing event payload")
		}
	default:
		return fmt.Errorf("unknown client message type %q", msg.Type)
	}
	return nil
}

func DecodeClientMessage(data []byte) (ClientMessage, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return ClientMessage{}, fmt.Errorf("invalid client message json: %w", err)
	}
	if envelope.Type == "" {
		return ClientMessage{}, fmt.Errorf("client message missing type")
	}

	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return ClientMessage{}, fmt.Errorf("invalid %s client message: %w", envelope.Type, err)
	}
	return msg, nil
}

func EncodeServerMessage(msg ServerMessage) ([]byte, error) {
	if err := validateServerMessageForEncode(msg); err != nil {
		return nil, err
	}
	return json.Marshal(msg)
}

func validateServerMessageForEncode(msg ServerMessage) error {
	switch msg.Type {
	case "Transition":
		if msg.StartVersion == nil {
			return fmt.Errorf("transition message missing startVersion")
		}
		if msg.EndVersion == nil {
			return fmt.Errorf("transition message missing endVersion")
		}
		if msg.Modifications == nil {
			return fmt.Errorf("transition message missing modifications")
		}
	case "TransitionChunk":
		if msg.Chunk == "" {
			return fmt.Errorf("transition chunk message missing chunk")
		}
		if msg.TransitionID == "" {
			return fmt.Errorf("transition chunk message missing transitionId")
		}
		if msg.TotalParts == 0 {
			return fmt.Errorf("transition chunk message missing totalParts")
		}
		if msg.PartNumber >= msg.TotalParts {
			return fmt.Errorf("transition chunk partNumber %d out of range for totalParts %d", msg.PartNumber, msg.TotalParts)
		}
	case "MutationResponse", "ActionResponse":
		if msg.Success == nil {
			return fmt.Errorf("%s missing success", msg.Type)
		}
		if !*msg.Success && len(msg.Result) == 0 && msg.Error == "" {
			return fmt.Errorf("%s error response missing result/error", msg.Type)
		}
		if msg.Type == "MutationResponse" && msg.TS != "" {
			if _, err := DecodeTimestamp(msg.TS); err != nil {
				return fmt.Errorf("mutation response invalid ts: %w", err)
			}
		}
	case "AuthError", "FatalError":
		if msg.Error == "" {
			return fmt.Errorf("%s message missing error", msg.Type)
		}
	case "Ping":
		return nil
	default:
		return fmt.Errorf("unknown server message type %q", msg.Type)
	}
	return nil
}

func DecodeServerMessage(data []byte) (ServerMessage, error) {
	var msg ServerMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}

func EncodeTimestamp(ts uint64) string {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, ts)
	return base64.StdEncoding.EncodeToString(bytes)
}

func DecodeTimestamp(value string) (uint64, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return 0, fmt.Errorf("timestamp must be base64: %w", err)
	}
	if len(decoded) != 8 {
		return 0, fmt.Errorf("timestamp must decode to 8 bytes, got %d", len(decoded))
	}
	return binary.LittleEndian.Uint64(decoded), nil
}
