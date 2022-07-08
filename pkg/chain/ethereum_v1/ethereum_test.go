package ethereum_v1

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/beacon/chain"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
)

// TestCalculateDKGResultHash validates if calculated DKG result hash matches
// expected one.
//
// Expected hashes have been calculated on-chain with:
// `keccak256(abi.encode(groupPubKey, misbehaved))`
func TestCalculateDKGResultHash(t *testing.T) {
	chain := &ethereumChain{}

	var tests = map[string]struct {
		dkgResult    *relaychain.DKGResult
		expectedHash string
	}{

		"dkg result with no misbehaving members": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{0x64},
				Misbehaved:     []byte{},
			},
			expectedHash: "f1918e8562236eb17adc8502332f4c9c82bc14e19bfc0aa10ab674ff75b3d2f3",
		},
		"dkg result with misbehaving members": {
			dkgResult: &relaychain.DKGResult{
				GroupPublicKey: []byte{0x64},
				Misbehaved:     []byte{0x03, 0x05},
			},
			expectedHash: "9b84bec611298ebcd371abd418e5716f511d7ff3f086cc574a84afe01afb02ec",
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
	signature1 := common.LeftPadBytes([]byte("marry"), 65)
	signature2 := common.LeftPadBytes([]byte("had"), 65)
	signature3 := common.LeftPadBytes([]byte("a"), 65)
	signature4 := common.LeftPadBytes([]byte("little"), 65)
	signature5 := common.LeftPadBytes([]byte("lamb"), 65)

	invalidSignature := common.LeftPadBytes([]byte("invalid"), 64)

	var tests = map[string]struct {
		signaturesMap map[chain.GroupMemberIndex][]byte
		expectedError error
	}{
		"one valid signature": {
			signaturesMap: map[uint8][]byte{
				1: signature1,
			},
		},
		"five valid signatures": {
			signaturesMap: map[chain.GroupMemberIndex][]byte{
				3: signature3,
				1: signature1,
				4: signature4,
				5: signature5,
				2: signature2,
			},
		},
		"invalid signature": {
			signaturesMap: map[chain.GroupMemberIndex][]byte{
				1: signature1,
				2: invalidSignature,
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

				if len(signaturesSlice) != (ethutil.SignatureSize * len(indicesSlice)) {
					t.Errorf(
						"invalid signatures slice size\nexpected: %v\nactual:   %v\n",
						(ethutil.SignatureSize * len(indicesSlice)),
						len(signaturesSlice),
					)
				}
			}

			for i, actualMemberIndex := range indicesSlice {
				memberIndex := chain.GroupMemberIndex(actualMemberIndex.Uint64())

				actualSignature := signaturesSlice[ethutil.SignatureSize*i : ethutil.SignatureSize*(i+1)]
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
