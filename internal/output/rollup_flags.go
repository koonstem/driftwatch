package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindRollupFlags attaches rollup-related flags to cmd.
func BindRollupFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("rollup", false, "enable drift rollup summary")
	cmd.Flags().String("rollup-group-by", "service", "group rollup by 'service' or 'field'")
	cmd.Flags().Int("rollup-top-n", 10, "number of top entries to include in rollup")
}

// RollupOptionsFromFlags reads rollup flags from cmd and returns RollupOptions.
func RollupOptionsFromFlags(cmd *cobra.Command) (RollupOptions, error) {
	enabled, err := cmd.Flags().GetBool("rollup")
	if err != nil {
		return RollupOptions{}, fmt.Errorf("rollup: %w", err)
	}
	groupBy, err := cmd.Flags().GetString("rollup-group-by")
	if err != nil {
		return RollupOptions{}, fmt.Errorf("rollup-group-by: %w", err)
	}
	topN, err := cmd.Flags().GetInt("rollup-top-n")
	if err != nil {
		return RollupOptions{}, fmt.Errorf("rollup-top-n: %w", err)
	}
	return RollupOptions{
		Enabled: enabled,
		GroupBy: groupBy,
		TopN:    topN,
	}, nil
}
