package tbtcpg

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// MovingFundsTask is a task that may produce a moving funds proposal.
type MovingFundsTask struct {
	chain    Chain
	btcChain bitcoin.Chain
}

func NewMovingFundsTask(
	chain Chain,
	btcChain bitcoin.Chain,
) *MovingFundsTask {
	return &MovingFundsTask{
		chain:    chain,
		btcChain: btcChain,
	}
}

func (mft *MovingFundsTask) Run(request *tbtc.CoordinationProposalRequest) (
	tbtc.CoordinationProposal,
	bool,
	error,
) {
	walletPublicKeyHash := request.WalletPublicKeyHash

	taskLogger := logger.With(
		zap.String("task", mft.ActionType().String()),
		zap.String("walletPKH", fmt.Sprintf("0x%x", walletPublicKeyHash)),
	)

	proposal, err := mft.ProposeMovingFunds(
		taskLogger,
		walletPublicKeyHash,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot prepare moving funds proposal: [%w]",
			err,
		)
	}

	return proposal, false, nil
}

func (mft *MovingFundsTask) ProposeMovingFunds(
	taskLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
) (*tbtc.MovingFundsProposal, error) {
	taskLogger.Infof("preparing a moving funds proposal")

	proposal := &tbtc.MovingFundsProposal{
		WalletPublicKeyHash: walletPublicKeyHash,
	}

	taskLogger.Infof("validating the moving funds proposal")
	if _, err := tbtc.ValidateMovingFundsProposal(
		taskLogger,
		proposal,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to verify moving funds proposal: %v",
			err,
		)
	}

	return proposal, nil
}

func (mft *MovingFundsTask) ActionType() tbtc.WalletActionType {
	return tbtc.ActionMovingFunds
}
