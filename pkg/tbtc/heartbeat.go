package tbtc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

const (
	// heartbeatProposalValidityBlocks determines the wallet heartbeat proposal
	// validity time expressed in blocks. In other words, this is the worst-case
	// time for a wallet heartbeat during which the wallet is busy and cannot
	// take another actions. The value of 300 blocks is roughly 1 hour, assuming
	// 12 seconds per block.
	heartbeatProposalValidityBlocks = 300
	// heartbeatRequestTimeoutSafetyMarginBlocks determines the duration of the
	// safety margin that must be preserved between the signing timeout
	// and the timeout of the entire heartbeat action. This safety
	// margin prevents against the case where signing completes too late and
	// another action has been already requested by the coordinator.
	// The value of 25 blocks is roughly 5 minutes, assuming 12 seconds per block.
	heartbeatRequestTimeoutSafetyMarginBlocks = 25
)

type HeartbeatProposal struct {
	Message [16]byte
}

func (hp *HeartbeatProposal) ActionType() WalletActionType {
	return ActionHeartbeat
}

func (hp *HeartbeatProposal) ValidityBlocks() uint64 {
	return heartbeatProposalValidityBlocks
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

// heartbeatAction is a walletAction implementation handling heartbeat requests
// from the wallet coordinator.
type heartbeatAction struct {
	logger log.StandardLogger
	chain  Chain

	executingWallet wallet
	signingExecutor heartbeatSigningExecutor

	proposal    *HeartbeatProposal
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
	startBlock uint64,
	expiryBlock uint64,
	waitForBlockFn waitForBlockFn,
) *heartbeatAction {
	return &heartbeatAction{
		logger:          logger,
		chain:           chain,
		executingWallet: executingWallet,
		signingExecutor: signingExecutor,
		proposal:        proposal,
		startBlock:      startBlock,
		expiryBlock:     expiryBlock,
		waitForBlockFn:  waitForBlockFn,
	}
}

func (ha *heartbeatAction) execute() error {
	// Do not execute the heartbeat action if the operator is unstaking.
	isUnstaking, err := ha.chain.IsOperatorUnstaking()
	if err != nil {
		return fmt.Errorf("failed to check if the operator is unstaking")
	}

	if isUnstaking {
		logger.Info(
			"quitting the heartbeat action without signing because the " +
				"operator is unstaking",
		)
		return nil
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(ha.wallet().publicKey)

	err = ha.chain.ValidateHeartbeatProposal(walletPublicKeyHash, ha.proposal)
	if err != nil {
		return fmt.Errorf("heartbeat proposal is invalid: [%v]", err)
	}

	messageBytes := bitcoin.ComputeHash(ha.proposal.Message[:])
	messageToSign := new(big.Int).SetBytes(messageBytes[:])

	// Just in case. This should never happen.
	if ha.expiryBlock < heartbeatRequestTimeoutSafetyMarginBlocks {
		return fmt.Errorf("invalid proposal expiry block")
	}

	heartbeatCtx, cancelHeartbeatCtx := withCancelOnBlock(
		context.Background(),
		ha.expiryBlock-heartbeatRequestTimeoutSafetyMarginBlocks,
		ha.waitForBlockFn,
	)
	defer cancelHeartbeatCtx()

	signature, _, _, err := ha.signingExecutor.sign(
		heartbeatCtx,
		messageToSign,
		ha.startBlock,
	)
	if err != nil {
		return fmt.Errorf("cannot sign heartbeat message: [%v]", err)
	}

	logger.Infof(
		"generated signature [%s] for heartbeat message [0x%x]",
		signature,
		ha.proposal.Message[:],
	)

	return nil
}

func (ha *heartbeatAction) wallet() wallet {
	return ha.executingWallet
}

func (ha *heartbeatAction) actionType() WalletActionType {
	return ActionHeartbeat
}
