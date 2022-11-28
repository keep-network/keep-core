package tbtc

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func assembleDepositSweepTransaction(
	wallet wallet,
	utxos []*bitcoin.UnspentTransactionOutput,
	fee int64,
) (*bitcoin.Transaction, error) {
	if len(utxos) < 1 {
		return nil, fmt.Errorf("at least one deposit is required")
	}

	// TODO: Handle the wallet's main UTXO.

	transaction := bitcoin.NewTransaction()
	totalInputsValue := int64(0)

	for _, utxo := range utxos {
		transaction.AddInput(utxo.Outpoint)

		// TODO: Add deposit's value to the totalInputsValue.
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)
	outputScript, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot compute output script: [%v]", err)
	}

	outputValue := totalInputsValue - fee

	transaction.AddOutput(outputScript, outputValue)

	return transaction, nil
}
