package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registerLoadFlag attaches the --load-env flag to the given command.
// When set, vaultlink will read the specified .env file and merge its
// existing values before writing, respecting the overwrite mode.
func registerLoadFlag(cmd *cobra.Command) {
	cmd.Flags().String("load-env", "", "path to an existing .env file to load before syncing")
}

// resolvedLoadEnvPath returns the --load-env flag value from the command.
func resolvedLoadEnvPath(cmd *cobra.Command) (string, error) {
	v, err := cmd.Flags().GetString("load-env")
	if err != nil {
		return "", fmt.Errorf("load-env flag: %w", err)
	}
	return v, nil
}
