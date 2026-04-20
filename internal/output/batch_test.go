package output_test

import (
	"fmt"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/output"
)

func makeBatchResults(n int) drift.Report {
	results := make([]drift.DriftResult, n)
	for i := range results {
		results[i] = drift.DriftResult{
			Service: fmt.Sprintf("svc-%d", i),
			Drifted: i%2 == 0,
		}
	}
	return drift.Report{Results: results}
}

type captureWriter struct {
	batches []drift.Report
}

func (c *captureWriter) Write(r drift.Report) error {
	c.batches = append(c.batches, r)
	return nil
}

func TestBatchWriter_Disabled_ForwardsImmediately(t *testing.T) {
	cap := &captureWriter{}
	bw, err := output.NewBatchWriter(output.BatchOptions{Enabled: false, BatchSize: 3}, cap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	report := makeBatchResults(5)
	if err := bw.Write(report); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if len(cap.batches) != 1 {
		t.Errorf("expected 1 batch, got %d", len(cap.batches))
	}
	if len(cap.batches[0].Results) != 5 {
		t.Errorf("expected 5 results, got %d", len(cap.batches[0].Results))
	}
}

func TestBatchWriter_Enabled_FlushesOnBatchSize(t *testing.T) {
	cap := &captureWriter{}
	bw, err := output.NewBatchWriter(output.BatchOptions{Enabled: true, BatchSize: 3}, cap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := bw.Write(makeBatchResults(6)); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if len(cap.batches) != 2 {
		t.Errorf("expected 2 batches, got %d", len(cap.batches))
	}
}

func TestBatchWriter_Flush_SendsRemainder(t *testing.T) {
	cap := &captureWriter{}
	bw, err := output.NewBatchWriter(output.BatchOptions{Enabled: true, BatchSize: 4}, cap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = bw.Write(makeBatchResults(5))
	if err := bw.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	total := 0
	for _, b := range cap.batches {
		total += len(b.Results)
	}
	if total != 5 {
		t.Errorf("expected 5 total results after flush, got %d", total)
	}
}

func TestBatchWriter_InvalidBatchSize_ReturnsError(t *testing.T) {
	cap := &captureWriter{}
	_, err := output.NewBatchWriter(output.BatchOptions{Enabled: true, BatchSize: 0}, cap)
	if err == nil {
		t.Error("expected error for batch size 0, got nil")
	}
}

func TestBatchWriter_NilNext_ReturnsError(t *testing.T) {
	_, err := output.NewBatchWriter(output.DefaultBatchOptions(), nil)
	if err == nil {
		t.Error("expected error for nil next writer")
	}
}
