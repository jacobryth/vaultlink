package encode

import (
	"encoding/base64"
	"fmt"
)

type Level string

const (
	LevelNone   Level = "none"
	LevelBase64 Level = "base64"
)

// Encoder applies value encoding to secrets.
type Encoder struct {
	level Level
}

// New returns an Encoder for the given level.
func New(level Level) (*Encoder, error) {
	switch level {
	case LevelNone, LevelBase64:
		return &Encoder{level: level}, nil
	default:
		return nil, fmt.Errorf("encode: unknown level %q", level)
	}
}

// Apply encodes secret values according to the configured level.
func (e *Encoder) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch e.level {
		case LevelBase64:
			out[k] = base64.StdEncoding.EncodeToString([]byte(v))
		default:
			out[k] = v
		}
	}
	return out
}

// Decode reverses base64 encoding for a single value.
func Decode(encoded string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("encode: decode failed: %w", err)
	}
	return string(b), nil
}
