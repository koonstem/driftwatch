package output

import (
	"fmt"
	"sort"

	"github.com/yourorg/driftwatch/internal/drift"
)

// AggregateOptions controls how results are aggregated.
type AggregateOptions struct {
	GroupBy string // "service" or "field"
}

// DefaultAggregateOptions returns sensible defaults.
func DefaultAggregateOptions() AggregateOptions {
	return AggregateOptions{GroupBy: "service"}
}

// AggregateResult holds a grouped summary.
type AggregateResult struct {
	Key        string
	Total      int
	Drifted    int
	Fields     []string
}

// NewAggregator returns a function that aggregates drift results.
func NewAggregator(opts AggregateOptions) func([]drift.DriftResult) ([]AggregateResult, error) {
	return func(results []drift.DriftResult) ([]AggregateResult, error) {
		switch opts.GroupBy {
		case "field":
			return aggregateByField(results), nil
		case "service", "":
			return aggregateByService(results), nil
		default:
			return nil, fmt.Errorf("unknown group-by value: %q", opts.GroupBy)
		}
	}
}

func aggregateByService(results []drift.DriftResult) []AggregateResult {
	out := make([]AggregateResult, 0, len(results))
	for _, r := range results {
		ar := AggregateResult{
			Key:   r.Service,
			Total: len(r.Fields),
		}
		for _, f := range r.Fields {
			if f.Drifted {
				ar.Drifted++
				ar.Fields = append(ar.Fields, f.Name)
			}
		}
		out = append(out, ar)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func aggregateByField(results []drift.DriftResult) []AggregateResult {
	counts := map[string]*AggregateResult{}
	for _, r := range results {
		for _, f := range r.Fields {
			if _, ok := counts[f.Name]; !ok {
				counts[f.Name] = &AggregateResult{Key: f.Name}
			}
			counts[f.Name].Total++
			if f.Drifted {
				counts[f.Name].Drifted++
			}
		}
	}
	out := make([]AggregateResult, 0, len(counts))
	for _, v := range counts {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}
