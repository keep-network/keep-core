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
		"dkg result has only group public key provided": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{100},
			},
			expectedHash: "62f670bf6f172ab82df59082f8255ccae11e0fd956be902f5601a5c3a12ba1a5",
		},
		"dkg result has only disqualified provided": {
			dkgResult: &relaychain.DKGResult{
				Disqualified: []byte{1, 0, 1, 0},
			},
			expectedHash: "22c8e49873c2173ae650f7a241d2808d068ae0d3a5121bac41e8597bd70459f4",
		},
		"dkg result has only inactive provided": {
			dkgResult: &relaychain.DKGResult{
				Inactive: []byte{0, 1, 1, 1},
			},
			expectedHash: "13cae25f320b3b54ba1b03faba0bb38e793b7289109e2ac00c30be39d40487a2",
		},
		"dkg result has all parameters provided": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{3, 40, 200},
				Disqualified:   []byte{1, 0, 1, 0},
				Inactive:       []byte{0, 1, 1, 0},
			},
			expectedHash: "4e8b56086bfc0ceb8c59a546c2fd38e5becb77c0a38bd74f21c57c6499603180",
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
			expectedHash: "d48bbfd2b4b22423d354a919f0f9b993a5e3fbd0c93cb6a68ec2e87709349900",
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
			expectedHash: "a79c258065c5e01c83afd0b581b47623d7e020e1f8288cb5c26d337fb5537adf",
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
