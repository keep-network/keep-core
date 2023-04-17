package tbtc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
)

const (
	// DKGSeedCachePeriod is the time period the cache maintains
	// the DKG seed corresponding to a DKG instance.
	DKGSeedCachePeriod = 7 * 24 * time.Hour
	// DKGResultHashCachePeriod is the time period the cache maintains
	// the given DKG result hash.
	DKGResultHashCachePeriod = 7 * 24 * time.Hour
	// DepositSweepProposalCachePeriod is the time period the cache maintains
	// the given deposit sweep proposal.
	DepositSweepProposalCachePeriod = 7 * 24 * time.Hour
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
// - DKG result submitted
// - Deposit sweep proposal submission
type deduplicator struct {
	dkgSeedCache              *cache.TimeCache
	dkgResultHashCache        *cache.TimeCache
	depositSweepProposalCache *cache.TimeCache
}

func newDeduplicator() *deduplicator {
	return &deduplicator{
		dkgSeedCache:              cache.NewTimeCache(DKGSeedCachePeriod),
		dkgResultHashCache:        cache.NewTimeCache(DKGResultHashCachePeriod),
		depositSweepProposalCache: cache.NewTimeCache(DepositSweepProposalCachePeriod),
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

// notifyDKGResultSubmitted notifies the client wants to start some actions
// upon the DKG result submission. It returns boolean indicating whether the
// client should proceed with the actions or ignore the event as a duplicate.
func (d *deduplicator) notifyDKGResultSubmitted(
	newDKGResultSeed *big.Int,
	newDKGResultHash DKGChainResultHash,
	newDKGResultBlock uint64,
) bool {
	d.dkgResultHashCache.Sweep()

	cacheKey := newDKGResultSeed.Text(16) +
		hex.EncodeToString(newDKGResultHash[:]) +
		strconv.Itoa(int(newDKGResultBlock))

	// If the key is not in the cache, that means the result was not handled
	// yet and the client should proceed with the execution.
	if !d.dkgResultHashCache.Has(cacheKey) {
		d.dkgResultHashCache.Add(cacheKey)
		return true
	}

	// Otherwise, the DKG result is a duplicate and the client should not
	// proceed with the execution.
	return false
}

// notifyDepositSweepProposalSubmitted notifies the client wants to start some
// actions upon the deposit sweep proposal submission. It returns boolean
// indicating whether the client should proceed with the actions or ignore the
// event as a duplicate.
func (d *deduplicator) notifyDepositSweepProposalSubmitted(
	newProposal *DepositSweepProposal,
) bool {
	d.depositSweepProposalCache.Sweep()

	// We build the cache key by hashing the concatenation of relevant fields
	// of the proposal. It may be tempting to extract that code into a general
	// "hash code" function exposed by the DepositSweepProposal type but this
	// is not necessarily a good idea. The deduplicator is responsible for
	// detecting duplicates and construction of cache keys is part of that job.
	// Extracting this logic outside would push that responsibility out of the
	// deduplicator control. That is dangerous as deduplication logic could
	// be implicitly changeable from the outside and lead to serious bugs.
	var buffer bytes.Buffer
	buffer.Write(newProposal.WalletPubKeyHash[:])
	for _, depositKey := range newProposal.DepositsKeys {
		buffer.Write(depositKey.FundingTxHash[:])
		fundingOutputIndex := make([]byte, 4)
		binary.BigEndian.PutUint32(fundingOutputIndex, depositKey.FundingOutputIndex)
		buffer.Write(fundingOutputIndex)
	}
	buffer.Write(newProposal.SweepTxFee.Bytes())

	bufferSha256 := sha256.Sum256(buffer.Bytes())
	cacheKey := hex.EncodeToString(bufferSha256[:])

	// If the key is not in the cache, that means the proposal was not handled
	// yet and the client should proceed with the execution.
	if !d.depositSweepProposalCache.Has(cacheKey) {
		d.depositSweepProposalCache.Add(cacheKey)
		return true
	}

	// Otherwise, the proposal is a duplicate and the client should not
	// proceed with the execution.
	return false
}
