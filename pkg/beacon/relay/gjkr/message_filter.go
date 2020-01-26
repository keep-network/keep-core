package gjkr

import (
	"crypto/ecdsa"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (em *SymmetricKeyGeneratingMember) MarkInactiveMembers(
	ephemeralPubKeyMessages []*EphemeralPublicKeyMessage,
) {
	filter := em.messageFilter()
	for _, message := range ephemeralPubKeyMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *CommitmentsVerifyingMember) MarkInactiveMembers(
	sharesMessages []*PeerSharesMessage,
	commitmentsMessages []*MemberCommitmentsMessage,
) {
	filter := cvm.messageFilter()
	for _, sharesMessage := range sharesMessages {
		for _, commitmentsMessage := range commitmentsMessages {
			if sharesMessage.senderID == commitmentsMessage.senderID {
				filter.MarkMemberAsActive(sharesMessage.senderID)
				break
			}
		}
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *SharesJustifyingMember) MarkInactiveMembers(
	sharesAccusationsMessages []*SecretSharesAccusationsMessage,
) {
	filter := cvm.messageFilter()
	for _, message := range sharesAccusationsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (sm *SharingMember) MarkInactiveMembers(
	keySharePointsMessages []*MemberPublicKeySharePointsMessage,
) {
	filter := sm.messageFilter()
	for _, message := range keySharePointsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (cvm *PointsJustifyingMember) MarkInactiveMembers(
	pointsAccusationsMessages []*PointsAccusationsMessage,
) {
	filter := cvm.messageFilter()
	for _, message := range pointsAccusationsMessages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

// MarkInactiveMembers takes all messages from the previous DKG protocol
// execution phase and marks all member who did not send a message as IA.
func (rm *ReconstructingMember) MarkInactiveMembers(
	messages []*MisbehavedEphemeralKeysMessage,
) {
	filter := rm.messageFilter()
	for _, message := range messages {
		filter.MarkMemberAsActive(message.senderID)
	}

	filter.FlushInactiveMembers()
}

func (mc *memberCore) messageFilter() *group.InactiveMemberFilter {
	return group.NewInactiveMemberFilter(mc.ID, mc.group)
}

func (mc *memberCore) IsSenderAccepted(senderID group.MemberIndex) bool {
	return mc.group.IsOperating(senderID)
}

// IsSenderValid checks if sender of the provided ProtocolMessage is in the
// group and uses appropriate group member index.
func (mc *memberCore) IsSenderValid(
	senderID group.MemberIndex,
	senderPublicKey *ecdsa.PublicKey,
) bool {
	return mc.membershipValidator.IsSelectedAtIndex(
		// At GJKR protocol, we index members from 1 but when they are selected
		// the the group, they are indexed from 0.
		int(senderID)-1,
		senderPublicKey,
	)
}
