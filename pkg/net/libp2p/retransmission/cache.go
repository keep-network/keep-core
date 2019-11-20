package retransmission

import (
	"container/list"
	"sync"
	"time"
)

// SynchronizedTimeCache provides a time cache with synchronization abilities.
type SynchronizedTimeCache struct {
	indexer  *list.List
	cache    map[string]time.Time
	timespan time.Duration
	mutex    sync.RWMutex
}

// NewSynchronizedTimeCache creates a new cache instance with provided timespan.
func NewSynchronizedTimeCache(timespan time.Duration) *SynchronizedTimeCache {
	stc := &SynchronizedTimeCache{
		indexer:  list.New(),
		cache:    make(map[string]time.Time),
		timespan: timespan,
	}

	go func() {
		time.Sleep(timespan)
		stc.sweep()
	}()

	return stc
}

// Add adds an entry to the cache. Returns `true` if entry was not present in
// the cache and was successfully added into it. Returns `false` if
// entry is already in the cache. This method is synchronized.
func (stc *SynchronizedTimeCache) Add(item string) bool {
	stc.mutex.Lock()
	defer stc.mutex.Unlock()

	_, ok := stc.cache[item]
	if ok {
		return false
	}

	stc.cache[item] = time.Now()
	stc.indexer.PushFront(item)
	return true
}

func (stc *SynchronizedTimeCache) sweep() {
	stc.mutex.Lock()
	defer stc.mutex.Unlock()

	for {
		back := stc.indexer.Back()
		if back == nil {
			return
		}

		item := back.Value.(string)
		itemTime, ok := stc.cache[item]
		if !ok {
			panic("inconsistent cache state")
		}

		if time.Since(itemTime) > stc.timespan {
			stc.indexer.Remove(back)
			delete(stc.cache, item)
		} else {
			return
		}
	}
}

// Has checks presence of an entry in the cache. Returns `true` if entry is
// present and `false` otherwise. This method is synchronized.
func (stc *SynchronizedTimeCache) Has(item string) bool {
	stc.mutex.RLock()
	defer stc.mutex.RUnlock()

	_, ok := stc.cache[item]
	return ok
}
