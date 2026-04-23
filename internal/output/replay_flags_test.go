package output

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newReplayCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindReplayFlags(cmd)
	return cmd
}

func TestBindReplayFlags_Defaults(t *testing.T) {
	cmd := newReplayCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := ReplayOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.File != "" {
		t.Errorf("expected empty file, got %q", opts.File)
	}
	if opts.Delay != 0 {
		t.Errorf("expected zero delay, got %s", opts.Delay)
	}
}

func TestBindReplayFlags_EnabledRequiresFile(t *testing.T) {
	cmd := newReplayCmd()
	_ = cmd.ParseFlags([]string{"--replay"})
	_, err := ReplayOptionsFromFlags(cmd)
	if err == nil {
		t.Error("expected error when --replay set without --replay-file")
	}
}

func TestBindReplayFlags_FullOptions(t *testing.T) {
	cmd := newReplayCmd()
	_ = cmd.ParseFlags([]string{"--replay", "--replay-file", "/tmp/h.json", "--replay-delay", "200ms"})
	opts, err := ReplayOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
	if opts.File != "/tmp/h.json" {
		t.Errorf("file: got %q, want %q", opts.File, "/tmp/h.json")
	}
	if opts.Delay != 200*time.Millisecond {
		t.Errorf("delay: got %s, want 200ms", opts.Delay)
	}
}

func TestIsReplayEnabled_False(t *testing.T) {
	cmd := newReplayCmd()
	_ = cmd.ParseFlags([]string{})
	if IsReplayEnabled(cmd) {
		t.Error("expected false")
	}
}

func TestIsReplayEnabled_True(t *testing.T) {
	cmd := newReplayCmd()
	_ = cmd.ParseFlags([]string{"--replay", "--replay-file", "x.json"})
	if !IsReplayEnabled(cmd) {
		t.Error("expected true")
	}
}

func TestValidateReplayDelay_Negative(t *testing.T) {
	if err := ValidateReplayDelay(-1 * time.Second); err == nil {
		t.Error("expected error for negative delay")
	}
}

func TestValidateReplayDelay_Zero(t *testing.T) {
	if err := ValidateReplayDelay(0); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
