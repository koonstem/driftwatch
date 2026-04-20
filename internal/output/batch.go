package output

import (
	"fmt"

	"github.com/user/driftwatch/internal/drift"
)

// DefaultBatchOptions returns sensible defaults for batching.
func DefaultBatchOptions() BatchOptions {
	return BatchOptions{
		Enabled:   false,
		BatchSize: 10,
	}
}

// BatchOptions controls how results are batched before forwarding.
type BatchOptions struct {
	Enabled   bool
	BatchSize int
}

// Validate returns an error if the options are invalid.
func (o BatchOptions) Validate() error {
	if o.Enabled && o.BatchSize < 1 {
		return fmt.Errorf("batch size must be at least 1, got %d", o.BatchSize)
	}
	return nil
}

// BatchWriter accumulates drift results and forwards them in batches.
type BatchWriter struct {
	opts    BatchOptions
	next    Writer
	buffer  []drift.DriftResult
}

// NewBatchWriter creates a BatchWriter that forwards to next in batches.
func NewBatchWriter(opts BatchOptions, next Writer) (*BatchWriter, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	if next == nil {
		return nil, fmt.Errorf("batch writer requires a non-nil next writer")
	}
	return &BatchWriter{opts: opts, next: next}, nil
}

// Write buffers results and flushes when the batch size is reached.
// If batching is disabled, results are forwarded immediately.
func (b *BatchWriter) Write(report drift.Report) error {
	if !b.opts.Enabled {
		return b.next.Write(report)
	}

	b.buffer = append(b.buffer, report.Results...)

	for len(b.buffer) >= b.opts.BatchSize {
		batch := b.buffer[:b.opts.BatchSize]
		b.buffer = b.buffer[b.opts.BatchSize:]
		if err := b.next.Write(drift.Report{Results: batch}); err != nil {
			return err
		}
	}
	return nil
}

// Flush forwards any remaining buffered results to the next writer.
func (b *BatchWriter) Flush() error {
	if len(b.buffer) == 0 {
		return nil
	}
	remaining := b.buffer
	b.buffer = nil
	return b.next.Write(drift.Report{Results: remaining})
}
