package config

import (
	"fmt"
	"os"

	"github.com/yourorg/vaultlink/internal/filter"
	"gopkg.in/yaml.v3"
)

// Config is the top-level vaultlink configuration.
type Config struct {
	Vault VaultConfig   `yaml:"vault"`
	Roles []filter.Role `yaml:"roles"`
	Output OutputConfig `yaml:"output"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	SecretPath string `yaml:"secret_path"`
}

// OutputConfig controls how the .env file is written.
type OutputConfig struct {
	File      string `yaml:"file"`
	Overwrite bool   `yaml:"overwrite"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("config: vault.address is required")
	}
	if c.Vault.SecretPath == "" {
		return fmt.Errorf("config: vault.secret_path is required")
	}
	if c.Output.File == "" {
		c.Output.File = ".env"
	}
	return nil
}
