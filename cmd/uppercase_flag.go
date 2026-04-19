package cmd

import (
	"github.com/spf13/cobra"
)

func registerUppercaseFlag(cmd *cobra.Command) {
	cmd.Flags().String("uppercase", "none", `Uppercase transformation level: none, keys, values, both`)
}

func resolvedUppercaseLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("uppercase")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
