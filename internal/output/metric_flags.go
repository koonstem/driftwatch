package output

import "github.com/spf13/cobra"

// BindMetricFlags registers metric-related flags onto cmd.
func BindMetricFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("metrics", false, "emit Prometheus-style text metrics")
	cmd.Flags().Bool("metrics-timestamp", true, "include Unix-ms timestamp in metric output")
	cmd.Flags().String("metrics-prefix", "driftwatch", "prefix for metric names")
}

// MetricOptionsFromFlags reads metric flags from cmd and returns MetricOptions.
func MetricOptionsFromFlags(cmd *cobra.Command) (MetricOptions, error) {
	opts := DefaultMetricOptions()

	enabled, err := cmd.Flags().GetBool("metrics")
	if err != nil {
		return opts, err
	}
	opts.Enabled = enabled

	ts, err := cmd.Flags().GetBool("metrics-timestamp")
	if err != nil {
		return opts, err
	}
	opts.Timestamp = ts

	prefix, err := cmd.Flags().GetString("metrics-prefix")
	if err != nil {
		return opts, err
	}
	if prefix != "" {
		opts.Prefix = prefix
	}

	return opts, nil
}
