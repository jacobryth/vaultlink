package cmd

import (
	"github.com/spf13/cobra"
)

func registerLinebreakFlag(cmd *cobra.Command) {
	cmd.Flags().String("linebreak", "none", `Normalize line endings in secret values: none, unix, windows`)
}

func resolvedLinebreakLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("linebreak")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
