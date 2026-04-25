package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Level controls encryption behaviour.
type Level string

const (
	LevelNone    Level = "none"
	LevelEncrypt Level = "encrypt"
	LevelDecrypt Level = "decrypt"
)

var validLevels = map[Level]struct{}{
	LevelNone:    {},
	LevelEncrypt: {},
	LevelDecrypt: {},
}

// Cipher applies AES-GCM encryption or decryption to secret values.
type Cipher struct {
	level Level
	key   []byte // must be 16, 24, or 32 bytes
}

// New returns a Cipher for the given level and base64-encoded key.
func New(level Level, b64Key string) (*Cipher, error) {
	if _, ok := validLevels[level]; !ok {
		return nil, fmt.Errorf("cipher: unknown level %q", level)
	}
	if level == LevelNone {
		return &Cipher{level: level}, nil
	}
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, fmt.Errorf("cipher: invalid key encoding: %w", err)
	}
	if l := len(key); l != 16 && l != 24 && l != 32 {
		return nil, fmt.Errorf("cipher: key must be 16, 24, or 32 bytes, got %d", l)
	}
	return &Cipher{level: level, key: key}, nil
}

// Level returns the configured encryption level.
func (c *Cipher) Level() Level {
	return c.level
}

// Apply transforms secret values according to the configured level.
func (c *Cipher) Apply(secrets map[string]string) (map[string]string, error) {
	if secrets == nil {
		return nil, nil
	}
	if c.level == LevelNone {
		return secrets, nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		var (transformed string; err error)
		if c.level == LevelEncrypt {
			transformed, err = c.encrypt(v)
		} else {
			transformed, err = c.decrypt(v)
		}
		if err != nil {
			return nil, fmt.Errorf("cipher: key %q: %w", k, err)
		}
		out[k] = transformed
	}
	return out, nil
}

func (c *Cipher) newGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func (c *Cipher) encrypt(plaintext string) (string, error) {
	gcm, err := c.newGCM()
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (c *Cipher) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	gcm, err := c.newGCM()
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(plain), nil
}
