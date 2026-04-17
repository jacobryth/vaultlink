package rotate

import (
	"testing"
)

func TestDetect_NewKey(t *testing.T) {
	r := New()
	prev := map[string]string{}
	curr := map[string]string{"API_KEY": "secret123"}
	records := r.Detect(prev, curr)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", records[0].Key)
	}
	if records[0].PrevHash != "" {
		t.Errorf("expected empty prev hash for new key")
	}
}

func TestDetect_ChangedValue(t *testing.T) {
	r := New()
	prev := map[string]string{"DB_PASS": "old"}
	curr := map[string]string{"DB_PASS": "new"}
	records := r.Detect(prev, curr)
	if len(records) != 1 {
	expected 1 record, got %d", len(records))
	}
	if records[0].PrevHash == records[0].NewHash {
		t.Error("expected prev and new hash to differ")
	}
}

func TestDetect_UnchangedValue(t *testing.T) {
	r := New()
	prev := map[string]string{"TOKEN": "same"}
	curr := map[string]string{"TOKEN": "same"}
	records := r.Detect(prev, curr)
	if len(records) != 0 {
		t.Errorf("expected no records for unchanged value, got %d", len(records))
	}
}

func TestDetect_EmptyMaps(t *testing.T) {
	r := New()
	records := r.Detect(map[string]string{}, map[string]string{})
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestSummary_NoRotations(t *testing.T) {
	s := Summary(nil)
	if s != "no rotations detected" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithRotations(t *testing.T) {
	r := New()
	records := r.Detect(
		map[string]string{"A": "old", "B": "x"},
		map[string]string{"A": "new", "B": "x", "C": "added"},
	)
	s := Summary(records)
	if s != "2 secret(s) rotated" {
		t.Errorf("unexpected summary: %s", s)
	}
}
