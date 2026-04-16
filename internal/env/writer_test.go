package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, true)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := readExisting(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	for k, v := range secrets {
		if got[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, got[k])
		}
	}
}

func TestWrite_MergeNoOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Write initial content
	if err := os.WriteFile(path, []byte("EXISTING_KEY=original\n"), 0600); err != nil {
		t.Fatal(err)
	}

	w := NewWriter(path, false)
	secrets := map[string]string{
		"EXISTING_KEY": "overridden",
		"NEW_KEY":      "newvalue",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := readExisting(path)
	if err != nil {
		t.Fatalf("failed to read merged file: %v", err)
	}

	if got["EXISTING_KEY"] != "original" {
		t.Errorf("expected EXISTING_KEY to remain %q, got %q", "original", got["EXISTING_KEY"])
	}
	if got["NEW_KEY"] != "newvalue" {
		t.Errorf("expected NEW_KEY=%q, got %q", "newvalue", got["NEW_KEY"])
	}
}

func TestWrite_OverwriteMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("OLD_KEY=oldval\n"), 0600); err != nil {
		t.Fatal(err)
	}

	w := NewWriter(path, true)
	if err := w.Write(map[string]string{"NEW_KEY": "newval"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := readExisting(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be absent in overwrite mode")
	}
}
