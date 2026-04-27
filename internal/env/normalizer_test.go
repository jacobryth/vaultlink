package env

import (
	"testing"
)

func TestNewNormalizer_ValidLevels(t *testing.T) {
	levels := []NormalizeLevel{
		NormalizeLevelNone,
		NormalizeLevelKeys,
		NormalizeLevelValues,
		NormalizeLevelBoth,
	}
	for _, lvl := range levels {
		_, err := NewNormalizer(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNewNormalizer_InvalidLevel(t *testing.T) {
	_, err := NewNormalizer("unknown")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApplyNormalizer_NilSecrets(t *testing.T) {
	n, _ := NewNormalizer(NormalizeLevelBoth)
	result := n.Apply(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestApplyNormalizer_NoNormalization(t *testing.T) {
	n, _ := NewNormalizer(NormalizeLevelNone)
	input := map[string]string{"my-key": "  value  "}
	out := n.Apply(input)
	if out["my-key"] != "  value  " {
		t.Errorf("expected unchanged value, got %q", out["my-key"])
	}
}

func TestApplyNormalizer_NormalizeKeys(t *testing.T) {
	n, _ := NewNormalizer(NormalizeLevelKeys)
	input := map[string]string{"my-key": "  hello  ", "db.host": "localhost"}
	out := n.Apply(input)
	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected key MY_KEY")
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST")
	}
	if out["MY_KEY"] != "  hello  " {
		t.Errorf("expected value unchanged, got %q", out["MY_KEY"])
	}
}

func TestApplyNormalizer_NormalizeValues(t *testing.T) {
	n, _ := NewNormalizer(NormalizeLevelValues)
	input := map[string]string{"KEY": "  trimmed  "}
	out := n.Apply(input)
	if out["KEY"] != "trimmed" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApplyNormalizer_NormalizeBoth(t *testing.T) {
	n, _ := NewNormalizer(NormalizeLevelBoth)
	input := map[string]string{"api.secret": "  s3cr3t  "}
	out := n.Apply(input)
	if out["API_SECRET"] != "s3cr3t" {
		t.Errorf("expected API_SECRET=s3cr3t, got %v", out)
	}
}
