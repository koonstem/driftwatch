package output

import (
	"errors"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
)

func makeCircuitResults() []drift.Result {
	return []drift.Result{
		{Service: "api", Drifted: true, Fields: []drift.FieldDiff{{Field: "image", Expected: "v1", Actual: "v2"}}},
	}
}

type failingWriter struct{ calls int }

func (f *failingWriter) Write(_ []drift.Result) error {
	f.calls++
	return errors.New("downstream error")
}

type countingWriter struct{ calls int }

func (c *countingWriter) Write(_ []drift.Result) error {
	c.calls++
	return nil
}

func TestCircuitWriter_Disabled_AlwaysForwards(t *testing.T) {
	next := &countingWriter{}
	w, _ := NewCircuitWriter(CircuitOptions{Enabled: false, MaxFailures: 1}, next)
	for i := 0; i < 5; i++ {
		_ = w.Write(makeCircuitResults())
	}
	if next.calls != 5 {
		t.Fatalf("expected 5 calls, got %d", next.calls)
	}
}

func TestCircuitWriter_OpensAfterMaxFailures(t *testing.T) {
	next := &failingWriter{}
	opts := CircuitOptions{Enabled: true, MaxFailures: 2, ResetTimeout: 10 * time.Second}
	w, _ := NewCircuitWriter(opts, next)

	_ = w.Write(makeCircuitResults()) // failure 1
	_ = w.Write(makeCircuitResults()) // failure 2 — circuit opens

	err := w.Write(makeCircuitResults()) // should be blocked
	if err == nil {
		t.Fatal("expected circuit open error, got nil")
	}
	if next.calls != 2 {
		t.Fatalf("expected exactly 2 downstream calls, got %d", next.calls)
	}
}

func TestCircuitWriter_ResetsAfterTimeout(t *testing.T) {
	next := &failingWriter{}
	opts := CircuitOptions{Enabled: true, MaxFailures: 1, ResetTimeout: 1 * time.Millisecond}
	w, _ := NewCircuitWriter(opts, next)

	_ = w.Write(makeCircuitResults()) // opens circuit
	time.Sleep(5 * time.Millisecond)

	// Half-open probe — still fails, re-opens
	err := w.Write(makeCircuitResults())
	if err == nil {
		t.Fatal("expected error on half-open probe that fails")
	}
}

func TestCircuitWriter_ClosesOnSuccess(t *testing.T) {
	calls := 0
	next := writerFunc(func(_ []drift.Result) error {
		calls++
		return nil
	})
	opts := CircuitOptions{Enabled: true, MaxFailures: 3, ResetTimeout: time.Second}
	w, _ := NewCircuitWriter(opts, next)

	for i := 0; i < 5; i++ {
		if err := w.Write(makeCircuitResults()); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}
	if calls != 5 {
		t.Fatalf("expected 5 calls, got %d", calls)
	}
}

func TestCircuitWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewCircuitWriter(DefaultCircuitOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil next writer")
	}
}
