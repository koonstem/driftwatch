package output

import (
	"errors"
	"testing"

	"github.com/yourusername/driftwatch/internal/drift"
)

func makeDedupeResults(service string, drifted bool) []drift.Result {
	return []drift.Result{
		{
			Service: service,
			Drifted: drifted,
			Fields:  []drift.FieldDiff{{Field: "image", Expected: "nginx:1", Actual: "nginx:2"}},
		},
	}
}

type countingWriter struct {
	calls  int
	lastFn func([]drift.Result) error
}

func (c *countingWriter) Write(r []drift.Result) error {
	c.calls++
	if c.lastFn != nil {
		return c.lastFn(r)
	}
	return nil
}

func TestDedupeWriter_Disabled_AlwaysForwards(t *testing.T) {
	inner := &countingWriter{}
	w := NewDedupeWriter(inner, DedupeOptions{Enabled: false})
	res := makeDedupeResults("svc-a", true)

	_ = w.Write(res)
	_ = w.Write(res)

	if inner.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls)
	}
}

func TestDedupeWriter_Enabled_SuppressesDuplicate(t *testing.T) {
	inner := &countingWriter{}
	w := NewDedupeWriter(inner, DedupeOptions{Enabled: true})
	res := makeDedupeResults("svc-a", true)

	_ = w.Write(res)
	_ = w.Write(res)

	if inner.calls != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls)
	}
}

func TestDedupeWriter_Enabled_ForwardsOnChange(t *testing.T) {
	inner := &countingWriter{}
	w := NewDedupeWriter(inner, DedupeOptions{Enabled: true})

	_ = w.Write(makeDedupeResults("svc-a", true))
	_ = w.Write(makeDedupeResults("svc-b", false))

	if inner.calls != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls)
	}
}

func TestDedupeWriter_Reset_AllowsImmediateWrite(t *testing.T) {
	inner := &countingWriter{}
	w := NewDedupeWriter(inner, DedupeOptions{Enabled: true})
	res := makeDedupeResults("svc-a", true)

	_ = w.Write(res)
	w.Reset()
	_ = w.Write(res)

	if inner.calls != 2 {
		t.Fatalf("expected 2 calls after reset, got %d", inner.calls)
	}
}

func TestDedupeWriter_PropagatesInnerError(t *testing.T) {
	expected := errors.New("write failed")
	inner := &countingWriter{lastFn: func(_ []drift.Result) error { return expected }}
	w := NewDedupeWriter(inner, DedupeOptions{Enabled: true})

	err := w.Write(makeDedupeResults("svc-a", true))
	if !errors.Is(err, expected) {
		t.Fatalf("expected propagated error, got %v", err)
	}
}
