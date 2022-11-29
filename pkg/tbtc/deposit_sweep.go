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
) (*bitcoin.TransactionBuilder, error) {
	if len(deposits) < 1 {
		return nil, fmt.Errorf("at least one deposit is required")
	}

	builder := bitcoin.NewTransactionBuilder()
	totalInputsValue := int64(0)

	if walletMainUtxo != nil {
		// TODO: Set proper scriptCode and witness.
		builder.AddInput(walletMainUtxo, nil, false)
		totalInputsValue += walletMainUtxo.Value
	}

	for _, deposit := range deposits {
		// TODO: Set proper scriptCode and witness.
		builder.AddInput(deposit.utxo, nil, false)
		totalInputsValue += deposit.utxo.Value
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)
	outputScript, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot compute output script: [%v]", err)
	}

	outputValue := totalInputsValue - fee

	builder.AddOutput(&bitcoin.TransactionOutput{
		Value:           outputValue,
		PublicKeyScript: outputScript,
	})

	return builder, nil
}
