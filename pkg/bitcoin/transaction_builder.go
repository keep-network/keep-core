package bitcoin

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
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

	hash := chainhash.Hash(utxo.Outpoint.TransactionHash)
	outpoint := wire.NewOutPoint(&hash, utxo.Outpoint.OutputIndex)

	// Deliberately set both `signatureScript` and `witness` arguments to nil
	// because at this point, the input does not contain any signature data.
	tb.transaction.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

	tb.sigHashArgs = append(tb.sigHashArgs, sigHashArgs)

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

	hash := chainhash.Hash(utxo.Outpoint.TransactionHash)
	outpoint := wire.NewOutPoint(&hash, utxo.Outpoint.OutputIndex)

	// Signature data required to unlock a P2SH/P2WSH UTXO needs the plain-text
	// redeem script to be placed as the last item of the `witness` field for
	// P2WSH or the `signatureScript` field for P2SH. Here we prepare to fulfill
	// that requirement by putting the redeem script to the correct field and
	// let the AddSignatures method prepend it with the actual signature
	// and public key.
	if sigHashArgs.witness {
		tb.transaction.AddTxIn(wire.NewTxIn(outpoint, nil, [][]byte{redeemScript}))
	} else {
		tb.transaction.AddTxIn(wire.NewTxIn(outpoint, redeemScript, nil))
	}

	tb.sigHashArgs = append(tb.sigHashArgs, sigHashArgs)

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

// AddSignatures adds signature data for transaction inputs and returns a
// signed Transaction instance. The signatures slice should have the same
// length as the transaction's input vector. The signature with given index
// should correspond to the input with the same index. Each signature
// should also contain a public key that can be used for verification, i.e.
// this should be the public key that corresponds to the private key used
// to produce the given signature.
func (tb *TransactionBuilder) AddSignatures(
	signatures []*struct {
		R, S      *big.Int
		publicKey *ecdsa.PublicKey
	},
) (*Transaction, error) {
	if len(signatures) != len(tb.transaction.TxIn) {
		return nil, fmt.Errorf("wrong signatures count")
	}

	for i, input := range tb.transaction.TxIn {
		signature := signatures[i]

		signatureBytes := append(
			(&btcec.Signature{R: signature.R, S: signature.S}).Serialize(),
			byte(txscript.SigHashAll),
		)
		publicKeyBytes := (*btcec.PublicKey)(
			signature.publicKey,
		).SerializeCompressed()

		sigHashArgs := tb.sigHashArgs[i]

		if sigHashArgs.witness {
			witness := wire.TxWitness{
				signatureBytes,
				publicKeyBytes,
			}

			// If the Witness field was pre-filled with data, put them at
			// the end of the final witness field. This is the case for
			// P2WSH inputs.
			if len(input.Witness) == 1 {
				witness = append(witness, input.Witness[0])
			}

			input.Witness = witness
		} else {
			builder := txscript.NewScriptBuilder().
				AddData(signatureBytes).
				AddData(publicKeyBytes)

			// If the SignatureScript field was pre-filled with data, put them
			// at the end of the final SignatureScript field. This is the case
			// for P2SH inputs.
			if len(input.SignatureScript) > 0 {
				builder.AddData(input.SignatureScript)
			}

			script, err := builder.Script()
			if err != nil {
				return nil, fmt.Errorf(
					"cannot build signature script for input [%v]: [%v]",
					i,
					err,
				)
			}

			input.SignatureScript = script
		}
	}

	return tb.stateToTransaction(), nil
}

// stateToTransaction converts the internal state of the builder into a Transaction.
func (tb *TransactionBuilder) stateToTransaction() *Transaction {
	inputs := make([]*TransactionInput, len(tb.transaction.TxIn))
	for i, input := range tb.transaction.TxIn {
		inputs[i] = &TransactionInput{
			Outpoint: &TransactionOutpoint{
				TransactionHash: Hash(input.PreviousOutPoint.Hash),
				OutputIndex:     input.PreviousOutPoint.Index,
			},
			SignatureScript: input.SignatureScript,
			Witness:         input.Witness,
			Sequence:        input.Sequence,
		}
	}

	outputs := make([]*TransactionOutput, len(tb.transaction.TxOut))
	for i, output := range tb.transaction.TxOut {
		outputs[i] = &TransactionOutput{
			Value:           output.Value,
			PublicKeyScript: output.PkScript,
		}
	}

	return &Transaction{
		Version:  tb.transaction.Version,
		Inputs:   inputs,
		Outputs:  outputs,
		Locktime: tb.transaction.LockTime,
	}
}

// TotalInputsValue returns the total value of transaction inputs.
func (tb *TransactionBuilder) TotalInputsValue() int64 {
	totalInputsValue := int64(0)

	for _, sigHashArgs := range tb.sigHashArgs {
		totalInputsValue += sigHashArgs.value
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
