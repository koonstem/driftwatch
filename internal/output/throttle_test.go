package output

import (
	"errors"
	"testing"
	"time"
)

func TestThrottleWriter_Disabled_AlwaysRenders(t *testing.T) {
	calls := 0
	render := func() error { calls++; return nil }

	w := NewThrottleWriter(DefaultThrottleOptions(), render)

	for i := 0; i < 5; i++ {
		called, err := w.Write()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Errorf("iteration %d: expected render to be called", i)
		}
	}

	if calls != 5 {
		t.Errorf("expected 5 render calls, got %d", calls)
	}
}

func TestThrottleWriter_Enabled_SuppressesWithinInterval(t *testing.T) {
	calls := 0
	render := func() error { calls++; return nil }

	opts := ThrottleOptions{Enabled: true, Interval: 10 * time.Second}
	w := NewThrottleWriter(opts, render)

	called1, err := w.Write()
	if err != nil || !called1 {
		t.Fatalf("first write should succeed: called=%v err=%v", called1, err)
	}

	called2, err := w.Write()
	if err != nil {
		t.Fatalf("unexpected error on second write: %v", err)
	}
	if called2 {
		t.Error("second write within interval should be suppressed")
	}

	if calls != 1 {
		t.Errorf("expected 1 render call, got %d", calls)
	}
}

func TestThrottleWriter_Enabled_AllowsAfterInterval(t *testing.T) {
	calls := 0
	render := func() error { calls++; return nil }

	opts := ThrottleOptions{Enabled: true, Interval: 1 * time.Millisecond}
	w := NewThrottleWriter(opts, render)

	w.Write() //nolint:errcheck
	time.Sleep(5 * time.Millisecond)

	called, err := w.Write()
	if err != nil || !called {
		t.Errorf("write after interval should be allowed: called=%v err=%v", called, err)
	}

	if calls != 2 {
		t.Errorf("expected 2 render calls, got %d", calls)
	}
}

func TestThrottleWriter_Reset_AllowsImmediateWrite(t *testing.T) {
	calls := 0
	render := func() error { calls++; return nil }

	opts := ThrottleOptions{Enabled: true, Interval: 10 * time.Second}
	w := NewThrottleWriter(opts, render)

	w.Write() //nolint:errcheck
	w.Reset()

	called, err := w.Write()
	if err != nil || !called {
		t.Errorf("write after reset should be allowed: called=%v err=%v", called, err)
	}

	if calls != 2 {
		t.Errorf("expected 2 render calls after reset, got %d", calls)
	}
}

func TestThrottleWriter_PropagatesRenderError(t *testing.T) {
	expected := errors.New("render failed")
	render := func() error { return expected }

	w := NewThrottleWriter(DefaultThrottleOptions(), render)

	_, err := w.Write()
	if !errors.Is(err, expected) {
		t.Errorf("expected render error to propagate, got: %v", err)
	}
}
