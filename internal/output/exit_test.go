package output

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/drift"
)

func makeReportWithDrift(hasDrift bool) *drift.Report {
	results := []drift.Result{}
	if hasDrift {
		results = append(results, drift.Result{
			ServiceName: "api",
			Drifted:     true,
			Reasons:     []string{"image mismatch"},
		})
	}
	return &drift.Report{Results: results}
}

func TestExitCode_NoDrift_FailEnabled(t *testing.T) {
	ec := NewExitCoder(true)
	report := makeReportWithDrift(false)
	if got := ec.Code(report); got != 0 {
		t.Errorf("expected exit code 0, got %d", got)
	}
}

func TestExitCode_WithDrift_FailEnabled(t *testing.T) {
	ec := NewExitCoder(true)
	report := makeReportWithDrift(true)
	if got := ec.Code(report); got != 1 {
		t.Errorf("expected exit code 1, got %d", got)
	}
}

func TestExitCode_WithDrift_FailDisabled(t *testing.T) {
	ec := NewExitCoder(false)
	report := makeReportWithDrift(true)
	if got := ec.Code(report); got != 0 {
		t.Errorf("expected exit code 0 when failOnDrift=false, got %d", got)
	}
}

func TestExitCode_NoDrift_FailDisabled(t *testing.T) {
	ec := NewExitCoder(false)
	report := makeReportWithDrift(false)
	if got := ec.Code(report); got != 0 {
		t.Errorf("expected exit code 0, got %d", got)
	}
}
