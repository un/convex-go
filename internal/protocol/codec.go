package protocol

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func EncodeClientMessage(msg ClientMessage) ([]byte, error) {
	return json.Marshal(msg)
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
