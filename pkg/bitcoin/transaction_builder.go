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
	chain       Chain
	transaction *wire.MsgTx
	sigHashArgs []*inputSigHashArgs
}

// NewTransactionBuilder constructs a new TransactionBuilder instance.
func NewTransactionBuilder(chain Chain) *TransactionBuilder {
	transaction := wire.NewMsgTx(wire.TxVersion)
	transaction.LockTime = 0

	return &TransactionBuilder{
		chain:       chain,
		transaction: transaction,
		sigHashArgs: make([]*inputSigHashArgs, 0),
	}
}

// AddPublicKeyHashInput adds an unsigned input pointing to a UTXO locked
// using a P2PKH or P2WPKH script.
func (tb *TransactionBuilder) AddPublicKeyHashInput(
	utxo *UnspentTransactionOutput,
) error {
	utxoScript, err := tb.getScript(utxo)
	if err != nil {
		return fmt.Errorf(
			"cannot get locking script for UTXO pointed "+
				"by the input: [%v]",
			err,
		)
	}

	class := txscript.GetScriptClass(utxoScript)
	isPublicKeyHashScript := class == txscript.PubKeyHashTy ||
		class == txscript.WitnessV0PubKeyHashTy
	if !isPublicKeyHashScript {
		return fmt.Errorf(
			"UTXO pointed by the input is not P2PKH/P2WPKH",
		)
	}

	// The UTXO was locked using a P2PKH/P2WPKH script so, the scriptCode
	// required to build the sighash is equivalent to that script. Worth
	// noting that the P2WPKH script is actually converted to the P2PKH script
	// when used as a scriptCode, according to BIP-0143. For reference see,
	// https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki#specification.
	// That conversion is handled within the `txscript.CalcWitnessSigHash` call.
	sigHashArgs := &inputSigHashArgs{
		value:      utxo.Value,
		scriptCode: utxoScript,
		witness:    txscript.IsWitnessProgram(utxoScript),
	}

	tb.addInput(utxo, sigHashArgs)

	return nil
}

// AddScriptHashInput adds an unsigned input pointing to a UTXO locked
// using a P2SH or P2WSH script. This function also requires the plain-text
// redeemScript whose hash was used to build the P2SH/P2WSH locking script.
func (tb *TransactionBuilder) AddScriptHashInput(
	utxo *UnspentTransactionOutput,
	redeemScript Script,
) error {
	utxoScript, err := tb.getScript(utxo)
	if err != nil {
		return fmt.Errorf(
			"cannot get locking script for UTXO pointed "+
				"by the input: [%v]",
			err,
		)
	}

	class := txscript.GetScriptClass(utxoScript)
	isPublicKeyHashScript := class == txscript.ScriptHashTy ||
		class == txscript.WitnessV0ScriptHashTy
	if !isPublicKeyHashScript {
		return fmt.Errorf(
			"UTXO pointed by the input is not P2SH/P2WSH",
		)
	}

	// The UTXO was locked using a P2SH/P2WSH script so, the scriptCode required
	// to build the sighash is equivalent to the plain-text redeem script whose
	// hash is included in the P2SH/P2WSH script.
	sigHashArgs := &inputSigHashArgs{
		value:      utxo.Value,
		scriptCode: redeemScript,
		witness:    txscript.IsWitnessProgram(utxoScript),
	}

	tb.addInput(utxo, sigHashArgs)

	return nil
}

// getScript gets the locking script (PublicKeyScript) for the given unspent
// transaction output.
func (tb *TransactionBuilder) getScript(
	utxo *UnspentTransactionOutput,
) (Script, error) {
	hash := utxo.Outpoint.TransactionHash
	transaction, err := tb.chain.GetTransaction(hash)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get transaction with hash [%s]: [%v]",
			hash.String(InternalByteOrder),
			err,
		)
	}

	return transaction.Outputs[utxo.Outpoint.OutputIndex].PublicKeyScript, nil
}

// addInput adds the given input along with their sighash arguments.
func (tb *TransactionBuilder) addInput(
	utxo *UnspentTransactionOutput,
	sigHashArgs *inputSigHashArgs,
) {
	hash := chainhash.Hash(utxo.Outpoint.TransactionHash)
	outpoint := wire.NewOutPoint(&hash, utxo.Outpoint.OutputIndex)

	// Deliberately set both `signatureScript` and `witness` arguments to nil
	// because at this point, the input does not contain any signature data.
	tb.transaction.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

	tb.sigHashArgs = append(tb.sigHashArgs, sigHashArgs)
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

// TotalInputsValue returns the total value of transaction inputs.
func (tb *TransactionBuilder) TotalInputsValue() int64 {
	totalInputsValue := int64(0)

	for _, args := range tb.sigHashArgs {
		totalInputsValue += args.value
	}

	return totalInputsValue
}

// inputSigHashArgs is a helper structure holding some arguments required to
// compute a sighash for the given input.
type inputSigHashArgs struct {
	// value denotes the satoshi value of the UTXO pointed by the given input.
	value int64
	// scriptCode is a component of the input's sighash and is the script that
	// is actually executed while unlocking the given UTXO. The scriptCode
	// depends on the script type that was used to lock the given UTXO.
	scriptCode []byte
	// witness denotes whether the given input point's to a UTXO locked using
	// a witness script.
	witness bool
}
