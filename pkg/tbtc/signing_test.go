package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
)

func TestSigningRetryLoop(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
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
			signingGroupMemberIndex: 4,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// The first attempt starts with all signing group members that
			// are randomly trimmed to the honest threshold count. As result,
			// members 1, 2, 3 and 7 are excluded for the first attempt.
			expectedLastAttempt: &signingAttemptParams{
				number:                 1,
				startBlock:             200,
				excludedMembersIndexes: []group.MemberIndex{1, 2, 3, 7},
			},
		},
		"IA error on initial attempts and honest threshold is maintained": {
			signingGroupMemberIndex: 8,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number == 1 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{4},
					}
				}

				if attempt.number == 2 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{6},
					}
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Members 4 and 6 were inactive so, they are excluded in the
			// last attempt. Additionally, members 5 and 9 are randomly
			// excluded in order to trim the attempt's signing group to the
			// honest threshold count.
			expectedLastAttempt: &signingAttemptParams{
				number:                 3,
				startBlock:             356, // 200 + 2 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{4, 5, 6, 9},
			},
		},
		"IA error on initial attempts and honest threshold is not maintained": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				// Members 2 and 3 are controlled by operators that hold
				// 5 members in the entire group. Excluding them will break
				// the honest threshold.
				if attempt.number == 1 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{2, 3},
					}
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Members 2 and 3 were inactive but excluding their operators drops
			// the group size below the honest threshold. We fall back to the
			// random algorithm that excludes members 3, 7, 8, and 10 for the
			// given seed.
			expectedLastAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             278, // 200 + 1 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{3, 7, 8, 10},
			},
		},
		"other error on initial attempts": {
			signingGroupMemberIndex: 3,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number == 1 || attempt.number == 2 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Since the error is not related with inactive members, we
			// use the random algorithm from the very beginning. It
			// excludes members 1, 2, 5 and 6 for the given seed.
			expectedLastAttempt: &signingAttemptParams{
				number:                 3,
				startBlock:             356, // 200 + 2 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 6},
			},
		},
		"other error then IA error on initial attempts": {
			signingGroupMemberIndex: 7,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number == 1 {
					return nil, fmt.Errorf("invalid data")
				}

				if attempt.number == 2 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{3},
					}
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// The random algorithm was used first so subsequent errors related
			// to inactive members are not taken into account. The random
			// algorithm excludes members 1, 2, 5 and 6 for the given
			// seed.
			expectedLastAttempt: &signingAttemptParams{
				number:                 3,
				startBlock:             356, // 200 + 2 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 6},
			},
		},
		"IA error on initial and later attempts": {
			signingGroupMemberIndex: 2,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				inactiveQueue := []group.MemberIndex{1, 2, 3, 4, 5}

				if attempt.number <= 5 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{
							inactiveQueue[attempt.number-1],
						},
					}
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// 5 attempts failed due to different single members who were inactive.
			// The 6th attempt should be made using the random retry that
			// excludes members 1, 4, 6, and 9 for the given seed.
			expectedLastAttempt: &signingAttemptParams{
				number:                 6,
				startBlock:             590, // 200 + 5 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 4, 6, 9},
			},
		},
		"IA error then other error on initial attempts": {
			signingGroupMemberIndex: 3,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number == 1 {
					return nil, &signing.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{9},
					}
				}

				if attempt.number == 2 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// First attempt fail due to member 9 who is inactive but the second
			// attempt fail due to another error so the random algorithm
			// should be used eventually and excludes members 1, 2, 5, and 6
			// for the given seed.
			expectedLastAttempt: &signingAttemptParams{
				number:                 3,
				startBlock:             356, // 200 + 2 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{1, 2, 5, 6},
			},
		},
		"other error on initial and later attempts": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				if attempt.number <= 15 {
					return nil, fmt.Errorf("invalid data")
				}

				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Random algorithm is used from the very beginning. The start block
			// for the 16th attempt can be calculated as follows: 200 + 15 * (73 + 5)
			// where 78 denotes a duration of an attempt (73 blocks plus 5
			// delay blocks).
			expectedLastAttempt: &signingAttemptParams{
				number:                 16,
				startBlock:             1370,
				excludedMembersIndexes: []group.MemberIndex{2, 5, 6, 7},
			},
		},
		"executing member excluded": {
			signingGroupMemberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			signingAttemptFn: func(attempt *signingAttemptParams) (*signing.Result, error) {
				return testResult, nil
			},
			expectedErr:    nil,
			expectedResult: testResult,
			// Member 1 is the executing one. The first attempt starts with all
			// signing group members that are randomly trimmed to the honest
			// threshold count. Member 1 is part of the excluded members list
			// so, it skips the first attempt and ends on attempt 2.
			expectedLastAttempt: &signingAttemptParams{
				number:                 2,
				startBlock:             278, // 200 + 1 * (73 + 5)
				excludedMembersIndexes: []group.MemberIndex{3, 7, 8, 10},
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
			expectedErr:         nil,
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

			var lastAttemptStartBlock uint64
			var lastAttempt *signingAttemptParams

			result, err := retryLoop.start(
				ctx,
				func(ctx context.Context, attemptStartBlock uint64) error {
					lastAttemptStartBlock = attemptStartBlock
					return nil
				},
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

			if test.expectedLastAttempt != nil {
				if test.expectedLastAttempt.startBlock != lastAttemptStartBlock {
					t.Errorf("unexpected last attempt start block\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
						test.expectedLastAttempt.startBlock,
						lastAttemptStartBlock,
					)
				}
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
