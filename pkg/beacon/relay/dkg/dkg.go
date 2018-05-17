package dkg

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &JoinMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &MemberCommitmentsMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &MemberShareMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &AccusationsMessage{} })
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &JustificationsMessage{} })
}

// ExecuteDKG runs the full distributed key generation lifecycle, given a
// broadcast channel to mediate it and a group size and threshold. It returns a
// threshold group member who is participating in the group if the generation
// was successful, and an error representing what went wrong if not.
func ExecuteDKG(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	groupSize int,
	threshold int,
) (*thresholdgroup.Member, error) {
	// FIXME Probably pass in a way to ask for a receiver's public key?
	// FIXME Need a way to time out in a given stage, especially the waiting
	//       ones.

	// Generate a nonzero memberID; loop until rand.Int31 returns something
	// other than 0, hopefully no more than once :)
	memberID := "0"
	for memberID = fmt.Sprintf("%v", rand.Int31()); memberID == "0"; {
	}
	fmt.Printf("[member:%v] Initializing member.\n", memberID)
	localMember := thresholdgroup.NewMember(memberID, threshold, groupSize)

	recvChan := make(chan net.Message)
	channel.Recv(func(msg net.Message) error {
		recvChan <- msg
		return nil
	})

	fmt.Printf("[member:%v] Waiting for join timeout...\n", memberID)
	blockCounter.WaitForBlocks(15)

	fmt.Printf("[member:%v] Broadcasting join.\n", memberID)
	err := channel.Send(&JoinMessage{&localMember.BlsID})
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast join: [%v]", err)
	}

	// Wait for all members.
	waiter := blockCounter.BlockWaiter(10)
	fmt.Printf("[member:%v] Waiting for other members...\n", memberID)
	sharingMember, err := waitForMemberIDs(&localMember, channel, recvChan)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all member ids: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for member join timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(15)
	fmt.Printf("[member:%v] Initiating commitment broadcast phase.\n", memberID)

	fmt.Printf("[member:%v] Broadcasting public commitment.\n", memberID)
	err = sendCommitments(channel, sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast commitments: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other commitments...\n", memberID)
	err = waitForCommitments(recvChan, sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all commitments: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for commitment timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(20)
	fmt.Printf("[member:%v] Sending private shares.\n", memberID)
	err = sendShares(channel, sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to send all private shares: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other shares...\n", memberID)
	justifyingMember, err := waitForShares(recvChan, sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all private shares: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for share exchange timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(15)
	fmt.Printf("[member:%v] Broadcasting accusations.\n", memberID)
	err = sendAccusations(channel, justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast accusations: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other accusations...\n", memberID)
	err = waitForAccusations(recvChan, justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all accusations: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for accusation timeout...\n", memberID)
	<-waiter

	fmt.Printf("[member:%v] Broadcasting justifications.\n", memberID)
	err = sendJustifications(channel, justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast justifications: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other justifications...\n", memberID)
	member, err := waitForJustifications(recvChan, justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all justifications: [%v]", err)
	}

	fmt.Printf("[member:%v] Finalized member.\n", memberID)
	return member, err
}

func registerMemberID(
	member *thresholdgroup.LocalMember,
	msg net.Message,
	broadcastChannel net.BroadcastChannel,
) error {
	switch joinMsg := msg.Payload().(type) {
	case *JoinMessage:
		err := broadcastChannel.RegisterIdentifier(msg.TransportSenderID(), joinMsg.id)
		if err != nil {
			return err
		}

		member.RegisterMemberID(joinMsg.id)
	}

	return nil
}

func waitForMemberIDs(
	member *thresholdgroup.LocalMember,
	broadcastChannel net.BroadcastChannel,
	recvChan <-chan net.Message,
) (*thresholdgroup.SharingMember, error) {
	for msg := range recvChan {
		switch msg.Payload().(type) {
		case *JoinMessage:
			err := registerMemberID(member, msg, broadcastChannel)
			if err != nil {
				return nil, err
			}

			if member.ReadyForSharing() {
				return member.InitializeSharing(), nil
			}
		}
	}

	return nil, fmt.Errorf("did not complete DKG member bootstrap")
}

func sendCommitments(channel net.BroadcastChannel, member *thresholdgroup.SharingMember) error {
	return channel.Send(&MemberCommitmentsMessage{&member.BlsID, member.Commitments()})
}

func waitForCommitments(
	recvChan <-chan net.Message,
	sharingMember *thresholdgroup.SharingMember,
) error {
done:
	for msg := range recvChan {
		switch commitmentMsg := msg.Payload().(type) {
		case *MemberCommitmentsMessage:
			if senderID, ok := msg.ProtocolSenderID().(*bls.ID); ok {
				if senderID.IsEqual(&sharingMember.BlsID) {
					continue
				}

				sharingMember.AddCommitmentsFromID(
					*commitmentMsg.id,
					commitmentMsg.Commitments)

				if sharingMember.CommitmentsComplete() {
					break done
				}
			} else {
				return fmt.Errorf(
					"unknown protocol sender id type [%v] for network id [%v]",
					reflect.TypeOf(msg.ProtocolSenderID()),
					msg.TransportSenderID())
			}
		}
	}

	return nil
}

func sendShares(channel net.BroadcastChannel, member *thresholdgroup.SharingMember) error {
	fmt.Printf("[member:%v] Despatching shares!\n", member.ID)
	for _, receiverID := range member.OtherMemberIDs() {
		share := member.SecretShareForID(receiverID)
		channel.SendTo(
			net.ProtocolIdentifier(receiverID),
			&MemberShareMessage{&member.BlsID, receiverID, &share})
	}
	fmt.Printf("[member:%v] Shares despatched!\n", member.ID)

	return nil
}

func waitForShares(
	recvChan <-chan net.Message,
	sharingMember *thresholdgroup.SharingMember,
) (*thresholdgroup.JustifyingMember, error) {
	for msg := range recvChan {
		switch shareMsg := msg.Payload().(type) {
		case *MemberShareMessage:
			if shareMsg.receiverID.IsEqual(&sharingMember.BlsID) {
				sharingMember.AddShareFromID(*shareMsg.id, *shareMsg.Share)

				if sharingMember.SharesComplete() {
					return sharingMember.InitializeJustification(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("did not complete share exchange")
}

func sendAccusations(channel net.BroadcastChannel, member *thresholdgroup.JustifyingMember) error {
	return channel.Send(&AccusationsMessage{&member.BlsID, member.AccusedIDs()})
}

func waitForAccusations(
	recvChan <-chan net.Message,
	justifyingMember *thresholdgroup.JustifyingMember,
) error {
	expectedAccusationCount := len(justifyingMember.OtherMemberIDs())
	seenAccusations := make(map[bls.ID]struct{}, expectedAccusationCount)
done:
	for msg := range recvChan {
		switch accusationMsg := msg.Payload().(type) {
		case *AccusationsMessage:
			if senderID, ok := msg.ProtocolSenderID().(*bls.ID); ok {
				if senderID.IsEqual(&justifyingMember.BlsID) {
					continue
				}

				for _, accusedID := range accusationMsg.accusedIDs {
					justifyingMember.AddAccusationFromID(
						accusationMsg.id,
						accusedID)
				}

				seenAccusations[*accusationMsg.id] = struct{}{}
				if len(seenAccusations) == expectedAccusationCount {
					break done
				}
			} else {
				return fmt.Errorf(
					"unknown protocol sender id for network id [%v]",
					msg.TransportSenderID())
			}
		}
	}

	return nil
}

func sendJustifications(channel net.BroadcastChannel, justifyingMember *thresholdgroup.JustifyingMember) error {
	return channel.Send(
		&JustificationsMessage{
			&justifyingMember.BlsID,
			justifyingMember.Justifications()})
}

func waitForJustifications(
	recvChan <-chan net.Message,
	justifyingMember *thresholdgroup.JustifyingMember,
) (*thresholdgroup.Member, error) {
	memberIDs := justifyingMember.OtherMemberIDs()
	seenJustifications := make(map[bls.ID]bool, len(memberIDs))
	for msg := range recvChan {
		switch justificationsMsg := msg.Payload().(type) {
		case *JustificationsMessage:
			if senderID, ok := msg.ProtocolSenderID().(*bls.ID); ok {
				if senderID.IsEqual(&justifyingMember.BlsID) {
					continue
				}

				for accuserID, justification := range justificationsMsg.justifications {
					justifyingMember.RecordJustificationFromID(
						*justificationsMsg.id,
						accuserID,
						justification)
				}

				seenJustifications[*justificationsMsg.id] = true
				if len(seenJustifications) == len(memberIDs) {
					return justifyingMember.FinalizeMember(), nil
				}
			} else {
				return nil, fmt.Errorf(
					"unknown protocol sender id for network id [%v]",
					msg.TransportSenderID())
			}
		}
	}

	return nil, fmt.Errorf("did not complete justification phase")
}
