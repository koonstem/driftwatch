package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeSnapshotResults() []drift.DriftResult {
	return []drift.DriftResult{
		{Service: "api", Drifted: true, Fields: []drift.FieldDrift{{Field: "image", Expected: "api:1", Actual: "api:2"}}},
		{Service: "worker", Drifted: false},
	}
}

func TestSnapshotWriter_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := makeSnapshotResults()
	w := NewSnapshotWriter(path)
	if err := w.Write(results); err != nil {
		t.Fatalf("Write: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if len(snap.Results) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.Results[0].Service != "api" {
		t.Errorf("expected service api, got %s", snap.Results[0].Service)
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestLoadSnapshot_FileNotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSnapshotWriter_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	w := NewSnapshotWriter(path)
	_ = w.Write(makeSnapshotResults())

	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := raw["captured_at"]; !ok {
		t.Error("expected captured_at field in JSON")
	}
}
