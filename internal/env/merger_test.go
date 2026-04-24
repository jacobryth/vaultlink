package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
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

func TestNewEnvMerger_InvalidMode(t *testing.T) {
	_, err := NewEnvMerger("bogus")
	if err == nil {
		t.Fatal("expected error for unknown merge mode")
	}
}

func TestMergeFiles_Overwrite(t *testing.T) {
	a := writeTempEnvFile(t, "KEY=first\nONLY_A=yes\n")
	b := writeTempEnvFile(t, "KEY=second\nONLY_B=yes\n")

	m, _ := NewEnvMerger(MergeModeOverwrite)
	result, err := m.MergeFiles([]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "second" {
		t.Errorf("expected KEY=second, got %q", result["KEY"])
	}
	if result["ONLY_A"] != "yes" || result["ONLY_B"] != "yes" {
		t.Errorf("expected both unique keys present")
	}
}

func TestMergeFiles_Keep(t *testing.T) {
	a := writeTempEnvFile(t, "KEY=first\n")
	b := writeTempEnvFile(t, "KEY=second\n")

	m, _ := NewEnvMerger(MergeModeKeep)
	result, err := m.MergeFiles([]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "first" {
		t.Errorf("expected KEY=first (keep mode), got %q", result["KEY"])
	}
}

func TestMergeFiles_ErrorOnDuplicate(t *testing.T) {
	a := writeTempEnvFile(t, "KEY=first\n")
	b := writeTempEnvFile(t, "KEY=second\n")

	m, _ := NewEnvMerger(MergeModeError)
	_, err := m.MergeFiles([]string{a, b})
	if err == nil {
		t.Fatal("expected error on duplicate key")
	}
}

func TestMergeFiles_MissingFile(t *testing.T) {
	m, _ := NewEnvMerger(MergeModeOverwrite)
	_, err := m.MergeFiles([]string{filepath.Join(t.TempDir(), "nonexistent.env")})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMergeFiles_EmptyList(t *testing.T) {
	m, _ := NewEnvMerger(MergeModeOverwrite)
	result, err := m.MergeFiles([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
