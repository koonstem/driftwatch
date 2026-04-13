package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeExportResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			Service: "api",
			Status:  drift.StatusDrifted,
			Fields: []drift.FieldDiff{
				{Field: "image", Expected: "api:v1", Actual: "api:v2"},
			},
		},
		{
			Service: "worker",
			Status:  drift.StatusOK,
			Fields:  nil,
		},
	}
}

func TestExporter_CSV_Headers(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter(ExportOptions{Format: "csv", Timestamp: false})
	if err := export(&buf, makeExportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv parse error: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("expected at least one record")
	}
	header := records[0]
	expected := []string{"service", "status", "field", "expected", "actual"}
	for i, col := range expected {
		if header[i] != col {
			t.Errorf("header[%d] = %q, want %q", i, header[i], col)
		}
	}
}

func TestExporter_CSV_Rows(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter(ExportOptions{Format: "csv", Timestamp: false})
	if err := export(&buf, makeExportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "api") {
		t.Error("expected service 'api' in CSV output")
	}
	if !strings.Contains(buf.String(), "image") {
		t.Error("expected field 'image' in CSV output")
	}
}

func TestExporter_CSV_Timestamp(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter(ExportOptions{Format: "csv", Timestamp: true})
	if err := export(&buf, makeExportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "generated_at") {
		t.Error("expected generated_at column in CSV output")
	}
}

func TestExporter_JSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter(ExportOptions{Format: "json", Timestamp: false})
	if err := export(&buf, makeExportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected at least one JSON row")
	}
	if _, ok := rows[0]["service"]; !ok {
		t.Error("expected 'service' key in JSON row")
	}
}

func TestExporter_JSON_NoTimestamp(t *testing.T) {
	var buf bytes.Buffer
	export := NewExporter(ExportOptions{Format: "json", Timestamp: false})
	if err := export(&buf, makeExportResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "generated_at") {
		t.Error("did not expect generated_at in JSON output when Timestamp=false")
	}
}
