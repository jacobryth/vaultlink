package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newEnvSortCmd() *cobra.Command {
	c := &cobra.Command{Use: "test"}
	registerEnvSortFlag(c)
	return c
}

func TestEnvSortFlag_Default(t *testing.T) {
	c := newEnvSortCmd()
	_ = c.Execute()
	got := resolvedEnvSortMode(c)
	if got != "none" {
		t.Fatalf("expected default 'none', got %q", got)
	}
}

func TestEnvSortFlag_SetAsc(t *testing.T) {
	c := newEnvSortCmd()
	_ = c.Flags().Set("env-sort", "asc")
	got := resolvedEnvSortMode(c)
	if got != "asc" {
		t.Fatalf("expected 'asc', got %q", got)
	}
}

func TestEnvSortFlag_SetDesc(t *testing.T) {
	c := newEnvSortCmd()
	_ = c.Flags().Set("env-sort", "desc")
	got := resolvedEnvSortMode(c)
	if got != "desc" {
		t.Fatalf("expected 'desc', got %q", got)
	}
}
