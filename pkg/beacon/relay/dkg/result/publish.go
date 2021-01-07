package result

import (
	"fmt"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform DKG result publication protocol interactions by registering all the
// required protocol message unmarshallers.
// The channel needs to be fully initialized before Publish is called.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &DKGResultHashSignatureMessage{}
	})
}

// Publish executes Phase 13 and 14 of DKG as a state machine. First, the
// chosen result is hashed, signed, and sent over a broadcast channel. Then, all
// other signatures and results are received and accounted for. Those that match
// our own result and added to the list of votes. Finally, we submit the result
// along with everyone's votes.
func Publish(
	memberIndex group.MemberIndex,
	dkgGroup *group.Group,
	membershipValidator group.MembershipValidator,
	result *gjkr.Result,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	signing chain.Signing,
	blockCounter chain.BlockCounter,
	startBlockHeight uint64,
) error {
	initialState := &resultSigningState{
		channel:                 channel,
		relayChain:              relayChain,
		signing:                 signing,
		blockCounter:            blockCounter,
		member:                  NewSigningMember(memberIndex, dkgGroup, membershipValidator),
		result:                  convertGjkrResult(result),
		signatureMessages:       make([]*DKGResultHashSignatureMessage, 0),
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(
		channel,
		blockCounter,
		initialState,
		2*dkgGroup.GroupSize(),
	)

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
