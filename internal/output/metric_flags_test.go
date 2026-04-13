package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMetricCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindMetricFlags(cmd)
	return cmd
}

func TestBindMetricFlags_Defaults(t *testing.T) {
	cmd := newMetricCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := MetricOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected metrics disabled by default")
	}
	if !opts.Timestamp {
		t.Error("expected timestamp enabled by default")
	}
	if opts.Prefix != "driftwatch" {
		t.Errorf("expected default prefix 'driftwatch', got %q", opts.Prefix)
	}
}

func TestBindMetricFlags_Enabled(t *testing.T) {
	cmd := newMetricCmd()
	_ = cmd.ParseFlags([]string{"--metrics"})
	opts, err := MetricOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected metrics enabled")
	}
}

func TestBindMetricFlags_CustomPrefix(t *testing.T) {
	cmd := newMetricCmd()
	_ = cmd.ParseFlags([]string{"--metrics-prefix", "myservice"})
	opts, err := MetricOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Prefix != "myservice" {
		t.Errorf("expected prefix 'myservice', got %q", opts.Prefix)
	}
}

func TestBindMetricFlags_NoTimestamp(t *testing.T) {
	cmd := newMetricCmd()
	_ = cmd.ParseFlags([]string{"--metrics-timestamp=false"})
	opts, err := MetricOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Timestamp {
		t.Error("expected timestamp disabled")
	}
}
