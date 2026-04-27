package env

import (
	"testing"
)

// TestEnvValidator_RoundTrip_KeysAndValues verifies that a clean map passes
// both key and value validation end-to-end.
func TestEnvValidator_RoundTrip_KeysAndValues(t *testing.T) {
	v, err := NewEnvValidator(EnvValidateModeBoth)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	secrets := map[string]string{
		"APP_NAME":    "vaultlink",
		"DB_PASSWORD": "s3cr3t!",
		"_INTERNAL":   "value",
	}

	violations := v.Validate(secrets)
	if len(violations) != 0 {
		t.Errorf("expected no violations for clean map, got: %v", violations)
	}
}

// TestEnvValidator_RoundTrip_DetectsAllBadKeys ensures every bad key in a
// mixed map is flagged exactly once.
func TestEnvValidator_RoundTrip_DetectsAllBadKeys(t *testing.T) {
	v, err := NewEnvValidator(EnvValidateModeKeys)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	secrets := map[string]string{
		"VALID_KEY":  "ok",
		"0STARTS":    "bad",
		"has space":  "bad",
		"ALSO_VALID": "ok",
	}

	violations := v.Validate(secrets)
	if len(violations) != 2 {
		t.Errorf("expected exactly 2 violations, got %d: %v", len(violations), violations)
	}
}
