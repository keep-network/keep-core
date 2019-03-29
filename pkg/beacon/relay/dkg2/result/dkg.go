package result

import (
	"fmt"
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// SignAndSubmit executes Phase 13 and 14 of DKG as a state machine. First, the
// chosen result is hashed, signed, and sent over a broadcast channel. Then, all
// other signatures and results are received and accounted for. Those that match
// our own result and added to the list of votes. Finally, we submit the result
// along with everyone's votes.
func SignAndSubmit(
	privateKey *operator.PrivateKey,
	playerIndex group.MemberIndex,
	requestID *big.Int,
	dkgGroup *group.Group,
	result *relayChain.DKGResult,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	blockCounter chain.BlockCounter,
) error {
	initialState := &resultSigningState{
		channel:           channel,
		relayChain:        relayChain,
		blockCounter:      blockCounter,
		member:            NewSigningMember(playerIndex, dkgGroup, privateKey),
		requestID:         requestID,
		result:            result,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	initializeChannel(channel)

	stateMachine := state.NewMachine(channel, blockCounter, initialState)

	lastState, err := stateMachine.Execute()
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
