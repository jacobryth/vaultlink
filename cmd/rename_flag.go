package cmd

import (
	"github.com/spf13/cobra"
	"vaultlink/internal/rename"
)

func registerRenameFlags(cmd *cobra.Command) {
	cmd.Flags().String("rename", "none", "Key rename strategy: none, snake, kebab, custom")
	cmd.Flags().StringSlice("rename-rules", nil, "Custom rename rules in FROM=TO format (used with --rename=custom)")
}

func resolvedRenameLevel(cmd *cobra.Command) rename.Level {
	v, _ := cmd.Flags().GetString("rename")
	return rename.Level(v)
}

func resolvedRenameRules(cmd *cobra.Command) []rename.Rule {
	raw, _ := cmd.Flags().GetStringSlice("rename-rules")
	var rules []rename.Rule
	for _, entry := range raw {
		for i := 0; i < len(entry); i++ {
			if entry[i] == '=' {
				rules = append(rules, rename.Rule{From: entry[:i], To: entry[i+1:]})
				break
			}
		}
	}
	return rules
}
