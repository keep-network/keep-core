package libp2p

import (
	"container/list"
	"sync"
	"time"
)

// TimeCache provides a time cache.
type TimeCache struct {
	indexer  *list.List
	cache    map[string]time.Time
	timespan time.Duration
	mutex    sync.RWMutex
}

// NewTimeCache creates a new cache instance with provided timespan.
func NewTimeCache(timespan time.Duration) *TimeCache {
	tc := &TimeCache{
		indexer:  list.New(),
		cache:    make(map[string]time.Time),
		timespan: timespan,
	}

	go func() {
		for {
			time.Sleep(timespan)
			tc.sweep()
		}
	}()

	return tc
}

// Add adds an entry to the cache. Returns `true` if entry was not present in
// the cache and was successfully added into it. Returns `false` if
// entry is already in the cache. This method is synchronized.
func (tc *TimeCache) Add(item string) bool {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	_, ok := tc.cache[item]
	if ok {
		return false
	}

	tc.cache[item] = time.Now()
	tc.indexer.PushFront(item)
	return true
}

func (tc *TimeCache) sweep() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	for {
		back := tc.indexer.Back()
		if back == nil {
			return
		}

		item := back.Value.(string)
		itemTime, ok := tc.cache[item]
		if !ok {
			logger.Fatalf(
				"inconsistent cache state - expected item %v is not present",
				item,
			)
			return
		}

		if time.Since(itemTime) > tc.timespan {
			tc.indexer.Remove(back)
			delete(tc.cache, item)
		} else {
			return
		}
	}
}

// Has checks presence of an entry in the cache. Returns `true` if entry is
// present and `false` otherwise.
func (tc *TimeCache) Has(item string) bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	_, ok := tc.cache[item]
	return ok
}
