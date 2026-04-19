package linebreak

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, level := range []string{LevelNone, LevelUnix, LevelWindows} {
		_, err := New(level)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", level, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("crlf")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	lb, _ := New(LevelUnix)
	if lb.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoLinebreak(t *testing.T) {
	lb, _ := New(LevelNone)
	input := map[string]string{"KEY": "value\r\nother"}
	out := lb.Apply(input)
	if out["KEY"] != "value\r\nother" {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_UnixNormalization(t *testing.T) {
	lb, _ := New(LevelUnix)
	input := map[string]string{"KEY": "line1\r\nline2\rline3"}
	out := lb.Apply(input)
	if out["KEY"] != "line1\nline2\nline3" {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}

func TestApply_WindowsNormalization(t *testing.T) {
	lb, _ := New(LevelWindows)
	input := map[string]string{"KEY": "line1\nline2\r\nline3"}
	out := lb.Apply(input)
	expected := "line1\r\nline2\r\nline3"
	if out["KEY"] != expected {
		t.Errorf("expected %q, got %q", expected, out["KEY"])
	}
}

func TestApply_AlreadyUnix(t *testing.T) {
	lb, _ := New(LevelUnix)
	input := map[string]string{"KEY": "line1\nline2"}
	out := lb.Apply(input)
	if out["KEY"] != "line1\nline2" {
		t.Errorf("unexpected value: %q", out["KEY"])
	}
}
