package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/filter"
)

func makeResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ServiceName: "api",
			Drifted:     true,
			Differences: []drift.Difference{{Field: "image", Expected: "api:1.0", Actual: "api:2.0"}},
		},
		{
			ServiceName: "worker",
			Drifted:     false,
			Differences: nil,
		},
		{
			ServiceName: "proxy",
			Drifted:     true,
			Differences: []drift.Difference{{Field: "image", Expected: "nginx:1.21", Actual: "nginx:1.23"}},
		},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	results := makeResults()
	out := filter.Filter(results, filter.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func TestFilter_OnlyDrifted(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{OnlyDrifted: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(out))
	}
	for _, r := range out {
		if !r.Drifted {
			t.Errorf("expected only drifted results, got non-drifted: %s", r.ServiceName)
		}
	}
}

func TestFilter_ByServiceName(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{Services: []string{"api", "worker"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestFilter_ByServiceName_CaseInsensitive(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{Services: []string{"API"}})
	if len(out) != 1 || out[0].ServiceName != "api" {
		t.Fatalf("expected 1 result for 'API', got %d", len(out))
	}
}

func TestFilter_LabelSelector(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{LabelSelector: "image=nginx"})
	if len(out) != 1 || out[0].ServiceName != "proxy" {
		t.Fatalf("expected proxy result via label selector, got %d results", len(out))
	}
}

func TestFilter_LabelSelector_Invalid(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{LabelSelector: "badformat"})
	if len(out) != 0 {
		t.Fatalf("expected 0 results for invalid selector, got %d", len(out))
	}
}

func TestFilter_Combined(t *testing.T) {
	out := filter.Filter(makeResults(), filter.Options{
		OnlyDrifted: true,
		Services:    []string{"api"},
	})
	if len(out) != 1 || out[0].ServiceName != "api" {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}
