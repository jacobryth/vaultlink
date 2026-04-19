package cmd

import (
	"github.com/spf13/cobra"
)

func registerTypecastFlag(cmd *cobra.Command) {
	cmd.Flags().String("typecast", "none", `Value type normalization level: none, string, infer`)
}

func resolvedTypecastLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("typecast")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
