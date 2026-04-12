package output

import (
	"fmt"
	"io"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// SummaryWriter writes a human-readable summary of a drift detection run.
type SummaryWriter struct {
	out io.Writer
}

// NewSummaryWriter returns a new SummaryWriter that writes to w.
func NewSummaryWriter(w io.Writer) *SummaryWriter {
	return &SummaryWriter{out: w}
}

// Write outputs a summary block for the given report.
func (s *SummaryWriter) Write(report drift.Report) error {
	total := len(report.Results)
	drifted := 0
	for _, r := range report.Results {
		if r.Drifted {
			drifted++
		}
	}
	clean := total - drifted

	timestamp := time.Now().UTC().Format(time.RFC3339)

	fmt.Fprintf(s.out, "\n--- Drift Detection Summary [%s] ---\n", timestamp)
	fmt.Fprintf(s.out, "  Services checked : %d\n", total)
	fmt.Fprintf(s.out, "  Clean            : %d\n", clean)
	fmt.Fprintf(s.out, "  Drifted          : %d\n", drifted)

	if drifted > 0 {
		fmt.Fprintln(s.out, "  Status           : DRIFT DETECTED")
	} else {
		fmt.Fprintln(s.out, "  Status           : OK")
	}
	fmt.Fprintln(s.out, "------------------------------------")

	return nil
}
