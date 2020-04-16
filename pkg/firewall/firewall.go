package firewall

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

// MinimumStakePolicy is a net.Firewall rule making sure the remote peer
// has a minimum stake of KEEP.
func MinimumStakePolicy(stakeMonitor chain.StakeMonitor) net.Firewall {
	return &minimumStakePolicy{stakeMonitor: stakeMonitor}
}

type minimumStakePolicy struct {
	stakeMonitor chain.StakeMonitor
}

func (msp *minimumStakePolicy) Validate(
	remotePeerPublicKey *ecdsa.PublicKey,
) error {
	networkPublicKey := key.NetworkPublic(*remotePeerPublicKey)
	hasMinimumStake, err := msp.stakeMonitor.HasMinimumStake(
		key.NetworkPubKeyToEthAddress(&networkPublicKey),
	)
	if err != nil {
		return fmt.Errorf(
			"could not validate remote peer's minimum stake: [%v]",
			err,
		)
	}

	if !hasMinimumStake {
		return fmt.Errorf("remote peer has no minimum stake")
	}

	return nil
}
