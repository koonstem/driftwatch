package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// HistoryEntry represents a single recorded drift detection run.
type HistoryEntry struct {
	Timestamp time.Time          `json:"timestamp"`
	Total     int                `json:"total"`
	Drifted   int                `json:"drifted"`
	Results   []drift.Result     `json:"results"`
}

// HistoryWriter writes drift results to a history log file.
type HistoryWriter struct {
	path    string
	maxSize int
}

// NewHistoryWriter returns a HistoryWriter that appends to the given file path.
// maxSize controls the maximum number of entries retained (0 = unlimited).
func NewHistoryWriter(path string, maxSize int) *HistoryWriter {
	return &HistoryWriter{path: path, maxSize: maxSize}
}

// Write appends a new entry for the given results to the history file.
func (h *HistoryWriter) Write(results []drift.Result) error {
	entries, err := LoadHistory(h.path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Total:     len(results),
		Drifted:   drifted,
		Results:   results,
	}
	entries = append(entries, entry)

	if h.maxSize > 0 && len(entries) > h.maxSize {
		entries = entries[len(entries)-h.maxSize:]
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(h.path, data, 0o644)
}

// LoadHistory reads all history entries from the given file.
func LoadHistory(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}
	return entries, nil
}

// SortHistoryByTime sorts entries ascending by timestamp.
func SortHistoryByTime(entries []HistoryEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
}
