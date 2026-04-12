package output

import (
	"fmt"
	"io"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// WatchWriter periodically re-renders drift results to a writer,
// clearing the screen between renders to simulate a live view.
type WatchWriter struct {
	out      io.Writer
	interval time.Duration
	render   func([]drift.Result) error
}

// NewWatchWriter creates a WatchWriter that calls renderFn on each tick.
func NewWatchWriter(out io.Writer, interval time.Duration, renderFn func([]drift.Result) error) *WatchWriter {
	return &WatchWriter{
		out:      out,
		interval: interval,
		render:   renderFn,
	}
}

// Run starts the watch loop. It calls fetchFn to obtain fresh results on
// each tick and writes them via the render function. The loop exits when
// stopCh is closed.
func (w *WatchWriter) Run(stopCh <-chan struct{}, fetchFn func() ([]drift.Result, error)) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		results, err := fetchFn()
		if err != nil {
			return fmt.Errorf("watch fetch: %w", err)
		}
		w.clearScreen()
		if err := w.render(results); err != nil {
			return fmt.Errorf("watch render: %w", err)
		}

		select {
		case <-stopCh:
			return nil
		case <-ticker.C:
		}
	}
}

func (w *WatchWriter) clearScreen() {
	// ANSI escape: move cursor to top-left and clear screen
	fmt.Fprint(w.out, "\033[H\033[2J")
}
