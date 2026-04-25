package cmd

import (
	"testing"

	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func newPatchCmd() *cobra.Command {
	c := &cobra.Command{Use: "test"}
	registerPatchFlag(c)
	return c
}

func TestPatchFlag_Default(t *testing.T) {
	c := newPatchCmd()
	_ = c.ParseFlags([]string{})
	mode, err := resolvedPatchMode(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != env.PatchModeUpsert {
		t.Errorf("expected upsert, got %q", mode)
	}
}

func TestPatchFlag_SetAdd(t *testing.T) {
	c := newPatchCmd()
	_ = c.ParseFlags([]string{"--patch-mode", "add"})
	mode, err := resolvedPatchMode(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != env.PatchModeAddOnly {
		t.Errorf("expected add, got %q", mode)
	}
}

func TestPatchFlag_SetRemove(t *testing.T) {
	c := newPatchCmd()
	_ = c.ParseFlags([]string{"--patch-mode", "remove"})
	mode, err := resolvedPatchMode(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != env.PatchModeRemove {
		t.Errorf("expected remove, got %q", mode)
	}
}

func TestPatchFlag_InvalidMode(t *testing.T) {
	c := newPatchCmd()
	_ = c.ParseFlags([]string{"--patch-mode", "invalid"})
	_, err := resolvedPatchMode(c)
	if err == nil {
		t.Fatal("expected error for invalid patch mode")
	}
}
