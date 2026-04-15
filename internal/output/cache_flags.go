package output

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// BindCacheFlags registers cache-related flags on cmd.
func BindCacheFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("cache", false, "enable result caching")
	cmd.Flags().Duration("cache-ttl", 30*time.Second, "how long to cache results")
}

// CacheOptionsFromFlags reads cache flags from cmd and returns CacheOptions.
func CacheOptionsFromFlags(cmd *cobra.Command) (CacheOptions, error) {
	enabled, err := cmd.Flags().GetBool("cache")
	if err != nil {
		return CacheOptions{}, fmt.Errorf("cache flag: %w", err)
	}
	ttl, err := cmd.Flags().GetDuration("cache-ttl")
	if err != nil {
		return CacheOptions{}, fmt.Errorf("cache-ttl flag: %w", err)
	}
	if ttl <= 0 {
		return CacheOptions{}, fmt.Errorf("cache-ttl must be positive, got %s", ttl)
	}
	return CacheOptions{Enabled: enabled, TTL: ttl}, nil
}
