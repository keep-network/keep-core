package gjkr

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestCombineReceivedShares(t *testing.T) {
	selfShareS := big.NewInt(9)
	selfShareT := big.NewInt(19)
	q := big.NewInt(59)

	receivedShareS := make(map[int]*big.Int)
	receivedShareT := make(map[int]*big.Int)
	// Simulate shares received from peer members.
	// Peer members IDs are in [100, 101, 102, 103, 104, 105] to differ them from
	// slice indices.
	for i := 0; i <= 5; i++ {
		receivedShareS[100+i] = big.NewInt(int64(10 + i))
		receivedShareT[100+i] = big.NewInt(int64(20 + i))
	}

	// 9 + 10 + 11 + 12 + 13 + 14 + 15 = 84 mod 59 = 25
	expectedShareS := big.NewInt(25)
	// 19 + 20 + 21 + 22 + 23 + 24 + 25 = 154 mod 59 = 36
	expectedShareT := big.NewInt(36)

	config := &DKG{Q: q}
	member := &QualifiedMember{
		SharesJustifyingMember: &SharesJustifyingMember{
			CommittingMember: &CommittingMember{
				memberCore: &memberCore{
					protocolConfig: config,
				},
				selfSecretShareS:     selfShareS,
				selfSecretShareT:     selfShareT,
				receivedValidSharesS: receivedShareS,
				receivedValidSharesT: receivedShareT,
			},
		},
	}

	member.CombineMemberShares()

	if member.masterPrivateKeyShare.Cmp(expectedShareS) != 0 {
		t.Errorf("incorrect combined shares S value\nexpected: %v\nactual:   %v\n",
			expectedShareS,
			member.masterPrivateKeyShare,
		)
	}
	if member.shareT.Cmp(expectedShareT) != 0 {
		t.Errorf("incorrect combined shares T value\nexpected: %v\nactual:   %v\n",
			expectedShareT,
			member.shareT,
		)
	}
}

func TestCalculatePublicCoefficients(t *testing.T) {
	secretCoefficients := []*big.Int{
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(2),
	}
	expectedPublicCoefficients := []*big.Int{
		big.NewInt(216), // 6^3 mod 1907 = 216
		big.NewInt(148), // 6^5 mod 1907 = 148
		big.NewInt(36),  // 6^2 mod 1907 = 36
	}

	config := &DKG{P: big.NewInt(1907), Q: big.NewInt(953)}

	// This test uses rand.Reader mock to get specific `g` value in `NewVSS`
	// initialization.
	mockRandomReader := testutils.NewMockRandReader(big.NewInt(6))
	vss, err := pedersen.NewVSS(mockRandomReader, config.P, config.Q)
	if err != nil {
		t.Fatalf("VSS initialization failed [%s]", err)
	}

	member := &SharingMember{
		QualifiedMember: &QualifiedMember{
			SharesJustifyingMember: &SharesJustifyingMember{
				CommittingMember: &CommittingMember{
					memberCore: &memberCore{
						protocolConfig: config,
					},
					vss:                vss,
					secretCoefficients: secretCoefficients,
				},
			},
		},
	}

	message := member.CalculatePublicKeySharePoints()

	if !reflect.DeepEqual(member.publicKeySharePoints, expectedPublicCoefficients) {
		t.Errorf("incorrect member's public shares\nexpected: %v\nactual:   %v\n",
			expectedPublicCoefficients,
			member.publicKeySharePoints,
		)
	}

	if !reflect.DeepEqual(message.publicKeySharePoints, expectedPublicCoefficients) {
		t.Errorf("incorrect public shares in message\nexpected: %v\nactual:   %v\n",
			expectedPublicCoefficients,
			message.publicKeySharePoints,
		)
	}
}

func TestCalculateAndVerifyPublicKeySharePoints(t *testing.T) {
	threshold := 3
	groupSize := 5

	sharingMembers, err := initializeSharingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("Group initialization failed [%s]", err)
	}

	sharingMember := sharingMembers[0]

	var tests = map[string]struct {
		modifyPublicKeySharePointsMessages func(messages []*MemberPublicKeySharePointsMessage)
		expectedError                      error
		expectedAccusedIDs                 []int
	}{
		"positive validation - no accusations": {
			expectedError: nil,
		},
		"negative validation - changed public key share - one accused member": {
			modifyPublicKeySharePointsMessages: func(messages []*MemberPublicKeySharePointsMessage) {
				messages[1].publicKeySharePoints[1] = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{3},
		},
		"negative validation - changed public key share - two accused members": {
			modifyPublicKeySharePointsMessages: func(messages []*MemberPublicKeySharePointsMessage) {
				messages[0].publicKeySharePoints[1] = big.NewInt(13)
				messages[3].publicKeySharePoints[1] = big.NewInt(18)
			},
			expectedError:      nil,
			expectedAccusedIDs: []int{2, 5},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			messages := make([]*MemberPublicKeySharePointsMessage, groupSize)

			for i, m := range sharingMembers {
				messages[i] = m.CalculatePublicKeySharePoints()
			}

			filteredMessages := filterMemberPublicKeySharePointsMessages(
				messages,
				sharingMember.ID,
			)

			if test.modifyPublicKeySharePointsMessages != nil {
				test.modifyPublicKeySharePointsMessages(filteredMessages)
			}

			accusedMessage, err := sharingMember.VerifyPublicKeySharePoints(filteredMessages)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"\nexpected: %s\nactual:   %s\n",
					test.expectedError,
					err,
				)
			}

			if !reflect.DeepEqual(accusedMessage.accusedIDs, test.expectedAccusedIDs) {
				t.Fatalf("incorrect accused IDs\nexpected: %v\nactual:   %v\n",
					test.expectedAccusedIDs,
					accusedMessage.accusedIDs,
				)
			}
		})
	}
}

func initializeSharingMembersGroup(threshold, groupSize int) ([]*SharingMember, error) {
	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, fmt.Errorf("group initialization failed [%s]", err)
	}

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {
		cm.secretCoefficients = make([]*big.Int, threshold+1)
		for i := 0; i < threshold+1; i++ {
			cm.secretCoefficients[i], err = crand.Int(crand.Reader, cm.protocolConfig.Q)
			if err != nil {
				return nil, fmt.Errorf("secret share generation failed [%s]", err)
			}
		}
		sharingMembers = append(sharingMembers, &SharingMember{
			QualifiedMember: &QualifiedMember{
				SharesJustifyingMember: &SharesJustifyingMember{
					CommittingMember: cm,
				},
			},
			receivedValidPeerPublicKeySharePoints: make(map[int][]*big.Int, groupSize-1),
		})
	}

	for _, sm := range sharingMembers {
		for _, cm := range committingMembers {
			sm.receivedValidSharesS[cm.ID] = evaluateMemberShare(sm.ID, cm.secretCoefficients)
		}
	}

	return sharingMembers, nil
}

func filterMemberPublicKeySharePointsMessages(
	messages []*MemberPublicKeySharePointsMessage, receiverID int,
) []*MemberPublicKeySharePointsMessage {
	var result []*MemberPublicKeySharePointsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
