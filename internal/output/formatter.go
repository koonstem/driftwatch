package output

import (
	"fmt"
	"io"
	"os"

	"github.com/driftwatch/internal/drift"
)

// Format enumerates supported output formats.
type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

// Formatter writes drift reports in the requested format.
type Formatter struct {
	format Format
	out    io.Writer
}

// NewFormatter creates a Formatter writing to w (defaults to os.Stdout).
func NewFormatter(format Format, w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, out: w}
}

// Write serialises the report to the configured output.
func (f *Formatter) Write(report *drift.Report) error {
	switch f.format {
	case FormatJSON:
		return writeJSON(f.out, report)
	case FormatTable:
		return writeTable(f.out, report)
	case FormatText:
		return writeText(f.out, report)
	default:
		return fmt.Errorf("unsupported format: %q", f.format)
	}
}
