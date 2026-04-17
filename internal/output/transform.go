package output

import (
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// TransformOptions controls how drift results are transformed before output.
type TransformOptions struct {
	UppercaseService bool
	TrimImageTag     bool
	RenameFields     map[string]string
}

// DefaultTransformOptions returns a no-op transform configuration.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		RenameFields: map[string]string{},
	}
}

// NewTransformer returns a Writer that applies field transformations before
// delegating to the wrapped writer.
func NewTransformer(opts TransformOptions, next Writer) Writer {
	return writerFunc(func(report drift.Report) error {
		transformed := applyTransforms(opts, report)
		return next.Write(transformed)
	})
}

func applyTransforms(opts TransformOptions, report drift.Report) drift.Report {
	out := drift.Report{
		Results: make([]drift.Result, len(report.Results)),
	}
	for i, r := range report.Results {
		if opts.UppercaseService {
			r.Service = strings.ToUpper(r.Service)
		}
		if opts.TrimImageTag {
			r.Fields = trimImageTags(r.Fields)
		}
		if len(opts.RenameFields) > 0 {
			r.Fields = renameFields(opts.RenameFields, r.Fields)
		}
		out.Results[i] = r
	}
	return out
}

func trimImageTags(fields []drift.Field) []drift.Field {
	out := make([]drift.Field, len(fields))
	for i, f := range fields {
		if f.Name == "image" {
			f.Actual = stripTag(f.Actual)
			f.Expected = stripTag(f.Expected)
		}
		out[i] = f
	}
	return out
}

func stripTag(image string) string {
	if idx := strings.LastIndex(image, ":"); idx != -1 {
		return image[:idx]
	}
	return image
}

func renameFields(mapping map[string]string, fields []drift.Field) []drift.Field {
	out := make([]drift.Field, len(fields))
	for i, f := range fields {
		if newName, ok := mapping[f.Name]; ok {
			f.Name = newName
		}
		out[i] = f
	}
	return out
}
