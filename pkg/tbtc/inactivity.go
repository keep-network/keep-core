package tbtc

import (
	"context"
	"math/big"
	"sync"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/inactivity"
)

type inactivityClaimExecutor struct {
	chain   Chain
	signers []*signer

	protocolLatch *generator.ProtocolLatch
}

// TODO Consider moving all inactivity-related code to pkg/protocol/inactivity.
func newInactivityClaimExecutor(
	chain Chain,
	signers []*signer,
) *inactivityClaimExecutor {
	return &inactivityClaimExecutor{
		chain:   chain,
		signers: signers,
	}
}

func (ice *inactivityClaimExecutor) publishClaim(
	inactiveMembersIndexes []group.MemberIndex,
	heartbeatFailed bool,
) error {
	// TODO: Build a claim and launch the publish function for all
	//       the signers. The value of `heartbeat` should be true and
	//       `inactiveMembersIndices` should be empty.

	wg := sync.WaitGroup{}
	wg.Add(len(ice.signers))

	for _, currentSigner := range ice.signers {
		ice.protocolLatch.Lock()
		defer ice.protocolLatch.Unlock()

		go func(signer *signer) {
			// TODO: Launch claim publishing for members.
		}(currentSigner)
	}

	return nil
}

func (ice *inactivityClaimExecutor) publish(
	ctx context.Context,
	inactivityLogger log.StandardLogger,
	seed *big.Int,
	memberIndex group.MemberIndex,
	broadcastChannel net.BroadcastChannel,
	groupSize int,
	dishonestThreshold int,
	membershipValidator *group.MembershipValidator,
	inactivityClaim *inactivity.Claim,
) error {
	return inactivity.Publish(
		ctx,
		inactivityLogger,
		seed.Text(16),
		memberIndex,
		broadcastChannel,
		groupSize,
		dishonestThreshold,
		membershipValidator,
		newInactivityClaimSigner(ice.chain),
		newInactivityClaimSubmitter(),
		inactivityClaim,
	)
}
