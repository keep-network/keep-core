package dkg

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	"github.com/ipfs/go-log"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-dkg")

// ExecuteDKG runs the full distributed key generation lifecycle.
func ExecuteDKG(
	seed *big.Int,
	index uint8, // starts with 0
	groupSize int,
	dishonestThreshold int,
	startBlockHeight uint64,
	blockCounter chain.BlockCounter,
	relayChain relayChain.Interface,
	signing chain.Signing,
	channel net.BroadcastChannel,
) (*ThresholdSigner, error) {
	// The staker index should begin with 1
	playerIndex := group.MemberIndex(index + 1)

	gjkr.RegisterUnmarshallers(channel)
	dkgResult.RegisterUnmarshallers(channel)

	gjkrResult, gjkrEndBlockHeight, err := gjkr.Execute(
		playerIndex,
		groupSize,
		blockCounter,
		channel,
		dishonestThreshold,
		seed,
		startBlockHeight,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] GJKR execution failed [%v]",
			playerIndex,
			err,
		)
	}

	dkgResultChannel := make(chan *event.DKGResultSubmission)
	dkgResultSubscription, err := relayChain.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			dkgResultChannel <- event
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] could not create DKG result subscription [%v]",
			playerIndex,
			err,
		)
	}
	defer dkgResultSubscription.Unsubscribe()

	err = dkgResult.Publish(
		playerIndex,
		gjkrResult.Group,
		gjkrResult,
		channel,
		relayChain,
		signing,
		blockCounter,
		gjkrEndBlockHeight,
	)
	if err != nil {
		logger.Warningf(
			"[member:%v] DKG result publication process failed [%v]; "+
				"checking conditional membership possibility",
			playerIndex,
			err,
		)

		// In case of DKG timeout, this context will prevent endless waiting.
		// DKG result should be published after 3 * 64 = 192 blocks.
		// Assuming even 30 seconds for a block, it can take a bit less than
		// 2 hours. Waiting 3 hours seems to be a reasonable value with a
		// security margin.
		ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Hour)
		defer cancelCtx()

		select {
		case dkgResultEvent := <-dkgResultChannel:
			if shouldStayInGroup(playerIndex, gjkrResult, dkgResultEvent) {
				logger.Debugf(
					"[member:%v] conditional membership is possible",
					playerIndex,
				)
			} else {
				return nil, fmt.Errorf(
					"[member:%v] DKG result publication process failed [%v] "+
						"and conditional membership is not possible",
					playerIndex,
					err,
				)
			}
		case <-ctx.Done():
			return nil, fmt.Errorf(
				"[member:%v] DKG result publication process failed [%v] "+
					"and conditional membership check timed out",
				playerIndex,
				err,
			)
		}
	}

	return &ThresholdSigner{
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
	}, nil
}

func shouldStayInGroup(
	memberIndex group.MemberIndex,
	gjkrResult *gjkr.Result,
	dkgResultEvent *event.DKGResultSubmission,
) bool {
	supportsSameGroupPublicKey := bytes.Equal(
		dkgResult.ConvertGjkrResult(gjkrResult).GroupPublicKey,
		dkgResultEvent.GroupPublicKey,
	)

	// If member didn't support the same group public key, it could not be
	// a conditional member of the group.
	if !supportsSameGroupPublicKey {
		return false
	}

	// If member is considered as misbehaved, it could not be a conditional
	// member of the group.
	for _, misbehaved := range dkgResultEvent.Misbehaved {
		if memberIndex == misbehaved {
			return false
		}
	}

	return true
}
