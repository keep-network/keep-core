package ethereum

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
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
			expectedHash: "f1918e8562236eb17adc8502332f4c9c82bc14e19bfc0aa10ab674ff75b3d2f3",
		},
		"dkg result has only disqualified provided": {
			dkgResult: &relaychain.DKGResult{
				Disqualified: []byte{1, 0, 1, 0},
			},
			expectedHash: "ddb76fe48db5426b0729b9b973ecf962f375fd5f453f404e26f8b4bec4c97760",
		},
		"dkg result has only inactive provided": {
			dkgResult: &relaychain.DKGResult{
				Inactive: []byte{0, 1, 1, 1},
			},
			expectedHash: "967ecff5d1dbe8e817013123a8fb1762edfcda5ad776d5a22a9ff1dbb274cf0e",
		},
		"dkg result has all parameters provided": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{3, 40, 200},
				Disqualified:   []byte{1, 0, 1, 0},
				Inactive:       []byte{0, 1, 1, 0},
			},
			expectedHash: "a57664d91d1bbc7920ed3e658aea2060cf1105748f70eaabc434cfdfaf9d769a",
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
			expectedHash: "22df70f7a92e04a2ab0e7266b836f34dcf38fe96d8a9d4ce67f14bfbcc1a926b",
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
			expectedHash: "2a2f24baccd2f9bc371e08fff5bba43c7d4301d7d93fc97f2f93eb3cda28ddb6",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
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

func TestConvertSignaturesToChainFormat(t *testing.T) {
	memberIndex1 := group.MemberIndex(1)
	memberIndex2 := group.MemberIndex(2)
	memberIndex3 := group.MemberIndex(3)
	memberIndex4 := group.MemberIndex(4)
	memberIndex5 := group.MemberIndex(5)

	signature1 := common.LeftPadBytes([]byte("marry"), 65)
	signature2 := common.LeftPadBytes([]byte("had"), 65)
	signature3 := common.LeftPadBytes([]byte("a"), 65)
	signature4 := common.LeftPadBytes([]byte("little"), 65)
	signature5 := common.LeftPadBytes([]byte("lamb"), 65)

	invalidSignature := common.LeftPadBytes([]byte("invalid"), 64)

	var tests = map[string]struct {
		signaturesMap map[group.MemberIndex][]byte
		expectedError error
	}{
		"one valid signature": {
			signaturesMap: map[group.MemberIndex][]byte{
				memberIndex1: signature1,
			},
		},
		"five valid signatures": {
			signaturesMap: map[group.MemberIndex][]byte{
				memberIndex3: signature3,
				memberIndex1: signature1,
				memberIndex4: signature4,
				memberIndex5: signature5,
				memberIndex2: signature2,
			},
		},
		"invalid signature": {
			signaturesMap: map[group.MemberIndex][]byte{
				memberIndex1: signature1,
				memberIndex2: invalidSignature,
			},
			expectedError: fmt.Errorf("invalid signature size for member [2] got [64]-bytes but required [65]-bytes"),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			indicesSlice, signaturesSlice, err :=
				convertSignaturesToChainFormat(test.signaturesMap)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf(
					"invalid error\nexpected: %v\nactual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if test.expectedError == nil {
				if len(indicesSlice) != len(test.signaturesMap) {
					t.Errorf(
						"invalid member indices slice length\nexpected: %v\nactual:   %v\n",
						len(test.signaturesMap),
						len(indicesSlice),
					)
				}

				if len(signaturesSlice) != (SignatureSize * len(indicesSlice)) {
					t.Errorf(
						"invalid signatures slice size\nexpected: %v\nactual:   %v\n",
						(SignatureSize * len(indicesSlice)),
						len(signaturesSlice),
					)
				}
			}

			for i, actualMemberIndex := range indicesSlice {
				memberIndex := group.MemberIndex(actualMemberIndex.Uint64())

				actualSignature := signaturesSlice[SignatureSize*i : SignatureSize*(i+1)]
				if !bytes.Equal(
					test.signaturesMap[memberIndex],
					actualSignature,
				) {
					t.Errorf(
						"invalid signatures for member %v\nexpected: %v\nactual:   %v\n",
						actualMemberIndex,
						test.signaturesMap[memberIndex],
						actualSignature,
					)
				}
			}
		})
	}
}
