package sync

import (
	"fmt"

	"github.com/user/vaultlink/internal/audit"
	"github.com/user/vaultlink/internal/config"
	"github.com/user/vaultlink/internal/diff"
	"github.com/user/vaultlink/internal/env"
	"github.com/user/vaultlink/internal/filter"
	"github.com/user/vaultlink/internal/snapshot"
	"github.com/user/vaultlink/internal/vault"
)

// Syncer orchestrates fetching, filtering, diffing, and writing secrets.
type Syncer struct {
	cfg     *config.Config
	logger  *audit.Logger
	manager *snapshot.Manager
}

// New creates a new Syncer instance.
func New(cfg *config.Config, logger *audit.Logger) *Syncer {
	return &Syncer{
		cfg:     cfg,
		logger:  logger,
		manager: snapshot.NewManager(cfg.OutputFile + ".snapshot"),
	}
}

// Run executes the full sync pipeline.
func (s *Syncer) Run() error {
	client, err := vault.NewClient(s.cfg.VaultAddress, s.cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("read secrets: %w", err)
	}

	f := filter.NewFilter(s.cfg.Role)
	filtered := f.Apply(secrets)

	previous, _ := s.manager.Load()
	changes := diff.Compare(previous, filtered)
	s.logger.Log(fmt.Sprintf("sync complete: %s", diff.Summary(changes)))

	if len(changes) == 0 {
		return nil
	}

	writer, err := env.NewWriter(s.cfg.OutputFile, s.cfg.Overwrite)
	if err != nil {
		return fmt.Errorf("env writer: %w", err)
	}
	if err := writer.Write(filtered); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	if err := s.manager.Save(filtered); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	return nil
}
