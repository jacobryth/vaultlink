package validate

import (
	"testing"
)

func TestValidate_NoViolations(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "^DB_", Required: true, ValuePattern: ".+"},
	})
	secrets := map[string]string{"DB_HOST": "localhost"}
	violations := v.Validate(secrets)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %+v", violations)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "^API_KEY$", Required: true},
	})
	violations := v.Validate(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "^API_KEY$" {
		t.Errorf("unexpected key: %s", violations[0].Key)
	}
}

func TestValidate_ValuePatternMismatch(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "^PORT$", ValuePattern: `^\d+$`},
	})
	secrets := map[string]string{"PORT": "not-a-number"}
	violations := v.Validate(secrets)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_RequiredEmptyValue(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "^SECRET$", Required: true},
	})
	secrets := map[string]string{"SECRET": "   "}
	violations := v.Validate(secrets)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_InvalidKeyPattern(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "[", Required: false},
	})
	violations := v.Validate(map[string]string{"KEY": "val"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for bad pattern, got %d", len(violations))
	}
}

func TestValidate_MultipleRules(t *testing.T) {
	v := New([]Rule{
		{KeyPattern: "^DB_HOST$", Required: true},
		{KeyPattern: "^DB_PORT$", Required: true, ValuePattern: `^\d+$`},
	})
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	violations := v.Validate(secrets)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %+v", violations)
	}
}
