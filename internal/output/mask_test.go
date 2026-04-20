package output

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
)

func makeMaskResults(field, expected, actual string) []drift.DriftResult {
	return []drift.DriftResult{
		{
			Service: "api",
			Drifted: true,
			Fields: []drift.FieldDiff{
				{Field: field, Expected: expected, Actual: actual},
			},
		},
	}
}

func TestMaskWriter_Disabled_PassesThrough(t *testing.T) {
	var got []drift.DriftResult
	opts := DefaultMaskOptions()
	opts.Enabled = false
	w, err := NewMaskWriter(opts, func(r []drift.DriftResult) error { got = r; return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeMaskResults("image", "nginx:1.19", "nginx:1.20")
	if err := w.Write(input); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if got[0].Fields[0].Expected != "nginx:1.19" {
		t.Errorf("expected unmasked value, got %q", got[0].Fields[0].Expected)
	}
}

func TestMaskWriter_Enabled_MasksTargetField(t *testing.T) {
	var got []drift.DriftResult
	opts := DefaultMaskOptions()
	opts.Enabled = true
	opts.Fields = []string{"image"}
	opts.MaskChar = "*"
	opts.MaskLength = 4
	w, _ := NewMaskWriter(opts, func(r []drift.DriftResult) error { got = r; return nil })
	input := makeMaskResults("image", "nginx:1.19", "nginx:1.20")
	if err := w.Write(input); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if got[0].Fields[0].Expected != "****" {
		t.Errorf("expected masked value, got %q", got[0].Fields[0].Expected)
	}
	if got[0].Fields[0].Actual != "****" {
		t.Errorf("expected masked actual, got %q", got[0].Fields[0].Actual)
	}
}

func TestMaskWriter_Enabled_DoesNotMaskOtherFields(t *testing.T) {
	var got []drift.DriftResult
	opts := DefaultMaskOptions()
	opts.Enabled = true
	opts.Fields = []string{"image"}
	w, _ := NewMaskWriter(opts, func(r []drift.DriftResult) error { got = r; return nil })
	input := makeMaskResults("replicas", "3", "2")
	if err := w.Write(input); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if got[0].Fields[0].Expected != "3" {
		t.Errorf("non-target field should not be masked, got %q", got[0].Fields[0].Expected)
	}
}

func TestMaskWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := NewMaskWriter(DefaultMaskOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil next writer")
	}
}

func TestMaskWriter_CaseInsensitiveFieldMatch(t *testing.T) {
	var got []drift.DriftResult
	opts := DefaultMaskOptions()
	opts.Enabled = true
	opts.Fields = []string{"IMAGE"}
	opts.MaskLength = 3
	w, _ := NewMaskWriter(opts, func(r []drift.DriftResult) error { got = r; return nil })
	input := makeMaskResults("image", "nginx:latest", "nginx:stable")
	if err := w.Write(input); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if got[0].Fields[0].Expected != "***" {
		t.Errorf("expected case-insensitive mask, got %q", got[0].Fields[0].Expected)
	}
}
