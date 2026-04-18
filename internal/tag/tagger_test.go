package tag

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelAll, LevelEnv} {
		_, err := New(lvl, "PRE_")
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bad", "PRE_")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_DefaultPrefix(t *testing.T) {
	tg, _ := New(LevelAll, "")
	if tg.prefix != "APP_" {
		t.Errorf("expected default prefix APP_, got %q", tg.prefix)
	}
}

func TestApply_NilSecrets(t *testing.T) {
	tg, _ := New(LevelNone, "X_")
	if tg.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoTag(t *testing.T) {
	tg, _ := New(LevelNone, "X_")
	in := map[string]string{"key": "val"}
	out := tg.Apply(in)
	if out["key"] != "val" {
		t.Errorf("expected key unchanged, got %v", out)
	}
}

func TestApply_AllTag(t *testing.T) {
	tg, _ := New(LevelAll, "PRE_")
	in := map[string]string{"KEY": "val", "lower": "v2"}
	out := tg.Apply(in)
	if _, ok := out["PRE_KEY"]; !ok {
		t.Error("expected PRE_KEY in output")
	}
	if _, ok := out["PRE_lower"]; !ok {
		t.Error("expected PRE_lower in output")
	}
}

func TestApply_EnvTag(t *testing.T) {
	tg, _ := New(LevelEnv, "PRE_")
	in := map[string]string{"UPPER": "v1", "lower": "v2"}
	out := tg.Apply(in)
	if _, ok := out["PRE_UPPER"]; !ok {
		t.Error("expected PRE_UPPER tagged")
	}
	if _, ok := out["lower"]; !ok {
		t.Error("expected lower unchanged")
	}
	if _, ok := out["PRE_lower"]; ok {
		t.Error("did not expect PRE_lower")
	}
}
