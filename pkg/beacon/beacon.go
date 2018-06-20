package beacon

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/membership"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
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
	relayChain relay.ChainInterface,
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

		_ = relayChain.SubmitGroupPublicKey(
			"test",
			member.GroupPublicKeyBytes(),
		).OnSuccess(func(data interface{}) {
			fmt.Printf(
				"Submission of public key: [%s].\n",
				data,
			)
		}).OnFailure(func(err error) {
			fmt.Printf(
				"Failed submission of public key: [%v].\n",
				err,
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

func checkChainParticipantState(relayChain relay.ChainInterface) (participantState, error) {
	// FIXME Zero in on the participant's current state per the chain.
	fmt.Println(relayChain)

	// FIXME This will return a real chain-based state: are we staked? What
	// FIXME groups are we in already, if any?
	return unstaked, nil
}

func libp2pConnected(relayChain relay.ChainInterface, handle chain.Handle) {
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
