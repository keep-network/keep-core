package beacon

import (
	"context"
	"fmt"
	"math/big"

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
	chainConfig.TicketInitialSubmissionTimeout = 2
	chainConfig.TicketReactiveSubmissionTimeout = 3
	chainConfig.TicketChallengeTimeout = 4
	chainConfig.MinimumStake = big.NewInt(40000000)
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
			fmt.Printf("Saw new relay entry request [%+v]\n", request)

			if request.PreviousValue == nil {
				request.PreviousValue = big.NewInt(0)
			}

			go node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				request.PreviousValue.Bytes(),
				request.RequestID,
				request.Seed,
				chainConfig.GroupSize,
			)
		})

		relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
			fmt.Printf("Saw new relay entry submitted [%+v]\n", entry)
			// new entry generated, try to join the group
			go node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				entry.Value.Bytes(),
				entry.RequestID,
				entry.Seed,
				chainConfig.GroupSize,
			)

			nextRequestID := new(big.Int).Add(entry.RequestID, big.NewInt(1))

			go node.GenerateRelayEntryIfEligible(
				nextRequestID,
				entry.Value,
				entry.Seed,
				relayChain,
			)
		})

		relayChain.OnDKGResultPublished(func(publishedResult *event.DKGResultPublication) {
			fmt.Printf("Saw new group registered [%+v]\n", publishedResult)
			node.RegisterGroup(
				publishedResult.RequestID.String(),
				publishedResult.GroupPublicKey,
			)
		})
	}

	<-ctx.Done()

	return nil
}

func checkParticipantState() (participantState, error) {
	return staked, nil
}
