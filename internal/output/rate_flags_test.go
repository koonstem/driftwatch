package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newRateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindRateFlags(cmd)
	return cmd
}

func TestBindRateFlags_Defaults(t *testing.T) {
	cmd := newRateCmd()
	_ = cmd.ParseFlags([]string{})
	opts := RateOptionsFromFlags(cmd)
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.MaxPerMin != 60 {
		t.Errorf("expected MaxPerMin=60, got %d", opts.MaxPerMin)
	}
	if opts.BurstSize != 5 {
		t.Errorf("expected BurstSize=5, got %d", opts.BurstSize)
	}
}

func TestBindRateFlags_Enabled(t *testing.T) {
	cmd := newRateCmd()
	_ = cmd.ParseFlags([]string{"--rate-limit"})
	opts := RateOptionsFromFlags(cmd)
	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestBindRateFlags_CustomValues(t *testing.T) {
	cmd := newRateCmd()
	_ = cmd.ParseFlags([]string{"--rate-limit", "--rate-max-per-min=30", "--rate-burst=10"})
	opts := RateOptionsFromFlags(cmd)
	if opts.MaxPerMin != 30 {
		t.Errorf("expected MaxPerMin=30, got %d", opts.MaxPerMin)
	}
	if opts.BurstSize != 10 {
		t.Errorf("expected BurstSize=10, got %d", opts.BurstSize)
	}
}

func TestDefaultRateOptions_Values(t *testing.T) {
	opts := DefaultRateOptions()
	if opts.Enabled {
		t.Error("default Enabled should be false")
	}
	if opts.MaxPerMin <= 0 {
		t.Error("default MaxPerMin should be positive")
	}
	if opts.BurstSize <= 0 {
		t.Error("default BurstSize should be positive")
	}
}
