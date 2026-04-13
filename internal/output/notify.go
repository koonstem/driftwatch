package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// NotifyOptions controls notification behaviour.
type NotifyOptions struct {
	Enabled   bool
	OnlyDrift bool
	Channels  []string
}

// DefaultNotifyOptions returns sensible defaults.
func DefaultNotifyOptions() NotifyOptions {
	return NotifyOptions{
		Enabled:   false,
		OnlyDrift: true,
		Channels:  []string{"stdout"},
	}
}

// NotifyWriter writes human-readable notification summaries to one or more
// channels (currently only stdout/stderr writers are supported).
type NotifyWriter struct {
	opts    NotifyOptions
	writers []io.Writer
}

// NewNotifyWriter creates a NotifyWriter that fans out to the supplied writers.
func NewNotifyWriter(opts NotifyOptions, writers ...io.Writer) *NotifyWriter {
	return &NotifyWriter{opts: opts, writers: writers}
}

// Notify emits a notification for the given results if conditions are met.
func (n *NotifyWriter) Notify(results []drift.DriftResult) error {
	if !n.opts.Enabled {
		return nil
	}

	hasDrift := false
	for _, r := range results {
		if r.Drifted {
			hasDrift = true
			break
		}
	}

	if n.opts.OnlyDrift && !hasDrift {
		return nil
	}

	msg := buildNotifyMessage(results)
	for _, w := range n.writers {
		if _, err := fmt.Fprintln(w, msg); err != nil {
			return fmt.Errorf("notify write: %w", err)
		}
	}
	return nil
}

func buildNotifyMessage(results []drift.DriftResult) string {
	var drifted []string
	for _, r := range results {
		if r.Drifted {
			drifted = append(drifted, r.ServiceName)
		}
	}
	if len(drifted) == 0 {
		return "[driftwatch] No drift detected."
	}
	return fmt.Sprintf("[driftwatch] Drift detected in %d service(s): %s",
		len(drifted), strings.Join(drifted, ", "))
}
