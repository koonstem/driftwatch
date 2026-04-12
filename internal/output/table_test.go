package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makeTableReport(drifted bool) *drift.Report {
	if !drifted {
		return &drift.Report{Results: []drift.Result{}}
	}
	return &drift.Report{
		Results: []drift.Result{
			{
				ServiceName: "api",
				Drifted:     true,
				Diffs: []drift.Diff{
					{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
					{Field: "replicas", Expected: "3", Actual: "2"},
				},
			},
		},
	}
}

func TestWriteTable_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	err := writeTable(&buf, makeTableReport(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteTable_WithDrift_HasHeaders(t *testing.T) {
	var buf bytes.Buffer
	err := writeTable(&buf, makeTableReport(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, hdr := range []string{"SERVICE", "FIELD", "EXPECTED", "ACTUAL"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output:\n%s", hdr, out)
		}
	}
}

func TestWriteTable_WithDrift_HasRows(t *testing.T) {
	var buf bytes.Buffer
	_ = writeTable(&buf, makeTableReport(true))
	out := buf.String()
	for _, want := range []string{"api", "image", "nginx:1.25", "nginx:1.24", "replicas", "3", "2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in table output:\n%s", want, out)
		}
	}
}

func TestWriteTable_WithDrift_HasSeparator(t *testing.T) {
	var buf bytes.Buffer
	_ = writeTable(&buf, makeTableReport(true))
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Fatal("expected at least header + separator")
	}
	if !strings.Contains(lines[1], "---") {
		t.Errorf("expected separator line, got: %s", lines[1])
	}
}
