package output

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeAuditResults(drifted bool) []drift.Result {
	return []drift.Result{
		{
			Service: "web",
			Drifted: drifted,
			Fields:  []drift.FieldDiff{{Field: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}},
		},
	}
}

func TestAuditWriter_Disabled_NoFile(t *testing.T) {
	dir := t.TempDir()
	opts := DefaultAuditOptions()
	opts.Enabled = false
	opts.FilePath = filepath.Join(dir, "audit.jsonl")

	w := NewAuditWriter(opts)
	if err := w.Write(makeAuditResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(opts.FilePath); !os.IsNotExist(err) {
		t.Error("expected no audit file to be created when disabled")
	}
}

func TestAuditWriter_Enabled_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	opts := DefaultAuditOptions()
	opts.Enabled = true
	opts.FilePath = filepath.Join(dir, "audit.jsonl")
	opts.RunID = "test-run-001"

	w := NewAuditWriter(opts)
	if err := w.Write(makeAuditResults(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(opts.FilePath); os.IsNotExist(err) {
		t.Fatal("expected audit file to be created")
	}
}

func TestAuditWriter_Enabled_ValidJSONL(t *testing.T) {
	dir := t.TempDir()
	opts := DefaultAuditOptions()
	opts.Enabled = true
	opts.FilePath = filepath.Join(dir, "audit.jsonl")
	opts.RunID = "test-run-002"

	w := NewAuditWriter(opts)
	if err := w.Write(makeAuditResults(true)); err != nil {
		t.Fatalf("write error: %v", err)
	}

	f, _ := os.Open(opts.FilePath)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line in audit file")
	}

	var entry AuditEntry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON in audit entry: %v", err)
	}

	if entry.RunID != "test-run-002" {
		t.Errorf("expected run_id %q, got %q", "test-run-002", entry.RunID)
	}
	if entry.DriftedCount != 1 {
		t.Errorf("expected drifted_count 1, got %d", entry.DriftedCount)
	}
}

func TestAuditWriter_Appends(t *testing.T) {
	dir := t.TempDir()
	opts := DefaultAuditOptions()
	opts.Enabled = true
	opts.FilePath = filepath.Join(dir, "audit.jsonl")

	w := NewAuditWriter(opts)
	_ = w.Write(makeAuditResults(false))
	_ = w.Write(makeAuditResults(true))

	f, _ := os.Open(opts.FilePath)
	defer f.Close()

	lines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
	}
	if lines != 2 {
		t.Errorf("expected 2 audit lines, got %d", lines)
	}
}
