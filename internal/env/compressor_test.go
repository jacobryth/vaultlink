package env

import (
	"testing"
)

func TestNewEnvCompressor_ValidModes(t *testing.T) {
	modes := []CompressMode{
		CompressModeNone,
		CompressModeKeys,
		CompressModeValues,
		CompressModeBoth,
	}
	for _, m := range modes {
		_, err := NewEnvCompressor(m)
		if err != nil {
			t.Errorf("expected no error for mode %q, got %v", m, err)
		}
	}
}

func TestNewEnvCompressor_InvalidMode(t *testing.T) {
	_, err := NewEnvCompressor("unknown")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestEnvCompress_NilSecrets(t *testing.T) {
	c, _ := NewEnvCompressor(CompressModeBoth)
	result := c.Compress(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestEnvCompress_NoneMode_PassThrough(t *testing.T) {
	c, _ := NewEnvCompressor(CompressModeNone)
	input := map[string]string{"KEY": "", "OTHER": "val"}
	out := c.Compress(input)
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestEnvCompress_Keys_RemovesBlankKey(t *testing.T) {
	c, _ := NewEnvCompressor(CompressModeKeys)
	input := map[string]string{"VALID": "yes", "   ": "blank-key"}
	out := c.Compress(input)
	if _, ok := out["   "]; ok {
		t.Error("expected blank key to be removed")
	}
	if out["VALID"] != "yes" {
		t.Error("expected VALID key to be preserved")
	}
}

func TestEnvCompress_Values_RemovesEmptyValue(t *testing.T) {
	c, _ := NewEnvCompressor(CompressModeValues)
	input := map[string]string{"EMPTY": "", "FILLED": "data"}
	out := c.Compress(input)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY key to be removed")
	}
	if out["FILLED"] != "data" {
		t.Error("expected FILLED key to be preserved")
	}
}

func TestEnvCompress_Both_RemovesEither(t *testing.T) {
	c, _ := NewEnvCompressor(CompressModeBoth)
	input := map[string]string{
		"GOOD":  "value",
		"EMPTY": "",
		"  ":    "blank-key",
	}
	out := c.Compress(input)
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if out["GOOD"] != "value" {
		t.Error("expected GOOD to be preserved")
	}
}
