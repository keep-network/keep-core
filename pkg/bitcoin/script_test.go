package bitcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"

	"github.com/keep-network/keep-core/internal/testutils"
)

func TestNewScriptFromVarLenData(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	var tests = map[string]struct {
		data           []byte
		expectedScript Script
		expectedErr    error
	}{
		"proper variable length data": {
			data:           fromHex("1600148db50eb52063ea9d98b3eac91489a90f738986f6"),
			expectedScript: fromHex("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
		},
		"nil variable length data": {
			data:        nil,
			expectedErr: fmt.Errorf("cannot read compact size uint: [EOF]"),
		},
		"variable length data with missing script": {
			data:        fromHex("16"),
			expectedErr: fmt.Errorf("malformed var len data"),
		},
		"variable length data with missing compact size uint": {
			data:        fromHex("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
			expectedErr: fmt.Errorf("malformed var len data"),
		},
		"variable length data with wrong compact size uint value": {
			data:        fromHex("1500148db50eb52063ea9d98b3eac91489a90f738986f6"),
			expectedErr: fmt.Errorf("malformed var len data"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			script, err := NewScriptFromVarLenData(test.data)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\nexpected: %+v\nactual:   %+v\n",
					test.expectedErr,
					err,
				)
			}

			testutils.AssertBytesEqual(t, test.expectedScript, script)
		})
	}
}

func TestScript_ToVarLenData(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	bytes20HashBytes, err := hex.DecodeString(
		"8db50eb52063ea9d98b3eac91489a90f738986f6",
	)
	if err != nil {
		t.Fatal(err)
	}
	var bytes20Hash [20]byte
	copy(bytes20Hash[:], bytes20HashBytes)

	bytes32HashBytes, err := hex.DecodeString(
		"86a303cdd2e2eab1d1679f1a813835dc5a1b65321077cdccaf08f98cbf04ca96",
	)
	if err != nil {
		t.Fatal(err)
	}
	var bytes32Hash [32]byte
	copy(bytes32Hash[:], bytes32HashBytes)

	var tests = map[string]struct {
		script        Script
		expectedBytes []byte
	}{
		"p2pkh script": {
			script: func() Script {
				result, err := PayToPublicKeyHash(bytes20Hash)
				if err != nil {
					t.Fatal(err)
				}
				return result
			}(),
			expectedBytes: fromHex("1976a9148db50eb52063ea9d98b3eac91489a90f738986f688ac"),
		},
		"p2wpkh script": {
			script: func() Script {
				result, err := PayToWitnessPublicKeyHash(bytes20Hash)
				if err != nil {
					t.Fatal(err)
				}
				return result
			}(),
			expectedBytes: fromHex("1600148db50eb52063ea9d98b3eac91489a90f738986f6"),
		},
		"p2sh script": {
			script: func() Script {
				result, err := PayToScriptHash(bytes20Hash)
				if err != nil {
					t.Fatal(err)
				}
				return result
			}(),
			expectedBytes: fromHex("17a9148db50eb52063ea9d98b3eac91489a90f738986f687"),
		},
		"p2wsh script": {
			script: func() Script {
				result, err := PayToWitnessScriptHash(bytes32Hash)
				if err != nil {
					t.Fatal(err)
				}
				return result
			}(),
			expectedBytes: fromHex("22002086a303cdd2e2eab1d1679f1a813835dc5a1b65321077cdccaf08f98cbf04ca96"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			bytes, err := test.script.ToVarLenData()
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBytesEqual(t, test.expectedBytes, bytes)
		})
	}
}

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
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	expectedResult, err := hex.DecodeString(
		"a9143ec459d0f3c29286ae5df5fcc421e2786024277e87",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(t, expectedResult, result[:])
}
