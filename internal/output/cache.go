package output

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// CacheOptions controls result caching behaviour.
type CacheOptions struct {
	Enabled bool
	TTL     time.Duration
}

// DefaultCacheOptions returns sensible defaults.
func DefaultCacheOptions() CacheOptions {
	return CacheOptions{
		Enabled: false,
		TTL:     30 * time.Second,
	}
}

type cacheEntry struct {
	results []drift.DriftResult
	stored  time.Time
}

// CacheWriter wraps a render function and memoises results by a hash of
// the input slice for the configured TTL.
type CacheWriter struct {
	opts  CacheOptions
	mu    sync.Mutex
	cache map[string]cacheEntry
}

// NewCacheWriter creates a CacheWriter with the given options.
func NewCacheWriter(opts CacheOptions) *CacheWriter {
	return &CacheWriter{
		opts:  opts,
		cache: make(map[string]cacheEntry),
	}
}

// Lookup returns cached results and true if a valid entry exists for the
// given results slice, otherwise it returns nil and false.
func (c *CacheWriter) Lookup(results []drift.DriftResult) ([]drift.DriftResult, bool) {
	if !c.opts.Enabled {
		return nil, false
	}
	key := hashResults(results)
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.stored) > c.opts.TTL {
		delete(c.cache, key)
		return nil, false
	}
	return entry.results, true
}

// Store saves results under a hash key derived from the slice contents.
func (c *CacheWriter) Store(results []drift.DriftResult) {
	if !c.opts.Enabled {
		return
	}
	key := hashResults(results)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{results: results, stored: time.Now()}
}

// Invalidate removes all cached entries.
func (c *CacheWriter) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]cacheEntry)
}

func hashResults(results []drift.DriftResult) string {
	b, _ := json.Marshal(results)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
