package gjkr

import (
	"github.com/keep-network/keep-core/pkg/chain/local"
)

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

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

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
