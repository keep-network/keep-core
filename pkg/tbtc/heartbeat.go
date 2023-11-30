package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"time"

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
	// heartbeatRequestConfirmationBlocks determines the block length of the
	// confirmation period on the host chain that is preserved after a heartbeat
	// request submission.
	heartbeatRequestConfirmationBlocks = 3
	// heartbeatRequestTimeoutSafetyMargin determines the duration of the
	// safety margin that must be preserved between the signing timeout
	// and the timeout of the entire heartbeat action. This safety
	// margin prevents against the case where signing completes too late and
	// another action has been already requested by the coordinator.
	heartbeatRequestTimeoutSafetyMargin = 5 * time.Minute
)

type HeartbeatProposal struct {
	Message []byte
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
	) (*tecdsa.Signature, uint64, error)
}

// heartbeatAction is a walletAction implementation handling heartbeat requests
// from the wallet coordinator.
type heartbeatAction struct {
	logger           log.StandardLogger
	executingWallet  wallet
	signingExecutor  heartbeatSigningExecutor
	message          []byte
	startBlock       uint64
	requestExpiresAt time.Time
}

func newHeartbeatAction(
	logger log.StandardLogger,
	executingWallet wallet,
	signingExecutor heartbeatSigningExecutor,
	message []byte,
	startBlock uint64,
	requestExpiresAt time.Time,
) *heartbeatAction {
	return &heartbeatAction{
		logger:           logger,
		executingWallet:  executingWallet,
		signingExecutor:  signingExecutor,
		message:          message,
		startBlock:       startBlock,
		requestExpiresAt: requestExpiresAt,
	}
}

func (ha *heartbeatAction) execute() error {
	// TODO: When implementing the moving funds action we should make sure
	// heartbeats are not executed by unstaking clients.

	messageBytes := bitcoin.ComputeHash(ha.message)
	messageToSign := new(big.Int).SetBytes(messageBytes[:])

	heartbeatCtx, cancelHeartbeatCtx := context.WithTimeout(
		context.Background(),
		time.Until(ha.requestExpiresAt.Add(-heartbeatRequestTimeoutSafetyMargin)),
	)
	defer cancelHeartbeatCtx()

	signature, _, err := ha.signingExecutor.sign(heartbeatCtx, messageToSign, ha.startBlock)
	if err != nil {
		return fmt.Errorf("cannot sign heartbeat message: [%v]", err)
	}

	logger.Infof(
		"generated signature [%s] for heartbeat message [0x%x]",
		signature,
		ha.message,
	)

	return nil
}

func (ha *heartbeatAction) wallet() wallet {
	return ha.executingWallet
}

func (ha *heartbeatAction) actionType() WalletActionType {
	return ActionHeartbeat
}
