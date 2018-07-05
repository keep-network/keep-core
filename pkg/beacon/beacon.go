package beacon

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"

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
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	netProvider net.Provider,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	curParticipantState, err := checkParticipantState()
	if err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state, aborting: [%s]", err))
	}

	// FIXME Nuke post-M1 when we plug in real staking stuff.
	proceed := &sync.WaitGroup{}
	proceed.Add(1)
	relayChain.AddStaker(netProvider.ID().String()).
		OnComplete(func(_ *event.StakerRegistration, err error) {
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Failed to register staker with id [%v].",
					netProvider.ID().String(),
				)
				curParticipantState = unstaked
				proceed.Done()
				return
			}

			proceed.Done()
		})
	proceed.Wait()

	switch curParticipantState {
	case unstaked:
		// check for stake command-line parameter to initialize staking?
		return fmt.Errorf("account is unstaked")
	default:
		node := relay.NewNode(
			netProvider.ID().String(),
			netProvider,
			blockCounter,
			chainConfig,
		)

		relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
			entryBigInt := &big.Int{}
			entryBigInt.SetBytes(entry.Value[:])
			node.JoinGroupIfEligible(relayChain, entry.RequestID, entryBigInt)
		})

		relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
			node.RegisterGroup(registration.RequestID, registration.GroupPublicKey)
		})

	}

	<-ctx.Done()

	return nil
}

func checkParticipantState() (participantState, error) {
	return staked, nil
}
