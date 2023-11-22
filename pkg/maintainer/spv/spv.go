package spv

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/tbtc"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/maintainer/btcdiff"
)

var logger = log.Logger("keep-maintainer-spv")

// The length of the Bitcoin difficulty epoch in blocks.
const difficultyEpochLength = 2016

func Initialize(
	ctx context.Context,
	config Config,
	spvChain Chain,
	btcDiffChain btcdiff.Chain,
	btcChain bitcoin.Chain,
) {
	spvMaintainer := &spvMaintainer{
		config:       config,
		spvChain:     spvChain,
		btcDiffChain: btcDiffChain,
		btcChain:     btcChain,
	}

	go spvMaintainer.startControlLoop(ctx)
}

// proofTypes holds the information about proof types supported by the
// SPV maintainer.
var proofTypes = map[tbtc.WalletActionType]struct {
	unprovenTransactionsGetter unprovenTransactionsGetter
	transactionProofSubmitter  transactionProofSubmitter
}{
	tbtc.ActionDepositSweep: {
		unprovenTransactionsGetter: getUnprovenDepositSweepTransactions,
		transactionProofSubmitter:  SubmitDepositSweepProof,
	},
	tbtc.ActionRedemption: {
		unprovenTransactionsGetter: getUnprovenRedemptionTransactions,
		transactionProofSubmitter:  SubmitRedemptionProof,
	},
}

type spvMaintainer struct {
	config       Config
	spvChain     Chain
	btcDiffChain btcdiff.Chain
	btcChain     bitcoin.Chain
}

func (sm *spvMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting SPV maintainer")

	defer func() {
		logger.Info("stopping SPV maintainer")
	}()

	for {
		err := sm.maintainSpv(ctx)
		if err != nil {
			logger.Errorf(
				"error while maintaining SPV: [%v]; restarting maintainer",
				err,
			)
		}

		select {
		case <-time.After(sm.config.RestartBackoffTime):
		case <-ctx.Done():
			return
		}
	}
}

func (sm *spvMaintainer) maintainSpv(ctx context.Context) error {
	for {
		for action, v := range proofTypes {
			logger.Infof("starting [%s] proof task execution...", action)

			if err := sm.proveTransactions(
				v.unprovenTransactionsGetter,
				v.transactionProofSubmitter,
			); err != nil {
				return fmt.Errorf(
					"error while proving [%s] transactions: [%v]",
					action,
					err,
				)
			}

			logger.Infof("[%s] proof task completed", action)
		}

		logger.Infof(
			"proof tasks completed; next run in [%s]",
			sm.config.IdleBackoffTime,
		)

		select {
		case <-time.After(sm.config.IdleBackoffTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// unprovenTransactionsGetter is a type representing a function that is
// used to get unproven Bitcoin transactions.
type unprovenTransactionsGetter func(
	historyDepth uint64,
	transactionLimit int,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (
	[]*bitcoin.Transaction,
	error,
)

// transactionProofSubmitter is a type representing a function that is used
// to submit the constructed SPV proof to the host chain.
type transactionProofSubmitter func(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
) error

// proveTransactions gets unproven Bitcoin transactions using the provided
// unprovenTransactionsGetter, build the SPV proofs, and submits them using
// the provided transactionProofSubmitter.
func (sm *spvMaintainer) proveTransactions(
	unprovenTransactionsGetter unprovenTransactionsGetter,
	transactionProofSubmitter transactionProofSubmitter,
) error {
	transactions, err := unprovenTransactionsGetter(
		sm.config.HistoryDepth,
		sm.config.TransactionLimit,
		sm.btcChain,
		sm.spvChain,
	)
	if err != nil {
		return fmt.Errorf("failed to get unproven transactions: [%v]", err)
	}

	logger.Infof("found [%d] unproven transaction(s)", len(transactions))

	for _, transaction := range transactions {
		// Print the transaction in the same endianness as block explorers do.
		transactionHashStr := transaction.Hash().Hex(bitcoin.ReversedByteOrder)

		logger.Infof(
			"proceeding with proof for transaction [%s]",
			transactionHashStr,
		)

		_, accumulatedConfirmations, requiredConfirmations, err := getProofInfo(
			transaction.Hash(),
			sm.btcChain,
			sm.spvChain,
			sm.btcDiffChain,
		)
		if err != nil {
			return fmt.Errorf("failed to get proof info: [%v]", err)
		}

		if accumulatedConfirmations < requiredConfirmations {
			// Skip the transaction as it has not accumulated enough
			// confirmations. It will be proven later.
			logger.Infof(
				"skipped proving transaction [%s]; transaction "+
					"has [%v/%v] confirmations",
				transactionHashStr,
				accumulatedConfirmations,
				requiredConfirmations,
			)
			continue
		}

		err = transactionProofSubmitter(
			transaction.Hash(),
			requiredConfirmations,
			sm.btcChain,
			sm.spvChain,
		)
		if err != nil {
			return err
		}

		logger.Infof(
			"successfully submitted proof for transaction [%s]",
			transactionHashStr,
		)
	}

	logger.Infof("finished round of proving transactions")

	return nil
}

func isInputCurrentWalletsMainUTXO(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	walletPublicKeyHash [20]byte,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (bool, error) {
	// Get the transaction the input originated from to calculate the input value.
	previousTransaction, err := btcChain.GetTransaction(fundingTxHash)
	if err != nil {
		return false, fmt.Errorf("failed to get previous transaction: [%v]", err)
	}
	fundingOutputValue := previousTransaction.Outputs[fundingOutputIndex].Value

	// Assume the input is the main UTXO and calculate hash.
	mainUtxoHash := spvChain.ComputeMainUtxoHash(&bitcoin.UnspentTransactionOutput{
		Outpoint: &bitcoin.TransactionOutpoint{
			TransactionHash: fundingTxHash,
			OutputIndex:     fundingOutputIndex,
		},
		Value: fundingOutputValue,
	})

	// Get the wallet and check if its main UTXO matches the calculated hash.
	wallet, err := spvChain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf("failed to get wallet: [%v]", err)
	}

	return bytes.Equal(mainUtxoHash[:], wallet.MainUtxoHash[:]), nil
}

// getProofInfo returns information about the SPV proof. It includes the
// information whether the transaction proof range is within the previous and
// current difficulty epochs as seen by the relay, the accumulated number of
// confirmations and the required number of confirmations.
func getProofInfo(
	transactionHash bitcoin.Hash,
	btcChain bitcoin.Chain,
	spvChain Chain,
	btcDiffChain btcdiff.Chain,
) (
	bool, uint, uint, error,
) {
	accumulatedConfirmations, err := btcChain.GetTransactionConfirmations(
		transactionHash,
	)
	if err != nil {
		return false, 0, 0, fmt.Errorf(
			"failed to get transaction confirmations: [%v]",
			err,
		)
	}

	txProofDifficultyFactor, err := spvChain.TxProofDifficultyFactor()
	if err != nil {
		return false, 0, 0, fmt.Errorf(
			"failed to get transaction proof difficulty factor: [%v]",
			err,
		)
	}

	// Testnet proofs often require longer header chains due to the frequent
	// difficulty changes on Bitcoin testnet.
	requiredConfirmations := uint(txProofDifficultyFactor.Uint64()) + 3

	return false, accumulatedConfirmations, requiredConfirmations, nil
}

// walletEvent is a type constraint representing wallet-related chain events.
type walletEvent interface {
	WalletPublicKeyHash() [20]byte
}

// uniqueWalletPublicKeyHashes parses the list of wallet-related events and
// returns a list of unique wallet public key hashes.
func uniqueWalletPublicKeyHashes[T walletEvent](events []T) [][20]byte {
	cache := make(map[string]struct{})
	var publicKeyHashes [][20]byte

	for _, event := range events {
		key := event.WalletPublicKeyHash()
		strKey := hex.EncodeToString(key[:])

		// Check for uniqueness
		if _, exists := cache[strKey]; !exists {
			cache[strKey] = struct{}{}
			publicKeyHashes = append(publicKeyHashes, key)
		}
	}

	return publicKeyHashes
}

// spvProofAssembler is a type representing a function that is used
// to assemble an SPV proof for the given transaction hash and confirmations
// count.
type spvProofAssembler func(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
) (*bitcoin.Transaction, *bitcoin.SpvProof, error)
