package output

import (
	"sync"
	"time"
)

// ThrottleOptions controls how frequently output is emitted.
type ThrottleOptions struct {
	Enabled  bool
	Interval time.Duration
}

// DefaultThrottleOptions returns sensible defaults (disabled).
func DefaultThrottleOptions() ThrottleOptions {
	return ThrottleOptions{
		Enabled:  false,
		Interval: 5 * time.Second,
	}
}

// ThrottleWriter wraps a render function and suppresses calls that arrive
// more frequently than the configured interval.
type ThrottleWriter struct {
	opts    ThrottleOptions
	mu      sync.Mutex
	lastRun time.Time
	render  func() error
}

// NewThrottleWriter creates a ThrottleWriter. render is the function that
// produces output; it will only be called if the throttle interval has elapsed
// since the last successful call (or if throttling is disabled).
func NewThrottleWriter(opts ThrottleOptions, render func() error) *ThrottleWriter {
	return &ThrottleWriter{
		opts:   opts,
		render: render,
	}
}

// Write attempts to invoke the underlying render function, honouring the
// throttle interval. It returns (true, nil) when render was called, and
// (false, nil) when the call was suppressed. Any error from render is
// propagated directly.
func (t *ThrottleWriter) Write() (bool, error) {
	if !t.opts.Enabled {
		if err := t.render(); err != nil {
			return false, err
		}
		return true, nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if !t.lastRun.IsZero() && now.Sub(t.lastRun) < t.opts.Interval {
		return false, nil
	}

	if err := t.render(); err != nil {
		return false, err
	}

	t.lastRun = now
	return true, nil
}

// Reset clears the last-run timestamp so the next Write call is always
// allowed through regardless of the interval.
func (t *ThrottleWriter) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastRun = time.Time{}
}
