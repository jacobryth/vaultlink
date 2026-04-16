package sync

import (
	"fmt"

	"github.com/vaultlink/internal/config"
	"github.com/vaultlink/internal/env"
	"github.com/vaultlink/internal/filter"
	"github.com/vaultlink/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to a .env file.
type Syncer struct {
	client *vault.Client
	writer *env.Writer
	filter *filter.Filter
	cfg    *config.Config
}

// New creates a new Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("sync: failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("sync: failed to create env writer: %w", err)
	}

	f := filter.NewFilter(cfg.Role, cfg.RolePrefixes)

	return &Syncer{
		client: client,
		writer: writer,
		filter: f,
		cfg:    cfg,
	}, nil
}

// Run performs the full sync: read secrets, filter by role, write to .env.
func (s *Syncer) Run() (int, error) {
	secrets, err := s.client.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return 0, fmt.Errorf("sync: failed to read secrets: %w", err)
	}

	filtered := s.filter.Apply(secrets)

	if err := s.writer.Write(filtered, s.cfg.Overwrite); err != nil {
		return 0, fmt.Errorf("sync: failed to write env file: %w", err)
	}

	return len(filtered), nil
}
