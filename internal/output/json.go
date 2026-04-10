package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/driftwatch/internal/drift"
)

type jsonReport struct {
	Drifted  bool             `json:"drifted"`
	Services []jsonServiceDrift `json:"services,omitempty"`
}

type jsonServiceDrift struct {
	Service  string `json:"service"`
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

func writeJSON(w io.Writer, report *drift.Report) error {
	out := jsonReport{Drifted: report.HasDrift()}

	for _, entry := range report.Entries {
		if entry.Drifted {
			out.Services = append(out.Services, jsonServiceDrift{
				Service:  entry.ServiceName,
				Field:    entry.Field,
				Expected: entry.Expected,
				Actual:   entry.Actual,
			})
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("output: json encode: %w", err)
	}
	return nil
}
