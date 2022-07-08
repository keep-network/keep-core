package firewall

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// Disabled is an empty Firewall implementation enforcing no rules
// on the connection.
var Disabled = &noFirewall{}

type noFirewall struct{}

func (nf *noFirewall) Validate(remotePeerPublicKey *operator.PublicKey) error {
	return nil
}

const (
	// PositiveMinimumStakeCachePeriod is the time period the cache maintains
	// the positive result of the last HasMinimumStake check.
	// We use the cache to minimize calls to Ethereum client.
	PositiveMinimumStakeCachePeriod = 12 * time.Hour

	// NegativeMinimumStakeCachePeriod is the time period the cache maintains
	// the negative result of the last HasMinimumStake check.
	// We use the cache to minimize calls to Ethereum client.
	NegativeMinimumStakeCachePeriod = 1 * time.Hour
)

var errNoMinimumStake = fmt.Errorf("remote peer has no minimum stake")

// MinimumStakePolicy is a net.Firewall rule making sure the remote peer
// has a stake delegation in the Threshold TokenStaking contract and the minimum
// authorization required by the application
func MinimumStakePolicy(stakeMonitor chain.StakeMonitor) net.Firewall {
	return &minimumStakePolicy{
		stakeMonitor: stakeMonitor,
	}
}

type minimumStakePolicy struct {
	stakeMonitor chain.StakeMonitor
}

func (msp *minimumStakePolicy) Validate(
	remotePeerPublicKey *operator.PublicKey,
) error {
	// TODO: Should caching results for given operators be used?
	hasMinimumStake, err := msp.stakeMonitor.HasMinimumStake(remotePeerPublicKey)
	if err != nil {
		return fmt.Errorf(
			"could not validate remote peer's minimum stake: [%v]",
			err,
		)
	}

	if !hasMinimumStake {
		return errNoMinimumStake
	}

	return nil
}
