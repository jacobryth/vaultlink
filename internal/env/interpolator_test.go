package env

import (
	"testing"
)

func TestNewInterpolator_ValidModes(t *testing.T) {
	modes := []InterpolateMode{InterpolateModeNone, InterpolateModeStrict, InterpolateModeLoose}
	for _, m := range modes {
		_, err := NewInterpolator(m)
		if err != nil {
			t.Errorf("expected no error for mode %q, got %v", m, err)
		}
	}
}

func TestNewInterpolator_InvalidMode(t *testing.T) {
	_, err := NewInterpolator("bogus")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestInterpolate_NilSecrets(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeLoose)
	out, err := ip.Apply(nil)
	if err != nil || out != nil {
		t.Fatalf("expected nil, nil; got %v, %v", out, err)
	}
}

func TestInterpolate_NoneMode_PassThrough(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeNone)
	secrets := map[string]string{"A": "${B}", "B": "hello"}
	out, err := ip.Apply(secrets)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "${B}" {
		t.Errorf("expected raw value, got %q", out["A"])
	}
}

func TestInterpolate_LooseMode_ResolvesRef(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeLoose)
	secrets := map[string]string{"HOST": "localhost", "DSN": "postgres://${HOST}/db"}
	out, err := ip.Apply(secrets)
	if err != nil {
		t.Fatal(err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("unexpected DSN: %q", out["DSN"])
	}
}

func TestInterpolate_LooseMode_LeavesUnresolved(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeLoose)
	secrets := map[string]string{"A": "${MISSING}"}
	out, err := ip.Apply(secrets)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "${MISSING}" {
		t.Errorf("expected raw ref, got %q", out["A"])
	}
}

func TestInterpolate_StrictMode_ErrorOnMissing(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeStrict)
	secrets := map[string]string{"A": "${NOPE}"}
	_, err := ip.Apply(secrets)
	if err == nil {
		t.Fatal("expected error for unresolved strict reference")
	}
}

func TestInterpolate_StrictMode_ResolvesPresent(t *testing.T) {
	ip, _ := NewInterpolator(InterpolateModeStrict)
	secrets := map[string]string{"PORT": "5432", "URL": "http://host:${PORT}"}
	out, err := ip.Apply(secrets)
	if err != nil {
		t.Fatal(err)
	}
	if out["URL"] != "http://host:5432" {
		t.Errorf("unexpected URL: %q", out["URL"])
	}
}
