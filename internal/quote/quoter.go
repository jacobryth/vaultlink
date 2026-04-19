package quote

import (
	"fmt"
	"strings"
)

type Level string

const (
	None   Level = "none"
	Double Level = "double"
	Single Level = "single"
)

var validLevels = map[Level]bool{
	None:   true,
	Double: true,
	Single: true,
}

type Quoter struct {
	level Level
}

func New(level Level) (*Quoter, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("unknown quote level %q: must be none, double, or single", level)
	}
	return &Quoter{level: level}, nil
}

func (q *Quoter) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if q.level == None {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = q.wrap(v)
	}
	return result
}

func (q *Quoter) wrap(v string) string {
	switch q.level {
	case Double:
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	case Single:
		return `'` + strings.ReplaceAll(v, `'`, `\'`) + `'`
	default:
		return v
	}
}
