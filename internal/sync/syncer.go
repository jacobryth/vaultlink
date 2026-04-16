package sync

import (
	"fmt"

	"github.com/user/vaultlink/internal/audit"
	"github.com/user/vaultlink/internal/config"
	"github.com/user/vaultlink/internal/env"
	"github.com/user/vaultlink/internal/filter"
	"github.com/user/vaultlink/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	logger *audit.Logger
}

// New creates a Syncer from the given config.
func New(cfg *config.Config) (*Syncer, error) {
	l, err := audit.NewLogger(cfg.AuditLog)
	if err != nil {
		return nil, fmt.Errorf("syncer: init audit logger: %w", err)
	}
	return &Syncer{cfg: cfg, logger: l}, nil
}

// Run executes the full sync pipeline.
func (s *Syncer) Run() error {
	defer s.logger.Close()

	c, err := vault.NewClient(s.cfg.VaultAddress, s.cfg.VaultToken)
	if err != nil {
		_ = s.logger.Log("sync", s.cfg.SecretPath, s.cfg.Role, nil, false, err.Error())
		return fmt.Errorf("syncer: vault client: %w", err)
	}

	secrets, err := c.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		_ = s.logger.Log("sync", s.cfg.SecretPath, s.cfg.Role, nil, false, err.Error())
		return fmt.Errorf("syncer: read secrets: %w", err)
	}

	f := filter.NewFilter(s.cfg.Role)
	filtered := f.Apply(secrets)

	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}

	w, err := env.NewWriter(s.cfg.OutputFile)
	if err != nil {
		_ = s.logger.Log("sync", s.cfg.SecretPath, s.cfg.Role, keys, false, err.Error())
		return fmt.Errorf("syncer: env writer: %w", err)
	}

	if err := w.Write(filtered, s.cfg.Overwrite); err != nil {
		_ = s.logger.Log("sync", s.cfg.SecretPath, s.cfg.Role, keys, false, err.Error())
		return fmt.Errorf("syncer: write env: %w", err)
	}

	return s.logger.Log("sync", s.cfg.SecretPath, s.cfg.Role, keys, true, "")
}
