package output

import "github.com/spf13/cobra"

// BindFilterWriterFlags registers filter-writer flags on cmd.
func BindFilterWriterFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("fw-only-drifted", false, "only forward drifted results")
	cmd.Flags().StringSlice("fw-services", nil, "only forward results for these services")
}

// FilterWriterOptionsFromFlags reads filter-writer options from cmd flags.
func FilterWriterOptionsFromFlags(cmd *cobra.Command) FilterWriterOptions {
	opts := DefaultFilterWriterOptions()
	if v, err := cmd.Flags().GetBool("fw-only-drifted"); err == nil {
		opts.OnlyDrifted = v
	}
	if v, err := cmd.Flags().GetStringSlice("fw-services"); err == nil {
		opts.Services = v
	}
	return opts
}
