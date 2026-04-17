package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vaultlink/internal/notify"
)

// notifyLevel holds the --notify flag value.
var notifyLevel string

// registerNotifyFlag attaches the --notify flag to a command.
func registerNotifyFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&notifyLevel,
		"notify",
		string(notify.LevelSummary),
		`notification verbosity: silent | summary | verbose`,
	)
}

// resolvedNotifyLevel returns the notify.Level from the flag value,
// falling back to LevelSummary for unrecognised values.
func resolvedNotifyLevel() notify.Level {
	switch notify.Level(notifyLevel) {
	case notify.LevelSilent, notify.LevelSummary, notify.LevelVerbose:
		return notify.Level(notifyLevel)
	default:
		return notify.LevelSummary
	}
}
