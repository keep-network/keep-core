package libp2p

import (
	"container/list"
	"sync"
	"time"
)

// timeCache provides a time cache safe for concurrent use by
// multiple goroutines without additional locking or coordination.
type timeCache struct {
	indexer  *list.List
	cache    map[string]time.Time
	timespan time.Duration
	mutex    sync.RWMutex
}

// NewTimeCache creates a new cache instance with provided timespan.
func newTimeCache(timespan time.Duration) *timeCache {
	return &timeCache{
		indexer:  list.New(),
		cache:    make(map[string]time.Time),
		timespan: timespan,
	}
}

// Add adds an entry to the cache. Returns `true` if entry was not present in
// the cache and was successfully added into it. Returns `false` if
// entry is already in the cache. This method is synchronized.
func (tc *timeCache) add(item string) bool {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	_, ok := tc.cache[item]
	if ok {
		return false
	}

	// sweep old entries
	for {
		back := tc.indexer.Back()
		if back == nil {
			break
		}

		item := back.Value.(string)
		itemTime, ok := tc.cache[item]
		if !ok {
			logger.Fatalf(
				"inconsistent cache state - expected item [%v] is not present",
				item,
			)
			break
		}

		if time.Since(itemTime) > tc.timespan {
			tc.indexer.Remove(back)
			delete(tc.cache, item)
		} else {
			break
		}
	}

	tc.cache[item] = time.Now()
	tc.indexer.PushFront(item)
	return true
}

// Has checks presence of an entry in the cache. Returns `true` if entry is
// present and `false` otherwise.
func (tc *timeCache) has(item string) bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	_, ok := tc.cache[item]
	return ok
}
