package cmd

import (
	"github.com/spf13/cobra"
)

func registerTimestampFlags(cmd *cobra.Command) {
	cmd.Flags().String("timestamp", "none", `Stamp secret keys with current date. Levels: none, suffix, prefix`)
	cmd.Flags().String("timestamp-format", "20060102", "Go time format string used for the timestamp stamp (default: 20060102)")
}

// resolvedTimestampLevel returns the --timestamp flag value from cmd.
func resolvedTimestampLevel(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("timestamp")
	if v == "" {
		return "none"
	}
	return v
}

// resolvedTimestampFormat returns the --timestamp-format flag value from cmd.
func resolvedTimestampFormat(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("timestamp-format")
	if v == "" {
		return "20060102"
	}
	return v
}
