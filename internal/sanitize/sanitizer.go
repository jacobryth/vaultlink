package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

// Level controls how aggressively values are sanitized.
type Level string

const (
	None    Level = "none"
	Strip   Level = "strip"   // remove non-printable / control characters
	Normalize Level = "normalize" // strip + collapse whitespace
)

var validLevels = map[Level]bool{
	None:      true,
	Strip:     true,
	Normalize: true,
}

// nonPrintable matches ASCII control characters (except tab/newline).
var nonPrintable = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

// multiSpace matches runs of whitespace (space/tab).
var multiSpace = regexp.MustCompile(`[ \t]{2,}`)

// Sanitizer cleans secret values according to the configured level.
type Sanitizer struct {
	level Level
}

// New creates a Sanitizer for the given level.
func New(level Level) (*Sanitizer, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("sanitize: unknown level %q (want none|strip|normalize)", level)
	}
	return &Sanitizer{level: level}, nil
}

// Apply returns a sanitized copy of secrets. The original map is not modified.
func (s *Sanitizer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if s.level == None {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = s.sanitizeValue(v)
	}
	return out
}

func (s *Sanitizer) sanitizeValue(v string) string {
	// Always strip control characters for strip and normalize.
	v = nonPrintable.ReplaceAllString(v, "")
	if s.level == Normalize {
		v = strings.TrimSpace(v)
		v = multiSpace.ReplaceAllString(v, " ")
	}
	return v
}
