package trim

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, level := range []string{LevelNone, LevelSpace, LevelAll} {
		_, err := New(level)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", level, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	tr, _ := New(LevelSpace)
	if tr.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoTrim(t *testing.T) {
	tr, _ := New(LevelNone)
	in := map[string]string{"KEY": "  value  "}
	out := tr.Apply(in)
	if out["KEY"] != "  value  " {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_SpaceTrim(t *testing.T) {
	tr, _ := New(LevelSpace)
	in := map[string]string{"KEY": "  hello world  "}
	out := tr.Apply(in)
	if out["KEY"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApply_AllTrim(t *testing.T) {
	tr, _ := New(LevelAll)
	in := map[string]string{"KEY": "  he llo\tworld\n  "}
	out := tr.Apply(in)
	if out["KEY"] != "helloworld" {
		t.Errorf("expected fully trimmed value, got %q", out["KEY"])
	}
}

func TestApply_MultipleKeys(t *testing.T) {
	tr, _ := New(LevelSpace)
	in := map[string]string{"A": "  foo  ", "B": "\tbar\t"}
	out := tr.Apply(in)
	if out["A"] != "foo" || out["B"] != "bar" {
		t.Errorf("unexpected values: A=%q B=%q", out["A"], out["B"])
	}
}
