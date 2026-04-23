package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// registerChainFlag adds the --filter-chain flag to the given command.
// The flag accepts a comma-separated list of filter step names to apply in order.
// Supported values: prefix, role, redact
func registerChainFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"filter-chain",
		"",
		`Ordered comma-separated filter steps to apply (e.g. "prefix,role,redact")`,
	)
}

// resolvedChainSteps returns the ordered list of filter step names from the
// --filter-chain flag on the given command. Returns nil if the flag is empty.
func resolvedChainSteps(cmd *cobra.Command) []string {
	raw, err := cmd.Flags().GetString("filter-chain")
	if err != nil || strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	steps := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			steps = append(steps, p)
		}
	}
	return steps
}
