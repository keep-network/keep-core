package tbtc

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

func TestDkgExecutor_RegisterSigner(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	const (
		groupSize          = 5
		groupQuorum        = 3
		honestThreshold    = 2
		dishonestThreshold = 3
	)

	localChain := Connect(groupSize, groupQuorum, honestThreshold)

	selectedOperators := []chain.Address{
		"0xAA",
		"0xBB",
		"0xCC",
		"0xDD",
		"0xEE",
	}

	var tests = map[string]struct {
		memberIndex               group.MemberIndex
		disqualifiedMemberIndexes []group.MemberIndex
		inactiveMemberIndexes     []group.MemberIndex

		expectedError                      error
		expectedFinalSigningGroupIndex     group.MemberIndex
		expectedFinalSigningGroupOperators []chain.Address
	}{
		"all members participating": {
			memberIndex:                        1,
			disqualifiedMemberIndexes:          nil,
			inactiveMemberIndexes:              nil,
			expectedFinalSigningGroupIndex:     1,
			expectedFinalSigningGroupOperators: selectedOperators,
		},
		"some member inactive": {
			memberIndex:                        3,
			disqualifiedMemberIndexes:          nil,
			inactiveMemberIndexes:              []group.MemberIndex{2, 5},
			expectedFinalSigningGroupIndex:     2,
			expectedFinalSigningGroupOperators: []chain.Address{"0xAA", "0xCC", "0xDD"},
		},
		"some members disqualified": {
			memberIndex:                        1,
			disqualifiedMemberIndexes:          []group.MemberIndex{2, 5},
			inactiveMemberIndexes:              nil,
			expectedError:                      nil,
			expectedFinalSigningGroupIndex:     1,
			expectedFinalSigningGroupOperators: []chain.Address{"0xAA", "0xCC", "0xDD"},
		},
		"the current member inactive": {
			memberIndex:               2,
			disqualifiedMemberIndexes: nil,
			inactiveMemberIndexes:     []group.MemberIndex{2, 5},
			expectedError:             fmt.Errorf("failed to resolve final signing group member index"),
		},
		"the current member disqualified": {
			memberIndex:               5,
			disqualifiedMemberIndexes: []group.MemberIndex{2, 5},
			inactiveMemberIndexes:     nil,
			expectedError:             fmt.Errorf("failed to resolve final signing group member index"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			persistenceHandle := &mockPersistenceHandle{}
			walletRegistry := newWalletRegistry(persistenceHandle)

			dkgExecutor := &dkgExecutor{
				// setting only the fields really needed for this test
				chain:          localChain,
				walletRegistry: walletRegistry,
			}

			group := group.NewGroup(dishonestThreshold, groupSize)
			for _, disqualifiedMember := range test.disqualifiedMemberIndexes {
				group.MarkMemberAsDisqualified(disqualifiedMember)
			}
			for _, inactiveMember := range test.inactiveMemberIndexes {
				group.MarkMemberAsInactive(inactiveMember)
			}

			result := &dkg.Result{
				Group:           group,
				PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
			}

			signer, err := dkgExecutor.registerSigner(result, test.memberIndex, selectedOperators)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedError,
					err,
				)
			}

			if test.expectedError != nil {
				if signer != nil {
					t.Errorf("expected nil signer")
				}

				// do not check the rest of assertions, the signer should be nil
				return
			}

			testutils.AssertIntsEqual(
				t,
				"final signing group index",
				int(test.expectedFinalSigningGroupIndex),
				int(signer.signingGroupMemberIndex),
			)

			if !reflect.DeepEqual(
				test.expectedFinalSigningGroupOperators,
				signer.wallet.signingGroupOperators,
			) {
				t.Errorf(
					"unexpected final signing group operators\n"+
						"expected: %v\n"+
						"actual:   %v\n",
					test.expectedFinalSigningGroupOperators,
					signer.wallet.signingGroupOperators,
				)
			}

			registeredSigners := walletRegistry.getSigners(
				result.PrivateKeyShare.PublicKey(),
			)

			testutils.AssertIntsEqual(
				t,
				"number of signers registered",
				1,
				len(registeredSigners),
			)
		})
	}
}

// TODO: Refactor this test scenario and add more paths.
func TestDkgExecutor_ExecuteDkgValidation(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	const (
		groupSize       = 5
		groupQuorum     = 3
		honestThreshold = 2
	)

	localChain := Connect(groupSize, groupQuorum, honestThreshold)

	operatorID, operatorAddress, err := localChain.operator()
	if err != nil {
		t.Fatal(err)
	}

	signatures := make(map[group.MemberIndex][]byte)
	var operatorsIDs chain.OperatorIDs
	var operatorsAddresses chain.Addresses

	for memberIndex := uint8(1); memberIndex <= groupSize; memberIndex++ {
		signatures[memberIndex] = []byte{memberIndex}
		operatorsIDs = append(operatorsIDs, operatorID)
		operatorsAddresses = append(operatorsAddresses, operatorAddress)
	}

	submitterMemberIndex := group.MemberIndex(2)

	tecdsaDkgResult := &dkg.Result{
		Group:           group.NewGroup(groupSize-honestThreshold, groupSize),
		PrivateKeyShare: tecdsa.NewPrivateKeyShare(testData[0]),
	}

	groupSelectionResult := &GroupSelectionResult{
		OperatorsIDs:       operatorsIDs,
		OperatorsAddresses: operatorsAddresses,
	}

	dkgResultSubmittedEventChan := make(chan *DKGResultSubmittedEvent, 1)
	_ = localChain.OnDKGResultSubmitted(
		func(event *DKGResultSubmittedEvent) {
			dkgResultSubmittedEventChan <- event
		},
	)

	err = localChain.startDKG()
	if err != nil {
		t.Fatal(err)
	}

	err = localChain.SubmitDKGResult(
		submitterMemberIndex,
		tecdsaDkgResult,
		signatures,
		groupSelectionResult,
	)
	if err != nil {
		t.Fatal(err)
	}

	dkgResultSubmittedEvent := <-dkgResultSubmittedEventChan

	waitForBlockFn := func(ctx context.Context, block uint64) error {
		blockCounter, err := localChain.BlockCounter()
		if err != nil {
			return err
		}

		wait, err := blockCounter.BlockHeightWaiter(block)
		if err != nil {
			return err
		}

		select {
		case <-wait:
		case <-ctx.Done():
		}

		return nil
	}

	// Setting only the fields really needed for this test.
	dkgExecutor := &dkgExecutor{
		operatorID:      operatorID,
		operatorAddress: operatorAddress,
		chain:           localChain,
		waitForBlockFn:  waitForBlockFn,
	}

	dkgResultApprovedEventChan := make(chan *DKGResultApprovedEvent, 1)
	_ = localChain.OnDKGResultApproved(
		func(event *DKGResultApprovedEvent) {
			dkgResultApprovedEventChan <- event
		},
	)

	dkgExecutor.executeDkgValidation(
		dkgResultSubmittedEvent.Seed,
		dkgResultSubmittedEvent.BlockNumber,
		dkgResultSubmittedEvent.Result,
		dkgResultSubmittedEvent.ResultHash,
	)

	var dkgResultApprovedEvent *DKGResultApprovedEvent
	select {
	case dkgResultApprovedEvent = <-dkgResultApprovedEventChan:
	case <-time.After(1 * time.Minute):
	}

	if dkgResultApprovedEvent == nil {
		t.Errorf("expected approval did not happen")
	}

	if dkgResultApprovedEvent != nil {
		testutils.AssertIntsEqual(
			t,
			"approval block",
			21,
			int(dkgResultApprovedEvent.BlockNumber),
		)
	}

	dkgState, err := localChain.GetDKGState()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertIntsEqual(t, "DKG state", int(Idle), int(dkgState))
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
