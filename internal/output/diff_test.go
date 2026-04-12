package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeDiffReport(drifted bool) drift.Report {
	field := drift.DriftedField{
		Field:    "image",
		Expected: "nginx:1.25",
		Actual:   "nginx:1.24",
	}
	result := drift.Result{
		ServiceName:   "web",
		ContainerName: "/web_1",
		Drifted:       drifted,
	}
	if drifted {
		result.DriftedFields = []drift.DriftedField{field}
	}
	return drift.Report{Results: []drift.Result{result}}
}

func TestDiffWriter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := NewDiffWriter(&buf, false)
	w.Write(makeDiffReport(false))

	out := buf.String()
	if !strings.Contains(out, "no drift detected") {
		t.Errorf("expected no-drift message, got: %q", out)
	}
}

func TestDiffWriter_WithDrift_HasHeader(t *testing.T) {
	var buf bytes.Buffer
	w := NewDiffWriter(&buf, false)
	w.Write(makeDiffReport(true))

	out := buf.String()
	if !strings.Contains(out, "--- declared: web") {
		t.Errorf("expected declared header, got: %q", out)
	}
	if !strings.Contains(out, "+++ running:  /web_1") {
		t.Errorf("expected running header, got: %q", out)
	}
}

func TestDiffWriter_WithDrift_HasFieldLines(t *testing.T) {
	var buf bytes.Buffer
	w := NewDiffWriter(&buf, false)
	w.Write(makeDiffReport(true))

	out := buf.String()
	if !strings.Contains(out, "- image: nginx:1.25") {
		t.Errorf("expected removed line, got: %q", out)
	}
	if !strings.Contains(out, "+ image: nginx:1.24") {
		t.Errorf("expected added line, got: %q", out)
	}
}

func TestDiffWriter_ColorEnabled_ContainsEscape(t *testing.T) {
	var buf bytes.Buffer
	w := NewDiffWriter(&buf, true)
	w.Write(makeDiffReport(true))

	out := buf.String()
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected ANSI escape codes with color enabled, got: %q", out)
	}
}
