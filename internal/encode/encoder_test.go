package encode

import (
	"encoding/base64"
	"testing"
)

func TestNew_ValidLevels(t *testing.T) {
	for _, l := range []Level{LevelNone, LevelBase64} {
		_, err := New(l)
		if err != nil {
			t.Errorf("expected no error for level %q, got %v", l, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("rot13")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	e, _ := New(LevelBase64)
	if e.Apply(nil) != nil {
		t.Fatal("expected nil output for nil input")
	}
}

func TestApply_NoEncoding(t *testing.T) {
	e, _ := New(LevelNone)
	in := map[string]string{"KEY": "value"}
	out := e.Apply(in)
	if out["KEY"] != "value" {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
}

func TestApply_Base64(t *testing.T) {
	e, _ := New(LevelBase64)
	in := map[string]string{"SECRET": "hunter2", "EMPTY": ""}
	out := e.Apply(in)

	want := base64.StdEncoding.EncodeToString([]byte("hunter2"))
	if out["SECRET"] != want {
		t.Errorf("expected %q, got %q", want, out["SECRET"])
	}
	wantEmpty := base64.StdEncoding.EncodeToString([]byte(""))
	if out["EMPTY"] != wantEmpty {
		t.Errorf("expected %q, got %q", wantEmpty, out["EMPTY"])
	}
}

func TestDecode_RoundTrip(t *testing.T) {
	original := "super-secret-value"
	e, _ := New(LevelBase64)
	encoded := e.Apply(map[string]string{"K": original})["K"]
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decoded != original {
		t.Errorf("expected %q, got %q", original, decoded)
	}
}

func TestDecode_InvalidInput(t *testing.T) {
	_, err := Decode("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}
