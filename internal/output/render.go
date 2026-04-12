package output

import (
	"fmt"
	"io"

	"github.com/yourorg/driftwatch/internal/drift"
)

// RenderOptions controls how drift results are rendered to output.
type RenderOptions struct {
	Format  string
	Color   bool
	Verbose bool
	Writer  io.Writer
}

// Renderer writes drift results in a configured format.
type Renderer struct {
	opts      RenderOptions
	colorizer *Colorizer
}

// NewRenderer creates a Renderer with the given options.
func NewRenderer(opts RenderOptions) *Renderer {
	return &Renderer{
		opts:      opts,
		colorizer: NewColorizer(opts.Color),
	}
}

// Render writes the drift report to the configured writer in the chosen format.
func (r *Renderer) Render(report *drift.Report) error {
	switch r.opts.Format {
	case "json":
		return writeJSON(r.opts.Writer, report)
	case "table":
		return writeTable(r.opts.Writer, report, r.colorizer)
	case "diff":
		w := NewDiffWriter(r.opts.Writer, r.colorizer)
		w.Write(report)
		return nil
	case "summary":
		sw := NewSummaryWriter(r.opts.Writer, r.colorizer)
		sw.Write(report)
		return nil
	case "text", "":
		f := NewFormatter(r.opts.Writer, r.colorizer)
		f.Write(report)
		return nil
	default:
		return fmt.Errorf("unknown output format: %q", r.opts.Format)
	}
}
