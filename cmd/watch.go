package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"vaultlink/internal/config"
	"vault"
	"vaultlink/internal/sync"
)

var watchInterval time.Duration

funcatchCmd.Flags().DurationVar(&watchInterval, "interval",*time.Second, "polling interval for vault synctregn	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Continuously sync secrets from Vault on a schedule",
	RunE:  runWatch,
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfgFile, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("watch: failed to load config: %w", err)
	}

	notifyLevel := resolvedNotifyLevel(cmd)

	syncer, err := sync.New(cfg, notifyLevel)
	if err != nil {
		return fmt.Errorf("watch: failed to create syncer: %w", err)
	}

	s, err := schedule.New(watchInterval, func(ctx context.Context) error {
		return syncer.Run(ctx)
	}, func(runErr error) {
		fmt.Fprintf(os.Stderr, "watch: sync error: %v\n", runErr)
	})
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stdout, "watch: starting sync every %s (ctrl+c to stop)\n", watchInterval)
	s.Run(ctx)
	fmt.Fprintln(os.Stdout, "watch: stopped")
	return nil
}
