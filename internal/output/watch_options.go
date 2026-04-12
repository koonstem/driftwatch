package output

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// WatchOptions holds configuration for watch mode.
type WatchOptions struct {
	Enabled  bool
	Interval time.Duration
}

// BindWatchFlags registers --watch and --watch-interval flags on cmd.
func BindWatchFlags(cmd *cobra.Command, opts *WatchOptions) {
	cmd.Flags().BoolVar(&opts.Enabled, "watch", false, "re-run detection on an interval and refresh output")
	cmd.Flags().DurationVar(&opts.Interval, "watch-interval", 5*time.Second, "how often to refresh in watch mode (e.g. 10s, 1m)")
}

// Validate returns an error if the options are inconsistent.
func (o *WatchOptions) Validate() error {
	if o.Enabled && o.Interval <= 0 {
		return fmt.Errorf("--watch-interval must be a positive duration")
	}
	return nil
}
