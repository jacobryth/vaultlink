package env

import (
	"fmt"
	"strings"
)

// FormatLevel controls how key=value pairs are formatted when writing.
type FormatLevel string

const (
	FormatNone    FormatLevel = "none"
	FormatExport  FormatLevel = "export"
	FormatInlined FormatLevel = "inlined"
)

var validFormatLevels = map[FormatLevel]bool{
	FormatNone:    true,
	FormatExport:  true,
	FormatInlined: true,
}

// Formatter applies a formatting convention to secret key=value pairs.
type Formatter struct {
	level FormatLevel
}

// NewFormatter returns a Formatter for the given level.
// Returns an error if the level is unknown.
func NewFormatter(level FormatLevel) (*Formatter, error) {
	if !validFormatLevels[level] {
		return nil, fmt.Errorf("env/formatter: unknown level %q; valid levels are: none, export, inlined", level)
	}
	return &Formatter{level: level}, nil
}

// Apply formats the given secrets map into a slice of strings.
// none     → KEY=VALUE
// export   → export KEY=VALUE
// inlined  → KEY=VALUE; (semicolon-separated, single line)
func (f *Formatter) Apply(secrets map[string]string) []string {
	if secrets == nil {
		return nil
	}

	lines := make([]string, 0, len(secrets))
	for k, v := range secrets {
		switch f.level {
		case FormatExport:
			lines = append(lines, fmt.Sprintf("export %s=%s", k, v))
		default:
			lines = append(lines, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if f.level == FormatInlined {
		return []string{strings.Join(lines, "; ")}
	}
	return lines
}
