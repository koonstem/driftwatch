package output

import (
	"errors"
	"testing"

	"github.com/user/driftwatch/internal/drift"
)

type stubWriter struct {
	called bool
	err    error
}

func (s *stubWriter) Write(_ []drift.DriftResult) error {
	s.called = true
	return s.err
}

func makePipelineResults() []drift.DriftResult {
	return []drift.DriftResult{{Service: "svc-a", Drifted: false}}
}

func TestPipeline_AllSucceed(t *testing.T) {
	a, b := &stubWriter{}, &stubWriter{}
	opts := DefaultPipelineOptions()
	opts.Writers = []Writer{a, b}
	p := NewPipeline(opts)
	if err := p.Run(makePipelineResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both writers to be called")
	}
}

func TestPipeline_StopOnError(t *testing.T) {
	a := &stubWriter{err: errors.New("fail")}
	b := &stubWriter{}
	opts := DefaultPipelineOptions()
	opts.StopOnError = true
	opts.Writers = []Writer{a, b}
	p := NewPipeline(opts)
	if err := p.Run(makePipelineResults()); err == nil {
		t.Fatal("expected error")
	}
	if b.called {
		t.Error("second writer should not be called when StopOnError=true")
	}
}

func TestPipeline_ContinueOnError(t *testing.T) {
	a := &stubWriter{err: errors.New("fail")}
	b := &stubWriter{}
	opts := DefaultPipelineOptions()
	opts.StopOnError = false
	opts.Writers = []Writer{a, b}
	p := NewPipeline(opts)
	if err := p.Run(makePipelineResults()); err == nil {
		t.Fatal("expected combined error")
	}
	if !b.called {
		t.Error("second writer should still be called when StopOnError=false")
	}
}

func TestPipeline_NoWriters_ReturnsNil(t *testing.T) {
	p := NewPipeline(DefaultPipelineOptions())
	if err := p.Run(makePipelineResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
