package redact

import "strings"

// Rule defines a redaction rule for secret keys.
type Rule struct {
	Keywords []string
}

// DefaultRule returns a Rule with common sensitive key patterns.
func DefaultRule() Rule {
	return Rule{
		Keywords: []string{"password", "secret", "token", "key", "private", "credential"},
	}
}

// IsSensitive returns true if the key matches any redaction keyword.
func (r Rule) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, kw := range r.Keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// Apply returns a copy of secrets with sensitive values replaced by "***".
func (r Rule) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if r.IsSensitive(k) {
			result[k] = "***"
		} else {
			result[k] = v
		}
	}
	return result
}
