package cmd

import (
	"github.com/spf13/cobra"
)

const defaultRotateLog = ""

// registerRotateFlags adds rotation-related flags to the given command.
func registerRotateFlags(cmd *cobra.Command) {
	cmd.Flags().String("rotate-log", defaultRotateLog,
		"path to append rotation records (JSON lines); leave empty to print to stdout")
}

// resolvedRotateLog returns the rotate-log flag value from the command.
// If the flag is not set or an error occurs, it returns the default value.
func resolvedRotateLog(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("rotate-log")
	if err != nil {
		return defaultRotateLog
	}
	return v
}
