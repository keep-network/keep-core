package tbtc

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func assembleDepositSweepTransaction(
	wallet wallet,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	deposits []*deposit,
	fee int64,
) (*bitcoin.Transaction, error) {
	if len(deposits) < 1 {
		return nil, fmt.Errorf("at least one deposit is required")
	}

	transaction := bitcoin.NewTransaction()
	totalInputsValue := int64(0)

	if walletMainUtxo != nil {
		transaction.AddInput(walletMainUtxo.Outpoint)
		totalInputsValue += walletMainUtxo.Value
	}

	for _, deposit := range deposits {
		transaction.AddInput(deposit.utxo.Outpoint)
		totalInputsValue += deposit.utxo.Value
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
