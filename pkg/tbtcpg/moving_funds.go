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

	liveWalletsCount, err := mft.chain.GetLiveWalletsCount()
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get Live wallets count: [%w]",
			err,
		)
	}

	walletMainUtxo, err := tbtc.DetermineWalletMainUtxo(
		walletPublicKeyHash,
		mft.chain,
		mft.btcChain,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get wallet's maun UTXO: [%w]",
			err,
		)
	}

	walletBalance := int64(0)
	if walletMainUtxo != nil {
		walletBalance = walletMainUtxo.Value
	}

	walletChainData, err := mft.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get source wallet's chain data: [%w]",
			err,
		)
	}

	ok, err := mft.CheckEligibility(
		taskLogger,
		liveWalletsCount,
		walletBalance,
		walletChainData,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot check wallet eligibility: [%w]",
			err,
		)
	}

	if !ok {
		taskLogger.Info("wallet not eligible for moving funds")
		return nil, false, nil
	}

	targetWallets, commitmentSubmitted, err := mft.FindTargetWallets(
		taskLogger,
		walletChainData,
	)
	if err != nil {
		return nil, false, fmt.Errorf("cannot find target wallets: [%w]", err)
	}

	proposal, err := mft.ProposeMovingFunds(
		taskLogger,
		targetWallets,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot prepare moving funds proposal: [%w]",
			err,
		)
	}

	if !commitmentSubmitted {
		mft.SubmitMovingFundsCommitment()
	}

	return proposal, false, nil
}

func (mft *MovingFundsTask) FindTargetWallets(
	taskLogger log.StandardLogger,
	walletChainData *tbtc.WalletChainData,
) ([][20]byte, bool, error) {
	if walletChainData.MovingFundsTargetWalletsCommitmentHash == [32]byte{} {
		taskLogger.Infof("Move funds commitment has not been submitted yet")
		// TODO: Find the target wallets among Live wallets.
	} else {
		taskLogger.Infof("Move funds commitment has already been submitted")
		// TODO: Find the target wallets from the emitted event.
	}
	return nil, false, nil
}

func (mft *MovingFundsTask) CheckEligibility(
	taskLogger log.StandardLogger,
	liveWalletsCount uint32,
	walletBalance int64,
	walletChainData *tbtc.WalletChainData,
) (bool, error) {
	if liveWalletsCount == 0 {
		taskLogger.Infof("there are no Live wallets available")
		return false, nil
	}

	if walletBalance <= 0 {
		// The wallet's balance cannot be `0`. Since we are dealing with
		// a signed integer we also check it's not negative just in case.
		taskLogger.Infof("source wallet does not have a positive balance")
		return false, nil
	}

	if walletChainData.State != tbtc.StateMovingFunds {
		taskLogger.Infof("source wallet not in MoveFunds state")
		return false, nil
	}

	if walletChainData.PendingRedemptionsValue > 0 {
		taskLogger.Infof("source wallet has pending redemptions")
		return false, nil
	}

	if walletChainData.PendingMovedFundsSweepRequestsCount > 0 {
		taskLogger.Infof("source wallet has pending moved funds sweep requests")
		return false, nil
	}

	return true, nil
}

func (mft *MovingFundsTask) SubmitMovingFundsCommitment() ([][20]byte, error) {
	// TODO: After submitting the transaction wait until it has at least one
	//       confirmation.
	return nil, nil
}

func (mft *MovingFundsTask) ProposeMovingFunds(
	taskLogger log.StandardLogger,
	targetWallets [][20]byte,
) (*tbtc.MovingFundsProposal, error) {
	taskLogger.Infof("preparing a moving funds proposal")

	proposal := &tbtc.MovingFundsProposal{
		TargetWallets: targetWallets,
		// TODO: Add fee
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
