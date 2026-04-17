package notify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vaultlink/internal/diff"
)

// Level controls verbosity of notifications.
type Level string

const (
	LevelSilent  Level = "silent"
	LevelSummary Level = "summary"
	LevelVerbose Level = "verbose"
)

// Notifier writes sync change notifications to an output.
type Notifier struct {
	out   io.Writer
	level Level
}

// New creates a Notifier writing to out at the given level.
// If out is nil, os.Stdout is used.
func New(out io.Writer, level Level) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	if level == "" {
		level = LevelSummary
	}
	return &Notifier{out: out, level: level}
}

// Notify prints change information based on the configured level.
func (n *Notifier) Notify(result diff.Result) {
	if n.level == LevelSilent {
		return
	}

	summary := diff.Summary(result)
	if n.level == LevelSummary {
		fmt.Fprintln(n.out, summary)
		return
	}

	// Verbose: print summary then per-key details.
	fmt.Fprintln(n.out, summary)
	if len(result.Added) > 0 {
		fmt.Fprintf(n.out, "  added:   %s\n", strings.Join(result.Added, ", "))
	}
	if len(result.Removed) > 0 {
		fmt.Fprintf(n.out, "  removed: %s\n", strings.Join(result.Removed, ", "))
	}
	if len(result.Updated) > 0 {
		fmt.Fprintf(n.out, "  updated: %s\n", strings.Join(result.Updated, ", "))
	}
}
