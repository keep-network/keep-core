// Package announcer contains an implementation of a generic protocol announcer
// that can be used to determine live participants of an interactive protocol
// before executing the given protocol session.
package announcer

import (
	"context"
	"fmt"
	"sort"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/announcer/gen/pb"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"google.golang.org/protobuf/proto"
)

// announcementMessage represents a message that is used to announce
// member's participation in the given session of the protocol.
type announcementMessage struct {
	senderID   group.MemberIndex
	protocolID string
	sessionID  string
}

func (am *announcementMessage) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.AnnouncementMessage{
		SenderID:   uint32(am.senderID),
		ProtocolID: am.protocolID,
		SessionID:  am.sessionID,
	})
}

func (am *announcementMessage) Unmarshal(bytes []byte) error {
	pbMessage := pb.AnnouncementMessage{}
	if err := proto.Unmarshal(bytes, &pbMessage); err != nil {
		return fmt.Errorf(
			"failed to unmarshal AnnouncementMessage: [%v]",
			err,
		)
	}

	if senderID := pbMessage.SenderID; senderID > group.MaxMemberIndex {
		return fmt.Errorf("invalid member index value: [%v]", senderID)
	} else {
		am.senderID = group.MemberIndex(senderID)
	}

	am.protocolID = pbMessage.ProtocolID
	am.sessionID = pbMessage.SessionID

	return nil
}

func (am *announcementMessage) Type() string {
	return "protocol_announcer/announcement_message"
}

// Announcer is an implementation of the protocol announcer that performs the
// readiness announcement over the provided broadcast channel.
type Announcer struct {
	protocolID          string
	groupSize           int
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
}

// New creates a new instance of the Announcer. It expects a unique protocol
// identifier, the size of the group performing the protocol, a broadcast
// channel configured to mediate between group members, and a membership
// validator configured to validate the group membership of announcements
// senders.
func New(
	protocolID string,
	groupSize int,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) *Announcer {
	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &announcementMessage{}
	})

	return &Announcer{
		protocolID:          protocolID,
		groupSize:           groupSize,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
	}
}

// Announce sends the member's readiness announcement for the given protocol
// session and listens for announcements from other group members. It returns a
// list of unique members indexes that are ready for the given attempt,
// including the executing member's index. The list is sorted in ascending order.
// This function blocks until the ctx passed as argument is done.
func (a *Announcer) Announce(
	ctx context.Context,
	memberIndex group.MemberIndex,
	sessionID string,
) (
	[]group.MemberIndex,
	error,
) {
	messagesChan := make(chan net.Message, a.groupSize)
	a.broadcastChannel.Recv(ctx, func(message net.Message) {
		messagesChan <- message
	})

	err := a.broadcastChannel.Send(ctx, &announcementMessage{
		senderID:   memberIndex,
		protocolID: a.protocolID,
		sessionID:  sessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot send announcement message: [%w]", err)
	}

	readyMembersIndexesSet := make(map[group.MemberIndex]bool)
	// Mark itself as ready.
	readyMembersIndexesSet[memberIndex] = true

loop:
	for {
		select {
		case netMessage := <-messagesChan:
			announcement, ok := netMessage.Payload().(*announcementMessage)
			if !ok {
				continue
			}

			if announcement.senderID == memberIndex {
				continue
			}

			if !a.membershipValidator.IsValidMembership(
				announcement.senderID,
				netMessage.SenderPublicKey(),
			) {
				continue
			}

			if announcement.protocolID != a.protocolID {
				continue
			}

			if announcement.sessionID != sessionID {
				continue
			}

			readyMembersIndexesSet[announcement.senderID] = true
		case <-ctx.Done():
			break loop
		}
	}

	readyMembersIndexes := make([]group.MemberIndex, 0)
	for memberIndex := range readyMembersIndexesSet {
		readyMembersIndexes = append(readyMembersIndexes, memberIndex)
	}

	sort.Slice(readyMembersIndexes, func(i, j int) bool {
		return readyMembersIndexes[i] < readyMembersIndexes[j]
	})

	return readyMembersIndexes, nil
}
