package output

import "github.com/driftwatch/internal/drift"

// FilterWriterOptions controls which results are forwarded to the inner writer.
type FilterWriterOptions struct {
	OnlyDrifted bool
	Services    []string
}

// DefaultFilterWriterOptions returns sensible defaults.
func DefaultFilterWriterOptions() FilterWriterOptions {
	return FilterWriterOptions{}
}

// FilterWriter wraps another Writer and drops results that do not match the
// configured criteria before forwarding.
type FilterWriter struct {
	opts  FilterWriterOptions
	inner Writer
}

// NewFilterWriter creates a FilterWriter that forwards matching results to inner.
func NewFilterWriter(opts FilterWriterOptions, inner Writer) *FilterWriter {
	return &FilterWriter{opts: opts, inner: inner}
}

// Write filters results then delegates to the inner writer.
func (f *FilterWriter) Write(report drift.Report) error {
	filtered := make([]drift.Result, 0, len(report.Results))
	for _, r := range report.Results {
		if f.opts.OnlyDrifted && !r.Drifted {
			continue
		}
		if len(f.opts.Services) > 0 && !containsService(f.opts.Services, r.Service) {
			continue
		}
		filtered = append(filtered, r)
	}
	report.Results = filtered
	return f.inner.Write(report)
}

func containsService(services []string, name string) bool {
	for _, s := range services {
		if s == name {
			return true
		}
	}
	return false
}
