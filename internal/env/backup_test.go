package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewBackup_ValidModes(t *testing.T) {
	for _, mode := range []BackupMode{BackupNone, BackupAlways, BackupOnWrite} {
		_, err := NewBackup(mode, "")
		if err != nil {
			t.Errorf("mode %q: unexpected error: %v", mode, err)
		}
	}
}

func TestNewBackup_InvalidMode(t *testing.T) {
	_, err := NewBackup("unknown", "")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewBackup_DefaultSuffix(t *testing.T) {
	b, _ := NewBackup(BackupAlways, "")
	if b.suffix != ".bak" {
		t.Errorf("expected default suffix .bak, got %q", b.suffix)
	}
}

func TestShouldBackup_None(t *testing.T) {
	b, _ := NewBackup(BackupNone, "")
	if b.ShouldBackup() {
		t.Error("expected false for BackupNone")
	}
}

func TestShouldBackup_Always(t *testing.T) {
	b, _ := NewBackup(BackupAlways, "")
	if !b.ShouldBackup() {
		t.Error("expected true for BackupAlways")
	}
}

func TestCreate_None_ReturnsEmpty(t *testing.T) {
	b, _ := NewBackup(BackupNone, "")
	path, err := b.Create("/any/file")
	if err != nil || path != "" {
		t.Errorf("expected empty path and nil err, got %q, %v", path, err)
	}
}

func TestCreate_FileNotExist_NoError(t *testing.T) {
	b, _ := NewBackup(BackupAlways, ".bak")
	path, err := b.Create("/nonexistent/file.env")
	if err != nil {
		t.Fatalf("expected nil error for missing source, got %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestCreate_WritesBackup(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	if err := os.WriteFile(src, []byte("KEY=val\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	b, _ := NewBackup(BackupAlways, ".bak")
	dest, err := b.Create(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest == "" {
		t.Fatal("expected non-empty backup path")
	}
	if !strings.HasSuffix(dest, ".bak") {
		t.Errorf("backup path %q should end with .bak", dest)
	}
	data, _ := os.ReadFile(dest)
	if string(data) != "KEY=val\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}
