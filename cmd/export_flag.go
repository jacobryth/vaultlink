package cmd

import (
	"github.com/spf13/cobra"
)

var exportFormat string
var exportOutput string

func registerExportFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&exportFormat, "export-format", "env", "Export format: env or json")
	cmd.Flags().StringVar(&exportOutput, "export-output", "", "Destination file path for exported secrets (overrides output_file)")
}

func resolvedExportFormat() string {
	if exportFormat == "" {
		return "env"
	}
	return exportFormat
}

func resolvedExportOutput(fallback string) string {
	if exportOutput != "" {
		return exportOutput
	}
	return fallback
}
