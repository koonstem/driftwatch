package output

import (
	"fmt"
	"sort"

	"github.com/user/driftwatch/internal/drift"
)

// RollupOptions controls how drift results are rolled up into a summary.
type RollupOptions struct {
	Enabled   bool
	GroupBy   string // "service" or "field"
	TopN      int
}

// DefaultRollupOptions returns sensible defaults.
func DefaultRollupOptions() RollupOptions {
	return RollupOptions{
		Enabled: false,
		GroupBy: "service",
		TopN:    10,
	}
}

// RollupEntry represents a single rolled-up group.
type RollupEntry struct {
	Key        string
	DriftCount int
	Services   []string
}

// NewRollupWriter returns a Writer that summarises results into rollup groups
// before forwarding to next.
func NewRollupWriter(opts RollupOptions, next Writer) (Writer, error) {
	if next == nil {
		return nil, fmt.Errorf("rollup: next writer must not be nil")
	}
	if opts.GroupBy != "service" && opts.GroupBy != "field" {
		return nil, fmt.Errorf("rollup: invalid group_by %q, must be 'service' or 'field'", opts.GroupBy)
	}
	return writerFunc(func(report drift.Report) error {
		if !opts.Enabled {
			return next.Write(report)
		}
		entries := rollup(report.Results, opts)
		report.Annotations = appendRollupAnnotations(report.Annotations, entries)
		return next.Write(report)
	}), nil
}

func rollup(results []drift.Result, opts RollupOptions) []RollupEntry {
	counts := map[string]*RollupEntry{}
	for _, r := range results {
		if !r.Drifted {
			continue
		}
		if opts.GroupBy == "service" {
			e := getOrCreate(counts, r.Service)
			e.DriftCount += len(r.Fields)
			e.Services = append(e.Services, r.Service)
		} else {
			for _, f := range r.Fields {
				e := getOrCreate(counts, f.Name)
				e.DriftCount++
				e.Services = appendUnique(e.Services, r.Service)
			}
		}
	}
	entries := make([]RollupEntry, 0, len(counts))
	for _, e := range counts {
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].DriftCount > entries[j].DriftCount
	})
	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}
	return entries
}

func getOrCreate(m map[string]*RollupEntry, key string) *RollupEntry {
	if e, ok := m[key]; ok {
		return e
	}
	e := &RollupEntry{Key: key}
	m[key] = e
	return e
}

func appendUnique(ss []string, s string) []string {
	for _, v := range ss {
		if v == s {
			return ss
		}
	}
	return append(ss, s)
}

func appendRollupAnnotations(existing map[string]string, entries []RollupEntry) map[string]string {
	if existing == nil {
		existing = map[string]string{}
	}
	for i, e := range entries {
		existing[fmt.Sprintf("rollup.%d.key", i)] = e.Key
		existing[fmt.Sprintf("rollup.%d.drift_count", i)] = fmt.Sprintf("%d", e.DriftCount)
	}
	return existing
}
