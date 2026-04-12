package output

import "github.com/driftwatch/driftwatch/internal/drift"

// BaselineDiff describes how a service's drift state changed vs a baseline.
type BaselineDiff struct {
	Service    string
	WasDrifted bool
	NowDrifted bool
	New        bool // service not present in baseline
	Gone       bool // service present in baseline but not in current results
}

// Changed returns true when the drift state has actually changed.
func (d BaselineDiff) Changed() bool {
	return d.New || d.Gone || d.WasDrifted != d.NowDrifted
}

// CompareToBaseline diffs current drift results against a saved baseline.
// It returns one BaselineDiff per service that appears in either set.
func CompareToBaseline(baseline *Baseline, results []drift.Result) []BaselineDiff {
	baselineMap := make(map[string]BaselineEntry, len(baseline.Entries))
	for _, e := range baseline.Entries {
		baselineMap[e.Service] = e
	}

	currentMap := make(map[string]drift.Result, len(results))
	for _, r := range results {
		currentMap[r.Service] = r
	}

	seen := make(map[string]struct{})
	var diffs []BaselineDiff

	for _, r := range results {
		seen[r.Service] = struct{}{}
		entry, exists := baselineMap[r.Service]
		diff := BaselineDiff{
			Service:    r.Service,
			NowDrifted: r.Drifted,
			New:        !exists,
		}
		if exists {
			diff.WasDrifted = entry.Drifted
		}
		diffs = append(diffs, diff)
	}

	for _, e := range baseline.Entries {
		if _, ok := seen[e.Service]; !ok {
			diffs = append(diffs, BaselineDiff{
				Service:    e.Service,
				WasDrifted: e.Drifted,
				Gone:       true,
			})
		}
	}

	return diffs
}
