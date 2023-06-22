package spv

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

var logger = log.Logger("keep-maintainer-spv")

func Initialize(
	ctx context.Context,
	config Config,
	chain Chain,
	btcChain bitcoin.Chain,
) {
	spvMaintainer := &spvMaintainer{
		config:   config,
		chain:    chain,
		btcChain: btcChain,
	}

	go spvMaintainer.startControlLoop(ctx)
}

type spvMaintainer struct {
	config   Config
	chain    Chain
	btcChain bitcoin.Chain
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
		case <-time.After(sm.config.RestartBackOffTime):
		case <-ctx.Done():
			return
		}
	}
}

func (sm *spvMaintainer) maintainSpv(ctx context.Context) error {
	logger.Infof("Maintaining SPV proof...")

	for {
		if err := sm.proveDepositSweepTransactions(); err != nil {
			return fmt.Errorf(
				"error while proving deposit sweep transactions: [%v]",
				err,
			)
		}

		// TODO: Add proving of other type of SPV transactions: redemption
		// transactions, moving funds transaction, etc.

		select {
		case <-time.After(sm.config.IdleBackOffTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (sm *spvMaintainer) proveDepositSweepTransactions() error {
	depositSweepTransactions, err := sm.getUnprovenDepositSweepTransactions()
	if err != nil {
		return fmt.Errorf(
			"failed to get unproven deposit sweep transactions: [%v]",
			err,
		)
	}

	fmt.Println("depositSweepTransactions: ", depositSweepTransactions)

	// TODO: Assemble the proof and submit to the Bridge
	return nil
}

func (sm *spvMaintainer) getUnprovenDepositSweepTransactions() (
	[]*bitcoin.Transaction,
	error,
) {
	// TODO: Limit how far in the past we are looking for the events.
	//       Possibly store latest checked height in memory or file.
	depositSweepTransactionProposals, err :=
		sm.chain.PastDepositSweepProposalSubmittedEvents(nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past deposit sweep proposal submitted events: [%v]",
			err,
		)
	}

	// There will often be multiple events emitted for a single wallet. Prepare
	// a list of unique wallet public key hashes.
	walletPublicKeyHashes := uniqueWalletPublicKeyHashes(
		depositSweepTransactionProposals,
	)

	unprovenDepositSweepTransactions := []*bitcoin.Transaction{}

	for _, walletPublicKeyHash := range walletPublicKeyHashes {
		// TODO: Should we check the wallet's state before attempting to submit
		//       the deposit sweep proof?
		// TODO: Think what the limit of transactions retrieved should be.
		walletTransactions, err := sm.btcChain.GetTransactionsForPublicKeyHash(
			walletPublicKeyHash,
			5,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get transactions for wallet: [%v]",
				err,
			)
		}

		for _, transaction := range walletTransactions {
			isUnprovenDepositSweepTransaction, err :=
				sm.isUnprovenDepositSweepTransaction(
					transaction,
					walletPublicKeyHash,
				)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to check if transaction is an unproven deposit sweep "+
						"transaction: [%v]",
					err,
				)
			}

			if isUnprovenDepositSweepTransaction {
				unprovenDepositSweepTransactions = append(
					unprovenDepositSweepTransactions,
					transaction,
				)
			}
		}
	}

	return unprovenDepositSweepTransactions, nil
}

func (sm *spvMaintainer) isUnprovenDepositSweepTransaction(
	transaction *bitcoin.Transaction,
	walletPublicKeyHash [20]byte,
) (bool, error) {
	// If the transaction does not have exactly one output, it cannot be a
	// deposit sweep transaction.
	if len(transaction.Outputs) != 1 {
		return false, nil
	}

	depositFound := false
	mainUtxoFound := false

	// Look at the transaction's inputs. All the inputs must be deposit inputs,
	// except for one input which can be the main UTXO.
	for _, input := range transaction.Inputs {
		fundingTransactionHash := input.Outpoint.TransactionHash
		fundingOutpointIndex := input.Outpoint.OutputIndex

		// Check if the input is a deposit input.
		deposit, err := sm.chain.Deposits(
			fundingTransactionHash,
			fundingOutpointIndex,
		)
		if err != nil {
			return false, fmt.Errorf("failed to get a deposit: [%v]", err)
		}

		if deposit.RevealedAt.Equal(time.Unix(0, 0)) {
			// The input is not a deposit input. The transaction can still be
			// a deposit sweep transaction, since the input may be the main UTXO.

			// Since the transaction can have at most one main UTXO input,
			// return immediately if the main UTXO has already been seen.
			if mainUtxoFound {
				return false, nil
			}

			// Check if the input is the main UTXO.
			isMainUtxo, err := sm.isInputWalletsMainUTXO(
				fundingTransactionHash,
				fundingOutpointIndex,
				walletPublicKeyHash,
			)
			if err != nil {
				return false, fmt.Errorf("failed to check is input is the main UTXO")
			}

			// The input was not the main UTXO.
			if !isMainUtxo {
				return false, nil
			}

			// The input was the main UTXO.
			mainUtxoFound = true
		} else {
			// The input is a deposit input. Check if it swept or not.
			if deposit.SweptAt.Equal(time.Unix(0, 0)) {
				// The input is a deposit and it's unswept.
				depositFound = true
			} else {
				// The input is a deposit, but it's already swept.
				// The transaction is probably a deposit sweep transaction, but
				// it's already proven.
				return false, nil
			}
		}
	}

	// All the inputs are either deposit or main UTXO inputs. As the final check
	// verify if at least one of them was a deposit input.
	return depositFound, nil
}

func (sm *spvMaintainer) isInputWalletsMainUTXO(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	walletPublicKeyHash [20]byte,
) (bool, error) {
	// Get the transaction the input originated from to calculate the input value.
	previousTransaction, err := sm.btcChain.GetTransaction(fundingTxHash)
	if err != nil {
		return false, fmt.Errorf("failed to get previous transaction: [%v]", err)
	}
	fundingOutputValue := previousTransaction.Outputs[fundingOutputIndex].Value

	// Assume the input is the main UTXO and calculate hash.
	mainUtxoHash := sm.chain.ComputeMainUtxoHash(&bitcoin.UnspentTransactionOutput{
		Outpoint: &bitcoin.TransactionOutpoint{
			TransactionHash: fundingTxHash,
			OutputIndex:     fundingOutputIndex,
		},
		Value: fundingOutputValue,
	})

	// Get the wallet and check if its main UTXO matches the calculated hash.
	// TODO: `GetWallet` may return an error if the wallet is not found.
	//       Should we return an error in such case or continue?
	wallet, err := sm.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf("failed to get wallet: [%v]", err)
	}

	return bytes.Equal(mainUtxoHash[:], wallet.MainUtxoHash[:]), nil
}

// uniqueWalletPublicKeyHashes parses the list of events and returns a list of
// unique wallet public key hashes.
func uniqueWalletPublicKeyHashes(
	events []*tbtc.DepositSweepProposalSubmittedEvent,
) [][20]byte {
	cache := make(map[string]struct{})
	var publicKeyHashes [][20]byte

	for _, event := range events {
		key := event.Proposal.WalletPublicKeyHash
		strKey := hex.EncodeToString(key[:])

		// Check for uniqueness
		if _, exists := cache[strKey]; !exists {
			cache[strKey] = struct{}{}
			publicKeyHashes = append(publicKeyHashes, key)
		}
	}

	return publicKeyHashes
}
