package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func TestNewLogger_Stdout(t *testing.T) {
	l, err := NewLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.f != os.Stdout {
		t.Error("expected stdout")
	}
	if err := l.Close(); err != nil {
		t.Errorf("close: %v", err)
	}
}

func TestNewLogger_File(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "audit-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	l, err := NewLogger(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Close()

	keys := []string{"DB_HOST", "DB_PASS"}
	if err := l.Log("sync", "secret/app", "backend", keys, true, ""); err != nil {
		t.Fatalf("log: %v", err)
	}
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected one line in log file")
	}
	var e Entry
	if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Operation != "sync" {
		t.Errorf("got op %q, want sync", e.Operation)
	}
	if e.Role != "backend" {
		t.Errorf("got role %q, want backend", e.Role)
	}
	if len(e.Keys) != 2 {
		t.Errorf("got %d keys, want 2", len(e.Keys))
	}
	if !e.Success {
		t.Error("expected success=true")
	}
}

func TestLog_InvalidPath(t *testing.T) {
	_, err := NewLogger("/nonexistent-dir/audit.log")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
