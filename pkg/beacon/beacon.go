package beacon

import (
	"context"
	"encoding/hex"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-beacon")

// Initialize kicks off the random beacon by initializing internal state,
// ensuring preconditions like staking are met, and then kicking off the
// internal random beacon implementation. Returns an error if this failed,
// otherwise enters a blocked loop.
func Initialize(
	ctx context.Context,
	stakingID string,
	chainHandle chain.Handle,
	netProvider net.Provider,
	persistence persistence.Handle,
) error {
	relayChain := chainHandle.ThresholdRelay()
	chainConfig, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	stakeMonitor, err := chainHandle.StakeMonitor()
	if err != nil {
		return err
	}

	staker, err := stakeMonitor.StakerFor(stakingID)
	if err != nil {
		return err
	}

	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return err
	}

	signing := chainHandle.Signing()

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
		logger.Infof(
			"new relay entry requested at block [%v] from group [%x] using "+
				"previous entry [%v] and seed [%v]",
			request.BlockNumber,
			request.GroupPublicKey,
			request.PreviousEntry,
			request.Seed,
		)

		if node.IsInGroup(request.GroupPublicKey) {
			go node.GenerateRelayEntry(
				request.PreviousEntry,
				request.Seed,
				relayChain,
				request.GroupPublicKey,
				request.BlockNumber,
			)
		} else {
			go node.MonitorRelayEntry(
				relayChain,
				request.BlockNumber,
				chainConfig,
			)
		}
	})

	relayChain.OnGroupSelectionStarted(func(event *event.GroupSelectionStart) {
		logger.Infof("group selection started: [%+v]", event)

		onGroupSelected := func(group *groupselection.Result) {
			for _, staker := range group.SelectedStakers {
				logger.Infof(
					"new candidate group member: [0x%v]",
					hex.EncodeToString(staker),
				)
			}
			node.JoinGroupIfEligible(
				relayChain,
				signing,
				group,
				event.NewEntry,
			)
		}

		go func() {
			err := groupselection.CandidateToNewGroup(
				relayChain,
				blockCounter,
				chainConfig,
				staker,
				event.NewEntry,
				event.BlockNumber,
				onGroupSelected,
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
