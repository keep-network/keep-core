package beacon

import (
	"context"
	"encoding/hex"
	"sync"

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

	pendingGroupSelections := &event.GroupSelectionTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	pendingRelayRequests := &event.RelayRequestTrack{
		Data:  make(map[string]bool),
		Mutex: &sync.Mutex{},
	}

	node.ResumeSigningIfEligible(relayChain, signing)

	relayChain.OnRelayEntryRequested(func(request *event.Request) {
		currentEntryStartBlock, err := relayChain.CurrentRequestStartBlock()
		if err != nil {
			logger.Warningf("could not check current request start block")
			return
		}

		// Ignore if the block number of the received request isn't equal to
		// the current entry start block from the chain. Most likely, there
		// are some problems with chain synchronization and multiple events
		// arrived at the same time but only the last should be handled.
		if request.BlockNumber != currentEntryStartBlock.Uint64() {
			logger.Warningf(
				"ignoring relay entry request with block number [%v]; "+
					"current entry start block stored on-chain is [%v]",
				request.BlockNumber,
				currentEntryStartBlock.Uint64(),
			)
			return
		}

		if node.IsInGroup(request.GroupPublicKey) {
			go func() {
				previousEntry := hex.EncodeToString(request.PreviousEntry[:])

				if ok := pendingRelayRequests.Add(previousEntry); !ok {
					logger.Errorf(
						"relay entry requested event with previous entry "+
							"[0x%x] has been registered already",
						request.PreviousEntry,
					)
					return
				}

				defer pendingRelayRequests.Remove(previousEntry)

				logger.Infof(
					"new relay entry requested at block [%v] from group "+
						"[0x%x] using previous entry [0x%x]",
					request.BlockNumber,
					request.GroupPublicKey,
					request.PreviousEntry,
				)

				node.GenerateRelayEntry(
					request.PreviousEntry,
					relayChain,
					signing,
					request.GroupPublicKey,
					request.BlockNumber,
				)
			}()
		} else {
			go node.ForwardSignatureShares(request.GroupPublicKey)
		}

		go node.MonitorRelayEntry(
			relayChain,
			request.BlockNumber,
			chainConfig,
		)
	})

	relayChain.OnGroupSelectionStarted(func(event *event.GroupSelectionStart) {
		onGroupSelected := func(group *groupselection.Result) {
			for index, staker := range group.SelectedStakers {
				logger.Infof(
					"new candidate group member [0x%v] with index [%v]",
					hex.EncodeToString(staker),
					index,
				)
			}
			node.JoinGroupIfEligible(
				relayChain,
				signing,
				group,
				event.NewEntry,
			)
		}

		newEntry := event.NewEntry.Text(16)
		go func() {
			if ok := pendingGroupSelections.Add(newEntry); !ok {
				logger.Errorf(
					"group selection event with seed [0x%x] has been registered already",
					event.NewEntry,
				)
				return
			}

			defer pendingGroupSelections.Remove(newEntry)

			logger.Infof(
				"group selection started with seed [0x%x] at block [%v]",
				event.NewEntry,
				event.BlockNumber,
			)

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
		logger.Infof(
			"new group with public key [0x%x] registered on-chain at block [%v]",
			registration.GroupPublicKey,
			registration.BlockNumber,
		)
		go groupRegistry.UnregisterStaleGroups()
	})

	return nil
}
