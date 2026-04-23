package output

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

func TestReplay_FullPipeline_OrderPreserved(t *testing.T) {
	entries := []HistoryEntry{
		{Timestamp: time.Now().Add(-2 * time.Minute), Results: []drift.DriftResult{{Service: "alpha", Status: "ok"}}},
		{Timestamp: time.Now().Add(-1 * time.Minute), Results: []drift.DriftResult{{Service: "beta", Status: "drifted"}}},
		{Timestamp: time.Now(), Results: []drift.DriftResult{{Service: "gamma", Status: "ok"}}},
	}

	data, _ := json.Marshal(entries)
	f, _ := os.CreateTemp(t.TempDir(), "replay-int-*.json")
	_ = os.WriteFile(f.Name(), data, 0o644)

	var order []string
	w := NewReplayWriter(ReplayOptions{Enabled: true, File: f.Name()}, func(r []drift.DriftResult) error {
		for _, res := range r {
			order = append(order, res.Service)
		}
		return nil
	})

	if err := w.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alpha", "beta", "gamma"}
	for i, svc := range want {
		if order[i] != svc {
			t.Errorf("position %d: got %q, want %q", i, order[i], svc)
		}
	}
}

func TestReplay_FullPipeline_EmptyHistory_NoError(t *testing.T) {
	data, _ := json.Marshal([]HistoryEntry{})
	f, _ := os.CreateTemp(t.TempDir(), "replay-empty-*.json")
	_ = os.WriteFile(f.Name(), data, 0o644)

	w := NewReplayWriter(ReplayOptions{Enabled: true, File: f.Name()}, func(_ []drift.DriftResult) error {
		return nil
	})
	if err := w.Run(); err != nil {
		t.Fatalf("unexpected error on empty history: %v", err)
	}
}
