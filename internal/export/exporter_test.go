package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []string{"env", "json", "ENV", "JSON"} {
		_, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %s, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_EnvFormat(t *testing.T) {
	e, _ := New("env")
	dest := filepath.Join(t.TempDir(), ".env.export")
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := e.Write(secrets, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	content := string(data)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", content)
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", content)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	e, _ := New("json")
	dest := filepath.Join(t.TempDir(), "secrets.json")
	secrets := map[string]string{"KEY": "value"}
	if err := e.Write(secrets, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	var out map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %v", out)
	}
}

func TestWrite_CreatesNestedDirs(t *testing.T) {
	e, _ := New("env")
	dest := filepath.Join(t.TempDir(), "a", "b", "c", ".env")
	if err := e.Write(map[string]string{"X": "1"}, dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
