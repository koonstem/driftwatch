package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// AuditEntry represents a single audit log record.
type AuditEntry struct {
	Timestamp   time.Time          `json:"timestamp"`
	RunID       string             `json:"run_id"`
	TotalChecked int              `json:"total_checked"`
	DriftedCount int              `json:"drifted_count"`
	Results     []drift.Result     `json:"results"`
}

// AuditOptions controls audit log behaviour.
type AuditOptions struct {
	Enabled  bool
	FilePath string
	RunID    string
}

// DefaultAuditOptions returns sensible defaults.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		Enabled:  false,
		FilePath: "driftwatch-audit.jsonl",
		RunID:    fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

// AuditWriter appends structured audit entries to a JSONL file.
type AuditWriter struct {
	opts AuditOptions
}

// NewAuditWriter creates an AuditWriter with the given options.
func NewAuditWriter(opts AuditOptions) *AuditWriter {
	return &AuditWriter{opts: opts}
}

// Write appends an audit entry for the provided results if enabled.
func (a *AuditWriter) Write(results []drift.Result) error {
	if !a.opts.Enabled {
		return nil
	}

	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	entry := AuditEntry{
		Timestamp:    time.Now().UTC(),
		RunID:        a.opts.RunID,
		TotalChecked: len(results),
		DriftedCount: drifted,
		Results:      results,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	f, err := os.OpenFile(a.opts.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}
