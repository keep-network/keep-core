package tbtc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/clientinfo"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/sortition"
)

// TODO: Unit tests for `tbtc.go`.

var logger = log.Logger("keep-tbtc")

// ProtocolName denotes the name of the protocol defined by this package.
const ProtocolName = "tbtc"

// GroupParameters is a structure grouping TBTC group parameters.
type GroupParameters struct {
	// GroupSize is the target size of a group in TBTC.
	GroupSize int
	// GroupQuorum is the minimum number of active participants behaving
	// according to the protocol needed to generate a group in TBTC. This value
	// is smaller than the GroupSize and bigger than the HonestThreshold.
	GroupQuorum int
	// HonestThreshold is the minimum number of active participants behaving
	// according to the protocol needed to generate a signature.
	HonestThreshold int
}

// DishonestThreshold is the maximum number of misbehaving participants for
// which it is still possible to generate a signature. Misbehaviour is any
// misconduct to the protocol, including inactivity.
func (gp *GroupParameters) DishonestThreshold() int {
	return gp.GroupSize - gp.HonestThreshold
}

const (
	DefaultPreParamsPoolSize              = 1000
	DefaultPreParamsGenerationTimeout     = 2 * time.Minute
	DefaultPreParamsGenerationDelay       = 10 * time.Second
	DefaultPreParamsGenerationConcurrency = 1
)

var DefaultKeyGenerationConcurrency = runtime.GOMAXPROCS(0)

// Config carries the config for tBTC protocol.
type Config struct {
	// The size of the pre-parameters pool for tECDSA.
	PreParamsPoolSize int
	// Timeout for pre-parameters generation for tECDSA.
	PreParamsGenerationTimeout time.Duration
	// The delay between generating new pre-params for tECDSA.
	PreParamsGenerationDelay time.Duration
	// Concurrency level for pre-parameters generation for tECDSA.
	PreParamsGenerationConcurrency int
	// Concurrency level for key-generation for tECDSA.
	KeyGenerationConcurrency int
}

// Initialize kicks off the TBTC by initializing internal state, ensuring
// preconditions like staking are met, and then kicking off the internal TBTC
// implementation. Returns an error if this failed.
func Initialize(
	ctx context.Context,
	chain Chain,
	btcChain bitcoin.Chain,
	netProvider net.Provider,
	keyStorePersistence persistence.ProtectedHandle,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
	config Config,
	clientInfo *clientinfo.Registry,
) error {
	groupParameters := &GroupParameters{
		GroupSize:       100,
		GroupQuorum:     90,
		HonestThreshold: 51,
	}

	node, err := newNode(
		groupParameters,
		chain,
		btcChain,
		netProvider,
		keyStorePersistence,
		workPersistence,
		scheduler,
		config,
	)
	if err != nil {
		return fmt.Errorf("cannot set up TBTC node: [%v]", err)
	}

	deduplicator := newDeduplicator()

	if clientInfo != nil {
		// only if client info endpoint is configured
		clientInfo.ObserveApplicationSource(
			"tbtc",
			map[string]clientinfo.Source{
				"pre_params_count": func() float64 {
					return float64(node.dkgExecutor.preParamsCount())
				},
			},
		)
	}

	err = sortition.MonitorPool(
		ctx,
		logger,
		chain,
		sortition.DefaultStatusCheckTick,
		sortition.NewConjunctionPolicy(
			sortition.NewBetaOperatorPolicy(chain, logger),
			&enoughPreParamsInPoolPolicy{
				node:   node,
				config: config,
			},
		),
	)
	if err != nil {
		return fmt.Errorf(
			"could not set up sortition pool monitoring: [%v]",
			err,
		)
	}

	_ = chain.OnDKGStarted(func(event *DKGStartedEvent) {
		go func() {
			if ok := deduplicator.notifyDKGStarted(
				event.Seed,
			); !ok {
				logger.Infof(
					"DKG started event with seed [0x%x] has been "+
						"already processed",
					event.Seed,
				)
				return
			}

			confirmationBlock := event.BlockNumber + dkgStartedConfirmationBlocks

			logger.Infof(
				"observed DKG started event with seed [0x%x] and "+
					"starting block [%v]; waiting for block [%v] to confirm",
				event.Seed,
				event.BlockNumber,
				confirmationBlock,
			)

			err := node.waitForBlockHeight(ctx, confirmationBlock)
			if err != nil {
				logger.Errorf("failed to confirm DKG started event: [%v]", err)
				return
			}

			dkgState, err := chain.GetDKGState()
			if err != nil {
				logger.Errorf("failed to check DKG state: [%v]", err)
				return
			}

			if dkgState == AwaitingResult {
				// Fetch all past DKG started events starting from one
				// confirmation period before the original event's block.
				// If there was a chain reorg, the event we received could be
				// moved to a block with a lower number than the one
				// we received.
				pastEvents, err := chain.PastDKGStartedEvents(
					&DKGStartedEventFilter{
						StartBlock: event.BlockNumber - dkgStartedConfirmationBlocks,
					},
				)
				if err != nil {
					logger.Errorf("failed to get past DKG started events: [%v]", err)
					return
				}

				// Should not happen but just in case.
				if len(pastEvents) == 0 {
					logger.Errorf("no past DKG started events")
					return
				}

				lastEvent := pastEvents[len(pastEvents)-1]

				logger.Infof(
					"DKG started with seed [0x%x] at block [%v]",
					lastEvent.Seed,
					lastEvent.BlockNumber,
				)

				// The off-chain protocol should be started as close as possible
				// to the current block or even further. Starting the off-chain
				// protocol with a past block will likely cause a failure of the
				// first attempt as the start block is used to synchronize
				// the announcements and the state machine. Here we ensure
				// a proper start point by delaying the execution by the
				// confirmation period length.
				node.joinDKGIfEligible(
					lastEvent.Seed,
					lastEvent.BlockNumber,
					dkgStartedConfirmationBlocks,
				)
			} else {
				logger.Infof(
					"DKG started event with seed [0x%x] and starting "+
						"block [%v] was not confirmed",
					event.Seed,
					event.BlockNumber,
				)
			}
		}()
	})

	_ = chain.OnDKGResultSubmitted(func(event *DKGResultSubmittedEvent) {
		go func() {
			if ok := deduplicator.notifyDKGResultSubmitted(
				event.Seed,
				event.ResultHash,
				event.BlockNumber,
			); !ok {
				logger.Warnf(
					"Result with hash [0x%x] for DKG with seed [0x%x] "+
						"and starting block [%v] has been already processed",
					event.ResultHash,
					event.Seed,
					event.BlockNumber,
				)
				return
			}

			logger.Infof(
				"Result with hash [0x%x] for DKG with seed [0x%x] "+
					"submitted at block [%v]",
				event.ResultHash,
				event.Seed,
				event.BlockNumber,
			)

			node.validateDKG(
				event.Seed,
				event.BlockNumber,
				event.Result,
				event.ResultHash,
			)
		}()
	})

	// Set up a handler of a wallet heartbeat request coming from the
	// WalletCoordinator on-chain contract. Once an event is seen, a handler
	// goroutine makes sure that the observed event is not duplicate, waits
	// a fixed confirmation period, and ensures the on-chain state justifies
	// the occurrence of the event. Once done, the original event is used
	// to trigger the heartbeat request action. The handler does not care about
	// possible subsequent events being a result of chain reorgs. This is because
	// the WalletCoordinator contract is just a coordination point based on
	// the chain consensus. If enough clients received the event, they should
	// follow it and execute a heartbeat signing. The message for that heartbeat
	// request was validated in the contract, so even if there was a reorg and
	// another event landed on the canonical chain later, the signature
	// will still be valid and not lead to fraud. The only reason the handler
	// waits a fixed confirmation period after receiving the coordination event
	// is to make sure the right type of action is executed given different types
	// of actions may have different lock times. We do not want to run into a
	// situation when the majority of clients execute heartbeat with N blocks
	// wallet lock time and the chain has M < N blocks wallet lock time because
	// the canonical chain - as a result of a reorg - is supposed to execute
	// e.g. redemption.
	_ = chain.OnHeartbeatRequestSubmitted(func(event *HeartbeatRequestSubmittedEvent) {
		go func() {
			walletPublicKeyHash := event.WalletPublicKeyHash
			message := event.Message

			if ok := deduplicator.notifyHeartbeatRequestSubmitted(
				walletPublicKeyHash,
				message,
			); !ok {
				logger.Infof(
					"heartbeat request for wallet PKH [0x%x] and message [0x%x] "+
						"has been already processed",
					walletPublicKeyHash,
					message,
				)
				return
			}

			confirmationBlock := event.BlockNumber +
				heartbeatRequestConfirmationBlocks

			logger.Infof(
				"observed heartbeat request for wallet PKH [0x%x] "+
					"at block [%v]; waiting for block [%v] to confirm",
				walletPublicKeyHash,
				event.BlockNumber,
				confirmationBlock,
			)

			err := node.waitForBlockHeight(ctx, confirmationBlock)
			if err != nil {
				logger.Errorf(
					"failed to confirm heartbeat request for "+
						"wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			expiresAt, cause, err := chain.GetWalletLock(
				walletPublicKeyHash,
			)
			if err != nil {
				logger.Errorf(
					"failed to get lock for wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			// The event is confirmed if the wallet is locked due to a heartbeat
			// action.
			if time.Now().Before(expiresAt) && cause == ActionHeartbeat {
				logger.Infof(
					"heartbeat request submitted for "+
						"wallet PKH [0x%x] at block [%v] by [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					event.Coordinator,
				)

				node.handleHeartbeatRequest(
					walletPublicKeyHash,
					message,
					expiresAt,
					event.BlockNumber,
					heartbeatRequestConfirmationBlocks,
				)
			} else {
				logger.Infof(
					"heartbeat request for wallet PKH [0x%x] "+
						"at block [%v] was not confirmed; existing wallet lock "+
						"has unexpected expiration time [%s] and/or cause [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					expiresAt,
					cause,
				)
			}
		}()
	})

	// Set up a handler of deposit sweep proposals coming from the
	// WalletCoordinator on-chain contract. Once an event is seen, a handler
	// goroutine makes sure that the observed event is not a duplicate, waits
	// a fixed confirmation period, and ensures the on-chain state justifies
	// the occurrence of the event. Once done, the original event is used
	// to trigger the deposit sweep action. The handler does not care about
	// possible subsequent events being result of chain reorgs. This is because
	// the WalletCoordinator contract is just a coordination point based on
	// the chain consensus. If enough clients received the event, they should
	// follow it and execute a signature. All input parameters for that
	// signature are validated, so even if there was a reorg and another event
	// landed on the canonical chain later, the first signature will still be
	// valid and approved by Bitcoin. The only reason the handler waits a
	// fixed confirmation period after receiving the coordination event is to
	// make sure the right type of action is executed given different types of
	// actions may have different lock times. We do not want to run into a
	// situation when the majority of clients execute sweep with N blocks wallet
	// lock time and the chain has M < N blocks wallet lock time because the
	// canonical chain - as a result of a reorg - is supposed to execute
	// e.g. redemption.
	_ = chain.OnDepositSweepProposalSubmitted(func(event *DepositSweepProposalSubmittedEvent) {
		go func() {
			walletPublicKeyHash := event.Proposal.WalletPublicKeyHash

			if ok := deduplicator.notifyDepositSweepProposalSubmitted(
				event.Proposal,
			); !ok {
				logger.Infof(
					"deposit sweep proposal for wallet PKH [0x%x] "+
						"has been already processed",
					walletPublicKeyHash,
				)
				return
			}

			confirmationBlock := event.BlockNumber +
				depositSweepProposalConfirmationBlocks

			logger.Infof(
				"observed deposit sweep proposal for wallet PKH [0x%x] "+
					"at block [%v]; waiting for block [%v] to confirm",
				walletPublicKeyHash,
				event.BlockNumber,
				confirmationBlock,
			)

			err := node.waitForBlockHeight(ctx, confirmationBlock)
			if err != nil {
				logger.Errorf(
					"failed to confirm deposit sweep proposal for "+
						"wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			expiresAt, cause, err := chain.GetWalletLock(
				walletPublicKeyHash,
			)
			if err != nil {
				logger.Errorf(
					"failed to get lock for wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			// The event is confirmed if the wallet is locked due to a deposit
			// sweep action.
			if time.Now().Before(expiresAt) && cause == ActionDepositSweep {
				logger.Infof(
					"deposit sweep proposal submitted for "+
						"wallet PKH [0x%x] at block [%v] by [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					event.Coordinator,
				)

				node.handleDepositSweepProposal(
					event.Proposal,
					expiresAt,
					event.BlockNumber,
					depositSweepProposalConfirmationBlocks,
				)
			} else {
				logger.Infof(
					"deposit sweep proposal for wallet PKH [0x%x] "+
						"at block [%v] was not confirmed; existing wallet lock "+
						"has unexpected expiration time [%s] and/or cause [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					expiresAt,
					cause,
				)
			}
		}()
	})

	// Set up a handler of redemption proposals coming from the
	// WalletCoordinator on-chain contract. Once an event is seen, a handler
	// goroutine makes sure that the observed event is not a duplicate, waits
	// a fixed confirmation period, and ensures the on-chain state justifies
	// the occurrence of the event. Once done, the original event is used
	// to trigger the redemption action. The handler does not care about
	// possible subsequent events being result of chain reorgs. This is because
	// the WalletCoordinator contract is just a coordination point based on
	// the chain consensus. If enough clients received the event, they should
	// follow it and execute a signature. All input parameters for that
	// signature are validated, so even if there was a reorg and another event
	// landed on the canonical chain later, the first signature will still be
	// valid and approved by Bitcoin. The only reason the handler waits a
	// fixed confirmation period after receiving the coordination event is to
	// make sure the right type of action is executed given different types of
	// actions may have different lock times. We do not want to run into a
	// situation when the majority of clients execute redemption with N blocks
	// wallet lock time and the chain has M < N blocks wallet lock time because
	// the canonical chain - as a result of a reorg - is supposed to execute
	// e.g. wallet heartbeat.
	_ = chain.OnRedemptionProposalSubmitted(func(event *RedemptionProposalSubmittedEvent) {
		go func() {
			walletPublicKeyHash := event.Proposal.WalletPublicKeyHash

			if ok := deduplicator.notifyRedemptionProposalSubmitted(
				event.Proposal,
			); !ok {
				logger.Infof(
					"redemption proposal for wallet PKH [0x%x] "+
						"has been already processed",
					walletPublicKeyHash,
				)
				return
			}

			confirmationBlock := event.BlockNumber +
				redemptionProposalConfirmationBlocks

			logger.Infof(
				"observed redemption proposal for wallet PKH [0x%x] "+
					"at block [%v]; waiting for block [%v] to confirm",
				walletPublicKeyHash,
				event.BlockNumber,
				confirmationBlock,
			)

			err := node.waitForBlockHeight(ctx, confirmationBlock)
			if err != nil {
				logger.Errorf(
					"failed to confirm redemption proposal for "+
						"wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			expiresAt, cause, err := chain.GetWalletLock(
				walletPublicKeyHash,
			)
			if err != nil {
				logger.Errorf(
					"failed to get lock for wallet PKH [0x%x]: [%v]",
					walletPublicKeyHash,
					err,
				)
				return
			}

			// The event is confirmed if the wallet is locked due to a
			// redemption action.
			if time.Now().Before(expiresAt) && cause == ActionRedemption {
				logger.Infof(
					"redemption proposal submitted for "+
						"wallet PKH [0x%x] at block [%v] by [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					event.Coordinator,
				)

				node.handleRedemptionProposal(
					event.Proposal,
					expiresAt,
					event.BlockNumber,
					redemptionProposalConfirmationBlocks,
				)
			} else {
				logger.Infof(
					"redemption proposal for wallet PKH [0x%x] "+
						"at block [%v] was not confirmed; existing wallet lock "+
						"has unexpected expiration time [%s] and/or cause [%v]",
					walletPublicKeyHash,
					event.BlockNumber,
					expiresAt,
					cause,
				)
			}
		}()
	})

	return nil
}

// enoughPreParamsInPoolPolicy is a policy that enforces the sufficient size
// of the DKG pre-parameters pool before joining the sortition pool.
type enoughPreParamsInPoolPolicy struct {
	node   *node
	config Config
}

func (eppip *enoughPreParamsInPoolPolicy) ShouldJoin() bool {
	paramsInPool := eppip.node.dkgExecutor.preParamsCount()
	poolSize := eppip.config.PreParamsPoolSize
	return paramsInPool >= poolSize
}
