package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeNotifyResults(drifted bool) []drift.DriftResult {
	fields := []drift.DriftedField{}
	if drifted {
		fields = append(fields, drift.DriftedField{Field: "image", Expected: "nginx:1.24", Actual: "nginx:1.25"})
	}
	return []drift.DriftResult{
		{ServiceName: "web", Drifted: drifted, Fields: fields},
	}
}

func TestNotifyWriter_Disabled_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultNotifyOptions()
	opts.Enabled = false
	w := NewNotifyWriter(opts, &buf)
	_ = w.Notify(makeNotifyResults(true))
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}
}

func TestNotifyWriter_Enabled_NoDrift_OnlyDrift(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultNotifyOptions()
	opts.Enabled = true
	opts.OnlyDrift = true
	w := NewNotifyWriter(opts, &buf)
	_ = w.Notify(makeNotifyResults(false))
	if buf.Len() != 0 {
		t.Errorf("expected no output when no drift and OnlyDrift=true, got %q", buf.String())
	}
}

func TestNotifyWriter_Enabled_NoDrift_AllEvents(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultNotifyOptions()
	opts.Enabled = true
	opts.OnlyDrift = false
	w := NewNotifyWriter(opts, &buf)
	_ = w.Notify(makeNotifyResults(false))
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got %q", buf.String())
	}
}

func TestNotifyWriter_Enabled_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultNotifyOptions()
	opts.Enabled = true
	opts.OnlyDrift = true
	w := NewNotifyWriter(opts, &buf)
	_ = w.Notify(makeNotifyResults(true))
	out := buf.String()
	if !strings.Contains(out, "web") {
		t.Errorf("expected service name in output, got %q", out)
	}
	if !strings.Contains(out, "Drift detected") {
		t.Errorf("expected drift message, got %q", out)
	}
}

func TestNotifyWriter_FansOutToMultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	opts := DefaultNotifyOptions()
	opts.Enabled = true
	opts.OnlyDrift = false
	w := NewNotifyWriter(opts, &buf1, &buf2)
	_ = w.Notify(makeNotifyResults(false))
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected output in both writers")
	}
}
