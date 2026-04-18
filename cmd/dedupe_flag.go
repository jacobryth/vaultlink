package cmd

import (
	"github.com/spf13/cobra"
)

func registerDedupeFlag(cmd *cobra.Command) {
	cmd.Flags().String("dedupe", "first", "Duplicate key resolution strategy: first, last, error")
}

func resolvedDedupeStrategy(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("dedupe")
	if err != nil || v == "" {
		return "first"
	}
	return v
}
