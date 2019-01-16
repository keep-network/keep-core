package dkg2

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
)

func TestConvertResult(t *testing.T) {
	groupSize := 5

	publicKey := new(bn256.G1).ScalarBaseMult(big.NewInt(2))
	var serializedPublicKey [32]byte
	copy(serializedPublicKey[:], publicKey.Marshal())

	var tests = map[string]struct {
		gjkrResult     *gjkr.Result
		expectedResult *relayChain.DKGResult
	}{
		"success: false, group public key: nil, DQ and IA: empty": {
			gjkrResult: &gjkr.Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []gjkr.MemberID{},
				Inactive:       []gjkr.MemberID{},
			},
			expectedResult: &relayChain.DKGResult{
				Success:        false,
				GroupPublicKey: [32]byte{},
				Disqualified:   []byte{0x00, 0x00, 0x00, 0x00, 0x00},
				Inactive:       []byte{0x00, 0x00, 0x00, 0x00, 0x00},
			},
		},
		"success: true, group public key: provided, DQ and IA: provided": {
			gjkrResult: &gjkr.Result{
				Success:        true,
				GroupPublicKey: publicKey,
				Disqualified:   []gjkr.MemberID{1, 3, 4},
				Inactive:       []gjkr.MemberID{5},
			},
			expectedResult: &relayChain.DKGResult{
				Success:        true,
				GroupPublicKey: serializedPublicKey,
				Disqualified:   []byte{0x01, 0x00, 0x01, 0x01, 0x00},
				Inactive:       []byte{0x00, 0x00, 0x00, 0x00, 0x01},
			},
		},
	}
	for _, test := range tests {
		convertedResult := convertResult(test.gjkrResult, groupSize)

		if !test.expectedResult.Equals(convertedResult) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, convertedResult)
		}
	}
}
