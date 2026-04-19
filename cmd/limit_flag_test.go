package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func newLimitCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerLimitFlags(cmd)
	return cmd
}

func TestLimitFlag_Default(t *testing.T) {
	cmd := newLimitCmd()
	cmd.ParseFlags([]string{})
	if got := resolvedLimitLevel(cmd); got != "none" {
		t.Errorf("expected none, got %s", got)
	}
	if got := resolvedLimitCount(cmd); got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

func TestLimitFlag_SetFirst(t *testing.T) {
	cmd := newLimitCmd()
	cmd.ParseFlags([]string{"--limit=first", "--limit-count=5"})
	if got := resolvedLimitLevel(cmd); got != "first" {
		t.Errorf("expected first, got %s", got)
	}
	if got := resolvedLimitCount(cmd); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestLimitFlag_SetLast(t *testing.T) {
	cmd := newLimitCmd()
	cmd.ParseFlags([]string{"--limit=last", "--limit-count=3"})
	if got := resolvedLimitLevel(cmd); got != "last" {
		t.Errorf("expected last, got %s", got)
	}
	if got := resolvedLimitCount(cmd); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}
