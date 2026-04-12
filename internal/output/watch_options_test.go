package output

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	return &cobra.Command{Use: "test"}
}

func TestBindWatchFlags_Defaults(t *testing.T) {
	cmd := newWatchCmd()
	var opts WatchOptions
	BindWatchFlags(cmd, &opts)

	if err := cmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if opts.Enabled {
		t.Error("expected watch disabled by default")
	}
	if opts.Interval != 5*time.Second {
		t.Errorf("expected default interval 5s, got %v", opts.Interval)
	}
}

func TestBindWatchFlags_Enabled(t *testing.T) {
	cmd := newWatchCmd()
	var opts WatchOptions
	BindWatchFlags(cmd, &opts)

	if err := cmd.ParseFlags([]string{"--watch", "--watch-interval", "15s"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	if !opts.Enabled {
		t.Error("expected watch enabled")
	}
	if opts.Interval != 15*time.Second {
		t.Errorf("expected interval 15s, got %v", opts.Interval)
	}
}

func TestWatchOptions_Validate_OK(t *testing.T) {
	opts := WatchOptions{Enabled: true, Interval: 10 * time.Second}
	if err := opts.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestWatchOptions_Validate_ZeroInterval(t *testing.T) {
	opts := WatchOptions{Enabled: true, Interval: 0}
	if err := opts.Validate(); err == nil {
		t.Error("expected validation error for zero interval")
	}
}

func TestWatchOptions_Validate_DisabledIgnoresInterval(t *testing.T) {
	opts := WatchOptions{Enabled: false, Interval: 0}
	if err := opts.Validate(); err != nil {
		t.Errorf("disabled watch should not validate interval, got: %v", err)
	}
}
