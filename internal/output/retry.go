package output

import (
	"fmt"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// RetryOptions configures retry behaviour for transient write failures.
type RetryOptions struct {
	Enabled     bool
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64
}

// DefaultRetryOptions returns sensible defaults for retry behaviour.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		Enabled:     false,
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Backoff:     2.0,
	}
}

// RetryWriter wraps a Renderer and retries on transient errors.
type RetryWriter struct {
	inner Renderer
	opts  RetryOptions
}

// NewRetryWriter creates a RetryWriter wrapping the given Renderer.
func NewRetryWriter(inner Renderer, opts RetryOptions) *RetryWriter {
	return &RetryWriter{inner: inner, opts: opts}
}

// Render attempts to render results, retrying on failure if enabled.
func (r *RetryWriter) Render(results []drift.Result) error {
	if !r.opts.Enabled {
		return r.inner.Render(results)
	}

	delay := r.opts.Delay
	var lastErr error

	for attempt := 1; attempt <= r.opts.MaxAttempts; attempt++ {
		if err := r.inner.Render(results); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < r.opts.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * r.opts.Backoff)
		}
	}

	return fmt.Errorf("render failed after %d attempts: %w", r.opts.MaxAttempts, lastErr)
}
