package cache

import (
	"github.com/whyrusleeping/timecache"

	"sync"
	"time"
)

// SynchronizedTimeCache provides a time cache with synchronization abilities.
type SynchronizedTimeCache struct {
	cacheMutex sync.Mutex
	cache      *timecache.TimeCache
}

// NewSynchronizedTimeCache creates a new cache instance with provided time span.
func NewSynchronizedTimeCache(span time.Duration) *SynchronizedTimeCache {
	return &SynchronizedTimeCache{
		cache: timecache.NewTimeCache(span),
	}
}

// Add adds an entry to the cache. Returns `true` if entry was not present in
// the cache and was successfully added into it. Returns `false` if
// entry is already in the cache. This method is synchronized.
func (mc *SynchronizedTimeCache) Add(id string) bool {
	mc.cacheMutex.Lock()
	defer mc.cacheMutex.Unlock()
	if mc.cache.Has(id) {
		return false
	}

	mc.cache.Add(id)
	return true
}

// Has checks presence of an entry in the cache. Returns `true` if entry is
// present and `false` otherwise. This method is synchronized.
func (mc *SynchronizedTimeCache) Has(id string) bool {
	mc.cacheMutex.Lock()
	defer mc.cacheMutex.Unlock()
	return mc.cache.Has(id)
}
