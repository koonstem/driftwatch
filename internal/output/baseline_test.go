package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{Service: "api", ExpectedImage: "api:1.0", Drifted: false},
		{Service: "worker", ExpectedImage: "worker:2.0", Drifted: true, Reasons: []string{"image mismatch"}},
	}
}

func TestBaselineWriter_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "baseline.json")
	w := NewBaselineWriter(path)

	if err := w.Write(makeResults()); err != nil {
		t.Fatalf("Write: %v", err)
	}

	bl, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if len(bl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(bl.Entries))
	}
	if bl.Entries[1].Service != "worker" || !bl.Entries[1].Drifted {
		t.Errorf("unexpected entry: %+v", bl.Entries[1])
	}
}

func TestLoadBaseline_FileNotFound(t *testing.T) {
	_, err := LoadBaseline(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	_, err := LoadBaseline(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompareToBaseline_NewService(t *testing.T) {
	bl := &Baseline{Entries: []BaselineEntry{}}
	diffs := CompareToBaseline(bl, makeResults())
	for _, d := range diffs {
		if !d.New {
			t.Errorf("expected service %q to be New", d.Service)
		}
	}
}

func TestCompareToBaseline_GoneService(t *testing.T) {
	bl := &Baseline{
		Entries: []BaselineEntry{
			{Service: "old-svc", Image: "old:1", Drifted: false},
		},
	}
	diffs := CompareToBaseline(bl, []drift.Result{})
	if len(diffs) != 1 || !diffs[0].Gone {
		t.Errorf("expected gone diff, got %+v", diffs)
	}
}

func TestCompareToBaseline_DriftChanged(t *testing.T) {
	bl := &Baseline{
		Entries: []BaselineEntry{
			{Service: "api", Drifted: false},
		},
	}
	current := []drift.Result{
		{Service: "api", Drifted: true},
	}
	diffs := CompareToBaseline(bl, current)
	if len(diffs) != 1 || !diffs[0].Changed() {
		t.Errorf("expected changed diff")
	}
}
