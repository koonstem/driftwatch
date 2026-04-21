package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newRollupCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindRollupFlags(cmd)
	return cmd
}

func TestBindRollupFlags_Defaults(t *testing.T) {
	cmd := newRollupCmd()
	_ = cmd.ParseFlags([]string{})
	opts, err := RollupOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected rollup disabled by default")
	}
	if opts.GroupBy != "service" {
		t.Errorf("expected default group_by 'service', got %q", opts.GroupBy)
	}
	if opts.TopN != 10 {
		t.Errorf("expected default top-n 10, got %d", opts.TopN)
	}
}

func TestBindRollupFlags_Enabled(t *testing.T) {
	cmd := newRollupCmd()
	_ = cmd.ParseFlags([]string{"--rollup"})
	opts, err := RollupOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Enabled {
		t.Error("expected rollup enabled")
	}
}

func TestBindRollupFlags_ByField(t *testing.T) {
	cmd := newRollupCmd()
	_ = cmd.ParseFlags([]string{"--rollup-group-by", "field"})
	opts, err := RollupOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.GroupBy != "field" {
		t.Errorf("expected group_by 'field', got %q", opts.GroupBy)
	}
}

func TestBindRollupFlags_CustomTopN(t *testing.T) {
	cmd := newRollupCmd()
	_ = cmd.ParseFlags([]string{"--rollup-top-n", "5"})
	opts, err := RollupOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.TopN != 5 {
		t.Errorf("expected top-n 5, got %d", opts.TopN)
	}
}

func TestDefaultRollupOptions_Values(t *testing.T) {
	opts := DefaultRollupOptions()
	if opts.Enabled {
		t.Error("default should be disabled")
	}
	if opts.GroupBy != "service" {
		t.Errorf("default group_by should be 'service', got %q", opts.GroupBy)
	}
	if opts.TopN != 10 {
		t.Errorf("default top-n should be 10, got %d", opts.TopN)
	}
}
