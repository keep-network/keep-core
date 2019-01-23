package dkg

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

type keyGenerationState interface {
	groupMember() thresholdgroup.BaseMember
	// activeBlocks is the period during which this state is active, in blocks.
	activeBlocks() int

	initiate() error
	receive(msg net.Message) error
	nextState() (keyGenerationState, error)
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *thresholdgroup.LocalMember
}

func (is *initializationState) groupMember() thresholdgroup.BaseMember {
	return is.member
}

func (is *initializationState) activeBlocks() int { return 1 }

func (is *initializationState) initiate() error {
	return nil
}

func (is *initializationState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for initialization state: [%#v]", msg)
}

func (is *initializationState) nextState() (keyGenerationState, error) {
	return &joinState{is.channel, is.member}, nil
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol. Join
// messages from other members are valid in this state, and when the member is
// ready and activePeriod has elapsed, it proceeds to commitmentState.
type joinState struct {
	channel net.BroadcastChannel
	member  *thresholdgroup.LocalMember
}

func (js *joinState) groupMember() thresholdgroup.BaseMember { return js.member }
func (js *joinState) activeBlocks() int                      { return 1 }

func (js *joinState) initiate() error {
	return js.channel.Send(&JoinMessage{js.member.BlsID.Raw()})
}

func (js *joinState) receive(msg net.Message) error {
	switch joinMsg := msg.Payload().(type) {
	case *JoinMessage:
		js.member.RegisterMemberID(joinMsg.id)

		return nil
	}

	return fmt.Errorf("unexpected message for join state: [%#v]", msg)
}

func (js *joinState) nextState() (keyGenerationState, error) {
	if js.member.MemberListComplete() {
		sharingMember := js.member.InitializeSharing()
		return &commitmentState{js.channel, sharingMember}, nil
	}

	return js, nil
}

type commitmentState struct {
	channel net.BroadcastChannel
	member  *thresholdgroup.SharingMember
}

func (cs *commitmentState) groupMember() thresholdgroup.BaseMember { return cs.member }
func (cs *commitmentState) activeBlocks() int                      { return 2 }

func (cs *commitmentState) initiate() error {
	return cs.channel.Send(&MemberCommitmentsMessage{
		cs.member.BlsID.Raw(),
		cs.member.Commitments(),
	})
}

func (cs *commitmentState) receive(msg net.Message) error {
	switch commitmentMsg := msg.Payload().(type) {
	case *MemberCommitmentsMessage:
		if senderID, ok := msg.ProtocolSenderID().(*thresholdgroup.BlsID); ok {
			if senderID.Raw().IsEqual(cs.member.BlsID.Raw()) {
				fmt.Printf("sender [%v]", cs.member.BlsID)
				return nil
			}

			cs.member.AddCommitmentsFromID(
				commitmentMsg.id,
				commitmentMsg.Commitments,
			)

			return nil
		}

		return fmt.Errorf(
			"unknown protocol sender id type [%T]",
			msg.ProtocolSenderID(),
		)
	}

	return fmt.Errorf("unexpected message for committing state: [%#v]", msg)
}

func (cs *commitmentState) nextState() (keyGenerationState, error) {
	if cs.member.CommitmentsComplete() {
		return &sharingState{cs.channel, cs.member}, nil
	}

	return cs, nil
}

type sharingState struct {
	channel net.BroadcastChannel
	member  *thresholdgroup.SharingMember
}

func (ss *sharingState) groupMember() thresholdgroup.BaseMember { return ss.member }
func (ss *sharingState) activeBlocks() int                      { return 3 }

func (ss *sharingState) initiate() error {
	for _, receiverID := range ss.member.OtherMemberIDs() {
		share := ss.member.SecretShareForID(receiverID)

		err := ss.channel.SendTo(
			net.ProtocolIdentifier(receiverID),
			&MemberShareMessage{ss.member.BlsID.Raw(), receiverID, share})

		if err != nil {
			return err
		}
	}

	return nil
}

func (ss *sharingState) receive(msg net.Message) error {
	switch shareMsg := msg.Payload().(type) {
	case *MemberShareMessage:
		if shareMsg.receiverID.IsEqual(ss.member.BlsID.Raw()) {
			ss.member.AddShareFromID(shareMsg.id, shareMsg.Share)
		}
		return nil
	}

	return fmt.Errorf("unexpected message for sharing state: [%#v]", msg)
}

func (ss *sharingState) nextState() (keyGenerationState, error) {
	if ss.member.SharesComplete() {
		justifyingMember := ss.member.InitializeJustification()
		return &accusingState{
			ss.channel,
			justifyingMember,
			make(map[bls.ID]struct{}),
			len(justifyingMember.OtherMemberIDs()),
		}, nil
	}

	return ss, nil
}

type accusingState struct {
	channel                 net.BroadcastChannel
	member                  *thresholdgroup.JustifyingMember
	seenAccusations         map[bls.ID]struct{}
	expectedAccusationCount int
}

func (as *accusingState) groupMember() thresholdgroup.BaseMember { return as.member }
func (as *accusingState) activeBlocks() int                      { return 1 }

func (as *accusingState) initiate() error {
	return as.channel.Send(&AccusationsMessage{
		as.member.BlsID.Raw(),
		as.member.AccusedIDs(),
	})
}

func (as *accusingState) receive(msg net.Message) error {
	switch accusationMsg := msg.Payload().(type) {
	case *AccusationsMessage:
		if senderID, ok := msg.ProtocolSenderID().(*thresholdgroup.BlsID); ok {
			if senderID.Raw().IsEqual(as.member.BlsID.Raw()) {
				return nil
			}

			for _, accusedID := range accusationMsg.accusedIDs {
				as.member.AddAccusationFromID(
					accusationMsg.id,
					accusedID,
				)
			}

			as.seenAccusations[*accusationMsg.id] = struct{}{}

			return nil
		}

		return fmt.Errorf(
			"unknown protocol sender id [%v]",
			msg.ProtocolSenderID(),
		)
	}

	return fmt.Errorf("unexpected message for justifying state: [%#v]", msg)
}

func (as *accusingState) nextState() (keyGenerationState, error) {
	if len(as.seenAccusations) == as.expectedAccusationCount {
		return &justifyingState{
			as.channel,
			as.member,
			make(map[bls.ID]struct{}),
			as.expectedAccusationCount,
		}, nil
	}

	return as, nil
}

type justifyingState struct {
	channel                    net.BroadcastChannel
	member                     *thresholdgroup.JustifyingMember
	seenJustifications         map[bls.ID]struct{}
	expectedJustificationCount int
}

func (js *justifyingState) groupMember() thresholdgroup.BaseMember { return js.member }
func (js *justifyingState) activeBlocks() int                      { return 1 }

func (js *justifyingState) initiate() error {
	return js.channel.Send(
		&JustificationsMessage{js.member.BlsID.Raw(), js.member.Justifications()})
}

func (js *justifyingState) receive(msg net.Message) error {
	switch justificationsMsg := msg.Payload().(type) {
	case *JustificationsMessage:
		if senderID, ok := msg.ProtocolSenderID().(*thresholdgroup.BlsID); ok {
			if senderID.Raw().IsEqual(js.member.BlsID.Raw()) {
				return nil
			}

			for accuserID, justification := range justificationsMsg.justifications {
				js.member.RecordJustificationFromID(
					justificationsMsg.id,
					&accuserID,
					justification,
				)
			}

			js.seenJustifications[*justificationsMsg.id] = struct{}{}

			return nil
		}

		return fmt.Errorf(
			"unknown protocol sender id [%v]",
			msg.ProtocolSenderID(),
		)
	}

	return fmt.Errorf("unexpected message for justifying state: [%#v]", msg)
}

func (js *justifyingState) nextState() (keyGenerationState, error) {
	if len(js.seenJustifications) == js.expectedJustificationCount {
		member, err := js.member.FinalizeMember()
		if err != nil {
			return nil, err
		}

		return &keyedState{member}, nil
	}

	return js, nil
}

type keyedState struct {
	member *thresholdgroup.Member
}

func (ks *keyedState) groupMember() thresholdgroup.BaseMember { return ks.member }
func (ks *keyedState) activeBlocks() int                      { return 0 }

func (ks *keyedState) initiate() error {
	return nil
}

func (ks *keyedState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for keyed state: [%#v]", msg)
}

func (ks *keyedState) nextState() (keyGenerationState, error) {
	return nil, nil
}
