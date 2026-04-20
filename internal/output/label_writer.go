package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yourorg/driftwatch/internal/drift"
)

// LabelWriterOptions controls label-based annotation of drift results.
type LabelWriterOptions struct {
	Enabled    bool
	Annotate   bool   // inject labels as extra fields in output
	FilterKey  string // only include results that have this label key
	FilterValue string // only include results whose label key matches this value
}

// DefaultLabelWriterOptions returns safe defaults.
func DefaultLabelWriterOptions() LabelWriterOptions {
	return LabelWriterOptions{
		Enabled:  false,
		Annotate: true,
	}
}

// LabelWriter filters and/or annotates drift results by container labels.
type LabelWriter struct {
	opts  LabelWriterOptions
	next  Writer
}

// NewLabelWriter constructs a LabelWriter that wraps next.
func NewLabelWriter(opts LabelWriterOptions, next Writer) *LabelWriter {
	return &LabelWriter{opts: opts, next: next}
}

// Write applies label filtering/annotation before delegating to the inner writer.
func (lw *LabelWriter) Write(w io.Writer, report drift.Report) error {
	if !lw.opts.Enabled {
		return lw.next.Write(w, report)
	}

	filtered := make([]drift.Result, 0, len(report.Results))
	for _, r := range report.Results {
		if lw.opts.FilterKey != "" {
			val, ok := r.Labels[lw.opts.FilterKey]
			if !ok {
				continue
			}
			if lw.opts.FilterValue != "" && val != lw.opts.FilterValue {
				continue
			}
		}
		if lw.opts.Annotate && len(r.Labels) > 0 {
			r = annotateResult(r)
		}
		filtered = append(filtered, r)
	}

	report.Results = filtered
	return lw.next.Write(w, report)
}

// annotateResult appends sorted label key=value pairs as an extra DriftField.
func annotateResult(r drift.Result) drift.Result {
	keys := make([]string, 0, len(r.Labels))
	for k := range r.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, r.Labels[k]))
	}

	r.Fields = append(r.Fields, drift.DriftField{
		Field:    "labels",
		Expected: "",
		Actual:   strings.Join(pairs, ", "),
	})
	return r
}
