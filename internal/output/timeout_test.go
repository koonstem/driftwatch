package output

import (
	"errors"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func makeTimeoutResults() []drift.DriftResult {
	return []drift.DriftResult{
		{ServiceName: "api", Drifted: false, Fields: nil},
	}
}

func TestTimeoutWriter_Disabled_CallsImmediately(t *testing.T) {
	called := false
	w, err := NewTimeoutWriter(DefaultTimeoutOptions(), func(_ []drift.DriftResult) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeTimeoutResults()); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if !called {
		t.Error("expected downstream to be called")
	}
}

func TestTimeoutWriter_Enabled_CompletesWithinDeadline(t *testing.T) {
	opts := TimeoutOptions{Enabled: true, Duration: 500 * time.Millisecond}
	w, err := NewTimeoutWriter(opts, func(_ []drift.DriftResult) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeTimeoutResults()); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
}

func TestTimeoutWriter_Enabled_ExceedsDeadline(t *testing.T) {
	opts := TimeoutOptions{Enabled: true, Duration: 20 * time.Millisecond}
	w, err := NewTimeoutWriter(opts, func(_ []drift.DriftResult) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = w.Write(makeTimeoutResults())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !containsSubstring(err.Error(), "timeout") {
		t.Errorf("expected error to mention timeout, got: %v", err)
	}
}

func TestTimeoutWriter_PropagatesDownstreamError(t *testing.T) {
	downstreamErr := errors.New("downstream failure")
	opts := TimeoutOptions{Enabled: true, Duration: 500 * time.Millisecond}
	w, err := NewTimeoutWriter(opts, func(_ []drift.DriftResult) error {
		return downstreamErr
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeTimeoutResults()); !errors.Is(err, downstreamErr) {
		t.Errorf("expected downstream error, got: %v", err)
	}
}

func TestTimeoutWriter_InvalidDuration_ReturnsError(t *testing.T) {
	opts := TimeoutOptions{Enabled: true, Duration: -1 * time.Second}
	_, err := NewTimeoutWriter(opts, func(_ []drift.DriftResult) error { return nil })
	if err == nil {
		t.Fatal("expected validation error for negative duration")
	}
}

func TestTimeoutWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewTimeoutWriter(DefaultTimeoutOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil next writer")
	}
}

// containsSubstring is a local helper to avoid importing strings in tests.
func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
