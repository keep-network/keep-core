package tbtcpg

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// ErrSweepTxFeeTooHigh is the error returned when the estimated fee exceeds the
// maximum fee allowed for the moved funds sweep transaction.
var ErrSweepTxFeeTooHigh = fmt.Errorf(
	"estimated fee exceeds the maximum fee",
)

// ErrNoPendingMovedFundsSweepRequests is the error returned when no pending
// moved funds sweep requests.
var ErrNoPendingMovedFundsSweepRequests = fmt.Errorf(
	"could not find any pending moved funds sweep request",
)

// MovedFundsSweepLookBackBlocks is the look-back period in blocks used
// when searching for submitted moving funds-related events. It's equal to
// 30 days assuming 12 seconds per block.
const MovedFundsSweepLookBackBlocks = uint64(216000)

// MovedFundsSweepTask is a task that may produce a moved funds sweep proposal.
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

	// Check if the wallet is eligible for moved funds sweep.
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

	// Find the transaction hash and output index of the UTXO that the source
	// wallet transferred to this wallet. If the wallet has more than one
	// pending moved funds sweep requests, only the data for the first
	// encountered request will be returned.
	movingFundsTxHash, movingFundsOutputIdx, err := mfst.FindMovingFundsTxData(
		walletPublicKeyHash,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot find moving funds transaction data: [%w]",
			err,
		)
	}

	// Check if the wallet already has a main UTXO.
	hasMainUtxo := walletChainData.MainUtxoHash != [32]byte{}

	proposal, err := mfst.ProposeMovedFundsSweep(
		taskLogger,
		walletPublicKeyHash,
		movingFundsTxHash,
		movingFundsOutputIdx,
		hasMainUtxo,
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

// FindMovingFundsTxData finds the transaction hash and output index for the
// unswept funds transferred from a source wallet to this wallet. It returns
// data for only one funds transfer. If there were more than one funds transfers
// the data for the first encountered unswept transfer is returned.
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
	if currentBlockNumber > MovedFundsSweepLookBackBlocks {
		filterStartBlock = currentBlockNumber - MovedFundsSweepLookBackBlocks
	}

	// Find all the wallets that recently committed to move funds to this wallet.
	// The returned data also contains information on the position on the target
	// wallets list.
	movingFundsCommitments, err := mfst.findMovingFundsCommitments(
		walletPublicKeyHash,
		filterStartBlock,
	)
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to find committing source wallets: [%w]",
			err,
		)
	}

	// Find the transaction outpoints that transfer funds to this wallet coming
	// from the recent moving funds transactions.
	movingFundsTxOutpoints, err := mfst.findMovingFundsTxOutpoints(
		movingFundsCommitments,
		filterStartBlock,
	)
	if err != nil {
		return bitcoin.Hash{}, 0, fmt.Errorf(
			"failed to find moving funds transactions: [%w]",
			err,
		)
	}

	// Find the first outpoint that represents an unproven moved funds sweep
	// request and return it.
	for _, outpoint := range movingFundsTxOutpoints {
		movedFundsSweepRequest, err := mfst.chain.GetMovedFundsSweepRequest(
			outpoint.TransactionHash,
			outpoint.OutputIndex,
		)
		if err != nil {
			return bitcoin.Hash{}, 0, fmt.Errorf(
				"failed to get moved funds sweep request: [%w]",
				err,
			)
		}

		if movedFundsSweepRequest.State == tbtc.MovedFundsStatePending {
			return outpoint.TransactionHash, outpoint.OutputIndex, nil
		}
	}

	// No pending moved funds sweep request could be found.
	return bitcoin.Hash{}, 0, ErrNoPendingMovedFundsSweepRequests
}

func (mfst *MovedFundsSweepTask) findMovingFundsCommitments(
	targetWalletPublicKeyHash [20]byte,
	filterStartBlock uint64,
) ([]movingFundsCommitmentData, error) {
	// Get all the recent moving funds commitment submitted events.
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

	// Iterate over the events and take data on the wallets that committed to
	// this wallet.
	commitmentData := []movingFundsCommitmentData{}
	for _, event := range events {
		for targetWalletIndex, targetWallet := range event.TargetWallets {
			if targetWallet == targetWalletPublicKeyHash {
				// If our wallet has been found on the target wallet list save
				// the committing wallet public key hash and the position of
				// our wallet on the target wallet list.
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

func (mfst *MovedFundsSweepTask) findMovingFundsTxOutpoints(
	commitmentData []movingFundsCommitmentData,
	filterStartBlock uint64,
) ([]bitcoin.TransactionOutpoint, error) {
	movingFundsWallets := [][20]byte{}
	for _, data := range commitmentData {
		movingFundsWallets = append(movingFundsWallets, data.committingWallet)
	}

	// Use all the wallets that committed to move funds to our wallet as a filter.
	filter := &tbtc.MovingFundsCompletedEventFilter{
		StartBlock:          filterStartBlock,
		WalletPublicKeyHash: movingFundsWallets,
	}

	// Retrieve all the moving funds completed events. The list will contain
	// information on all the moving funds transactions to our wallet that have
	// been proven in the Bridge.
	events, err := mfst.chain.PastMovingFundsCompletedEvents(
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past moving funds completed events: [%w]",
			err,
		)
	}

	movingFundsTxOutpoints := []bitcoin.TransactionOutpoint{}

	// Iterate over all the events and prepare outpoints of Bitcoin transactions
	// that transferred funds to our wallets. These outpoints represent moved
	// funds sweep requests. Some of them are already processed, some are still
	// pending.
	for _, event := range events {
		movingFundsTxHash := event.MovingFundsTxHash

		movingFundsTxOutpointIdx := uint32(0)
		for _, commitment := range commitmentData {
			if commitment.committingWallet == event.WalletPublicKeyHash {
				movingFundsTxOutpointIdx = commitment.targetWalletIndex
			}
		}

		movingFundsTxOutpoints = append(
			movingFundsTxOutpoints,
			bitcoin.TransactionOutpoint{
				TransactionHash: movingFundsTxHash,
				OutputIndex:     movingFundsTxOutpointIdx,
			},
		)
	}

	return movingFundsTxOutpoints, nil
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
	movingFundsTxOutputIndex uint32,
	hasMainUtxo bool,
	fee int64,
) (*tbtc.MovedFundsSweepProposal, error) {
	taskLogger.Infof("preparing a moved funds sweep proposal")

	// Estimate fee if it's missing.
	if fee <= 0 {
		taskLogger.Infof("estimating moved funds sweep transaction fee")

		_, _, _, _, _, _, _, sweepTxMaxTotalFee, _, _, _, err := mfst.chain.GetMovingFundsParameters()
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get moved funds sweep tx max total fee: [%w]",
				err,
			)
		}

		estimatedFee, err := EstimateMovedFundsSweepFee(
			mfst.btcChain,
			hasMainUtxo,
			sweepTxMaxTotalFee,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot estimate moving funds transaction fee: [%w]",
				err,
			)
		}

		fee = estimatedFee
	}

	taskLogger.Infof("moved funds sweep transaction fee: [%d]", fee)

	proposal := &tbtc.MovedFundsSweepProposal{
		MovingFundsTxHash:        movingFundsTxHash,
		MovingFundsTxOutputIndex: movingFundsTxOutputIndex,
		SweepTxFee:               big.NewInt(fee),
	}

	taskLogger.Infof("validating the moved funds sweep proposal")

	if err := tbtc.ValidateMovedFundsSweepProposal(
		taskLogger,
		walletPublicKeyHash,
		proposal,
		mfst.chain,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to verify moved funds sweep proposal: [%w]",
			err,
		)
	}

	return proposal, nil
}

// EstimateMovedFundsSweepFee estimates fee for the moved funds sweep transaction
// that merges the received main UTXO from the source wallets with the current
// wallet's main UTXO.
func EstimateMovedFundsSweepFee(
	btcChain bitcoin.Chain,
	hasMainUtxo bool,
	sweepTxMaxTotalFee uint64,
) (int64, error) {
	// The transaction always has an input coming from the moved funds
	// transferred by the source wallet. Additionally, it may have the second
	// input which is the wallet main UTXO.
	inputCount := 1
	if hasMainUtxo {
		inputCount++
	}

	sizeEstimator := bitcoin.NewTransactionSizeEstimator().
		AddPublicKeyHashInputs(inputCount, true).
		AddPublicKeyHashOutputs(1, true)

	transactionSize, err := sizeEstimator.VirtualSize()
	if err != nil {
		return 0, fmt.Errorf(
			"cannot estimate transaction virtual size: [%v]",
			err,
		)
	}

	feeEstimator := bitcoin.NewTransactionFeeEstimator(btcChain)

	totalFee, err := feeEstimator.EstimateFee(transactionSize)
	if err != nil {
		return 0, fmt.Errorf("cannot estimate transaction fee: [%v]", err)
	}

	if uint64(totalFee) > sweepTxMaxTotalFee {
		return 0, ErrSweepTxFeeTooHigh
	}

	return totalFee, nil
}
