package rotate

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestWriteLog_EmptyRecords(t *testing.T) {
	err := WriteLog("/tmp/should_not_be_created_vaultlink.log", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// file should not exist
	if _, statErr := os.Stat("/tmp/should_not_be_created_vaultlink.log"); !os.IsNotExist(statErr) {
		os.Remove("/tmp/should_not_be_created_vaultlink.log")
	}
}

func TestWriteLog_CreatesFile(t *testing.T) {
	f, err := os.CreateTemp("", "rotate_log_*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	records := []Record{
		{Key: "SECRET", RotatedAt: time.Now().UTC(), PrevHash: "aaa", NewHash: "bbb"},
	}
	if err := WriteLog(f.Name(), records); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, _ := os.Open(f.Name())
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var entries []LogEntry
	for scanner.Scan() {
		var e LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Fatalf("invalid JSON line: %v", err)
		}
		entries = append(entries, e)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", entries[0].Key)
	}
}

func TestWriteLog_InvalidPath(t *testing.T) {
	records := []Record{
		{Key: "X", RotatedAt: time.Now().UTC(), PrevHash: "", NewHash: "abc"},
	}
	err := WriteLog("/nonexistent_dir/rotate.log", records)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
