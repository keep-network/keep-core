package bitcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

// Script represents an arbitrary Bitcoin script, NOT prepended with the
// byte-length of the script
type Script []byte

// NewScriptFromVarLenData construct a Script instance based on the provided
// variable length data prepended with a CompactSizeUint.
func NewScriptFromVarLenData(varLenData []byte) (Script, error) {
	// Extract the CompactSizeUint value that holds the byte length of the script.
	// Also, extract the byte length of the CompactSizeUint itself.
	scriptByteLength, compactByteLength, err := readCompactSizeUint(varLenData)
	if err != nil {
		return nil, fmt.Errorf("cannot read compact size uint: [%v]", err)
	}

	// Make sure the combined byte length of the script and the byte length
	// of the CompactSizeUint matches the total byte length of the variable
	// length data. Otherwise, the input data slice is malformed.
	if uint64(scriptByteLength)+uint64(compactByteLength) != uint64(len(varLenData)) {
		return nil, fmt.Errorf("malformed var len data")
	}

	// Extract the actual script by omitting the leading CompactSizeUint.
	return varLenData[compactByteLength:], nil
}

// ToVarLenData converts the Script to a byte array prepended with a
// CompactSizeUint holding the script's byte length.
func (s Script) ToVarLenData() ([]byte, error) {
	compactBytes, err := writeCompactSizeUint(CompactSizeUint(len(s)))
	if err != nil {
		return nil, fmt.Errorf("cannot write compact size uint: [%v]", err)
	}

	return append(compactBytes, s...), nil
}

// PublicKeyHash constructs the 20-byte public key hash by applying SHA-256
// then RIPEMD-160 on the provided ECDSA public key.
func PublicKeyHash(publicKey *ecdsa.PublicKey) [20]byte {
	publicKeyBytes := elliptic.MarshalCompressed(
		publicKey.Curve,
		publicKey.X,
		publicKey.Y,
	)

	publicKeyHashBytes := btcutil.Hash160(publicKeyBytes)

	var result [20]byte
	copy(result[:], publicKeyHashBytes)

	return result
}

// PayToWitnessPublicKeyHash constructs a P2WPKH script for the provided
// 20-byte public key hash. The function assumes the provided public key hash
// is valid.
func PayToWitnessPublicKeyHash(publicKeyHash [20]byte) (Script, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).
		AddData(publicKeyHash[:]).
		Script()
}

// PayToPublicKeyHash constructs a P2PKH script for the provided 20-byte public
// key hash. The function assumes the provided public key hash is valid.
func PayToPublicKeyHash(publicKeyHash [20]byte) (Script, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(publicKeyHash[:]).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		Script()
}

// WitnessScriptHash constructs the 32-byte witness script hash by applying
// single SHA-256 on the provided Script.
func WitnessScriptHash(script Script) [32]byte {
	return sha256.Sum256(script)
}

// ScriptHash constructs the 20-byte script hash by applying SHA-256 then
// RIPEMD-160 on the provided Script.
func ScriptHash(script Script) [20]byte {
	scriptHashBytes := btcutil.Hash160(script)

	var result [20]byte
	copy(result[:], scriptHashBytes)

	return result
}

// PayToWitnessScriptHash constructs a P2WSH script for the provided 32-byte
// witness script hash. The function assumes the provided script hash is valid.
func PayToWitnessScriptHash(witnessScriptHash [32]byte) (Script, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).
		AddData(witnessScriptHash[:]).
		Script()
}

// PayToScriptHash constructs a P2SH script for the provided 20-byte script
// hash. The function assumes the provided script hash is valid.
func PayToScriptHash(scriptHash [20]byte) (Script, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_HASH160).
		AddData(scriptHash[:]).
		AddOp(txscript.OP_EQUAL).
		Script()
}
