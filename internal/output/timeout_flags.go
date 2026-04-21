package output

import (
	"time"

	"github.com/spf13/cobra"
)

// BindTimeoutFlags registers timeout-related flags on cmd.
func BindTimeoutFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("timeout", false, "enable write timeout for downstream writers")
	cmd.Flags().Duration("timeout-duration", 10*time.Second, "maximum duration allowed for a write operation")
}

// TimeoutOptionsFromFlags reads the registered timeout flags from cmd.
func TimeoutOptionsFromFlags(cmd *cobra.Command) TimeoutOptions {
	enabled, _ := cmd.Flags().GetBool("timeout")
	duration, _ := cmd.Flags().GetDuration("timeout-duration")
	return TimeoutOptions{
		Enabled:  enabled,
		Duration: duration,
	}
}
