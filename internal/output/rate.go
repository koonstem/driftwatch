package output

import (
	"fmt"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// RateOptions controls rate-limiting behaviour for writers.
type RateOptions struct {
	Enabled     bool
	MaxPerMin   int
	BurstSize   int
}

// DefaultRateOptions returns sensible defaults (disabled).
func DefaultRateOptions() RateOptions {
	return RateOptions{
		Enabled:   false,
		MaxPerMin: 60,
		BurstSize: 5,
	}
}

// RateWriter wraps another Writer and enforces a token-bucket rate limit.
type RateWriter struct {
	opts    RateOptions
	next    Writer
	tokens  int
	lastFil time.Time
}

// NewRateWriter constructs a RateWriter. When disabled it is a transparent pass-through.
func NewRateWriter(opts RateOptions, next Writer) (*RateWriter, error) {
	if next == nil {
		return nil, fmt.Errorf("rate: next writer must not be nil")
	}
	if opts.MaxPerMin <= 0 {
		return nil, fmt.Errorf("rate: MaxPerMin must be > 0")
	}
	if opts.BurstSize <= 0 {
		return nil, fmt.Errorf("rate: BurstSize must be > 0")
	}
	return &RateWriter{
		opts:    opts,
		next:    next,
		tokens:  opts.BurstSize,
		lastFil: time.Now(),
	}, nil
}

// Write forwards results to the next writer, subject to rate limiting.
func (r *RateWriter) Write(results []drift.Result) error {
	if !r.opts.Enabled {
		return r.next.Write(results)
	}
	r.refill()
	if r.tokens <= 0 {
		return fmt.Errorf("rate: limit exceeded (%d/min), write suppressed", r.opts.MaxPerMin)
	}
	r.tokens--
	return r.next.Write(results)
}

// refill adds tokens proportional to elapsed time since last fill.
func (r *RateWriter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastFil)
	added := int(elapsed.Minutes() * float64(r.opts.MaxPerMin))
	if added > 0 {
		r.tokens += added
		if r.tokens > r.opts.BurstSize {
			r.tokens = r.opts.BurstSize
		}
		r.lastFil = now
	}
}

// Available returns the number of tokens currently available in the bucket.
// When rate limiting is disabled, it always returns the configured BurstSize.
func (r *RateWriter) Available() int {
	if !r.opts.Enabled {
		return r.opts.BurstSize
	}
	r.refill()
	return r.tokens
}
