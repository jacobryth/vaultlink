package comment

import (
	"fmt"
	"strings"
)

// Level controls how comments are added to env output.
type Level string

const (
	LevelNone   Level = "none"
	LevelInline Level = "inline"
	LevelBlock  Level = "block"
)

var validLevels = map[Level]bool{
	LevelNone:   true,
	LevelInline: true,
	LevelBlock:  true,
}

// Commenter annotates secret keys with comments.
type Commenter struct {
	level  Level
	source string
}

// New returns a Commenter for the given level and source label.
func New(level Level, source string) (*Commenter, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("comment: unknown level %q", level)
	}
	return &Commenter{level: level, source: source}, nil
}

// Apply returns a new map with comments embedded as pseudo-keys or
// returns the original map unchanged when level is none.
func (c *Commenter) Apply(secrets map[string]string) map[string]string {
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
		switch c.level {
		case LevelInline:
			// Store annotation alongside value separated by a sentinel.
			out[k] = fmt.Sprintf("%s # source:%s", v, c.source)
		case LevelBlock:
			// Prefix key with a comment marker key that writers can detect.
			commentKey := fmt.Sprintf("#%s", strings.ToUpper(k))
			out[commentKey] = fmt.Sprintf("source:%s", c.source)
			out[k] = v
		}
	}
	return out
}
