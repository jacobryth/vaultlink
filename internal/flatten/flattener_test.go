package flatten

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelUnderscore, LevelDot} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("slash")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	f, _ := New(LevelUnderscore)
	if f.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoFlatten(t *testing.T) {
	f, _ := New(LevelNone)
	in := map[string]string{"db/host": "localhost", "db/port": "5432"}
	out := f.Apply(in)
	if out["db/host"] != "localhost" || out["db/port"] != "5432" {
		t.Fatal("expected keys unchanged for level none")
	}
}

func TestApply_Underscore(t *testing.T) {
	f, _ := New(LevelUnderscore)
	in := map[string]string{"db/host": "localhost", "app/secret/key": "abc"}
	out := f.Apply(in)
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host, got %v", out)
	}
	if out["app_secret_key"] != "abc" {
		t.Errorf("expected app_secret_key, got %v", out)
	}
}

func TestApply_Dot(t *testing.T) {
	f, _ := New(LevelDot)
	in := map[string]string{"db/host": "localhost"}
	out := f.Apply(in)
	if out["db.host"] != "localhost" {
		t.Errorf("expected db.host, got %v", out)
	}
}

func TestApply_NoSlash_Unchanged(t *testing.T) {
	f, _ := New(LevelUnderscore)
	in := map[string]string{"PLAIN_KEY": "value"}
	out := f.Apply(in)
	if out["PLAIN_KEY"] != "value" {
		t.Errorf("expected PLAIN_KEY unchanged, got %v", out)
	}
}
