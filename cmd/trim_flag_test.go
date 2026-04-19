package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTrimCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerTrimFlag(cmd)
	return cmd
}

func TestTrimFlag_Default(t *testing.T) {
	cmd := newTrimCmd()
	_ = cmd.Execute()
	val, err := cmd.Flags().GetString("trim")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "none" {
		t.Errorf("expected default 'none', got %q", val)
	}
}

func TestTrimFlag_SetSpace(t *testing.T) {
	cmd := newTrimCmd()
	_ = cmd.Flags().Set("trim", "space")
	if resolvedTrimLevel(cmd) != "space" {
		t.Error("expected 'space'")
	}
}

func TestTrimFlag_SetAll(t *testing.T) {
	cmd := newTrimCmd()
	_ = cmd.Flags().Set("trim", "all")
	if resolvedTrimLevel(cmd) != "all" {
		t.Error("expected 'all'")
	}
}
