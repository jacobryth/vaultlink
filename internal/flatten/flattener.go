package flatten

import (
	"fmt"
	"strings"
)

// Level controls how nested keys are flattened.
type Level string

const (
	LevelNone      Level = "none"
	LevelUnderscore Level = "underscore"
	LevelDot        Level = "dot"
)

// Flattener converts nested key segments into flat env-style keys.
type Flattener struct {
	level     Level
	separator string
}

// New returns a Flattener for the given level.
func New(level Level) (*Flattener, error) {
	switch level {
	case LevelNone:
		return &Flattener{level: level, separator: ""}, nil
	case LevelUnderscore:
		return &Flattener{level: level, separator: "_"}, nil
	case LevelDot:
		return &Flattener{level: level, separator: "."}, nil
	default:
		return nil, fmt.Errorf("flatten: unknown level %q", level)
	}
}

// Apply flattens keys in the provided secrets map.
// Keys containing "/" are joined using the configured separator.
// If level is none, the map is returned unchanged.
func (f *Flattener) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if f.level == LevelNone {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		parts := strings.Split(k, "/")
		flat := strings.Join(parts, f.separator)
		out[flat] = v
	}
	return out
}
