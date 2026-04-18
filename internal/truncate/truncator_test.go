package truncate

import (
	"testing"
)

func TestNew_DefaultsToNone(t *testing.T) {
	tr := New("unknown")
	if tr.level != LevelNone {
		t.Errorf("expected LevelNone, got %s", tr.level)
	}
}

func TestNew_KnownLevels(t *testing.T) {
	for _, l := range []Level{LevelShort, LevelTiny} {
		tr := New(l)
		if tr.level != l {
			t.Errorf("expected %s, got %s", l, tr.level)
		}
	}
}

func TestApply_NilSecrets(t *testing.T) {
	tr := New(LevelShort)
	if tr.Apply(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestApply_NoTruncation(t *testing.T) {
	tr := New(LevelNone)
	in := map[string]string{"KEY": "supersecretvalue"}
	out := tr.Apply(in)
	if out["KEY"] != "supersecretvalue" {
		t.Errorf("unexpected truncation: %s", out["KEY"])
	}
}

func TestApply_ShortTruncation(t *testing.T) {
	tr := New(LevelShort)
	in := map[string]string{"KEY": "supersecretvalue"}
	out := tr.Apply(in)
	if out["KEY"] != "supersecr..." {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_TinyTruncation(t *testing.T) {
	tr := New(LevelTiny)
	in := map[string]string{"KEY": "supersecretvalue"}
	out := tr.Apply(in)
	if out["KEY"] != "sup..." {
		t.Errorf("unexpected value: %s", out["KEY"])
	}
}

func TestApply_ShortValue_NoTruncation(t *testing.T) {
	tr := New(LevelShort)
	in := map[string]string{"KEY": "abc"}
	out := tr.Apply(in)
	if out["KEY"] != "abc" {
		t.Errorf("short value should not be truncated, got: %s", out["KEY"])
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	tr := New(LevelShort)
	out := tr.Apply(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty map")
	}
}
