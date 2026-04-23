package merge

import (
	"testing"
)

func TestNew_ValidStrategies(t *testing.T) {
	for _, s := range []Strategy{StrategyNone, StrategyOverwrite, StrategyKeepFirst} {
		_, err := New(s)
		if err != nil {
			t.Errorf("expected no error for strategy %q, got %v", s, err)
		}
	}
}

func TestNew_InvalidStrategy(t *testing.T) {
	_, err := New("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestApply_NilLayers(t *testing.T) {
	m, _ := New(StrategyOverwrite)
	out := m.Apply()
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestApply_None_ReturnsFirstLayer(t *testing.T) {
	m, _ := New(StrategyNone)
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "99", "C": "3"}
	out := m.Apply(a, b)
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected result for none strategy: %v", out)
	}
	if _, ok := out["C"]; ok {
		t.Error("none strategy should not include keys from second layer")
	}
}

func TestApply_Overwrite(t *testing.T) {
	m, _ := New(StrategyOverwrite)
	a := map[string]string{"KEY": "old", "ONLY_A": "yes"}
	b := map[string]string{"KEY": "new", "ONLY_B": "yes"}
	out := m.Apply(a, b)
	if out["KEY"] != "new" {
		t.Errorf("expected overwrite to win, got %q", out["KEY"])
	}
	if out["ONLY_A"] != "yes" || out["ONLY_B"] != "yes" {
		t.Errorf("expected all keys present: %v", out)
	}
}

func TestApply_KeepFirst(t *testing.T) {
	m, _ := New(StrategyKeepFirst)
	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "override", "EXTRA": "val"}
	out := m.Apply(a, b)
	if out["KEY"] != "original" {
		t.Errorf("expected keep-first to preserve original, got %q", out["KEY"])
	}
	if out["EXTRA"] != "val" {
		t.Errorf("expected new-only key to be present")
	}
}

func TestApply_None_NilFirstLayer(t *testing.T) {
	m, _ := New(StrategyNone)
	out := m.Apply(nil)
	if out == nil || len(out) != 0 {
		t.Errorf("expected empty non-nil map for nil first layer, got %v", out)
	}
}
