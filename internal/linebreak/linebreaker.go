package linebreak

import (
	"fmt"
	"strings"
)

const (
	LevelNone  = "none"
	LevelUnix  = "unix"
	LevelWindows = "windows"
)

var validLevels = map[string]bool{
	LevelNone:    true,
	LevelUnix:    true,
	LevelWindows: true,
}

// Linebreaker normalizes line endings in secret values.
type Linebreaker struct {
	level string
}

// New creates a Linebreaker with the given level.
func New(level string) (*Linebreaker, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("linebreak: unknown level %q, must be one of: none, unix, windows", level)
	}
	return &Linebreaker{level: level}, nil
}

// Apply normalizes line endings in each secret value based on the configured level.
func (l *Linebreaker) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if l.level == LevelNone {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		// Normalize to unix first
		normalized := strings.ReplaceAll(v, "\r\n", "\n")
		normalized = strings.ReplaceAll(normalized, "\r", "\n")
		if l.level == LevelWindows {
			normalized = strings.ReplaceAll(normalized, "\n", "\r\n")
		}
		result[k] = normalized
	}
	return result
}
