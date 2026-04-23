package output

import (
	"fmt"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeEnvelopeResults() []drift.DriftResult {
	return []drift.DriftResult{
		{Service: "api", Drifted: true, Fields: map[string]string{"image": "nginx:1.24"}},
		{Service: "worker", Drifted: false, Fields: map[string]string{}},
	}
}

func fixedTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestEnvelopeWriter_Disabled_PassesThrough(t *testing.T) {
	var got []drift.DriftResult
	w, err := NewEnvelopeWriter(DefaultEnvelopeOptions(), func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeEnvelopeResults()
	if err := w.Write(input); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if len(got) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(got))
	}
	if _, ok := got[0].Fields["envelope.timestamp"]; ok {
		t.Error("expected no envelope metadata when disabled")
	}
}

func TestEnvelopeWriter_Enabled_AnnotatesFields(t *testing.T) {
	opts := EnvelopeOptions{
		Enabled: true,
		Source:  "ci",
		RunID:   "run-42",
		AddedAt: fixedTime,
	}
	var got []drift.DriftResult
	w, _ := NewEnvelopeWriter(opts, func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	_ = w.Write(makeEnvelopeResults())

	for _, r := range got {
		if r.Fields["envelope.run_id"] != "run-42" {
			t.Errorf("service %s: expected run_id=run-42, got %q", r.Service, r.Fields["envelope.run_id"])
		}
		if r.Fields["envelope.source"] != "ci" {
			t.Errorf("service %s: expected source=ci, got %q", r.Service, r.Fields["envelope.source"])
		}
		if r.Fields["envelope.timestamp"] != "2024-06-01T12:00:00Z" {
			t.Errorf("service %s: unexpected timestamp %q", r.Service, r.Fields["envelope.timestamp"])
		}
	}
}

func TestEnvelopeWriter_Enabled_PreservesExistingFields(t *testing.T) {
	opts := EnvelopeOptions{Enabled: true, RunID: "x", AddedAt: fixedTime}
	var got []drift.DriftResult
	w, _ := NewEnvelopeWriter(opts, func(r []drift.DriftResult) error {
		got = r
		return nil
	})
	_ = w.Write(makeEnvelopeResults())
	if got[0].Fields["image"] != "nginx:1.24" {
		t.Error("existing fields should not be overwritten")
	}
}

func TestEnvelopeWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewEnvelopeWriter(DefaultEnvelopeOptions(), nil)
	if err == nil {
		t.Error("expected error for nil next writer")
	}
}

func TestEnvelopeWriter_PropagatesNextError(t *testing.T) {
	opts := EnvelopeOptions{Enabled: true, AddedAt: fixedTime}
	w, _ := NewEnvelopeWriter(opts, func(_ []drift.DriftResult) error {
		return fmt.Errorf("downstream failure")
	})
	if err := w.Write(makeEnvelopeResults()); err == nil {
		t.Error("expected propagated error from next writer")
	}
}
