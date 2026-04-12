package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func makeSummaryReport(drifted bool) drift.Report {
	results := []drift.Result{
		{
			Service:  "api",
			Drifted:  drifted,
			Reasons:  []string{},
		},
		{
			Service: "worker",
			Drifted: false,
			Reasons: []string{},
		},
	}
	return drift.Report{Results: results}
}

func TestSummaryWriter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := NewSummaryWriter(&buf)
	report := makeSummaryReport(false)

	if err := w.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Services checked : 2") {
		t.Errorf("expected services count 2, got:\n%s", out)
	}
	if !strings.Contains(out, "Drifted          : 0") {
		t.Errorf("expected 0 drifted, got:\n%s", out)
	}
	if !strings.Contains(out, "Status           : OK") {
		t.Errorf("expected OK status, got:\n%s", out)
	}
}

func TestSummaryWriter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	w := NewSummaryWriter(&buf)
	report := makeSummaryReport(true)

	if err := w.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Drifted          : 1") {
		t.Errorf("expected 1 drifted, got:\n%s", out)
	}
	if !strings.Contains(out, "Status           : DRIFT DETECTED") {
		t.Errorf("expected DRIFT DETECTED status, got:\n%s", out)
	}
	if !strings.Contains(out, "Clean            : 1") {
		t.Errorf("expected 1 clean, got:\n%s", out)
	}
}

func TestSummaryWriter_HeaderAndFooter(t *testing.T) {
	var buf bytes.Buffer
	w := NewSummaryWriter(&buf)
	report := makeSummaryReport(false)

	_ = w.Write(report)
	out := buf.String()

	if !strings.Contains(out, "Drift Detection Summary") {
		t.Errorf("expected header line, got:\n%s", out)
	}
	if !strings.Contains(out, "----") {
		t.Errorf("expected footer separator, got:\n%s", out)
	}
}
