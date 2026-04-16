package output

import "github.com/spf13/cobra"

// BindAggregateFlags registers aggregate flags on a cobra command.
func BindAggregateFlags(cmd *cobra.Command) {
	cmd.Flags().String("aggregate-by", "service", `Group drift results by field: "service" or "field"`)
}

// AggregateOptionsFromFlags reads aggregate options from cobra flags.
func AggregateOptionsFromFlags(cmd *cobra.Command) (AggregateOptions, error) {
	groupBy, err := cmd.Flags().GetString("aggregate-by")
	if err != nil {
		return DefaultAggregateOptions(), err
	}
	return AggregateOptions{GroupBy: groupBy}, nil
}
