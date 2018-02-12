package dkg

import (
	"fmt"

	"github.com/keep-network/go-dfinity-crypto/bls"
	"github.com/keep-network/go-dfinity-crypto/rand"
	"github.com/keep-network/keep-core/go/beacon/broadcast"
	"github.com/keep-network/keep-core/go/beacon/chain"
	"github.com/keep-network/keep-core/go/thresholdgroup"
)

// JoinMessage is an empty message payload indicating a member has joined. The
// sender is the joining member. It is expected to be broadcast.
type JoinMessage struct{}

// MemberCommitmentsMessage is a message payload that carries the sender's
// public commitments during distributed key generation. It is expected to be
// broadcast.
type MemberCommitmentsMessage struct {
	Commitments []bls.PublicKey
}

// MemberShareMessage is a message payload that carries the sender's private
// share for the recipient during distributed key generation. It is expected to
// be communicated in encrypted fashion to the recipient over a broadcast
// channel.
type MemberShareMessage struct {
	Share bls.SecretKey
}

// AccusationsMessage is a message payload that carries all of the sender's
// accusations against other members of the threshold group. If all other
// members behaved honestly from the sender's point of view, this message should
// be broadcast but with an empty slice of `accusedIDs`. It is expected to be
// broadcast.
type AccusationsMessage struct {
	accusedIDs []bls.ID
}

// JustificationsMessage is a message payload that carries all of the sender's
// justifications in response to other threshold group members' accusations. If
// no other member accused the sender, this message should be broadcast but with
// an empty map of `justifications`. It is expected to be broadcast.
type JustificationsMessage struct {
	justifications map[bls.ID]bls.SecretKey
}

// Execute runs the full distributed key generation lifecycle, given a broadcast
// channel to mediate it and a group size and threshold. It returns a threshold
// group member who is participating in the group if the generation was
// successful, and an error representing what went wrong if not.
func Execute(blockCounter chain.BlockCounter, channel broadcast.Channel, groupSize int, threshold int) (*thresholdgroup.Member, error) {
	// FIXME Probably pass in a way to ask for a receiver's public key?
	// FIXME Need a way to time out in a given stage, especially the waiting
	//       ones.

	memberID := rand.NewRand().String()
	fmt.Printf("[member:%v] Initializing member.\n", memberID)
	localMember := thresholdgroup.NewMember(memberID, threshold)

	recvChan := channel.RecvChan()

	fmt.Printf("[member:%v] Waiting for join timeout...\n", memberID)
	blockCounter.WaitForBlocks(5)

	fmt.Printf("[member:%v] Broadcasting join.\n", memberID)
	channel.Send(broadcast.NewBroadcastMessage(localMember.BlsID, JoinMessage{}))

	// Wait for all members.
	waiter := blockCounter.BlockWaiter(3)
	fmt.Printf("[member:%v] Waiting for other members...\n", memberID)
	memberIDs, err := waitForMemberIDs(&localMember.BlsID, recvChan, groupSize)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all member ids: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for member join timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(3)
	fmt.Printf("[member:%v] Initiating commitment broadcast phase.\n", memberID)
	sharingMember := localMember.InitializeSharing(memberIDs)

	fmt.Printf("[member:%v] Broadcasting public commitment.\n", memberID)
	err = sendCommitments(channel, &sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast commitments: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other commitments...\n", memberID)
	err = waitForCommitments(&localMember.BlsID, recvChan, &sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all commitments: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for commitment timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(5)
	fmt.Printf("[member:%v] Sending private shares.\n", memberID)
	err = sendShares(channel, &sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to send all private shares: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other shares...\n", memberID)
	err = waitForShares(&sharingMember.BlsID, recvChan, &sharingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all private shares: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for share exchange timeout...\n", memberID)
	<-waiter

	waiter = blockCounter.BlockWaiter(3)
	fmt.Printf("[member:%v] Initiating accusation/justification phase.\n", memberID)
	justifyingMember := sharingMember.InitializeJustification()
	fmt.Printf("[member:%v] Broadcasting accusations.\n", memberID)
	err = sendAccusations(channel, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast accusations: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other accusations...\n", memberID)
	err = waitForAccusations(recvChan, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all accusations: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for accusation timeout...\n", memberID)
	<-waiter

	fmt.Printf("[member:%v] Broadcasting justifications.\n", memberID)
	err = sendJustifications(channel, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast justifications: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other justifications...\n", memberID)
	err = waitForJustifications(recvChan, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all justifications: [%v]", err)
	}

	fmt.Printf("[member:%v] Finalizing member.\n", memberID)
	member := justifyingMember.FinalizeMember()
	return &member, nil
}

func waitForMemberIDs(myID *bls.ID, recvChan <-chan broadcast.Message, groupSize int) ([]bls.ID, error) {
	memberIDs := make([]bls.ID, 0, groupSize)

done:
	for msg := range recvChan {
		switch msg.Data.(type) {
		case JoinMessage:
			if msg.Sender.IsEqual(myID) {
				continue
			}

			memberIDs = append(memberIDs, msg.Sender)

			if len(memberIDs) == groupSize-1 {
				break done
			}
		}
	}

	return memberIDs, nil
}

func sendCommitments(channel broadcast.Channel, member *thresholdgroup.SharingMember) error {
	channel.Send(broadcast.NewBroadcastMessage(member.BlsID, MemberCommitmentsMessage{member.Commitments()}))

	return nil
}

func waitForCommitments(myID *bls.ID, recvChan <-chan broadcast.Message, sharingMember *thresholdgroup.SharingMember) error {
done:
	for msg := range recvChan {
		switch commitmentMsg := msg.Data.(type) {
		case MemberCommitmentsMessage:
			if msg.Sender.IsEqual(myID) {
				continue
			}

			sharingMember.AddCommitmentsFromID(msg.Sender, commitmentMsg.Commitments)

			if sharingMember.CommitmentsComplete() {
				break done
			}
		}
	}

	return nil
}

func sendShares(channel broadcast.Channel, member *thresholdgroup.SharingMember) error {
	for _, receiverID := range member.OtherMemberIDs() {
		share := member.SecretShareForID(receiverID)
		fmt.Printf("[member:%v] Despatching a share!\n", member.ID)
		channel.Send(broadcast.NewPrivateMessage(member.BlsID, receiverID, MemberShareMessage{share}))
	}

	return nil
}

func waitForShares(myID *bls.ID, recvChan <-chan broadcast.Message, sharingMember *thresholdgroup.SharingMember) error {
done:
	for msg := range recvChan {
		switch shareMsg := msg.Data.(type) {
		case MemberShareMessage:
			if msg.Receiver.IsEqual(myID) {
				fmt.Printf("[member:%v] Received one id from [%v].\n", myID.GetHexString(), msg.Sender.GetHexString())
				sharingMember.AddShareFromID(msg.Sender, shareMsg.Share)

				if sharingMember.SharesComplete() {
					break done
				}
			}
		}
	}

	return nil
}

func sendAccusations(channel broadcast.Channel, member *thresholdgroup.JustifyingMember) error {
	channel.Send(broadcast.NewBroadcastMessage(member.BlsID, AccusationsMessage{member.AccusedIDs()}))

	return nil
}

func waitForAccusations(recvChan <-chan broadcast.Message, justifyingMember *thresholdgroup.JustifyingMember) error {
	memberIDs := justifyingMember.OtherMemberIDs()
	seenAccusations := make(map[bls.ID]bool, len(memberIDs))
done:
	for msg := range recvChan {
		switch accusationMsg := msg.Data.(type) {
		case AccusationsMessage:
			for _, accusedID := range accusationMsg.accusedIDs {
				justifyingMember.AddAccusationFromID(msg.Sender, accusedID)
			}

			seenAccusations[msg.Sender] = true
			if len(seenAccusations) == len(memberIDs) {
				break done
			}
		}
	}

	return nil
}

func sendJustifications(channel broadcast.Channel, justifyingMember *thresholdgroup.JustifyingMember) error {
	channel.Send(
		broadcast.NewBroadcastMessage(
			justifyingMember.BlsID,
			JustificationsMessage{justifyingMember.Justifications()}))

	return nil
}

func waitForJustifications(recvChan <-chan broadcast.Message, justifyingMember *thresholdgroup.JustifyingMember) error {
	memberIDs := justifyingMember.OtherMemberIDs()
	seenJustifications := make(map[bls.ID]bool, len(memberIDs))
done:
	for msg := range recvChan {
		switch justificationsMsg := msg.Data.(type) {
		case JustificationsMessage:
			for accuserID, justification := range justificationsMsg.justifications {
				justifyingMember.RecordJustificationFromID(msg.Sender, accuserID, justification)
			}

			seenJustifications[msg.Sender] = true
			if len(seenJustifications) == len(memberIDs) {
				break done
			}
		}
	}

	return nil
}
