package tecdsa

import (
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"math/big"
	"testing"
)

func TestSignatureEquals(t *testing.T) {
	var tests = map[string]struct {
		signature1     *Signature
		signature2     *Signature
		expectedResult bool
	}{
		"both signatures are nil": {
			signature1:     nil,
			signature2:     nil,
			expectedResult: true,
		},
		"first signature is nil": {
			signature1: nil,
			signature2: &Signature{
				R:          nil,
				S:          nil,
				RecoveryID: 0,
			},
			expectedResult: false,
		},
		"second signature is nil": {
			signature1: &Signature{
				R:          nil,
				S:          nil,
				RecoveryID: 0,
			},
			signature2:     nil,
			expectedResult: false,
		},
		"both signatures' fields are non-nil and same": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: true,
		},
		"both signatures' R are nil": {
			signature1: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: true,
		},
		"both signatures' R are nil but recovery IDs are different": {
			signature1: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 3,
			},
			expectedResult: false,
		},
		"first signature's R is nil": {
			signature1: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"second signature's R is nil": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          nil,
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"both signatures' R are non-nil but different": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(101),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"both signatures' S are nil": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 2,
			},
			expectedResult: true,
		},
		"both signatures' S are nil but recovery IDs are different": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 3,
			},
			expectedResult: false,
		},
		"first signature's S is nil": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"second signature's S is nil": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          nil,
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"both signatures' S are non-nil but different": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(201),
				RecoveryID: 2,
			},
			expectedResult: false,
		},
		"both signatures' recovery IDs are different": {
			signature1: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 2,
			},
			signature2: &Signature{
				R:          big.NewInt(100),
				S:          big.NewInt(200),
				RecoveryID: 3,
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.signature1.Equals(test.signature2)

			testutils.AssertBoolsEqual(
				t,
				"equals result",
				test.expectedResult,
				result,
			)
		})
	}
}
