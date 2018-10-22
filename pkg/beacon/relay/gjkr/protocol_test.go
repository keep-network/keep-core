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

func TestCalculateSharesAndCommitments(t *testing.T) {
	threshold := 3
	groupSize := 5

	members, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member := members[0]
	peerSharesMessages, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
	if err != nil {
		t.Fatalf("shares and commitments calculation failed [%s]", err)
	}

	if len(member.secretCoefficients) != (threshold + 1) {
		t.Fatalf("\nexpected: secret coefficients number %v\nactual:   %v\n",
			threshold+1,
			len(member.secretCoefficients),
		)
	}
	if len(peerSharesMessages) != (groupSize - 1) {
		t.Fatalf("\nexpected: peer shares messages number %v\nactual:   %v\n",
			groupSize-1,
			len(peerSharesMessages),
		)
	}

	if len(commitmentsMessage.commitments) != (threshold + 1) {
		t.Fatalf("\nexpected: calculated commitments number %v\nactual:   %v\n",
			threshold+1,
			len(commitmentsMessage.commitments),
		)
	}
}

func TestSharesAndCommitmentsCalculationAndVerification(t *testing.T) {
	threshold := 3
	groupSize := 5

	var tests = map[string]struct {
		modifyPeerShareMessages   func(messages []*PeerSharesMessage) []*PeerSharesMessage
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage) []*MemberCommitmentsMessage
		expectedError             error
		expectedAccusedIDs        int
	}{
		"positive validation": {
			expectedError:      nil,
			expectedAccusedIDs: 0,
		},
		"negative validation - changed random share": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) []*PeerSharesMessage {
				messages[0].shareT = big.NewInt(13)
				return messages
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) []*MemberCommitmentsMessage {
				messages[0].commitments[0] = big.NewInt(13)
				return messages
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeCommittingMembersGroup(threshold, groupSize)
			if err != nil {
				t.Fatalf("group initialization failed [%s]", err)
			}
			currentMember := members[0]

			var peerSharesMessages []*PeerSharesMessage
			var commitmentsMessages []*MemberCommitmentsMessage

			for _, member := range members {
				peerSharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
				if err != nil {
					t.Fatalf("shares and commitments calculation failed [%s]", err)
				}
				peerSharesMessages = append(peerSharesMessages, peerSharesMessage...)
				commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
			}

			filteredPeerSharesMessages := filterPeerSharesMessage(peerSharesMessages, currentMember.ID)
			filteredMemberCommitmentsMessages := filterMemberCommitmentsMessages(commitmentsMessages, currentMember.ID)

			if test.modifyPeerShareMessages != nil {
				filteredPeerSharesMessages = test.modifyPeerShareMessages(filteredPeerSharesMessages)
			}
			if test.modifyCommitmentsMessages != nil {
				filteredMemberCommitmentsMessages = test.modifyCommitmentsMessages(filteredMemberCommitmentsMessages)
			}

			accusedMessage, err := currentMember.VerifyReceivedSharesAndCommitmentsMessages(
				filteredPeerSharesMessages,
				filteredMemberCommitmentsMessages,
			)
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"\nexpected: [%v]\nactual:   [%v]\n",
					test.expectedError,
					err,
				)
			}

			if len(accusedMessage.accusedIDs) != test.expectedAccusedIDs {
				t.Fatalf("\nexpected: accused member's IDs %v\nactual:   %v\n",
					test.expectedAccusedIDs,
					accusedMessage.accusedIDs,
				)
			}

			expectedReceivedSharesLength := groupSize - 1 - test.expectedAccusedIDs
			if len(currentMember.receivedSharesS) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: received shares S %v\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedSharesS),
				)
			}
			if len(currentMember.receivedSharesT) != expectedReceivedSharesLength {
				t.Fatalf("\nexpected: received shares T %v\nactual:   %v\n",
					expectedReceivedSharesLength,
					len(currentMember.receivedSharesT),
				)
			}
		})
	}
}

func TestCombineReceivedShares(t *testing.T) {
	receivedShareS := make(map[int]*big.Int)
	receivedShareT := make(map[int]*big.Int)
	for i := 0; i <= 5; i++ {
		receivedShareS[100+i] = big.NewInt(int64(i))
		receivedShareT[100+i] = big.NewInt(int64(10 + i))
	}

	expectedShareS := big.NewInt(15)
	expectedShareT := big.NewInt(75)

	config, err := predefinedDKGconfig()
	if err != nil {
		t.Fatalf("DKG Config initialization failed [%s]", err)
	}
	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			receivedSharesS: receivedShareS,
			receivedSharesT: receivedShareT,
		},
	}

	member.CombineReceivedShares()

	if member.shareS.Cmp(expectedShareS) != 0 {
		t.Errorf("combined shares s %s doesn't match expected %s",
			member.shareS,
			expectedShareS)
	}
	if member.shareT.Cmp(expectedShareT) != 0 {
		t.Errorf("combined shares t %s doesn't match expected %s",
			member.shareT,
			expectedShareT)
	}
}

func TestCalculatePublicCoefficients(t *testing.T) {
	secretCoefficients := []*big.Int{
		big.NewInt(3),
		big.NewInt(5),
		big.NewInt(2),
	}
	expectedPublicCoefficients := []*big.Int{
		big.NewInt(216),
		big.NewInt(148),
		big.NewInt(36),
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
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			vss:                vss,
			secretCoefficients: secretCoefficients,
		},
	}

	message := member.CalculatePublicCoefficients()

	if !reflect.DeepEqual(member.publicCoefficients, expectedPublicCoefficients) {
		t.Errorf("public shares for member doesn't match expected\nactual: %s\nexpected: %s",
			member.publicCoefficients,
			expectedPublicCoefficients,
		)
	}

	if !reflect.DeepEqual(message.publicCoefficients, expectedPublicCoefficients) {
		t.Errorf("public shares in message doesn't match expected\nactual: %s\nexpected: %s",
			message.publicCoefficients,
			expectedPublicCoefficients,
		)
	}
}

func TestCalculateAndVerifyPublicCoefficients(t *testing.T) {
	threshold := 3
	groupSize := 5

	sharingMembers, err := initializeSharingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	sharingMember := sharingMembers[0]
	var tests = map[string]struct {
		modifyPublicCoefficientsMessages func(messages []*MemberPublicCoefficientsMessage)
		expectedError                    error
		expectedAccusedIDs               int
	}{
		"positive validation": {
			expectedError:      nil,
			expectedAccusedIDs: 0,
		},
		"negative validation - changed public key share": {
			modifyPublicCoefficientsMessages: func(messages []*MemberPublicCoefficientsMessage) {
				messages[1].publicCoefficients[1] = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			messages := make([]*MemberPublicCoefficientsMessage, groupSize)

			for i, m := range sharingMembers {
				messages[i] = m.CalculatePublicCoefficients()
			}

			filteredMessages := filterMemberPublicCoefficientsMessages(
				messages,
				sharingMember.ID,
			)

			if test.modifyPublicCoefficientsMessages != nil {
				test.modifyPublicCoefficientsMessages(filteredMessages)
			}

			accusedMessage, err := sharingMember.VerifyPublicCoefficients(filteredMessages)

			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"expected: %v\nactual: %v\n",
					test.expectedError,
					err,
				)
			}

			if len(accusedMessage.accusedIDs) != test.expectedAccusedIDs {
				t.Fatalf("expecting %d accused member's IDs but received %d",
					test.expectedAccusedIDs,
					accusedMessage.accusedIDs,
				)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	threshold := 3
	groupSize := 5

	committingMembers, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	var peerSharesMessages []*PeerSharesMessage
	var commitmentsMessages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		peerSharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		peerSharesMessages = append(peerSharesMessages, peerSharesMessage...)
		commitmentsMessages = append(commitmentsMessages, commitmentsMessage)
	}

	committingMember := committingMembers[0]

	accusedSecretSharesMessage, err := committingMember.VerifyReceivedSharesAndCommitmentsMessages(
		filterPeerSharesMessage(peerSharesMessages, committingMember.ID),
		filterMemberCommitmentsMessages(commitmentsMessages, committingMember.ID),
	)
	if err != nil {
		t.Fatalf("shares and commitments verification failed [%s]", err)
	}

	if len(accusedSecretSharesMessage.accusedIDs) > 0 {
		t.Fatalf("found accused members but was not expecting to")
	}

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {
		sharingMembers = append(sharingMembers, &SharingMember{
			CommittingMember: cm,
		})
	}

	sharingMember := sharingMembers[0]
	if len(sharingMember.receivedSharesS) != groupSize-1 {
		t.Fatalf("received shares number %d doesn't match expected number %d",
			len(sharingMember.receivedSharesS),
			groupSize-1,
		)
	}

	for _, member := range sharingMembers {
		member.CombineReceivedShares()
	}

	secondMessages := make([]*MemberPublicCoefficientsMessage, groupSize)
	for i, member := range sharingMembers {
		secondMessages[i] = member.CalculatePublicCoefficients()
	}
	accusedCoefficientsMessage, err := sharingMember.VerifyPublicCoefficients(
		filterMemberPublicCoefficientsMessages(secondMessages, sharingMember.ID),
	)
	if err != nil {
		t.Fatalf("phase8 failed [%s]", err)
	}
	if len(accusedCoefficientsMessage.accusedIDs) > 0 {
		t.Fatalf("something wrong %v", accusedCoefficientsMessage.accusedIDs)
	}
}

func TestGeneratePolynomial(t *testing.T) {
	degree := 3
	config := &DKG{P: big.NewInt(100), Q: big.NewInt(9)}

	coefficients, err := generatePolynomial(degree, config)
	if err != nil {
		t.Fatalf("unexpected error [%s]", err)
	}

	if len(coefficients) != degree+1 {
		t.Fatalf("\nexpected: coefficients number %d\nactual:   %d\n",
			degree+1,
			len(coefficients),
		)
	}
	for _, c := range coefficients {
		if c.Cmp(big.NewInt(0)) <= 0 || c.Cmp(config.Q) >= 0 {
			t.Fatalf("\nexpected: coefficient between 0 and %d\nactual:   %v\n",
				config.Q,
				c,
			)
		}
	}
}

func initializeCommittingMembersGroup(threshold, groupSize int) ([]*CommittingMember, error) {
	config, err := predefinedDKGconfig()
	if err != nil {
		return nil, fmt.Errorf("DKG Config initialization failed [%s]", err)
	}

	vss, err := pedersen.NewVSS(crand.Reader, config.P, config.Q)
	if err != nil {
		return nil, fmt.Errorf("VSS initialization failed [%s]", err)
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*CommittingMember

	for i := 1; i <= groupSize; i++ {
		id := i
		members = append(members, &CommittingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: config,
			},
			vss:             vss,
			receivedSharesS: make(map[int]*big.Int),
			receivedSharesT: make(map[int]*big.Int),
		})
		group.RegisterMemberID(id)
	}
	return members, nil
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
			CommittingMember: cm,
		})
	}

	for _, sm := range sharingMembers {
		for _, cm := range committingMembers {
			sm.receivedSharesS[cm.ID] = evaluateMemberShare(sm.ID, cm.secretCoefficients)
		}
	}

	return sharingMembers, nil
}

func filterPeerSharesMessage(
	messages []*PeerSharesMessage,
	receiverID int,
) []*PeerSharesMessage {
	var result []*PeerSharesMessage
	for _, msg := range messages {
		if msg.senderID != receiverID &&
			msg.receiverID == receiverID {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberCommitmentsMessages(
	messages []*MemberCommitmentsMessage,
	receiverID int,
) []*MemberCommitmentsMessage {
	var result []*MemberCommitmentsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberPublicCoefficientsMessages(
	messages []*MemberPublicCoefficientsMessage, receiverID int,
) []*MemberPublicCoefficientsMessage {
	var result []*MemberPublicCoefficientsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
