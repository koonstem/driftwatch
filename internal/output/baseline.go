package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// BaselineEntry represents a snapshot of a single service's drift state.
type BaselineEntry struct {
	Service   string            `json:"service"`
	Image     string            `json:"image"`
	Drifted   bool              `json:"drifted"`
	Reasons   []string          `json:"reasons,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Baseline is the top-level structure written to a baseline file.
type Baseline struct {
	GeneratedAt time.Time       `json:"generated_at"`
	Entries     []BaselineEntry `json:"entries"`
}

// BaselineWriter writes drift results to a baseline snapshot file.
type BaselineWriter struct {
	path string
}

// NewBaselineWriter creates a BaselineWriter that writes to the given path.
func NewBaselineWriter(path string) *BaselineWriter {
	return &BaselineWriter{path: path}
}

// Write serialises the drift results into a baseline JSON file.
func (b *BaselineWriter) Write(results []drift.Result) error {
	entries := make([]BaselineEntry, 0, len(results))
	for _, r := range results {
		entries = append(entries, BaselineEntry{
			Service: r.Service,
			Image:   r.ExpectedImage,
			Drifted: r.Drifted,
			Reasons: r.Reasons,
			Labels:  r.Labels,
		})
	}

	bl := Baseline{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
	}

	data, err := json.MarshalIndent(bl, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}

	if err := os.WriteFile(b.path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", b.path, err)
	}
	return nil
}

// LoadBaseline reads a previously written baseline file.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var bl Baseline
	if err := json.Unmarshal(data, &bl); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &bl, nil
}
