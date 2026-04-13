package output

import (
	"fmt"
	"io"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// RenderFormat enumerates the supported output formats.
type RenderFormat string

const (
	FormatText     RenderFormat = "text"
	FormatJSON     RenderFormat = "json"
	FormatTable    RenderFormat = "table"
	FormatSummary  RenderFormat = "summary"
	FormatDiff     RenderFormat = "diff"
	FormatTemplate RenderFormat = "template"
)

// RenderOptions configures the Renderer.
type RenderOptions struct {
	Format   RenderFormat
	Color    bool
	Template TemplateOptions
}

// Renderer dispatches drift results to the appropriate output writer.
type Renderer struct {
	w    io.Writer
	opts RenderOptions
}

// NewRenderer creates a Renderer with the given options.
func NewRenderer(w io.Writer, opts RenderOptions) *Renderer {
	return &Renderer{w: w, opts: opts}
}

// Render writes results in the configured format.
func (r *Renderer) Render(results []drift.Result) error {
	switch r.opts.Format {
	case FormatJSON:
		return writeJSON(r.w, results)
	case FormatTable:
		return writeTable(r.w, results, r.opts.Color)
	case FormatSummary:
		sw := NewSummaryWriter(r.w)
		return sw.Write(results)
	case FormatDiff:
		dw := NewDiffWriter(r.w, r.opts.Color)
		return dw.Write(results)
	case FormatTemplate:
		tw, err := NewTemplateWriter(r.w, r.opts.Template)
		if err != nil {
			return err
		}
		return tw.Write(results)
	case FormatText, "":
		return writeText(r.w, results, r.opts.Color)
	default:
		return fmt.Errorf("renderer: unknown format %q", r.opts.Format)
	}
}
