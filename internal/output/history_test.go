package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeHistoryResults(drifted bool) []drift.Result {
	return []drift.Result{
		{
			Service: "web",
			Drifted: drifted,
			Fields:  []drift.Field{{Name: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}},
		},
	}
}

func TestHistoryWriter_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	w := NewHistoryWriter(path, 0)

	if err := w.Write(makeHistoryResults(true)); err != nil {
		t.Fatalf("Write: %v", err)
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Drifted != 1 {
		t.Errorf("expected Drifted=1, got %d", entries[0].Drifted)
	}
	if entries[0].Total != 1 {
		t.Errorf("expected Total=1, got %d", entries[0].Total)
	}
}

func TestHistoryWriter_Appends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	w := NewHistoryWriter(path, 0)

	_ = w.Write(makeHistoryResults(false))
	_ = w.Write(makeHistoryResults(true))

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHistoryWriter_MaxSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	w := NewHistoryWriter(path, 2)

	for i := 0; i < 5; i++ {
		_ = w.Write(makeHistoryResults(false))
	}

	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries after max trim, got %d", len(entries))
	}
}

func TestLoadHistory_FileNotFound(t *testing.T) {
	_, err := LoadHistory("/nonexistent/path/history.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadHistory_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadHistory(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSortHistoryByTime(t *testing.T) {
	now := time.Now()
	entries := []HistoryEntry{
		{Timestamp: now.Add(2 * time.Hour)},
		{Timestamp: now},
		{Timestamp: now.Add(time.Hour)},
	}
	SortHistoryByTime(entries)
	if !entries[0].Timestamp.Equal(now) {
		t.Errorf("expected earliest first")
	}
}

func TestHistoryEntry_JSONFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	w := NewHistoryWriter(path, 0)
	_ = w.Write(makeHistoryResults(true))

	data, _ := os.ReadFile(path)
	var raw []map[string]interface{}
	_ = json.Unmarshal(data, &raw)

	for _, key := range []string{"timestamp", "total", "drifted", "results"} {
		if _, ok := raw[0][key]; !ok {
			t.Errorf("missing JSON field: %s", key)
		}
	}
}
