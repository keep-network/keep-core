package tbtc

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

const (
	// heartbeatTotalProposalValidityBlocks determines the total wallet
	// heartbeat proposal validity time expressed in blocks. In other words,
	// this is the worst-case time for a wallet heartbeat during which the
	// wallet is busy and cannot take another actions. It includes the total
	// duration needed to perform both signing the heartbeat message and
	// optionally notifying about operator inactivity if the heartbeat failed.
	// The value of 600 blocks is roughly 2 hours, assuming 12 seconds per block.
	heartbeatTotalProposalValidityBlocks = 600
	// heartbeatInactivityClaimValidityBlocks determines the duration that needs
	// to be preserved for the optional notification about operator inactivity
	// that follows a failed heartbeat signing.
	heartbeatInactivityClaimValidityBlocks = 300
	// heartbeatTimeoutSafetyMarginBlocks determines the duration of the safety
	// margin that must be preserved between the timeout of operator inactivity
	// notification and the timeout of the entire heartbeat action. This safety
	// margin prevents against the case where signing completes too late and
	// another action has been already requested by the coordinator. The value
	// of 25 blocks is roughly 5 minutes, assuming 12 seconds per block.
	heartbeatTimeoutSafetyMarginBlocks = 25
	// heartbeatSigningMinimumActiveOperators determines the minimum number of
	// active operators during signing for a heartbeat to be considered valid.
	heartbeatSigningMinimumActiveOperators = 70
	// heartbeatConsecutiveFailuresThreshold determines the number of consecutive
	// heartbeat failures required to trigger inactivity operator notification.
	heartbeatConsecutiveFailureThreshold = 3
)

type HeartbeatProposal struct {
	Message [16]byte
}

func (hp *HeartbeatProposal) ActionType() WalletActionType {
	return ActionHeartbeat
}

func (hp *HeartbeatProposal) ValidityBlocks() uint64 {
	return heartbeatTotalProposalValidityBlocks
}

// heartbeatSigningExecutor is an interface meant to decouple the specific
// implementation of the signing executor from the heartbeat action.
type heartbeatSigningExecutor interface {
	sign(
		ctx context.Context,
		message *big.Int,
		startBlock uint64,
	) (*tecdsa.Signature, uint32, uint64, error)
}

// heartbeatInactivityClaimExecutor is an interface meant to decouple the
// specific implementation of the inactivity claim executor from the heartbeat
// action.
type heartbeatInactivityClaimExecutor interface {
	claimInactivity(
		ctx context.Context,
		inactiveMembersIndexes []group.MemberIndex,
		heartbeatFailed bool,
		sessionID *big.Int,
	) error
}

// heartbeatAction is a walletAction implementation handling heartbeat requests
// from the wallet coordinator.
type heartbeatAction struct {
	logger log.StandardLogger
	chain  Chain

	executingWallet wallet
	signingExecutor heartbeatSigningExecutor

	proposal       *HeartbeatProposal
	failureCounter *heartbeatFailureCounter

	inactivityClaimExecutor heartbeatInactivityClaimExecutor

	startBlock  uint64
	expiryBlock uint64

	waitForBlockFn waitForBlockFn
}

func newHeartbeatAction(
	logger log.StandardLogger,
	chain Chain,
	executingWallet wallet,
	signingExecutor heartbeatSigningExecutor,
	proposal *HeartbeatProposal,
	failureCounter *heartbeatFailureCounter,
	inactivityClaimExecutor heartbeatInactivityClaimExecutor,
	startBlock uint64,
	expiryBlock uint64,
	waitForBlockFn waitForBlockFn,
) *heartbeatAction {
	return &heartbeatAction{
		logger:                  logger,
		chain:                   chain,
		executingWallet:         executingWallet,
		signingExecutor:         signingExecutor,
		proposal:                proposal,
		failureCounter:          failureCounter,
		inactivityClaimExecutor: inactivityClaimExecutor,
		startBlock:              startBlock,
		expiryBlock:             expiryBlock,
		waitForBlockFn:          waitForBlockFn,
	}
}

func (ha *heartbeatAction) execute() error {
	// Do not execute the heartbeat action if the operator is unstaking.
	isUnstaking, err := ha.chain.IsOperatorUnstaking()
	if err != nil {
		return fmt.Errorf("failed to check if the operator is unstaking")
	}

	if isUnstaking {
		ha.logger.Warn(
			"quitting the heartbeat action without signing because the " +
				"operator is unstaking",
		)
		return nil
	}

	walletPublicKey := ha.wallet().publicKey
	walletPublicKeyHash := bitcoin.PublicKeyHash(walletPublicKey)
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		return fmt.Errorf("failed to unmarshal wallet public key: [%v]", err)
	}

	walletKey := hex.EncodeToString(walletPublicKeyBytes)

	err = ha.chain.ValidateHeartbeatProposal(walletPublicKeyHash, ha.proposal)
	if err != nil {
		return fmt.Errorf("heartbeat proposal is invalid: [%v]", err)
	}

	messageBytes := bitcoin.ComputeHash(ha.proposal.Message[:])
	messageToSign := new(big.Int).SetBytes(messageBytes[:])

	// Just in case. This should never happen.
	if ha.expiryBlock < heartbeatInactivityClaimValidityBlocks {
		return fmt.Errorf("invalid proposal expiry block")
	}

	heartbeatSigningCtx, cancelHeartbeatSigningCtx := withCancelOnBlock(
		context.Background(),
		ha.expiryBlock-heartbeatInactivityClaimValidityBlocks,
		ha.waitForBlockFn,
	)
	defer cancelHeartbeatSigningCtx()

	signature, activeOperatorsCount, _, err := ha.signingExecutor.sign(
		heartbeatSigningCtx,
		messageToSign,
		ha.startBlock,
	)

	// If there was no error and the number of active operators during signing
	// was enough, we can consider the heartbeat procedure as successful.
	if err == nil && activeOperatorsCount >= heartbeatSigningMinimumActiveOperators {
		ha.logger.Infof(
			"successfully generated signature [%s] for heartbeat message [0x%x]",
			signature,
			ha.proposal.Message[:],
		)

		// Reset the counter for consecutive heartbeat failure.
		ha.failureCounter.reset(walletKey)

		return nil
	}

	// If there was an error or the number of active operators during signing
	// was not enough, we must consider the heartbeat procedure as a failure.
	ha.logger.Warnf(
		"heartbeat failed; [%d/%d] operators participated; the process "+
			"returned [%v] as error",
		activeOperatorsCount,
		heartbeatSigningMinimumActiveOperators,
		err,
	)

	// Increment the heartbeat failure counter.
	ha.failureCounter.increment(walletKey)

	// If the number of consecutive heartbeat failures does not exceed the
	// threshold do not notify about operator inactivity.
	if ha.failureCounter.get(walletKey) < heartbeatConsecutiveFailureThreshold {
		ha.logger.Warnf(
			"leaving without notifying about operator inactivity; current "+
				"heartbeat failure count is [%d]",
			ha.failureCounter.get(walletKey),
		)
		return nil
	}

	heartbeatInactivityCtx, cancelHeartbeatInactivityCtx := withCancelOnBlock(
		context.Background(),
		ha.expiryBlock-heartbeatTimeoutSafetyMarginBlocks,
		ha.waitForBlockFn,
	)
	defer cancelHeartbeatInactivityCtx()

	// The value of consecutive heartbeat failures exceeds the threshold.
	// Proceed with operator inactivity notification.
	err = ha.inactivityClaimExecutor.claimInactivity(
		heartbeatInactivityCtx,
		// Leave the list of inactive operators empty even if some operators
		// were inactive during signing heartbeat. The inactive operators could
		// simply be in the process of unstaking and therefore should not be
		// punished.
		[]group.MemberIndex{},
		true,
		messageToSign,
	)
	if err != nil {
		return fmt.Errorf(
			"error while notifying about operator inactivity [%v]]",
			err,
		)
	}

	return nil
}

func (ha *heartbeatAction) wallet() wallet {
	return ha.executingWallet
}

func (ha *heartbeatAction) actionType() WalletActionType {
	return ActionHeartbeat
}

// heartbeatFailureCounter holds counters keeping track of consecutive
// heartbeat failures. Each wallet has a separate counter. The key used in
// the map is the uncompressed public key (with 04 prefix) of the wallet.
type heartbeatFailureCounter struct {
	mutex    sync.Mutex
	counters map[string]uint
}

func newHeartbeatFailureCounter() *heartbeatFailureCounter {
	return &heartbeatFailureCounter{
		counters: make(map[string]uint),
	}
}

func (hfc *heartbeatFailureCounter) increment(walletPublicKey string) {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	hfc.counters[walletPublicKey]++

}

func (hfc *heartbeatFailureCounter) reset(walletPublicKey string) {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	hfc.counters[walletPublicKey] = 0
}

func (hfc *heartbeatFailureCounter) get(walletPublicKey string) uint {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	return hfc.counters[walletPublicKey]
}
