package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateFlag_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	registerValidateFlag(cmd)

	val := resolvedValidateRulesFile()
	if val != "" {
		t.Errorf("expected empty default, got %q", val)
	}
}

func TestValidateFlag_SetValue(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	registerValidateFlag(cmd)

	if err := cmd.Flags().Set("validate-rules", "/etc/rules.yaml"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	val := resolvedValidateRulesFile()
	if val != "/etc/rules.yaml" {
		t.Errorf("expected /etc/rules.yaml, got %q", val)
	}
}
