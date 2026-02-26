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
	var msg ClientMessage
	err := json.Unmarshal(data, &msg)
	return msg, err
}

func EncodeServerMessage(msg ServerMessage) ([]byte, error) {
	return json.Marshal(msg)
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
