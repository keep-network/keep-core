package dkg

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ipfs/go-log"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/beacon/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/group"
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
	membershipValidator group.MembershipValidator,
	startBlockHeight uint64,
	blockCounter chain.BlockCounter,
	beaconChain beaconchain.Interface,
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
		membershipValidator,
		startBlockHeight,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] GJKR execution failed [%v]",
			playerIndex,
			err,
		)
	}

	startPublicationBlockHeight := gjkrEndBlockHeight

	dkgResultChannel := make(chan *event.DKGResultSubmission)
	dkgResultSubscription := beaconChain.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			dkgResultChannel <- event
		},
	)
	defer dkgResultSubscription.Unsubscribe()

	err = dkgResult.Publish(
		playerIndex,
		gjkrResult.Group,
		membershipValidator,
		gjkrResult,
		channel,
		beaconChain,
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
			beaconChain,
			blockCounter,
		); err != nil {
			return nil, err
		}
	}

	return &ThresholdSigner{
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
		groupPublicKeyShares: gjkrResult.GroupPublicKeyShares(),
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
	beaconChain beaconchain.Interface,
	blockCounter chain.BlockCounter,
) error {
	dkgResultEvent, err := waitForDkgResultEvent(
		dkgResultChannel,
		startPublicationBlockHeight,
		beaconChain,
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
	beaconChain beaconchain.Interface,
	blockCounter chain.BlockCounter,
) (*event.DKGResultSubmission, error) {
	config := beaconChain.GetConfig()

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
