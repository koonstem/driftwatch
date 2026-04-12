package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// PagerOptions controls paged output behaviour.
type PagerOptions struct {
	PageSize int
	ShowPageInfo bool
}

// DefaultPagerOptions returns sensible defaults.
func DefaultPagerOptions() PagerOptions {
	return PagerOptions{
		PageSize:     20,
		ShowPageInfo: true,
	}
}

// PagerWriter writes drift results in pages.
type PagerWriter struct {
	w    io.Writer
	opts PagerOptions
	color *Colorizer
}

// NewPagerWriter creates a PagerWriter that writes to w.
func NewPagerWriter(w io.Writer, opts PagerOptions, color *Colorizer) *PagerWriter {
	return &PagerWriter{w: w, opts: opts, color: color}
}

// WritePage writes a single page of results starting at offset.
// It returns the number of items written and whether more pages remain.
func (p *PagerWriter) WritePage(results []drift.DriftResult, page int) (written int, hasMore bool) {
	if p.opts.PageSize <= 0 {
		p.opts.PageSize = 20
	}

	start := page * p.opts.PageSize
	if start >= len(results) {
		return 0, false
	}

	end := start + p.opts.PageSize
	if end > len(results) {
		end = len(results)
	}

	slice := results[start:end]
	total := len(results)
	hasMore = end < total

	if p.opts.ShowPageInfo {
		header := fmt.Sprintf("Page %d — showing %d–%d of %d results\n",
			page+1, start+1, end, total)
		fmt.Fprint(p.w, p.color.Bold(header))
		fmt.Fprintln(p.w, strings.Repeat("─", 48))
	}

	for _, r := range slice {
		status := "ok"
		if r.Drifted {
			status = p.color.Red("DRIFTED")
		} else {
			status = p.color.Green("ok")
		}
		fmt.Fprintf(p.w, "  %-30s %s\n", r.ServiceName, status)
		for _, f := range r.Fields {
			fmt.Fprintf(p.w, "      %-20s expected=%s actual=%s\n",
				f.Field, f.Expected, f.Actual)
		}
		written++
	}

	if p.opts.ShowPageInfo && hasMore {
		fmt.Fprintln(p.w, strings.Repeat("─", 48))
		fmt.Fprintf(p.w, "  %s\n", p.color.Bold("(more results on next page)"))
	}

	return written, hasMore
}
