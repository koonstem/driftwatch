package output

import (
	"github.com/yourorg/driftwatch/internal/drift"
)

// SnapshotDiff describes changes between two snapshots.
type SnapshotDiff struct {
	// New contains services that appear in current but not in previous.
	New []drift.DriftResult
	// Resolved contains services that were drifted in previous but are clean now.
	Resolved []drift.DriftResult
	// Changed contains services whose drift fields changed between snapshots.
	Changed []drift.DriftResult
	// Unchanged contains services with identical drift state in both snapshots.
	Unchanged []drift.DriftResult
}

// CompareSnapshots returns a diff between a previous and current snapshot.
func CompareSnapshots(previous, current *Snapshot) SnapshotDiff {
	prev := indexByService(previous.Results)
	var diff SnapshotDiff

	for _, cur := range current.Results {
		old, exists := prev[cur.Service]
		if !exists {
			diff.New = append(diff.New, cur)
			continue
		}
		delete(prev, cur.Service)

		if fieldsChanged(old, cur) {
			diff.Changed = append(diff.Changed, cur)
		} else {
			diff.Unchanged = append(diff.Unchanged, cur)
		}
	}

	// Anything left in prev was present before but absent now — resolved.
	for _, r := range prev {
		diff.Resolved = append(diff.Resolved, r)
	}

	return diff
}

func indexByService(results []drift.DriftResult) map[string]drift.DriftResult {
	m := make(map[string]drift.DriftResult, len(results))
	for _, r := range results {
		m[r.Service] = r
	}
	return m
}

func fieldsChanged(a, b drift.DriftResult) bool {
	if a.Drifted != b.Drifted || len(a.Fields) != len(b.Fields) {
		return true
	}
	aFields := make(map[string]string, len(a.Fields))
	for _, f := range a.Fields {
		aFields[f.Field] = f.Actual
	}
	for _, f := range b.Fields {
		if aFields[f.Field] != f.Actual {
			return true
		}
	}
	return false
}
