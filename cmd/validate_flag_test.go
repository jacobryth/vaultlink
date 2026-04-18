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

func TestValidateFlag_OverwriteValue(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	registerValidateFlag(cmd)

	if err := cmd.Flags().Set("validate-rules", "/first/rules.yaml"); err != nil {
		t.Fatalf("failed to set initial flag value: %v", err)
	}
	if err := cmd.Flags().Set("validate-rules", "/second/rules.yaml"); err != nil {
		t.Fatalf("failed to overwrite flag value: %v", err)
	}

	val := resolvedValidateRulesFile()
	if val != "/second/rules.yaml" {
		t.Errorf("expected /second/rules.yaml after overwrite, got %q", val)
	}
}
