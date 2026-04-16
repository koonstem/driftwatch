package output

import "github.com/spf13/cobra"

// PipelineFlagOptions holds CLI-parsed pipeline settings.
type PipelineFlagOptions struct {
	StopOnError bool
}

// BindPipelineFlags registers pipeline flags on a cobra command.
func BindPipelineFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("pipeline-stop-on-error", false, "stop pipeline execution on first writer error")
}

// PipelineOptionsFromFlags reads pipeline options from a cobra command.
func PipelineOptionsFromFlags(cmd *cobra.Command) PipelineOptions {
	stop, _ := cmd.Flags().GetBool("pipeline-stop-on-error")
	opts := DefaultPipelineOptions()
	opts.StopOnError = stop
	return opts
}
