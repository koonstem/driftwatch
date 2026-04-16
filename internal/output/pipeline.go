package output

import "github.com/user/driftwatch/internal/drift"

// PipelineOptions controls which writers are chained together.
type PipelineOptions struct {
	StopOnError bool
	Writers     []Writer
}

// DefaultPipelineOptions returns safe defaults.
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		StopOnError: false,
	}
}

// Writer is the shared interface for all output writers.
type Writer interface {
	Write(results []drift.DriftResult) error
}

// NewPipeline builds an ordered pipeline of writers.
// Each writer receives the same results slice.
// If StopOnError is true, the first error halts execution.
func NewPipeline(opts PipelineOptions) *Pipeline {
	return &Pipeline{opts: opts}
}

// Pipeline executes an ordered chain of writers.
type Pipeline struct {
	opts PipelineOptions
}

// Run executes all writers in order.
func (p *Pipeline) Run(results []drift.DriftResult) error {
	var errs []error
	for _, w := range p.opts.Writers {
		if err := w.Write(results); err != nil {
			if p.opts.StopOnError {
				return err
			}
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return combineErrors(errs)
	}
	return nil
}
