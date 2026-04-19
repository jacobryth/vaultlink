package limit

import (
	"fmt"
	"sort"
)

type Level string

const (
	None  Level = "none"
	First Level = "first"
	Last  Level = "last"
)

var validLevels = map[Level]bool{
	None:  true,
	First: true,
	Last:  true,
}

// Limiter restricts the number of secrets returned.
type Limiter struct {
	level Level
	count int
}

// New creates a Limiter with the given level and count.
// count is ignored when level is None.
func New(level Level, count int) (*Limiter, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("limit: unknown level %q", level)
	}
	if level != None && count <= 0 {
		return nil, fmt.Errorf("limit: count must be positive, got %d", count)
	}
	return &Limiter{level: level, count: count}, nil
}

// Apply returns a subset of secrets based on the configured level.
func (l *Limiter) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if l.level == None {
		return secrets
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if l.level == Last {
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	if l.count < len(keys) {
		keys = keys[:l.count]
	}

	result := make(map[string]string, len(keys))
	for _, k := range keys {
		result[k] = secrets[k]
	}
	return result
}
