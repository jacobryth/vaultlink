package cmd

import (
	"testing"

	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func newCompressCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerCompressFlag(cmd)
	return cmd
}

func TestCompressFlag_Default(t *testing.T) {
	cmd := newCompressCmd()
	_ = cmd.ParseFlags([]string{})
	mode := resolvedCompressMode(cmd)
	if mode != env.CompressModeNone {
		t.Errorf("expected %q, got %q", env.CompressModeNone, mode)
	}
}

func TestCompressFlag_SetKeys(t *testing.T) {
	cmd := newCompressCmd()
	_ = cmd.ParseFlags([]string{"--compress", "keys"})
	mode := resolvedCompressMode(cmd)
	if mode != env.CompressModeKeys {
		t.Errorf("expected %q, got %q", env.CompressModeKeys, mode)
	}
}

func TestCompressFlag_SetValues(t *testing.T) {
	cmd := newCompressCmd()
	_ = cmd.ParseFlags([]string{"--compress", "values"})
	mode := resolvedCompressMode(cmd)
	if mode != env.CompressModeValues {
		t.Errorf("expected %q, got %q", env.CompressModeValues, mode)
	}
}

func TestCompressFlag_SetBoth(t *testing.T) {
	cmd := newCompressCmd()
	_ = cmd.ParseFlags([]string{"--compress", "both"})
	mode := resolvedCompressMode(cmd)
	if mode != env.CompressModeBoth {
		t.Errorf("expected %q, got %q", env.CompressModeBoth, mode)
	}
}
