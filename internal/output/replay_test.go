package output

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeReplayHistory(t *testing.T, entries []HistoryEntry) string {
	t.Helper()
	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal history: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "replay-*.json")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	_ = f.Close()
	if err := os.WriteFile(f.Name(), data, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return f.Name()
}

func TestReplayWriter_Disabled_DoesNothing(t *testing.T) {
	called := false
	w := NewReplayWriter(DefaultReplayOptions(), func(_ []drift.DriftResult) error {
		called = true
		return nil
	})
	if err := w.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("downstream should not be called when disabled")
	}
}

func TestReplayWriter_Enabled_NoFile_ReturnsError(t *testing.T) {
	w := NewReplayWriter(ReplayOptions{Enabled: true, File: ""}, func(_ []drift.DriftResult) error { return nil })
	if err := w.Run(); err == nil {
		t.Error("expected error when file is empty")
	}
}

func TestReplayWriter_Enabled_ReplayesAllSnapshots(t *testing.T) {
	results1 := []drift.DriftResult{{Service: "web", Status: "ok"}}
	results2 := []drift.DriftResult{{Service: "db", Status: "drifted"}}
	entries := []HistoryEntry{
		{Timestamp: time.Now(), Results: results1},
		{Timestamp: time.Now(), Results: results2},
	}
	path := makeReplayHistory(t, entries)

	var received [][]drift.DriftResult
	w := NewReplayWriter(ReplayOptions{Enabled: true, File: path}, func(r []drift.DriftResult) error {
		received = append(received, r)
		return nil
	})
	if err := w.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(received))
	}
	if received[0][0].Service != "web" {
		t.Errorf("first snapshot service: got %q, want %q", received[0][0].Service, "web")
	}
}

func TestLoadReplayFile_FileNotFound(t *testing.T) {
	_, err := LoadReplayFile("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadReplayFile_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	_ = os.WriteFile(f.Name(), []byte("not json"), 0o644)
	_, err := LoadReplayFile(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
