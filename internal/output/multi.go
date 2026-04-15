package output

import (
	"fmt"
	"io"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// MultiWriter fans out a single Write call to multiple writers in order.
// If any writer returns an error the remaining writers are still called and
// all errors are collected and returned as a combined error.
type MultiWriter struct {
	writers []func([]drift.Result, io.Writer) error
	out     io.Writer
}

// MultiOptions configures a MultiWriter.
type MultiOptions struct {
	// StopOnError causes the fan-out to halt at the first writer that errors.
	StopOnError bool
}

// DefaultMultiOptions returns safe defaults.
func DefaultMultiOptions() MultiOptions {
	return MultiOptions{
		StopOnError: false,
	}
}

// NewMultiWriter creates a MultiWriter that writes to out using each provided
// writer function.
func NewMultiWriter(opts MultiOptions, out io.Writer, writers ...func([]drift.Result, io.Writer) error) *MultiWriter {
	return &MultiWriter{
		writers: writers,
		out:     out,
	}
}

// Write calls every registered writer function with results and the shared
// output writer.  Errors are accumulated; the first error is returned when
// StopOnError is true.
func (m *MultiWriter) Write(results []drift.Result) error {
	var errs []error
	for _, w := range m.writers {
		if err := w(results, m.out); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return combineErrors(errs)
}

func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	}
	msg := "multiple writer errors:"
	for _, e := range errs {
		msg += fmt.Sprintf(" [%s]", e.Error())
	}
	return fmt.Errorf("%s", msg)
}
