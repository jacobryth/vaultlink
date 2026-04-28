package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func registerCompressFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"compress",
		string(env.CompressModeNone),
		`Remove blank/empty entries from secrets. Modes: none, keys, values, both.
  none   - no compression (default)
  keys   - remove entries with blank key names
  values - remove entries with empty values
  both   - remove entries with blank keys or empty values`,
	)
}

func resolvedCompressMode(cmd *cobra.Command) env.CompressMode {
	v, err := cmd.Flags().GetString("compress")
	if err != nil || v == "" {
		return env.CompressModeNone
	}
	return env.CompressMode(v)
}
