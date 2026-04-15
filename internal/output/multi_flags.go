package output

import "github.com/spf13/cobra"

// BindMultiFlags registers multi-writer flags on cmd.
func BindMultiFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("multi-stop-on-error", false, "stop fan-out to subsequent writers after the first error")
}

// MultiOptionsFromFlags reads flag values from cmd and returns MultiOptions.
func MultiOptionsFromFlags(cmd *cobra.Command) MultiOptions {
	opts := DefaultMultiOptions()
	if v, err := cmd.Flags().GetBool("multi-stop-on-error"); err == nil {
		opts.StopOnError = v
	}
	return opts
}
