package trim

import (
	"fmt"
	"strings"
)

const (
	LevelNone    = "none"
	LevelSpace   = "space"
	LevelAll     = "all"
)

var validLevels = map[string]bool{
	LevelNone:  true,
	LevelSpace: true,
	LevelAll:   true,
}

// Trimmer trims whitespace from secret values.
type Trimmer struct {
	level string
}

// New creates a new Trimmer with the given level.
func New(level string) (*Trimmer, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("trim: unknown level %q, must be one of: none, space, all", level)
	}
	return &Trimmer{level: level}, nil
}

// Apply trims secret values according to the configured level.
func (t *Trimmer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if t.level == LevelNone {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch t.level {
		case LevelSpace:
			out[k] = strings.TrimSpace(v)
		case LevelAll:
			out[k] = strings.Map(func(r rune) rune {
				if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
					return -1
				}
				return r
			}, v)
		}
	}
	return out
}
