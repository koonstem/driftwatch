package output_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/user/driftwatch/internal/output"
)

func newBatchCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	output.BindBatchFlags(cmd)
	return cmd
}

func TestBindBatchFlags_Defaults(t *testing.T) {
	cmd := newBatchCmd()
	_ = cmd.ParseFlags([]string{})
	opts := output.BatchOptionsFromFlags(cmd)
	if opts.Enabled {
		t.Error("expected batch disabled by default")
	}
	if opts.BatchSize != 10 {
		t.Errorf("expected default batch size 10, got %d", opts.BatchSize)
	}
}

func TestBindBatchFlags_Enabled(t *testing.T) {
	cmd := newBatchCmd()
	_ = cmd.ParseFlags([]string{"--batch"})
	opts := output.BatchOptionsFromFlags(cmd)
	if !opts.Enabled {
		t.Error("expected batch enabled")
	}
}

func TestBindBatchFlags_CustomSize(t *testing.T) {
	cmd := newBatchCmd()
	_ = cmd.ParseFlags([]string{"--batch", "--batch-size=25"})
	opts := output.BatchOptionsFromFlags(cmd)
	if opts.BatchSize != 25 {
		t.Errorf("expected batch size 25, got %d", opts.BatchSize)
	}
}

func TestDefaultBatchOptions_Values(t *testing.T) {
	opts := output.DefaultBatchOptions()
	if opts.Enabled {
		t.Error("expected default disabled")
	}
	if opts.BatchSize != 10 {
		t.Errorf("expected default batch size 10, got %d", opts.BatchSize)
	}
}
