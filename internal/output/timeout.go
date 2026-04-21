package output

import (
	"context"
	"fmt"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// DefaultTimeoutOptions returns sensible defaults for the timeout writer.
func DefaultTimeoutOptions() TimeoutOptions {
	return TimeoutOptions{
		Enabled: false,
		Duration: 10 * time.Second,
	}
}

// TimeoutOptions controls how long a downstream writer is allowed to run.
type TimeoutOptions struct {
	Enabled  bool
	Duration time.Duration
}

// Validate returns an error if the options are invalid.
func (o TimeoutOptions) Validate() error {
	if o.Enabled && o.Duration <= 0 {
		return fmt.Errorf("timeout duration must be positive, got %s", o.Duration)
	}
	return nil
}

// TimeoutWriter wraps a downstream Writer and cancels it if it exceeds the
// configured duration.
type TimeoutWriter struct {
	opts TimeoutOptions
	next func([]drift.DriftResult) error
}

// NewTimeoutWriter creates a TimeoutWriter that delegates to next.
func NewTimeoutWriter(opts TimeoutOptions, next func([]drift.DriftResult) error) (*TimeoutWriter, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	if next == nil {
		return nil, fmt.Errorf("timeout: next writer must not be nil")
	}
	return &TimeoutWriter{opts: opts, next: next}, nil
}

// Write calls the downstream writer, enforcing the timeout when enabled.
func (w *TimeoutWriter) Write(results []drift.DriftResult) error {
	if !w.opts.Enabled {
		return w.next(results)
	}

	ctx, cancel := context.WithTimeout(context.Background(), w.opts.Duration)
	defer cancel()

	type result struct{ err error }
	ch := make(chan result, 1)

	go func() {
		ch <- result{err: w.next(results)}
	}()

	select {
	case r := <-ch:
		return r.err
	case <-ctx.Done():
		return fmt.Errorf("timeout: writer exceeded %s deadline", w.opts.Duration)
	}
}
