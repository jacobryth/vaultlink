package defaults

import "fmt"

// Level controls how defaults are applied.
type Level string

const (
	None    Level = "none"
	Missing Level = "missing"
	Empty   Level = "empty"
	Both    Level = "both"
)

var validLevels = map[Level]bool{
	None:    true,
	Missing: true,
	Empty:   true,
	Both:    true,
}

// Defaulter fills in default values for secrets based on a rule map.
type Defaulter struct {
	level    Level
	defaults map[string]string
}

// New creates a Defaulter with the given level and default rules.
// defaults maps key names to their fallback values.
func New(level Level, defaults map[string]string) (*Defaulter, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("defaults: unknown level %q", level)
	}
	if (level == Missing || level == Empty || level == Both) && len(defaults) == 0 {
		return nil, fmt.Errorf("defaults: level %q requires at least one default rule", level)
	}
	return &Defaulter{level: level, defaults: defaults}, nil
}

// Apply fills in default values according to the configured level.
// Missing: only set if the key is absent from secrets.
// Empty:   only set if the key exists but has an empty value.
// Both:    set if key is absent or value is empty.
// None:    returns secrets unchanged.
func (d *Defaulter) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if d.level == None {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for key, fallback := range d.defaults {
		val, exists := out[key]
		switch d.level {
		case Missing:
			if !exists {
				out[key] = fallback
			}
		case Empty:
			if exists && val == "" {
				out[key] = fallback
			}
		case Both:
			if !exists || val == "" {
				out[key] = fallback
			}
		}
	}
	return out
}
