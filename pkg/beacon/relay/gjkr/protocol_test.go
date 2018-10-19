package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
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

	members, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}
	currentMember := members[0]

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
				messages[1].shareT = big.NewInt(13)
				return messages
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) []*MemberCommitmentsMessage {
				messages[1].commitments[0] = big.NewInt(13)
				return messages
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
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
	var messages []*MemberCommitmentsMessage
	for _, member := range committingMembers {
		peerSharesMessage, commitmentsMessage, err := member.CalculateMembersSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		peerSharesMessages = append(peerSharesMessages, peerSharesMessage...)
		messages = append(messages, commitmentsMessage)
	}

	committingMember := committingMembers[0]

	accusedMessage, err := committingMember.VerifyReceivedSharesAndCommitmentsMessages(
		filterPeerSharesMessage(peerSharesMessages, committingMember.ID),
		filterMemberCommitmentsMessages(messages, committingMember.ID),
	)
	if err != nil {
		t.Fatalf("shares and commitments verification failed [%s]", err)
	}

	if len(accusedMessage.accusedIDs) > 0 {
		t.Fatalf("found accused members but was not expecting to")
	}
}

func TestGeneratePolynomial(t *testing.T) {
	degree := 3
	config := &config.DKG{P: big.NewInt(100), Q: big.NewInt(9)}

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
	config, err := config.PredefinedDKGconfig()
	if err != nil {
		return nil, fmt.Errorf("DKG Config initialization failed [%s]", err)
	}

	vss, err := pedersen.NewVSS(config.P, config.Q)
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
