package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func registerDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().String("default-level", "none",
		"Default-fill level: none, missing, empty, both")
	cmd.Flags().StringSlice("default-rules", nil,
		"Default rules as KEY=VALUE pairs (e.g. DB_PORT=5432)")
}

func resolvedDefaultLevel(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("default-level")
	return strings.TrimSpace(v)
}

func resolvedDefaultRules(cmd *cobra.Command) map[string]string {
	pairs, _ := cmd.Flags().GetStringSlice("default-rules")
	rules := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			rules[strings.TrimSpace(parts[0])] = parts[1]
		}
	}
	return rules
}
