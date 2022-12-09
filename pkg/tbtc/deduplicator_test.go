package tbtc

import (
	"encoding/hex"
	"github.com/keep-network/keep-common/pkg/cache"
	"math/big"
	"testing"
	"time"
)

const testDKGSeedCachePeriod = 1 * time.Second
const testDKGResultHashCachePeriod = 1 * time.Second

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

func TestNotifyDKGResultSubmitted(t *testing.T) {
	deduplicator := deduplicator{
		dkgResultHashCache: cache.NewTimeCache(testDKGResultHashCachePeriod),
	}

	hash1Bytes, err := hex.DecodeString("92327ddff69a2b8c7ae787c5d590a2f14586089e6339e942d56e82aa42052cd9")
	if err != nil {
		t.Fatal(err)
	}
	var hash1 [32]byte
	copy(hash1[:], hash1Bytes)

	hash2Bytes, err := hex.DecodeString("23c0062913c4614bdff07f94475ceb4c585df53f71611776c3521ed8f8785913")
	if err != nil {
		t.Fatal(err)
	}
	var hash2 [32]byte
	copy(hash2[:], hash2Bytes)

	// Add the first hash.
	canProcess := deduplicator.notifyDKGResultSubmitted(hash1)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the second hash.
	canProcess = deduplicator.notifyDKGResultSubmitted(hash2)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the first hash before caching period elapses.
	canProcess = deduplicator.notifyDKGResultSubmitted(hash1)
	if canProcess {
		t.Fatal("should not be allowed to process")
	}

	// Wait until caching period elapses.
	time.Sleep(testDKGResultHashCachePeriod)

	// Add the first hash again.
	canProcess = deduplicator.notifyDKGResultSubmitted(hash1)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}
}
