package tbtc

import (
	"github.com/keep-network/keep-common/pkg/cache"
	"math/big"
	"time"
)

const (
	// DKGSeedCachePeriod is the time period the cache maintains
	// the DKG seed corresponding to a DKG instance.
	DKGSeedCachePeriod = 7 * 24 * time.Hour
)

// deduplicator decides whether the given event should be handled by the
// client or not.
//
// Event subscription may emit the same event two or more times. The same event
// can be emitted right after it's been emitted for the first time. The same
// event can also be emitted a long time after it's been emitted for the first
// time. It is deduplicator's responsibility to decide whether the given
// event is a duplicate and should be ignored or if it is not a duplicate and
// should be handled.
//
// Those events are supported:
// - DKG started
type deduplicator struct {
	dkgSeedCache *cache.TimeCache
}

func newDeduplicator() *deduplicator {
	return &deduplicator{
		dkgSeedCache: cache.NewTimeCache(DKGSeedCachePeriod),
	}
}

// notifyDKGStarted notifies the client wants to start the distributed key
// generation upon receiving an event. It returns boolean indicating whether the
// client should proceed with the execution or ignore the event as a duplicate.
func (d *deduplicator) notifyDKGStarted(
	newDKGSeed *big.Int,
) bool {
	d.dkgSeedCache.Sweep()

	// The cache key is the hexadecimal representation of the seed.
	cacheKey := newDKGSeed.Text(16)
	// If the key is not in the cache, that means the seed was not handled
	// yet and the client should proceed with the execution.
	if !d.dkgSeedCache.Has(cacheKey) {
		d.dkgSeedCache.Add(cacheKey)
		return true
	}

	// Otherwise, the DKG seed is a duplicate and the client should not proceed
	// with the execution.
	return false
}
