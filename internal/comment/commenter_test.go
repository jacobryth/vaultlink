package comment

import (
	"strings"
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{LevelNone, LevelInline, LevelBlock} {
		_, err := New(lvl, "vault")
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("fancy", "vault")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	c, _ := New(LevelNone, "vault")
	if c.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoComment(t *testing.T) {
	c, _ := New(LevelNone, "vault")
	in := map[string]string{"KEY": "value"}
	out := c.Apply(in)
	if out["KEY"] != "value" {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_InlineComment(t *testing.T) {
	c, _ := New(LevelInline, "vault")
	in := map[string]string{"DB_PASS": "secret"}
	out := c.Apply(in)
	if !strings.Contains(out["DB_PASS"], "# source:vault") {
		t.Errorf("expected inline comment, got: %s", out["DB_PASS"])
	}
}

func TestApply_BlockComment(t *testing.T) {
	c, _ := New(LevelBlock, "prod-vault")
	in := map[string]string{"API_KEY": "abc123"}
	out := c.Apply(in)

	if out["API_KEY"] != "abc123" {
		t.Errorf("original key missing or wrong value")
	}
	commentKey := "#API_KEY"
	if !strings.Contains(out[commentKey], "source:prod-vault") {
		t.Errorf("expected block comment key %q, got: %v", commentKey, out[commentKey])
	}
}

func TestApply_BlockComment_KeyCount(t *testing.T) {
	c, _ := New(LevelBlock, "vault")
	in := map[string]string{"X": "1", "Y": "2"}
	out := c.Apply(in)
	if len(out) != 4 {
		t.Errorf("expected 4 keys (2 original + 2 comment), got %d", len(out))
	}
}
