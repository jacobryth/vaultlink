package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_BasicPairs(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	l := NewLoader(p)
	got, err := l.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", got)
	}
}

func TestLoad_SkipsCommentAndBlank(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nKEY=value\n")
	got, err := NewLoader(p).Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["KEY"] != "value" {
		t.Errorf("expected single key 'KEY', got: %v", got)
	}
}

func TestLoad_StripInlineComment(t *testing.T) {
	p := writeTempEnv(t, "SECRET=abc123 # this is a secret\n")
	got, err := NewLoader(p).Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SECRET"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", got["SECRET"])
	}
}

func TestLoad_StripQuotes(t *testing.T) {
	p := writeTempEnv(t, `DB_URL="postgres://localhost/db"` + "\n")
	got, err := NewLoader(p).Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_URL"] != "postgres://localhost/db" {
		t.Errorf("unexpected value: %q", got["DB_URL"])
	}
}

func TestLoad_FileNotExist_ReturnsEmpty(t *testing.T) {
	l := NewLoader("/nonexistent/.env")
	got, err := l.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got: %v", got)
	}
}

func TestLoad_MalformedLine_ReturnsError(t *testing.T) {
	p := writeTempEnv(t, "NOEQUALS\n")
	_, err := NewLoader(p).Load()
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestKeys_ReturnsList(t *testing.T) {
	p := writeTempEnv(t, "A=1\nB=2\nC=3\n")
	keys, err := NewLoader(p).Keys()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(keys), keys)
	}
}
