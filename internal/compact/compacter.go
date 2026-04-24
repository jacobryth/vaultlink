package compact

import "fmt"

// Level controls how compaction is applied to secret values.
type Level string

const (
	LevelNone    Level = "none"
	LevelBlank   Level = "blank"
	LevelAll     Level = "all"
)

var validLevels = map[Level]bool{
	LevelNone:  true,
	LevelBlank: true,
	LevelAll:   true,
}

// Compacter removes empty or whitespace-only entries from a secrets map.
type Compacter struct {
	level Level
}

// New returns a Compacter for the given level.
// LevelNone is a no-op.
// LevelBlank removes keys whose values are empty strings.
// LevelAll removes keys whose values are empty or whitespace-only.
func New(level Level) (*Compacter, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("compact: unknown level %q; valid levels are none, blank, all", level)
	}
	return &Compacter{level: level}, nil
}

// Apply returns a new map with entries removed according to the compaction level.
func (c *Compacter) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if c.level == LevelNone {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if c.level == LevelBlank && v == "" {
			continue
		}
		if c.level == LevelAll && isBlankOrEmpty(v) {
			continue
		}
		out[k] = v
	}
	return out
}

func isBlankOrEmpty(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
