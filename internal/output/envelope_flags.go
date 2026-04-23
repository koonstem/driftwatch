package output

import (
	"github.com/spf13/cobra"
)

// BindEnvelopeFlags attaches envelope-related flags to a cobra command.
func BindEnvelopeFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("envelope", false, "wrap results with run metadata envelope")
	cmd.Flags().String("envelope-source", "", "source identifier to include in envelope")
	cmd.Flags().String("envelope-run-id", "", "run ID to include in envelope (defaults to empty)")
}

// EnvelopeOptionsFromFlags builds EnvelopeOptions from parsed cobra flags.
func EnvelopeOptionsFromFlags(cmd *cobra.Command) EnvelopeOptions {
	opts := DefaultEnvelopeOptions()

	if v, err := cmd.Flags().GetBool("envelope"); err == nil {
		opts.Enabled = v
	}
	if v, err := cmd.Flags().GetString("envelope-source"); err == nil {
		opts.Source = v
	}
	if v, err := cmd.Flags().GetString("envelope-run-id"); err == nil {
		opts.RunID = v
	}

	return opts
}
