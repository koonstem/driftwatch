package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
)

func makePagerResults(n int, driftedIndexes ...int) []drift.DriftResult {
	drifted := make(map[int]bool, len(driftedIndexes))
	for _, i := range driftedIndexes {
		drifted[i] = true
	}
	results := make([]drift.DriftResult, n)
	for i := range results {
		results[i] = drift.DriftResult{
			ServiceName: fmt.Sprintf("service-%02d", i+1),
			Drifted:     drifted[i],
		}
		if drifted[i] {
			results[i].Fields = []drift.FieldDiff{
				{Field: "image", Expected: "v1", Actual: "v2"},
			}
		}
	}
	return results
}

func newTestPager(pageSize int, showInfo bool) (*PagerWriter, *bytes.Buffer) {
	var buf bytes.Buffer
	color := NewColorizer(false)
	opts := PagerOptions{PageSize: pageSize, ShowPageInfo: showInfo}
	return NewPagerWriter(&buf, opts, color), &buf
}

func TestPagerWriter_SinglePage(t *testing.T) {
	pw, buf := newTestPager(20, true)
	results := makePagerResults(5)

	written, hasMore := pw.WritePage(results, 0)

	if written != 5 {
		t.Errorf("expected 5 written, got %d", written)
	}
	if hasMore {
		t.Error("expected no more pages")
	}
	if !strings.Contains(buf.String(), "service-01") {
		t.Error("expected service-01 in output")
	}
}

func TestPagerWriter_MultiPage_HasMore(t *testing.T) {
	pw, _ := newTestPager(3, false)
	results := makePagerResults(7)

	written, hasMore := pw.WritePage(results, 0)

	if written != 3 {
		t.Errorf("expected 3 written, got %d", written)
	}
	if !hasMore {
		t.Error("expected more pages")
	}
}

func TestPagerWriter_SecondPage(t *testing.T) {
	pw, buf := newTestPager(3, false)
	results := makePagerResults(7)

	pw.WritePage(results, 1)

	if !strings.Contains(buf.String(), "service-04") {
		t.Error("expected service-04 on page 2")
	}
}

func TestPagerWriter_PageBeyondEnd(t *testing.T) {
	pw, buf := newTestPager(5, true)
	results := makePagerResults(3)

	written, hasMore := pw.WritePage(results, 5)

	if written != 0 {
		t.Errorf("expected 0 written, got %d", written)
	}
	if hasMore {
		t.Error("expected no more pages")
	}
	if buf.Len() != 0 {
		t.Error("expected empty output for out-of-range page")
	}
}

func TestPagerWriter_ShowPageInfo(t *testing.T) {
	pw, buf := newTestPager(2, true)
	results := makePagerResults(5)

	pw.WritePage(results, 0)

	out := buf.String()
	if !strings.Contains(out, "Page 1") {
		t.Error("expected page info header")
	}
	if !strings.Contains(out, "more results") {
		t.Error("expected 'more results' footer")
	}
}
