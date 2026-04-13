package output

import (
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makeSortResults() []drift.Result {
	return []drift.Result{
		{Service: "zebra", Drifted: false, Fields: nil},
		{Service: "alpha", Drifted: true, Fields: []drift.FieldDiff{{Field: "image"}}},
		{Service: "mango", Drifted: true, Fields: []drift.FieldDiff{{Field: "env"}}},
	}
}

func TestSortResults_ByService(t *testing.T) {
	results := SortResults(makeSortResults(), SortOptions{By: SortByService})
	if results[0].Service != "alpha" || results[1].Service != "mango" || results[2].Service != "zebra" {
		t.Errorf("unexpected order: %v", serviceNames(results))
	}
}

func TestSortResults_ByService_Reverse(t *testing.T) {
	results := SortResults(makeSortResults(), SortOptions{By: SortByService, Reverse: true})
	if results[0].Service != "zebra" {
		t.Errorf("expected zebra first, got %s", results[0].Service)
	}
}

func TestSortResults_ByStatus_DriftedFirst(t *testing.T) {
	results := SortResults(makeSortResults(), SortOptions{By: SortByStatus})
	if !results[0].Drifted {
		t.Errorf("expected drifted result first")
	}
	if results[2].Drifted {
		t.Errorf("expected non-drifted result last")
	}
}

func TestSortResults_ByField(t *testing.T) {
	results := SortResults(makeSortResults(), SortOptions{By: SortByField})
	// zebra has no fields → firstField returns "", sorts before "env" and "image"
	if results[0].Service != "zebra" {
		t.Errorf("expected zebra first (empty field), got %s", results[0].Service)
	}
}

func TestSortResults_DoesNotMutateInput(t *testing.T) {
	original := makeSortResults()
	originalFirst := original[0].Service
	SortResults(original, SortOptions{By: SortByService})
	if original[0].Service != originalFirst {
		t.Errorf("input slice was mutated")
	}
}

func serviceNames(results []drift.Result) []string {
	names := make([]string, len(results))
	for i, r := range results {
		names[i] = r.Service
	}
	return names
}
