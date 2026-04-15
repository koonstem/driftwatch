package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindAuditFlags registers audit-related flags on the given command.
func BindAuditFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("audit", false, "enable audit log output")
	cmd.Flags().String("audit-file", "driftwatch-audit.jsonl", "path to audit log file (JSONL)")
	cmd.Flags().String("audit-run-id", "", "custom run identifier for audit entries (default: timestamp-based)")
}

// AuditOptionsFromFlags builds AuditOptions from bound cobra flags.
func AuditOptionsFromFlags(cmd *cobra.Command) (AuditOptions, error) {
	opts := DefaultAuditOptions()

	enabled, err := cmd.Flags().GetBool("audit")
	if err != nil {
		return opts, fmt.Errorf("audit flag: %w", err)
	}
	opts.Enabled = enabled

	filePath, err := cmd.Flags().GetString("audit-file")
	if err != nil {
		return opts, fmt.Errorf("audit-file flag: %w", err)
	}
	if filePath != "" {
		opts.FilePath = filePath
	}

	runID, err := cmd.Flags().GetString("audit-run-id")
	if err != nil {
		return opts, fmt.Errorf("audit-run-id flag: %w", err)
	}
	if runID != "" {
		opts.RunID = runID
	}

	return opts, nil
}
