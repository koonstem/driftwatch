package output

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// BindReplayFlags attaches replay-related flags to a cobra command.
func BindReplayFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("replay", false, "replay drift results from a history file")
	cmd.Flags().String("replay-file", "", "path to history file to replay")
	cmd.Flags().Duration("replay-delay", 0, "delay between replayed snapshots (e.g. 500ms)")
}

// ReplayOptionsFromFlags builds ReplayOptions from parsed cobra flags.
func ReplayOptionsFromFlags(cmd *cobra.Command) (ReplayOptions, error) {
	enabled, err := cmd.Flags().GetBool("replay")
	if err != nil {
		return ReplayOptions{}, fmt.Errorf("replay flag: %w", err)
	}
	file, err := cmd.Flags().GetString("replay-file")
	if err != nil {
		return ReplayOptions{}, fmt.Errorf("replay-file flag: %w", err)
	}
	delay, err := cmd.Flags().GetDuration("replay-delay")
	if err != nil {
		return ReplayOptions{}, fmt.Errorf("replay-delay flag: %w", err)
	}
	if enabled && file == "" {
		return ReplayOptions{}, fmt.Errorf("--replay-file is required when --replay is set")
	}
	return ReplayOptions{
		Enabled: enabled,
		File:    file,
		Delay:   delay,
	}, nil
}

// IsReplayEnabled returns true when the replay flag is set.
func IsReplayEnabled(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("replay")
	return v
}

// ValidateReplayDelay ensures the delay is non-negative.
func ValidateReplayDelay(d time.Duration) error {
	if d < 0 {
		return fmt.Errorf("replay-delay must be non-negative, got %s", d)
	}
	return nil
}
