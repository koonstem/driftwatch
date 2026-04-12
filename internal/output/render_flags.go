package output

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// RenderFlags holds CLI flag values that control output rendering.
type RenderFlags struct {
	Format  string
	Color   bool
	Verbose bool
}

// BindRenderFlags registers output-related flags on a cobra command.
func BindRenderFlags(cmd *cobra.Command, flags *RenderFlags) {
	cmd.Flags().StringVarP(
		&flags.Format, "output", "o", "text",
		`Output format: text, json, table, diff, summary`,
	)
	cmd.Flags().BoolVar(
		&flags.Color, "color", true,
		"Enable color output (disable with --color=false)",
	)
	cmd.Flags().BoolVarP(
		&flags.Verbose, "verbose", "v", false,
		"Enable verbose output",
	)
}

// RendererFromFlags builds a Renderer from parsed flag values, writing to w.
// If w is nil, os.Stdout is used.
func RendererFromFlags(flags RenderFlags, w io.Writer) *Renderer {
	if w == nil {
		w = os.Stdout
	}
	return NewRenderer(RenderOptions{
		Format:  flags.Format,
		Color:   flags.Color,
		Verbose: flags.Verbose,
		Writer:  w,
	})
}
