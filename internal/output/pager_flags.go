package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindPagerFlags registers pager-related flags on cmd and returns a
// function that resolves a PagerOptions from those flags.
func BindPagerFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page-size", 20, "number of results to show per page")
	cmd.Flags().Bool("no-page-info", false, "suppress page header and footer")
}

// PagerOptionsFromFlags builds PagerOptions from flags bound by BindPagerFlags.
func PagerOptionsFromFlags(cmd *cobra.Command) (PagerOptions, error) {
	pageSize, err := cmd.Flags().GetInt("page-size")
	if err != nil {
		return PagerOptions{}, fmt.Errorf("page-size flag: %w", err)
	}
	if pageSize <= 0 {
		return PagerOptions{}, fmt.Errorf("page-size must be greater than zero, got %d", pageSize)
	}

	noPageInfo, err := cmd.Flags().GetBool("no-page-info")
	if err != nil {
		return PagerOptions{}, fmt.Errorf("no-page-info flag: %w", err)
	}

	return PagerOptions{
		PageSize:     pageSize,
		ShowPageInfo: !noPageInfo,
	}, nil
}
