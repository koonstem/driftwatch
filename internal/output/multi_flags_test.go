package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMultiCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindMultiFlags(cmd)
	return cmd
}

func TestBindMultiFlags_Defaults(t *testing.T) {
	cmd := newMultiCmd()
	_ = cmd.ParseFlags([]string{})
	opts := MultiOptionsFromFlags(cmdopOnError {
		t.Error("expected StopOnError to default to false")
	}
}

func TestBindMultiFlags_StopOnError(t *testing.T) {
	cmd := newMultiCmd()
	_ = cmd.ParseFlags([]string{"--multi-stop-on-error"})
	opts := MultiOptionsFromFlags(cmd)
	if !opts.StopOnError {
		t.Error("expected StopOnError to be true")
	}
}

func TestDefaultMultiOptions_Values(t *testing.T) {
	opts := DefaultMultiOptions()
	if opts.StopOnError {
		t.Error("default StopOnError should be false")
	}
}
