package env

import (
	"fmt"
	"regexp"
	"strings"
)

// EnvValidateMode controls how env file key/value pairs are validated.
type EnvValidateMode string

const (
	EnvValidateModeNone   EnvValidateMode = "none"
	EnvValidateModeKeys   EnvValidateMode = "keys"
	EnvValidateModeValues EnvValidateMode = "values"
	EnvValidateModeBoth   EnvValidateMode = "both"
)

var validEnvValidateModes = map[EnvValidateMode]bool{
	EnvValidateModeNone:   true,
	EnvValidateModeKeys:   true,
	EnvValidateModeValues: true,
	EnvValidateModeBoth:   true,
}

// validKeyPattern matches POSIX-compliant env variable names.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// EnvValidator validates key/value pairs in an env map.
type EnvValidator struct {
	mode EnvValidateMode
}

// NewEnvValidator constructs an EnvValidator for the given mode.
func NewEnvValidator(mode EnvValidateMode) (*EnvValidator, error) {
	if !validEnvValidateModes[mode] {
		return nil, fmt.Errorf("env/validator: unknown mode %q", mode)
	}
	return &EnvValidator{mode: mode}, nil
}

// Validate checks the provided secrets map and returns a list of violations.
// Returns nil when mode is none or no violations are found.
func (v *EnvValidator) Validate(secrets map[string]string) []string {
	if secrets == nil || v.mode == EnvValidateModeNone {
		return nil
	}

	var violations []string

	for k, val := range secrets {
		if v.mode == EnvValidateModeKeys || v.mode == EnvValidateModeBoth {
			if !validKeyPattern.MatchString(k) {
				violations = append(violations, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
			}
		}
		if v.mode == EnvValidateModeValues || v.mode == EnvValidateModeBoth {
			if strings.ContainsAny(val, "\x00\r") {
				violations = append(violations, fmt.Sprintf("invalid value for key %q: contains control characters", k))
			}
		}
	}

	return violations
}
