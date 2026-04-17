package mask

import "strings"

// Level controls how much of a secret value is masked.
type Level string

const (
	LevelFull    Level = "full"    // ****
	LevelPartial Level = "partial" // ab****
	LevelNone    Level = "none"    // no masking
)

// Masker applies value masking to secret maps.
type Masker struct {
	Level Level
}

// New returns a Masker with the given level.
func New(level Level) *Masker {
	switch level {
	case LevelFull, LevelPartial, LevelNone:
		return &Masker{Level: level}
	default:
		return &Masker{Level: LevelFull}
	}
}

// Apply returns a copy of secrets with values masked according to the level.
func (m *Masker) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.maskValue(v)
	}
	return out
}

func (m *Masker) maskValue(v string) string {
	switch m.Level {
	case LevelNone:
		return v
	case LevelPartial:
		if len(v) <= 2 {
			return strings.Repeat("*", len(v))
		}
		return v[:2] + strings.Repeat("*", len(v)-2)
	default: // full
		return strings.Repeat("*", len(v))
	}
}
