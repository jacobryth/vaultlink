package sync

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultlink/internal/config"
)

func mockVaultServer(t *testing.T, payload string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(payload))
	}))
}

func newTestConfig(t *testing.T, vaultURL, outFile string) *config.Config {
	t.Helper()
	return &config.Config{
		VaultAddress: vaultURL,
		VaultToken:   "test-token",
		SecretPath:   "secret/data/app",
		OutputFile:   outFile,
		Role:         "",
		RolePrefixes: []string{},
		Overwrite:    true,
	}
}

func TestRun_Success(t *testing.T) {
	server := mockVaultServer(t, `{"data":{"data":{"APP_KEY":"abc","DB_PASS":"secret"}}}`)
	defer server.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	cfg := newTestConfig(t, server.URL, outFile)

	syncer, err := New(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	count, err := syncer.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 secrets written, got %d", count)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	contents := string(data)
	if len(contents) == 0 {
		t.Error("expected non-empty .env file")
	}
}

func TestRun_InvalidVaultAddress(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{
		VaultAddress: "",
		VaultToken:   "token",
		SecretPath:   "secret/data/app",
		OutputFile:   filepath.Join(tmpDir, ".env"),
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}
