package dkg

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

var (
	groupSize          = 10
	dishonestThreshold = 5
)

func TestPhase3and4(t *testing.T) {
	config, err := config.PredefinedDKGconfig()
	if err != nil {
		t.Fatalf("DKG Config initialization failed [%s]", err)
	}

	vss, err := pedersen.NewVSS(config.P, config.Q)
	if err != nil {
		t.Fatalf("VSS initialization failed [%s]", err)
	}

	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: dishonestThreshold,
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
			vss: vss,
		})
		group.RegisterMemberID(id)
	}

	messages := make([]*MemberCommitmentsMessage, groupSize)
	for i, member := range members {
		messages[i], err = member.CalculateSharesAndCommitments()
		if err != nil {
			t.Fatalf("phase3 failed [%s]", err)
		}
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

func filterPeerSharesMessage(
	messages []*PeerSharesMessage,
	receiverID *big.Int,
) []*PeerSharesMessage {
	var result []*PeerSharesMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
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
