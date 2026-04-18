package tag

import (
	"fmt"
	"strings"
)

// Level controls which keys get tagged.
type Level string

const (
	LevelNone Level = "none"
	LevelAll  Level = "all"
	LevelEnv  Level = "env"
)

// Tagger adds a prefix tag to secret keys.
type Tagger struct {
	level  Level
	prefix string
}

// New returns a Tagger. prefix is the string prepended to keys.
func New(level Level, prefix string) (*Tagger, error) {
	switch level {
	case LevelNone, LevelAll, LevelEnv:
		// valid
	default:
		return nil, fmt.Errorf("tag: unknown level %q", level)
	}
	if prefix == "" {
		prefix = "APP_"
	}
	return &Tagger{level: level, prefix: prefix}, nil
}

// Apply returns a new map with keys tagged according to the level.
func (t *Tagger) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch t.level {
		case LevelAll:
			out[t.prefix+k] = v
		case LevelEnv:
			if strings.ToUpper(k) == k {
				out[t.prefix+k] = v
			} else {
				out[k] = v
			}
		default:
			out[k] = v
		}
	}
	return out
}
