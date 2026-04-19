package group

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{None, Prefix, Custom} {
		rules := map[string]string{"APP_": "app"}
		_, err := New(lvl, "_", rules)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bogus", "_", nil)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_CustomRequiresRules(t *testing.T) {
	_, err := New(Custom, "_", nil)
	if err == nil {
		t.Fatal("expected error when custom level has no rules")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	g, _ := New(None, "", nil)
	out := g.Apply(nil)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}

func TestApply_NoGroup(t *testing.T) {
	g, _ := New(None, "", nil)
	out := g.Apply(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if len(out["default"]) != 2 {
		t.Errorf("expected 2 keys in default bucket, got %d", len(out["default"]))
	}
}

func TestApply_PrefixGroup(t *testing.T) {
	g, _ := New(Prefix, "_", nil)
	secrets := map[string]string{
		"APP_SECRET": "s1",
		"APP_KEY":    "s2",
		"DB_HOST":    "localhost",
		"NOPREFIX":   "val",
	}
	out := g.Apply(secrets)
	if len(out["APP"]) != 2 {
		t.Errorf("expected 2 in APP group, got %d", len(out["APP"]))
	}
	if len(out["DB"]) != 1 {
		t.Errorf("expected 1 in DB group, got %d", len(out["DB"]))
	}
	if len(out["default"]) != 1 {
		t.Errorf("expected 1 in default group, got %d", len(out["default"]))
	}
}

func TestApply_CustomGroup(t *testing.T) {
	rules := map[string]string{"DB_": "database", "APP_": "application"}
	g, _ := New(Custom, "", rules)
	secrets := map[string]string{
		"DB_HOST":   "localhost",
		"APP_TOKEN": "abc",
		"OTHER":     "xyz",
	}
	out := g.Apply(secrets)
	if out["database"]["DB_HOST"] != "localhost" {
		t.Error("expected DB_HOST in database group")
	}
	if out["application"]["APP_TOKEN"] != "abc" {
		t.Error("expected APP_TOKEN in application group")
	}
	if out["default"]["OTHER"] != "xyz" {
		t.Error("expected OTHER in default group")
	}
}
