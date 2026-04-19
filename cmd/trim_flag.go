package cmd

import (
	"github.com/spf13/cobra"
)

func registerTrimFlag(cmd *cobra.Command) {
	cmd.Flags().String("trim", "none", `Trim whitespace from secret values. Levels: none, space, all`)
}

func resolvedTrimLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("trim")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
