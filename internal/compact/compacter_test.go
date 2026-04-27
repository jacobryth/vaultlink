package compact

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelBlank, LevelAll} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("expected error for invalid level, got nil")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	c, _ := New(LevelBlank)
	if got := c.Apply(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestApply_NoCompaction(t *testing.T) {
	c, _ := New(LevelNone)
	input := map[string]string{"A": "", "B": "val", "C": "  "}
	got := c.Apply(input)
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d", len(got))
	}
}

func TestApply_BlankRemovesEmpty(t *testing.T) {
	c, _ := New(LevelBlank)
	input := map[string]string{"A": "", "B": "value", "C": "  "}
	got := c.Apply(input)
	if _, ok := got["A"]; ok {
		t.Error("expected key A to be removed")
	}
	if _, ok := got["B"]; !ok {
		t.Error("expected key B to be kept")
	}
	// whitespace-only value survives blank mode
	if _, ok := got["C"]; !ok {
		t.Error("expected key C to be kept in blank mode")
	}
}

func TestApply_AllRemovesWhitespace(t *testing.T) {
	c, _ := New(LevelAll)
	input := map[string]string{"A": "", "B": "value", "C": "   ", "D": "\t\n"}
	got := c.Apply(input)
	if _, ok := got["A"]; ok {
		t.Error("expected key A to be removed")
	}
	if _, ok := got["C"]; ok {
		t.Error("expected key C to be removed")
	}
	if _, ok := got["D"]; ok {
		t.Error("expected key D to be removed")
	}
	if _, ok := got["B"]; !ok {
		t.Error("expected key B to be kept")
	}
}

func TestApply_PreservesOriginal(t *testing.T) {
	c, _ := New(LevelAll)
	input := map[string]string{"X": "keep"}
	got := c.Apply(input)
	got["X"] = "mutated"
	if input["X"] != "keep" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelBlank, LevelAll} {
		c, _ := New(lvl)
		got := c.Apply(map[string]string{})
		if len(got) != 0 {
			t.Errorf("level %q: expected empty map, got %d entries", lvl, len(got))
		}
	}
}
