package ethereum

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// TestCalculateDKGResultHash validates if calculated DKG result hash matches
// expected one.
//
// Expected hashes has been calculated on-chain with:
// `keccak256(abi.encode(success, groupPubKey, disqualified, inactive))`.
func TestCalculateDKGResultHash(t *testing.T) {
	chain := &ethereumChain{}

	var tests = map[string]struct {
		dkgResult    *relaychain.DKGResult
		expectedHash string
	}{
		"dkg result has only success provided": {
			dkgResult: &relaychain.DKGResult{
				Success: true,
			},
			expectedHash: "3550dcd2ae8c05d5ca00d5fcd004b6aaf7d3b90f64e0e6078ae62f1f7a2e3f27",
		},
		"dkg result has only group public key provided": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{100},
			},
			expectedHash: "132d7294d5d6a5434295d5f5d611bd210261f4e1360cfeb1054e30e69747ea16",
		},
		"dkg result has only disqualified provided": {
			dkgResult: &relaychain.DKGResult{
				Disqualified: []byte{1, 0, 1, 0},
			},
			expectedHash: "19120029436cfa37c800dce5ca006c1392b7c1d3d3dce50b91e44d1dc2e41c82",
		},
		"dkg result has only inactive provided": {
			dkgResult: &relaychain.DKGResult{
				Inactive: []byte{0, 1, 1, 1},
			},
			expectedHash: "941ef3e2fb3c006ea88533e1f3adf7ce974d114d30161e6a6209428a39747009",
		},
		"dkg result has all parameters provided": {
			dkgResult: &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: []byte{3, 40, 200},
				Disqualified:   []byte{1, 0, 1, 0},
				Inactive:       []byte{0, 1, 1, 0},
			},
			expectedHash: "d89ec37a4636214269a8aa8639691dc3c39157111a7321dddfcc6f6146b77fd3",
		},
		"dkg result has disqualified longer than 32 bytes": {
			dkgResult: &relaychain.DKGResult{
				Disqualified: []byte{
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				},
			},
			expectedHash: "b390bfe240f6cf522b4d866b069f08a96bf9d2c6a7cac6574291a8353e05562a",
		},
		"dkg result has group public key longer than 64 bytes": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{
					33, 249, 72, 108, 111, 44, 64, 58, 107, 112,
					108, 74, 214, 170, 149, 99, 212, 2, 48, 137,
					146, 12, 128, 8, 103, 47, 13, 161, 14, 126,
					5, 151, 0, 199, 90, 57, 31, 29, 175, 197,
					158, 45, 138, 205, 82, 95, 171, 104, 246, 8,
					203, 130, 138, 115, 72, 108, 232, 87, 129, 161,
					39, 228, 55, 222, 94, 238, 85, 128, 137, 187,
					27, 252, 25, 38, 201, 41, 127, 179, 75, 112,
				},
			},
			expectedHash: "86706e29471a9a32ca95525a8f50ec9dbd0751e6eff09f0bf3630023ea464545",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// expectedEncodedResult := common.Hex2Bytes(test.expectedSerializedResult)
			expectedHash := common.Hex2Bytes(test.expectedHash)

			actualHash, err := chain.CalculateDKGResultHash(test.dkgResult)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expectedHash, actualHash[:]) {
				t.Errorf(
					"\nexpected: %v\nactual:   %x\n",
					test.expectedHash,
					actualHash,
				)
			}
		})
	}
}
