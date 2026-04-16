package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockVaultServer(t *testing.T, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("failed to encode mock response: %v", err)
		}
	}))
}

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestReadSecrets_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
		},
	}
	srv := mockVaultServer(t, "/v1/secret/data/myapp", payload)
	defer srv.Close()

	client, err := NewClient(Config{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	secrets, err := client.ReadSecrets(context.Background(), "myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}
