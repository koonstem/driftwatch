package output

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// MetricOptions controls metric output behaviour.
type MetricOptions struct {
	Enabled   bool
	Timestamp bool
	Prefix    string
}

// DefaultMetricOptions returns sensible defaults.
func DefaultMetricOptions() MetricOptions {
	return MetricOptions{
		Enabled:   false,
		Timestamp: true,
		Prefix:    "driftwatch",
	}
}

// MetricWriter emits Prometheus-style text metrics for drift results.
type MetricWriter struct {
	opts MetricOptions
	w    io.Writer
}

// NewMetricWriter creates a MetricWriter that writes to w.
func NewMetricWriter(w io.Writer, opts MetricOptions) *MetricWriter {
	return &MetricWriter{opts: opts, w: w}
}

// Write emits gauge metrics derived from the drift report.
func (m *MetricWriter) Write(report drift.Report) error {
	if !m.opts.Enabled {
		return nil
	}

	total := len(report.Results)
	drifted := 0
	for _, r := range report.Results {
		if r.Drifted {
			drifted++
		}
	}

	prefix := m.opts.Prefix
	ts := ""
	if m.opts.Timestamp {
		ts = fmt.Sprintf(" %d", time.Now().UnixMilli())
	}

	lines := []string{
		fmt.Sprintf("# HELP %s_services_total Total number of services checked.", prefix),
		fmt.Sprintf("# TYPE %s_services_total gauge", prefix),
		fmt.Sprintf("%s_services_total %d%s", prefix, total, ts),
		fmt.Sprintf("# HELP %s_drifted_total Number of services with drift detected.", prefix),
		fmt.Sprintf("# TYPE %s_drifted_total gauge", prefix),
		fmt.Sprintf("%s_drifted_total %d%s", prefix, drifted, ts),
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(m.w, line); err != nil {
			return err
		}
	}
	return nil
}
