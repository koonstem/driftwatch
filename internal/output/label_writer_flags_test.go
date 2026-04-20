package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newLabelWriterCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindLabelWriterFlags(cmd)
	return cmd
}

func TestBindLabelWriterFlags_Defaults(t *testing.T) {
	cmd := newLabelWriterCmd()
	_ = cmd.ParseFlags([]string{})
	opts := LabelWriterOptionsFromFlags(cmd)

	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if !opts.Annotate {
		t.Error("expected Annotate=true by default")
	}
	if opts.FilterKey != "" {
		t.Errorf("expected empty FilterKey, got %q", opts.FilterKey)
	}
	if opts.FilterValue != "" {
		t.Errorf("expected empty FilterValue, got %q", opts.FilterValue)
	}
}

func TestBindLabelWriterFlags_Enabled(t *testing.T) {
	cmd := newLabelWriterCmd()
	_ = cmd.ParseFlags([]string{"--label-filter"})
	opts := LabelWriterOptionsFromFlags(cmd)

	if !opts.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestBindLabelWriterFlags_KeyAndValue(t *testing.T) {
	cmd := newLabelWriterCmd()
	_ = cmd.ParseFlags([]string{"--label-key", "env", "--label-value", "prod"})
	opts := LabelWriterOptionsFromFlags(cmd)

	if opts.FilterKey != "env" {
		t.Errorf("expected FilterKey='env', got %q", opts.FilterKey)
	}
	if opts.FilterValue != "prod" {
		t.Errorf("expected FilterValue='prod', got %q", opts.FilterValue)
	}
}

func TestBindLabelWriterFlags_AnnotateFalse(t *testing.T) {
	cmd := newLabelWriterCmd()
	_ = cmd.ParseFlags([]string{"--label-annotate=false"})
	opts := LabelWriterOptionsFromFlags(cmd)

	if opts.Annotate {
		t.Error("expected Annotate=false")
	}
}

func TestDefaultLabelWriterOptions_Values(t *testing.T) {
	opts := DefaultLabelWriterOptions()
	if opts.Enabled {
		t.Error("default Enabled should be false")
	}
	if !opts.Annotate {
		t.Error("default Annotate should be true")
	}
}
