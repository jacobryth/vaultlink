package env

import (
	"testing"
)

func TestNewEnvSorter_ValidModes(t *testing.T) {
	for _, m := range []string{"none", "asc", "desc"} {
		s, err := NewEnvSorter(m)
		if err != nil {
			t.Fatalf("mode %q: unexpected error: %v", m, err)
		}
		if s == nil {
			t.Fatalf("mode %q: expected non-nil sorter", m)
		}
	}
}

func TestNewEnvSorter_InvalidMode(t *testing.T) {
	_, err := NewEnvSorter("random")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestEnvSort_NilInput(t *testing.T) {
	s, _ := NewEnvSorter("asc")
	if got := s.Apply(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestEnvSort_NoneMode(t *testing.T) {
	s, _ := NewEnvSorter("none")
	input := []string{"Z=1", "A=2", "M=3"}
	got := s.Apply(input)
	for i, v := range input {
		if got[i] != v {
			t.Fatalf("none mode should not reorder: got %v", got)
		}
	}
}

func TestEnvSort_Asc(t *testing.T) {
	s, _ := NewEnvSorter("asc")
	input := []string{"Z=1", "A=2", "M=3"}
	got := s.Apply(input)
	want := []string{"A=2", "M=3", "Z=1"}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("asc: index %d want %q got %q", i, v, got[i])
		}
	}
}

func TestEnvSort_Desc(t *testing.T) {
	s, _ := NewEnvSorter("desc")
	input := []string{"A=1", "Z=2", "M=3"}
	got := s.Apply(input)
	want := []string{"Z=2", "M=3", "A=1"}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("desc: index %d want %q got %q", i, v, got[i])
		}
	}
}

func TestEnvSort_PassthroughNonKV(t *testing.T) {
	s, _ := NewEnvSorter("asc")
	input := []string{"# comment", "Z=1", "", "A=2"}
	got := s.Apply(input)
	// KV lines sorted first, then passthrough lines appended
	if got[0] != "A=2" || got[1] != "Z=1" {
		t.Fatalf("unexpected order: %v", got)
	}
	if got[2] != "# comment" && got[3] != "" {
		t.Fatalf("passthrough lines missing: %v", got)
	}
}
