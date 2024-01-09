package tbtcpg

import (
	"fmt"
	"math/big"
	"sort"

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
			"cannot get wallet's maun UTXO: [%w]",
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

	targetWallets, commitmentSubmitted, err := mft.FindTargetWallets(
		taskLogger,
		walletPublicKeyHash,
		walletChainData,
		uint64(walletBalance),
		liveWalletsCount,
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
		_, err := mft.SubmitMovingFundsCommitment()
		if err != nil {
			return nil, false, fmt.Errorf(
				"error while submitting moving funds commitment: [%w]",
				err,
			)
		}
	}

	return proposal, false, nil
}

func (mft *MovingFundsTask) FindTargetWallets(
	taskLogger log.StandardLogger,
	sourceWalletPublicKeyHash [20]byte,
	walletChainData *tbtc.WalletChainData,
	walletBalance uint64,
	liveWalletsCount uint32,
) ([][20]byte, bool, error) {
	if walletChainData.MovingFundsTargetWalletsCommitmentHash == [32]byte{} {
		targetWallets, err := mft.findNewTargetWallets(
			taskLogger,
			sourceWalletPublicKeyHash,
			walletBalance,
			liveWalletsCount,
		)
		return targetWallets, false, err
	}

	targetWallets, err := mft.retrieveCommittedTargetWallets(
		taskLogger,
		sourceWalletPublicKeyHash,
	)
	return targetWallets, true, err
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
		return nil, fmt.Errorf(
			"wallet max BTC transfer must be positive: [%w]", err,
		)
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
			"failed to get enough target wallets: required [%v]; gathered [%v]",
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

	// Look back one month when searching for the commitment event.
	filterLookBackSeconds := uint64(30 * 24 * 60 * 60)
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
				"target wallet [0x%x] is not Live",
				targetWallet,
			)
		}
	}

	return targetWallets, nil
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
