package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writePatchTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNewPatcher_InvalidMode(t *testing.T) {
	_, err := NewPatcher("bogus")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewPatcher_ValidModes(t *testing.T) {
	for _, m := range []PatchMode{PatchModeUpsert, PatchModeAddOnly, PatchModeRemove} {
		_, err := NewPatcher(m)
		if err != nil {
			t.Errorf("mode %q: unexpected error: %v", m, err)
		}
	}
}

func TestPatch_Upsert(t *testing.T) {
	path := writePatchTempEnv(t, "FOO=old\nBAR=keep\n")
	p, _ := NewPatcher(PatchModeUpsert)
	if err := p.Patch(path, map[string]string{"FOO": "new", "BAZ": "added"}); err != nil {
		t.Fatal(err)
	}
	m, _ := readEnvFile(path)
	if m["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %q", m["FOO"])
	}
	if m["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %q", m["BAR"])
	}
	if m["BAZ"] != "added" {
		t.Errorf("expected BAZ=added, got %q", m["BAZ"])
	}
}

func TestPatch_AddOnly(t *testing.T) {
	path := writePatchTempEnv(t, "FOO=original\n")
	p, _ := NewPatcher(PatchModeAddOnly)
	if err := p.Patch(path, map[string]string{"FOO": "ignored", "NEW": "yes"}); err != nil {
		t.Fatal(err)
	}
	m, _ := readEnvFile(path)
	if m["FOO"] != "original" {
		t.Errorf("expected FOO unchanged, got %q", m["FOO"])
	}
	if m["NEW"] != "yes" {
		t.Errorf("expected NEW=yes, got %q", m["NEW"])
	}
}

func TestPatch_Remove(t *testing.T) {
	path := writePatchTempEnv(t, "FOO=val\nBAR=val\n")
	p, _ := NewPatcher(PatchModeRemove)
	if err := p.Patch(path, map[string]string{"FOO": ""}); err != nil {
		t.Fatal(err)
	}
	m, _ := readEnvFile(path)
	if _, ok := m["FOO"]; ok {
		t.Error("expected FOO to be removed")
	}
	if m["BAR"] != "val" {
		t.Errorf("expected BAR=val, got %q", m["BAR"])
	}
}

func TestPatch_FileNotExist_Upsert(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.env")
	p, _ := NewPatcher(PatchModeUpsert)
	if err := p.Patch(path, map[string]string{"KEY": "val"}); err != nil {
		t.Fatal(err)
	}
	m, _ := readEnvFile(path)
	if m["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", m["KEY"])
	}
}
