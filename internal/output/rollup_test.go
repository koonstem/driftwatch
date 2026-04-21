package output

import (
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

func makeRollupResults() []drift.Result {
	return []drift.Result{
		{
			Service: "api",
			Drifted: true,
			Fields:  []drift.Field{{Name: "image"}, {Name: "env"}},
		},
		{
			Service: "worker",
			Drifted: true,
			Fields:  []drift.Field{{Name: "image"}},
		},
		{
			Service: "db",
			Drifted: false,
		},
	}
}

func TestRollupWriter_Disabled_ForwardsUnchanged(t *testing.T) {
	opts := DefaultRollupOptions()
	opts.Enabled = false
	var got drift.Report
	w, err := NewRollupWriter(opts, writerFunc(func(r drift.Report) error {
		got = r
		return nil
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := drift.Report{Results: makeRollupResults()}
	if err := w.Write(report); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if len(got.Annotations) != 0 {
		t.Errorf("expected no annotations when disabled, got %v", got.Annotations)
	}
}

func TestRollupWriter_Enabled_ByService_AnnotatesTopEntries(t *testing.T) {
	opts := DefaultRollupOptions()
	opts.Enabled = true
	opts.GroupBy = "service"
	var got drift.Report
	w, _ := NewRollupWriter(opts, writerFunc(func(r drift.Report) error {
		got = r
		return nil
	}))
	report := drift.Report{Results: makeRollupResults()}
	if err := w.Write(report); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if got.Annotations["rollup.0.key"] != "api" {
		t.Errorf("expected top rollup key 'api', got %q", got.Annotations["rollup.0.key"])
	}
	if got.Annotations["rollup.0.drift_count"] != "2" {
		t.Errorf("expected drift_count '2', got %q", got.Annotations["rollup.0.drift_count"])
	}
}

func TestRollupWriter_Enabled_ByField_GroupsCorrectly(t *testing.T) {
	opts := DefaultRollupOptions()
	opts.Enabled = true
	opts.GroupBy = "field"
	var got drift.Report
	w, _ := NewRollupWriter(opts, writerFunc(func(r drift.Report) error {
		got = r
		return nil
	}))
	report := drift.Report{Results: makeRollupResults()}
	if err := w.Write(report); err != nil {
		t.Fatalf("write error: %v", err)
	}
	// "image" appears in both api and worker => count 2, should be top
	if got.Annotations["rollup.0.key"] != "image" {
		t.Errorf("expected top field 'image', got %q", got.Annotations["rollup.0.key"])
	}
}

func TestRollupWriter_NilNext_ReturnsError(t *testing.T) {
	opts := DefaultRollupOptions()
	_, err := NewRollupWriter(opts, nil)
	if err == nil {
		t.Fatal("expected error for nil next writer")
	}
}

func TestRollupWriter_InvalidGroupBy_ReturnsError(t *testing.T) {
	opts := DefaultRollupOptions()
	opts.GroupBy = "invalid"
	_, err := NewRollupWriter(opts, writerFunc(func(r drift.Report) error { return nil }))
	if err == nil {
		t.Fatal("expected error for invalid group_by")
	}
}
