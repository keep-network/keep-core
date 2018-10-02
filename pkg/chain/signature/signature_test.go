package signature

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

// FIXME This needs to be moved to chain/ethereum/keydecrypt_test.go
func TestKeyFileDecryption(t *testing.T) {
	goodKeyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	badKeyFile := "./testdata/nonexistent-file.booyan"

	tests := map[string]struct {
		keyFile      string
		password     string
		errorMessage string
	}{
		"good password": {
			keyFile:      goodKeyFile,
			password:     "password",
			errorMessage: "",
		},
		"bad file": {
			keyFile:  badKeyFile,
			password: "",
			errorMessage: fmt.Sprintf(
				"unable to read KeyFile %v [open %v: no such file or directory]",
				badKeyFile,
				badKeyFile,
			),
		},
		"bad password": {
			keyFile:  goodKeyFile,
			password: "nanananana",
			errorMessage: fmt.Sprintf(
				"unable to decrypt %v [could not decrypt key with given passphrase]",
				goodKeyFile,
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			_, err := ethereum.DecryptKeyFile(test.keyFile, test.password)
			message := ""
			if err != nil {
				message = err.Error()
			}

			if message != test.errorMessage {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.errorMessage,
					err,
				)
			}
		})
	}
}

func TestSign(t *testing.T) {
	// Test for one address + message = good signature.
	// Test for different address + message = different good signature.
	// Don't test for failure but note that the only failures are deep in C-land.
}

func TestRecoverPublicKey(t *testing.T) {
	hexPubkey := "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e"

	tests := map[string]struct {
		message              string
		signature            string
		errorMessage         string
		expectedHexPublicKey *string
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
			if hexKey != *test.expectedHexPublicKey {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedHexPublicKey,
					hexKey,
				)
			}
		})
	}
}

func TestPublicKeyToAddress(t *testing.T) {
	// public key is nil
	// public key is bad somehow?
	// public key is good, address is correct
}
