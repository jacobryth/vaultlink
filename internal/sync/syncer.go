package sync

import (
	"fmt"

	"github.com/user/vaultlink/internal/audit"
	"github.com/user/vaultlink/internal/config"
	"github.com/user/vaultlink/internal/env"
	"github.com/user/vaultlink/internal/filter"
	"github.com/user/vaultlink/internal/snapshot"
	"github.com/user/vaultlink/internal/vault"
)

// Syncer orchestrates reading from Vault and writing to .env.
type Syncer struct {
	cfg      *config.Config
	logger   *audit.Logger
	snapshot *snapshot.Manager
}

// New creates a Syncer from the given config and logger.
func New(cfg *config.Config, logger *audit.Logger) *Syncer {
	return &Syncer{
		cfg:      cfg,
		logger:   logger,
		snapshot: snapshot.NewManager(".vaultlink.snapshot.json"),
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

	changed, err := s.snapshot.HasChanged(filtered)
	if err != nil {
		return fmt.Errorf("snapshot check: %w", err)
	}
	if !changed {
		s.logger.Log("no changes detected, skipping write")
		return nil
	}

	w := env.NewWriter(s.cfg.OutputFile, s.cfg.Overwrite)
	if err := w.Write(filtered); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	if err := s.snapshot.Save(&snapshot.Snapshot{
		SecretPath: s.cfg.SecretPath,
		Keys:       keys,
		Checksum:   filtered,
	}); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	s.logger.Log(fmt.Sprintf("synced %d keys to %s", len(filtered), s.cfg.OutputFile))
	return nil
}
