package cmd

import (
	"github.com/spf13/cobra"
	"vaultlink/internal/truncate"
)

func registerTruncateFlag(cmd *cobra.Command) {
	cmd.Flags().String("truncate", "none", "Truncate secret values in output: none, short, tiny")
}

func resolvedTruncateLevel(cmd *cobra.Command) truncate.Level {
	val, err := cmd.Flags().GetString("truncate")
	if err != nil {
		return truncate.LevelNone
	}
	switch truncate.Level(val) {
	case truncate.LevelShort, truncate.LevelTiny:
		return truncate.Level(val)
	default:
		return truncate.LevelNone
	}
}
