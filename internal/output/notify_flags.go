package output

import "github.com/spf13/cobra"

// BindNotifyFlags registers notification flags on cmd.
func BindNotifyFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("notify", false, "enable drift notifications")
	cmd.Flags().Bool("notify-only-drift", true, "only notify when drift is detected")
	cmd.Flags().StringSlice("notify-channels", []string{"stdout"}, "notification channels (stdout)")
}

// NotifyOptionsFromFlags builds NotifyOptions from parsed cobra flags.
func NotifyOptionsFromFlags(cmd *cobra.Command) (NotifyOptions, error) {
	enabled, err := cmd.Flags().GetBool("notify")
	if err != nil {
		return NotifyOptions{}, err
	}
	onlyDrift, err := cmd.Flags().GetBool("notify-only-drift")
	if err != nil {
		return NotifyOptions{}, err
	}
	channels, err := cmd.Flags().GetStringSlice("notify-channels")
	if err != nil {
		return NotifyOptions{}, err
	}
	return NotifyOptions{
		Enabled:   enabled,
		OnlyDrift: onlyDrift,
		Channels:  channels,
	}, nil
}
