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

// GroupInfo holds the essential information about the group which should
// perform DKG.
type GroupInfo struct {
	Size                int
	DishonestThreshold  int
	MembershipValidator group.MembershipValidator
}

// ExecuteDKG runs the distributed key generation protocol.
func ExecuteDKG(
	seed *big.Int,
	index int, // starts with 0
	groupInfo *GroupInfo,
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
		groupInfo.Size,
		groupInfo.DishonestThreshold,
		groupInfo.MembershipValidator,
		blockCounter,
		channel,
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
		gjkrResult.Group,
		groupInfo.MembershipValidator,
		gjkrResult,
		channel,
		relayChain,
		signing,
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
		memberIndex:          playerIndex,
		groupPublicKey:       gjkrResult.GroupPublicKey,
		groupPrivateKeyShare: gjkrResult.GroupPrivateKeyShare,
	}, nil
}
