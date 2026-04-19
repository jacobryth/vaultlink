package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newUppercaseCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerUppercaseFlag(cmd)
	return cmd
}

func TestUppercaseFlag_Default(t *testing.T) {
	cmd := newUppercaseCmd()
	_ = cmd.Execute()
	if v := resolvedUppercaseLevel(cmd); v != "none" {
		t.Fatalf("expected none, got %q", v)
	}
}

func TestUppercaseFlag_SetKeys(t *testing.T) {
	cmd := newUppercaseCmd()
	cmd.SetArgs([]string{"--uppercase=keys"})
	_ = cmd.Execute()
	if v := resolvedUppercaseLevel(cmd); v != "keys" {
		t.Fatalf("expected keys, got %q", v)
	}
}

func TestUppercaseFlag_SetBoth(t *testing.T) {
	cmd := newUppercaseCmd()
	cmd.SetArgs([]string{"--uppercase=both"})
	_ = cmd.Execute()
	if v := resolvedUppercaseLevel(cmd); v != "both" {
		t.Fatalf("expected both, got %q", v)
	}
}
