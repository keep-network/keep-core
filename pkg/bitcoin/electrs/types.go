package electrs

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type transaction struct {
	TxID     string `json:"txid"`
	Version  int32  `json:"version"`
	Locktime uint32 `json:"locktime"`
	Vin      []vin  `json:"vin"`
	Vout     []vout `json:"vout"`
}

type vin struct {
	TxID      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig string `json:"scriptsig"`
	Sequence  uint32 `json:"sequence"`
}

type vout struct {
	Value        int64  `json:"value"`
	ScriptPubkey string `json:"scriptpubkey"`
}

func (t *transaction) convert() (bitcoin.Transaction, error) {
	result := bitcoin.Transaction{
		Version:  t.Version,
		Locktime: t.Locktime,
	}

	for _, vin := range t.Vin {
		txHash, err := bitcoin.NewHashFromString(vin.TxID, bitcoin.ReversedByteOrder)
		if err != nil {
			return result, fmt.Errorf(
				"failed to decode transaction hash from [%s]: %w",
				vin.TxID,
				err,
			)
		}

		input := &bitcoin.TransactionInput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: txHash,
				OutputIndex:     vin.Vout,
			},
			SignatureScript: []byte(vin.ScriptSig),
			Sequence:        vin.Sequence,
		}

		result.Inputs = append(result.Inputs, input)
	}

	for _, vout := range t.Vout {
		output := &bitcoin.TransactionOutput{
			Value:           vout.Value,
			PublicKeyScript: []byte(vout.ScriptPubkey),
		}

		result.Outputs = append(result.Outputs, output)
	}

	return result, nil
}
