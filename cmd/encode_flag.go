package cmd

import (
	"github.com/spf13/cobra"
)

func registerEncodeFlag(cmd *cobra.Command) {
	cmd.Flags().String("encode", "none", `Value encoding level: none, base64`)
}

func resolvedEncodeLevel(cmd *cobra.Command) string {
	v, err := cmd.Flags().GetString("encode")
	if err != nil || v == "" {
		return "none"
	}
	return v
}
