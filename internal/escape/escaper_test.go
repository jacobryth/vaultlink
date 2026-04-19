package escape

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, l := range []Level{None, Shell, Newline} {
		_, err := New(l)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", l, err)
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
	e, _ := New(Shell)
	if e.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoEscape(t *testing.T) {
	e, _ := New(None)
	in := map[string]string{"KEY": "value with $dollar"}
	out := e.Apply(in)
	if out["KEY"] != "value with $dollar" {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_ShellEscape(t *testing.T) {
	e, _ := New(Shell)
	in := map[string]string{
		"A": `say "hello"`,
		"B": `cost $5`,
		"C": `back\slash`,
	}
	out := e.Apply(in)
	if out["A"] != `say \"hello\"` {
		t.Errorf("unexpected A: %s", out["A"])
	}
	if out["B"] != `cost \$5` {
		t.Errorf("unexpected B: %s", out["B"])
	}
	if out["C"] != `back\\slash` {
		t.Errorf("unexpected C: %s", out["C"])
	}
}

func TestApply_NewlineEscape(t *testing.T) {
	e, _ := New(Newline)
	in := map[string]string{"KEY": "line1\nline2\r"}
	out := e.Apply(in)
	expected := `line1\nline2\r`
	if out["KEY"] != expected {
		t.Errorf("expected %q got %q", expected, out["KEY"])
	}
}
