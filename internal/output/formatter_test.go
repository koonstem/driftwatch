package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makeReport(drifted bool) *drift.Report {
	if !drifted {
		return &drift.Report{Results: []drift.Result{}}
	}
	return &drift.Report{
		Results: []drift.Result{
			{
				ServiceName: "web",
				Drifted:     true,
				Diffs: []drift.Diff{
					{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
				},
			},
		},
	}
}

func TestFormatter_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatText, &buf)
	if err := f.Write(makeReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestFormatter_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatText, &buf)
	_ = f.Write(makeReport(true))
	if !strings.Contains(buf.String(), "web") {
		t.Errorf("expected service name in output: %s", buf.String())
	}
}

func TestFormatter_JSONNoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatJSON, &buf)
	_ = f.Write(makeReport(false))
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestFormatter_JSONWithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatJSON, &buf)
	_ = f.Write(makeReport(true))
	if !strings.Contains(buf.String(), "web") {
		t.Errorf("expected service name in JSON: %s", buf.String())
	}
}

func TestFormatter_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatTable, &buf)
	if err := f.Write(makeReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SERVICE") {
		t.Errorf("expected table header in output: %s", buf.String())
	}
}

func TestFormatter_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter("xml", &buf)
	if err := f.Write(makeReport(false)); err == nil {
		t.Error("expected error for unsupported format")
	}
}
