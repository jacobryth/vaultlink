package sample

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelRandom, LevelNth} {
		_, err := New(lvl, 2)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bogus", 1)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_ZeroN_NonNone(t *testing.T) {
	_, err := New(LevelRandom, 0)
	if err == nil {
		t.Fatal("expected error when n < 1 for random level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	s, _ := New(LevelNone, 1)
	if s.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoSample(t *testing.T) {
	s, _ := New(LevelNone, 1)
	in := map[string]string{"A": "1", "B": "2"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_RandomSample(t *testing.T) {
	s, _ := New(LevelRandom, 2)
	in := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	for k := range out {
		if _, ok := in[k]; !ok {
			t.Errorf("sampled key %q not in original", k)
		}
	}
}

func TestApply_RandomSample_ExceedsSize(t *testing.T) {
	s, _ := New(LevelRandom, 100)
	in := map[string]string{"X": "1", "Y": "2"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected all 2 keys, got %d", len(out))
	}
}

func TestApply_NthSample(t *testing.T) {
	s, _ := New(LevelNth, 2)
	in := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	out := s.Apply(in)
	// Every 2nd key (indices 0, 2) → 2 keys
	if len(out) == 0 {
		t.Fatal("expected at least one key from nth sampling")
	}
	for k := range out {
		if _, ok := in[k]; !ok {
			t.Errorf("nth key %q not in original", k)
		}
	}
}
