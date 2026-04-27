package env

import (
	"testing"
)

func TestNewEnvValidator_ValidModes(t *testing.T) {
	modes := []EnvValidateMode{
		EnvValidateModeNone,
		EnvValidateModeKeys,
		EnvValidateModeValues,
		EnvValidateModeBoth,
	}
	for _, m := range modes {
		_, err := NewEnvValidator(m)
		if err != nil {
			t.Errorf("expected no error for mode %q, got %v", m, err)
		}
	}
}

func TestNewEnvValidator_InvalidMode(t *testing.T) {
	_, err := NewEnvValidator("strict")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestEnvValidate_NilSecrets(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeBoth)
	violations := v.Validate(nil)
	if violations != nil {
		t.Errorf("expected nil violations for nil input, got %v", violations)
	}
}

func TestEnvValidate_NoneMode_PassThrough(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeNone)
	secrets := map[string]string{"123bad": "value"}
	violations := v.Validate(secrets)
	if violations != nil {
		t.Errorf("expected no violations in none mode, got %v", violations)
	}
}

func TestEnvValidate_InvalidKey(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeKeys)
	secrets := map[string]string{"123invalid": "ok", "VALID_KEY": "ok"}
	violations := v.Validate(secrets)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d: %v", len(violations), violations)
	}
}

func TestEnvValidate_InvalidValue(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeValues)
	secrets := map[string]string{"KEY": "bad\x00value"}
	violations := v.Validate(secrets)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestEnvValidate_BothMode_NoViolations(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeBoth)
	secrets := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	violations := v.Validate(secrets)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestEnvValidate_BothMode_MultipleViolations(t *testing.T) {
	v, _ := NewEnvValidator(EnvValidateModeBoth)
	secrets := map[string]string{"bad key!": "val\r\n", "OK_KEY": "clean"}
	violations := v.Validate(secrets)
	if len(violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d: %v", len(violations), violations)
	}
}
