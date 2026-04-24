package env

import (
	"strings"
	"testing"
)

func TestNewFormatter_ValidLevels(t *testing.T) {
	levels := []FormatLevel{FormatNone, FormatExport, FormatInlined}
	for _, lvl := range levels {
		_, err := NewFormatter(lvl)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNewFormatter_InvalidLevel(t *testing.T) {
	_, err := NewFormatter("fancy")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
	if !strings.Contains(err.Error(), "unknown level") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApply_NilSecrets(t *testing.T) {
	f, _ := NewFormatter(FormatNone)
	result := f.Apply(nil)
	if result != nil {
		t.Errorf("expected nil for nil secrets, got %v", result)
	}
}

func TestApply_NoFormat(t *testing.T) {
	f, _ := NewFormatter(FormatNone)
	secrets := map[string]string{"KEY": "value"}
	lines := f.Apply(secrets)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "KEY=value" {
		t.Errorf("expected KEY=value, got %q", lines[0])
	}
}

func TestApply_ExportFormat(t *testing.T) {
	f, _ := NewFormatter(FormatExport)
	secrets := map[string]string{"TOKEN": "abc123"}
	lines := f.Apply(secrets)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != "export TOKEN=abc123" {
		t.Errorf("expected 'export TOKEN=abc123', got %q", lines[0])
	}
}

func TestApply_InlinedFormat(t *testing.T) {
	f, _ := NewFormatter(FormatInlined)
	secrets := map[string]string{"A": "1", "B": "2"}
	lines := f.Apply(secrets)
	if len(lines) != 1 {
		t.Fatalf("expected 1 inlined line, got %d", len(lines))
	}
	// Both pairs should appear in the single line separated by "; "
	if !strings.Contains(lines[0], ";") {
		t.Errorf("expected semicolon separator in inlined output, got %q", lines[0])
	}
}
