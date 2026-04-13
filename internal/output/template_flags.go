package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindTemplateFlags registers template-related flags onto the given command.
func BindTemplateFlags(cmd *cobra.Command) {
	cmd.Flags().String("template-file", "", "path to a Go template file for custom output")
	cmd.Flags().String("template", "", "inline Go template string for custom output")
}

// TemplateOptionsFromFlags reads template flags from the command and returns
// a TemplateOptions struct. Returns an error if both or neither source is set.
func TemplateOptionsFromFlags(cmd *cobra.Command) (TemplateOptions, error) {
	file, _ := cmd.Flags().GetString("template-file")
	inline, _ := cmd.Flags().GetString("template")

	if file != "" && inline != "" {
		return TemplateOptions{}, fmt.Errorf("template: --template-file and --template are mutually exclusive")
	}

	return TemplateOptions{
		TemplatePath: file,
		TemplateStr:  inline,
	}, nil
}

// IsTemplateEnabled returns true when any template flag has been set.
func IsTemplateEnabled(cmd *cobra.Command) bool {
	file, _ := cmd.Flags().GetString("template-file")
	inline, _ := cmd.Flags().GetString("template")
	return file != "" || inline != ""
}
