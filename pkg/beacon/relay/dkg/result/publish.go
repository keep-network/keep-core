package result

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// Publish executes Phase 13 and 14 of DKG as a state machine. First, the
// chosen result is hashed, signed, and sent over a broadcast channel. Then, all
// other signatures and results are received and accounted for. Those that match
// our own result and added to the list of votes. Finally, we submit the result
// along with everyone's votes.
func Publish(
	playerIndex group.MemberIndex,
	signingId *big.Int,
	dkgGroup *group.Group,
	result *gjkr.Result,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	blockCounter chain.BlockCounter,
	startBlockHeight uint64,
) error {
	privateKey, _ := relayChain.GetKeys()
	initialState := &resultSigningState{
		channel:                 channel,
		relayChain:              relayChain,
		blockCounter:            blockCounter,
		member:                  NewSigningMember(playerIndex, dkgGroup, privateKey),
		signingId:               signingId,
		result:                  convertResult(result, dkgGroup.GroupSize()),
		signatureMessages:       make([]*DKGResultHashSignatureMessage, 0),
		signingStartBlockHeight: startBlockHeight,
	}

	initializeChannel(channel)

	stateMachine := state.NewMachine(channel, blockCounter, initialState)

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

// initializeChannel initializes a given broadcast channel to be able to
// perform distributed key generation interactions.
func initializeChannel(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &DKGResultHashSignatureMessage{}
	})
}
