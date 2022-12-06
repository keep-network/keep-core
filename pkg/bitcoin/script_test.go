package bitcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"testing"
)

func TestPublicKeyHash(t *testing.T) {
	// An arbitrary uncompressed public key.
	publicKeyBytes, err := hex.DecodeString(
		"04989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9" +
			"d218b65e7d91c752f7b22eaceb771a9af3a6f3d3f010a5d471a1aeef7d7713af",
	)
	if err != nil {
		t.Fatal(err)
	}

	x, y := elliptic.Unmarshal(btcec.S256(), publicKeyBytes)
	publicKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}

	result := PublicKeyHash(publicKey)

	expectedResult, err := hex.DecodeString(
		"8db50eb52063ea9d98b3eac91489a90f738986f6",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestPayToWitnessPublicKeyHash(t *testing.T) {
	// The 20-byte public key hash, same as the output of TestPublicKeyHash.
	publicKeyHashBytes, err := hex.DecodeString(
		"8db50eb52063ea9d98b3eac91489a90f738986f6",
	)
	if err != nil {
		t.Fatal(err)
	}

	var publicKeyHash [20]byte
	copy(publicKeyHash[:], publicKeyHashBytes)

	result, err := PayToWitnessPublicKeyHash(publicKeyHash)

	expectedResult, err := hex.DecodeString(
		"00148db50eb52063ea9d98b3eac91489a90f738986f6",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestPayToPublicKeyHash(t *testing.T) {
	// The 20-byte public key hash, same as the output of TestPublicKeyHash.
	publicKeyHashBytes, err := hex.DecodeString(
		"8db50eb52063ea9d98b3eac91489a90f738986f6",
	)
	if err != nil {
		t.Fatal(err)
	}

	var publicKeyHash [20]byte
	copy(publicKeyHash[:], publicKeyHashBytes)

	result, err := PayToPublicKeyHash(publicKeyHash)

	expectedResult, err := hex.DecodeString(
		"76a9148db50eb52063ea9d98b3eac91489a90f738986f688ac",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestWitnessScriptHash(t *testing.T) {
	// An arbitrary redeem script.
	script, err := hex.DecodeString(
		"14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d0003952375" +
			"76a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e25" +
			"7eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac68",
	)
	if err != nil {
		t.Fatal(err)
	}

	result := WitnessScriptHash(script)

	expectedResult, err := hex.DecodeString(
		"86a303cdd2e2eab1d1679f1a813835dc5a1b65321077cdccaf08f98cbf04ca96",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestScriptHash(t *testing.T) {
	// An arbitrary redeem script.
	script, err := hex.DecodeString(
		"14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d0003952375" +
			"76a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e25" +
			"7eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac68",
	)
	if err != nil {
		t.Fatal(err)
	}

	result := ScriptHash(script)

	expectedResult, err := hex.DecodeString(
		"3ec459d0f3c29286ae5df5fcc421e2786024277e",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestPayToWitnessScriptHash(t *testing.T) {
	// The 32-byte script hash, same as the output of TestWitnessScriptHash.
	scriptHashBytes, err := hex.DecodeString(
		"86a303cdd2e2eab1d1679f1a813835dc5a1b65321077cdccaf08f98cbf04ca96",
	)
	if err != nil {
		t.Fatal(err)
	}

	var scriptHash [32]byte
	copy(scriptHash[:], scriptHashBytes)

	result, err := PayToWitnessScriptHash(scriptHash)

	expectedResult, err := hex.DecodeString(
		"002086a303cdd2e2eab1d1679f1a813835dc5a1b65321077cdccaf08f98cbf04ca96",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}

func TestPayToScriptHash(t *testing.T) {
	// The 20-byte script hash, same as the output of TestScriptHash.
	scriptHashBytes, err := hex.DecodeString(
		"3ec459d0f3c29286ae5df5fcc421e2786024277e",
	)
	if err != nil {
		t.Fatal(err)
	}

	var scriptHash [20]byte
	copy(scriptHash[:], scriptHashBytes)

	result, err := PayToScriptHash(scriptHash)

	expectedResult, err := hex.DecodeString(
		"a9143ec459d0f3c29286ae5df5fcc421e2786024277e87",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}
