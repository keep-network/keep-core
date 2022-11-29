package bitcoin

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"math/big"
)

// TransactionBuilder is a component that is responsible for the whole
// transaction creation process. It assembles an unsigned transaction,
// prepares it for signing, and applies the given signatures in order to
// produce a full-fledged signed transaction that can be spread across
// the Bitcoin network.
type TransactionBuilder struct {
	transaction *wire.MsgTx
	sigHashArgs []*inputSigHashArgs
}

// NewTransactionBuilder constructs a new TransactionBuilder instance.
func NewTransactionBuilder() *TransactionBuilder {
	transaction := wire.NewMsgTx(wire.TxVersion)
	transaction.LockTime = 0

	return &TransactionBuilder{
		transaction: transaction,
		sigHashArgs: make([]*inputSigHashArgs, 0),
	}
}

// AddInput adds a new unsigned transaction's input based on the provided UTXO
// and stores some additional data required to construct a proper input's
// signature hash needed to produce a signature unlocking that UTXO.
//
// The scriptCode parameter represents the script that is actually executed
// while unlocking the given UTXO. The scriptCode depends on the script type
// that was used to lock the given UTXO. Specifically:
// - If the UTXO was locked with a P2PKH/P2WPKH script, the scriptCode is
//   equivalent to the original non-SegWit PublicKeyScript format of the P2PKH
//   output, i.e. `<0x76a914> <20-byte PKH> <0x88ac>`. This applies to P2WPKH
//   outputs as well, even though they use SegWit PublicKeyScript format, i.e.
//   `<0x0014> <20-byte PKH>`. This is because the P2WPKH SegWit PublicKeyScript
//   format is actually translated to the P2PKH non-SegWit PublicKeyScript
//   format. However, this function accepts both PublicKeyScript formats and do
//   necessary translations underneath.
// - If the UTXO was locked with a P2SH/P2WSH script, the scriptCode is
//   equivalent to the plain-text redeem script whose hash is included
//   in the PublicKeyScript format.
//
// The witness parameter tells whether the added input points to a SegWit
// output or not.
func (tb *TransactionBuilder) AddInput(
	utxo *UnspentTransactionOutput,
	scriptCode []byte,
	witness bool,
) {
	hash := chainhash.Hash(utxo.Outpoint.TransactionHash)
	outpoint := wire.NewOutPoint(&hash, utxo.Outpoint.OutputIndex)

	// Deliberately set both `signatureScript` and `witness` arguments to nil
	// because at this point, the input does not contain any signature data.
	tb.transaction.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

	tb.sigHashArgs = append(tb.sigHashArgs, &inputSigHashArgs{
		value:      utxo.Value,
		scriptCode: scriptCode,
		witness:    witness,
	})
}

// AddOutput adds a new transaction's output.
func (tb *TransactionBuilder) AddOutput(output *TransactionOutput) {
	tb.transaction.AddTxOut(wire.NewTxOut(output.Value, output.PublicKeyScript))
}

// SigHashes computes the signature hashes for all transaction inputs. Elements
// of the returned slice are ordered in the same way as the transaction inputs
// they correspond to. That is, an element at the given index matches the input
// with the same index.
func (tb *TransactionBuilder) SigHashes() ([]*big.Int, error) {
	sigHashes := make([]*big.Int, len(tb.transaction.TxIn))

	// Calculation of sighashes for witness inputs can be faster as common
	// sighash fragments can be pre-computed upfront and reused.
	witnessSigHashFragments := txscript.NewTxSigHashes(tb.transaction)

	for i := range tb.transaction.TxIn {
		sigHashArgs := tb.sigHashArgs[i]

		var sigHashBytes []byte
		var err error

		if sigHashArgs.witness {
			sigHashBytes, err = txscript.CalcWitnessSigHash(
				sigHashArgs.scriptCode,
				witnessSigHashFragments,
				txscript.SigHashAll,
				tb.transaction,
				i,
				sigHashArgs.value,
			)
		} else {
			sigHashBytes, err = txscript.CalcSignatureHash(
				sigHashArgs.scriptCode,
				txscript.SigHashAll,
				tb.transaction,
				i,
			)
		}

		if err != nil {
			return nil, fmt.Errorf(
				"cannot calculate sighash for input [%v]: [%v]",
				i,
				err,
			)
		}

		sigHashes[i] = new(big.Int).SetBytes(sigHashBytes)
	}

	return sigHashes, nil
}

// inputSigHashArgs is a helper structure holding some arguments required to
// compute a sighash for the given input.
type inputSigHashArgs struct {
	value      int64
	scriptCode []byte
	witness    bool
}
