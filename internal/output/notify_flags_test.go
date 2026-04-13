package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newNotifyCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindNotifyFlags(cmd)
	return cmd
}

func TestBindNotifyFlags_Defaults(t *testing.T) {
	cmd := newNotifyCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := NotifyOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if !opts.OnlyDrift {
		t.Error("expected OnlyDrift=true by default")
	}
	if len(opts.Channels) != 1 || opts.Channels[0] != "stdout" {
		t.Errorf("expected channels=[stdout], got %v", opts.Channels)
	}
}

func TestBindNotifyFlags_Enabled(t *testing.T) {
	cmd := newNotifyCmd()
	_ = cmd.ParseFlags([]string{"--notify"})
	opts, err := NotifyOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestBindNotifyFlags_AllEvents(t *testing.T) {
	cmd := newNotifyCmd()
	_ = cmd.ParseFlags([]string{"--notify", "--notify-only-drift=false"})
	opts, err := NotifyOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.OnlyDrift {
		t.Error("expected OnlyDrift=false")
	}
}

func TestBindNotifyFlags_CustomChannels(t *testing.T) {
	cmd := newNotifyCmd()
	_ = cmd.ParseFlags([]string{"--notify-channels=stdout,stderr"})
	opts, err := NotifyOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(opts.Channels))
	}
}
