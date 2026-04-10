package drift

import (
	"fmt"
	"io"
	"strings"
)

// Reporter formats drift results for human-readable output.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a summary of all drift results to the reporter's writer.
func (r *Reporter) Print(results []Result) {
	driftCount := 0
	for _, res := range results {
		if res.Drifted {
			driftCount++
		}
	}

	fmt.Fprintf(r.w, "Drift report — %d service(s) checked, %d drifted\n", len(results), driftCount)
	fmt.Fprintln(r.w, strings.Repeat("-", 50))

	for _, res := range results {
		if res.Drifted {
			fmt.Fprintf(r.w, "[DRIFT]  %s\n", res.ServiceName)
			for _, reason := range res.Reasons {
				fmt.Fprintf(r.w, "         • %s\n", reason)
			}
		} else {
			fmt.Fprintf(r.w, "[OK]     %s\n", res.ServiceName)
		}
	}
}

// HasDrift returns true if any result indicates drift.
func HasDrift(results []Result) bool {
	for _, r := range results {
		if r.Drifted {
			return true
		}
	}
	return false
}
