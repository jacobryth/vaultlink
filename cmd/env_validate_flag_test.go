package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newEnvValidateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerEnvValidateFlag(cmd)
	return cmd
}

func TestEnvValidateFlag_Default(t *testing.T) {
	cmd := newEnvValidateCmd()
	_ = cmd.Execute()
	val, err := cmd.Flags().GetString("env-validate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "none" {
		t.Errorf("expected default 'none', got %q", val)
	}
}

func TestEnvValidateFlag_SetKeys(t *testing.T) {
	cmd := newEnvValidateCmd()
	_ = cmd.Flags().Set("env-validate", "keys")
	mode, err := resolvedEnvValidateMode(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "keys" {
		t.Errorf("expected 'keys', got %q", mode)
	}
}

func TestEnvValidateFlag_SetBoth(t *testing.T) {
	cmd := newEnvValidateCmd()
	_ = cmd.Flags().Set("env-validate", "both")
	mode, err := resolvedEnvValidateMode(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "both" {
		t.Errorf("expected 'both', got %q", mode)
	}
}

func TestEnvValidateFlag_InvalidMode(t *testing.T) {
	cmd := newEnvValidateCmd()
	_ = cmd.Flags().Set("env-validate", "strict")
	_, err := resolvedEnvValidateMode(cmd)
	if err == nil {
		t.Fatal("expected error for unknown mode 'strict'")
	}
}
