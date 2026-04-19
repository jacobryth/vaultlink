package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newLinebreakCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerLinebreakFlag(cmd)
	return cmd
}

func TestLinebreakFlag_Default(t *testing.T) {
	cmd := newLinebreakCmd()
	_ = cmd.Execute()
	if got := resolvedLinebreakLevel(cmd); got != "none" {
		t.Errorf("expected 'none', got %q", got)
	}
}

func TestLinebreakFlag_SetUnix(t *testing.T) {
	cmd := newLinebreakCmd()
	cmd.SetArgs([]string{"--linebreak", "unix"})
	_ = cmd.Execute()
	if got := resolvedLinebreakLevel(cmd); got != "unix" {
		t.Errorf("expected 'unix', got %q", got)
	}
}

func TestLinebreakFlag_SetWindows(t *testing.T) {
	cmd := newLinebreakCmd()
	cmd.SetArgs([]string{"--linebreak", "windows"})
	_ = cmd.Execute()
	if got := resolvedLinebreakLevel(cmd); got != "windows" {
		t.Errorf("expected 'windows', got %q", got)
	}
}
