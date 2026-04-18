package dedupe

import (
	"testing"
)

func TestNew_ValidStrategies(t *testing.T) {
	for _, s := range []string{"first", "last", "error", ""} {
		_, err := New(s)
		if err != nil {
			t.Errorf("expected no error for strategy %q, got %v", s, err)
		}
	}
}

func TestNew_InvalidStrategy(t *testing.T) {
	_, err := New("random")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestApply_NilPairs(t *testing.T) {
	d, _ := New("first")
	out, err := d.Apply(nil)
	if err != nil || len(out) != 0 {
		t.Fatal("expected empty map and no error")
	}
}

func TestApply_First(t *testing.T) {
	d, _ := New("first")
	pairs := []KV{{"KEY", "a"}, {"KEY", "b"}}
	out, _ := d.Apply(pairs)
	if out["KEY"] != "a" {
		t.Errorf("expected 'a', got %q", out["KEY"])
	}
}

func TestApply_Last(t *testing.T) {
	d, _ := New("last")
	pairs := []KV{{"KEY", "a"}, {"KEY", "b"}}
	out, _ := d.Apply(pairs)
	if out["KEY"] != "b" {
		t.Errorf("expected 'b', got %q", out["KEY"])
	}
}

func TestApply_ErrorOnDuplicate(t *testing.T) {
	d, _ := New("error")
	pairs := []KV{{"KEY", "a"}, {"KEY", "b"}}
	_, err := d.Apply(pairs)
	if err == nil {
		t.Fatal("expected error on duplicate key")
	}
}

func TestApply_NoDuplicates(t *testing.T) {
	d, _ := New("error")
	pairs := []KV{{"A", "1"}, {"B", "2"}}
	out, err := d.Apply(pairs)
	if err != nil || out["A"] != "1" || out["B"] != "2" {
		t.Fatal("unexpected result")
	}
}
