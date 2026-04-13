package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

// TestMetricWriter_TimestampPresent verifies that a Unix-ms timestamp suffix
// is appended to each metric line when Timestamp is true.
func TestMetricWriter_TimestampPresent(t *testing.T) {
	var buf bytes.Buffer
	opts := MetricOptions{Enabled: true, Timestamp: true, Prefix: "driftwatch"}
	w := NewMetricWriter(&buf, opts)

	report := drift.Report{
		Results: []drift.Result{
			{Service: "api", Drifted: true},
			{Service: "db", Drifted: false},
		},
	}
	if err := w.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		// metric lines should have 3 fields: name value timestamp
		parts := strings.Fields(line)
		if len(parts) != 3 {
			t.Errorf("expected 3 fields in metric line %q, got %d", line, len(parts))
		}
	}
}

// TestMetricWriter_FullPipeline exercises the flags → options → writer path.
func TestMetricWriter_FullPipeline(t *testing.T) {
	cmd := newMetricCmd()
	_ = cmd.ParseFlags([]string{"--metrics", "--metrics-prefix", "ci", "--metrics-timestamp=false"})

	opts, err := MetricOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("flag parse error: %v", err)
	}

	var buf bytes.Buffer
	w := NewMetricWriter(&buf, opts)

	report := drift.Report{
		Results: []drift.Result{
			{Service: "svc-a", Drifted: true},
		},
	}
	if err := w.Write(report); err != nil {
		t.Fatalf("write error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ci_services_total 1") {
		t.Errorf("expected ci_services_total 1 in output:\n%s", out)
	}
	if !strings.Contains(out, "ci_drifted_total 1") {
		t.Errorf("expected ci_drifted_total 1 in output:\n%s", out)
	}
}
