package cipher_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/yourusername/vaultlink/internal/cipher"
)

// 32-byte AES-256 key encoded as base64
var testKey = base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))

func TestNew_ValidLevels(t *testing.T) {
	for _, lvl := range []cipher.Level{cipher.LevelNone, cipher.LevelEncrypt, cipher.LevelDecrypt} {
		_, err := cipher.New(lvl, testKey)
		if err != nil {
			t.Errorf("level %q: unexpected error: %v", lvl, err)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := cipher.New("scramble", testKey)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestNew_InvalidKey(t *testing.T) {
	_, err := cipher.New(cipher.LevelEncrypt, "not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for bad key encoding")
	}
}

func TestNew_WrongKeyLength(t *testing.T) {
	shortKey := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	_, err := cipher.New(cipher.LevelEncrypt, shortKey)
	if err == nil {
		t.Fatal("expected error for wrong key length")
	}
}

func TestApply_NilSecrets(t *testing.T) {
	c, _ := cipher.New(cipher.LevelNone, testKey)
	out, err := c.Apply(nil)
	if err != nil || out != nil {
		t.Fatalf("expected nil,nil got %v,%v", out, err)
	}
}

func TestApply_NoEncryption(t *testing.T) {
	c, _ := cipher.New(cipher.LevelNone, testKey)
	in := map[string]string{"KEY": "value"}
	out, err := c.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected passthrough, got %q", out["KEY"])
	}
}

func TestApply_EncryptThenDecrypt(t *testing.T) {
	enc, _ := cipher.New(cipher.LevelEncrypt, testKey)
	dec, _ := cipher.New(cipher.LevelDecrypt, testKey)

	secrets := map[string]string{
		"DB_PASS": "supersecret",
		"API_KEY": "abc123",
	}

	encrypted, err := enc.Apply(secrets)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	for k, v := range encrypted {
		if v == secrets[k] {
			t.Errorf("key %q: value was not changed after encryption", k)
		}
	}

	decrypted, err := dec.Apply(encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	for k, want := range secrets {
		if decrypted[k] != want {
			t.Errorf("key %q: got %q want %q", k, decrypted[k], want)
		}
	}
}

func TestApply_DecryptBadCiphertext(t *testing.T) {
	dec, _ := cipher.New(cipher.LevelDecrypt, testKey)
	_, err := dec.Apply(map[string]string{"X": "notbase64!!!"})
	if err == nil || !strings.Contains(err.Error(), "X") {
		t.Fatalf("expected error mentioning key, got %v", err)
	}
}
