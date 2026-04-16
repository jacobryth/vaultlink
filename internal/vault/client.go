package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	api    *vaultapi.Client
	Mount  string
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	Mount   string
}

// NewClient creates and authenticates a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = os.Getenv("VAULT_ADDR")
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("VAULT_TOKEN")
	}
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or pass --address)")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}

	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = cfg.Address

	c, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}
	c.SetToken(cfg.Token)

	mount := cfg.Mount
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: c, Mount: mount}, nil
}

// ReadSecrets reads all key-value pairs from the given secret path.
func (c *Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	fullPath := fmt.Sprintf("%s/data/%s", c.Mount, path)
	secret, err := c.api.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at path %q", fullPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret format at %q", fullPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
