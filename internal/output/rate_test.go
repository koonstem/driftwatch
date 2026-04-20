package output

import (
	"errors"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeRateResults() []drift.Result {
	return []drift.Result{
		{Service: "svc-a", Status: drift.StatusMatch},
	}
}

func TestRateWriter_Disabled_AlwaysPasses(t *testing.T) {
	called := 0
	next := WriterFunc(func(r []drift.Result) error { called++; return nil })
	opts := DefaultRateOptions()
	opts.Enabled = false
	rw, err := NewRateWriter(opts, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 20; i++ {
		if err := rw.Write(makeRateResults()); err != nil {
			t.Fatalf("write %d failed: %v", i, err)
		}
	}
	if called != 20 {
		t.Errorf("expected 20 calls, got %d", called)
	}
}

func TestRateWriter_Enabled_AllowsUpToBurst(t *testing.T) {
	called := 0
	next := WriterFunc(func(r []drift.Result) error { called++; return nil })
	opts := RateOptions{Enabled: true, MaxPerMin: 60, BurstSize: 3}
	rw, err := NewRateWriter(opts, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := rw.Write(makeRateResults()); err != nil {
			t.Errorf("write %d should succeed: %v", i, err)
		}
	}
	if called != 3 {
		t.Errorf("expected 3 calls, got %d", called)
	}
}

func TestRateWriter_Enabled_BlocksAfterBurst(t *testing.T) {
	next := WriterFunc(func(r []drift.Result) error { return nil })
	opts := RateOptions{Enabled: true, MaxPerMin: 60, BurstSize: 2}
	rw, err := NewRateWriter(opts, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = rw.Write(makeRateResults())
	_ = rw.Write(makeRateResults())
	if err := rw.Write(makeRateResults()); err == nil {
		t.Error("expected rate limit error on third write")
	}
}

func TestRateWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewRateWriter(DefaultRateOptions(), nil)
	if err == nil {
		t.Error("expected error for nil next writer")
	}
}

func TestRateWriter_PropagatesNextError(t *testing.T) {
	next := WriterFunc(func(r []drift.Result) error { return errors.New("downstream failure") })
	opts := RateOptions{Enabled: true, MaxPerMin: 60, BurstSize: 5}
	rw, _ := NewRateWriter(opts, next)
	if err := rw.Write(makeRateResults()); err == nil {
		t.Error("expected downstream error to propagate")
	}
}
