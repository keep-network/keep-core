package dkg

import (
	"bytes"
	"fmt"
	"math/big"

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

	startPublicationBlockHeight := gjkrEndBlockHeight

	err = dkgResult.Publish(
		playerIndex,
		gjkrResult.Group,
		gjkrResult,
		channel,
		relayChain,
		signing,
		blockCounter,
		startPublicationBlockHeight,
	)
	if err != nil {
		// Result publication failed. It means that either the result this
		// member proposed is not supported by the majority of group members or
		// that the chain interaction failed. In either case, we observe the
		// chain for the result published by any other group member and based
		// on that, we decide whether we should stay in the final group
		// or drop our membership.
		logger.Warningf(
			"[member:%v] DKG result publication process failed [%v]",
			playerIndex,
			err,
		)

		dkgResultEvent, err := waitForDkgResultEvent(
			dkgResultChannel,
			startPublicationBlockHeight,
			relayChain,
			blockCounter,
		)
		if err != nil {
			return nil, err
		}

		if !shouldStayInGroup(playerIndex, gjkrResult, dkgResultEvent) {
			return nil, fmt.Errorf(
				"[member:%v] could not stay in the group",
				playerIndex,
			)
		}
	}

	return &ThresholdSigner{
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
	}, nil
}

func waitForDkgResultEvent(
	dkgResultChannel chan *event.DKGResultSubmission,
	startPublicationBlockHeight uint64,
	relayChain relayChain.Interface,
	blockCounter chain.BlockCounter,
) (*event.DKGResultSubmission, error) {
	config, err := relayChain.GetConfig()
	if err != nil {
		return nil, err
	}

	timeoutBlock := startPublicationBlockHeight +
		dkgResult.SigningStateBlocks() +
		(uint64(config.GroupSize) * config.ResultPublicationBlockStep)

	timeoutBlockChannel, err := blockCounter.BlockHeightWaiter(timeoutBlock)
	if err != nil {
		return nil, err
	}

	select {
	case dkgResultEvent := <-dkgResultChannel:
		return dkgResultEvent, nil
	case <-timeoutBlockChannel:
		return nil, fmt.Errorf("DKG result publication timed out")
	}
}

func shouldStayInGroup(
	memberIndex group.MemberIndex,
	gjkrResult *gjkr.Result,
	dkgResultEvent *event.DKGResultSubmission,
) bool {
	groupPublicKey, err := gjkrResult.GroupPublicKeyBytes()
	if err != nil {
		return false
	}

	supportsSameGroupPublicKey := bytes.Equal(
		groupPublicKey,
		dkgResultEvent.GroupPublicKey,
	)

	// If member didn't support the same group public key,
	// it could not stay in the group.
	if !supportsSameGroupPublicKey {
		return false
	}

	// If member is considered as misbehaved, it could not stay in the group.
	for _, misbehaved := range dkgResultEvent.Misbehaved {
		if memberIndex == misbehaved {
			return false
		}
	}

	return true
}
