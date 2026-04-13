package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindExportFlags attaches export-related flags to cmd.
func BindExportFlags(cmd *cobra.Command) {
	cmd.Flags().String("export-format", "csv", `export format: csv or json`)
	cmd.Flags().Bool("export-timestamp", true, "include generated_at timestamp in export")
	cmd.Flags().String("export-path", "", "write export to this file path (empty = stdout)")
}

// ExportOptionsFromFlags reads export flags from cmd and returns ExportOptions.
func ExportOptionsFromFlags(cmd *cobra.Command) (ExportOptions, string, error) {
	fmt_, err := cmd.Flags().GetString("export-format")
	if err != nil {
		return ExportOptions{}, "", err
	}
	if fmt_ != "csv" && fmt_ != "json" {
		return ExportOptions{}, "", fmt.Errorf("export: unsupported format %q (must be csv or json)", fmt_)
	}
	ts, err := cmd.Flags().GetBool("export-timestamp")
	if err != nil {
		return ExportOptions{}, "", err
	}
	path, err := cmd.Flags().GetString("export-path")
	if err != nil {
		return ExportOptions{}, "", err
	}
	return ExportOptions{Format: fmt_, Timestamp: ts}, path, nil
}
