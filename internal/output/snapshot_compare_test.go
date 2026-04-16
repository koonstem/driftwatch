package output

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

func snap(results []drift.DriftResult) *Snapshot {
	return &Snapshot{CapturedAt: time.Now().UTC(), Results: results}
}

func TestCompareSnapshots_NewService(t *testing.T) {
	prev := snap([]drift.DriftResult{})
	curr := snap([]drift.DriftResult{{Service: "api", Drifted: true}})

	diff := CompareSnapshots(prev, curr)
	if len(diff.New) != 1 || diff.New[0].Service != "api" {
		t.Errorf("expected 1 new service, got %+v", diff.New)
	}
}

func TestCompareSnapshots_ResolvedService(t *testing.T) {
	prev := snap([]drift.DriftResult{{Service: "api", Drifted: true}})
	curr := snap([]drift.DriftResult{})

	diff := CompareSnapshots(prev, curr)
	if len(diff.Resolved) != 1 || diff.Resolved[0].Service != "api" {
		t.Errorf("expected 1 resolved service, got %+v", diff.Resolved)
	}
}

func TestCompareSnapshots_ChangedService(t *testing.T) {
	prev := snap([]drift.DriftResult{
		{Service: "api", Drifted: true, Fields: []drift.FieldDrift{{Field: "image", Actual: "api:1"}}},
	})
	curr := snap([]drift.DriftResult{
		{Service: "api", Drifted: true, Fields: []drift.FieldDrift{{Field: "image", Actual: "api:2"}}},
	})

	diff := CompareSnapshots(prev, curr)
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed service, got %+v", diff.Changed)
	}
}

func TestCompareSnapshots_UnchangedService(t *testing.T) {
	result := drift.DriftResult{Service: "worker", Drifted: false}
	prev := snap([]drift.DriftResult{result})
	curr := snap([]drift.DriftResult{result})

	diff := CompareSnapshots(prev, curr)
	if len(diff.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged service, got %+v", diff.Unchanged)
	}
	if len(diff.New) != 0 || len(diff.Changed) != 0 || len(diff.Resolved) != 0 {
		t.Error("expected no new/changed/resolved services")
	}
}

func TestCompareSnapshots_EmptySnapshots(t *testing.T) {
	prev := snap([]drift.DriftResult{})
	curr := snap([]drift.DriftResult{})

	diff := CompareSnapshots(prev, curr)
	if len(diff.New) != 0 || len(diff.Resolved) != 0 || len(diff.Changed) != 0 || len(diff.Unchanged) != 0 {
		t.Errorf("expected empty diff for empty snapshots, got %+v", diff)
	}
}
