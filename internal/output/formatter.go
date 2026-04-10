package output

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes drift reports to an output stream.
type Formatter struct {
	format Format
	out    io.Writer
}

// NewFormatter creates a Formatter writing to the given writer.
// If w is nil, os.Stdout is used.
func NewFormatter(format Format, w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, out: w}
}

// Write renders the report according to the configured format.
func (f *Formatter) Write(report *drift.Report) error {
	switch f.format {
	case FormatJSON:
		return writeJSON(f.out, report)
	default:
		return writeText(f.out, report)
	}
}

func writeText(w io.Writer, report *drift.Report) error {
	if !report.HasDrift() {
		fmt.Fprintln(w, "✓ No drift detected.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tFIELD\tEXPECTED\tACTUAL")
	for _, entry := range report.Entries {
		if entry.Drifted {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				entry.ServiceName, entry.Field, entry.Expected, entry.Actual)
		}
	}
	return tw.Flush()
}
