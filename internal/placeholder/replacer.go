package placeholder

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\$\{([^}]+)\}`)

// Level controls placeholder replacement behaviour.
type Level string

const (
	LevelNone   Level = "none"
	LevelSelf   Level = "self"   // replace placeholders using secrets from the same map
	LevelStrict Level = "strict" // like self but errors on missing keys
)

// Replacer resolves ${KEY} placeholders inside secret values.
type Replacer struct {
	level Level
}

// New returns a Replacer for the given level.
func New(level Level) (*Replacer, error) {
	switch level {
	case LevelNone, LevelSelf, LevelStrict:
		return &Replacer{level: level}, nil
	default:
		return nil, fmt.Errorf("placeholder: unknown level %q", level)
	}
}

// Apply resolves placeholders in each value using other entries in secrets.
func (r *Replacer) Apply(secrets map[string]string) (map[string]string, error) {
	if secrets == nil {
		return nil, nil
	}
	if r.level == LevelNone {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result, nil
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		replaced, err := r.resolve(v, secrets)
		if err != nil {
			return nil, fmt.Errorf("placeholder: key %q: %w", k, err)
		}
		result[k] = replaced
	}
	return result, nil
}

func (r *Replacer) resolve(value string, secrets map[string]string) (string, error) {
	var resolveErr error
	resolved := placeholderRe.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		if val, ok := secrets[key]; ok {
			return val
		}
		if r.level == LevelStrict {
			resolveErr = fmt.Errorf("missing key %q", key)
			return match
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return resolved, nil
}
