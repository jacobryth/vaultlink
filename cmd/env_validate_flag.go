package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultEnvValidateMode = "none"

func registerEnvValidateFlag(cmd *cobra.Command) {
	cmd.Flags().String(
		"env-validate",
		defaultEnvValidateMode,
		`Validate env key/value pairs before writing. Modes: none, keys, values, both`,
	)
}

func resolvedEnvValidateMode(cmd *cobra.Command) (string, error) {
	val, err := cmd.Flags().GetString("env-validate")
	if err != nil {
		return "", fmt.Errorf("env-validate flag: %w", err)
	}
	switch val {
	case "none", "keys", "values", "both":
		return val, nil
	default:
		return "", fmt.Errorf("env-validate: unknown mode %q (want none|keys|values|both)", val)
	}
}
