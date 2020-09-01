package group

import (
	"github.com/ipfs/go-log"
)

var logger = log.Logger("keep-message-filter")

// MessageFiltering interface defines method allowing to filter out messages
// from members that are not part of the group or were marked as IA or DQ.
type MessageFiltering interface {

	// IsSenderAccepted returns true if the message from the given sender should be
	// accepted for further processing. Otherwise, function returns false.
	// Message from the given sender is allowed only if that member is a properly
	// operating group member - it was not DQ or IA so far.
	IsSenderAccepted(senderID MemberIndex) bool

	// IsSenderValid returns true if the message from the given sender should be
	// accepted for further processing. Otherwise, function returns false.
	// IsSenderValid checks if sender of the provided ProtocolMessage is in the
	// group and uses appropriate group member index.
	IsSenderValid(senderID MemberIndex, senderPublicKey []byte) bool
}

// ProtocolMessage is a common interface for all messages of GJKR DKG protocol.
type ProtocolMessage interface {
	// SenderID returns protocol-level identifier of the message sender.
	SenderID() MemberIndex
}

// InactiveMemberFilter is a proxy facilitates filtering out inactive members
// in the given phase and registering their final list in DKG Group.
type InactiveMemberFilter struct {
	selfMemberID MemberIndex
	group        *Group

	phaseActiveMembers []MemberIndex
}

// NewInactiveMemberFilter creates a new instance of InactiveMemberFilter.
// It accepts member index of the current member (the one which will be
// filtering out other group members for inactivity) and the reference to Group
// to which all those members belong.
func NewInactiveMemberFilter(
	selfMemberIndex MemberIndex,
	group *Group,
) *InactiveMemberFilter {
	return &InactiveMemberFilter{
		selfMemberID:       selfMemberIndex,
		group:              group,
		phaseActiveMembers: make([]MemberIndex, 0),
	}
}

// MarkMemberAsActive marks member with the given index as active in the given
// phase.
func (mf *InactiveMemberFilter) MarkMemberAsActive(memberID MemberIndex) {
	mf.phaseActiveMembers = append(mf.phaseActiveMembers, memberID)
}

// FlushInactiveMembers takes all members who were not previously marked as
// active and flushes them to DKG group as inactive members.
func (mf *InactiveMemberFilter) FlushInactiveMembers() {
	isActive := func(id MemberIndex) bool {
		if id == mf.selfMemberID {
			return true
		}

		for _, activeMemberID := range mf.phaseActiveMembers {
			if activeMemberID == id {
				return true
			}
		}

		return false
	}

	for _, operatingMemberID := range mf.group.OperatingMemberIDs() {
		if !isActive(operatingMemberID) {
			logger.Warningf(
				"[member:%v] marking member [%v] as inactive",
				mf.selfMemberID,
				operatingMemberID,
			)
			mf.group.MarkMemberAsInactive(operatingMemberID)
		}
	}
}

// IsMessageFromSelf is an auxiliary function determining whether the given
// ProtocolMessage is from the current member itself.
func IsMessageFromSelf(memberIndex MemberIndex, message ProtocolMessage) bool {
	if message.SenderID() == memberIndex {
		return true
	}

	return false
}

// IsSenderValid checks if sender of the provided ProtocolMessage is in the
// group and uses appropriate group member index.
func IsSenderValid(
	filter MessageFiltering,
	message ProtocolMessage,
	senderPublicKey []byte,
) bool {
	return filter.IsSenderValid(message.SenderID(), senderPublicKey)
}

// IsSenderAccepted determines if sender of the given ProtocoLMessage is
// accepted by group (not marked as inactive or disqualified).
func IsSenderAccepted(filter MessageFiltering, message ProtocolMessage) bool {
	return filter.IsSenderAccepted(message.SenderID())
}
