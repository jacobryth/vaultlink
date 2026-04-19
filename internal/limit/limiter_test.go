package limit

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{None, First, Last} {
		count := 2
		if lvl == None {
			count = 0
		}
		_, err := New(lvl, count)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("middle", 3)
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestNew_ZeroCount(t *testing.T) {
	_, err := New(First, 0)
	if err == nil {
		t.Fatal("expected error for zero count with First level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	l, _ := New(None, 0)
	if l.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoLimit(t *testing.T) {
	l, _ := New(None, 0)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := l.Apply(secrets)
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestApply_First(t *testing.T) {
	l, _ := New(First, 2)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := l.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["A"]; !ok {
		t.Error("expected key A in first-2 result")
	}
	if _, ok := out["B"]; !ok {
		t.Error("expected key B in first-2 result")
	}
}

func TestApply_Last(t *testing.T) {
	l, _ := New(Last, 2)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out := l.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["C"]; !ok {
		t.Error("expected key C in last-2 result")
	}
	if _, ok := out["B"]; !ok {
		t.Error("expected key B in last-2 result")
	}
}

func TestApply_CountExceedsKeys(t *testing.T) {
	l, _ := New(First, 10)
	secrets := map[string]string{"X": "1", "Y": "2"}
	out := l.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}
