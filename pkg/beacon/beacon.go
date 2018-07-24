package beacon

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

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

	node := relay.NewNode(
		netProvider.ID().String(),
		netProvider,
		blockCounter,
		chainConfig,
	)

	// FIXME Nuke post-M1 when we plug in real staking stuff.
	var (
		proceed   sync.WaitGroup
		stakingID = netProvider.ID().String()[:32]
	)

	proceed.Add(1)
	relayChain.AddStaker(stakingID).
		OnComplete(func(stake *event.StakerRegistration, err error) {
			defer proceed.Done()

			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Failed to register staker with id [%v].",
					stakingID,
				)
				curParticipantState = unstaked
				return
			}
		})
	proceed.Wait()

	switch curParticipantState {
	case unstaked:
		// check for stake command-line parameter to initialize staking?
		return fmt.Errorf("account is unstaked")
	default:
		relayChain.OnStakerAdded(func(staker *event.StakerRegistration) {
			node.AddStaker(staker.Index, staker.GroupMemberID)
		})

		// Retry until we can sync our staking list
		syncStakingListWithRetry(&node, relayChain)

		relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
			entryBigInt := (&big.Int{}).SetBytes(entry.Value[:])
			node.JoinGroupIfEligible(relayChain, entry.RequestID, entryBigInt)
		})

		relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
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

func syncStakingListWithRetry(node *relay.Node, relayChain relaychain.Interface) {
	for {
		t := time.NewTimer(1)
		defer t.Stop()

		select {
		case <-t.C:
			list, err := relayChain.GetStakerList()
			if err != nil {
				fmt.Printf(
					"failed to sync staking list: [%v], retrying...\n",
					err,
				)

				// FIXME: exponential backoff
				t.Reset(3 * time.Second)
				continue
			}

			node.SyncStakingList(list)

			// exit this loop when we've successfully synced
			return
		}
	}
}
