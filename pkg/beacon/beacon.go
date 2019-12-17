package beacon

import (
	"context"
	"encoding/hex"
	"math/big"

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

	go func() {
		currentBlock, err := blockCounter.CurrentBlock()
		if err != nil {
			logger.Errorf("could not get current block [%v]", err)
		}

		logger.Infof("current block is [%v]", currentBlock)

		waiter, err := blockCounter.BlockHeightWaiter(currentBlock + 1)
		if err != nil {
			logger.Errorf("could not create block height waiter [%v]", err)
		}

		for {
			select {
			case b := <-waiter:
				logger.Infof("got block [%v]", b)
				currentBlock = b
				waiter, err = blockCounter.BlockHeightWaiter(currentBlock + 1)
				if err != nil {
					logger.Errorf("could not create block height waiter [%v]", err)
					break
				}

				chainBlock, err := chainHandle.ThresholdRelay().TestBlock(big.NewInt(int64(currentBlock)))
				if err != nil {
					logger.Errorf("failed executing on-chain block check [%v]", err)
					break
				}

				logger.Infof("confirmed of-chain block [%v] against current chain block [%v]", currentBlock, chainBlock)
			}
		}
	}()

	relayChain.OnRelayEntryRequested(func(request *event.Request) {
		logger.Infof(
			"new relay entry requested at block [%v] from group [0x%x] using "+
				"previous entry [0x%x]",
			request.BlockNumber,
			request.GroupPublicKey,
			request.PreviousEntry,
		)

		if node.IsInGroup(request.GroupPublicKey) {
			go node.GenerateRelayEntry(
				request.PreviousEntry,
				relayChain,
				request.GroupPublicKey,
				request.BlockNumber,
			)
		}

		go node.MonitorRelayEntry(
			relayChain,
			request.BlockNumber,
			chainConfig,
		)
	})

	relayChain.OnGroupSelectionStarted(func(event *event.GroupSelectionStart) {
		logger.Infof(
			"group selection started with seed [0x%v] at block [%v]",
			event.NewEntry.Text(16),
			event.BlockNumber,
		)

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
		logger.Infof(
			"new group with public key [0x%x] registered on-chain at block [%v]",
			registration.GroupPublicKey,
			registration.BlockNumber,
		)
		go groupRegistry.UnregisterStaleGroups()
	})

	return nil
}
