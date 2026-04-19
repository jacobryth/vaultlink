package uppercase

import (
	"fmt"
	"strings"
)

const (
	LevelNone  = "none"
	LevelKeys  = "keys"
	LevelValues = "values"
	LevelBoth  = "both"
)

var validLevels = map[string]bool{
	LevelNone:   true,
	LevelKeys:   true,
	LevelValues: true,
	LevelBoth:   true,
}

// Uppercaser applies uppercase transformation to secret keys and/or values.
type Uppercaser struct {
	level string
}

// New creates a new Uppercaser with the given level.
func New(level string) (*Uppercaser, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("uppercase: unknown level %q", level)
	}
	return &Uppercaser{level: level}, nil
}

// Apply transforms the secrets map according to the configured level.
func (u *Uppercaser) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if u.level == LevelNone {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := k
		newVal := v
		if u.level == LevelKeys || u.level == LevelBoth {
			newKey = strings.ToUpper(k)
		}
		if u.level == LevelValues || u.level == LevelBoth {
			newVal = strings.ToUpper(v)
		}
		result[newKey] = newVal
	}
	return result
}
