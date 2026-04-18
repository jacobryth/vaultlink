package truncate

import "strings"

// Level controls how values are truncated in output.
type Level string

const (
	LevelNone  Level = "none"
	LevelShort Level = "short" // first 8 chars
	LevelTiny  Level = "tiny"  // first 3 chars
)

// Truncator applies value truncation to secrets.
type Truncator struct {
	level  Level
	suffix string
}

// New returns a Truncator for the given level.
// Unknown levels default to LevelNone.
func New(level Level) *Truncator {
	switch level {
	case LevelShort, LevelTiny:
		return &Truncator{level: level, suffix: "..."}
	default:
		return &Truncator{level: LevelNone}
	}
}

// Apply returns a copy of secrets with values truncated.
func (t *Truncator) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = t.truncate(v)
	}
	return out
}

func (t *Truncator) truncate(v string) string {
	var max int
	switch t.level {
	case LevelShort:
		max = 8
	case LevelTiny:
		max = 3
	default:
		return v
	}
	if len(v) <= max {
		return v
	}
	return strings.TrimRight(v[:max], " ") + t.suffix
}
