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

// newTimeCache creates a new cache instance with provided timespan.
func newTimeCache(timespan time.Duration) *timeCache {
	return &timeCache{
		indexer:  list.New(),
		cache:    make(map[string]time.Time),
		timespan: timespan,
	}
}

// add adds an entry to the cache. If entry already exists, it resets its
// timestamp. Before new entry is added, all outdated entries are removed from
// the cache.
func (tc *timeCache) add(item string) {
	tc.sweep()

	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	if _, isPresent := tc.cache[item]; isPresent {
		// if item is already present, move its index to the front
		for index := tc.indexer.Front(); index != nil; index = index.Next() {
			if index.Value.(string) == item {
				tc.indexer.MoveToFront(index)
				break
			}
		}
	} else {
		tc.indexer.PushFront(item)
	}

	tc.cache[item] = time.Now()
}

// sweep removes all outdated entries from the cache.
func (tc *timeCache) sweep() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

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
}

// has checks presence of an entry in the cache. Returns `true` if entry is
// present and `false` otherwise.
func (tc *timeCache) has(item string) bool {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	_, ok := tc.cache[item]
	return ok
}
