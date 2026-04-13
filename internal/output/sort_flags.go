package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindSortFlags attaches sort-related flags to a cobra command.
func BindSortFlags(cmd *cobra.Command) {
	cmd.Flags().String("sort-by", string(SortByService), `sort results by field: "service", "status", or "field"`)
	cmd.Flags().Bool("sort-reverse", false, "reverse the sort order")
}

// SortOptionsFromFlags reads sort options from a cobra command's flags.
func SortOptionsFromFlags(cmd *cobra.Command) (SortOptions, error) {
	byStr, err := cmd.Flags().GetString("sort-by")
	if err != nil {
		return SortOptions{}, fmt.Errorf("reading --sort-by: %w", err)
	}

	reverse, err := cmd.Flags().GetBool("sort-reverse")
	if err != nil {
		return SortOptions{}, fmt.Errorf("reading --sort-reverse: %w", err)
	}

	field := SortField(byStr)
	switch field {
	case SortByService, SortByStatus, SortByField:
		// valid
	default:
		return SortOptions{}, fmt.Errorf("unknown --sort-by value %q: must be service, status, or field", byStr)
	}

	return SortOptions{By: field, Reverse: reverse}, nil
}
