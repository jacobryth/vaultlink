package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultlink/internal/config"
	"github.com/vaultlink/internal/sync"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Sync Vault secrets to a local .env file",
	RunE:  runSync,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("config", "c", "vaultlink.yaml", "Path to config file")
	runCmd.Flags().Bool("overwrite", false, "Overwrite existing keys in .env file")
}

func runSync(cmd *cobra.Command, _ []string) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	overwrite, err := cmd.Flags().GetBool("overwrite")
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if overwrite {
		cfg.Overwrite = true
	}

	syncer, err := sync.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize syncer: %w", err)
	}

	count, err := syncer.Run()
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✓ Synced %d secret(s) to %s\n", count, cfg.OutputFile)
	return nil
}
