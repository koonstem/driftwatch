package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeWatchResults(drifted bool) []drift.Result {
	status := drift.StatusOK
	if drifted {
		status = drift.StatusDrifted
	}
	return []drift.Result{
		{ServiceName: "api", Status: status},
	}
}

func TestWatchWriter_RunsRenderOnFirstTick(t *testing.T) {
	var buf bytes.Buffer
	renderCount := 0

	ww := NewWatchWriter(&buf, 100*time.Millisecond, func(results []drift.Result) error {
		renderCount++
		return nil
	})

	stopCh := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		done <- ww.Run(stopCh, func() ([]drift.Result, error) {
			return makeWatchResults(false), nil
		})
	}()

	time.Sleep(60 * time.Millisecond)
	close(stopCh)

	if err := <-done; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if renderCount < 1 {
		t.Errorf("expected at least 1 render call, got %d", renderCount)
	}
}

func TestWatchWriter_ClearsScreen(t *testing.T) {
	var buf bytes.Buffer
	ww := NewWatchWriter(&buf, 200*time.Millisecond, func(_ []drift.Result) error { return nil })

	stopCh := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- ww.Run(stopCh, func() ([]drift.Result, error) {
			return nil, nil
		})
	}()
	time.Sleep(50 * time.Millisecond)
	close(stopCh)
	<-done

	if !strings.Contains(buf.String(), "\033[H\033[2J") {
		t.Error("expected ANSI clear-screen escape in output")
	}
}

func TestWatchWriter_FetchError_StopsLoop(t *testing.T) {
	var buf bytes.Buffer
	ww := NewWatchWriter(&buf, 50*time.Millisecond, func(_ []drift.Result) error { return nil })

	stopCh := make(chan struct{})
	err := ww.Run(stopCh, func() ([]drift.Result, error) {
		return nil, errors.New("source unavailable")
	})

	if err == nil {
		t.Fatal("expected error from fetch failure, got nil")
	}
	if !strings.Contains(err.Error(), "source unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWatchWriter_RenderError_StopsLoop(t *testing.T) {
	var buf bytes.Buffer
	ww := NewWatchWriter(&buf, 50*time.Millisecond, func(_ []drift.Result) error {
		return errors.New("render failed")
	})

	stopCh := make(chan struct{})
	err := ww.Run(stopCh, func() ([]drift.Result, error) {
		return makeWatchResults(true), nil
	})

	if err == nil {
		t.Fatal("expected error from render failure, got nil")
	}
	if !strings.Contains(err.Error(), "render failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
