package operator

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"math/big"
	"reflect"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	curveImplementation := btcec.S256()

	privateKey, publicKey, err := GenerateKeyPair(curveImplementation)
	if err != nil {
		t.Fatal(err)
	}

	expectedCurve := Secp256k1
	actualCurve := publicKey.Curve
	if expectedCurve != actualCurve {
		t.Errorf(
			"unexpected curve name\nexpected: %v\nactual:   %v",
			expectedCurve,
			actualCurve,
		)
	}

	if !reflect.DeepEqual(privateKey.PublicKey, *publicKey) {
		t.Errorf("private key contains wrong public key")
	}

	if !curveImplementation.IsOnCurve(publicKey.X, publicKey.Y) {
		t.Errorf("public key coordinates are not on the curve")
	}

	expectedX, expectedY := curveImplementation.ScalarBaseMult(privateKey.D.Bytes())
	if expectedX.Cmp(publicKey.X) != 0 {
		t.Errorf(
			"unexpected public key X coordinate\nexpected: %v\nactual:   %v",
			expectedX,
			publicKey.X,
		)
	}
	if expectedY.Cmp(publicKey.Y) != 0 {
		t.Errorf(
			"unexpected public key Y coordinate\nexpected: %v\nactual:   %v",
			expectedY,
			publicKey.Y,
		)
	}
}

func TestGenerateKeyPair_UnsupportedCurve(t *testing.T) {
	curveImplementation := elliptic.P256()

	_, _, err := GenerateKeyPair(curveImplementation)

	expectedError := fmt.Errorf(
		"cannot parse curve name: [unknown curve: [%v]]",
		curveImplementation.Params().Name,
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestMarshalUncompressed(t *testing.T) {
	x, ok := new(big.Int).SetString(
		"6f3d3cc22b5a3ab0cd9f56500b0abd104476a9d4b7a55fef000fee30ba4a7768",
		16,
	)
	if !ok {
		t.Fatal("cannot set X coordinate")
	}

	y, ok := new(big.Int).SetString(
		"d9d487bb049778b75614b5971a05a70d4bdf19d675d16c41e8b0da718d637c80",
		16,
	)
	if !ok {
		t.Fatal("cannot set Y coordinate")
	}

	publicKey := &PublicKey{
		Curve: Secp256k1,
		X:     x,
		Y:     y,
	}

	publicKeyBytes := MarshalUncompressed(publicKey)

	expectedPublicKeyHex :=
		"046f3d3cc22b5a3ab0cd9f56500b0abd104476a9d4b7a55fef000fee30ba4a7768" +
			"d9d487bb049778b75614b5971a05a70d4bdf19d675d16c41e8b0da718d637c80"
	actualPublicKeyHex := hex.EncodeToString(publicKeyBytes)
	if expectedPublicKeyHex != actualPublicKeyHex {
		t.Errorf(
			"unexpected marshaled public key\nexpected: %v\nactual:   %v\n",
			expectedPublicKeyHex,
			actualPublicKeyHex,
		)
	}
}

func TestMarshalCompressed_EvenYCoordinate(t *testing.T) {
	x, ok := new(big.Int).SetString(
		"6f3d3cc22b5a3ab0cd9f56500b0abd104476a9d4b7a55fef000fee30ba4a7768",
		16,
	)
	if !ok {
		t.Fatal("cannot set X coordinate")
	}

	y, ok := new(big.Int).SetString(
		"d9d487bb049778b75614b5971a05a70d4bdf19d675d16c41e8b0da718d637c80",
		16,
	)
	if !ok {
		t.Fatal("cannot set Y coordinate")
	}

	publicKey := &PublicKey{
		Curve: Secp256k1,
		X:     x,
		Y:     y,
	}

	expectedPublicKeyHex :=
		"026f3d3cc22b5a3ab0cd9f56500b0abd104476a9d4b7a55fef000fee30ba4a7768"
	actualPublicKeyHex := publicKey.String()
	if expectedPublicKeyHex != actualPublicKeyHex {
		t.Errorf(
			"unexpected marshaled public key\nexpected: %v\nactual:   %v\n",
			expectedPublicKeyHex,
			actualPublicKeyHex,
		)
	}
}

func TestMarshalCompressed_OddYCoordinate(t *testing.T) {
	x, ok := new(big.Int).SetString(
		"063b948fffa3220bf302ee2874cfd26802f2fab2a3e6ababcb66ffb6d21218b3",
		16,
	)
	if !ok {
		t.Fatal("cannot set X coordinate")
	}

	y, ok := new(big.Int).SetString(
		"9a06178ec1ec41813b54bcbc7c4cd8061983f87948a4452e07aafd91992369ed",
		16,
	)
	if !ok {
		t.Fatal("cannot set Y coordinate")
	}

	publicKey := &PublicKey{
		Curve: Secp256k1,
		X:     x,
		Y:     y,
	}

	expectedPublicKeyHex :=
		"03063b948fffa3220bf302ee2874cfd26802f2fab2a3e6ababcb66ffb6d21218b3"
	actualPublicKeyHex := publicKey.String()
	if expectedPublicKeyHex != actualPublicKeyHex {
		t.Errorf(
			"unexpected marshaled public key\nexpected: %v\nactual:   %v\n",
			expectedPublicKeyHex,
			actualPublicKeyHex,
		)
	}
}
