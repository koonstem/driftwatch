package output

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeMultiResults() []drift.Result {
	return []drift.Result{
		{Service: "svc-a", Drifted: false, Fields: nil},
		{Service: "svc-b", Drifted: true, Fields: []drift.FieldDiff{{Field: "image", Declared: "nginx:1.24", Running: "nginx:1.25"}}},
	}
}

func TestMultiWriter_AllSucceed(t *testing.T) {
	var buf bytes.Buffer
	called := 0
	w := func(_ []drift.Result, _ io.Writer) error { called++; return nil }
	mw := NewMultiWriter(DefaultMultiOptions(), &buf, w, w, w)
	if err := mw.Write(makeMultiResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 3 {
		t.Fatalf("expected 3 calls, got %d", called)
	}
}

func TestMultiWriter_OneError_OthersStillCalled(t *testing.T) {
	var buf bytes.Buffer
	called := 0
	ok := func(_ []drift.Result, _ io.Writer) error { called++; return nil }
	bad := func(_ []drift.Result, _ io.Writer) error { called++; return errors.New("boom") }
	mw := NewMultiWriter(DefaultMultiOptions(), &buf, ok, bad, ok)
	if err := mw.Write(makeMultiResults()); err == nil {
		t.Fatal("expected error")
	}
	if called != 3 {
		t.Fatalf("expected all 3 writers called, got %d", called)
	}
}

func TestMultiWriter_MultipleErrors_CombinesMessages(t *testing.T) {
	var buf bytes.Buffer
	bad := func(_ []drift.Result, _ io.Writer) error { return errors.New("fail") }
	mw := NewMultiWriter(DefaultMultiOptions(), &buf, bad, bad)
	err := mw.Write(makeMultiResults())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !strings.Contains(err.Error(), "multiple writer errors") {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestMultiWriter_NoWriters_ReturnsNil(t *testing.T) {
	var buf bytes.Buffer
	mw := NewMultiWriter(DefaultMultiOptions(), &buf)
	if err := mw.Write(makeMultiResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultiWriter_EmptyResults_NoError(t *testing.T) {
	var buf bytes.Buffer
	called := false
	w := func(r []drift.Result, _ io.Writer) error {
		called = true
		if len(r) != 0 {
			return errors.New("expected empty")
		}
		return nil
	}
	mw := NewMultiWriter(DefaultMultiOptions(), &buf, w)
	if err := mw.Write([]drift.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("writer was not called")
	}
}
