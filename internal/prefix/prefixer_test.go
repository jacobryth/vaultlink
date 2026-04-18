package prefix

import (
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	_, err := New(LevelNone, "", "")
	if err != nil {
		t.Fatalf("expected no error for none, got %v", err)
	}
	_, err = New(LevelEnv, "staging", "")
	if err != nil {
		t.Fatalf("expected no error for env, got %v", err)
	}
	_, err = New(LevelCustom, "", "MY_")
	if err != nil {
		t.Fatalf("expected no error for custom, got %v", err)
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bad", "", "")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_EnvRequiresName(t *testing.T) {
	_, err := New(LevelEnv, "", "")
	if err == nil {
		t.Fatal("expected error when env name is empty")
	}
}

func TestNew_CustomRequiresPrefix(t *testing.T) {
	_, err := New(LevelCustom, "", "")
	if err == nil {
		t.Fatal("expected error when custom prefix is empty")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	p, _ := New(LevelNone, "", "")
	if p.Apply(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApply_NoPrefix(t *testing.T) {
	p, _ := New(LevelNone, "", "")
	out := p.Apply(map[string]string{"KEY": "val"})
	if out["KEY"] != "val" {
		t.Fatalf("expected KEY unchanged, got %v", out)
	}
}

func TestApply_EnvPrefix(t *testing.T) {
	p, _ := New(LevelEnv, "prod", "")
	out := p.Apply(map[string]string{"DB_HOST": "localhost"})
	if _, ok := out["PROD_DB_HOST"]; !ok {
		t.Fatalf("expected PROD_DB_HOST in output, got %v", out)
	}
}

func TestApply_CustomPrefix(t *testing.T) {
	p, _ := New(LevelCustom, "", "APP_")
	out := p.Apply(map[string]string{"SECRET": "abc"})
	if _, ok := out["APP_SECRET"]; !ok {
		t.Fatalf("expected APP_SECRET in output, got %v", out)
	}
}
