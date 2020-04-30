package firewall

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

// Disabled is an empty Firewall implementation enforcing no rules
// on the connection.
var Disabled = &noFirewall{}

type noFirewall struct{}

func (nf *noFirewall) Validate(remotePeerPublicKey *ecdsa.PublicKey) error {
	return nil
}

// MinimumStakeCachePeriod is the time period the cache maintains the result of
// the last HasMinimumStake result. We use the cache to minimize calls to
// Ethereum client.
const MinimumStakeCachePeriod = 12 * time.Hour

var errNoMinimumStake = fmt.Errorf("remote peer has no minimum stake")

// MinimumStakePolicy is a net.Firewall rule making sure the remote peer
// has a minimum stake of KEEP.
func MinimumStakePolicy(stakeMonitor chain.StakeMonitor) net.Firewall {
	return &minimumStakePolicy{
		stakeMonitor: stakeMonitor,
		cache:        cache.NewTimeCache(MinimumStakeCachePeriod),
	}
}

type minimumStakePolicy struct {
	stakeMonitor chain.StakeMonitor
	cache        *cache.TimeCache
}

func (msp *minimumStakePolicy) Validate(
	remotePeerPublicKey *ecdsa.PublicKey,
) error {
	networkPublicKey := key.NetworkPublic(*remotePeerPublicKey)
	address := key.NetworkPubKeyToEthAddress(&networkPublicKey)

	// First, check in the in-memory time cache to minimize hits to ETH client.
	// If the Keep client with the given chain address is in the cache it means
	// it has had a minimum stake the last HasMinimumStake was executed and
	// caching period has not elapsed yet.
	//
	// If the caching period elapsed, this check will return false and we
	// have to ask the chain about the current status.
	//
	// Similarly, if the client has no minimum stake the last time
	// HasMinimumStake was executed, we have to ask the chain about the current
	// status.
	msp.cache.Sweep()
	if msp.cache.Has(address) {
		return nil
	}

	hasMinimumStake, err := msp.stakeMonitor.HasMinimumStake(address)
	if err != nil {
		return fmt.Errorf(
			"could not validate remote peer's minimum stake: [%v]",
			err,
		)
	}

	if !hasMinimumStake {
		return errNoMinimumStake
	}

	// Add this address to the cache. We'll not hit HasMinimumStake again
	// for the entire caching period.
	msp.cache.Add(address)

	return nil
}
