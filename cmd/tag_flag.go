package cmd

import (
	"github.com/spf13/cobra"
	"vaultlink/internal/tag"
)

func registerTagFlags(cmd *cobra.Command) {
	cmd.Flags().String("tag-level", "none", "Key tagging level: none, all, env")
	cmd.Flags().String("tag-prefix", "APP_", "Prefix to prepend to tagged keys")
}

func resolvedTagLevel(cmd *cobra.Command) tag.Level {
	v, _ := cmd.Flags().GetString("tag-level")
	return tag.Level(v)
}

func resolvedTagPrefix(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("tag-prefix")
	return v
}
