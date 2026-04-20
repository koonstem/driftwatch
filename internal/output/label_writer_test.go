package output

import (
	"bytes"
	"io"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

// captureWriter records the report passed to it.
type captureWriter struct {
	report drift.Report
}

func (c *captureWriter) Write(_ io.Writer, r drift.Report) error {
	c.report = r
	return nil
}

func makeLabelReport() drift.Report {
	return drift.Report{
		Results: []drift.Result{
			{
				Service: "api",
				Drifted: true,
				Labels:  map[string]string{"env": "prod", "team": "backend"},
				Fields:  []drift.DriftField{{Field: "image", Expected: "api:1", Actual: "api:2"}},
			},
			{
				Service: "worker",
				Drifted: false,
				Labels:  map[string]string{"env": "staging"},
			},
		},
	}
}

func TestLabelWriter_Disabled_ForwardsAll(t *testing.T) {
	cap := &captureWriter{}
	lw := NewLabelWriter(LabelWriterOptions{Enabled: false}, cap)
	report := makeLabelReport()
	if err := lw.Write(&bytes.Buffer{}, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(cap.report.Results))
	}
}

func TestLabelWriter_FilterByKey_MatchesOnly(t *testing.T) {
	cap := &captureWriter{}
	opts := LabelWriterOptions{Enabled: true, FilterKey: "env", FilterValue: "prod", Annotate: false}
	lw := NewLabelWriter(opts, cap)
	if err := lw.Write(&bytes.Buffer{}, makeLabelReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.report.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(cap.report.Results))
	}
	if cap.report.Results[0].Service != "api" {
		t.Errorf("expected service 'api', got %q", cap.report.Results[0].Service)
	}
}

func TestLabelWriter_Annotate_AddsLabelField(t *testing.T) {
	cap := &captureWriter{}
	opts := LabelWriterOptions{Enabled: true, Annotate: true}
	lw := NewLabelWriter(opts, cap)
	if err := lw.Write(&bytes.Buffer{}, makeLabelReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := cap.report.Results[0]
	found := false
	for _, f := range result.Fields {
		if f.Field == "labels" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'labels' field in annotated result")
	}
}

func TestLabelWriter_FilterKeyOnly_NoValueFilter(t *testing.T) {
	cap := &captureWriter{}
	opts := LabelWriterOptions{Enabled: true, FilterKey: "env", Annotate: false}
	lw := NewLabelWriter(opts, cap)
	if err := lw.Write(&bytes.Buffer{}, makeLabelReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// both results have "env" label
	if len(cap.report.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(cap.report.Results))
	}
}

func TestLabelWriter_NoMatchingKey_FiltersAll(t *testing.T) {
	cap := &captureWriter{}
	opts := LabelWriterOptions{Enabled: true, FilterKey: "nonexistent", Annotate: false}
	lw := NewLabelWriter(opts, cap)
	if err := lw.Write(&bytes.Buffer{}, makeLabelReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.report.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(cap.report.Results))
	}
}
