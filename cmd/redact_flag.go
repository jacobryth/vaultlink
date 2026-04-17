package cmd

import (
	"github.com/spf13/cobra"
)

var redactEnabled bool

// registerRedactFlag attaches the --redact flag to a command.
func registerRedactFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&redactEnabled, "redact", true, "Redact sensitive values in output (default: true)")
}
