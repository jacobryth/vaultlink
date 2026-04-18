package cmd

import (
	"github.com/spf13/cobra"
)

var validateRulesFile string

func registerValidateFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&validateRulesFile,
		"validate-rules",
		"",
		"Path to a YAML file defining validation rules for secrets",
	)
}

func resolvedValidateRulesFile() string {
	return validateRulesFile
}
