package output

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
)

func makeTransformReport() drift.Report {
	return drift.Report{
		Results: []drift.Result{
			{
				Service: "web",
				Drifted: true,
				Fields: []drift.Field{
					{Name: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"},
					{Name: "replicas", Expected: "3", Actual: "2"},
				},
			},
		},
	}
}

func TestTransformer_UppercaseService(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.UppercaseService = true

	var got drift.Report
	w := NewTransformer(opts, writerFunc(func(r drift.Report) error { got = r; return nil }))
	_ = w.Write(makeTransformReport())

	if got.Results[0].Service != "WEB" {
		t.Errorf("expected WEB, got %s", got.Results[0].Service)
	}
}

func TestTransformer_TrimImageTag(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.TrimImageTag = true

	var got drift.Report
	w := NewTransformer(opts, writerFunc(func(r drift.Report) error { got = r; return nil }))
	_ = w.Write(makeTransformReport())

	field := got.Results[0].Fields[0]
	if field.Expected != "nginx" {
		t.Errorf("expected 'nginx', got %s", field.Expected)
	}
	if field.Actual != "nginx" {
		t.Errorf("expected 'nginx', got %s", field.Actual)
	}
}

func TestTransformer_RenameFields(t *testing.T) {
	opts := DefaultTransformOptions()
	opts.RenameFields = map[string]string{"image": "container_image"}

	var got drift.Report
	w := NewTransformer(opts, writerFunc(func(r drift.Report) error { got = r; return nil }))
	_ = w.Write(makeTransformReport())

	if got.Results[0].Fields[0].Name != "container_image" {
		t.Errorf("expected container_image, got %s", got.Results[0].Fields[0].Name)
	}
}

func TestTransformer_NoOp(t *testing.T) {
	opts := DefaultTransformOptions()

	var got drift.Report
	w := NewTransformer(opts, writerFunc(func(r drift.Report) error { got = r; return nil }))
	_ = w.Write(makeTransformReport())

	if got.Results[0].Service != "web" {
		t.Errorf("expected web, got %s", got.Results[0].Service)
	}
	if got.Results[0].Fields[0].Expected != "nginx:1.25" {
		t.Errorf("expected nginx:1.25, got %s", got.Results[0].Fields[0].Expected)
	}
}
