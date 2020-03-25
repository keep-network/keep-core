package dkg

import (
	"bytes"
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/ipfs/go-log"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-dkg")

const publicationDelayBlocks = 20

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

	logger.Debugf(
		"[member:%v] GJKR ended at block [%v]; "+
			"waiting [%v] blocks before DKG result publication",
		playerIndex,
		gjkrEndBlockHeight,
		publicationDelayBlocks,
	)

	// Calculation of group public key shares is time-expensive.
	// It takes around 3 minutes for each member when group size is 64.
	// To avoid desynchronization between members it is triggered in
	// a separate goroutine. Apart from that, result publication is delayed
	// to reduce the gap between successful on-chain group registration and
	// the moment when threshold signers are ready. This is important
	// because during this gap the group is blind for incoming relay entry
	// requests.
	groupPublicKeySharesChan := make(chan map[group.MemberIndex]*bn256.G2)
	go func() {
		logger.Debugf(
			"[member:%v] starting group public key shares calculation",
			playerIndex,
		)

		groupPublicKeyShares := gjkrResult.GroupPublicKeyShares()

		logger.Debugf(
			"[member:%v] group public key shares calculated",
			playerIndex,
		)

		groupPublicKeySharesChan <- groupPublicKeyShares
	}()

	startPublicationBlockHeight := gjkrEndBlockHeight + publicationDelayBlocks

	logger.Debugf(
		"[member:%v] DKG result publication scheduled for block [%v]",
		playerIndex,
		startPublicationBlockHeight,
	)

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

		if err := decideMemberFate(
			playerIndex,
			gjkrResult,
			dkgResultChannel,
			startPublicationBlockHeight,
			relayChain,
			blockCounter,
		); err != nil {
			return nil, err
		}
	}

	// Wait for group public key shares calculation outcome.
	groupPublicKeyShares := <-groupPublicKeySharesChan

	return &ThresholdSigner{
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
		groupPublicKeyShares: groupPublicKeyShares,
	}, nil
}

// decideMemberFate decides what the member will do in case it failed
// publishing its DKG result. Member can stay in the group if it
// supports the same group public key as the one registered on-chain and
// the member is not considered as misbehaving by the group.
func decideMemberFate(
	playerIndex group.MemberIndex,
	gjkrResult *gjkr.Result,
	dkgResultChannel chan *event.DKGResultSubmission,
	startPublicationBlockHeight uint64,
	relayChain relayChain.Interface,
	blockCounter chain.BlockCounter,
) error {
	dkgResultEvent, err := waitForDkgResultEvent(
		dkgResultChannel,
		startPublicationBlockHeight,
		relayChain,
		blockCounter,
	)
	if err != nil {
		return err
	}

	groupPublicKey, err := gjkrResult.GroupPublicKeyBytes()
	if err != nil {
		return err
	}

	// If member don't support the same group public key, it could not stay
	// in the group.
	if !bytes.Equal(groupPublicKey, dkgResultEvent.GroupPublicKey) {
		return fmt.Errorf(
			"[member:%v] could not stay in the group because "+
				"member do not support the same group public key",
			playerIndex,
		)
	}

	// If member is considered as misbehaved, it could not stay in the group.
	for _, misbehaved := range dkgResultEvent.Misbehaved {
		if playerIndex == misbehaved {
			return fmt.Errorf(
				"[member:%v] could not stay in the group because "+
					"member is considered as misbehaving",
				playerIndex,
			)
		}
	}

	return nil
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
		dkgResult.PrePublicationBlocks() +
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
