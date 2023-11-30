package wallet

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

const depositScriptByteSize = 92

// DepositReference holds some data allowing to identify and refer to a deposit.
type DepositReference struct {
	FundingTxHash      bitcoin.Hash
	FundingOutputIndex uint32
	RevealBlock        uint64
}

// Deposit holds some detailed data about a deposit.
type Deposit struct {
	DepositReference

	WalletPublicKeyHash [20]byte
	ScriptType          bitcoin.ScriptType
	DepositKey          string
	IsSwept             bool
	AmountBtc           float64
	Confirmations       uint
}

// FindDeposits finds deposits according to the given criteria.
func FindDeposits(
	chain Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits int,
	skipSwept bool,
	skipUnconfirmed bool,
) ([]*Deposit, error) {
	logger.Infof("reading revealed deposits from chain...")

	filter := &tbtc.DepositRevealedEventFilter{}
	if walletPublicKeyHash != [20]byte{} {
		filter.WalletPublicKeyHash = [][20]byte{walletPublicKeyHash}
	}

	depositRevealedEvents, err := chain.PastDepositRevealedEvents(filter)
	if err != nil {
		return []*Deposit{}, fmt.Errorf(
			"failed to get past deposit revealed events: [%w]",
			err,
		)
	}

	logger.Infof("found %d DepositRevealed events", len(depositRevealedEvents))

	// Take the oldest first
	sort.SliceStable(depositRevealedEvents, func(i, j int) bool {
		return depositRevealedEvents[i].BlockNumber < depositRevealedEvents[j].BlockNumber
	})

	logger.Infof("getting deposits details...")

	resultSliceCapacity := len(depositRevealedEvents)
	if maxNumberOfDeposits > 0 {
		resultSliceCapacity = maxNumberOfDeposits
	}

	result := make([]*Deposit, 0, resultSliceCapacity)
	for i, event := range depositRevealedEvents {
		if len(result) == cap(result) {
			break
		}

		logger.Debugf("getting details of deposit %d/%d", i+1, len(depositRevealedEvents))

		depositKey := chain.BuildDepositKey(event.FundingTxHash, event.FundingOutputIndex)

		depositRequest, found, err := chain.GetDepositRequest(
			event.FundingTxHash,
			event.FundingOutputIndex,
		)
		if err != nil {
			return result, fmt.Errorf(
				"failed to get deposit request: [%w]",
				err,
			)
		}

		if !found {
			return nil, fmt.Errorf(
				"no deposit request for key [0x%x]",
				depositKey.Text(16),
			)
		}

		isSwept := depositRequest.SweptAt.Unix() != 0
		if skipSwept && isSwept {
			logger.Debugf("deposit %d/%d is already swept", i+1, len(depositRevealedEvents))
			continue
		}

		confirmations, err := btcChain.GetTransactionConfirmations(event.FundingTxHash)
		if err != nil {
			logger.Errorf(
				"failed to get bitcoin transaction confirmations: [%v]",
				err,
			)
		}

		if skipUnconfirmed && confirmations < tbtc.DepositSweepRequiredFundingTxConfirmations {
			logger.Debugf(
				"deposit %d/%d funding transaction doesn't have enough confirmations: %d/%d",
				i+1, len(depositRevealedEvents),
				confirmations, tbtc.DepositSweepRequiredFundingTxConfirmations)
			continue
		}

		var scriptType bitcoin.ScriptType
		depositTransaction, err := btcChain.GetTransaction(
			event.FundingTxHash,
		)
		if err != nil {
			logger.Errorf(
				"failed to get deposit transaction data: [%v]",
				err,
			)
		} else {
			publicKeyScript := depositTransaction.Outputs[event.FundingOutputIndex].PublicKeyScript
			scriptType = bitcoin.GetScriptType(publicKeyScript)
		}

		result = append(
			result,
			&Deposit{
				DepositReference: DepositReference{
					FundingTxHash:      event.FundingTxHash,
					FundingOutputIndex: event.FundingOutputIndex,
					RevealBlock:        event.BlockNumber,
				},
				WalletPublicKeyHash: event.WalletPublicKeyHash,
				ScriptType:          scriptType,
				DepositKey:          hexutils.Encode(depositKey.Bytes()),
				IsSwept:             isSwept,
				AmountBtc:           convertSatToBtc(float64(depositRequest.Amount)),
				Confirmations:       confirmations,
			},
		)
	}

	return result, nil
}

// FindDepositsToSweep finds deposits that can be swept.
// If a wallet public key hash is provided, it will find unswept deposits for the
// given wallet. If a wallet public key hash is nil, it will check all wallets
// starting from the oldest one to find a first wallet containing unswept deposits
// and return those deposits.
// maxNumberOfDeposits is used as a ceiling for the number of deposits in the
// result. If number of discovered deposits meets the maxNumberOfDeposits the
// function will stop fetching more deposits.
// This function will return a wallet public key hash and a list of deposits from
// the wallet that can be swept.
// Deposits with insufficient number of funding transaction confirmations will
// not be taken into consideration for sweeping.
// The result will not mix deposits for different wallets.
//
// TODO: Cache immutable data
func FindDepositsToSweep(
	chain Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits uint16,
) ([20]byte, []*DepositReference, error) {
	logger.Infof("deposit sweep max size: %d", maxNumberOfDeposits)

	getDepositsToSweepFromWallet := func(walletToSweep [20]byte) ([]*Deposit, error) {
		unsweptDeposits, err := FindDeposits(
			chain,
			btcChain,
			walletToSweep,
			int(maxNumberOfDeposits),
			true,
			true,
		)
		if err != nil {
			return nil,
				fmt.Errorf(
					"failed to get deposits for [%s] wallet: [%w]",
					hexutils.Encode(walletToSweep[:]),
					err,
				)
		}
		return unsweptDeposits, nil
	}

	var depositsToSweep []*Deposit
	// If walletPublicKeyHash is not provided we need to find a wallet that has
	// unswept deposits.
	if walletPublicKeyHash == [20]byte{} {
		walletRegisteredEvents, err := chain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to get registered wallets: [%w]", err)
		}

		// Take the oldest first
		sort.SliceStable(walletRegisteredEvents, func(i, j int) bool {
			return walletRegisteredEvents[i].BlockNumber < walletRegisteredEvents[j].BlockNumber
		})

		sweepingWallets := walletRegisteredEvents
		// Only two the most recently created wallets are sweeping.
		if len(walletRegisteredEvents) >= 2 {
			sweepingWallets = walletRegisteredEvents[len(walletRegisteredEvents)-2:]
		}

		for _, registeredWallet := range sweepingWallets {
			logger.Infof(
				"fetching deposits from wallet [%s]...",
				hexutils.Encode(registeredWallet.WalletPublicKeyHash[:]),
			)

			unsweptDeposits, err := getDepositsToSweepFromWallet(
				registeredWallet.WalletPublicKeyHash,
			)
			if err != nil {
				return [20]byte{}, nil, err
			}

			// Check if there are any unswept deposits in this wallet. If so
			// sweep this wallet and don't check the other wallet.
			if len(unsweptDeposits) > 0 {
				walletPublicKeyHash = registeredWallet.WalletPublicKeyHash
				depositsToSweep = unsweptDeposits
				break
			}
		}
	} else {
		logger.Infof(
			"fetching deposits from wallet [%s]...",
			hexutils.Encode(walletPublicKeyHash[:]),
		)
		unsweptDeposits, err := getDepositsToSweepFromWallet(
			walletPublicKeyHash,
		)
		if err != nil {
			return [20]byte{}, nil, err
		}
		depositsToSweep = unsweptDeposits
	}

	if len(depositsToSweep) == 0 {
		return [20]byte{}, nil, nil
	}

	logger.Infof(
		"found [%d] deposits to sweep for wallet [%s]",
		len(depositsToSweep),
		hexutils.Encode(walletPublicKeyHash[:]),
	)

	for i, deposit := range depositsToSweep {
		logger.Infof(
			"deposit [%d/%d] - %s",
			i+1,
			len(depositsToSweep),
			fmt.Sprintf(
				"depositKey: [%s], reveal block: [%d], funding transaction: [%s], output index: [%d]",
				deposit.DepositKey,
				deposit.RevealBlock,
				deposit.FundingTxHash.Hex(bitcoin.ReversedByteOrder),
				deposit.FundingOutputIndex,
			))
	}

	depositsRefs := make([]*DepositReference, len(depositsToSweep))
	for i, deposit := range depositsToSweep {
		depositsRefs[i] = &DepositReference{
			FundingTxHash:      deposit.FundingTxHash,
			FundingOutputIndex: deposit.FundingOutputIndex,
			RevealBlock:        deposit.RevealBlock,
		}
	}

	return walletPublicKeyHash, depositsRefs, nil
}

// ProposeDepositsSweep handles deposit sweep proposal request submission.
func ProposeDepositsSweep(
	chain Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	fee int64,
	deposits []*DepositReference,
	dryRun bool,
) error {
	if len(deposits) == 0 {
		return fmt.Errorf("deposits list is empty")
	}

	// Estimate fee if it's missing.
	if fee <= 0 {
		logger.Infof("estimating sweep transaction fee...")
		var err error
		_, _, perDepositMaxFee, _, err := chain.GetDepositParameters()
		if err != nil {
			return fmt.Errorf("cannot get deposit tx max fee: [%w]", err)
		}

		estimatedFee, _, err := estimateDepositsSweepFee(
			btcChain,
			len(deposits),
			perDepositMaxFee,
		)
		if err != nil {
			return fmt.Errorf("cannot estimate sweep transaction fee: [%v]", err)
		}

		fee = estimatedFee
	}

	logger.Infof("sweep transaction fee: [%d]", fee)

	logger.Infof("preparing a deposit sweep proposal...")

	depositsKeys := make([]struct {
		FundingTxHash      bitcoin.Hash
		FundingOutputIndex uint32
	}, len(deposits))

	depositsRevealBlocks := make([]*big.Int, len(deposits))

	for i, deposit := range deposits {
		depositsKeys[i] = struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			FundingTxHash:      deposit.FundingTxHash,
			FundingOutputIndex: deposit.FundingOutputIndex,
		}
		depositsRevealBlocks[i] = big.NewInt(int64(deposit.RevealBlock))
	}

	proposal := &tbtc.DepositSweepProposal{
		WalletPublicKeyHash:  walletPublicKeyHash,
		DepositsKeys:         depositsKeys,
		SweepTxFee:           big.NewInt(fee),
		DepositsRevealBlocks: depositsRevealBlocks,
	}

	logger.Infof("validating the deposit sweep proposal...")
	if _, err := tbtc.ValidateDepositSweepProposal(
		logger,
		proposal,
		tbtc.DepositSweepRequiredFundingTxConfirmations,
		chain,
		btcChain,
	); err != nil {
		return fmt.Errorf("failed to verify deposit sweep proposal: %v", err)
	}

	if !dryRun {
		logger.Infof("submitting the deposit sweep proposal...")
		if err := chain.SubmitDepositSweepProposalWithReimbursement(proposal); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proposal: %v", err)
		}
	}

	return nil
}

// EstimateDepositsSweepFee computes the total fee for the Bitcoin deposits
// sweep transaction for the given depositsCount. If the provided depositsCount
// is 0, this function computes the total fee for Bitcoin deposits sweep
// transactions containing a various number of input deposits, from 1 up to the
// maximum count allowed by the WalletCoordinator contract. Computed fees for
// specific deposits counts are returned as a map.
//
// While making estimations, this function assumes a sweep transaction
// consists of:
//   - 1 P2WPKH input being the current wallet main UTXO. That means the produced
//     fees may be overestimated for the very first sweep transaction of
//     each wallet.
//   - N P2WSH inputs representing the deposits. Worth noting that real
//     transactions may contain legacy P2SH deposits as well so produced fees may
//     be underestimated in some rare cases.
//   - 1 P2WPKH output
//
// If any of the estimated fees exceed the maximum fee allowed by the Bridge
// contract, an error is returned as result.
func EstimateDepositsSweepFee(
	chain Chain,
	btcChain bitcoin.Chain,
	depositsCount int,
) (
	map[int]struct {
		TotalFee       int64
		SatPerVByteFee int64
	},
	error,
) {
	_, _, perDepositMaxFee, _, err := chain.GetDepositParameters()
	if err != nil {
		return nil, fmt.Errorf("cannot get deposit tx max fee: [%v]", err)
	}

	fees := make(map[int]struct {
		TotalFee       int64
		SatPerVByteFee int64
	})
	var depositsCountKeys []int

	if depositsCount > 0 {
		depositsCountKeys = append(depositsCountKeys, depositsCount)
	} else {
		sweepMaxSize, err := chain.GetDepositSweepMaxSize()
		if err != nil {
			return nil, fmt.Errorf("cannot get sweep max size: [%v]", sweepMaxSize)
		}

		for i := 1; i <= int(sweepMaxSize); i++ {
			depositsCountKeys = append(depositsCountKeys, i)
		}
	}

	for _, depositsCountKey := range depositsCountKeys {
		totalFee, satPerVByteFee, err := estimateDepositsSweepFee(
			btcChain,
			depositsCountKey,
			perDepositMaxFee,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot estimate fee for deposits count [%v]: [%v]",
				depositsCountKey,
				err,
			)
		}

		fees[depositsCountKey] = struct {
			TotalFee       int64
			SatPerVByteFee int64
		}{
			TotalFee:       totalFee,
			SatPerVByteFee: satPerVByteFee,
		}
	}

	return fees, nil
}

func estimateDepositsSweepFee(
	btcChain bitcoin.Chain,
	depositsCount int,
	perDepositMaxFee uint64,
) (int64, int64, error) {
	transactionSize, err := bitcoin.NewTransactionSizeEstimator().
		// 1 P2WPKH main UTXO input.
		AddPublicKeyHashInputs(1, true).
		// depositsCount P2WSH deposit inputs.
		AddScriptHashInputs(depositsCount, depositScriptByteSize, true).
		// 1 P2WPKH output.
		AddPublicKeyHashOutputs(1, true).
		VirtualSize()
	if err != nil {
		return 0, 0, fmt.Errorf("cannot estimate transaction virtual size: [%v]", err)
	}

	feeEstimator := bitcoin.NewTransactionFeeEstimator(btcChain)

	totalFee, err := feeEstimator.EstimateFee(transactionSize)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot estimate transaction fee: [%v]", err)
	}

	// Compute the maximum possible total fee for the entire sweep transaction.
	totalMaxFee := uint64(depositsCount) * perDepositMaxFee

	if uint64(totalFee) > totalMaxFee {
		return 0, 0, fmt.Errorf("estimated fee exceeds the maximum fee")
	}

	// Compute the actual sat/vbyte fee for informational purposes.
	satPerVByteFee := math.Round(float64(totalFee) / float64(transactionSize))

	return totalFee, int64(satPerVByteFee), nil
}

func convertSatToBtc(sats float64) float64 {
	return sats / float64(100000000)
}
