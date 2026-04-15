package output

import (
	"errors"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

type countingRenderer struct {
	calls    int
	failUntil int
	err      error
}

func (c *countingRenderer) Render(_ []drift.Result) error {
	c.calls++
	if c.calls <= c.failUntil {
		return c.err
	}
	return nil
}

func newRetryOpts(enabled bool, maxAttempts int) RetryOptions {
	return RetryOptions{
		Enabled:     enabled,
		MaxAttempts: maxAttempts,
		Delay:       time.Millisecond,
		Backoff:     1.0,
	}
}

func TestRetryWriter_Disabled_CallsOnce(t *testing.T) {
	inner := &countingRenderer{failUntil: 2, err: errors.New("fail")}
	w := NewRetryWriter(inner, newRetryOpts(false, 3))

	err := w.Render(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if inner.calls != 1 {
		t.Errorf("expected 1 call, got %d", inner.calls)
	}
}

func TestRetryWriter_Enabled_SucceedsOnRetry(t *testing.T) {
	inner := &countingRenderer{failUntil: 2, err: errors.New("transient")}
	w := NewRetryWriter(inner, newRetryOpts(true, 3))

	if err := w.Render(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 3 {
		t.Errorf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryWriter_Enabled_ExceedsMaxAttempts(t *testing.T) {
	inner := &countingRenderer{failUntil: 10, err: errors.New("persistent")}
	w := NewRetryWriter(inner, newRetryOpts(true, 3))

	err := w.Render(nil)
	if err == nil {
		t.Fatal("expected error after max attempts")
	}
	if inner.calls != 3 {
		t.Errorf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryWriter_Enabled_SucceedsFirstAttempt(t *testing.T) {
	inner := &countingRenderer{failUntil: 0}
	w := NewRetryWriter(inner, newRetryOpts(true, 3))

	if err := w.Render(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 1 {
		t.Errorf("expected 1 call, got %d", inner.calls)
	}
}

func TestDefaultRetryOptions(t *testing.T) {
	opts := DefaultRetryOptions()
	if opts.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if opts.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", opts.MaxAttempts)
	}
	if opts.Backoff != 2.0 {
		t.Errorf("expected Backoff=2.0, got %f", opts.Backoff)
	}
}
