package tbtcpg

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// Use the worst-case 126-byte deposit script with embedded extra data for estimation.
// This will ensure that deposit sweep transaction fees are not underestimated.
const depositScriptByteSize = 126

// DepositSweepTask is a task that may produce a deposit sweep proposal.
type DepositSweepTask struct {
	chain    Chain
	btcChain bitcoin.Chain
}

func NewDepositSweepTask(
	chain Chain,
	btcChain bitcoin.Chain,
) *DepositSweepTask {
	return &DepositSweepTask{
		chain:    chain,
		btcChain: btcChain,
	}
}

func (dst *DepositSweepTask) Run(request *tbtc.CoordinationProposalRequest) (
	tbtc.CoordinationProposal,
	bool,
	error,
) {
	walletPublicKeyHash := request.WalletPublicKeyHash

	taskLogger := logger.With(
		zap.String("task", dst.ActionType().String()),
		zap.String("walletPKH", fmt.Sprintf("0x%x", walletPublicKeyHash)),
	)

	depositSweepMaxSize, err := dst.chain.GetDepositSweepMaxSize()
	if err != nil {
		return nil, false, fmt.Errorf(
			"failed to get deposit sweep max size: [%w]",
			err,
		)
	}

	deposits, err := dst.FindDepositsToSweep(
		taskLogger,
		walletPublicKeyHash,
		depositSweepMaxSize,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot find deposits to sweep: [%w]",
			err,
		)
	}

	if len(deposits) == 0 {
		taskLogger.Info("no deposits to sweep")
		return nil, false, nil
	}

	proposal, err := dst.ProposeDepositsSweep(
		taskLogger,
		walletPublicKeyHash,
		deposits,
		0,
	)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot prepare deposit sweep proposal: [%w]",
			err,
		)
	}

	return proposal, true, nil
}

func (dst *DepositSweepTask) ActionType() tbtc.WalletActionType {
	return tbtc.ActionDepositSweep
}

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
	return findDeposits(
		logger,
		chain,
		btcChain,
		walletPublicKeyHash,
		maxNumberOfDeposits,
		skipSwept,
		skipUnconfirmed,
	)
}

// findDeposits finds deposits according to the given criteria.
func findDeposits(
	fnLogger log.StandardLogger,
	chain Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits int,
	skipSwept bool,
	skipUnconfirmed bool,
) ([]*Deposit, error) {
	fnLogger.Infof("reading revealed deposits from chain")

	depositMinAgeSeconds, err := chain.GetDepositMinAge()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get deposit minimum age: [%w]",
			err,
		)
	}
	depositMinAge := time.Duration(depositMinAgeSeconds) * time.Second

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

	fnLogger.Infof("found [%d] DepositRevealed events", len(depositRevealedEvents))

	// Take the oldest first
	sort.SliceStable(depositRevealedEvents, func(i, j int) bool {
		return depositRevealedEvents[i].BlockNumber < depositRevealedEvents[j].BlockNumber
	})

	fnLogger.Infof("getting deposits details")

	resultSliceCapacity := len(depositRevealedEvents)
	if maxNumberOfDeposits > 0 {
		resultSliceCapacity = maxNumberOfDeposits
	}

	// Capture time now for computations.
	timeNow := time.Now()

	result := make([]*Deposit, 0, resultSliceCapacity)
	for _, event := range depositRevealedEvents {
		if len(result) == cap(result) {
			break
		}

		depositKey := chain.BuildDepositKey(event.FundingTxHash, event.FundingOutputIndex)
		depositKeyStr := depositKey.Text(16)

		fnLogger.Debugf("getting details of deposit [%s]", depositKeyStr)

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
				"no deposit request for key [%s]",
				depositKeyStr,
			)
		}

		matureAt := depositRequest.RevealedAt.Add(depositMinAge)
		if !timeNow.After(matureAt) {
			fnLogger.Infof("deposit [%s] is not old enough", depositKeyStr)
			continue
		}

		isSwept := depositRequest.SweptAt.Unix() != 0
		if skipSwept && isSwept {
			fnLogger.Debugf("deposit [%s] is already swept", depositKeyStr)
			continue
		}

		confirmations, err := btcChain.GetTransactionConfirmations(event.FundingTxHash)
		if err != nil {
			fnLogger.Errorf(
				"failed to get bitcoin transaction confirmations: [%v]",
				err,
			)
		}

		if skipUnconfirmed && confirmations < tbtc.DepositSweepRequiredFundingTxConfirmations {
			fnLogger.Debugf(
				"deposit [%s] funding transaction doesn't have enough confirmations: [%d/%d]",
				depositKeyStr,
				confirmations,
				tbtc.DepositSweepRequiredFundingTxConfirmations,
			)
			continue
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
// maxNumberOfDeposits is used as a ceiling for the number of deposits in the
// result. If number of discovered deposits meets the maxNumberOfDeposits the
// function will stop fetching more deposits.
// This function will return a list of deposits from the wallet that can be swept.
// Deposits with insufficient number of funding transaction confirmations will
// not be taken into consideration for sweeping.
//
// TODO: Cache immutable data
func (dst *DepositSweepTask) FindDepositsToSweep(
	taskLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits uint16,
) ([]*DepositReference, error) {
	if walletPublicKeyHash == [20]byte{} {
		return nil, fmt.Errorf("wallet public key hash is required")
	}

	taskLogger.Infof("fetching max [%d] deposits", maxNumberOfDeposits)

	unsweptDeposits, err := findDeposits(
		taskLogger,
		dst.chain,
		dst.btcChain,
		walletPublicKeyHash,
		int(maxNumberOfDeposits),
		true,
		true,
	)
	if err != nil {
		return nil, err
	}

	depositsToSweep := unsweptDeposits

	if len(depositsToSweep) == 0 {
		return nil, nil
	}

	taskLogger.Infof(
		"found [%d] deposits to sweep",
		len(depositsToSweep),
	)

	for _, deposit := range depositsToSweep {
		taskLogger.Infof(
			"deposit [%s] - [%s]",
			deposit.DepositKey,
			fmt.Sprintf(
				"reveal block: [%d], funding transaction: [%s], output index: [%d]",
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

	return depositsRefs, nil
}

// ProposeDepositsSweep returns a deposit sweep proposal.
func (dst *DepositSweepTask) ProposeDepositsSweep(
	taskLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	deposits []*DepositReference,
	fee int64,
) (*tbtc.DepositSweepProposal, error) {
	if len(deposits) == 0 {
		return nil, fmt.Errorf("deposits list is empty")
	}

	taskLogger.Infof("preparing a deposit sweep proposal")

	// Estimate fee if it's missing.
	if fee <= 0 {
		taskLogger.Infof("estimating sweep transaction fee")
		var err error
		_, _, perDepositMaxFee, _, err := dst.chain.GetDepositParameters()
		if err != nil {
			return nil, fmt.Errorf("cannot get deposit tx max fee: [%w]", err)
		}

		estimatedFee, _, err := estimateDepositsSweepFee(
			dst.btcChain,
			len(deposits),
			perDepositMaxFee,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot estimate sweep transaction fee: [%v]", err)
		}

		fee = estimatedFee
	}

	taskLogger.Infof("sweep transaction fee: [%d]", fee)

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
		DepositsKeys:         depositsKeys,
		SweepTxFee:           big.NewInt(fee),
		DepositsRevealBlocks: depositsRevealBlocks,
	}

	taskLogger.Infof("validating the deposit sweep proposal")

	if _, err := tbtc.ValidateDepositSweepProposal(
		taskLogger,
		walletPublicKeyHash,
		proposal,
		tbtc.DepositSweepRequiredFundingTxConfirmations,
		dst.chain,
		dst.btcChain,
	); err != nil {
		return nil, fmt.Errorf("failed to verify deposit sweep proposal: %v", err)
	}

	return proposal, nil
}

// EstimateDepositsSweepFee computes the total fee for the Bitcoin deposits
// sweep transaction for the given depositsCount. If the provided depositsCount
// is 0, this function computes the total fee for Bitcoin deposits sweep
// transactions containing a various number of input deposits, from 1 up to the
// maximum count allowed by the WalletProposalValidator contract. Computed fees for
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
