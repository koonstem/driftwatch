package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

const (
	colService = "SERVICE"
	colField   = "FIELD"
	colWant    = "EXPECTED"
	colGot     = "ACTUAL"
)

func writeTable(w io.Writer, report *drift.Report) error {
	if len(report.Results) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}

	rows := collectRows(report)
	widths := columnWidths(rows)

	writeTableRow(w, widths, colService, colField, colWant, colGot)
	writeTableSep(w, widths)

	for _, row := range rows {
		writeTableRow(w, widths, row[0], row[1], row[2], row[3])
	}

	return nil
}

func collectRows(report *drift.Report) [][]string {
	var rows [][]string
	for _, r := range report.Results {
		for _, d := range r.Diffs {
			rows = append(rows, []string{r.ServiceName, d.Field, d.Expected, d.Actual})
		}
	}
	return rows
}

func columnWidths(rows [][]string) [4]int {
	w := [4]int{len(colService), len(colField), len(colWant), len(colGot)}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > w[i] {
				w[i] = len(cell)
			}
		}
	}
	return w
}

func writeTableRow(w io.Writer, widths [4]int, a, b, c, d string) {
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s\n",
		widths[0], a,
		widths[1], b,
		widths[2], c,
		widths[3], d,
	)
}

func writeTableSep(w io.Writer, widths [4]int) {
	parts := make([]string, 4)
	for i, n := range widths {
		parts[i] = strings.Repeat("-", n)
	}
	fmt.Fprintf(w, "%s  %s  %s  %s\n", parts[0], parts[1], parts[2], parts[3])
}
