package spv

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// SubmitRedemptionProof prepares redemption proof for the given transaction
// and submits it to the on-chain contract. If the number of required
// confirmations is `0`, an error is returned.
//
// TODO: Expose this function through the maintainer-cli tool.
func SubmitRedemptionProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
) error {
	panic("not implemented yet")
}

func getUnprovenRedemptionTransactions(
	historyDepth uint64,
	transactionLimit int,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (
	[]*bitcoin.Transaction,
	error,
) {
	blockCounter, err := spvChain.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("failed to get block counter: [%v]", err)
	}

	currentBlock, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get current block: [%v]", err)
	}

	// Calculate the starting block of the range in which the events will be
	// searched for.
	startBlock := currentBlock - historyDepth

	redemptionProposals, err :=
		spvChain.PastRedemptionProposalSubmittedEvents(
			&tbtc.RedemptionProposalSubmittedEventFilter{
				StartBlock: startBlock,
			},
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past redemption proposal submitted events: [%v]",
			err,
		)
	}

	// There will often be multiple events emitted for a single wallet. Prepare
	// a list of unique wallet public key hashes.
	walletPublicKeyHashes := uniqueWalletPublicKeyHashes(
		redemptionProposals,
	)

	var unprovenRedemptionTransactions []*bitcoin.Transaction

	for _, walletPublicKeyHash := range walletPublicKeyHashes {
		wallet, err := spvChain.GetWallet(walletPublicKeyHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet: [%v]", err)
		}

		if wallet.State != tbtc.StateLive &&
			wallet.State != tbtc.StateMovingFunds {
			// The wallet can only submit redemption proofs if it's `Live` or
			// `MovingFunds`. If the state is different skip it.
			logger.Infof(
				"skipped proving redemption transactions for wallet [%x] "+
					"because of wallet state [%v]",
				walletPublicKeyHash,
				wallet.State,
			)
			continue
		}

		walletTransactions, err := btcChain.GetTransactionsForPublicKeyHash(
			walletPublicKeyHash,
			transactionLimit,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get transactions for wallet: [%v]",
				err,
			)
		}

		for _, transaction := range walletTransactions {
			isUnproven, err :=
				isUnprovenRedemptionTransaction(
					transaction,
					walletPublicKeyHash,
					btcChain,
					spvChain,
				)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to check if transaction is an unproven redemption "+
						"transaction: [%v]",
					err,
				)
			}

			if isUnproven {
				unprovenRedemptionTransactions = append(
					unprovenRedemptionTransactions,
					transaction,
				)
			}
		}
	}

	return unprovenRedemptionTransactions, nil
}

func isUnprovenRedemptionTransaction(
	transaction *bitcoin.Transaction,
	walletPublicKeyHash [20]byte,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (bool, error) {
	// If the transaction does not have exactly one input, it cannot be a
	// redemption transaction.
	if len(transaction.Inputs) != 1 {
		return false, nil
	}

	singleInput := transaction.Inputs[0]

	// Check whether the single input is the current wallet main UTXO.
	isMainUtxo, err := isInputCurrentWalletsMainUTXO(
		singleInput.Outpoint.TransactionHash,
		singleInput.Outpoint.OutputIndex,
		walletPublicKeyHash,
		btcChain,
		spvChain,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check if input is the main UTXO",
		)
	}

	// If the single input is not the current main UTXO of the wallet, the
	// transaction is either a redemption transaction that is already
	// proven or it's not a redemption transaction at all.
	if !isMainUtxo {
		return false, nil
	}

	changeFound := false

	// Look at the transaction's outputs. All the outputs must be pending
	// redemption requests, except for one optional change output.
	for _, output := range transaction.Outputs {
		// First, check if the given output is a change (if it wasn't
		// found yet).
		if !changeFound {
			isChange, err := isWalletChangeOutput(walletPublicKeyHash, output)
			if err != nil {
				return false, fmt.Errorf(
					"failed to check if output is wallet change: [%v]",
					err,
				)
			}

			if isChange {
				changeFound = true
				continue
			}
		}

		// If the given output is not a change, it must be a pending redemption
		// request.
		_, err := spvChain.GetPendingRedemptionRequest(
			walletPublicKeyHash,
			output.PublicKeyScript,
		)
		if err != nil {
			if errors.Is(err, tbtc.ErrPendingRedemptionRequestNotFound) {
				// This output is neither a change nor a pending request.
				// That means this is not a redemption transaction.
				return false, nil
			} else {
				return false, fmt.Errorf(
					"failed to get pending redemption request: [%w]",
					err,
				)
			}
		}
	}

	return true, nil
}

func isWalletChangeOutput(
	walletPublicKeyHash [20]byte,
	output *bitcoin.TransactionOutput,
) (bool, error) {
	walletP2PKH, err := bitcoin.PayToPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf("cannot construct P2PKH for wallet: [%v]", err)
	}
	walletP2WPKH, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf("cannot construct P2WPKH for wallet: [%v]", err)
	}

	script := output.PublicKeyScript
	return bytes.Equal(script, walletP2PKH) || bytes.Equal(script, walletP2WPKH), nil
}
