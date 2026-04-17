package redact_test

import (
	"testing"

	"github.com/user/vaultlink/internal/redact"
)

func TestIsSensitive_MatchesKeyword(t *testing.T) {
	r := redact.DefaultRule()
	keys := []string{"DB_PASSWORD", "API_TOKEN", "SECRET_KEY", "PRIVATE_KEY", "aws_secret"}
	for _, k := range keys {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_SafeKey(t *testing.T) {
	r := redact.DefaultRule()
	if r.IsSensitive("APP_ENV") {
		t.Error("expected APP_ENV to not be sensitive")
	}
}

func TestApply_RedactsSensitive(t *testing.T) {
	r := redact.DefaultRule()
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_ENV":     "production",
	}
	out := r.Apply(secrets)
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected redacted value, got %q", out["DB_PASSWORD"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected plain value, got %q", out["APP_ENV"])
	}
}

func TestApply_NilInput(t *testing.T) {
	r := redact.DefaultRule()
	out := r.Apply(nil)
	if len(out) != 0 {
		t.Error("expected empty map for nil input")
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	r := redact.DefaultRule()
	out := r.Apply(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty result")
	}
}
