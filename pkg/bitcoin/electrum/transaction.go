package electrum

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/v2/wire"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// decodeTransaction deserializes a transaction from the hexadecimal serialized
// string to a btcd message format.
func decodeTransaction(rawTx string) (*wire.MsgTx, error) {
	headerBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode a hex string: [%w]", err)
	}

	buf := bytes.NewBuffer(headerBytes)

	var t wire.MsgTx
	if err := t.Deserialize(buf); err != nil {
		return nil, fmt.Errorf("failed to deserialize a transaction: [%w]", err)
	}

	return &t, nil
}

// convertRawTransaction transforms a transaction provided in the hexadecimal serialized
// string to the format expected by the bitcoin.Chain interface.
func convertRawTransaction(rawTx string) (*bitcoin.Transaction, error) {
	t, err := decodeTransaction(rawTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode a transaction: [%w]", err)
	}

	result := &bitcoin.Transaction{
		Version:  int32(t.Version),
		Locktime: t.LockTime,
	}

	for _, vin := range t.TxIn {
		input := &bitcoin.TransactionInput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: bitcoin.Hash(vin.PreviousOutPoint.Hash),
				OutputIndex:     vin.PreviousOutPoint.Index,
			},
			SignatureScript: vin.SignatureScript,
			Sequence:        vin.Sequence,
		}

		result.Inputs = append(result.Inputs, input)
	}

	for _, vout := range t.TxOut {
		output := &bitcoin.TransactionOutput{
			Value:           vout.Value,
			PublicKeyScript: vout.PkScript,
		}

		result.Outputs = append(result.Outputs, output)
	}

	return result, nil
}
