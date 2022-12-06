package tbtc

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

func TestDecideSigningGroupMemberFate(t *testing.T) {
	chainConfig := &ChainConfig{
		GroupSize:       10,
		GroupQuorum:     8,
		HonestThreshold: 6,
	}

	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	memberIndex := group.MemberIndex(1)

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
		ctxFn                          func() (context.Context, context.CancelFunc)
		resultSubmittedEvent           *DKGResultSubmittedEvent
		expectedOperatingMemberIndexes []group.MemberIndex
		expectedError                  error
	}{
		"member supports the published result": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{7, 10},
				BlockNumber:         5,
			},
			// should return operating members according to the published result
			expectedOperatingMemberIndexes: []group.MemberIndex{1, 2, 3, 4, 5, 6, 8, 9},
		},
		"member supports different public key": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: []byte{0x00, 0x01}, // different result
				Misbehaved:          []byte{7, 10},
				BlockNumber:         5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member does not support the same group public key",
				memberIndex,
			),
		},
		"member considered as misbehaved": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			resultSubmittedEvent: &DKGResultSubmittedEvent{
				MemberIndex:         2,
				GroupPublicKeyBytes: resultGroupPublicKeyBytes,
				Misbehaved:          []byte{memberIndex}, // member considered as misbehaved
				BlockNumber:         5,
			},
			expectedError: fmt.Errorf(
				"[member:%v] could not stay in the group because the "+
					"member is considered as misbehaving",
				memberIndex,
			),
		},
		"publication timeout exceeded": {
			ctxFn: func() (context.Context, context.CancelFunc) {
				ctx, cancelCtx := context.WithCancel(context.Background())
				// Cancel the context deliberately.
				cancelCtx()
				return ctx, cancelCtx
			},
			resultSubmittedEvent: nil, // the result is not published at all
			expectedError:        fmt.Errorf("result publication timed out"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			dkgResultChannel := make(chan *DKGResultSubmittedEvent, 1)

			if test.resultSubmittedEvent != nil {
				dkgResultChannel <- test.resultSubmittedEvent
			}

			ctx, cancelCtx := test.ctxFn()
			defer cancelCtx()

			operatingMemberIndexes, err := decideSigningGroupMemberFate(
				ctx,
				memberIndex,
				dkgResultChannel,
				result,
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
