package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newPipelineCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindPipelineFlags(cmd)
	return cmd
}

func TestBindPipelineFlags_Defaults(t *testing.T) {
	cmd := newPipelineCmd()
	_ = cmd.ParseFlags([]string{})
	opts := PipelineOptionsFromFlags(cmd)
	if opts.StopOnError {
		t.Error("expected StopOnError=false by default")
	}
}

func TestBindPipelineFlags_StopOnError(t *testing.T) {
	cmd := newPipelineCmd()
	_ = cmd.ParseFlags([]string{"--pipeline-stop-on-error"})
	opts := PipelineOptionsFromFlags(cmd)
	if !opts.StopOnError {
		t.Error("expected StopOnError=true")
	}
}

func TestDefaultPipelineOptions_Values(t *testing.T) {
	opts := DefaultPipelineOptions()
	if opts.StopOnError {
		t.Error("default StopOnError should be false")
	}
	if len(opts.Writers) != 0 {
		t.Error("default Writers should be empty")
	}
}
