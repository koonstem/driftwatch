package output_test

import (
	"errors"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func makeFilterReport() drift.Report {
	return drift.Report{
		Results: []drift.Result{
			{Service: "alpha", Drifted: true},
			{Service: "beta", Drifted: false},
			{Service: "gamma", Drifted: true},
		},
	}
}

type captureWriter struct {
	report drift.Report
}

func (c *captureWriter) Write(r drift.Report) error {
	c.report = r
	return nil
}

func TestFilterWriter_NoOptions_ForwardsAll(t *testing.T) {
	cap := &captureWriter{}
	fw := output.NewFilterWriter(output.DefaultFilterWriterOptions(), cap)
	_ = fw.Write(makeFilterReport())
	if len(cap.report.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(cap.report.Results))
	}
}

func TestFilterWriter_OnlyDrifted(t *testing.T) {
	cap := &captureWriter{}
	opts := output.DefaultFilterWriterOptions()
	opts.OnlyDrifted = true
	fw := output.NewFilterWriter(opts, cap)
	_ = fw.Write(makeFilterReport())
	if len(cap.report.Results) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(cap.report.Results))
	}
}

func TestFilterWriter_ByService(t *testing.T) {
	cap := &captureWriter{}
	opts := output.DefaultFilterWriterOptions()
	opts.Services = []string{"beta"}
	fw := output.NewFilterWriter(opts, cap)
	_ = fw.Write(makeFilterReport())
	if len(cap.report.Results) != 1 || cap.report.Results[0].Service != "beta" {
		t.Fatalf("expected only beta, got %+v", cap.report.Results)
	}
}

func TestFilterWriter_PropagatesInnerError(t *testing.T) {
	errInner := errors.New("inner failure")
	bad := &errWriter{err: errInner}
	fw := output.NewFilterWriter(output.DefaultFilterWriterOptions(), bad)
	if err := fw.Write(makeFilterReport()); !errors.Is(err, errInner) {
		t.Fatalf("expected inner error, got %v", err)
	}
}

type errWriter struct{ err error }

func (e *errWriter) Write(_ drift.Report) error { return e.err }
