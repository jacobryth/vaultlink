package timestamp

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []Level{None, Suffix, Prefix} {
		_, err := New(lvl, "")
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("bogus", "")
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	s, _ := New(Suffix, "")
	if got := s.Apply(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestApply_NoStamp(t *testing.T) {
	s, _ := New(None, "")
	input := map[string]string{"KEY": "val"}
	out := s.Apply(input)
	if _, ok := out["KEY"]; !ok {
		t.Error("expected KEY to be unchanged")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApply_SuffixStamp(t *testing.T) {
	s, _ := New(Suffix, "20060102")
	input := map[string]string{"DB_PASS": "secret"}
	out := s.Apply(input)

	stamp := time.Now().UTC().Format("20060102")
	expected := fmt.Sprintf("DB_PASS_%s", stamp)
	if _, ok := out[expected]; !ok {
		t.Errorf("expected key %q in output, got keys: %v", expected, keysOf(out))
	}
}

func TestApply_PrefixStamp(t *testing.T) {
	s, _ := New(Prefix, "20060102")
	input := map[string]string{"API_KEY": "abc"}
	out := s.Apply(input)

	stamp := time.Now().UTC().Format("20060102")
	for k := range out {
		if !strings.HasPrefix(k, stamp+"_") {
			t.Errorf("expected key to start with %q, got %q", stamp+"_", k)
		}
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	s, _ := New(Suffix, "")
	input := map[string]string{"X": "1"}
	s.Apply(input)
	if _, ok := input["X"]; !ok {
		t.Error("original map should not be modified")
	}
}

func keysOf(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
