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

	movingFundsCommitmentData, err := mfst.findMovingFundsCommitmentData(
		walletPublicKeyHash,
		filterStartBlock,
	)
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to find committing source wallets: [%w]",
			err,
		)
	}

	movingFundsTxData, err := mfst.findMovingFundsTransactions(
		movingFundsCommitmentData,
		filterStartBlock,
	)
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to find moving funds transactions: [%w]",
			err,
		)
	}

	for _, data := range movingFundsTxData {
		movedFundsSweepRequest, err := mfst.chain.GetMovedFundsSweepRequest(
			data.TransactionHash,
			data.OutputIndex,
		)
		if err != nil {
			return bitcoin.Hash{}, 0, fmt.Errorf(
				"failed to get moved funds sweep request: [%w]",
				err,
			)
		}

		if movedFundsSweepRequest.State == tbtc.MovedFundsStatePending {
			return data.TransactionHash, data.OutputIndex, nil
		}
	}

	return bitcoin.Hash{}, 0, fmt.Errorf(
		"could not find any pending moving funds request",
	)
}

func (mfst *MovedFundsSweepTask) findMovingFundsCommitmentData(
	walletPublicKeyHash [20]byte,
	filterStartBlock uint64,
) ([]movingFundsCommitmentData, error) {
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

	commitmentData := []movingFundsCommitmentData{}
	for _, event := range events {
		for targetWalletIndex, targetWallet := range event.TargetWallets {
			if targetWallet == walletPublicKeyHash {
				commitmentData = append(
					commitmentData,
					movingFundsCommitmentData{
						committingWallet:  event.WalletPublicKeyHash,
						targetWalletIndex: uint32(targetWalletIndex),
					},
				)
				break
			}
		}
	}

	return commitmentData, nil
}

func (mfst *MovedFundsSweepTask) findMovingFundsTransactions(
	commitmentData []movingFundsCommitmentData,
	filterStartBlock uint64,
) ([]bitcoin.TransactionOutpoint, error) {
	movingFundsWallets := [][20]byte{}
	for _, data := range commitmentData {
		movingFundsWallets = append(movingFundsWallets, data.committingWallet)
	}

	filter := &tbtc.MovingFundsCompletedEventFilter{
		StartBlock:          filterStartBlock,
		WalletPublicKeyHash: movingFundsWallets,
	}

	events, err := mfst.chain.PastMovingFundsCompletedEvents(
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past moving funds completed events: [%w]",
			err,
		)
	}

	movingFundsTxData := []bitcoin.TransactionOutpoint{}

	for _, event := range events {
		movingFundsTxHash := event.MovingFundsTxHash

		movingFundsTxOutpointIdx := uint32(0)
		for _, commitment := range commitmentData {
			if commitment.committingWallet == event.WalletPublicKeyHash {
				movingFundsTxOutpointIdx = commitment.targetWalletIndex
			}
		}

		movingFundsTxData = append(
			movingFundsTxData,
			bitcoin.TransactionOutpoint{
				TransactionHash: movingFundsTxHash,
				OutputIndex:     movingFundsTxOutpointIdx,
			},
		)
	}

	return movingFundsTxData, nil
}

type movingFundsCommitmentData struct {
	committingWallet  [20]byte
	targetWalletIndex uint32
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
