package dkg

import (
	"fmt"
	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
	tbtcchain "github.com/keep-network/keep-core/pkg/tbtc/chain"
)

// TODO: Description
func Publish(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	dkgGroup *group.Group,
	membershipValidator *group.MembershipValidator,
	result *tbtcchain.DKGResult,
	channel net.BroadcastChannel,
	tbtcChain tbtcchain.Chain,
	blockCounter chain.BlockCounter,
	startBlockHeight uint64,
) error {
	initialState := &resultSigningState{
		channel:                 channel,
		tbtcChain:               tbtcChain,
		blockCounter:            blockCounter,
		member:                  NewSigningMember(logger, memberIndex, dkgGroup, membershipValidator),
		result:                  result,
		signatureMessages:       make([]*DKGResultHashSignatureMessage, 0),
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(logger, channel, blockCounter, initialState)

	lastState, _, err := stateMachine.Execute(startBlockHeight)
	if err != nil {
		return err
	}

	_, ok := lastState.(*resultSubmissionState)
	if !ok {
		return fmt.Errorf("execution ended on state %T", lastState)
	}

	return nil
}
