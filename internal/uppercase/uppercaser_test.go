package uppercase

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, level := range []string{LevelNone, LevelKeys, LevelValues, LevelBoth} {
		_, err := New(level)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", level, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("invalid")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	u, _ := New(LevelKeys)
	if u.Apply(nil) != nil {
		t.Fatal("expected nil result for nil input")
	}
}

func TestApply_NoUppercase(t *testing.T) {
	u, _ := New(LevelNone)
	input := map[string]string{"key": "value"}
	out := u.Apply(input)
	if out["key"] != "value" {
		t.Fatalf("expected unchanged, got %v", out)
	}
}

func TestApply_UpperKeys(t *testing.T) {
	u, _ := New(LevelKeys)
	input := map[string]string{"db_host": "localhost"}
	out := u.Apply(input)
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatalf("expected uppercased key, got %v", out)
	}
	if out["DB_HOST"] != "localhost" {
		t.Fatalf("expected value unchanged, got %v", out["DB_HOST"])
	}
}

func TestApply_UpperValues(t *testing.T) {
	u, _ := New(LevelValues)
	input := map[string]string{"mode": "debug"}
	out := u.Apply(input)
	if out["mode"] != "DEBUG" {
		t.Fatalf("expected uppercased value, got %v", out["mode"])
	}
}

func TestApply_UpperBoth(t *testing.T) {
	u, _ := New(LevelBoth)
	input := map[string]string{"env": "production"}
	out := u.Apply(input)
	if _, ok := out["ENV"]; !ok {
		t.Fatalf("expected uppercased key")
	}
	if out["ENV"] != "PRODUCTION" {
		t.Fatalf("expected uppercased value, got %v", out["ENV"])
	}
}
