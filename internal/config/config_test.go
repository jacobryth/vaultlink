package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultlink/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultlink.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
vault:
  address: http://127.0.0.1:8200
  token: root
  secret_path: secret/data/app
roles:
  - name: backend
    prefixes: ["DB_", "APP_"]
output:
  file: .env
  overwrite: false
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected address: %s", cfg.Vault.Address)
	}
	if len(cfg.Roles) != 1 || cfg.Roles[0].Name != "backend" {
		t.Errorf("unexpected roles: %+v", cfg.Roles)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTemp(t, `
vault:
  secret_path: secret/data/app
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault.address")
	}
}

func TestLoad_DefaultOutputFile(t *testing.T) {
	path := writeTemp(t, `
vault:
  address: http://127.0.0.1:8200
  secret_path: secret/data/app
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Output.File != ".env" {
		t.Errorf("expected default output file '.env', got %s", cfg.Output.File)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/vaultlink.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
