package tbtc

import (
	"github.com/keep-network/keep-common/pkg/cache"
	"math/big"
	"testing"
	"time"
)

const testDKGSeedCachePeriod = 1 * time.Second

func TestNotifyDKGStarted(t *testing.T) {
	deduplicator := deduplicator{
		dkgSeedCache: cache.NewTimeCache(testDKGSeedCachePeriod),
	}

	seed1 := big.NewInt(100)
	seed2 := big.NewInt(200)

	// Add the first seed.
	canJoinDKG := deduplicator.notifyDKGStarted(seed1)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}

	// Add the second seed.
	canJoinDKG = deduplicator.notifyDKGStarted(seed2)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}

	// Add the first seed before caching period elapses.
	canJoinDKG = deduplicator.notifyDKGStarted(seed1)
	if canJoinDKG {
		t.Fatal("should not be allowed to join DKG")
	}

	// Wait until caching period elapses.
	time.Sleep(testDKGSeedCachePeriod)

	// Add the first seed again.
	canJoinDKG = deduplicator.notifyDKGStarted(seed1)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}
}
