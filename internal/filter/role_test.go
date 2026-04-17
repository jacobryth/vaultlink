package filter_test

import (
	"testing"

	"github.com/yourorg/vaultlink/internal/filter"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"AWS_KEY":     "AKIA123",
	"AWS_SECRET":  "abc",
	"APP_DEBUG":   "true",
}

func TestApply_KnownRole(t *testing.T) {
	f := filter.NewFilter([]filter.Role{
		{Name: "backend", Prefixes: []string{"DB_", "APP_"}},
	})

	result := f.Apply("backend", testSecrets)

	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
	if _, ok := result["AWS_KEY"]; ok {
		t.Error("AWS_KEY should not be present for backend role")
	}
}

func TestApply_UnknownRole(t *testing.T) {
	f := filter.NewFilter([]filter.Role{})
	result := f.Apply("nonexistent", testSecrets)
	if len(result) != 0 {
		t.Errorf("expected empty map for unknown role, got %d keys", len(result))
	}
}

func TestApply_EmptyPrefixes_AllowsAll(t *testing.T) {
	f := filter.NewFilter([]filter.Role{
		{Name: "admin", Prefixes: []string{}},
	})
	result := f.Apply("admin", testSecrets)
	if len(result) != len(testSecrets) {
		t.Errorf("expected all %d keys, got %d", len(testSecrets), len(result))
	}
}

func TestHasRole(t *testing.T) {
	f := filter.NewFilter([]filter.Role{
		{Name: "ops", Prefixes: []string{"AWS_"}},
	})
	if !f.HasRole("ops") {
		t.Error("expected HasRole to return true for 'ops'")
	}
	if f.HasRole("dev") {
		t.Error("expected HasRole to return false for 'dev'")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	f := filter.NewFilter([]filter.Role{
		{Name: "backend", Prefixes: []string{"DB_"}},
	})
	result := f.Apply("backend", nil)
	if len(result) != 0 {
		t.Errorf("expected empty map for nil secrets, got %d keys", len(result))
	}
}
