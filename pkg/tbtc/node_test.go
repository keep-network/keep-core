package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

func TestDkgRetryLoop(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	selectedOperators := chain.Addresses{
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

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	testResult := &dkg.Result{
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	var tests = map[string]struct {
		memberIndex               group.MemberIndex
		ctxFn                     func() (context.Context, context.CancelFunc)
		dkgAttemptFn              dkgAttemptFn
		expectedErr               error
		expectedExecutionEndBlock uint64
		expectedResult            *dkg.Result
		expectedLastAttempt       *dkgAttemptParams
	}{
		"success on initial attempt": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 310, // 200 + 110
			expectedResult:            testResult,
			expectedLastAttempt: &dkgAttemptParams{
				index:           1,
				startBlock:      200,
				excludedMembers: []group.MemberIndex{},
			},
		},
		"IA error on initial attempts and quorum is maintained": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index == 1 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{4},
					}
				}

				if attempt.index == 2 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{6},
					}
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 540, // 430 + 110 = 540
			expectedResult:            testResult,
			// Members 4 and 6 should be excluded in the last attempt as
			// they were inactive.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      430, // 200 + 110 + 5 + 110 + 5
				excludedMembers: []group.MemberIndex{4, 6},
			},
		},
		"IA error on initial attempts and quorum is not maintained": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				// Member 3 is controlled by an operator that controls 3 members
				// in total. Excluding that operator will break the group quorum.
				if attempt.index == 1 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{3},
					}
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 425, // 315 + 110
			expectedResult:            testResult,
			// Member 3 was inactive but excluding their operator drops the
			// group size below the quorum. We fall back to the random algorithm
			// that excludes member 4 for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           2,
				startBlock:      315, // 200 + 110 + 5
				excludedMembers: []group.MemberIndex{4},
			},
		},
		"other error on initial attempts": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index == 1 || attempt.index == 2 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 540, // 430 + 110
			expectedResult:            testResult,
			// Since the error is not related with inactive members, we
			// use the random algorithm from the very beginning. It
			// excludes members 2 and 5 (same operator) for the given
			// seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      430, // 200 + 110 + 5 + 110 + 5
				excludedMembers: []group.MemberIndex{2, 5},
			},
		},
		"other error then IA error on initial attempts": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index == 1 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				if attempt.index == 2 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{3},
					}
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 540, // 430 + 110
			expectedResult:            testResult,
			// The random algorithm was used first so subsequent errors related
			// to inactive members are not taken into account. The random
			// algorithm excludes members 2 and 5 (same operator) for the given
			// seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      430, // 200 + 110 + 5 + 110 + 5
				excludedMembers: []group.MemberIndex{2, 5},
			},
		},
		"IA error on initial and later attempts": {
			memberIndex: 2,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				inactiveQueue := []group.MemberIndex{1, 4, 6, 7, 9}

				if attempt.index <= 5 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{
							inactiveQueue[attempt.index-1],
						},
					}
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 885, // 775 + 110
			expectedResult:            testResult,
			// 5 attempts failed due to different single members who were inactive.
			// The 6th attempt should be made using the random retry that
			// returns member 9 for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           6,
				startBlock:      775, // 200 + 5 * (110 + 5)
				excludedMembers: []group.MemberIndex{9},
			},
		},
		"IA error then other error on initial attempts": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index == 1 {
					return nil, 0, &dkg.InactiveMembersError{
						InactiveMembersIndexes: []group.MemberIndex{2},
					}
				}

				if attempt.index == 2 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 540, // 430 + 110
			expectedResult:            testResult,
			// First attempt fail due to member 2 who is inactive but the second
			// attempt fail due to another error so the random algorithm
			// should be used eventually and return members 2 and 5
			// (same operator) for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      430, // 200 + 110 + 5 + 110 + 5
				excludedMembers: []group.MemberIndex{2, 5},
			},
		},
		"other error on initial and later attempts": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index <= 15 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 2130, // 2020 + 110
			expectedResult:            testResult,
			// Random algorithm is used from the very beginning. We also
			// observe a delay blocks bump on the 10th attempt which is
			// 100 blocks instead of 5. That said, the start block for the 16th
			// attempt can be calculated as follows:
			// 200 + 115 + 115 + 115 + 115 + 115 + 115 + 115 + 115 + 210 + 115 + 115 + 115 + 115 + 115 + 115
			// where all 115 denotes a duration of a normal attempt (110 blocks
			// plus 5 delay blocks) and 210 is the duration of the 10th attempt
			// (110 + 100 bumped delay blocks).
			expectedLastAttempt: &dkgAttemptParams{
				index:           16,
				startBlock:      2020,
				excludedMembers: []group.MemberIndex{7, 9},
			},
		},
		"executing member excluded": {
			memberIndex: 6,
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				if attempt.index <= 5 {
					return nil, 0, fmt.Errorf("invalid data")
				}

				return testResult, attempt.startBlock + dkg.ProtocolBlocks(), nil
			},
			expectedErr:               nil,
			expectedExecutionEndBlock: 1000, // 775+110+5+110
			expectedResult:            testResult,
			// Member 6 is the executing one. First 5 attempts fail and are
			// retried using the random algorithm. The 6th attempt does not
			// return an error but member 6 is excluded for this attempt so,
			// member 6 skips attempt 6 and succeeds on attempt 7.
			expectedLastAttempt: &dkgAttemptParams{
				index:           7,
				startBlock:      890, // 200 + 6 * (110 + 5)
				excludedMembers: []group.MemberIndex{7},
			},
		},
		"loop context done": {
			memberIndex: 1,
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			dkgAttemptFn: func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
				return nil, 0, fmt.Errorf("invalid data")
			},
			expectedErr:               fmt.Errorf("dkg retry loop received stop signal on attempt [1]"),
			expectedResult:            nil,
			expectedExecutionEndBlock: 0,
			expectedLastAttempt:       nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			retryLoop := newDkgRetryLoop(
				big.NewInt(100),
				200,
				test.memberIndex,
				selectedOperators,
				chainConfig,
			)

			// Given the small group size, we never reach the original
			// bump frequency which is 100. Here we make it smaller in order
			// to test its behavior.
			retryLoop.delayBlocksBumpFrequency = 10

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			var lastAttempt *dkgAttemptParams

			result, executionEndBlock, err := retryLoop.start(
				ctx,
				func(params *dkgAttemptParams) (*dkg.Result, uint64, error) {
					lastAttempt = params
					return test.dkgAttemptFn(params)
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

			if test.expectedExecutionEndBlock != executionEndBlock {
				t.Errorf(
					"unexpected execution end block\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					test.expectedExecutionEndBlock,
					executionEndBlock,
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
