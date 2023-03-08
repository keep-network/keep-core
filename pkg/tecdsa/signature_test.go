package tecdsa

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestSignatureString(t *testing.T) {
	var tests = map[string]struct {
		signature      *Signature
		expectedResult string
	}{
		"R and S values with leading zeros": {
			signature: &Signature{
				R:          big.NewInt(495),
				S:          big.NewInt(253885074),
				RecoveryID: 2,
			},
			expectedResult: "R: 0x01ef, S: 0x0f21fa92, RecoveryID: 2",
		},
		"R and S values with trailing zeros": {
			signature: &Signature{
				R:          big.NewInt(1308688384),
				S:          big.NewInt(4096),
				RecoveryID: 1,
			},
			expectedResult: "R: 0x4e010000, S: 0x1000, RecoveryID: 1",
		},
		"R and S values with no leading nor trailing zero": {
			signature: &Signature{
				R:          big.NewInt(12300),
				S:          big.NewInt(231),
				RecoveryID: 0,
			},
			expectedResult: "R: 0x300c, S: 0xe7, RecoveryID: 0",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.signature.String()

			testutils.AssertStringsEqual(
				t,
				"equals result",
				test.expectedResult,
				result,
			)
		})
	}
}

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
