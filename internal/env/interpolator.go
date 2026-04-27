package env

import (
	"fmt"
	"regexp"
	"strings"
)

// InterpolateMode controls how unresolved references are handled.
type InterpolateMode string

const (
	InterpolateModeNone   InterpolateMode = "none"
	InterpolateModeStrict InterpolateMode = "strict"
	InterpolateModeLoose  InterpolateMode = "loose"
)

var validInterpolateModes = map[InterpolateMode]bool{
	InterpolateModeNone:   true,
	InterpolateModeStrict: true,
	InterpolateModeLoose:  true,
}

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolator resolves ${KEY} references within secret values.
type Interpolator struct {
	mode InterpolateMode
}

// NewInterpolator creates an Interpolator with the given mode.
func NewInterpolator(mode InterpolateMode) (*Interpolator, error) {
	if !validInterpolateModes[mode] {
		return nil, fmt.Errorf("interpolator: unknown mode %q", mode)
	}
	return &Interpolator{mode: mode}, nil
}

// Apply resolves ${KEY} references in each value using the same map as context.
// In strict mode, unresolved references return an error.
// In loose mode, unresolved references are left as-is.
func (i *Interpolator) Apply(secrets map[string]string) (map[string]string, error) {
	if i.mode == InterpolateModeNone {
		return secrets, nil
	}
	if secrets == nil {
		return nil, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := i.resolve(v, secrets)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func (i *Interpolator) resolve(value string, ctx map[string]string) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := strings.TrimSpace(match[2 : len(match)-1])
		if val, ok := ctx[key]; ok {
			return val
		}
		if i.mode == InterpolateModeStrict {
			resolveErr = fmt.Errorf("unresolved reference ${%s}", key)
			return match
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
