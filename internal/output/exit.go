package output

import (
	"os"

	"github.com/yourusername/driftwatch/internal/drift"
)

// ExitCoder determines the appropriate OS exit code based on drift results.
type ExitCoder struct {
	failOnDrift bool
}

// NewExitCoder creates an ExitCoder. When failOnDrift is true, the process
// will exit with code 1 if any drift is detected.
func NewExitCoder(failOnDrift bool) *ExitCoder {
	return &ExitCoder{failOnDrift: failOnDrift}
}

// Code returns the exit code for the given report.
// Returns 0 if no drift or failOnDrift is disabled, 1 if drift was found.
func (e *ExitCoder) Code(report *drift.Report) int {
	if e.failOnDrift && drift.HasDrift(report) {
		return 1
	}
	return 0
}

// Exit calls os.Exit with the appropriate code for the given report.
func (e *ExitCoder) Exit(report *drift.Report) {
	os.Exit(e.Code(report))
}
