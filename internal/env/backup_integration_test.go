package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBackup_RoundTrip_AlwaysMode verifies that repeated Create calls
// produce distinct timestamped backup files without overwriting each other.
func TestBackup_RoundTrip_AlwaysMode(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")

	if err := os.WriteFile(src, []byte("A=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	b, err := NewBackup(BackupAlways, ".bak")
	if err != nil {
		t.Fatal(err)
	}

	p1, err := b.Create(src)
	if err != nil || p1 == "" {
		t.Fatalf("first backup failed: path=%q err=%v", p1, err)
	}

	// Update source and create second backup.
	if err := os.WriteFile(src, []byte("A=2\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	p2, err := b.Create(src)
	if err != nil || p2 == "" {
		t.Fatalf("second backup failed: path=%q err=%v", p2, err)
	}

	d1, _ := os.ReadFile(p1)
	d2, _ := os.ReadFile(p2)
	if string(d1) == string(d2) {
		t.Error("expected backup contents to differ between writes")
	}
}

// TestBackup_OnWrite_SuffixPreserved ensures the custom suffix is used.
func TestBackup_OnWrite_SuffixPreserved(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, ".env")
	if err := os.WriteFile(src, []byte("X=y\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	b, _ := NewBackup(BackupOnWrite, ".backup")
	dest, err := b.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(dest, ".backup") {
		t.Errorf("expected .backup suffix in %q", dest)
	}
}
