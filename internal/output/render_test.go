package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeRenderReport(drifted bool) *drift.Report {
	results := []drift.Result{
		{
			Service:  "api",
			Drifted:  drifted,
			Expected: "nginx:1.25",
			Actual:   "nginx:1.24",
			Fields:   []drift.FieldDiff{{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}},
		},
	}
	if !drifted {
		results[0].Actual = "nginx:1.25"
		results[0].Fields = nil
	}
	return &drift.Report{Results: results}
}

func TestRenderer_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "text", Writer: &buf})
	if err := r.Render(makeRenderReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderer_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "json", Writer: &buf})
	if err := r.Render(makeRenderReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"service\"") {
		t.Error("expected JSON output to contain 'service' key")
	}
}

func TestRenderer_TableFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "table", Writer: &buf})
	if err := r.Render(makeRenderReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty table output")
	}
}

func TestRenderer_SummaryFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "summary", Writer: &buf})
	if err := r.Render(makeRenderReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty summary output")
	}
}

func TestRenderer_DiffFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "diff", Writer: &buf})
	if err := r.Render(makeRenderReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty diff output")
	}
}

func TestRenderer_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "xml", Writer: &buf})
	err := r.Render(makeRenderReport(false))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("expected error to mention format name, got: %v", err)
	}
}

func TestRenderer_EmptyFormat_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(RenderOptions{Format: "", Writer: &buf})
	if err := r.Render(makeRenderReport(false)); err != nil {
		t.Fatalf("unexpected error for empty format: %v", err)
	}
}
