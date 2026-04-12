package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Snapshot represents a point-in-time capture of drift detection results.
type Snapshot struct {
	CapturedAt time.Time          `json:"captured_at"`
	Results    []drift.DriftResult `json:"results"`
}

// SnapshotWriter writes drift results to a snapshot file.
type SnapshotWriter struct {
	path string
}

// NewSnapshotWriter returns a SnapshotWriter that writes to the given path.
func NewSnapshotWriter(path string) *SnapshotWriter {
	return &SnapshotWriter{path: path}
}

// Write serialises results as a timestamped snapshot JSON file.
func (s *SnapshotWriter) Write(results []drift.DriftResult) error {
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", s.path, err)
	}

	return nil
}

// LoadSnapshot reads and deserialises a snapshot from disk.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &snap, nil
}
