package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func registerPatchFlag(cmd *cobra.Command) {
	cmd.Flags().String("patch-mode", "upsert",
		`Patch mode for targeted .env updates: upsert, add, remove`)
}

// resolvedPatchMode returns the validated PatchMode from the command flags.
func resolvedPatchMode(cmd *cobra.Command) (env.PatchMode, error) {
	raw, err := cmd.Flags().GetString("patch-mode")
	if err != nil {
		return "", err
	}
	m := env.PatchMode(raw)
	_, err = env.NewPatcher(m)
	if err != nil {
		return "", err
	}
	return m, nil
}
