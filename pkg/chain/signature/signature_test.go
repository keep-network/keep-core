package signature

import (
	"testing"

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
		"of successful signature without errors": {
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
		"of invalid password on decrypiton of a signature file": {
			in:                []byte{01, 02, 03, 04},
			addr:              "6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:       "01020304",
			expectedEIP55Addr: "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:    "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:           "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:          "nanananana",
			expectError:       true,
		},
		"with a valid KeyFile password, but an invalid signature": {
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
					t.Errorf("failed to returne an error [%v] \n", err)
				}
				return
			}
			msg, sig, err := Sign(key, test.in)
			if test.expectError {
				if err == nil {
					t.Errorf("failed to returne an error [%v] \n", err)
				}
				return
			}
			if err != nil {
				t.Errorf("returned an error [%v] \n", err)
			}
			if msg != test.expectedMsg {
				t.Errorf("expected %s got %s\n", test.expectedMsg, msg)
			}
			val, err := VerifySignature(test.addr, sig, msg)
			if test.expectNotToValidate {
				if val.IsValid {
					t.Errorf("should not have validated but did\n")
				}
				return
			}
			if err != nil {
				t.Errorf("falied to verify [%v] \n", err)
			}
			if val.RecoveredAddress != test.expectedEIP55Addr {
				t.Errorf("expected %s got %s\n", test.expectedEIP55Addr, val.RecoveredAddress)
			}
			if val.RecoveredPublicKey != test.expectedPubKey {
				t.Errorf("expected %s got %s\n", test.expectedEIP55Addr, val.RecoveredPublicKey)
			}

			if got := MessageHasValidSignature(test.addr, sig, msg); got == test.expectNotToValidate {
				t.Errorf("validation error, expected %v, got %v\n", !test.expectNotToValidate, got)
			}

			pk, err := GetPublicKey(test.addr, sig, msg)
			if err == nil {
				if pk != test.expectedPubKey {
					t.Errorf("invalid public key, expected %v, got %v\n", test.expectedPubKey, pk)
				}
			}

			pkEcdsa, _ := GetPublicKeyECDSA(sig, msg)
			ca := PublicKeyToAddress(pkEcdsa)
			encCa := EncodeAddressToEIP55(ca)
			if encCa != test.expectedEIP55Addr {
				t.Errorf("error, expected %s, got %s\n", test.expectedEIP55Addr, encCa)
			}

		})
	}

}
