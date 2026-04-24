package cmd

import (
	"github.com/spf13/cobra"
)

// registerSanitizeFlag attaches the --sanitize flag to the given command.
func registerSanitizeFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"sanitize",
		"none",
		`sanitize secret values before writing (none|strip|normalize).
  none      – no sanitization
  strip     – remove non-printable / control characters
  normalize – strip + trim and collapse whitespace`,
	)
}

// resolvedSanitizeLevel returns the sanitize level from the command flags.
func resolvedSanitizeLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("sanitize")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
