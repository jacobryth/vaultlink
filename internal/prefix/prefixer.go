package prefix

import (
	"fmt"
	"strings"
)

// Level controls how prefixing is applied.
type Level string

const (
	LevelNone   Level = "none"
	LevelEnv    Level = "env"
	LevelCustom Level = "custom"
)

// Prefixer adds a prefix to secret keys.
type Prefixer struct {
	level  Level
	prefix string
}

// New creates a Prefixer. level must be one of: none, env, custom.
// For "env" level, prefix is derived from the env argument.
// For "custom" level, prefix is used as-is.
func New(level Level, env, custom string) (*Prefixer, error) {
	switch level {
	case LevelNone:
		return &Prefixer{level: LevelNone}, nil
	case LevelEnv:
		if env == "" {
			return nil, fmt.Errorf("prefix: env name required for level 'env'")
		}
		return &Prefixer{level: LevelEnv, prefix: strings.ToUpper(env) + "_"}, nil
	case LevelCustom:
		if custom == "" {
			return nil, fmt.Errorf("prefix: custom prefix required for level 'custom'")
		}
		return &Prefixer{level: LevelCustom, prefix: custom}, nil
	default:
		return nil, fmt.Errorf("prefix: unknown level %q", level)
	}
}

// Apply returns a new map with prefixed keys.
func (p *Prefixer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if p.level == LevelNone {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[p.prefix+k] = v
	}
	return result
}
