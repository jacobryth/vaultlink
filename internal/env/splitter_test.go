package env

import (
	"testing"
)

func TestNewSplitter_ValidModes(t *testing.T) {
	modes := []struct {
		mode  SplitMode
		delim string
	}{
		{SplitNone, ""},
		{SplitPrefix, "_"},
		{SplitAlpha, ""},
	}
	for _, tc := range modes {
		_, err := NewSplitter(tc.mode, tc.delim)
		if err != nil {
			t.Errorf("mode %q: unexpected error: %v", tc.mode, err)
		}
	}
}

func TestNewSplitter_InvalidMode(t *testing.T) {
	_, err := NewSplitter("bogus", "")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewSplitter_PrefixRequiresDelim(t *testing.T) {
	_, err := NewSplitter(SplitPrefix, "")
	if err == nil {
		t.Fatal("expected error when prefix mode has empty delimiter")
	}
}

func TestSplit_NilSecrets(t *testing.T) {
	s, _ := NewSplitter(SplitNone, "")
	out := s.Split(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %v", out)
	}
}

func TestSplit_NoneMode(t *testing.T) {
	s, _ := NewSplitter(SplitNone, "")
	out := s.Split(map[string]string{"FOO": "1", "BAR": "2"})
	if len(out) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(out))
	}
	if len(out["default"]) != 2 {
		t.Fatalf("expected 2 entries in default bucket")
	}
}

func TestSplit_PrefixMode(t *testing.T) {
	s, _ := NewSplitter(SplitPrefix, "_")
	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
		"NOPREFIX": "val",
	}
	out := s.Split(secrets)
	if len(out["app"]) != 2 {
		t.Errorf("expected 2 in app bucket, got %d", len(out["app"]))
	}
	if len(out["db"]) != 1 {
		t.Errorf("expected 1 in db bucket, got %d", len(out["db"]))
	}
	if _, ok := out["default"]["NOPREFIX"]; !ok {
		t.Error("expected NOPREFIX in default bucket")
	}
}

func TestSplit_AlphaMode(t *testing.T) {
	s, _ := NewSplitter(SplitAlpha, "")
	secrets := map[string]string{
		"ALPHA": "a",
		"BETA":  "b",
		"BRACE": "c",
	}
	out := s.Split(secrets)
	if len(out["a"]) != 1 {
		t.Errorf("expected 1 in 'a' bucket, got %d", len(out["a"]))
	}
	if len(out["b"]) != 2 {
		t.Errorf("expected 2 in 'b' bucket, got %d", len(out["b"]))
	}
}
