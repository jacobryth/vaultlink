package group

import (
	"fmt"
	"strings"
)

// Level controls how secrets are grouped.
type Level string

const (
	None   Level = "none"
	Prefix Level = "prefix"
	Custom Level = "custom"
)

var knownLevels = map[Level]bool{
	None:   true,
	Prefix: true,
	Custom: true,
}

// Grouper groups secrets into named buckets.
type Grouper struct {
	level     Level
	delimiter string
	rules     map[string]string // key prefix -> group name
}

// New returns a Grouper for the given level.
// delimiter is used for Prefix level (e.g. "_").
// rules is required for Custom level: map of key-prefix to group name.
func New(level Level, delimiter string, rules map[string]string) (*Grouper, error) {
	if !knownLevels[level] {
		return nil, fmt.Errorf("group: unknown level %q", level)
	}
	if level == Custom && len(rules) == 0 {
		return nil, fmt.Errorf("group: custom level requires at least one rule")
	}
	if level == Prefix && delimiter == "" {
		delimiter = "_"
	}
	return &Grouper{level: level, delimiter: delimiter, rules: rules}, nil
}

// Apply groups the given secrets map into named buckets.
// Secrets that don't match any group are placed under "default".
func (g *Grouper) Apply(secrets map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	if secrets == nil {
		return result
	}
	for k, v := range secrets {
		bucket := g.resolve(k)
		if result[bucket] == nil {
			result[bucket] = make(map[string]string)
		}
		result[bucket][k] = v
	}
	return result
}

func (g *Grouper) resolve(key string) string {
	switch g.level {
	case Prefix:
		if idx := strings.Index(key, g.delimiter); idx > 0 {
			return key[:idx]
		}
		return "default"
	case Custom:
		for prefix, name := range g.rules {
			if strings.HasPrefix(key, prefix) {
				return name
			}
		}
		return "default"
	default:
		return "default"
	}
}
