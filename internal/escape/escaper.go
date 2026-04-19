package escape

import (
	"fmt"
	"strings"
)

type Level string

const (
	None    Level = "none"
	Shell   Level = "shell"
	Newline Level = "newline"
)

var knownLevels = map[Level]bool{
	None:    true,
	Shell:   true,
	Newline: true,
}

type Escaper struct {
	level Level
}

func New(level Level) (*Escaper, error) {
	if !knownLevels[level] {
		return nil, fmt.Errorf("escape: unknown level %q", level)
	}
	return &Escaper{level: level}, nil
}

func (e *Escaper) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if e.level == None {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = e.escape(v)
	}
	return out
}

func (e *Escaper) escape(val string) string {
	switch e.level {
	case Shell:
		v := strings.ReplaceAll(val, `\`, `\\`)
		v = strings.ReplaceAll(v, `"`, `\"`)
		v = strings.ReplaceAll(v, `$`, `\$`)
		v = strings.ReplaceAll(v, "`", "\\`")
		return v
	case Newline:
		v := strings.ReplaceAll(val, "\n", `\n`)
		v = strings.ReplaceAll(v, "\r", `\r`)
		return v
	}
	return val
}
