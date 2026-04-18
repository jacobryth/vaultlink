package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for a secret key or value.
type Rule struct {
	KeyPattern   string
	Required     bool
	ValuePattern string
}

// Violation holds a single validation failure.
type Violation struct {
	Key     string
	Message string
}

// Validator checks secrets against a set of rules.
type Validator struct {
	rules []Rule
}

// New creates a Validator with the given rules.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// Validate checks the provided secrets map and returns any violations.
func (v *Validator) Validate(secrets map[string]string) []Violation {
	var violations []Violation

	for _, rule := range v.rules {
		keyRe, err := regexp.Compile(rule.KeyPattern)
		if err != nil {
			violations = append(violations, Violation{Key: rule.KeyPattern, Message: "invalid key pattern: " + err.Error()})
			continue
		}

		matched := false
		for k, val := range secrets {
			if !keyRe.MatchString(k) {
				continue
			}
			matched = true
			if rule.ValuePattern != "" {
				valRe, err := regexp.Compile(rule.ValuePattern)
				if err != nil {
					violations = append(violations, Violation{Key: k, Message: "invalid value pattern: " + err.Error()})
					continue
				}
				if !valRe.MatchString(val) {
					violations = append(violations, Violation{Key: k, Message: fmt.Sprintf("value does not match pattern %q", rule.ValuePattern)})
				}
			}
			if rule.Required && strings.TrimSpace(val) == "" {
				violations = append(violations, Violation{Key: k, Message: "required value is empty"})
			}
		}

		if rule.Required && !matched {
			violations = append(violations, Violation{Key: rule.KeyPattern, Message: "required key not found"})
		}
	}

	return violations
}
