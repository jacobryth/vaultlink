package sort

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []string{"none", "asc", "desc"} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("random")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	s, _ := New("asc")
	if s.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoSort(t *testing.T) {
	s, _ := New("none")
	input := map[string]string{"B": "2", "A": "1"}
	out := s.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_Asc(t *testing.T) {
	s, _ := New("asc")
	input := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	out := s.Apply(input)
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}

func TestApply_Desc(t *testing.T) {
	s, _ := New("desc")
	input := map[string]string{"ZEBRA": "z", "ALPHA": "a"}
	out := s.Apply(input)
	if out["ZEBRA"] != "z" || out["ALPHA"] != "a" {
		t.Fatal("values should be preserved after sort")
	}
}
