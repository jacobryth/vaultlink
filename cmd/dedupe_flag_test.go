package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newDedupeCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerDedupeFlag(cmd)
	return cmd
}

func TestDedupeFlag_Default(t *testing.T) {
	cmd := newDedupeCmd()
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
	cmd.Execute()
	if v := resolvedDedupeStrategy(cmd); v != "first" {
		t.Errorf("expected 'first', got %q", v)
	}
}

func TestDedupeFlag_SetLast(t *testing.T) {
	cmd := newDedupeCmd()
	cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
	cmd.SetArgs([]string{"--dedupe=last"})
	cmd.Execute()
	if v := resolvedDedupeStrategy(cmd); v != "last" {
		t.Errorf("expected 'last', got %q", v)
	}
}
