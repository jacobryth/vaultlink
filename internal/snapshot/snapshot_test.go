package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestSaveAndLoad(t *testing.T) {
	m := NewManager(tempFile(t))
	s := &Snapshot{
		SecretPath: "secret/data/app",
		Keys:       []string{"DB_HOST", "DB_PASS"},
		Checksum:   map[string]string{"DB_HOST": "abc", "DB_PASS": "xyz"},
	}
	if err := m.Save(s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.SecretPath != s.SecretPath {
		t.Errorf("expected %s, got %s", s.SecretPath, loaded.SecretPath)
	}
	if len(loaded.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(loaded.Keys))
	}
}

func TestLoad_FileNotExist(t *testing.T) {
	m := NewManager("/tmp/nonexistent_snapshot_xyz.json")
	s, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestHasChanged_NoPrevious(t *testing.T) {
	m := NewManager(tempFile(t))
	changed, err := m.HasChanged(map[string]string{"KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true when no previous snapshot")
	}
}

func TestHasChanged_Identical(t *testing.T) {
	f := tempFile(t)
	m := NewManager(f)
	data := map[string]string{"A": "1", "B": "2"}
	_ = m.Save(&Snapshot{Checksum: data})
	changed, err := m.HasChanged(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false for identical data")
	}
}

func TestHasChanged_Modified(t *testing.T) {
	f := tempFile(t)
	m := NewManager(f)
	_ = m.Save(&Snapshot{Checksum: map[string]string{"A": "old"}})
	changed, err := m.HasChanged(map[string]string{"A": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true when value differs")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	m := NewManager("/nonexistent_dir/snapshot.json")
	err := m.Save(&Snapshot{})
	if err == nil {
		t.Error("expected error for invalid path")
		os.Remove("/nonexistent_dir/snapshot.json")
	}
}
