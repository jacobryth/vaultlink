package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRedactFlag_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	registerRedactFlag(cmd)

	f := cmd.Flags().Lookup("redact")
	if f == nil {
		t.Fatal("expected --redact flag to be registered")
	}
	if f.DefValue != "true" {
		t.Errorf("expected default true, got %s", f.DefValue)
	}
}

func TestRedactFlag_SetFalse(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	registerRedactFlag(cmd)

	err := cmd.Flags().Set("redact", "false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if redactEnabled {
		t.Error("expected redactEnabled to be false")
	}
}
