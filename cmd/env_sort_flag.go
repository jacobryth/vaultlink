package cmd

import (
	"github.com/spf13/cobra"
)

const defaultEnvSortMode = "none"

// registerEnvSortFlag attaches the --env-sort flag to cmd.
func registerEnvSortFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"env-sort",
		defaultEnvSortMode,
		`sort order applied to .env keys before writing (none|asc|desc)`,
	)
}

// resolvedEnvSortMode returns the effective env-sort mode from cmd flags.
// Falls back to the default when the flag is absent.
func resolvedEnvSortMode(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("env-sort")
	if err != nil || v == "" {
		return defaultEnvSortMode
	}
	return v
}
