package cmd

import (
	"github.com/spf13/cobra"
)

var validateRulesFile string

// registerValidateFlag registers the --validate-rules flag on the given command.
// The flag accepts a path to a YAML file defining validation rules for secrets.
func registerValidateFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&validateRulesFile,
		"validate-rules",
		"",
		"Path to a YAML file defining validation rules for secrets",
	)
}

// resolvedValidateRulesFile returns the path to the validation rules file
// as provided via the --validate-rules flag.
func resolvedValidateRulesFile() string {
	return validateRulesFile
}

// hasValidateRulesFile reports whether a validation rules file has been specified.
func hasValidateRulesFile() bool {
	return validateRulesFile != ""
}
