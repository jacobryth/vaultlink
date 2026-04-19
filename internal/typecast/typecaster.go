package typecast

import "fmt"

// Level defines how values are cast.
type Level string

const (
	None   Level = "none"
	String Level = "string"
	Infer  Level = "infer"
)

var validLevels = map[Level]bool{
	None: true, String: true, Infer: true,
}

// Typecaster casts secret values to inferred types as strings.
type Typecaster struct {
	level Level
}

// New returns a Typecaster for the given level.
func New(level Level) (*Typecaster, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("typecast: unknown level %q", level)
	}
	return &Typecaster{level: level}, nil
}

// Apply processes secrets and returns type-annotated or normalized values.
func (t *Typecaster) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if t.level == None || t.level == String {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	// Infer: normalize booleans and numbers to canonical form.
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = inferCast(v)
	}
	return out
}

func inferCast(v string) string {
	switch v {
	case "true", "True", "TRUE":
		return "true"
	case "false", "False", "FALSE":
		return "false"
	}
	return v
}
