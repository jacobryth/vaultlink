package checksum

import (
	"strings"
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{None, MD5, SHA256} {
		_, err := New(lvl)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("crc32")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
	if !strings.Contains(err.Error(), "crc32") {
		t.Errorf("error should mention the bad level, got: %v", err)
	}
}

func TestCompute_None(t *testing.T) {
	c, _ := New(None)
	if got := c.Compute(map[string]string{"KEY": "val"}); got != "" {
		t.Errorf("expected empty string for None level, got %q", got)
	}
}

func TestCompute_NilSecrets(t *testing.T) {
	for _, lvl := range []Level{MD5, SHA256} {
		c, _ := New(lvl)
		if got := c.Compute(nil); got != "" {
			t.Errorf("level %q: expected empty string for nil secrets, got %q", lvl, got)
		}
	}
}

func TestCompute_MD5_Deterministic(t *testing.T) {
	c, _ := New(MD5)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	a := c.Compute(secrets)
	b := c.Compute(secrets)
	if a != b {
		t.Errorf("MD5 not deterministic: %q != %q", a, b)
	}
	if len(a) != 32 {
		t.Errorf("expected 32-char MD5 hex, got len %d", len(a))
	}
}

func TestCompute_SHA256_Deterministic(t *testing.T) {
	c, _ := New(SHA256)
	secrets := map[string]string{"ALPHA": "1", "BETA": "2"}
	a := c.Compute(secrets)
	b := c.Compute(secrets)
	if a != b {
		t.Errorf("SHA256 not deterministic: %q != %q", a, b)
	}
	if len(a) != 64 {
		t.Errorf("expected 64-char SHA256 hex, got len %d", len(a))
	}
}

func TestCompute_DifferentInputs_DifferentHashes(t *testing.T) {
	c, _ := New(SHA256)
	a := c.Compute(map[string]string{"KEY": "value1"})
	b := c.Compute(map[string]string{"KEY": "value2"})
	if a == b {
		t.Error("different inputs produced the same hash")
	}
}
