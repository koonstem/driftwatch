package filter

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Options holds filtering criteria for drift results.
type Options struct {
	Services []string
	OnlyDrifted bool
	LabelSelector string
}

// Filter applies the given Options to a slice of DriftResult,
// returning only the results that match all specified criteria.
func Filter(results []drift.DriftResult, opts Options) []drift.DriftResult {
	var out []drift.DriftResult
	for _, r := range results {
		if opts.OnlyDrifted && !r.Drifted {
			continue
		}
		if len(opts.Services) > 0 && !matchesService(r.ServiceName, opts.Services) {
			continue
		}
		if opts.LabelSelector != "" && !matchesLabel(r, opts.LabelSelector) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func matchesService(name string, services []string) bool {
	for _, s := range services {
		if strings.EqualFold(s, name) {
			return true
		}
	}
	return false
}

func matchesLabel(r drift.DriftResult, selector string) bool {
	parts := strings.SplitN(selector, "=", 2)
	if len(parts) != 2 {
		return false
	}
	key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	for _, d := range r.Differences {
		if strings.EqualFold(d.Field, key) && strings.Contains(d.Actual, val) {
			return true
		}
	}
	return false
}
