package output_test

import (
	"testing"

	"github.com/driftwatch/internal/output"
	"github.com/spf13/cobra"
)

func newFilterWriterCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	output.BindFilterWriterFlags(cmd)
	return cmd
}

func TestBindFilterWriterFlags_Defaults(t *testing.T) {
	cmd := newFilterWriterCmd()
	_ = cmd.ParseFlags([]string{})
	opts := output.FilterWriterOptionsFromFlags(cmd)
	if opts.OnlyDrifted {
		t.Error("expected OnlyDrifted false by default")
	}
	if len(opts.Services) != 0 {
		t.Errorf("expected empty services, got %v", opts.Services)
	}
}

func TestBindFilterWriterFlags_OnlyDrifted(t *testing.T) {
	cmd := newFilterWriterCmd()
	_ = cmd.ParseFlags([]string{"--fw-only-drifted"})
	opts := output.FilterWriterOptionsFromFlags(cmd)
	if !opts.OnlyDrifted {
		t.Error("expected OnlyDrifted true")
	}
}

func TestBindFilterWriterFlags_Services(t *testing.T) {
	cmd := newFilterWriterCmd()
	_ = cmd.ParseFlags([]string{"--fw-services", "alpha,gamma"})
	opts := output.FilterWriterOptionsFromFlags(cmd)
	if len(opts.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(opts.Services))
	}
	if opts.Services[0] != "alpha" || opts.Services[1] != "gamma" {
		t.Errorf("unexpected services: %v", opts.Services)
	}
}
