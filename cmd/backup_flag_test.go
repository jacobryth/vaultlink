package cmd

import (
	"testing"

	"github.com/spf13/cobra"

	"vaultlink/internal/env"
)

func newBackupCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	registerBackupFlags(cmd)
	return cmd
}

func TestBackupFlag_DefaultMode(t *testing.T) {
	cmd := newBackupCmd()
	_ = cmd.ParseFlags([]string{})
	if got := resolvedBackupMode(cmd); got != env.BackupNone {
		t.Errorf("expected %q, got %q", env.BackupNone, got)
	}
}

func TestBackupFlag_DefaultSuffix(t *testing.T) {
	cmd := newBackupCmd()
	_ = cmd.ParseFlags([]string{})
	if got := resolvedBackupSuffix(cmd); got != ".bak" {
		t.Errorf("expected .bak, got %q", got)
	}
}

func TestBackupFlag_SetAlways(t *testing.T) {
	cmd := newBackupCmd()
	_ = cmd.ParseFlags([]string{"--backup-mode", "always"})
	if got := resolvedBackupMode(cmd); got != env.BackupAlways {
		t.Errorf("expected %q, got %q", env.BackupAlways, got)
	}
}

func TestBackupFlag_SetOnWrite(t *testing.T) {
	cmd := newBackupCmd()
	_ = cmd.ParseFlags([]string{"--backup-mode", "on_write"})
	if got := resolvedBackupMode(cmd); got != env.BackupOnWrite {
		t.Errorf("expected %q, got %q", env.BackupOnWrite, got)
	}
}

func TestBackupFlag_CustomSuffix(t *testing.T) {
	cmd := newBackupCmd()
	_ = cmd.ParseFlags([]string{"--backup-suffix", ".backup"})
	if got := resolvedBackupSuffix(cmd); got != ".backup" {
		t.Errorf("expected .backup, got %q", got)
	}
}
