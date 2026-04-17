package notify

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultlink/internal/diff"
	"github.com/user/vaultlink/internal/redact"
)

// Level controls how much detail is printed.
type Level string

const (
	Silent  Level = "silent"
	Summary Level = "summary"
	Verbose Level = "verbose"
)

// Notifier prints sync results to an output writer.
type Notifier struct {
	level  Level
	output io.Writer
	rule   redact.Rule
}

// New creates a Notifier with the given level. Defaults to os.Stdout.
func New(level Level, output io.Writer) *Notifier {
	if output == nil {
		output = os.Stdout
	}
	return &Notifier{level: level, output: output, rule: redact.DefaultRule()}
}

// Notify prints the diff result according to the configured level.
func (n *Notifier) Notify(result diff.Result) {
	switch n.level {
	case Silent:
		return
	case Summary:
		fmt.Fprintln(n.output, diff.Summary(result))
	case Verbose:
		fmt.Fprintln(n.output, diff.Summary(result))
		for _, k := range result.Added {
			fmt.Fprintf(n.output, "  [+] %s\n", k)
		}
		for _, k := range result.Removed {
			fmt.Fprintf(n.output, "  [-] %s\n", k)
		}
		for _, k := range result.Updated {
			fmt.Fprintf(n.output, "  [~] %s\n", k)
		}
	default:
		fmt.Fprintln(n.output, diff.Summary(result))
	}
}
