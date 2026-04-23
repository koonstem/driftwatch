package output

import (
	"fmt"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// EnvelopeOptions controls metadata wrapping around drift results.
type EnvelopeOptions struct {
	Enabled   bool
	Source    string
	RunID     string
	AddedAt   func() time.Time
}

// DefaultEnvelopeOptions returns sensible defaults.
func DefaultEnvelopeOptions() EnvelopeOptions {
	return EnvelopeOptions{
		Enabled: false,
		AddedAt: time.Now,
	}
}

// Envelope wraps a drift result with metadata fields.
type Envelope struct {
	RunID     string            `json:"run_id,omitempty"`
	Source    string            `json:"source,omitempty"`
	Timestamp string            `json:"timestamp"`
	Result    drift.DriftResult `json:"result"`
}

// EnvelopeWriter wraps each DriftResult with metadata before forwarding.
type EnvelopeWriter struct {
	opts EnvelopeOptions
	next func([]drift.DriftResult) error
}

// NewEnvelopeWriter constructs an EnvelopeWriter.
func NewEnvelopeWriter(opts EnvelopeOptions, next func([]drift.DriftResult) error) (*EnvelopeWriter, error) {
	if next == nil {
		return nil, fmt.Errorf("envelope: next writer must not be nil")
	}
	return &EnvelopeWriter{opts: opts, next: next}, nil
}

// Write annotates results with envelope metadata and forwards them.
func (e *EnvelopeWriter) Write(results []drift.DriftResult) error {
	if !e.opts.Enabled {
		return e.next(results)
	}

	ts := e.opts.AddedAt().UTC().Format(time.RFC3339)
	annotated := make([]drift.DriftResult, len(results))
	for i, r := range results {
		copy := r
		if copy.Fields == nil {
			copy.Fields = map[string]string{}
		}
		copy.Fields["envelope.run_id"] = e.opts.RunID
		copy.Fields["envelope.source"] = e.opts.Source
		copy.Fields["envelope.timestamp"] = ts
		annotated[i] = copy
	}
	return e.next(annotated)
}
