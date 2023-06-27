package spv

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// SubmitDepositSweepProof prepares deposit sweep proof for the given
// transaction and submits it to the onchain contract. If the number of required
// confirmations is `0`, the value is read from the chain.
func SubmitDepositSweepProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	btcChain bitcoin.Chain,
	spvChain Chain,
) error {
	if requiredConfirmations == 0 {
		// The caller did not specify the number of confirmations. Use the
		// default value stored onchain.
		txProofDifficulty, err := spvChain.TxProofDifficultyFactor()
		if err != nil {
			return fmt.Errorf(
				"failed to get transaction proof difficulty factor: [%v]",
				err,
			)
		}

		requiredConfirmations = uint(txProofDifficulty.Int64())
	}

	transaction, proof, err := bitcoin.AssembleSpvProof(
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

	mainUTXO, vault, err := parseTransactionInputs(
		btcChain,
		spvChain,
		transaction,
	)
	if err != nil {
		return fmt.Errorf(
			"error while parsing transaction inputs: [%v]",
			err,
		)
	}

	if err := spvChain.SubmitDepositSweepProofWithReimbursement(
		transaction,
		proof,
		mainUTXO,
		vault,
	); err != nil {
		return fmt.Errorf(
			"failed to submit deposit sweep proof with reimbursement: [%v]",
			err,
		)
	}

	return nil
}

// parseTransactionInputs parses the transaction's inputs and returns the main
// UTXO and the vault.
func parseTransactionInputs(
	btcChain bitcoin.Chain,
	spvChain Chain,
	transaction *bitcoin.Transaction,
) (
	bitcoin.UnspentTransactionOutput,
	common.Address,
	error,
) {
	// Represents the main UTXO of the deposit sweep transaction. Nil if there
	// was no main UTXO.
	var mainUTXO *bitcoin.UnspentTransactionOutput = nil

	// Stores the vault address of the deposits. Each deposit should have the
	// same value of vault. The zero-filled value indicates there was no vault
	// value set for the deposits.
	var vault = common.Address{}

	// This flag checks if at least one deposit input has been found during
	// deposit processing.
	var depositAlreadyProcessed = false

	// Perform a sanity check: a deposit sweep transaction must have exactly one
	// output.
	if len(transaction.Outputs) != 1 {
		return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
			"deposit sweep transaction has more than one output",
		)
	}

	for _, input := range transaction.Inputs {
		outpointTransactionHash := input.Outpoint.TransactionHash
		outpointIndex := input.Outpoint.OutputIndex

		previousTransaction, err := btcChain.GetTransaction(
			outpointTransactionHash,
		)
		if err != nil {
			return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
				"failed to get previous transaction: [%v]",
				err,
			)
		}

		publicKeyScript := previousTransaction.Outputs[outpointIndex].PublicKeyScript
		value := previousTransaction.Outputs[outpointIndex].Value
		scriptClass := txscript.GetScriptClass(publicKeyScript)

		if scriptClass == txscript.PubKeyHashTy ||
			scriptClass == txscript.WitnessV0PubKeyHashTy {
			// The input is P2PKH or P2WPKH, so we found main UTXO. There should
			// be at most one main UTXO. If any input of this kind has already
			// been found, report an error.
			if mainUTXO == nil {
				mainUTXO = &bitcoin.UnspentTransactionOutput{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: outpointTransactionHash,
						OutputIndex:     outpointIndex,
					},
					Value: value,
				}
			} else {
				return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
					"deposit sweep transaction has more than one non-deposit " +
						"inputs",
				)
			}
		} else if scriptClass == txscript.ScriptHashTy ||
			scriptClass == txscript.WitnessV0ScriptHashTy {
			// The input is P2SH or P2WSH, so we found a deposit input. All
			// the deposits should have the same vault set or no vault at all.
			// If the vault if different than the vault from any previous
			// deposit input, report an error.
			deposit, err := spvChain.GetDepositRequest(
				outpointTransactionHash,
				outpointIndex,
			)
			if err != nil {
				return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
					"failed to get deposit request: [%v]",
					err,
				)
			}

			// Call successful, but deposit not found.
			if deposit.RevealedAt.Unix() == 0 {
				return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
					"deposit not found: [%v]",
					err,
				)
			}

			if depositAlreadyProcessed {
				if vault != convertVaultAddress(deposit.Vault) {
					return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
						"swept deposits have different vaults",
					)
				}
			} else {
				// The first deposit input has been encountered, save the vault
				// and require subsequent deposit inputs to have the same vault.
				vault = convertVaultAddress(deposit.Vault)
				depositAlreadyProcessed = true
			}
		} else {
			// The type of the input is neither P2PKH, P2WPKH, P2SH or P2WSH.
			// Report an error.
			return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
				"deposit sweep transaction has incorrect input types",
			)
		}
	}

	// If none of the input was main UTXO, return zero-filled main UTXO.
	if mainUTXO == nil {
		mainUTXO = &bitcoin.UnspentTransactionOutput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: bitcoin.Hash{},
				OutputIndex:     0,
			},
			Value: 0,
		}
	}

	return *mainUTXO, vault, nil
}

// convertVaultAddress converts the vault's address from the chain's address
// to the common address. If the chain's address is nil, the returned value
// is zero-filled.
func convertVaultAddress(vault *chain.Address) common.Address {
	if vault == nil {
		return common.Address{}
	}

	return common.HexToAddress(string(*vault))
}
