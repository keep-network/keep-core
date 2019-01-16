package beacon

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
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
	stakingID string,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	stakeMonitor chain.StakeMonitor,
	netProvider net.Provider,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}
	chainConfig.TicketReactiveSubmissionTimeout = 2
	chainConfig.TicketChallengeTimeout = 3
	fmt.Printf("Got chainconfig [%+v]\n", chainConfig)

	curParticipantState, err := checkParticipantState()
	if err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state, aborting: [%s]", err))
	}

	staker, err := stakeMonitor.StakerFor(stakingID)
	if err != nil {
		return err
	}

	node := relay.NewNode(
		staker,
		netProvider,
		blockCounter,
		chainConfig,
	)

	switch curParticipantState {
	case unstaked:
		// check for stake command-line parameter to initialize staking?
		return fmt.Errorf("account is unstaked")
	default:
		relayChain.OnRelayEntryRequested(func(request *event.Request) {
			fmt.Printf("New entry requested [%+v]\n", request)
			node.GenerateRelayEntryIfEligible(request, relayChain)
		})

		relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
			fmt.Printf("Saw new relay entry [%+v]\n", entry)
			// new entry generated, try to join the group
			node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				entry.Value.Bytes(),
				entry.RequestID,
				entry.Seed,
			)
		})

		relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
			fmt.Printf("New group registered [%+v]\n", registration)
			node.RegisterGroup(
				registration.RequestID.String(),
				registration.GroupPublicKey,
			)
		})
	}

	<-ctx.Done()

	return nil
}

func checkParticipantState() (participantState, error) {
	return staked, nil
}
