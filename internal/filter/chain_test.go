package filter

import (
	"strings"
	"testing"
)

func TestNewChain_NoSteps(t *testing.T) {
	c := NewChain()
	if c.Len() != 0 {
		t.Fatalf("expected 0 steps, got %d", c.Len())
	}
}

func TestApply_NilSecrets(t *testing.T) {
	c := NewChain()
	out := c.Apply(nil)
	if out == nil || len(out) != 0 {
		t.Fatalf("expected empty map for nil input, got %v", out)
	}
}

func TestApply_NoSteps_PassThrough(t *testing.T) {
	c := NewChain()
	in := map[string]string{"KEY": "val", "OTHER": "x"}
	out := c.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_SingleStep_Filters(t *testing.T) {
	onlyUpper := func(s map[string]string) map[string]string {
		result := map[string]string{}
		for k, v := range s {
			if k == strings.ToUpper(k) {
				result[k] = v
			}
		}
		return result
	}
	c := NewChain(onlyUpper)
	in := map[string]string{"UPPER": "a", "lower": "b", "Mixed": "c"}
	out := c.Apply(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["UPPER"]; !ok {
		t.Error("expected UPPER key in output")
	}
}

func TestApply_MultipleSteps_Composed(t *testing.T) {
	keepA := func(s map[string]string) map[string]string {
		r := map[string]string{}
		for k, v := range s {
			if strings.HasPrefix(k, "A") {
				r[k] = v
			}
		}
		return r
	}
	keepLong := func(s map[string]string) map[string]string {
		r := map[string]string{}
		for k, v := range s {
			if len(v) > 3 {
				r[k] = v
			}
		}
		return r
	}
	c := NewChain(keepA, keepLong)
	in := map[string]string{
		"ALPHA": "longvalue",
		"APPLE": "no",
		"BETA":  "longvalue",
	}
	out := c.Apply(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(out), out)
	}
	if _, ok := out["ALPHA"]; !ok {
		t.Error("expected ALPHA in output")
	}
}

func TestApply_StepReturnsEmpty_ShortCircuits(t *testing.T) {
	calls := 0
	emptyStep := func(s map[string]string) map[string]string {
		return map[string]string{}
	}
	countStep := func(s map[string]string) map[string]string {
		calls++
		return s
	}
	c := NewChain(emptyStep, countStep)
	out := c.Apply(map[string]string{"K": "v"})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
	if calls != 0 {
		t.Errorf("expected countStep not to be called, but was called %d time(s)", calls)
	}
}

func TestChain_Len(t *testing.T) {
	f := func(s map[string]string) map[string]string { return s }
	c := NewChain(f, f, f)
	if c.Len() != 3 {
		t.Fatalf("expected 3, got %d", c.Len())
	}
}
