package defaults

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	levels := []Level{None, Missing, Empty, Both}
	for _, lvl := range levels {
		rules := map[string]string{"KEY": "val"}
		if lvl == None {
			rules = nil
		}
		_, err := New(lvl, rules)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("unknown", map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestNew_LevelRequiresRules(t *testing.T) {
	for _, lvl := range []Level{Missing, Empty, Both} {
		_, err := New(lvl, nil)
		if err == nil {
			t.Errorf("expected error for level %q with no rules", lvl)
		}
	}
}

func TestApply_NilSecrets(t *testing.T) {
	d, _ := New(None, nil)
	if d.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoDefault(t *testing.T) {
	d, _ := New(None, nil)
	secrets := map[string]string{"A": "1"}
	out := d.Apply(secrets)
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestApply_MissingLevel(t *testing.T) {
	d, _ := New(Missing, map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	secrets := map[string]string{"DB_HOST": "prod-host"}
	out := d.Apply(secrets)
	if out["DB_HOST"] != "prod-host" {
		t.Errorf("expected existing key untouched, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected missing key filled, got %q", out["DB_PORT"])
	}
}

func TestApply_EmptyLevel(t *testing.T) {
	d, _ := New(Empty, map[string]string{"TOKEN": "default-token"})
	secrets := map[string]string{"TOKEN": ""}
	out := d.Apply(secrets)
	if out["TOKEN"] != "default-token" {
		t.Errorf("expected empty key filled, got %q", out["TOKEN"])
	}
}

func TestApply_BothLevel(t *testing.T) {
	d, _ := New(Both, map[string]string{"X": "default-x", "Y": "default-y"})
	secrets := map[string]string{"X": ""}
	out := d.Apply(secrets)
	if out["X"] != "default-x" {
		t.Errorf("expected empty key filled, got %q", out["X"])
	}
	if out["Y"] != "default-y" {
		t.Errorf("expected missing key filled, got %q", out["Y"])
	}
}
