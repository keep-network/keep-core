package coordinator

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type btcTransactions = struct {
	FundingTxHash      bitcoin.Hash
	FundingOutputIndex uint32
}

const requiredFundingTxConfirmations = uint(6)

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
		if err := tbtcChain.SubmitDepositSweepProposal(proposal); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proposal: %v", err)
		}
	}

	return nil
}

func parseDeposits(depositsStrings []string) ([]btcTransactions, []*big.Int, error) {
	depositsKeys := make([]btcTransactions, len(depositsStrings))
	depositsRevealBlocks := make([]*big.Int, len(depositsStrings))

	for i, depositString := range depositsStrings {
		matched := depositsFormatRegexp.FindStringSubmatch(depositString)
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

		depositsKeys[i] = btcTransactions{
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
