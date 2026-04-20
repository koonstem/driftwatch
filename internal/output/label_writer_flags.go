package output

import (
	"github.com/spf13/cobra"
)

// BindLabelWriterFlags registers label-writer flags on cmd.
func BindLabelWriterFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("label-filter", false, "Enable label-based filtering and annotation")
	cmd.Flags().Bool("label-annotate", true, "Annotate output with container labels")
	cmd.Flags().String("label-key", "", "Only include results that have this label key")
	cmd.Flags().String("label-value", "", "Only include results whose label key matches this value")
}

// LabelWriterOptionsFromFlags reads label-writer flags from cmd.
func LabelWriterOptionsFromFlags(cmd *cobra.Command) LabelWriterOptions {
	enabled, _ := cmd.Flags().GetBool("label-filter")
	annotate, _ := cmd.Flags().GetBool("label-annotate")
	key, _ := cmd.Flags().GetString("label-key")
	value, _ := cmd.Flags().GetString("label-value")

	return LabelWriterOptions{
		Enabled:     enabled,
		Annotate:    annotate,
		FilterKey:   key,
		FilterValue: value,
	}
}
