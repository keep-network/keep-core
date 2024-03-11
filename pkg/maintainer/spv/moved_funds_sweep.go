package spv

import (
	"bytes"
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// SubmitMovedFundsSweepProof prepares moved funds sweep proof for the given
// transaction and submits it to the on-chain contract. If the number of
// required confirmations is `0`, an error is returned.
func SubmitMovedFundsSweepProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
) error {
	return submitMovedFundsSweepProof(
		transactionHash,
		requiredConfirmations,
		btcChain,
		spvChain,
		bitcoin.AssembleSpvProof,
	)
}

func submitMovedFundsSweepProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
	spvProofAssembler spvProofAssembler,
) error {
	if requiredConfirmations == 0 {
		return fmt.Errorf(
			"provided required confirmations count must be greater than 0",
		)
	}

	transaction, proof, err := spvProofAssembler(
		transactionHash,
		requiredConfirmations,
		btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to assemble transaction spv proof: [%v]",
			err,
		)
	}

	mainUTXO, err := parseMovedFundsSweepTransactionInputs(
		btcChain,
		transaction,
	)
	if err != nil {
		return fmt.Errorf(
			"error while parsing transaction inputs: [%v]",
			err,
		)
	}

	if err := spvChain.SubmitMovedFundsSweepProofWithReimbursement(
		transaction,
		proof,
		mainUTXO,
	); err != nil {
		return fmt.Errorf(
			"failed to submit moved funds sweep proof with reimbursement: [%v]",
			err,
		)
	}

	return nil
}

// parseMovedFundsSweepTransactionInputs parses the transaction's inputs and returns
// the wallet's main UTXO.
func parseMovedFundsSweepTransactionInputs(
	btcChain bitcoin.Chain,
	transaction *bitcoin.Transaction,
) (bitcoin.UnspentTransactionOutput, error) {
	// Perform a sanity check: a moved funds sweep transaction must have one or
	// two inputs.
	if len(transaction.Inputs) != 1 && len(transaction.Inputs) != 2 {
		return bitcoin.UnspentTransactionOutput{}, fmt.Errorf(
			"moved funds sweep transaction has incorrect number of inputs",
		)
	}

	// If the transaction has only one input, it means the wallet does not have
	// the main UTXO yet. Return zero-filled value.
	if len(transaction.Inputs) == 1 {
		return bitcoin.UnspentTransactionOutput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: bitcoin.Hash{},
				OutputIndex:     0,
			},
			Value: 0,
		}, nil
	}

	// If the transaction has two inputs, the second input is the wallet's main
	// UTXO.
	input := transaction.Inputs[1]

	// Get data of the input transaction whose output is spent by the moved
	// funds sweep transaction.
	inputTx, err := btcChain.GetTransaction(input.Outpoint.TransactionHash)
	if err != nil {
		return bitcoin.UnspentTransactionOutput{}, fmt.Errorf(
			"cannot get input transaction data: [%v]",
			err,
		)
	}

	// Get the specific output spent by the moved funds sweep transaction.
	spentOutput := inputTx.Outputs[input.Outpoint.OutputIndex]

	// Build the main UTXO object based on available data.
	mainUtxo := bitcoin.UnspentTransactionOutput{
		Outpoint: input.Outpoint,
		Value:    spentOutput.Value,
	}

	return mainUtxo, nil
}

func getUnprovenMovedFundsSweepTransactions(
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

	events, err :=
		spvChain.PastMovingFundsCommitmentSubmittedEvents(
			&tbtc.MovingFundsCommitmentSubmittedEventFilter{
				StartBlock: startBlock,
			},
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past moving funds commitment submitted events: [%v]",
			err,
		)
	}

	// Any wallet that was among target wallets recently could have created
	// a moved funds sweep transaction.
	targetWallets := [][20]byte{}
	for _, event := range events {
		targetWallets = append(targetWallets, event.TargetWallets...)
	}

	// Some target wallets may appear on the list multiple times. It can happen
	// if multiple source wallets used the same target wallet. Make a list
	// of unique wallets.
	walletPublicKeyHashes := uniqueKeyHashes(targetWallets)

	unprovenMovedFundsSweepTransactions := []*bitcoin.Transaction{}

	// Should we check state of the wallet?
	// Should we check if the wallet has pending moved funds sweep request?

	for _, walletPublicKeyHash := range walletPublicKeyHashes {
		// When wallet makes a moved funds sweep transaction, it transfers
		// funds to itself. Therefore we can search all the transactions that
		// pay to the wallet's public key hash.
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
				isUnprovenMovedFundsSweepTransaction(
					transaction,
					walletPublicKeyHash,
					btcChain,
					spvChain,
				)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to check if transaction is an unproven moved "+
						"funds sweep transaction: [%v]",
					err,
				)
			}

			if isUnproven {
				unprovenMovedFundsSweepTransactions = append(
					unprovenMovedFundsSweepTransactions,
					transaction,
				)

				// A wallet can have only one unproven moved funds sweep
				// transaction at a time. If we found such transaction, we don't
				// have to look at this wallet's transactions anymore.
				break
			}
		}
	}

	return unprovenMovedFundsSweepTransactions, nil
}

func isUnprovenMovedFundsSweepTransaction(
	transaction *bitcoin.Transaction,
	walletPublicKeyHash [20]byte,
	btcChain bitcoin.Chain,
	spvChain Chain,
) (bool, error) {
	// A moved funds sweep transaction must have one or two inputs.
	if len(transaction.Inputs) != 1 && len(transaction.Inputs) != 2 {
		return false, nil
	}

	// A moved funds sweep transaction must have exactly one output.
	if len(transaction.Outputs) != 1 {
		return false, nil
	}

	// The first input must point to a pending moved funds sweep request.
	requestTransactionHash := transaction.Inputs[0].Outpoint.TransactionHash
	requestOutputIndex := transaction.Inputs[0].Outpoint.OutputIndex

	movedFundsSweepRequest, err := spvChain.GetMovedFundsSweepRequest(
		requestTransactionHash,
		requestOutputIndex,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to get moved funds sweep request: [%v]",
			err,
		)
	}

	if movedFundsSweepRequest.State != tbtc.MovedFundsStatePending {
		return false, nil
	}

	// If there is the second input it must refer to the current wallet's main
	// UTXO.
	if len(transaction.Inputs) == 2 {
		fundingTransactionHash := transaction.Inputs[1].Outpoint.TransactionHash
		fundingOutpointIndex := transaction.Inputs[1].Outpoint.OutputIndex

		isMainUtxo, err := isInputCurrentWalletsMainUTXO(
			fundingTransactionHash,
			fundingOutpointIndex,
			walletPublicKeyHash,
			btcChain,
			spvChain,
		)
		if err != nil {
			return false, fmt.Errorf(
				"failed to check if input is the main UTXO: [%v]",
				err,
			)
		}

		// The input is not the current main UTXO of the wallet.
		// The transaction cannot be an unproven moved funds sweep transaction.
		if !isMainUtxo {
			return false, nil
		}
	}

	// If the transaction is a moved funds sweep transaction the output must
	// transfer funds to the wallet itself.
	output := transaction.Outputs[0]

	p2pkh, err := bitcoin.PayToPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf(
			"failed to compute p2pkh script for transaction output: [%v]",
			err,
		)
	}
	p2wpkh, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf(
			"failed to compute p2wpkh script for transaction output: [%v]",
			err,
		)
	}

	return bytes.Equal(output.PublicKeyScript, p2pkh) ||
		bytes.Equal(output.PublicKeyScript, p2wpkh), nil
}
