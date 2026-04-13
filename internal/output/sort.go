package output

import (
	"sort"

	"github.com/driftwatch/internal/drift"
)

// SortField defines the field to sort drift results by.
type SortField string

const (
	SortByService SortField = "service"
	SortByStatus  SortField = "status"
	SortByField   SortField = "field"
)

// SortOptions controls how results are sorted.
type SortOptions struct {
	By      SortField
	Reverse bool
}

// DefaultSortOptions returns the default sort configuration.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		By:      SortByService,
		Reverse: false,
	}
}

// SortResults returns a sorted copy of the provided drift results.
func SortResults(results []drift.Result, opts SortOptions) []drift.Result {
	copied := make([]drift.Result, len(results))
	copy(copied, results)

	sort.SliceStable(copied, func(i, j int) bool {
		var less bool
		switch opts.By {
		case SortByStatus:
			less = statusRank(copied[i]) < statusRank(copied[j])
		case SortByField:
			less = firstField(copied[i]) < firstField(copied[j])
		default:
			less = copied[i].Service < copied[j].Service
		}
		if opts.Reverse {
			return !less
		}
		return less
	})

	return copied
}

func statusRank(r drift.Result) int {
	if r.Drifted {
		return 0
	}
	return 1
}

func firstField(r drift.Result) string {
	if len(r.Fields) > 0 {
		return r.Fields[0].Field
	}
	return ""
}
