package output_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
)

func TestSortResults_FullPipeline_ServiceAsc(t *testing.T) {
	input := []drift.Result{
		{Service: "svc-c", Drifted: true},
		{Service: "svc-a", Drifted: false},
		{Service: "svc-b", Drifted: true},
	}

	opts := output.DefaultSortOptions()
	sorted := output.SortResults(input, opts)

	expected := []string{"svc-a", "svc-b", "svc-c"}
	for i, r := range sorted {
		if r.Service != expected[i] {
			t.Errorf("position %d: expected %q got %q", i, expected[i], r.Service)
		}
	}
}

func TestSortResults_FullPipeline_StatusDesc(t *testing.T) {
	input := []drift.Result{
		{Service: "ok-svc", Drifted: false},
		{Service: "bad-svc", Drifted: true},
	}

	opts := output.SortOptions{By: output.SortByStatus, Reverse: true}
	sorted := output.SortResults(input, opts)

	// reverse of drifted-first → non-drifted first
	if sorted[0].Drifted {
		t.Errorf("expected non-drifted first in reverse status sort")
	}
}
