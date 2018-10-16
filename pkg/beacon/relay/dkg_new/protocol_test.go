package dkg

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"

	"github.com/golang/go/src/crypto/rand"
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

	if len(member.secretShares) != (threshold + 1) {
		t.Fatalf("generated coefficients A number %d doesn't match expected number %d",
			len(member.secretShares),
			threshold+1,
		)
	}
	if len(peerSharesMessages) != (groupSize - 1) {
		t.Fatalf("peer shares messages number %d doesn't match expected %d",
			len(peerSharesMessages),
			groupSize-1,
		)
	}

	if len(commitmentsMessage.commitments) != (threshold + 1) {
		t.Fatalf("calculated commitments number %d doesn't match expected number %d",
			len(member.secretShares),
			threshold+1,
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

	currentMember := members[0]

	var tests = map[string]struct {
		modifyPeerShareMessages   func(messages []*PeerSharesMessage)
		modifyCommitmentsMessages func(messages []*MemberCommitmentsMessage)
		expectedError             error
		expectedAccusedIDs        int
	}{
		"positive validation": {
			expectedError:      nil,
			expectedAccusedIDs: 0,
		},
		"negative validation - changed random share": {
			modifyPeerShareMessages: func(messages []*PeerSharesMessage) {
				messages[1].randomShare = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
		"negative validation - changed commitment": {
			modifyCommitmentsMessages: func(messages []*MemberCommitmentsMessage) {
				messages[1].commitments[0] = big.NewInt(13)
			},
			expectedError:      nil,
			expectedAccusedIDs: 1,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			filteredPeerSharesMessages := filterPeerSharesMessage(peerSharesMessages, currentMember.ID)
			filteredMemberCommitmentsMessages := filterMemberCommitmentsMessages(commitmentsMessages, currentMember.ID)

			if test.modifyPeerShareMessages != nil {
				test.modifyPeerShareMessages(filteredPeerSharesMessages)
			}
			if test.modifyCommitmentsMessages != nil {
				test.modifyCommitmentsMessages(filteredMemberCommitmentsMessages)
			}

			accusedMessage, err := currentMember.VerifyReceivedSharesAndCommitmentsMessages(
				filteredPeerSharesMessages,
				filteredMemberCommitmentsMessages,
			)
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

func TestCombineReceivedShares(t *testing.T) {
	receivedSecretShares := make(map[*big.Int]*big.Int)
	receivedRandomShares := make(map[*big.Int]*big.Int)
	for i := 0; i <= 5; i++ {
		receivedSecretShares[big.NewInt(int64(100+i))] = big.NewInt(int64(i))
		receivedRandomShares[big.NewInt(int64(100+i))] = big.NewInt(int64(10 + i))
	}

	expectedPrivateKeyShare := big.NewInt(15)
	expectedPrivateRandomShare := big.NewInt(75)

	config, err := config.PredefinedDKGconfig()
	if err != nil {
		t.Fatalf("DKG Config initialization failed [%s]", err)
	}
	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			receivedSecretShares: receivedSecretShares,
			receivedRandomShares: receivedRandomShares,
		},
	}

	member.CombineReceivedShares()

	if member.privateKeyShare.Cmp(expectedPrivateKeyShare) != 0 {
		t.Errorf("combined secret shares %s doesn't match expected %s",
			member.privateKeyShare,
			expectedPrivateKeyShare)
	}
	if member.privateRandomShare.Cmp(expectedPrivateRandomShare) != 0 {
		t.Errorf("combined random shares %s doesn't match expected %s",
			member.privateRandomShare,
			expectedPrivateRandomShare)
	}
}

func TestCalculatePublicKeyShares(t *testing.T) {
	secretShares := []*big.Int{
		big.NewInt(3),
		big.NewInt(1),
		big.NewInt(5),
	}
	expectedPublicShares := []*big.Int{
		big.NewInt(46656),
		big.NewInt(36),
		big.NewInt(60466176),
	}

	mockRandomReader := testutils.NewMockRandReader(big.NewInt(6))

	config, err := config.PredefinedDKGconfig()
	if err != nil {
		t.Fatalf("DKG Config initialization failed [%s]", err)
	}
	vss, err := pedersen.NewVSS(mockRandomReader, config.P, config.Q)
	if err != nil {
		t.Fatalf("VSS initialization failed [%s]", err)
	}

	member := &SharingMember{
		CommittingMember: &CommittingMember{
			memberCore: &memberCore{
				protocolConfig: config,
			},
			vss:          vss,
			secretShares: secretShares,
		},
	}

	memberPublicKeySharesMessage := member.CalculatePublicKeyShares()

	if !reflect.DeepEqual(member.publicShares, expectedPublicShares) {
		t.Errorf("public shares for member doesn't match expected\nactual: %s\nexpected: %s",
			memberPublicKeySharesMessage.publicKeyShares,
			expectedPublicShares,
		)
	}

	if !reflect.DeepEqual(memberPublicKeySharesMessage.publicKeyShares, expectedPublicShares) {
		t.Errorf("public shares in message doesn't match expected\nactual: %s\nexpected: %s",
			memberPublicKeySharesMessage.publicKeyShares,
			expectedPublicShares,
		)
	}
}

func TestRoundTrip(t *testing.T) {
	threshold := 5
	groupSize := 10

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

	var sharingMembers []*SharingMember
	for _, cm := range committingMembers {

		sharingMembers = append(sharingMembers, &SharingMember{
			CommittingMember: cm,
		})
	}

	sharingMember := sharingMembers[0]
	if len(sharingMember.receivedSecretShares) != groupSize {
		t.Fatalf("received shares number %d doesn't match expected number %d", len(sharingMember.receivedSecretShares), groupSize-1)
	}

	for _, member := range sharingMembers {
		member.CombineReceivedShares()
	}

	secondMessages := make([]*MemberPublicKeySharesMessage, groupSize)
	for i, member := range sharingMembers {
		secondMessages[i] = member.CalculatePublicKeyShares()
	}
	secondAccusedMessage, err := sharingMember.VerifyPublicKeyShares(secondMessages)
	if err != nil {
		t.Fatalf("phase8 failed [%s]", err)
	}
	if len(secondAccusedMessage.accusedIDs) > 0 {
		t.Fatalf("something wrong %v", secondAccusedMessage.accusedIDs)
	}
}

func initializeCommittingMembersGroup(threshold, groupSize int) ([]*CommittingMember, error) {
	config, err := config.PredefinedDKGconfig()
	if err != nil {
		return nil, fmt.Errorf("DKG Config initialization failed [%s]", err)
	}

	vss, err := pedersen.NewVSS(rand.Reader, config.P, config.Q)
	if err != nil {
		return nil, fmt.Errorf("VSS initialization failed [%s]", err)
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*CommittingMember

	for i := 1; i <= groupSize; i++ {
		id := big.NewInt(int64(i))
		members = append(members, &CommittingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: config,
			},
			vss:                  vss,
			receivedSecretShares: make(map[*big.Int]*big.Int),
			receivedRandomShares: make(map[*big.Int]*big.Int),
		})
		group.RegisterMemberID(id)
	}
	return members, nil
}

func filterPeerSharesMessage(
	messages []*PeerSharesMessage,
	receiverID *big.Int,
) []*PeerSharesMessage {
	var result []*PeerSharesMessage
	for _, msg := range messages {
		if msg.senderID.Cmp(receiverID) != 0 &&
			msg.receiverID.Cmp(receiverID) == 0 {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberCommitmentsMessages(
	messages []*MemberCommitmentsMessage,
	receiverID *big.Int,
) []*MemberCommitmentsMessage {
	var result []*MemberCommitmentsMessage
	for _, msg := range messages {
		if msg.senderID.Cmp(receiverID) != 0 {
			result = append(result, msg)
		}
	}
	return result
}

func filterMemberPublicKeySharesMessageMessages(
	messages []*MemberPublicKeySharesMessage, receiverID *big.Int,
) []*MemberPublicKeySharesMessage {
	var result []*MemberPublicKeySharesMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
