package cmd

import (
	"fmt"
	"os"

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

// validateRulesFileExists checks that the validation rules file specified via
// --validate-rules exists and is a regular file. It returns an error if no file
// was specified or if the path cannot be stat'd.
func validateRulesFileExists() error {
	if validateRulesFile == "" {
		return fmt.Errorf("no validation rules file specified; use --validate-rules to provide one")
	}
	info, err := os.Stat(validateRulesFile)
	if err != nil {
		return fmt.Errorf("validation rules file %q not found: %w", validateRulesFile, err)
	}
	if info.IsDir() {
		return fmt.Errorf("validation rules file %q is a directory, not a file", validateRulesFile)
	}
	return nil
}
