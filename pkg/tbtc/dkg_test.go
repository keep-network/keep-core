package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain/local_v1"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"golang.org/x/crypto/sha3"
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
			expectedExecutionEndBlock: 325, // 200 + 125
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
			expectedExecutionEndBlock: 585, // 460 + 125
			expectedResult:            testResult,
			// Members 4 and 6 should be excluded in the last attempt as
			// they were inactive.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      460, // 200 + 125 + 5 + 125 + 5
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
			expectedExecutionEndBlock: 455, // 330 + 125
			expectedResult:            testResult,
			// Member 3 was inactive but excluding their operator drops the
			// group size below the quorum. We fall back to the random algorithm
			// that excludes member 4 for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           2,
				startBlock:      330, // 200 + 125 + 5
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
			expectedExecutionEndBlock: 585, // 460 + 125
			expectedResult:            testResult,
			// Since the error is not related with inactive members, we
			// use the random algorithm from the very beginning. It
			// excludes members 2 and 5 (same operator) for the given
			// seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      460, // 200 + 125 + 5 + 125 + 5
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
			expectedExecutionEndBlock: 585, // 460 + 125
			expectedResult:            testResult,
			// The random algorithm was used first so subsequent errors related
			// to inactive members are not taken into account. The random
			// algorithm excludes members 2 and 5 (same operator) for the given
			// seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      460, // 200 + 125 + 5 + 125 + 5
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
			expectedExecutionEndBlock: 975, // 850 + 125
			expectedResult:            testResult,
			// 5 attempts failed due to different single members who were inactive.
			// The 6th attempt should be made using the random retry that
			// returns member 9 for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           6,
				startBlock:      850, // 200 + 5 * (125 + 5)
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
			expectedExecutionEndBlock: 585, // 460 + 125
			expectedResult:            testResult,
			// First attempt fail due to member 2 who is inactive but the second
			// attempt fail due to another error so the random algorithm
			// should be used eventually and return members 2 and 5
			// (same operator) for the given seed.
			expectedLastAttempt: &dkgAttemptParams{
				index:           3,
				startBlock:      460, // 200 + 125 + 5 + 125 + 5
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
			expectedExecutionEndBlock: 2370, // 2245 + 125
			expectedResult:            testResult,
			// Random algorithm is used from the very beginning. We also
			// observe a delay blocks bump on the 10th attempt which is
			// 100 blocks instead of 5. That said, the start block for the 16th
			// attempt can be calculated as follows:
			// 200 + 130 + 130 + 130 + 130 + 130 + 130 + 130 + 130 + 225 + 130 + 130 + 130 + 130 + 130 + 130
			// where all 130 denotes a duration of a normal attempt (125 blocks
			// plus 5 delay blocks) and 225 is the duration of the 10th attempt
			// (125 + 100 bumped delay blocks).
			expectedLastAttempt: &dkgAttemptParams{
				index:           16,
				startBlock:      2245,
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
			expectedExecutionEndBlock: 1105, // 850 + 125 + 5 + 125
			expectedResult:            testResult,
			// Member 6 is the executing one. First 5 attempts fail and are
			// retried using the random algorithm. The 6th attempt does not
			// return an error but member 6 is excluded for this attempt so,
			// member 6 skips attempt 6 and succeeds on attempt 7.
			expectedLastAttempt: &dkgAttemptParams{
				index:           7,
				startBlock:      980, // 200 + 6 * (125 + 5)
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

func TestDecideSigningGroupMemberFate(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	blockCounter, err := local_v1.BlockCounter()
	if err != nil {
		t.Fatal(err)
	}

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	memberIndex := group.MemberIndex(1)
	publicationStartBlock := uint64(0)

	result := &dkg.Result{
		Group: group.NewGroup(
			chainConfig.GroupSize-chainConfig.HonestThreshold,
			chainConfig.GroupSize,
		),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	resultGroupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		resultSubmittedEvent           *DKGResultSubmittedEvent
		expectedOperatingMemberIndexes []group.MemberIndex
		expectedError                  error
	}{
		"member supports the published result": {
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{7, 10},
				BlockNumber:         publicationStartBlock + 5,
			},
			// should return operating members according to the published result
			expectedOperatingMemberIndexes: []group.MemberIndex{1, 2, 3, 4, 5, 6, 8, 9},
		},
		"member supports different public key": {
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: []byte{0x00, 0x01}, // different result
				Misbehaved:          []byte{7, 10},
				BlockNumber:         publicationStartBlock + 5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member does not support the same group public key",
				memberIndex,
			),
		},
		"member considered as misbehaved": {
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{memberIndex}, // member considered as misbehaved
				BlockNumber:         publicationStartBlock + 5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member is considered as misbehaving",
				memberIndex,
			),
		},
		"publication timeout exceeded": {
			resultSubmittedEvent: nil, // the result is not published at all
			expectedError:        fmt.Errorf("ECDSA DKG result publication timed out"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			dkgResultChannel := make(chan *DKGResultSubmittedEvent, 1)

			if test.resultSubmittedEvent != nil {
				dkgResultChannel <- test.resultSubmittedEvent
			}

			operatingMemberIndexes, err := decideSigningGroupMemberFate(
				memberIndex,
				dkgResultChannel,
				publicationStartBlock,
				result,
				chainConfig,
				blockCounter,
			)

			if !reflect.DeepEqual(
				test.expectedOperatingMemberIndexes,
				operatingMemberIndexes,
			) {
				t.Errorf(
					"unexpected operating member indexes\n"+
						"expected: [%v]\n"+
						"actual:   [%v]",
					test.expectedOperatingMemberIndexes,
					operatingMemberIndexes,
				)
			}

			if !reflect.DeepEqual(
				test.expectedError,
				err,
			) {
				t.Errorf(
					"unexpected error\n"+
						"expected: [%v]\n"+
						"actual:   [%v]",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestFinalSigningGroup(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       5,
		GroupQuorum:     3,
		HonestThreshold: 2,
	}

	selectedOperators := []chain.Address{
		"0xAA",
		"0xBB",
		"0xCC",
		"0xDD",
		"0xEE",
	}

	var tests = map[string]struct {
		selectedOperators           []chain.Address
		operatingMembersIndexes     []group.MemberIndex
		expectedFinalOperators      []chain.Address
		expectedFinalMembersIndexes map[group.MemberIndex]group.MemberIndex
		expectedError               error
	}{
		"selected operators count not equal to the group size": {
			selectedOperators:       selectedOperators[:4],
			operatingMembersIndexes: []group.MemberIndex{1, 2, 3, 4, 5},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
		"all selected operators are operating": {
			selectedOperators:           selectedOperators,
			operatingMembersIndexes:     []group.MemberIndex{5, 4, 3, 2, 1},
			expectedFinalOperators:      selectedOperators,
			expectedFinalMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 2: 2, 3: 3, 4: 4, 5: 5},
		},
		"honest majority of selected operators are operating": {
			selectedOperators:           selectedOperators,
			operatingMembersIndexes:     []group.MemberIndex{5, 1, 3},
			expectedFinalOperators:      []chain.Address{"0xAA", "0xCC", "0xEE"},
			expectedFinalMembersIndexes: map[group.MemberIndex]group.MemberIndex{1: 1, 3: 2, 5: 3},
		},
		"less than honest majority of selected operators are operating": {
			selectedOperators:       selectedOperators,
			operatingMembersIndexes: []group.MemberIndex{5, 1},
			expectedError:           fmt.Errorf("invalid input parameters"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualFinalOperators, actualFinalMembersIndexes, err :=
				finalSigningGroup(
					test.selectedOperators,
					test.operatingMembersIndexes,
					chainConfig,
				)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(
				test.expectedFinalOperators,
				actualFinalOperators,
			) {
				t.Errorf(
					"unexpected final operators\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedFinalOperators,
					actualFinalOperators,
				)
			}

			if !reflect.DeepEqual(
				test.expectedFinalMembersIndexes,
				actualFinalMembersIndexes,
			) {
				t.Errorf(
					"unexpected final members indexes\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedFinalMembersIndexes,
					actualFinalMembersIndexes,
				)
			}
		})
	}
}

func TestSignResult_SigningSuccessful(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}
	result := &dkg.Result{
		Group:           group.NewGroup(32, 64),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	signedResult, err := dkgResultSigner.SignResult(result)
	if err != nil {
		t.Fatal(err)
	}

	expectedPublicKey := chain.Signing().PublicKey()
	if !reflect.DeepEqual(
		expectedPublicKey,
		signedResult.PublicKey,
	) {
		t.Errorf(
			"unexpected public key\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedPublicKey,
			signedResult.PublicKey,
		)
	}

	expectedDKGResultHash := dkg.ResultHash(
		sha3.Sum256([]byte(fmt.Sprint(result))),
	)
	if expectedDKGResultHash != signedResult.ResultHash {
		t.Errorf(
			"unexpected result hash\n"+
				"expected: %v\n"+
				"actual:   %v\n",
			expectedDKGResultHash,
			signedResult.ResultHash,
		)
	}

	// Since signature is different on every run (even if the same private key
	// and result hash are used), simply verify if it's correct
	signatureVerification, err := chain.Signing().Verify(
		signedResult.ResultHash[:],
		signedResult.Signature,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !signatureVerification {
		t.Errorf(
			"Signature [0x%x] was not generated properly for the result hash "+
				"[0x%x]",
			signedResult.Signature,
			signedResult.ResultHash,
		)
	}
}

func TestSignResult_ErrorDuringDkgResultHashCalculation(t *testing.T) {
	chain := Connect(5, 4, 3)
	dkgResultSigner := newDkgResultSigner(chain)

	// Use nil as the DKG result to cause hash calculation error
	_, err := dkgResultSigner.SignResult(nil)

	expectedError := fmt.Errorf(
		"dkg result hash calculation failed [%w]",
		errNilDKGResult,
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}
