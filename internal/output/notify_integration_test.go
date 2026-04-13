package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func TestNotify_FullPipeline_DriftedServices(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "api", Drifted: true, Fields: []drift.DriftedField{
			{Field: "image", Expected: "api:1.0", Actual: "api:2.0"},
		}},
		{ServiceName: "db", Drifted: false},
	}

	var buf bytes.Buffer
	opts := NotifyOptions{Enabled: true, OnlyDrift: true, Channels: []string{"stdout"}}
	w := NewNotifyWriter(opts, &buf)

	if err := w.Notify(results); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Errorf("expected drifted service 'api' in output, got: %q", out)
	}
	if strings.Contains(out, "db") {
		t.Errorf("unexpected non-drifted service 'db' in output, got: %q", out)
	}
}

func TestNotify_FullPipeline_NoDriftAllEvents(t *testing.T) {
	results := []drift.DriftResult{
		{ServiceName: "worker", Drifted: false},
	}

	var buf bytes.Buffer
	opts := NotifyOptions{Enabled: true, OnlyDrift: false, Channels: []string{"stdout"}}
	w := NewNotifyWriter(opts, &buf)

	if err := w.Notify(results); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected no-drift message, got: %q", out)
	}
}
