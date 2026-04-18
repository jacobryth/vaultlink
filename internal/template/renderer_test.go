package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_ReplacesPlaceholders(t *testing.T) {
	r := New(false)
	out, err := r.Render("host={{DB_HOST}} port={{DB_PORT}}", map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "host=localhost port=5432" {
		t.Errorf("got %q", out)
	}
}

func TestRender_MissingKey_NonStrict(t *testing.T) {
	r := New(false)
	out, err := r.Render("val={{MISSING}}", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "val={{MISSING}}" {
		t.Errorf("got %q", out)
	}
}

func TestRender_MissingKey_Strict(t *testing.T) {
	r := New(true)
	_, err := r.Render("val={{MISSING}}", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unresolved placeholder")
	}
}

func TestRender_EmptySecrets(t *testing.T) {
	r := New(false)
	out, err := r.Render("static text", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "static text" {
		t.Errorf("got %q", out)
	}
}

func TestRenderFile_Success(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	_ = os.WriteFile(p, []byte("token={{TOKEN}}"), 0600)

	r := New(true)
	out, err := r.RenderFile(p, map[string]string{"TOKEN": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "token=abc123" {
		t.Errorf("got %q", out)
	}
}

func TestRenderFile_NotFound(t *testing.T) {
	r := New(false)
	_, err := r.RenderFile("/nonexistent/file.tmpl", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
