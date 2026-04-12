package output

import (
	"fmt"
	"io"
	"strings"
)

// ProgressWriter writes step-by-step progress messages during drift detection.
type ProgressWriter struct {
	w       io.Writer
	verbose bool
	color   *Colorizer
}

// NewProgressWriter creates a new ProgressWriter.
func NewProgressWriter(w io.Writer, verbose bool, color *Colorizer) *ProgressWriter {
	return &ProgressWriter{w: w, verbose: verbose, color: color}
}

// StepStart prints the beginning of a named step.
func (p *ProgressWriter) StepStart(name string) {
	if !p.verbose {
		return
	}
	label := p.color.Apply("→ "+name+"...", ColorCyan)
	fmt.Fprintln(p.w, label)
}

// StepDone prints a completion message for a named step.
func (p *ProgressWriter) StepDone(name string) {
	if !p.verbose {
		return
	}
	label := p.color.Apply("✓ "+name, ColorGreen)
	fmt.Fprintln(p.w, label)
}

// StepWarn prints a warning message for a named step.
func (p *ProgressWriter) StepWarn(name, reason string) {
	if !p.verbose {
		return
	}
	parts := []string{"⚠ ", name}
	if reason != "" {
		parts = append(parts, ": "+reason)
	}
	label := p.color.Apply(strings.Join(parts, ""), ColorYellow)
	fmt.Fprintln(p.w, label)
}

// Summary prints a one-line summary regardless of verbosity.
func (p *ProgressWriter) Summary(total, drifted int) {
	line := fmt.Sprintf("Checked %d service(s), %d drifted.", total, drifted)
	if drifted > 0 {
		line = p.color.Apply(line, ColorRed)
	} else {
		line = p.color.Apply(line, ColorGreen)
	}
	fmt.Fprintln(p.w, line)
}
