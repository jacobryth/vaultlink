package quote

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, l := range []Level{None, Double, Single} {
		_, err := New(l)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", l, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("backtick")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	q, _ := New(Double)
	if q.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoQuote(t *testing.T) {
	q, _ := New(None)
	in := map[string]string{"KEY": "value"}
	out := q.Apply(in)
	if out["KEY"] != "value" {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_DoubleQuote(t *testing.T) {
	q, _ := New(Double)
	in := map[string]string{"KEY": "hello world"}
	out := q.Apply(in)
	if out["KEY"] != `"hello world"` {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}

func TestApply_SingleQuote(t *testing.T) {
	q, _ := New(Single)
	in := map[string]string{"KEY": "hello"}
	out := q.Apply(in)
	if out["KEY"] != `'hello'` {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}

func TestApply_DoubleQuote_EscapesInner(t *testing.T) {
	q, _ := New(Double)
	in := map[string]string{"KEY": `say "hi"`}
	out := q.Apply(in)
	expected := `"say \"hi\""`
	if out["KEY"] != expected {
		t.Errorf("expected %q, got %q", expected, out["KEY"])
	}
}
