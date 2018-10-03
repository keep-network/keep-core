package signature

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

func TestSign(t *testing.T) {
	keyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	keyPassword := "password"

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
		"correct signature 2": {
			message:              "Ethereum is a lot of fun.",
			errorMessage:         "",
			expectedMessageInHex: "457468657265756d2069732061206c6f74206f662066756e2e",
			expectedSignature:    "c37e29996a39a237f46a3eeea8be2707d37e455ef29ffca10188089b1f47bbab010d020a093d3c617d65e9c1bb6cb50c964accf2c215ea979e69c90d0d66eab400",
		},
		"correct signature 3": {
			message:              "Move $1,000,000 from Tim's account to my account right now!",
			errorMessage:         "",
			expectedMessageInHex: "4d6f76652024312c3030302c3030302066726f6d2054696d2773206163636f756e7420746f206d79206163636f756e74207269676874206e6f7721",
			expectedSignature:    "7220cfea62ee991206156c6301d49bc16d841bff387d81acb8a1cad5a84c801c37871f7f4a37a070e42aa3836512e483ba908924031079224facb5c2dfe2ecc401",
		},
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
	}

	key, err := ethereum.DecryptKeyFile(keyFile, keyPassword)
	if err != nil {
		t.Fatalf("Missing %s or invalid [%s]\n", keyFile, err)
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

			// Did we get the correct message back?
			if messageInHex != test.expectedMessageInHex {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedMessageInHex,
					messageInHex,
				)
			}

			// Is the signature correct
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
		"correct signature for pubkey 2": {
			message:              "457468657265756d2069732061206c6f74206f662066756e2e",
			signature:            "c37e29996a39a237f46a3eeea8be2707d37e455ef29ffca10188089b1f47bbab010d020a093d3c617d65e9c1bb6cb50c964accf2c215ea979e69c90d0d66eab400",
			errorMessage:         "",
			expectedHexPublicKey: &hexPubkey,
			signatureCorrect:     true,
		},
		"correct signature for pubkey 3": {
			message:              "4d6f76652024312c3030302c3030302066726f6d2054696d2773206163636f756e7420746f206d79206163636f756e74207269676874206e6f7721",
			signature:            "7220cfea62ee991206156c6301d49bc16d841bff387d81acb8a1cad5a84c801c37871f7f4a37a070e42aa3836512e483ba908924031079224facb5c2dfe2ecc401",
			errorMessage:         "",
			expectedHexPublicKey: &hexPubkey,
			signatureCorrect:     true,
		},
		"in-correct signature for pubkey": {
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
		"valid address 2": {
			expectedAddressHex: "0xDb180Da9A8982C7Bb75Ca40039f959CB959c62e8",
			keyFile:            "./testdata/UTC--2018-08-27T00-03-51.111292084Z--Db180Da9A8982C7Bb75Ca40039f959CB959c62e8",
			keyPassword:        "vEbeMJ/kP9mN2gdI",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			key, err := ethereum.DecryptKeyFile(test.keyFile, test.keyPassword)
			if err != nil {
				t.Fatalf("Unable to read keyfile [%s] error [%v]\n", test.keyFile, err)
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
	keyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	keyPassword := "password"

	tests := map[string]struct {
		message       string
		errorMessage  string
		expectedValid bool
		signature     string
	}{
		"verify correct signature": {
			message:       "3031303230333034",
			signature:     "2844b7b1b57a020623c70c842c5795dce6bc61531dac75b5246c5825c44644b44fc0160fc82ccfdac1463407e7a2ff474beaf30d41674a9ee72838d39b0e5fec01",
			errorMessage:  "",
			expectedValid: true,
		},
		"verify correct signature 2": {
			message:       "457468657265756d2069732061206c6f74206f662066756e2e",
			signature:     "c37e29996a39a237f46a3eeea8be2707d37e455ef29ffca10188089b1f47bbab010d020a093d3c617d65e9c1bb6cb50c964accf2c215ea979e69c90d0d66eab400",
			errorMessage:  "",
			expectedValid: true,
		},
		"verify correct signature 3": {
			message:       "4d6f76652024312c3030302c3030302066726f6d2054696d2773206163636f756e7420746f206d79206163636f756e74207269676874206e6f7721",
			signature:     "7220cfea62ee991206156c6301d49bc16d841bff387d81acb8a1cad5a84c801c37871f7f4a37a070e42aa3836512e483ba908924031079224facb5c2dfe2ecc401",
			errorMessage:  "",
			expectedValid: true,
		},
	}

	key, err := ethereum.DecryptKeyFile(keyFile, keyPassword)
	if err != nil {
		t.Fatalf("Missing %s or invalid [%s]\n", keyFile, err)
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			// func VerifySignatureWithPubKey(pubkey *ecdsa.PublicKey, sig, msg string) (bool, error) {
			ok, err := VerifySignatureWithPubKey(&key.PrivateKey.PublicKey, test.signature, test.message)

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
