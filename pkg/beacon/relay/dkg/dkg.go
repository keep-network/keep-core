package dkg

import (
	"fmt"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/ipfs/go-log"

	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
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
		return nil, fmt.Errorf(
			"[member:%v] DKG result publication process failed [%v]",
			playerIndex,
			err,
		)
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
