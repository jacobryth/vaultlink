package typecast

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, l := range []Level{None, String, Infer} {
		_, err := New(l)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", l, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	tc, _ := New(Infer)
	if tc.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoTypecast(t *testing.T) {
	tc, _ := New(None)
	in := map[string]string{"KEY": "True"}
	out := tc.Apply(in)
	if out["KEY"] != "True" {
		t.Errorf("expected 'True', got %q", out["KEY"])
	}
}

func TestApply_InferBooleanTrue(t *testing.T) {
	tc, _ := New(Infer)
	for _, raw := range []string{"true", "True", "TRUE"} {
		out := tc.Apply(map[string]string{"K": raw})
		if out["K"] != "true" {
			t.Errorf("expected 'true' for input %q, got %q", raw, out["K"])
		}
	}
}

func TestApply_InferBooleanFalse(t *testing.T) {
	tc, _ := New(Infer)
	for _, raw := range []string{"false", "False", "FALSE"} {
		out := tc.Apply(map[string]string{"K": raw})
		if out["K"] != "false" {
			t.Errorf("expected 'false' for input %q, got %q", raw, out["K"])
		}
	}
}

func TestApply_InferPassthrough(t *testing.T) {
	tc, _ := New(Infer)
	in := map[string]string{"URL": "http://example.com"}
	out := tc.Apply(in)
	if out["URL"] != "http://example.com" {
		t.Errorf("unexpected value: %q", out["URL"])
	}
}
