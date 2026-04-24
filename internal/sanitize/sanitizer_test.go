package sanitize

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{None, Strip, Normalize} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("aggressive")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	s, _ := New(Strip)
	if got := s.Apply(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestApply_NoSanitize(t *testing.T) {
	s, _ := New(None)
	input := map[string]string{"KEY": "hello\x01world"}
	out := s.Apply(input)
	if out["KEY"] != "hello\x01world" {
		t.Errorf("None level should not modify values, got %q", out["KEY"])
	}
}

func TestApply_StripControlChars(t *testing.T) {
	s, _ := New(Strip)
	input := map[string]string{
		"A": "hello\x00world",
		"B": "clean",
		"C": "tab\there", // tab is allowed
	}
	out := s.Apply(input)
	if out["A"] != "helloworld" {
		t.Errorf("expected control char removed, got %q", out["A"])
	}
	if out["B"] != "clean" {
		t.Errorf("expected unchanged, got %q", out["B"])
	}
	if out["C"] != "tab\there" {
		t.Errorf("expected tab preserved, got %q", out["C"])
	}
}

func TestApply_NormalizeWhitespace(t *testing.T) {
	s, _ := New(Normalize)
	input := map[string]string{
		"A": "  hello   world  ",
		"B": "no  extra",
		"C": "\x01hidden\x02",
	}
	out := s.Apply(input)
	if out["A"] != "hello world" {
		t.Errorf("expected normalized, got %q", out["A"])
	}
	if out["B"] != "no extra" {
		t.Errorf("expected collapsed spaces, got %q", out["B"])
	}
	if out["C"] != "hidden" {
		t.Errorf("expected control chars stripped, got %q", out["C"])
	}
}

func TestApply_PreservesAllKeys(t *testing.T) {
	s, _ := New(Strip)
	input := map[string]string{"X": "a", "Y": "b", "Z": "c"}
	out := s.Apply(input)
	if len(out) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(out))
	}
}
