package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/merge"
)

// registerMergeFlag attaches the --merge flag to the given command.
func registerMergeFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"merge",
		string(merge.StrategyNone),
		`secret merge strategy when combining multiple sources (none|overwrite|keep-first)`,
	)
}

// resolvedMergeStrategy returns the validated merge.Strategy from the flag,
// defaulting to StrategyNone on any error.
func resolvedMergeStrategy(cmd *cobra.Command) merge.Strategy {
	raw, err := cmd.Flags().GetString("merge")
	if err != nil {
		return merge.StrategyNone
	}
	s := merge.Strategy(raw)
	if _, err := merge.New(s); err != nil {
		return merge.StrategyNone
	}
	return s
}
