package convex

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "strconv"
)

type Value struct {
    raw any
}

func NewValue(v any) Value {
    return Value{raw: v}
}

func NewNullValue() Value {
    return Value{raw: nil}
}

func (v Value) Raw() any {
    return v.raw
}

func (v Value) MarshalJSON() ([]byte, error) {
    encoded, err := encodeRaw(v.raw)
    if err != nil {
        return nil, err
    }
    return json.Marshal(encoded)
}

func (v *Value) UnmarshalJSON(data []byte) error {
    var decoded any
    if err := json.Unmarshal(data, &decoded); err != nil {
        return err
    }
    value, err := decodeRaw(decoded)
    if err != nil {
        return err
    }
    v.raw = value
    return nil
}

func encodeRaw(input any) (any, error) {
    switch t := input.(type) {
    case nil, bool, string:
        return t, nil
    case int:
        return map[string]any{"$integer": strconv.FormatInt(int64(t), 10)}, nil
    case int64:
        return map[string]any{"$integer": strconv.FormatInt(t, 10)}, nil
    case float64:
        if math.IsNaN(t) {
            return map[string]any{"$float": "NaN"}, nil
        }
        if math.IsInf(t, 1) {
            return map[string]any{"$float": "Infinity"}, nil
        }
        if math.IsInf(t, -1) {
            return map[string]any{"$float": "-Infinity"}, nil
        }
        if t == 0 && math.Signbit(t) {
            return map[string]any{"$float": "-0"}, nil
        }
        return t, nil
    case []byte:
        return map[string]any{"$bytes": base64.StdEncoding.EncodeToString(t)}, nil
    case []any:
        out := make([]any, 0, len(t))
        for _, item := range t {
            encoded, err := encodeRaw(item)
            if err != nil {
                return nil, err
            }
            out = append(out, encoded)
        }
        return out, nil
    case map[string]any:
        out := map[string]any{}
        for key, item := range t {
            encoded, err := encodeRaw(item)
            if err != nil {
                return nil, err
            }
            out[key] = encoded
        }
        return out, nil
    case Value:
        return encodeRaw(t.raw)
    default:
        return nil, fmt.Errorf("unsupported value type %T", input)
    }
}

func decodeRaw(input any) (any, error) {
    switch t := input.(type) {
    case nil, bool, string, float64:
        return t, nil
    case []any:
        out := make([]any, 0, len(t))
        for _, item := range t {
            decoded, err := decodeRaw(item)
            if err != nil {
                return nil, err
            }
            out = append(out, decoded)
        }
        return out, nil
    case map[string]any:
        if len(t) == 1 {
            if integer, ok := t["$integer"]; ok {
                s, ok := integer.(string)
                if !ok {
                    return nil, errors.New("$integer expects string")
                }
                parsed, err := strconv.ParseInt(s, 10, 64)
                if err != nil {
                    return nil, err
                }
                return parsed, nil
            }
            if f, ok := t["$float"]; ok {
                s, ok := f.(string)
                if !ok {
                    return nil, errors.New("$float expects string")
                }
                switch s {
                case "NaN":
                    return math.NaN(), nil
                case "Infinity":
                    return math.Inf(1), nil
                case "-Infinity":
                    return math.Inf(-1), nil
                case "-0":
                    return math.Copysign(0, -1), nil
                default:
                    return nil, fmt.Errorf("unsupported $float value %q", s)
                }
            }
            if b, ok := t["$bytes"]; ok {
                s, ok := b.(string)
                if !ok {
                    return nil, errors.New("$bytes expects string")
                }
                parsed, err := base64.StdEncoding.DecodeString(s)
                if err != nil {
                    return nil, err
                }
                return parsed, nil
            }
            if _, ok := t["$set"]; ok {
                return nil, errors.New("unsupported Convex type: $set")
            }
            if _, ok := t["$map"]; ok {
                return nil, errors.New("unsupported Convex type: $map")
            }
        }
        out := map[string]any{}
        for key, item := range t {
            decoded, err := decodeRaw(item)
            if err != nil {
                return nil, err
            }
            out[key] = decoded
        }
        return out, nil
    default:
        return nil, fmt.Errorf("unsupported json value type %T", input)
    }
}
