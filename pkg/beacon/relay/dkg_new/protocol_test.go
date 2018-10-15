package dkg

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

func TestCalculateSharesAndCommitments(t *testing.T) {
	threshold := 4
	groupSize := 10

	members, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	member := members[0]
	peerSharesMessages, commitmentsMessage, err := member.CalculateSharesAndCommitments()
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

func TestPhase3and4(t *testing.T) {
	threshold := 1
	groupSize := 2

	members, err := initializeCommittingMembersGroup(threshold, groupSize)
	if err != nil {
		t.Fatalf("group initialization failed [%s]", err)
	}

	var peerSharesMessages []*PeerSharesMessage
	var messages []*MemberCommitmentsMessage
	for _, member := range members {
		peerSharesMessage, commitmentsMessage, err := member.CalculateSharesAndCommitments()
		if err != nil {
			t.Fatalf("shares and commitments calculation failed [%s]", err)
		}
		peerSharesMessages = append(peerSharesMessages, peerSharesMessage...)
		messages = append(messages, commitmentsMessage)
	}

	if len(messages) != groupSize {
		t.Fatalf("generated messages number %d doesn't match expected number %d", len(messages), groupSize)
	}

	currentMember := members[0]

	accusedMessage, err := currentMember.VerifySharesAndCommitments(
		filterPeerSharesMessage(peerSharesMessages, currentMember.ID),
		filterMemberCommitmentsMessages(messages, currentMember.ID),
	)
	if err != nil {
		t.Fatalf("shares and commitments verification failed [%s]", err)
	}

	if len(accusedMessage.accusedIDs) > 0 {
		t.Fatalf("found accused members but was not expecting to")
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
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
