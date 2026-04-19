package placeholder

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelSelf, LevelStrict} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bogus")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	r, _ := New(LevelSelf)
	out, err := r.Apply(nil)
	if err != nil || out != nil {
		t.Fatalf("expected nil, nil; got %v, %v", out, err)
	}
}

func TestApply_NoReplace(t *testing.T) {
	r, _ := New(LevelNone)
	in := map[string]string{"FOO": "bar", "BAZ": "${FOO}"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["BAZ"] != "${FOO}" {
		t.Errorf("expected placeholder untouched, got %q", out["BAZ"])
	}
}

func TestApply_SelfReplace(t *testing.T) {
	r, _ := New(LevelSelf)
	in := map[string]string{"HOST": "localhost", "DSN": "postgres://${HOST}/db"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("unexpected DSN: %q", out["DSN"])
	}
}

func TestApply_MissingKey_NonStrict(t *testing.T) {
	r, _ := New(LevelSelf)
	in := map[string]string{"VAL": "${MISSING}"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["VAL"] != "${MISSING}" {
		t.Errorf("expected placeholder kept, got %q", out["VAL"])
	}
}

func TestApply_MissingKey_Strict(t *testing.T) {
	r, _ := New(LevelStrict)
	in := map[string]string{"VAL": "${MISSING}"}
	_, err := r.Apply(in)
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}
