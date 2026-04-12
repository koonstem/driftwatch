package output

import (
	"fmt"
	"io"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// DiffWriter writes a unified-diff-style view of drifted fields.
type DiffWriter struct {
	out      io.Writer
	colorize *Colorizer
}

// NewDiffWriter returns a DiffWriter that emits diff output to w.
func NewDiffWriter(w io.Writer, color bool) *DiffWriter {
	return &DiffWriter{
		out:      w,
		colorize: NewColorizer(color),
	}
}

// Write emits a diff-style block for every drifted result in the report.
func (d *DiffWriter) Write(report drift.Report) {
	if !report.HasDrift() {
		fmt.Fprintln(d.out, d.colorize.Green("✔ no drift detected"))
		return
	}

	for _, r := range report.Results {
		if !r.Drifted {
			continue
		}

		header := fmt.Sprintf("--- declared: %s", r.ServiceName)
		actual := fmt.Sprintf("+++ running:  %s", r.ContainerName)
		fmt.Fprintln(d.out, d.colorize.Bold(header))
		fmt.Fprintln(d.out, d.colorize.Bold(actual))

		for _, field := range r.DriftedFields {
			removed := fmt.Sprintf("- %s: %s", field.Field, field.Expected)
			added := fmt.Sprintf("+ %s: %s", field.Field, field.Actual)
			fmt.Fprintln(d.out, d.colorize.Red(removed))
			fmt.Fprintln(d.out, d.colorize.Green(added))
		}

		fmt.Fprintln(d.out)
	}
}
