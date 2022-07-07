package beacon

import (
	"encoding/hex"
	"time"

	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/persistence"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/beacon/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-beacon")

// Initialize kicks off the random beacon by initializing internal state,
// ensuring preconditions like staking are met, and then kicking off the
// internal random beacon implementation. Returns an error if this failed,
// otherwise enters a blocked loop.
func Initialize(
	operatorPublicKey *operator.PublicKey,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	signing chain.Signing,
	netProvider net.Provider,
	persistence persistence.Handle,
) error {
	chainConfig := relayChain.GetConfig()

	groupRegistry := registry.NewGroupRegistry(relayChain, persistence)
	groupRegistry.LoadExistingGroups()

	node := NewNode(
		operatorPublicKey,
		netProvider,
		blockCounter,
		chainConfig,
		groupRegistry,
	)

	eventDeduplicator := event.NewDeduplicator(relayChain)

	node.ResumeSigningIfEligible(relayChain, signing)

	_ = relayChain.OnRelayEntryRequested(func(request *event.Request) {
		onConfirmed := func() {
			if node.IsInGroup(request.GroupPublicKey) {
				go func() {
					shouldProcess, err := eventDeduplicator.NotifyRelayEntryStarted(
						request.BlockNumber,
						hex.EncodeToString(request.PreviousEntry[:]),
					)
					if err != nil {
						logger.Errorf(
							"could not determine whether relay entry "+
								"requested event with previous entry [0x%x] "+
								"and starting block [%v] is a duplicate: [%v]",
							request.PreviousEntry,
							request.BlockNumber,
							err,
						)
						return
					}

					if !shouldProcess {
						logger.Warningf(
							"relay entry requested event with previous "+
								"entry [0x%x] and starting block [%v] has been "+
								"already processed",
							request.PreviousEntry,
							request.BlockNumber,
						)
						return
					}

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
		}

		currentRelayRequestConfirmationRetries := 30
		currentRelayRequestConfirmationDelay := time.Second

		confirmCurrentRelayRequest(
			request.BlockNumber,
			relayChain,
			onConfirmed,
			currentRelayRequestConfirmationRetries,
			currentRelayRequestConfirmationDelay,
		)
	})

	_ = relayChain.OnDKGStarted(func(event *event.DKGStarted) {
		go func() {
			if ok := eventDeduplicator.NotifyDKGStarted(
				event.Seed,
			); !ok {
				logger.Warningf(
					"DKG started event with seed [0x%x] and "+
						"starting block [%v] has been already processed",
					event.Seed,
					event.BlockNumber,
				)
				return
			}

			logger.Infof(
				"DKG started with seed [0x%x] at block [%v]",
				event.Seed,
				event.BlockNumber,
			)

			node.JoinDKGIfEligible(
				relayChain,
				signing,
				event.Seed,
				event.BlockNumber,
			)
		}()
	})

	_ = relayChain.OnGroupRegistered(func(registration *event.GroupRegistration) {
		logger.Infof(
			"new group with public key [0x%x] registered on-chain at block [%v]",
			registration.GroupPublicKey,
			registration.BlockNumber,
		)
		go groupRegistry.UnregisterStaleGroups(registration.GroupPublicKey)
	})

	return nil
}

// Before we start relay entry signing process we need to confirm the current
// relay request start block on the chain. This is to avoid having the client
// participating in an old relay request signing that has already completed
// with the rest of the signing group member clients and the result has been
// already published to the chain.
//
// Such situation may happen when the current client received multiple blocks
// at once after a longer delay and in those blocks to relay request events
// to the same signing group were emitted.
//
// The confirmation mechanism has built-in retries. We can retry in case of an
// error but also when the expected request start block does not match the one
// currently registered on the chain. Such situation may happen for Infura-like
// setup when two or more chain clients are behind a load balancer and they do
// not have their state in sync yet.
func confirmCurrentRelayRequest(
	expectedRequestStartBlock uint64,
	chain relaychain.RelayEntryInterface,
	onConfirmed func(),
	maxRetries int,
	delay time.Duration,
) {
	for i := 1; ; i++ {
		currentRequestStartBlockBigInt, err := chain.CurrentRequestStartBlock()
		if err != nil {
			if i == maxRetries {
				logger.Errorf(
					"could not check current request start block: [%v]; "+
						"giving up after [%v] retries",
					err,
					maxRetries,
				)
				return
			}

			logger.Warningf(
				"could not check current request start block: [%v]; "+
					"will retry after [%v]",
				err,
				delay,
			)
			time.Sleep(delay)
			continue
		}

		currentRequestStartBlock := currentRequestStartBlockBigInt.Uint64()

		if currentRequestStartBlock == expectedRequestStartBlock {
			onConfirmed()
			return
		} else if currentRequestStartBlock > expectedRequestStartBlock {
			logger.Infof(
				"the currently pending relay request started at block [%v]; "+
					"skipping the execution of the old relay request from block [%v]",
				currentRequestStartBlock,
				expectedRequestStartBlock,
			)
			return
		} else if i == maxRetries {
			// This scenario usually happens when an entry was submitted very
			// fast before this node receives an event and is able to confirm a
			// request ID.
			if currentRequestStartBlock == 0 {
				logger.Warningf(
					"there is no entry in progress; "+
						"current request start block is 0 "+
						"giving up after [%v] retries",
					maxRetries,
				)
			} else {
				logger.Errorf(
					"could not confirm the expected relay request starting block; "+
						"the most recent one obtained from chain is [%v] and the "+
						"expected one is [%v]; giving up after [%v] retries",
					currentRequestStartBlock,
					expectedRequestStartBlock,
					maxRetries,
				)
			}
			return
		} else {
			logger.Infof(
				"received unexpected pending relay request start block [%v] "+
					"while the expected was [%v]; will retry after [%v]",
				currentRequestStartBlock,
				expectedRequestStartBlock,
				delay,
			)
			time.Sleep(delay)
		}
	}
}
