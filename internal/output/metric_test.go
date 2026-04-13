package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeMetricReport(driftedCount, cleanCount int) drift.Report {
	var results []drift.Result
	for i := 0; i < driftedCount; i++ {
		results = append(results, drift.Result{Service: "svc", Drifted: true})
	}
	for i := 0; i < cleanCount; i++ {
		results = append(results, drift.Result{Service: "ok", Drifted: false})
	}
	return drift.Report{Results: results}
}

func TestMetricWriter_Disabled(t *testing.T) {
	var buf bytes.Buffer
	w := NewMetricWriter(&buf, MetricOptions{Enabled: false})
	if err := w.Write(makeMetricReport(2, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}
}

func TestMetricWriter_Enabled_ContainsGauges(t *testing.T) {
	var buf bytes.Buffer
	opts := MetricOptions{Enabled: true, Timestamp: false, Prefix: "driftwatch"}
	w := NewMetricWriter(&buf, opts)
	if err := w.Write(makeMetricReport(2, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{
		"driftwatch_services_total 3",
		"driftwatch_drifted_total 2",
		"# TYPE driftwatch_services_total gauge",
		"# TYPE driftwatch_drifted_total gauge",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestMetricWriter_CustomPrefix(t *testing.T) {
	var buf bytes.Buffer
	opts := MetricOptions{Enabled: true, Timestamp: false, Prefix: "myapp"}
	w := NewMetricWriter(&buf, opts)
	if err := w.Write(makeMetricReport(0, 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "myapp_services_total 3") {
		t.Errorf("expected custom prefix in output, got:\n%s", out)
	}
	if strings.Contains(out, "driftwatch_") {
		t.Errorf("default prefix should not appear with custom prefix set")
	}
}

func TestMetricWriter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	opts := MetricOptions{Enabled: true, Timestamp: false, Prefix: "driftwatch"}
	w := NewMetricWriter(&buf, opts)
	if err := w.Write(makeMetricReport(0, 5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "driftwatch_drifted_total 0") {
		t.Errorf("expected zero drifted metric, got:\n%s", out)
	}
}
