package cmd

import (
	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func registerBackupFlags(cmd *cobra.Command) {
	cmd.Flags().String("backup-mode", "none", `backup mode before writing env file: none, always, on_write`)
	cmd.Flags().String("backup-suffix", ".bak", `suffix appended to backup filename`)
}

func resolvedBackupMode(cmd *cobra.Command) env.BackupMode {
	v, _ := cmd.Flags().GetString("backup-mode")
	return env.BackupMode(v)
}

func resolvedBackupSuffix(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("backup-suffix")
	return v
}
