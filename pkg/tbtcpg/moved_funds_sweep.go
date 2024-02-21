package tbtcpg

import (
	"fmt"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type MovedFundsSweepTask struct {
	chain    Chain
	btcChain bitcoin.Chain
}

func NewMovedFundsSweepTask(
	chain Chain,
	btcChain bitcoin.Chain,
) *MovedFundsSweepTask {
	return &MovedFundsSweepTask{
		chain:    chain,
		btcChain: btcChain,
	}
}

func (mfst *MovedFundsSweepTask) Run(request *tbtc.CoordinationProposalRequest) (
	tbtc.CoordinationProposal,
	bool,
	error,
) {
	walletPublicKeyHash := request.WalletPublicKeyHash

	taskLogger := logger.With(
		zap.String("task", mfst.ActionType().String()),
		zap.String("walletPKH", fmt.Sprintf("0x%x", walletPublicKeyHash)),
	)

	walletChainData, err := mfst.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get wallet's chain data: [%w]",
			err,
		)
	}

	if walletChainData.State != tbtc.StateLive &&
		walletChainData.State != tbtc.StateMovingFunds {
		taskLogger.Infof("wallet not in Live or MoveFunds state")
		return nil, false, nil
	}

	if walletChainData.PendingMovedFundsSweepRequestsCount == 0 {
		taskLogger.Infof("wallet has no pending moved funds sweep requests")
		return nil, false, nil
	}

	movingFundsTxHash, movingFundsOutputIdx, err := mfst.FindMovingFundsTxData(
		walletPublicKeyHash,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot find moving funds transaction data: [%w]",
			err,
		)
	}

	proposal, err := mfst.ProposeMovedFundsSweep(
		taskLogger,
		walletPublicKeyHash,
		movingFundsTxHash,
		movingFundsOutputIdx,
		0,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot prepare moved funds sweep proposal: [%w]",
			err,
		)
	}

	return proposal, true, nil
}

func (mfst *MovedFundsSweepTask) FindMovingFundsTxData(
	walletPublicKeyHash [20]byte,
) (
	txHash bitcoin.Hash,
	txOutputIdx uint32,
	err error,
) {
	blockCounter, err := mfst.chain.BlockCounter()
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to get block counter: [%w]",
			err,
		)
	}

	currentBlockNumber, err := blockCounter.CurrentBlock()
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to get current block number: [%w]",
			err,
		)
	}

	filterStartBlock := uint64(0)
	if currentBlockNumber > MovingFundsCommitmentLookBackBlocks {
		filterStartBlock = currentBlockNumber - MovingFundsCommitmentLookBackBlocks
	}

	_, err = mfst.findCommittingWallets(
		walletPublicKeyHash,
		filterStartBlock,
	)
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to find committing source wallets : [%w]",
			err,
		)
	}

	// TODO: Continue with implementation
	return bitcoin.Hash{}, 0, nil
}

func (mfst *MovedFundsSweepTask) findCommittingWallets(
	walletPublicKeyHash [20]byte,
	filterStartBlock uint64,
) ([][20]byte, error) {
	filter := &tbtc.MovingFundsCommitmentSubmittedEventFilter{
		StartBlock: filterStartBlock,
	}

	events, err := mfst.chain.PastMovingFundsCommitmentSubmittedEvents(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past moving funds commitment submitted events: [%w]",
			err,
		)
	}

	committingWallets := [][20]byte{}
	for _, event := range events {
		for _, targetWallet := range event.TargetWallets {
			if targetWallet == walletPublicKeyHash {
				committingWallets = append(
					committingWallets,
					event.WalletPublicKeyHash,
				)
			}
		}
	}

	return committingWallets, nil
}

func (mfst *MovedFundsSweepTask) ActionType() tbtc.WalletActionType {
	return tbtc.ActionMovedFundsSweep
}

func (mfst *MovedFundsSweepTask) ProposeMovedFundsSweep(
	taskLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	movingFundsTxHash bitcoin.Hash,
	movingFundsTxIndex uint32,
	fee int64,
) (*tbtc.MovedFundsSweepProposal, error) {
	// TODO: Implement
	return nil, nil
}
