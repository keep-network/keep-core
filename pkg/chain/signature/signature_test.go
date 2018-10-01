package signature

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

func TestSignature(t *testing.T) {

	tests := map[string]struct {
		in                  []byte
		addr                string
		expectedMsg         string
		expectedEIP55Addr   string
		expectedPubKey      string
		keyFile             string
		password            string
		expectError         bool
		expectNotToValidate bool
	}{
		"successful signature": {
			in:                  []byte{01, 02, 03, 04},
			addr:                "6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:         "01020304",
			expectedEIP55Addr:   "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:      "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:             "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:            "password",
			expectError:         false,
			expectNotToValidate: false,
		},
		"signature file with invalid password": {
			in:                []byte{01, 02, 03, 04},
			addr:              "6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:       "01020304",
			expectedEIP55Addr: "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:    "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:           "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:          "nanananana",
			expectError:       true,
		},
		"invalid signature": {
			in:                  []byte{01, 02, 03, 04},
			addr:                "9ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:         "01020304",
			expectedEIP55Addr:   "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:      "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:             "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:            "password",
			expectError:         false,
			expectNotToValidate: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			key, err := ethereum.DecryptKeyFile(test.keyFile, test.password)
			if test.expectError {
				if err == nil {
					t.Errorf("failed to return an error when error expected. Should have [%v]\n",
						test.expectError)
				}
				return
			}
			msg, sig, err := Sign(key, test.in)
			if test.expectError {
				if err == nil {
					t.Errorf("failed to return an error when error expected. Should have [%v]\n",
						test.expectError)
				}
				return
			}
			if err != nil {
				t.Errorf("returned an error [%v] \n", err)
			}
			if msg != test.expectedMsg {
				t.Errorf("Message invalid\nexpected: [%s]\nactual : [%s]\n", test.expectedMsg, msg)
			}

			// From sig/msg to a public key.
			recoveredPubkey, err := RecoverPublicKey(sig, msg)
			RecoveredPublicKey := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
			if RecoveredPublicKey != test.expectedPubKey {
				t.Errorf("invalid recovered public key\nexpected: [%s]\nactual : [%s]\n",
					test.expectedPubKey, RecoveredPublicKey)
			}

			// corrupt the msg/signature so it will not validate.
			if test.expectNotToValidate {
				msg = msg + "00"
			}

			// Determine if the original key is signer for sig/msg
			isValid, err := VerifySignatureWithPubKey(&key.PrivateKey.PublicKey, sig, msg)

			if test.expectNotToValidate && isValid {
				t.Errorf("should not have validated but did\n")
			} else if !test.expectNotToValidate && !isValid {
				t.Errorf("should have validated, but did not\n")
			}

			// Test that we get back the correct address for the recovered public key.
			addr := PublicKeyToAddress(recoveredPubkey)
			RecoveredAddress := addr.Hex()
			if RecoveredAddress != test.expectedEIP55Addr {
				t.Errorf("Invalid recovered address\nexpected: [%s] actual : [%s]\n",
					test.expectedEIP55Addr, RecoveredAddress)
			}

		})
	}

}
