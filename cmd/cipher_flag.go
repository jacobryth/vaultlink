package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yourusername/vaultlink/internal/cipher"
)

func registerCipherFlags(cmd *cobra.Command) {
	cmd.Flags().String("cipher", "none", `cipher level: none | encrypt | decrypt`)
	cmd.Flags().String("cipher-key", "", "base64-encoded AES key (16, 24, or 32 bytes)")
}

func resolvedCipherLevel(cmd *cobra.Command) cipher.Level {
	v, _ := cmd.Flags().GetString("cipher")
	return cipher.Level(v)
}

func resolvedCipherKey(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("cipher-key")
	return v
}
