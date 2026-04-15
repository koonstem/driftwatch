package output

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/yourusername/driftwatch/internal/drift"
)

// DedupeOptions controls deduplication behaviour.
type DedupeOptions struct {
	Enabled bool
}

// DefaultDedupeOptions returns sensible defaults (disabled).
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{Enabled: false}
}

// DedupeWriter wraps a Writer and suppresses renders whose result set
// is identical to the previously rendered one (by content hash).
type DedupeWriter struct {
	opts  DedupeOptions
	inner Writer
	mu    sync.Mutex
	last  string
}

// NewDedupeWriter creates a DedupeWriter.
func NewDedupeWriter(inner Writer, opts DedupeOptions) *DedupeWriter {
	return &DedupeWriter{opts: opts, inner: inner}
}

// Write renders results only when they differ from the last render.
func (d *DedupeWriter) Write(results []drift.Result) error {
	if !d.opts.Enabled {
		return d.inner.Write(results)
	}

	h, err := hashResults(results)
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if h == d.last {
		return nil
	}
	d.last = h
	return d.inner.Write(results)
}

// Reset clears the stored hash so the next write is always forwarded.
func (d *DedupeWriter) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = ""
}

func hashResults(results []drift.Result) (string, error) {
	b, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
