package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindSampleFlags registers sample-related flags onto cmd.
func BindSampleFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("sample", false, "enable probabilistic result sampling")
	cmd.Flags().Float64("sample-rate", 1.0, "fraction of results to forward (0.0–1.0)")
}

// SampleOptionsFromFlags reads sample options from cmd flags.
func SampleOptionsFromFlags(cmd *cobra.Command) (SampleOptions, error) {
	opts := DefaultSampleOptions()

	enabled, err := cmd.Flags().GetBool("sample")
	if err != nil {
		return opts, fmt.Errorf("reading --sample: %w", err)
	}
	opts.Enabled = enabled

	rate, err := cmd.Flags().GetFloat64("sample-rate")
	if err != nil {
		return opts, fmt.Errorf("reading --sample-rate: %w", err)
	}
	if rate < 0 || rate > 1 {
		return opts, fmt.Errorf("--sample-rate must be between 0.0 and 1.0, got %f", rate)
	}
	opts.Rate = rate

	return opts, nil
}
