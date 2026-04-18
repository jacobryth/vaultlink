package cmd

import (
	"github.com/spf13/cobra"
)

func registerSortFlag(cmd *cobra.Command) {
	cmd.Flags().String("sort", "none", "Sort secrets by key: none, asc, desc")
}

func resolvedSortLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("sort")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
