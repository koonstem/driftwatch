package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func makeReport(drifted bool) *drift.Report {
	r := &drift.Report{}
	if drifted {
		r.Entries = []drift.Entry{
			{ServiceName: "api", Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24", Drifted: true},
		}
	}
	return r
}

func TestFormatter_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	if err := f.Write(makeReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestFormatter_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	if err := f.Write(makeReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"api", "image", "nginx:1.25", "nginx:1.24"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestFormatter_JSONNoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatJSON, &buf)
	if err := f.Write(makeReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"drifted": false`) {
		t.Errorf("expected drifted=false in JSON, got: %q", buf.String())
	}
}

func TestFormatter_JSONWithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatJSON, &buf)
	if err := f.Write(makeReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{`"drifted": true`, `"service": "api"`, `"field": "image"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %q", want, out)
		}
	}
}

func TestNewFormatter_NilWriterUsesStdout(t *testing.T) {
	// Should not panic when w is nil
	f := output.NewFormatter(output.FormatText, nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}
