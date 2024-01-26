package tbtcpg

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"go.uber.org/zap"
)

var (
	// ErrMaxBtcTransferZero is the error returned when wallet max BTC transfer
	// parameter is zero.
	ErrMaxBtcTransferZero = fmt.Errorf(
		"wallet max BTC transfer must be positive",
	)

	// ErrNotEnoughTargetWallets is the error returned when the number of
	// gathered target wallets does not match the required target wallets count.
	ErrNotEnoughTargetWallets = fmt.Errorf("not enough target wallets")

	// ErrWrongCommitmentHash is the error returned when the hash calculated
	// from retrieved target wallets does not match the committed hash.
	ErrWrongCommitmentHash = fmt.Errorf(
		"target wallets hash must match commitment hash",
	)

	// ErrTargetWalletNotLive is the error returned when a target wallet is not
	// in the Live state.
	ErrTargetWalletNotLive = fmt.Errorf("target wallet is not live")

	// ErrNoExecutingOperator is the error returned when the task executing
	// operator is not found among the wallet operator IDs.
	ErrNoExecutingOperator = fmt.Errorf(
		"task executing operator not found among wallet operators",
	)

	// ErrTransactionNotIncluded is the error returned when the commitment
	// submission transaction was not included in the Ethereum blockchain.
	ErrTransactionNotIncluded = fmt.Errorf(
		"transaction not included in blockchain",
	)

	// ErrFeeTooHigh is the error returned when the estimated fee exceeds the
	// maximum fee allowed for the moving funds transaction.
	ErrFeeTooHigh = fmt.Errorf("estimated fee exceeds the maximum fee")
)

// MovingFundsCommitmentLookBackPeriod is the look-back period used when
// searching for submitted moving funds commitment events.
const MovingFundsCommitmentLookBackPeriod = 30 * 24 * time.Hour // 1 month

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

	// Check if the wallet is eligible for moving funds.
	walletChainData, err := mft.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get source wallet's chain data: [%w]",
			err,
		)
	}

	if walletChainData.State != tbtc.StateMovingFunds {
		taskLogger.Infof("source wallet not in MoveFunds state")
		return nil, false, nil
	}

	if walletChainData.PendingRedemptionsValue > 0 {
		taskLogger.Infof("source wallet has pending redemptions")
		return nil, false, nil
	}

	if walletChainData.PendingMovedFundsSweepRequestsCount > 0 {
		taskLogger.Infof("source wallet has pending moved funds sweep requests")
		return nil, false, nil
	}

	walletMainUtxo, err := tbtc.DetermineWalletMainUtxo(
		walletPublicKeyHash,
		mft.chain,
		mft.btcChain,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get wallet's main UTXO: [%w]",
			err,
		)
	}

	walletBalance := int64(0)
	if walletMainUtxo != nil {
		walletBalance = walletMainUtxo.Value
	}

	if walletBalance <= 0 {
		// The wallet's balance cannot be `0`. Since we are dealing with
		// a signed integer we also check it's not negative just in case.
		taskLogger.Infof("source wallet does not have a positive balance")
		return nil, false, nil
	}

	liveWalletsCount, err := mft.chain.GetLiveWalletsCount()
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get Live wallets count: [%w]",
			err,
		)
	}

	if liveWalletsCount == 0 {
		taskLogger.Infof("there are no Live wallets available")
		return nil, false, nil
	}

	targetWalletsCommitmentHash :=
		walletChainData.MovingFundsTargetWalletsCommitmentHash

	targetWallets, commitmentSubmitted, err := mft.FindTargetWallets(
		taskLogger,
		walletPublicKeyHash,
		targetWalletsCommitmentHash,
		uint64(walletBalance),
		liveWalletsCount,
	)
	if err != nil {
		return nil, false, fmt.Errorf("cannot find target wallets: [%w]", err)
	}

	proposal, err := mft.ProposeMovingFunds(
		taskLogger,
		walletPublicKeyHash,
		walletMainUtxo,
		targetWallets,
		0,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot prepare moving funds proposal: [%w]",
			err,
		)
	}

	if !commitmentSubmitted {
		walletMemberIDs, walletMemberIndex, err := mft.GetWalletMembersInfo(
			request.WalletOperators,
			request.ExecutingOperator,
		)
		if err != nil {
			return nil, false, fmt.Errorf(
				"cannot get wallet members IDs: [%w]",
				err,
			)
		}

		err = mft.SubmitMovingFundsCommitment(
			walletPublicKeyHash,
			walletMainUtxo,
			walletMemberIDs,
			walletMemberIndex,
			targetWallets,
		)
		if err != nil {
			return nil, false, fmt.Errorf(
				"error while submitting moving funds commitment: [%w]",
				err,
			)
		}
	}

	return proposal, false, nil
}

// FindTargetWallets returns a list of target wallets for the moving funds
// procedure. If the source wallet has not submitted moving funds commitment yet
// a new list of target wallets is prepared. If the source wallet has already
// submitted the commitment, the returned target wallet list is prepared based
// on the submitted commitment event.
func (mft *MovingFundsTask) FindTargetWallets(
	taskLogger log.StandardLogger,
	sourceWalletPublicKeyHash [20]byte,
	targetWalletsCommitmentHash [32]byte,
	walletBalance uint64,
	liveWalletsCount uint32,
) ([][20]byte, bool, error) {
	if targetWalletsCommitmentHash == [32]byte{} {
		targetWallets, err := mft.findNewTargetWallets(
			taskLogger,
			sourceWalletPublicKeyHash,
			walletBalance,
			liveWalletsCount,
		)

		return targetWallets, false, err
	} else {
		targetWallets, err := mft.retrieveCommittedTargetWallets(
			taskLogger,
			sourceWalletPublicKeyHash,
			targetWalletsCommitmentHash,
		)

		return targetWallets, true, err
	}
}

func (mft *MovingFundsTask) findNewTargetWallets(
	taskLogger log.StandardLogger,
	sourceWalletPublicKeyHash [20]byte,
	walletBalance uint64,
	liveWalletsCount uint32,
) ([][20]byte, error) {
	taskLogger.Infof(
		"commitment not submitted yet; looking for new target wallets",
	)

	_, _, _, _, _, walletMaxBtcTransfer, _, err := mft.chain.GetWalletParameters()
	if err != nil {
		return nil, fmt.Errorf("cannot get wallet parameters: [%w]", err)
	}

	if walletMaxBtcTransfer == 0 {
		return nil, ErrMaxBtcTransferZero
	}

	ceilingDivide := func(x, y uint64) uint64 {
		// The divisor must be positive, but we do not need to check it as
		// this function will be executed with wallet max BTC transfer as
		// the divisor and we already ensured it is positive.
		return (x + y - 1) / y
	}
	min := func(x, y uint64) uint64 {
		if x < y {
			return x
		}
		return y
	}

	targetWalletsCount := min(
		uint64(liveWalletsCount),
		ceilingDivide(walletBalance, walletMaxBtcTransfer),
	)

	// Prepare a list of target wallets using the new wallets registration
	// events. Retrieve only the necessary number of live wallets.
	// The iteration is started from the end of the list as the newest wallets
	// are located there and have highest the chance of being Live.
	events, err := mft.chain.PastNewWalletRegisteredEvents(nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past new wallet registered events: [%v]",
			err,
		)
	}

	targetWallets := make([][20]byte, 0)

	for i := len(events) - 1; i >= 0; i-- {
		walletPubKeyHash := events[i].WalletPublicKeyHash
		if walletPubKeyHash == sourceWalletPublicKeyHash {
			// Just in case make sure not to include the source wallet
			// itself.
			continue
		}

		wallet, err := mft.chain.GetWallet(walletPubKeyHash)
		if err != nil {
			taskLogger.Errorf(
				"failed to get wallet data for wallet with PKH [0x%x]: [%v]",
				walletPubKeyHash,
				err,
			)
			continue
		}

		if wallet.State == tbtc.StateLive {
			targetWallets = append(targetWallets, walletPubKeyHash)
		}
		if len(targetWallets) == int(targetWalletsCount) {
			// Stop the iteration if enough live wallets have been gathered.
			break
		}
	}

	if len(targetWallets) != int(targetWalletsCount) {
		return nil, fmt.Errorf(
			"%w: required [%v] target wallets; gathered [%v]",
			ErrNotEnoughTargetWallets,
			targetWalletsCount,
			len(targetWallets),
		)
	}

	// Sort the target wallets according to their numerical representation
	// as the on-chain contract expects.
	sort.Slice(targetWallets, func(i, j int) bool {
		bigIntI := new(big.Int).SetBytes(targetWallets[i][:])
		bigIntJ := new(big.Int).SetBytes(targetWallets[j][:])
		return bigIntI.Cmp(bigIntJ) < 0
	})

	logger.Infof("gathered [%v] target wallets", len(targetWallets))

	return targetWallets, nil
}

func (mft *MovingFundsTask) retrieveCommittedTargetWallets(
	taskLogger log.StandardLogger,
	sourceWalletPublicKeyHash [20]byte,
	targetWalletsCommitmentHash [32]byte,
) ([][20]byte, error) {
	taskLogger.Infof(
		"commitment already submitted; retrieving committed target wallets",
	)

	blockCounter, err := mft.chain.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get block counter: [%w]",
			err,
		)
	}

	currentBlockNumber, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get current block number: [%w]",
			err,
		)
	}

	filterLookBackSeconds := uint64(MovingFundsCommitmentLookBackPeriod.Seconds())
	filterLookBackBlocks := filterLookBackSeconds / uint64(mft.chain.AverageBlockTime().Seconds())
	filterStartBlock := currentBlockNumber - filterLookBackBlocks

	filter := &tbtc.MovingFundsCommitmentSubmittedEventFilter{
		StartBlock:          filterStartBlock,
		WalletPublicKeyHash: [][20]byte{sourceWalletPublicKeyHash},
	}

	events, err := mft.chain.PastMovingFundsCommitmentSubmittedEvents(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past moving funds commitment submitted events: [%w]",
			err,
		)
	}

	// Moving funds commitment can be submitted only once for a given wallet.
	// Check just in case.
	if len(events) != 1 {
		return nil, fmt.Errorf(
			"unexpected number of moving funds commitment submitted events: [%v]",
			len(events),
		)
	}

	targetWallets := events[0].TargetWallets

	// Make sure all the target wallets are Live.
	for _, targetWallet := range targetWallets {
		targetWalletChainData, err := mft.chain.GetWallet(targetWallet)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get wallet data for target wallet [0x%x]: [%w]",
				targetWallet,
				err,
			)
		}

		if targetWalletChainData.State != tbtc.StateLive {
			return nil, fmt.Errorf(
				"%w: [0x%x]",
				ErrTargetWalletNotLive,
				targetWallet,
			)
		}
	}

	// Just in case check if the hash of the target wallets matches the moving
	// funds target wallets commitment hash.
	packedWallets := []byte{}
	for _, wallet := range targetWallets {
		packedWallets = append(packedWallets, wallet[:]...)
	}

	calculatedHash := crypto.Keccak256(packedWallets)

	if !bytes.Equal(calculatedHash, targetWalletsCommitmentHash[:]) {
		return nil, ErrWrongCommitmentHash
	}

	return targetWallets, nil
}

// GetWalletMembersInfo returns the wallet member IDs based on the provided
// wallet operator addresses. Additionally, it returns the position of the
// moving funds task execution operator on the list.
func (mft *MovingFundsTask) GetWalletMembersInfo(
	walletOperators []chain.Address,
	executingOperator chain.Address,
) ([]uint32, uint32, error) {
	walletMemberIDs := make([]uint32, 0)
	walletMemberIndex := 0

	for index, operatorAddress := range walletOperators {
		operatorID, err := mft.chain.GetOperatorID(operatorAddress)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get operator ID: [%w]", err)
		}

		// Increment the index by 1 as operator indexing starts at 1, not 0.
		// This ensures the operator's position is correctly identified in the
		// range [1, walletOperators.length].
		if operatorAddress == executingOperator {
			walletMemberIndex = index + 1
		}

		walletMemberIDs = append(walletMemberIDs, operatorID)
	}

	// The task executing operator must always be on the wallet operators list.
	if walletMemberIndex == 0 {
		return nil, 0, ErrNoExecutingOperator
	}

	return walletMemberIDs, uint32(walletMemberIndex), nil
}

// SubmitMovingFundsCommitment submits the moving funds commitment and waits
// until the transaction has entered the Ethereum blockchain.
func (mft *MovingFundsTask) SubmitMovingFundsCommitment(
	walletPublicKeyHash [20]byte,
	walletMainUTXO *bitcoin.UnspentTransactionOutput,
	walletMembersIDs []uint32,
	walletMemberIndex uint32,
	targetWallets [][20]byte,
) error {
	err := mft.chain.SubmitMovingFundsCommitment(
		walletPublicKeyHash,
		*walletMainUTXO,
		walletMembersIDs,
		walletMemberIndex,
		targetWallets,
	)
	if err != nil {
		return fmt.Errorf(
			"error while submitting moving funds commitment to chain: [%w]",
			err,
		)
	}

	blockCounter, err := mft.chain.BlockCounter()
	if err != nil {
		return fmt.Errorf("error getting block counter [%w]", err)
	}

	currentBlock, err := blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("error getting current block [%w]", err)
	}

	// To verify the commitment transaction has entered the Ethereum blockchain
	// check that the commitment hash is not zero.
	stateCheck := func() (bool, error) {
		walletData, err := mft.chain.GetWallet(walletPublicKeyHash)
		if err != nil {
			return false, err
		}

		return walletData.MovingFundsTargetWalletsCommitmentHash != [32]byte{}, nil
	}

	// Wait `5` blocks since the current block and perform the transaction state
	// check. If the transaction has not entered the blockchain, consider it an
	// error.
	result, err := ethereum.WaitForBlockConfirmations(
		blockCounter,
		currentBlock,
		5,
		stateCheck,
	)
	if err != nil {
		return fmt.Errorf(
			"error while waiting for transaction confirmation [%w]",
			err,
		)
	}

	if !result {
		return ErrTransactionNotIncluded
	}

	return nil
}

// ProposeMovingFunds returns a moving funds proposal.
func (mft *MovingFundsTask) ProposeMovingFunds(
	taskLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	mainUTXO *bitcoin.UnspentTransactionOutput,
	targetWallets [][20]byte,
	fee int64,
) (*tbtc.MovingFundsProposal, error) {
	if len(targetWallets) == 0 {
		return nil, fmt.Errorf("target wallets list is empty")
	}

	taskLogger.Infof("preparing a moving funds proposal")

	// Estimate fee if it's missing.
	if fee <= 0 {
		taskLogger.Infof("estimating moving funds transaction fee")

		txMaxTotalFee, _, _, _, _, _, _, _, _, _, _, err := mft.chain.GetMovingFundsParameters()
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get moving funds tx max total fee: [%w]",
				err,
			)
		}

		estimatedFee, err := EstimateMovingFundsFee(
			mft.btcChain,
			len(targetWallets),
			txMaxTotalFee,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot estimate moving funds transaction fee: [%w]",
				err,
			)
		}

		fee = estimatedFee
	}

	proposal := &tbtc.MovingFundsProposal{
		TargetWallets:    targetWallets,
		MovingFundsTxFee: big.NewInt(fee),
	}

	taskLogger.Infof("validating the moving funds proposal")
	if _, err := tbtc.ValidateMovingFundsProposal(
		taskLogger,
		walletPublicKeyHash,
		mainUTXO,
		proposal,
		mft.chain,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to verify moving funds proposal: [%w]",
			err,
		)
	}

	return proposal, nil
}

func (mft *MovingFundsTask) ActionType() tbtc.WalletActionType {
	return tbtc.ActionMovingFunds
}

// EstimateMovingFundsFee estimates fee for the moving funds transaction that
// moves funds from the source wallet to target wallets.
func EstimateMovingFundsFee(
	btcChain bitcoin.Chain,
	targetWalletsCount int,
	txMaxTotalFee uint64,
) (int64, error) {
	sizeEstimator := bitcoin.NewTransactionSizeEstimator().
		AddPublicKeyHashInputs(1, true).
		AddPublicKeyHashOutputs(targetWalletsCount, true)

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

	if uint64(totalFee) > txMaxTotalFee {
		return 0, ErrFeeTooHigh
	}

	return totalFee, nil
}
