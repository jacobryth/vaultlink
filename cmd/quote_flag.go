package cmd

import (
	"github.com/spf13/cobra"
	"vaultlink/internal/quote"
)

func registerQuoteFlag(cmd *cobra.Command) {
	cmd.Flags().String("quote", "none", "Quote style for secret values: none, double, single")
}

func resolvedQuoteLevel(cmd *cobra.Command) quote.Level {
	val, err := cmd.Flags().GetString("quote")
	if err != nil || val == "" {
		return quote.None
	}
	return quote.Level(val)
}
