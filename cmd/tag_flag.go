package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vaultlink/internal/tag"
)

func registerTagFlags(cmd *cobra.Command) {
	cmd.Flags().String("tag-level", "none", "Key tagging level: none, all, env")
	cmd.Flags().String("tag-prefix", "APP_", "Prefix to prepend to tagged keys")
}

func resolvedTagLevel(cmd *cobra.Command) (tag.Level, error) {
	v, _ := cmd.Flags().GetString("tag-level")
	level := tag.Level(v)
	switch level {
	case tag.LevelNone, tag.LevelAll, tag.LevelEnv:
		return level, nil
	default:
		return tag.LevelNone, fmt.Errorf("invalid tag-level %q: must be one of none, all, env", v)
	}
}

func resolvedTagPrefix(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("tag-prefix")
	return v
}
