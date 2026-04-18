package rename

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, level := range []Level{LevelNone, LevelSnake, LevelKebab} {
		_, err := New(level, nil)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", level, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("camel", nil)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_CustomRequiresRules(t *testing.T) {
	_, err := New(LevelCustom, nil)
	if err == nil {
		t.Fatal("expected error for custom level without rules")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	r, _ := New(LevelNone, nil)
	if r.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoRename(t *testing.T) {
	r, _ := New(LevelNone, nil)
	in := map[string]string{"MY_KEY": "value"}
	out := r.Apply(in)
	if out["MY_KEY"] != "value" {
		t.Errorf("expected key unchanged, got %v", out)
	}
}

func TestApply_Snake(t *testing.T) {
	r, _ := New(LevelSnake, nil)
	in := map[string]string{"my-key": "v", "other-key": "x"}
	out := r.Apply(in)
	if out["MY_KEY"] != "v" || out["OTHER_KEY"] != "x" {
		t.Errorf("unexpected snake output: %v", out)
	}
}

func TestApply_Kebab(t *testing.T) {
	r, _ := New(LevelKebab, nil)
	in := map[string]string{"MY_KEY": "v"}
	out := r.Apply(in)
	if out["my-key"] != "v" {
		t.Errorf("unexpected kebab output: %v", out)
	}
}

func TestApply_Custom(t *testing.T) {
	rules := []Rule{{From: "OLD_KEY", To: "NEW_KEY"}}
	r, _ := New(LevelCustom, rules)
	in := map[string]string{"OLD_KEY": "val", "KEEP": "x"}
	out := r.Apply(in)
	if out["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY, got %v", out)
	}
	if out["KEEP"] != "x" {
		t.Errorf("expected KEEP unchanged, got %v", out)
	}
}
