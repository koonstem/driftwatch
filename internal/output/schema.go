package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// SchemaOptions controls schema export behaviour.
type SchemaOptions struct {
	Pretty  bool
	Version string
}

// DefaultSchemaOptions returns sensible defaults.
func DefaultSchemaOptions() SchemaOptions {
	return SchemaOptions{
		Pretty:  true,
		Version: "v1",
	}
}

// schemaField describes a single drift field observed across results.
type schemaField struct {
	Name     string `json:"name"`
	Observed int    `json:"observed_count"`
}

// schemaDoc is the top-level schema document written to the writer.
type schemaDoc struct {
	Version  string        `json:"version"`
	Services []string      `json:"services"`
	Fields   []schemaField `json:"drift_fields"`
}

// NewSchemaWriter writes a schema summary of observed drift fields to w.
func NewSchemaWriter(w io.Writer, opts SchemaOptions) func([]drift.Result) error {
	return func(results []drift.Result) error {
		serviceSet := map[string]struct{}{}
		fieldCount := map[string]int{}

		for _, r := range results {
			serviceSet[r.Service] = struct{}{}
			for _, d := range r.Diffs {
				fieldCount[d.Field]++
			}
		}

		services := make([]string, 0, len(serviceSet))
		for s := range serviceSet {
			services = append(services, s)
		}
		sort.Strings(services)

		fields := make([]schemaField, 0, len(fieldCount))
		for name, count := range fieldCount {
			fields = append(fields, schemaField{Name: name, Observed: count})
		}
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		doc := schemaDoc{
			Version:  opts.Version,
			Services: services,
			Fields:   fields,
		}

		var data []byte
		var err error
		if opts.Pretty {
			data, err = json.MarshalIndent(doc, "", "  ")
		} else {
			data, err = json.Marshal(doc)
		}
		if err != nil {
			return fmt.Errorf("schema: marshal: %w", err)
		}

		_, err = fmt.Fprintln(w, string(data))
		return err
	}
}
