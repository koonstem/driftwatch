package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// ExportOptions controls export behaviour.
type ExportOptions struct {
	Format    string // "csv" or "json"
	Timestamp bool   // include a generated_at column/field
}

// DefaultExportOptions returns sensible defaults.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:    "csv",
		Timestamp: true,
	}
}

// NewExporter returns a function that writes results to w in the chosen format.
func NewExporter(opts ExportOptions) func(io.Writer, []drift.DriftResult) error {
	switch opts.Format {
	case "json":
		return func(w io.Writer, results []drift.DriftResult) error {
			return exportJSON(w, results, opts.Timestamp)
		}
	default:
		return func(w io.Writer, results []drift.DriftResult) error {
			return exportCSV(w, results, opts.Timestamp)
		}
	}
}

func exportCSV(w io.Writer, results []drift.DriftResult, ts bool) error {
	cw := csv.NewWriter(w)
	header := []string{"service", "status", "field", "expected", "actual"}
	if ts {
		header = append(header, "generated_at")
	}
	if err := cw.Write(header); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	for _, r := range results {
		if len(r.Fields) == 0 {
			row := []string{r.Service, string(r.Status), "", "", ""}
			if ts {
				row = append(row, now)
			}
			if err := cw.Write(row); err != nil {
				return err
			}
			continue
		}
		for _, f := range r.Fields {
			row := []string{r.Service, string(r.Status), f.Field, f.Expected, f.Actual}
			if ts {
				row = append(row, now)
			}
			if err := cw.Write(row); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportJSON(w io.Writer, results []drift.DriftResult, ts bool) error {
	type row struct {
		Service     string `json:"service"`
		Status      string `json:"status"`
		Field       string `json:"field,omitempty"`
		Expected    string `json:"expected,omitempty"`
		Actual      string `json:"actual,omitempty"`
		GeneratedAt string `json:"generated_at,omitempty"`
	}
	now := time.Now().UTC().Format(time.RFC3339)
	var rows []row
	for _, r := range results {
		if len(r.Fields) == 0 {
			entry := row{Service: r.Service, Status: string(r.Status)}
			if ts {
				entry.GeneratedAt = now
			}
			rows = append(rows, entry)
			continue
		}
		for _, f := range r.Fields {
			entry := row{Service: r.Service, Status: string(r.Status), Field: f.Field, Expected: f.Expected, Actual: f.Actual}
			if ts {
				entry.GeneratedAt = now
			}
			rows = append(rows, entry)
		}
	}
	return writeJSON(w, rows)
}

// ExportPath writes results to the named file path.
func ExportPath(path string, opts ExportOptions, results []drift.DriftResult) error {
	f, err := createFile(path)
	if err != nil {
		return fmt.Errorf("export: open %s: %w", path, err)
	}
	defer f.Close()
	return NewExporter(opts)(f, results)
}
