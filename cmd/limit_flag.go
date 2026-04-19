package cmd

import (
	"github.com/spf13/cobra"
)

func registerLimitFlags(cmd *cobra.Command) {
	cmd.Flags().String("limit", "none", "Limit secrets returned: none, first, last")
	cmd.Flags().Int("limit-count", 10, "Number of secrets to return when limit is first or last")
}

func resolvedLimitLevel(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("limit")
	if v == "" {
		return "none"
	}
	return v
}

func resolvedLimitCount(cmd *cobra.Command) int {
	v, _ := cmd.Flags().GetInt("limit-count")
	if v <= 0 {
		return 10
	}
	return v
}
