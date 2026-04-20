package output

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/drift"
)

func makeSampleResults(n int) []drift.DriftResult {
	out := make([]drift.DriftResult, n)
	for i := range out {
		out[i] = drift.DriftResult{Service: fmt.Sprintf("svc-%d", i)}
	}
	return out
}

func TestSampleWriter_Disabled_ForwardsAll(t *testing.T) {
	var got []drift.DriftResult
	w, err := NewSampleWriter(DefaultSampleOptions(), func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeSampleResults(10)
	if err := w.Write(input); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if len(got) != 10 {
		t.Errorf("expected 10 results, got %d", len(got))
	}
}

func TestSampleWriter_Enabled_RateZero_ForwardsNone(t *testing.T) {
	var got []drift.DriftResult
	opts := SampleOptions{Enabled: true, Rate: 0.0, Seed: 42}
	w, err := NewSampleWriter(opts, func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Write(makeSampleResults(100)); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 results at rate 0.0, got %d", len(got))
	}
}

func TestSampleWriter_Enabled_RateOne_ForwardsAll(t *testing.T) {
	var got []drift.DriftResult
	opts := SampleOptions{Enabled: true, Rate: 1.0, Seed: 42}
	w, _ := NewSampleWriter(opts, func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	input := makeSampleResults(50)
	w.Write(input)
	if len(got) != 50 {
		t.Errorf("expected 50 results at rate 1.0, got %d", len(got))
	}
}

func TestSampleWriter_InvalidRate_ReturnsError(t *testing.T) {
	opts := SampleOptions{Enabled: true, Rate: 1.5}
	_, err := NewSampleWriter(opts, func(r []drift.DriftResult) error { return nil })
	if err == nil {
		t.Fatal("expected error for rate > 1.0")
	}
}

func TestSampleWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewSampleWriter(DefaultSampleOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil next")
	}
}
