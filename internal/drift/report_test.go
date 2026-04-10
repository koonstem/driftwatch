package drift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/runner"
	"github.com/yourorg/driftwatch/internal/source"
)

func TestReporter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)

	results := []drift.Result{
		{ServiceName: "api", Drifted: false},
	}
	r.Print(results)

	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got:\n%s", out)
	}
	if strings.Contains(out, "[DRIFT]") {
		t.Errorf("unexpected [DRIFT] in output")
	}
}

func TestReporter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)

	results := []drift.Result{
		{
			ServiceName: "worker",
			Expected:    source.ServiceSpec{Name: "worker", Image: "app:2.0"},
			Actual:      runner.ContainerInfo{Name: "worker", Image: "app:1.9"},
			Drifted:     true,
			Reasons:     []string{`image mismatch: expected "app:2.0", got "app:1.9"`},
		},
	}
	r.Print(results)

	out := buf.String()
	if !strings.Contains(out, "[DRIFT]") {
		t.Errorf("expected [DRIFT] in output, got:\n%s", out)
	}
	if !strings.Contains(out, "image mismatch") {
		t.Errorf("expected reason in output, got:\n%s", out)
	}
}

func TestReporter_Summary(t *testing.T) {
	var buf bytes.Buffer
	r := drift.NewReporter(&buf)

	results := []drift.Result{
		{ServiceName: "api", Drifted: false},
		{ServiceName: "db", Drifted: true, Reasons: []string{"status mismatch"}},
	}
	r.Print(results)

	out := buf.String()
	if !strings.Contains(out, "2 service(s) checked, 1 drifted") {
		t.Errorf("unexpected summary line:\n%s", out)
	}
}

func TestHasDrift(t *testing.T) {
	if drift.HasDrift([]drift.Result{{Drifted: false}}) {
		t.Error("expected no drift")
	}
	if !drift.HasDrift([]drift.Result{{Drifted: true}}) {
		t.Error("expected drift")
	}
}
