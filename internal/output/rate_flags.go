package output

import (
	"github.com/spf13/cobra"
)

// BindRateFlags registers rate-limiting flags on cmd.
func BindRateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("rate-limit", false, "Enable output rate limiting")
	cmd.Flags().Int("rate-max-per-min", 60, "Maximum writes per minute")
	cmd.Flags().Int("rate-burst", 5, "Burst size for rate limiter")
}

// RateOptionsFromFlags reads rate options from parsed cobra flags.
func RateOptionsFromFlags(cmd *cobra.Command) RateOptions {
	opts := DefaultRateOptions()
	if v, err := cmd.Flags().GetBool("rate-limit"); err == nil {
		opts.Enabled = v
	}
	if v, err := cmd.Flags().GetInt("rate-max-per-min"); err == nil && v > 0 {
		opts.MaxPerMin = v
	}
	if v, err := cmd.Flags().GetInt("rate-burst"); err == nil && v > 0 {
		opts.BurstSize = v
	}
	return opts
}
