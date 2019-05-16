package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
)

func TestSign(t *testing.T) {
	keyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	keyPassword := "password"

	// Notes:
	//
	// Any error failures come from somewhare in C land.  This makes an
	// error case very hard to test.  (Also errors from C depend on
	// go-ethereum's compile time flags as different C implementations
	// are supported - so the errors would not be consistent.  The
	// errors would depend on the set of OS libraries installed and
	// the underlying architecture of the system.)
	//
	// The data for #2 and #3 signatures were generated from NaCl - this
	// provides independent verification that the signature is correct.
	tests := map[string]struct {
		message              string
		errorMessage         string
		expectedMessageInHex string
		expectedSignature    string
	}{
		"correct signature": {
			message:              "01020304",
			errorMessage:         "",
			expectedMessageInHex: "3031303230333034",
			expectedSignature:    "2844b7b1b57a020623c70c842c5795dce6bc61531dac75b5246c5825c44644b44fc0160fc82ccfdac1463407e7a2ff474beaf30d41674a9ee72838d39b0e5fec01",
		},
		"correct NaCl signature": {
			message:              "Ethereum is a lot of fun.",
			errorMessage:         "",
			expectedMessageInHex: "457468657265756d2069732061206c6f74206f662066756e2e",
			expectedSignature:    "c37e29996a39a237f46a3eeea8be2707d37e455ef29ffca10188089b1f47bbab010d020a093d3c617d65e9c1bb6cb50c964accf2c215ea979e69c90d0d66eab400",
		},
	}

	key, err := ethutil.DecryptKeyFile(keyFile, keyPassword)
	if err != nil {
		t.Fatalf("Failed to read key file [%s] [%v]\n", keyFile, err)
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			messageInHex, signature, err := Sign(key, []byte(test.message))
			if err != nil {
				if test.errorMessage == "" || err.Error() != test.errorMessage {
					t.Errorf(
						"\nexpected: [%v]\nactual:   [%v]",
						test.errorMessage,
						err,
					)
				}
				return
			}

			if messageInHex != test.expectedMessageInHex {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedMessageInHex,
					messageInHex,
				)
			}

			if signature != test.expectedSignature {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedSignature,
					signature,
				)
			}

		})
	}
}

func TestRecoverPublicKey(t *testing.T) {
	hexPubkey := "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e"

	tests := map[string]struct {
		message              string
		signature            string
		errorMessage         string
		expectedHexPublicKey *string
		signatureCorrect     bool
	}{
		"bad message": {
			message:              "this is a test",
			signature:            "should be hex",
			errorMessage:         "failed to decode hex message to bytes: [encoding/hex: invalid byte: U+0074 't']",
			expectedHexPublicKey: nil,
		},
		"bad signature": {
			message:              "01020304",
			signature:            "should be hex",
			errorMessage:         "failed to decode hex signature to bytes: [encoding/hex: invalid byte: U+0073 's']",
			expectedHexPublicKey: nil,
		},
		"correct signature for pubkey": {
			message:              "01020304",
			signature:            "b6d61b98d0722a249c9cad3e16de3626d4969cef56ab12e9efb3ef00a4f9356e5a25574aed6447d3d2797ee8afb71b8b7ff68c0f2cfe8fa437d145f16a192fb201",
			errorMessage:         "",
			expectedHexPublicKey: &hexPubkey,
			signatureCorrect:     true,
		},
		"correct NaCL verified signature for pubkey": {
			message:              "457468657265756d2069732061206c6f74206f662066756e2e",
			signature:            "c37e29996a39a237f46a3eeea8be2707d37e455ef29ffca10188089b1f47bbab010d020a093d3c617d65e9c1bb6cb50c964accf2c215ea979e69c90d0d66eab400",
			errorMessage:         "",
			expectedHexPublicKey: &hexPubkey,
			signatureCorrect:     true,
		},
		"incorrect signature for pubkey": {
			message:              "4d6f76652024312c3030302c3030302066726f6d2054696d2773206163636f756e7420746f206d79206163636f756e74207269676874206e6f7721",
			signature:            "b6d61b98d0722a249c9cad3e16de3626d4969cef56ab12e9efb3ef00a4f9356e5a25574aed6447d3d2797ee8afb71b8b7ff68c0f2cfe8fa437d145f16a192fb201",
			errorMessage:         "",
			expectedHexPublicKey: &hexPubkey,
			signatureCorrect:     false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publicKey, err := RecoverPublicKey(test.signature, test.message)

			if err != nil {
				if test.errorMessage == "" || err.Error() != test.errorMessage {
					t.Errorf(
						"\nexpected: [%v]\nactual:   [%v]",
						test.errorMessage,
						err,
					)
				}
				return
			}

			hexKey := hex.EncodeToString(crypto.FromECDSAPub(publicKey))
			if test.signatureCorrect {
				if hexKey != *test.expectedHexPublicKey {
					t.Errorf(
						"\nexpected: [%v]\nactual:   [%v]",
						*test.expectedHexPublicKey,
						hexKey,
					)
				}
			} else {
				if hexKey == *test.expectedHexPublicKey {
					t.Errorf(
						"\nexpected: [%v]\nactual:   [%v]",
						*test.expectedHexPublicKey,
						hexKey,
					)
				}
			}
		})
	}
}

func TestPublicKeyToAddress(t *testing.T) {
	tests := map[string]struct {
		expectedAddressHex string
		keyFile            string
		keyPassword        string
	}{
		"valid address": {
			expectedAddressHex: "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			keyFile:            "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			keyPassword:        "password",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			key, err := ethutil.DecryptKeyFile(test.keyFile, test.keyPassword)
			if err != nil {
				t.Fatalf("Failed to read key file [%s] [%v]\n", test.keyFile, err)
				return
			}
			address := PublicKeyToAddress(&key.PrivateKey.PublicKey)
			addressHex := address.Hex()

			// Did we get the correct message back?
			if addressHex != test.expectedAddressHex {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedAddressHex,
					addressHex,
				)
			}

		})
	}
}

func TestVerifySignatureWithPubKey(t *testing.T) {

	tests := map[string]struct {
		message       string
		errorMessage  string
		expectedValid bool
		pubKeyJSON    string
		signature     string
	}{
		"verify correct signature": {
			message:       "3031303230333034",
			signature:     "2844b7b1b57a020623c70c842c5795dce6bc61531dac75b5246c5825c44644b44fc0160fc82ccfdac1463407e7a2ff474beaf30d41674a9ee72838d39b0e5fec01",
			pubKeyJSON:    `{"Curve":{"P":115792089237316195423570985008687907853269984665640564039457584007908834671663,"N":115792089237316195423570985008687907852837564279074904382605163141518161494337,"B":7,"Gx":55066263022277343669578718895168534326250603453777594175500187360389116729240,"Gy":32670510020758816978083085130507043184471273380659243275938904335757337482424,"BitSize":256},"X":59583986354391045197850449728092779822530178549585112526757356439863210932224,"Y":84561927522964405625378711571517351344267579854412957918810298793357853572174}`,
			errorMessage:  "",
			expectedValid: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			var pubkey ecdsa.PublicKey
			json.Unmarshal([]byte(test.pubKeyJSON), &pubkey)
			ok, err := VerifySignatureWithPubKey(&pubkey, test.signature, test.message)
			if err != nil {
				if test.errorMessage == "" || err.Error() != test.errorMessage {
					t.Errorf(
						"\nexpected: [%v]\nactual:   [%v]",
						test.errorMessage,
						err,
					)
				}
				return
			}

			if ok != test.expectedValid {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedValid,
					ok,
				)
			}

		})
	}
}
