package env

import (
	"testing"
)

func TestNewEnvDeduplicator_ValidModes(t *testing.T) {
	for _, mode := range []DedupeMode{DedupeModeNone, DedupeModeFirst, DedupeModeLast, DedupeModeError} {
		_, err := NewEnvDeduplicator(mode)
		if err != nil {
			t.Errorf("expected no error for mode %q, got %v", mode, err)
		}
	}
}

func TestNewEnvDeduplicator_InvalidMode(t *testing.T) {
	_, err := NewEnvDeduplicator("bogus")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestEnvDedupe_NilInput(t *testing.T) {
	d, _ := NewEnvDeduplicator(DedupeModeFirst)
	out, err := d.Apply(nil)
	if err != nil || out != nil {
		t.Fatalf("expected nil, nil; got %v, %v", out, err)
	}
}

func TestEnvDedupe_None_PassThrough(t *testing.T) {
	d, _ := NewEnvDeduplicator(DedupeModeNone)
	input := [][2]string{{"A", "1"}, {"A", "2"}}
	out, err := d.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(out))
	}
}

func TestEnvDedupe_First_KeepsFirst(t *testing.T) {
	d, _ := NewEnvDeduplicator(DedupeModeFirst)
	input := [][2]string{{"KEY", "first"}, {"OTHER", "x"}, {"KEY", "second"}}
	out, err := d.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(out))
	}
	if out[0][1] != "first" {
		t.Errorf("expected 'first', got %q", out[0][1])
	}
}

func TestEnvDedupe_Last_KeepsLast(t *testing.T) {
	d, _ := NewEnvDeduplicator(DedupeModeLast)
	input := [][2]string{{"KEY", "first"}, {"KEY", "second"}, {"OTHER", "y"}}
	out, err := d.Apply(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(out))
	}
	if out[0][1] != "second" {
		t.Errorf("expected 'second', got %q", out[0][1])
	}
}

func TestEnvDedupe_Error_ReturnsDuplicateError(t *testing.T) {
	d, _ := NewEnvDeduplicator(DedupeModeError)
	input := [][2]string{{"KEY", "v1"}, {"KEY", "v2"}}
	_, err := d.Apply(input)
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
}
