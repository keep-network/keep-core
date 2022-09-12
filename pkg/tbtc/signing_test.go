package tbtc

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"math/big"
	"reflect"
	"testing"
)

func TestSigningRetryLoop(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	signingGroupOperators := chain.Addresses{
		"address-1",
		"address-2",
		"address-8",
		"address-4",
		"address-2",
		"address-6",
		"address-7",
		"address-8",
		"address-9",
		"address-8",
	}

	testResult := &signing.Result{
		Signature: &tecdsa.Signature{
			R:          big.NewInt(300),
			S:          big.NewInt(400),
			RecoveryID: 2,
		},
	}

	var tests = map[string]struct {
		signingGroupMemberIndex group.MemberIndex
		ctxFn                   func() (context.Context, context.CancelFunc)
		signingAttemptFn        signingAttemptFn
		expectedErr             error
		expectedResult          *signing.Result
		expectedLastAttempt     *signingAttemptParams
	}{
		"success on initial attempt": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// The signing random retry algorithm invoked with the test seed
			// excludes 4 members (6 is the honest threshold) from the first
			// attempt: 3, 7, 8 and 10.
			expectedLastAttempt: &signingAttemptParams{
				index:           1,
				startBlock:      200,
				excludedMembers: []group.MemberIndex{3, 7, 8, 10},
			},
		},
		"error on initial attempt": {
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.index <= 1 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 4 is the executing one. The first attempt fails and
			// the signing random retry algorithm invoked with the test seed
			// excludes 3 members (6 is the honest threshold) from the second
			// attempt: 1, 2 and 5. The additional exclusion round that trims
			// the included members list to the honest threshold size adds
			// member 7 to the final excluded members list.
			expectedLastAttempt: &signingAttemptParams{
				index:           2,
				startBlock:      278, // 200 + 1 * (73 + 5)
				excludedMembers: []group.MemberIndex{1, 2, 5, 7},
			},
		},
		"executing member excluded": {
			signingGroupMemberIndex: 2,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.index <= 5 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 2 is the executing one. First 5 attempts fail and are
			// retried using the random algorithm. The 6th attempt does not
			// return an error but member 2 is excluded for this attempt so,
			// member 2 skips attempt 6 and ends on attempt 7.
			expectedLastAttempt: &signingAttemptParams{
				index:           7,
				startBlock:      668, // 200 + 6 * (73 + 5)
				excludedMembers: []group.MemberIndex{1, 4, 6, 9},
			},
		},
		"loop context done": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return nil, fmt.Errorf("invalid data")
			},
			expectedErr:         fmt.Errorf("signing retry loop received stop signal on attempt [1]"),
			expectedResult:      nil,
			expectedLastAttempt: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			retryLoop := newSigningRetryLoop(
				big.NewInt(100),
				200,
				test.signingGroupMemberIndex,
				signingGroupOperators,
				chainConfig,
			)

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			var lastAttempt *signingAttemptParams

			result, err := retryLoop.start(
				ctx,
				func(params *signingAttemptParams) (*signing.Result, error) {
					lastAttempt = params
					return test.signingAttemptFn(params)
				},
			)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedErr,
					err,
				)
			}

			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf(
					"unexpected result\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedResult,
					result,
				)
			}

			if !reflect.DeepEqual(test.expectedLastAttempt, lastAttempt) {
				t.Errorf(
					"unexpected last attempt\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedLastAttempt,
					lastAttempt,
				)
			}
		})
	}
}
