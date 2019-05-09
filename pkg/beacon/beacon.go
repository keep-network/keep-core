package beacon

import (
	"context"
	"fmt"
	"os"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupregistry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
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

	staker, err := stakeMonitor.StakerFor(stakingID)
	if err != nil {
		return err
	}

	groupRegistry := groupregistry.NewGroupRegistry(relayChain)

	node := relay.NewNode(
		staker,
		netProvider,
		blockCounter,
		chainConfig,
		groupRegistry,
	)

	relayChain.OnRelayEntryRequested(func(request *event.Request) {
		fmt.Printf("New relay entry requested [%+v]\n", request)

		go node.GenerateRelayEntryIfEligible(
			request.RequestID,
			request.PreviousEntry,
			request.Seed,
			relayChain,
			request.GroupPublicKey,
			request.BlockNumber,
		)
	})

	relayChain.OnRelayEntryGenerated(func(entry *event.Entry) {
		fmt.Printf("New relay entry generated [%+v]\n", entry)

		go func() {
			err := node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				entry.Value.Bytes(),
				entry.RequestID,
				entry.Seed,
				entry.BlockNumber,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Tickets submission failed: [%v]\n", err)
			}
		}()
	})

	relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
		fmt.Printf("New group registered on chain [%+v]\n", registration)
		go groupRegistry.UnregisterDeletedGroups()
	})

	return nil
}
