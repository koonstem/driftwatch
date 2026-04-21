package output

import (
	"time"

	"github.com/spf13/cobra"
)

// BindCircuitFlags attaches circuit breaker flags to the given command.
func BindCircuitFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("circuit-enabled", false, "enable circuit breaker for downstream writers")
	cmd.Flags().Int("circuit-max-failures", DefaultCircuitOptions().MaxFailures, "consecutive failures before opening the circuit")
	cmd.Flags().Duration("circuit-reset-timeout", DefaultCircuitOptions().ResetTimeout, "duration to wait before attempting half-open probe")
}

// CircuitOptionsFromFlags reads circuit breaker options from parsed cobra flags.
func CircuitOptionsFromFlags(cmd *cobra.Command) CircuitOptions {
	opts := DefaultCircuitOptions()

	if v, err := cmd.Flags().GetBool("circuit-enabled"); err == nil {
		opts.Enabled = v
	}
	if v, err := cmd.Flags().GetInt("circuit-max-failures"); err == nil {
		opts.MaxFailures = v
	}
	if v, err := cmd.Flags().GetDuration("circuit-reset-timeout"); err == nil {
		opts.ResetTimeout = v
	}

	return opts
}

// IsCircuitEnabled returns true when the circuit-enabled flag is set.
func IsCircuitEnabled(cmd *cobra.Command) bool {
	v, err := cmd.Flags().GetBool("circuit-enabled")
	if err != nil {
		return false
	}
	return v
}

// DefaultCircuitDuration is a convenience alias used in tests.
const DefaultCircuitDuration = 30 * time.Second
