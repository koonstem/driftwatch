package output

import (
	"fmt"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// DefaultMaskOptions returns safe defaults for the masker.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Enabled:    false,
		MaskChar:   "*",
		MaskLength: 6,
		Fields:     []string{"image"},
	}
}

// MaskOptions controls how field values are masked in drift results.
type MaskOptions struct {
	Enabled    bool
	MaskChar   string
	MaskLength int
	Fields     []string
}

// MaskWriter masks specified field values before forwarding to the next writer.
type MaskWriter struct {
	opts MaskOptions
	next func([]drift.DriftResult) error
}

// NewMaskWriter creates a MaskWriter that masks configured fields.
func NewMaskWriter(opts MaskOptions, next func([]drift.DriftResult) error) (*MaskWriter, error) {
	if next == nil {
		return nil, fmt.Errorf("mask: next writer must not be nil")
	}
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}
	if opts.MaskLength <= 0 {
		opts.MaskLength = 6
	}
	return &MaskWriter{opts: opts, next: next}, nil
}

// Write masks field values in results if enabled, then forwards to next.
func (m *MaskWriter) Write(results []drift.DriftResult) error {
	if !m.opts.Enabled {
		return m.next(results)
	}
	masked := make([]drift.DriftResult, len(results))
	for i, r := range results {
		masked[i] = r
		masked[i].Fields = maskFields(r.Fields, m.opts.Fields, m.opts.MaskChar, m.opts.MaskLength)
	}
	return m.next(masked)
}

func maskFields(fields []drift.FieldDiff, targets []string, char string, length int) []drift.FieldDiff {
	targetSet := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		targetSet[strings.ToLower(t)] = struct{}{}
	}
	out := make([]drift.FieldDiff, len(fields))
	for i, f := range fields {
		if _, ok := targetSet[strings.ToLower(f.Field)]; ok {
			f.Expected = strings.Repeat(char, length)
			f.Actual = strings.Repeat(char, length)
		}
		out[i] = f
	}
	return out
}
