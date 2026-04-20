package output

import (
	"github.com/spf13/cobra"
)

// BindBatchFlags registers batch-related flags on the given command.
func BindBatchFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("batch", false, "enable batched output")
	cmd.Flags().Int("batch-size", 10, "number of results per batch")
}

// BatchOptionsFromFlags reads flag values from the command and returns BatchOptions.
func BatchOptionsFromFlags(cmd *cobra.Command) BatchOptions {
	enabled, _ := cmd.Flags().GetBool("batch")
	size, _ := cmd.Flags().GetInt("batch-size")
	return BatchOptions{
		Enabled:   enabled,
		BatchSize: size,
	}
}
