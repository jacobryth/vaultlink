package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func registerInterpolateFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"interpolate",
		"none",
		`Resolve \${KEY} references in secret values. Modes: none, loose, strict`,
	)
}

func resolvedInterpolateMode(cmd *cobra.Command) env.InterpolateMode {
	v, _ := cmd.Flags().GetString("interpolate")
	if v == "" {
		return env.InterpolateModeNone
	}
	return env.InterpolateMode(v)
}
