package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// DefaultReplayOptions returns sensible defaults for ReplayWriter.
func DefaultReplayOptions() ReplayOptions {
	return ReplayOptions{
		Enabled: false,
		File:    "",
		Delay:   0,
	}
}

// ReplayOptions controls how historical drift results are replayed.
type ReplayOptions struct {
	Enabled bool
	File    string
	Delay   time.Duration
}

// ReplayWriter reads a previously saved history file and replays each
// snapshot through the downstream writer, optionally inserting a delay
// between entries.
type ReplayWriter struct {
	opts ReplayOptions
	next func([]drift.DriftResult) error
}

// NewReplayWriter constructs a ReplayWriter.
func NewReplayWriter(opts ReplayOptions, next func([]drift.DriftResult) error) *ReplayWriter {
	return &ReplayWriter{opts: opts, next: next}
}

// Run loads the history file and replays each snapshot in order.
func (r *ReplayWriter) Run() error {
	if !r.opts.Enabled {
		return nil
	}
	if r.opts.File == "" {
		return fmt.Errorf("replay: no file specified")
	}

	entries, err := LoadHistory(r.opts.File)
	if err != nil {
		return fmt.Errorf("replay: load history: %w", err)
	}

	for _, entry := range entries {
		if err := r.next(entry.Results); err != nil {
			return fmt.Errorf("replay: downstream error: %w", err)
		}
		if r.opts.Delay > 0 {
			time.Sleep(r.opts.Delay)
		}
	}
	return nil
}

// LoadReplayFile is a convenience wrapper that reads raw history JSON from
// disk and returns the slice of HistoryEntry values.
func LoadReplayFile(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("replay: file not found: %s", path)
		}
		return nil, err
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("replay: invalid JSON: %w", err)
	}
	return entries, nil
}
