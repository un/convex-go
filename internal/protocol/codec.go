package protocol

import "encoding/json"

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
