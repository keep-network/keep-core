package dkg

import (
	"context"
	"fmt"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/ipfs/go-log"

	"math/big"
	"time"

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
	publicationDelayTime time.Duration,
	publicationDelayStep uint64,
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
			"waiting [%v] before DKG result publication",
		playerIndex,
		gjkrEndBlockHeight,
		publicationDelayTime,
	)

	// Calculations of group public key shares are time-expensive.
	// They takes around 3 minutes for each member when group size is 64.
	// To avoid desynchronization between members they are triggered in
	// a separate goroutine.
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

	// In order to keep members synchronized while group public key shares
	// are calculated, each member waits the same amount of time before
	// proceeding with DKG result publication. It is important to choose an
	// amount of time which will allow to complete calculations with high
	// probability.
	ctx, cancel := context.WithTimeout(context.Background(), publicationDelayTime)
	defer cancel()

	startPublicationBlockHeight := gjkrEndBlockHeight
	for {
		startPublicationBlockHeight += publicationDelayStep
		err := blockCounter.WaitForBlockHeight(startPublicationBlockHeight)
		if err != nil {
			return nil, fmt.Errorf(
				"[member:%v] wait for block height [%v] failed: [%v]",
				playerIndex,
				startPublicationBlockHeight,
				err,
			)
		}

		if ctx.Err() != nil {
			startPublicationBlockHeight += publicationDelayStep
			break
		}
	}

	var groupPublicKeyShares map[group.MemberIndex]*bn256.G2

	select {
	case groupPublicKeyShares = <-groupPublicKeySharesChan:
		// Group public key shares have been calculated on time.
	default:
		return nil, fmt.Errorf(
			"[member:%v] group public key shares have not been calculated on time",
			playerIndex,
		)
	}

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

	return &ThresholdSigner{
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
		groupPublicKeyShares: groupPublicKeyShares,
	}, nil
}
