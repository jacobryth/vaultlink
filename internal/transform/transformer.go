package transform

import "strings"

// Level defines the transformation level applied to secret values.
type Level string

const (
	LevelNone   Level = "none"
	LevelUpper  Level = "upper"
	LevelLower  Level = "lower"
	LevelTrim   Level = "trim"
)

// Transformer applies string transformations to secret values.
type Transformer struct {
	level Level
}

// New creates a Transformer for the given level.
// Defaults to LevelNone if the level is unrecognized.
func New(level Level) (*Transformer, error) {
	switch level {
	case LevelNone, LevelUpper, LevelLower, LevelTrim:
		return &Transformer{level: level}, nil
	default:
		return nil, fmt.Errorf("transform: unknown level %q", level)
	}
}

// Apply returns a new map with values transformed according to the level.
func (t *Transformer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch t.level {
		case LevelUpper:
			out[k] = strings.ToUpper(v)
		case LevelLower:
			out[k] = strings.ToLower(v)
		case LevelTrim:
			out[k] = strings.TrimSpace(v)
		default:
			out[k] = v
		}
	}
	return out
}
