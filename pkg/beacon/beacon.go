package beacon

import (
	"context"
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/membership"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

type participantState int

const (
	unstaked participantState = iota
	staked
	waitingForGroup
	inIncompleteGroup
	inCompleteGroup
	inInitializingGroup
	inInitializedGroup
	inActivatingGroup
	inActiveGroup
)

// Initialize kicks off the random beacon by initializing internal state,
// ensuring preconditions like staking are met, and then kicking off the
// internal random beacon implementation. Returns an error if this failed,
// otherwise enters a blocked loop.
func Initialize(
	ctx context.Context,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	netProvider net.Provider,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	channel, err := netProvider.ChannelFor("test")
	if err != nil {
		return err
	}

	curParticipantState, err := checkParticipantState()
	if err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state, aborting: [%s]", err))
	}

	switch curParticipantState {
	case unstaked:
		// check for stake command-line parameter to initialize staking?
		return fmt.Errorf("account is unstaked")
	default:
		member, err := dkg.ExecuteDKG(
			blockCounter,
			channel,
			chainConfig.GroupSize,
			chainConfig.Threshold,
		)
		if err != nil {
			return err
		}

		err = relayChain.SubmitGroupPublicKey("test", member.GroupPublicKeyBytes())
		if err != nil {
			return err
		}

		relayChain.OnGroupPublicKeySubmissionFailed(func(id string, errorMessage string) {
			fmt.Printf(
				"Failed submission of public key %s: [%s].\n",
				id,
				errorMessage,
			)
		})
		relayChain.OnGroupPublicKeySubmitted(func(id string, activationBlock *big.Int) {
			fmt.Printf(
				"Public key submitted for %s; activating at block %v.\n",
				id,
				activationBlock,
			)
		})

		fmt.Printf(
			"Submitting public key for member %s, group %s\n",
			member.MemberID(),
			"test",
		)
	}

	<-ctx.Done()

	return nil
}

func checkParticipantState() (participantState, error) {
	return staked, nil
}

func checkChainParticipantState(relayChain relaychain.Interface) (participantState, error) {
	// FIXME Zero in on the participant's current state per the chain.
	fmt.Println(relayChain)

	// FIXME This will return a real chain-based state: are we staked? What
	// FIXME groups are we in already, if any?
	return unstaked, nil
}

func libp2pConnected(relayChain relaychain.Interface, handle chain.Handle) {
	if participantState, err := checkChainParticipantState(relayChain); err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state from libp2p, aborting: [%s]", err))
	} else {
		switch participantState {
		case staked:
			membership.WaitForGroup()
		case waitingForGroup:
			membership.WaitForGroup()
		case inIncompleteGroup:
			membership.WaitForGroupCompletion()
		case inCompleteGroup:
			membership.InitializeMembership()
		case inInitializingGroup:
			membership.InitializeMembership()
		case inInitializedGroup:
			membership.ActivateMembership()
		case inActivatingGroup:
			membership.ActivateMembership()
		case inActiveGroup:
			// FIXME We should have a non-empty state at this point ;)
			entry.ServeRequests(relay.EmptyState())
		default:
			panic(fmt.Sprintf("Unexpected participant state [%d].", participantState))
		}
	}
}
