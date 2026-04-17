package notify

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultlink/internal/diff"
)

func result(added, removed, updated []string) diff.Result {
	return diff.Result{Added: added, Removed: removed, Updated: updated}
}

func TestNotify_Silent(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, LevelSilent)
	n.Notify(result([]string{"KEY"}, nil, nil))
	if buf.Len() != 0 {
		t.Errorf("expected no output in silent mode, got %q", buf.String())
	}
}

func TestNotify_Summary(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, LevelSummary)
	n.Notify(result([]string{"A"}, []string{"B"}, nil))
	out := buf.String()
	if !strings.Contains(out, "added=1") || !strings.Contains(out, "removed=1") {
		t.Errorf("unexpected summary output: %q", out)
	}
	// Should NOT contain per-key detail lines.
	if strings.Contains(out, "added:") {
		t.Errorf("summary mode should not print per-key details")
	}
}

func TestNotify_Verbose(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, LevelVerbose)
	n.Notify(result([]string{"NEW_KEY"}, nil, []string{"EXISTING"}))
	out := buf.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("verbose mode should list added keys, got %q", out)
	}
	if !strings.Contains(out, "EXISTING") {
		t.Errorf("verbose mode should list updated keys, got %q", out)
	}
}

func TestNotify_DefaultsToStdout(t *testing.T) {
	// Just ensure no panic when out is nil.
	n := New(nil, LevelSilent)
	n.Notify(result(nil, nil, nil))
}

func TestNotify_EmptyLevel_DefaultsSummary(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf, "")
	if n.level != LevelSummary {
		t.Errorf("expected default level %q, got %q", LevelSummary, n.level)
	}
}
