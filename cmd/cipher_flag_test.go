package cmd_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultlink/cmd"
	"github.com/yourusername/vaultlink/internal/cipher"
)

func newCipherCmd() *cobra.Command {
	c := &cobra.Command{Use: "test"}
	cmd.registerCipherFlags(c)
	return c
}

func TestCipherFlag_Default(t *testing.T) {
	c := newCipherCmd()
	_ = c.ParseFlags([]string{})
	if got := cmd.resolvedCipherLevel(c); got != cipher.LevelNone {
		t.Errorf("expected %q, got %q", cipher.LevelNone, got)
	}
}

func TestCipherFlag_SetEncrypt(t *testing.T) {
	c := newCipherCmd()
	_ = c.ParseFlags([]string{"--cipher", "encrypt"})
	if got := cmd.resolvedCipherLevel(c); got != cipher.LevelEncrypt {
		t.Errorf("expected %q, got %q", cipher.LevelEncrypt, got)
	}
}

func TestCipherFlag_SetDecrypt(t *testing.T) {
	c := newCipherCmd()
	_ = c.ParseFlags([]string{"--cipher", "decrypt"})
	if got := cmd.resolvedCipherLevel(c); got != cipher.LevelDecrypt {
		t.Errorf("expected %q, got %q", cipher.LevelDecrypt, got)
	}
}

func TestCipherKeyFlag_Default(t *testing.T) {
	c := newCipherCmd()
	_ = c.ParseFlags([]string{})
	if got := cmd.resolvedCipherKey(c); got != "" {
		t.Errorf("expected empty default key, got %q", got)
	}
}

func TestCipherKeyFlag_SetValue(t *testing.T) {
	c := newCipherCmd()
	_ = c.ParseFlags([]string{"--cipher-key", "c29tZWtleQ=="})
	if got := cmd.resolvedCipherKey(c); got != "c29tZWtleQ==" {
		t.Errorf("unexpected key value: %q", got)
	}
}
