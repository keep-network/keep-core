package tbtc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log/v2"
)

const (
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

// heartbeatAction is a walletAction implementation handling heartbeat requests
// from the wallet coordinator.
type heartbeatAction struct {
	logger           log.StandardLogger
	executingWallet  wallet
	signingExecutor  *signingExecutor
	message          []byte
	startBlock       uint64
	requestExpiresAt time.Time
}

func newHeartbeatAction(
	logger log.StandardLogger,
	executingWallet wallet,
	signingExecutor *signingExecutor,
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
	preimageSha256 := sha256.Sum256(ha.message)
	messageBytes := sha256.Sum256(preimageSha256[:])
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
	return Heartbeat
}
