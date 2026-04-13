package output

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// TemplateOptions holds configuration for the template writer.
type TemplateOptions struct {
	TemplatePath string
	TemplateStr  string
}

// DefaultTemplateOptions returns sensible defaults.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{}
}

// NewTemplateWriter renders drift results using a user-supplied Go template.
func NewTemplateWriter(w io.Writer, opts TemplateOptions) (*TemplateWriter, error) {
	var tmpl *template.Template
	var err error

	switch {
	case opts.TemplatePath != "":
		tmpl, err = template.ParseFiles(opts.TemplatePath)
		if err != nil {
			return nil, fmt.Errorf("parse template file: %w", err)
		}
	case opts.TemplateStr != "":
		tmpl, err = template.New("inline").Parse(opts.TemplateStr)
		if err != nil {
			return nil, fmt.Errorf("parse inline template: %w", err)
		}
	default:
		return nil, fmt.Errorf("template: no template source provided")
	}

	return &TemplateWriter{w: w, tmpl: tmpl}, nil
}

// TemplateWriter writes drift results using a Go template.
type TemplateWriter struct {
	w    io.Writer
	tmpl *template.Template
}

// TemplateData is the data passed to the template.
type TemplateData struct {
	Results  []drift.Result
	Total    int
	Drifted  int
	OK       int
}

// Write renders the template with the provided results.
func (t *TemplateWriter) Write(results []drift.Result) error {
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}

	data := TemplateData{
		Results: results,
		Total:   len(results),
		Drifted: drifted,
		OK:      len(results) - drifted,
	}

	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	_, err := fmt.Fprint(t.w, buf.String())
	return err
}
