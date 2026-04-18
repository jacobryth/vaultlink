package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newSortCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerSortFlag(cmd)
	return cmd
}

func TestSortFlag_Default(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.Execute()
	if got := resolvedSortLevel(cmd); got != "none" {
		t.Errorf("expected 'none', got %q", got)
	}
}

func TestSortFlag_SetAsc(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.Flags().Set("sort", "asc")
	if got := resolvedSortLevel(cmd); got != "asc" {
		t.Errorf("expected 'asc', got %q", got)
	}
}

func TestSortFlag_SetDesc(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.Flags().Set("sort", "desc")
	if got := resolvedSortLevel(cmd); got != "desc" {
		t.Errorf("expected 'desc', got %q", got)
	}
}
