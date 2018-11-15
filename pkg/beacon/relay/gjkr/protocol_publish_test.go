package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPrepareResult(t *testing.T) {
	threshold := 4
	groupSize := 8

	members, err := initializePublishingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("%s", err)
	}

	publishingMember := members[0]

	var tests = map[string]struct {
		disqualifiedMemberIDs []int
		inactiveMemberIDs     []int
		expectedResult        *result.Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       nil,
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []int{2},
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   []int{2},
				Inactive:       nil,
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []int{3, 7},
			expectedResult: &result.Result{
				Success:        true,
				GroupPublicKey: big.NewInt(123), // TODO: Use group public key after Phase 12 is merged
				Disqualified:   nil,
				Inactive:       []int{3, 7},
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []int{2},
			inactiveMemberIDs:     []int{3, 7},
			expectedResult: &result.Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   []int{2},
				Inactive:       nil, // in case of failure only disqualified members are slashed
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []int{3, 5, 7},
			expectedResult: &result.Result{
				Success:        false,
				GroupPublicKey: nil,
				Disqualified:   nil,
				Inactive:       nil, // in case of failure only disqualified members are slashed
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publishingMember.group.disqualifiedMemberIDs = test.disqualifiedMemberIDs
			publishingMember.group.inactiveMemberIDs = test.inactiveMemberIDs

			publishingMember.PrepareResult()

			if !reflect.DeepEqual(test.expectedResult, publishingMember.result) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, publishingMember.result)
			}
		})
	}
}

func initializePublishingMembersGroup(threshold, groupSize int) ([]*PublishingMember, error) {
	chainHandle := local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, err
	}

	initialBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	dkg := &DKG{
		chain: &Chain{
			handle:                   chainHandle,
			expectedProtocolDuration: 3,                  // T_dkg
			blockStep:                2,                  // T_step
			initialBlockHeight:       initialBlockHeight, // T_init
		},
	}

}

func initializePublishingMembersGroup(
	threshold, groupSize, expectedProtocolDuration, blockStep int,
) ([]*PublishingMember, error) {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	chain, err := initChain(threshold, groupSize, expectedProtocolDuration, blockStep)
	if err != nil {
		return nil, err
	}

	dkg := &DKG{chain: chain}

	var members []*PublishingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members,
			&PublishingMember{
				PointsJustifyingMember: &PointsJustifyingMember{
					SharingMember: &SharingMember{
						QualifiedMember: &QualifiedMember{
							SharesJustifyingMember: &SharesJustifyingMember{
								CommittingMember: &CommittingMember{
									memberCore: &memberCore{
										ID:             id,
										group:          group,
										protocolConfig: dkg,
									},
								},
							},
						},
					},
				},
			})
		group.RegisterMemberID(id)
	}
	return members, nil
}

func initChain(
	threshold, groupSize, expectedProtocolDuration, blockStep int,
) (*Chain, error) {
	chainHandle := local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, err
	}

	initialBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	return &Chain{
		handle:                   chainHandle,
		expectedProtocolDuration: expectedProtocolDuration, // T_dkg
		blockStep:                blockStep,                // T_step
		initialBlockHeight:       initialBlockHeight,       // T_init
	}, nil
}
