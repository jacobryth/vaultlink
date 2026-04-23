package timestamp

import (
	"fmt"
	"time"
)

// Level controls how timestamps are applied to secret keys.
type Level string

const (
	None   Level = "none"
	Suffix Level = "suffix"
	Prefix Level = "prefix"
)

var validLevels = map[Level]bool{
	None:   true,
	Suffix: true,
	Prefix: true,
}

// Stamper appends or prepends a timestamp marker to secret keys.
type Stamper struct {
	level  Level
	format string
}

// New creates a Stamper with the given level and time format.
// format defaults to "20060102" if empty.
func New(level Level, format string) (*Stamper, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("timestamp: unknown level %q", level)
	}
	if format == "" {
		format = "20060102"
	}
	return &Stamper{level: level, format: format}, nil
}

// Apply returns a new map with timestamp markers added to keys according to
// the configured level. The original map is not modified.
func (s *Stamper) Apply(secrets map[string]string) map[string]string {
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

	stamp := time.Now().UTC().Format(s.format)
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		var newKey string
		switch s.level {
		case Suffix:
			newKey = fmt.Sprintf("%s_%s", k, stamp)
		case Prefix:
			newKey = fmt.Sprintf("%s_%s", stamp, k)
		default:
			newKey = k
		}
		out[newKey] = v
	}
	return out
}
