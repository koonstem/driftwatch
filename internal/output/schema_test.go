package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
	"github.com/spf13/cobra"
)

func makeSchemaResults() []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Drifted: true,
			Diffs: []drift.Diff{
				{Field: "image", Expected: "api:1", Actual: "api:2"},
				{Field: "replicas", Expected: "2", Actual: "3"},
			},
		},
		{
			Service: "worker",
			Drifted: true,
			Diffs: []drift.Diff{
				{Field: "image", Expected: "worker:1", Actual: "worker:2"},
			},
		},
		{Service: "cache", Drifted: false},
	}
}

func TestSchemaWriter_ContainsVersion(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultSchemaOptions()
	w := NewSchemaWriter(&buf, opts)
	if err := w(makeSchemaResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"version"`) {
		t.Error("expected version field in schema output")
	}
}

func TestSchemaWriter_ServicesAlphabetical(t *testing.T) {
	var buf bytes.Buffer
	w := NewSchemaWriter(&buf, DefaultSchemaOptions())
	_ = w(makeSchemaResults())

	var doc schemaDoc
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(doc.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(doc.Services))
	}
	if doc.Services[0] != "api" || doc.Services[1] != "cache" || doc.Services[2] != "worker" {
		t.Errorf("services not sorted: %v", doc.Services)
	}
}

func TestSchemaWriter_FieldCounts(t *testing.T) {
	var buf bytes.Buffer
	w := NewSchemaWriter(&buf, DefaultSchemaOptions())
	_ = w(makeSchemaResults())

	var doc schemaDoc
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &doc)

	counts := map[string]int{}
	for _, f := range doc.Fields {
		counts[f.Name] = f.Observed
	}
	if counts["image"] != 2 {
		t.Errorf("expected image count 2, got %d", counts["image"])
	}
	if counts["replicas"] != 1 {
		t.Errorf("expected replicas count 1, got %d", counts["replicas"])
	}
}

func TestSchemaWriter_NoDrift_EmptyFields(t *testing.T) {
	var buf bytes.Buffer
	w := NewSchemaWriter(&buf, DefaultSchemaOptions())
	_ = w([]drift.Result{{Service: "svc", Drifted: false}})

	var doc schemaDoc
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &doc)
	if len(doc.Fields) != 0 {
		t.Errorf("expected no fields, got %v", doc.Fields)
	}
}

func newSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindSchemaFlags(cmd)
	return cmd
}

func TestBindSchemaFlags_Defaults(t *testing.T) {
	cmd := newSchemaCmd()
	_ = cmd.Execute()
	opts, err := SchemaOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Pretty {
		t.Error("expected pretty=true by default")
	}
	if opts.Version != "v1" {
		t.Errorf("expected version v1, got %s", opts.Version)
	}
}

func TestBindSchemaFlags_EmptyVersionError(t *testing.T) {
	cmd := newSchemaCmd()
	_ = cmd.Flags().Set("schema-version", "")
	_, err := SchemaOptionsFromFlags(cmd)
	if err == nil {
		t.Error("expected error for empty version")
	}
}
