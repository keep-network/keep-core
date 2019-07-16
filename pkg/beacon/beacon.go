package beacon

import (
	"context"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/persistence"
)

var logger = log.Logger("keep-beacon")

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
	persistence persistence.Handle,
) error {
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	staker, err := stakeMonitor.StakerFor(stakingID)
	if err != nil {
		return err
	}

	groupRegistry := registry.NewGroupRegistry(relayChain, persistence)
	groupRegistry.LoadExistingGroups()

	node := relay.NewNode(
		staker,
		netProvider,
		blockCounter,
		chainConfig,
		groupRegistry,
	)

	relayChain.OnSignatureRequested(func(request *event.Request) {
		logger.Infof("new relay entry requested: [%+v]", request)

		go node.GenerateRelayEntryIfEligible(
			request.SigningId,
			request.PreviousEntry,
			request.Seed,
			relayChain,
			request.GroupPublicKey,
			request.BlockNumber,
		)
	})

	relayChain.OnGroupSelectionStarted(func(event *event.GroupSelectionStart) {
		logger.Infof("group selection started: [%+v]", event)

		go func() {
			err := node.SubmitTicketsForGroupSelection(
				relayChain,
				blockCounter,
				event.NewEntry.Bytes(),
				event.Seed,
				event.BlockNumber,
			)
			if err != nil {
				logger.Errorf("Tickets submission failed: [%v]", err)
			}
		}()
	})

	relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
		logger.Infof("new group registered on chain: [%+v]", registration)
		go groupRegistry.UnregisterStaleGroups()
	})

	return nil
}
