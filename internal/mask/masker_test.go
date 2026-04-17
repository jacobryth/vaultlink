package mask

import (
	"testing"
)

func TestNew_DefaultsToFull(t *testing.T) {
	m := New("invalid")
	if m.Level != LevelFull {
		t.Errorf("expected full, got %s", m.Level)
	}
}

func TestApply_FullMask(t *testing.T) {
	m := New(LevelFull)
	out := m.Apply(map[string]string{"KEY": "secret"})
	if out["KEY"] != "******" {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_PartialMask(t *testing.T) {
	m := New(LevelPartial)
	out := m.Apply(map[string]string{"KEY": "secret"})
	if out["KEY"] != "se****" {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_PartialMask_ShortValue(t *testing.T) {
	m := New(LevelPartial)
	out := m.Apply(map[string]string{"K": "a"})
	if out["K"] != "*" {
		t.Errorf("unexpected value: %s", out["K"])
	}
}

func TestApply_NoMask(t *testing.T) {
	m := New(LevelNone)
	out := m.Apply(map[string]string{"KEY": "secret"})
	if out["KEY"] != "secret" {
		t.Errorf("expected plaintext, got %s", out["KEY"])
	}
}

func TestApply_NilInput(t *testing.T) {
	m := New(LevelFull)
	if m.Apply(nil) != nil {
		t.Error("expected nil output for nil input")
	}
}

func TestApply_EmptyMap(t *testing.T) {
	m := New(LevelFull)
	out := m.Apply(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty map")
	}
}
