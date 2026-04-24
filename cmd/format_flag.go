package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

// registerFormatFlag attaches the --format flag to the given command.
func registerFormatFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"format",
		string(env.FormatNone),
		`Output format for env lines: none (KEY=VALUE), export (export KEY=VALUE), inlined (single line)`,
	)
}

// resolvedFormatLevel reads and validates the --format flag from the command.
// Falls back to FormatNone on an empty or missing value.
func resolvedFormatLevel(cmd *cobra.Command) (env.FormatLevel, error) {
	raw, err := cmd.Flags().GetString("format")
	if err != nil || raw == "" {
		return env.FormatNone, nil
	}
	level := env.FormatLevel(raw)
	// Validate by attempting construction.
	if _, err := env.NewFormatter(level); err != nil {
		return "", err
	}
	return level, nil
}
