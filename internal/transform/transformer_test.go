package transform

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	levels := []Level{LevelNone, LevelUpper, LevelLower, LevelTrim}
	for _, l := range levels {
		_, err := New(l)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", l, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New(Level("base64"))
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	tr, _ := New(LevelUpper)
	if tr.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoTransform(t *testing.T) {
	tr, _ := New(LevelNone)
	in := map[string]string{"KEY": "  Value  "}
	out := tr.Apply(in)
	if out["KEY"] != "  Value  " {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_Upper(t *testing.T) {
	tr, _ := New(LevelUpper)
	out := tr.Apply(map[string]string{"k": "hello"})
	if out["k"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["k"])
	}
}

func TestApply_Lower(t *testing.T) {
	tr, _ := New(LevelLower)
	out := tr.Apply(map[string]string{"k": "WORLD"})
	if out["k"] != "world" {
		t.Errorf("expected world, got %q", out["k"])
	}
}

func TestApply_Trim(t *testing.T) {
	tr, _ := New(LevelTrim)
	out := tr.Apply(map[string]string{"k": "  spaced  "})
	if out["k"] != "spaced" {
		t.Errorf("expected 'spaced', got %q", out["k"])
	}
}
