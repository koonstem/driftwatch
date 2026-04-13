package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeTemplateResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Drifted: false, Fields: nil},
		{
			Service: "worker",
			Drifted: true,
			Fields: []drift.FieldDiff{
				{Field: "image", Expected: "worker:1", Actual: "worker:2"},
			},
		},
	}
}

func TestTemplateWriter_InlineNoDrift(t *testing.T) {
	var buf bytes.Buffer
	tmplStr := "total={{.Total}} drifted={{.Drifted}}"
	w, err := NewTemplateWriter(&buf, TemplateOptions{TemplateStr: tmplStr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write([]drift.Result{{Service: "api", Drifted: false}}); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if !strings.Contains(buf.String(), "total=1") {
		t.Errorf("expected total=1, got %q", buf.String())
	}
	if !strings.Contains(buf.String(), "drifted=0") {
		t.Errorf("expected drifted=0, got %q", buf.String())
	}
}

func TestTemplateWriter_InlineWithDrift(t *testing.T) {
	var buf bytes.Buffer
	tmplStr := `{{range .Results}}{{if .Drifted}}DRIFT:{{.Service}}{{end}}{{end}}`
	w, err := NewTemplateWriter(&buf, TemplateOptions{TemplateStr: tmplStr})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeTemplateResults()); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if !strings.Contains(buf.String(), "DRIFT:worker") {
		t.Errorf("expected DRIFT:worker in output, got %q", buf.String())
	}
}

func TestTemplateWriter_FileTemplate(t *testing.T) {
	tmplContent := `services={{ .Total }}`
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.txt")
	if err := os.WriteFile(path, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	w, err := NewTemplateWriter(&buf, TemplateOptions{TemplatePath: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeTemplateResults()); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if !strings.Contains(buf.String(), "services=2") {
		t.Errorf("expected services=2, got %q", buf.String())
	}
}

func TestTemplateWriter_NoSource_ReturnsError(t *testing.T) {
	_, err := NewTemplateWriter(&bytes.Buffer{}, TemplateOptions{})
	if err == nil {
		t.Fatal("expected error for missing template source")
	}
}

func TestTemplateWriter_InvalidInline_ReturnsError(t *testing.T) {
	_, err := NewTemplateWriter(&bytes.Buffer{}, TemplateOptions{TemplateStr: "{{.Unclosed"})
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}
