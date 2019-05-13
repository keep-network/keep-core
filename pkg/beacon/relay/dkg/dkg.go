package dkg

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	dkgResult "github.com/keep-network/keep-core/pkg/beacon/relay/dkg/result"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// ExecuteDKG runs the full distributed key generation lifecycle.
func ExecuteDKG(
	requestID *big.Int,
	seed *big.Int,
	index int, // starts with 0
	groupSize int,
	threshold int,
	startBlockHeight uint64,
	blockCounter chain.BlockCounter,
	relayChain relayChain.Interface,
	channel net.BroadcastChannel,
) (*ThresholdSigner, error) {
	// The staker index should begin with 1
	playerIndex := group.MemberIndex(index + 1)
	err := playerIndex.Validate()
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] could not start DKG: [%v]",
			playerIndex,
			err,
		)
	}

	gjkrResult, gjkrEndBlockHeight, err := gjkr.Execute(
		playerIndex,
		blockCounter,
		channel,
		threshold,
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

	err = dkgResult.Publish(
		playerIndex,
		requestID,
		gjkrResult.Group,
		gjkrResult,
		channel,
		relayChain,
		blockCounter,
		gjkrEndBlockHeight,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"[member:%v] DKG result publication process failed [%v]",
			playerIndex,
			err,
		)
	}

	return &ThresholdSigner{
		MemberIndex:          playerIndex,
		GroupPublicKey:       gjkrResult.GroupPublicKey,
		GroupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
	}, nil
}
