package rename

import (
	"fmt"
	"strings"
)

// Level controls how keys are renamed.
type Level string

const (
	LevelNone   Level = "none"
	LevelSnake  Level = "snake"
	LevelKebab  Level = "kebab"
	LevelCustom Level = "custom"
)

// Rule defines a custom rename mapping.
type Rule struct {
	From string
	To   string
}

// Renamer renames secret keys based on a configured level.
type Renamer struct {
	level Level
	rules []Rule
}

// New creates a Renamer. For LevelCustom, rules must be provided.
func New(level Level, rules []Rule) (*Renamer, error) {
	switch level {
	case LevelNone, LevelSnake, LevelKebab:
		return &Renamer{level: level}, nil
	case LevelCustom:
		if len(rules) == 0 {
			return nil, fmt.Errorf("rename: custom level requires at least one rule")
		}
		return &Renamer{level: level, rules: rules}, nil
	default:
		return nil, fmt.Errorf("rename: unknown level %q", level)
	}
}

// Apply renames keys in the secrets map and returns a new map.
func (r *Renamer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[r.renameKey(k)] = v
	}
	return out
}

func (r *Renamer) renameKey(key string) string {
	switch r.level {
	case LevelSnake:
		return strings.ReplaceAll(strings.ToUpper(key), "-", "_")
	case LevelKebab:
		return strings.ReplaceAll(strings.ToLower(key), "_", "-")
	case LevelCustom:
		for _, rule := range r.rules {
			if rule.From == key {
				return rule.To
			}
		}
		return key
	default:
		return key
	}
}
