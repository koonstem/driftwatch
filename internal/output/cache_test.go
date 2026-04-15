package output

import (
	"testing"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

func makeCacheResults(drifted bool) []drift.DriftResult {
	status := drift.StatusMatch
	if drifted {
		status = drift.StatusDrifted
	}
	return []drift.DriftResult{
		{Service: "svc-a", Status: status},
	}
}

func TestCacheWriter_Disabled_NeverHits(t *testing.T) {
	c := NewCacheWriter(CacheOptions{Enabled: false, TTL: time.Minute})
	results := makeCacheResults(false)
	c.Store(results)
	_, ok := c.Lookup(results)
	if ok {
		t.Fatal("expected no cache hit when disabled")
	}
}

func TestCacheWriter_Enabled_HitsAfterStore(t *testing.T) {
	c := NewCacheWriter(CacheOptions{Enabled: true, TTL: time.Minute})
	results := makeCacheResults(true)
	c.Store(results)
	got, ok := c.Lookup(results)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestCacheWriter_Enabled_MissOnDifferentInput(t *testing.T) {
	c := NewCacheWriter(CacheOptions{Enabled: true, TTL: time.Minute})
	c.Store(makeCacheResults(false))
	_, ok := c.Lookup(makeCacheResults(true))
	if ok {
		t.Fatal("expected cache miss for different input")
	}
}

func TestCacheWriter_Enabled_ExpiredEntry(t *testing.T) {
	c := NewCacheWriter(CacheOptions{Enabled: true, TTL: time.Millisecond})
	results := makeCacheResults(false)
	c.Store(results)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Lookup(results)
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestCacheWriter_Invalidate_ClearsAll(t *testing.T) {
	c := NewCacheWriter(CacheOptions{Enabled: true, TTL: time.Minute})
	results := makeCacheResults(false)
	c.Store(results)
	c.Invalidate()
	_, ok := c.Lookup(results)
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestCacheOptionsFromFlags_Defaults(t *testing.T) {
	cmd := newSortCmd() // reuse any cobra cmd helper
	BindCacheFlags(cmd)
	opts, err := CacheOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Enabled {
		t.Error("expected cache disabled by default")
	}
	if opts.TTL != 30*time.Second {
		t.Errorf("expected 30s TTL, got %s", opts.TTL)
	}
}
