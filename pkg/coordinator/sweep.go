package coordinator

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type btcTransaction = struct {
	FundingTxHash      bitcoin.Hash
	FundingOutputIndex uint32
}

const requiredFundingTxConfirmations = uint(6)
const depositScriptByteSize = 92

var (
	DepositsFormatDescription = `Deposits details should be provided as strings containing:
  - bitcoin transaction hash (unprefixed bitcoin transaction hash in reverse (RPC) order),
  - bitcoin transaction output index,
  - ethereum block number when the deposit was revealed to the chain.
The properties should be separated by semicolons, in the following format:
` + depositsFormatPattern + `
e.g. bd99d1d0a61fd104925d9b7ac997958aa8af570418b3fde091f7bfc561608865:1:8392394
`
	depositsFormatPattern = "<unprefixed bitcoin transaction hash>:<bitcoin transaction output index>:<ethereum reveal block number>"
	depositsFormatRegexp  = regexp.MustCompile(`^([[:xdigit:]]+):(\d+):(\d+)$`)
)

// ProposeDepositsSweep handles deposit sweep proposal request submission.
func ProposeDepositsSweep(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletStr string,
	fee int64,
	depositsString []string,
	dryRun bool,
) error {
	walletPublicKeyHash, err := hexToWalletPublicKeyHash(walletStr)
	if err != nil {
		return fmt.Errorf("failed extract wallet public key hash: %v", err)
	}

	btcTransactions, depositsRevealBlocks, err := parseDeposits(depositsString)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	proposal := &tbtc.DepositSweepProposal{
		WalletPublicKeyHash:  walletPublicKeyHash,
		DepositsKeys:         btcTransactions,
		SweepTxFee:           big.NewInt(fee),
		DepositsRevealBlocks: depositsRevealBlocks,
	}

	logger.Infof("validating the proposal...")
	if _, err := tbtc.ValidateDepositSweepProposal(
		logger,
		proposal,
		requiredFundingTxConfirmations,
		tbtcChain,
		btcChain,
	); err != nil {
		return fmt.Errorf("failed to verify deposit sweep proposal: %v", err)
	}

	if !dryRun {
		logger.Infof("submitting the proposal...")
		if err := tbtcChain.SubmitDepositSweepProposalWithReimbursement(proposal); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proposal: %v", err)
		}
	}

	return nil
}

func parseDeposits(depositsStrings []string) ([]btcTransaction, []*big.Int, error) {
	depositsKeys := make([]btcTransaction, len(depositsStrings))
	depositsRevealBlocks := make([]*big.Int, len(depositsStrings))

	for i, depositString := range depositsStrings {
		matched := depositsFormatRegexp.FindStringSubmatch(depositString)
		// Check if number of resolved entries match expected number of groups
		// for the given regexp.
		if len(matched) != 4 {
			return nil, nil, fmt.Errorf("failed to parse deposit: [%s]", depositString)
		}

		txHash, err := bitcoin.NewHashFromString(matched[1], bitcoin.ReversedByteOrder)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid bitcoin transaction hash [%s]: %v", matched[1], err)

		}

		outputIndex, err := strconv.ParseInt(matched[2], 10, 32)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid bitcoin transaction output index [%s]: %v", matched[2], err)
		}

		revealBlock, err := strconv.ParseInt(matched[3], 10, 32)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid reveal block number [%s]: %v", matched[3], err)
		}

		depositsKeys[i] = btcTransaction{
			FundingTxHash:      txHash,
			FundingOutputIndex: uint32(outputIndex),
		}

		depositsRevealBlocks[i] = big.NewInt(revealBlock)
	}

	return depositsKeys, depositsRevealBlocks, nil
}

// ValidateDepositString validates format of the string containing deposit details.
func ValidateDepositString(depositString string) error {
	if !depositsFormatRegexp.MatchString(depositString) {
		return fmt.Errorf(
			"[%s] doesn't match pattern: %s",
			depositString,
			depositsFormatPattern,
		)
	}
	return nil
}

// EstimateDepositsSweepFee computes the total fee for the Bitcoin deposits
// sweep transaction for the given depositsCount. If the provided depositsCount
// is 0, this function computes the total fee for Bitcoin deposits sweep
// transactions containing a various number of input deposits, from 1 up to the
// maximum count allowed by the WalletCoordinator contract. Computed fees for
// specific deposits counts are printed as table to the standard output,
// for example:
//
// ---------------------------------------------
// deposits count total fee (satoshis) sat/vbyte
//              1                  201         1
//              2                  292         1
//              3                  384         1
// ---------------------------------------------
//
// While making estimations, this function assumes a sweep transaction
// consists of:
// - 1 P2WPKH input being the current wallet main UTXO. That means the produced
//   fees may be overestimated for the very first sweep transaction of
//   each wallet.
// - N P2WSH inputs representing the deposits. Worth noting that real
//   transactions may contain legacy P2SH deposits as well so produced fees may
//   be underestimated in some rare cases.
// - 1 P2WPKH output
//
// If any of the estimated fees exceed the maximum fee allowed by the Bridge
// contract, the maximum fee is returned as result.
func EstimateDepositsSweepFee(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	depositsCount int,
) error {
	_, _, perDepositMaxFee, _, err := tbtcChain.GetDepositParameters()
	if err != nil {
		return fmt.Errorf("cannot get deposit tx max fee: [%v]", err)
	}

	fees := make(map[int]struct {
		totalFee       int64
		satPerVByteFee int64
	})
	var depositsCountKeys []int

	if depositsCount > 0 {
		depositsCountKeys = append(depositsCountKeys, depositsCount)
	} else {
		sweepMaxSize, err := tbtcChain.GetDepositSweepMaxSize()
		if err != nil {
			return fmt.Errorf("cannot get sweep max size: [%v]", sweepMaxSize)
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
			return fmt.Errorf(
				"cannot estimate fee for deposits count [%v]: [%v]",
				depositsCountKey,
				err,
			)
		}

		fees[depositsCountKey] = struct {
			totalFee       int64
			satPerVByteFee int64
		}{
			totalFee:       totalFee,
			satPerVByteFee: satPerVByteFee,
		}
	}

	err = printDepositsSweepFeeTable(fees)
	if err != nil {
		return fmt.Errorf("cannot print fees table: [%v]", err)
	}

	return nil
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
	// Make sure the proposed total fee does not exceed the maximum possible total fee.
	totalFee = int64(math.Min(float64(totalFee), float64(totalMaxFee)))
	// Compute the actual sat/vbyte fee for informational purposes.
	satPerVByteFee := math.Round(float64(totalFee) / float64(transactionSize))

	return totalFee, int64(satPerVByteFee), nil
}

func printDepositsSweepFeeTable(
	fees map[int]struct {
		totalFee       int64
		satPerVByteFee int64
	},
) error {
	writer := tabwriter.NewWriter(
		os.Stdout,
		2,
		4,
		1,
		' ',
		tabwriter.AlignRight,
	)

	_, err := fmt.Fprintf(writer, "deposits count\ttotal fee (satoshis)\tsat/vbyte\t\n")
	if err != nil {
		return err
	}

	var depositsCountKeys []int
	for depositsCountKey := range fees {
		depositsCountKeys = append(depositsCountKeys, depositsCountKey)
	}

	sort.Slice(depositsCountKeys, func(i, j int) bool {
		return depositsCountKeys[i] < depositsCountKeys[j]
	})

	for _, depositsCountKey := range depositsCountKeys {
		_, err := fmt.Fprintf(
			writer,
			"%v\t%v\t%v\t\n",
			depositsCountKey,
			fees[depositsCountKey].totalFee,
			fees[depositsCountKey].satPerVByteFee,
		)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush the writer: %v", err)
	}

	return nil
}
