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
	if report == nil {
		return fmt.Errorf("cannot write nil report")
	}
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

// ValidFormats returns all supported Format values.
func ValidFormats() []Format {
	return []Format{FormatText, FormatJSON, FormatTable}
}

// ParseFormat converts a string to a Format, returning an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	f := Format(s)
	for _, v := range ValidFormats() {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("unknown format %q: valid options are text, json, table", s)
}
