package output

import (
	"math/rand"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

// DefaultSampleOptions returns conservative sampling defaults.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{
		Enabled: false,
		Rate:    1.0,
		Seed:    time.Now().UnixNano(),
	}
}

// SampleOptions controls probabilistic result sampling.
type SampleOptions struct {
	Enabled bool
	Rate    float64 // 0.0–1.0; fraction of results to forward
	Seed    int64
}

// SampleWriter randomly drops results based on a configured rate.
type SampleWriter struct {
	opts SampleOptions
	next func([]drift.DriftResult) error
	rng  *rand.Rand
}

// NewSampleWriter creates a SampleWriter that forwards sampled results to next.
func NewSampleWriter(opts SampleOptions, next func([]drift.DriftResult) error) (*SampleWriter, error) {
	if next == nil {
		return nil, errNilNext
	}
	if opts.Rate < 0 || opts.Rate > 1 {
		return nil, fmt.Errorf("sample rate must be between 0.0 and 1.0, got %f", opts.Rate)
	}
	return &SampleWriter{
		opts: opts,
		next: next,
		rng:  rand.New(rand.NewSource(opts.Seed)),
	}, nil
}

// Write forwards a sampled subset of results to the next writer.
func (s *SampleWriter) Write(results []drift.DriftResult) error {
	if !s.opts.Enabled {
		return s.next(results)
	}
	sampled := results[:0:0]
	for _, r := range results {
		if s.rng.Float64() < s.opts.Rate {
			sampled = append(sampled, r)
		}
	}
	return s.next(sampled)
}
