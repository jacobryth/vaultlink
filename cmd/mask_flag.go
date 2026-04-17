package cmd

import (
	"github.com/spf13/cobra"
	"vaultlink/internal/mask"
)

// registerMaskFlag adds --mask-level to the given command.
func registerMaskFlag(cmd *cobra.Command) {
	cmd.Flags().String("mask-level", "full", "Secret masking level for output: full, partial, none")
}

// resolvedMaskLevel returns the mask.Level from the command flag.
func resolvedMaskLevel(cmd *cobra.Command) mask.Level {
	v, err := cmd.Flags().GetString("mask-level")
	if err != nil || v == "" {
		return mask.LevelFull
	}
	return mask.Level(v)
}
