package output

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeAggResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			Service: "api",
			Fields: []drift.FieldResult{
				{Name: "image", Drifted: true},
				{Name: "env", Drifted: false},
			},
		},
		{
			Service: "worker",
			Fields: []drift.FieldResult{
				{Name: "image", Drifted: true},
				{Name: "image", Drifted: true},
			},
		},
	}
}

func TestAggregator_ByService(t *testing.T) {
	agg := NewAggregator(AggregateOptions{GroupBy: "service"})
	results, err := agg(makeAggResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Key != "api" {
		t.Errorf("expected first key api, got %s", results[0].Key)
	}
	if results[0].Drifted != 1 {
		t.Errorf("expected 1 drifted for api, got %d", results[0].Drifted)
	}
}

func TestAggregator_ByField(t *testing.T) {
	agg := NewAggregator(AggregateOptions{GroupBy: "field"})
	results, err := agg(makeAggResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := map[string]int{}
	for _, r := range results {
		keys[r.Key] = r.Drifted
	}
	if keys["image"] != 3 {
		t.Errorf("expected 3 drifted image fields, got %d", keys["image"])
	}
	if keys["env"] != 0 {
		t.Errorf("expected 0 drifted env fields, got %d", keys["env"])
	}
}

func TestAggregator_InvalidGroupBy(t *testing.T) {
	agg := NewAggregator(AggregateOptions{GroupBy: "unknown"})
	_, err := agg(makeAggResults())
	if err == nil {
		t.Fatal("expected error for unknown group-by")
	}
}

func TestDefaultAggregateOptions(t *testing.T) {
	opts := DefaultAggregateOptions()
	if opts.GroupBy != "service" {
		t.Errorf("expected default group-by service, got %s", opts.GroupBy)
	}
}
