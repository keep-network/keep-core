package dkg

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/rand"
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
func ExecuteDKG(blockCounter chain.BlockCounter, channel net.BroadcastChannel, groupSize int, threshold int) (*thresholdgroup.Member, error) {
	// FIXME Probably pass in a way to ask for a receiver's public key?
	// FIXME Need a way to time out in a given stage, especially the waiting
	//       ones.

	// Generate a nonzero memberID; loop until rand.NewRand returns something
	// other than 0, hopefully no more than once :)
	memberID := "0"
	for memberID = rand.NewRand().String(); memberID == "0"; {
	}
	fmt.Printf("[member:%v] Initializing member.\n", memberID)
	localMember := thresholdgroup.NewMember(memberID, threshold)

	recvChan := make(chan interface{})
	channel.Recv(func(msg interface{}) error {
		recvChan <- msg
		return nil
	})

	fmt.Printf("[member:%v] Waiting for join timeout...\n", memberID)
	blockCounter.WaitForBlocks(15)

	fmt.Printf("[member:%v] Broadcasting join.\n", memberID)
	channel.Send(&JoinMessage{&localMember.BlsID})

	// Wait for all members.
	waiter := blockCounter.BlockWaiter(10)
	fmt.Printf("[member:%v] Waiting for other members...\n", memberID)
	memberIDs, err := waitForMemberIDs(&localMember.BlsID, recvChan, groupSize)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all member ids: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for member join timeout...\n", memberID)
	<-waiter

	fmt.Printf("[member:%v] Saw IDs: %v\n", memberID, len(memberIDs))

	waiter = blockCounter.BlockWaiter(15)
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

	waiter = blockCounter.BlockWaiter(20)
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

	waiter = blockCounter.BlockWaiter(15)
	fmt.Printf("[member:%v] Initiating accusation/justification phase.\n", memberID)
	justifyingMember := sharingMember.InitializeJustification()
	fmt.Printf("[member:%v] Broadcasting accusations.\n", memberID)
	err = sendAccusations(channel, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast accusations: [%v]", err)
	}

	fmt.Printf("[member:%v] Waiting for other accusations...\n", memberID)
	err = waitForAccusations(&justifyingMember.BlsID, recvChan, &justifyingMember)
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
	err = waitForJustifications(&justifyingMember.BlsID, recvChan, &justifyingMember)
	if err != nil {
		return nil, fmt.Errorf("failed to receive all justifications: [%v]", err)
	}

	fmt.Printf("[member:%v] Finalizing member.\n", memberID)
	member := justifyingMember.FinalizeMember()
	return &member, nil
}

func waitForMemberIDs(myID *bls.ID, recvChan <-chan interface{}, groupSize int) ([]bls.ID, error) {
	memberIDs := make([]bls.ID, 0, groupSize)

done:
	for msg := range recvChan {
		switch joinMsg := msg.(type) {
		case JoinMessage:
			if joinMsg.id.IsEqual(myID) {
				continue
			}

			memberIDs = append(memberIDs, *joinMsg.id)

			if len(memberIDs) == groupSize-1 {
				break done
			}
		}
	}

	return memberIDs, nil
}

func sendCommitments(channel net.BroadcastChannel, member *thresholdgroup.SharingMember) error {
	channel.Send(&MemberCommitmentsMessage{&member.BlsID, member.Commitments()})

	return nil
}

func waitForCommitments(myID *bls.ID, recvChan <-chan interface{}, sharingMember *thresholdgroup.SharingMember) error {
done:
	for msg := range recvChan {
		switch commitmentMsg := msg.(type) {
		case MemberCommitmentsMessage:
			if commitmentMsg.id.IsEqual(myID) {
				continue
			}

			sharingMember.AddCommitmentsFromID(*commitmentMsg.id, commitmentMsg.Commitments)

			if sharingMember.CommitmentsComplete() {
				break done
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
			net.ClientIdentifier(receiverID.GetHexString()),
			&MemberShareMessage{&member.BlsID, &receiverID, &share})
	}
	fmt.Printf("[member:%v] Shares despatched!\n", member.ID)

	return nil
}

func waitForShares(myID *bls.ID, recvChan <-chan interface{}, sharingMember *thresholdgroup.SharingMember) error {
done:
	for msg := range recvChan {
		switch shareMsg := msg.(type) {
		case MemberShareMessage:
			if shareMsg.receiverID.IsEqual(myID) {
				sharingMember.AddShareFromID(*shareMsg.id, *shareMsg.Share)

				if sharingMember.SharesComplete() {
					break done
				}
			}
		}
	}

	return nil
}

func sendAccusations(channel net.BroadcastChannel, member *thresholdgroup.JustifyingMember) error {
	channel.Send(&AccusationsMessage{&member.BlsID, member.AccusedIDs()})

	return nil
}

func waitForAccusations(myID *bls.ID, recvChan <-chan interface{}, justifyingMember *thresholdgroup.JustifyingMember) error {
	memberIDs := justifyingMember.OtherMemberIDs()
	seenAccusations := make(map[bls.ID]bool, len(memberIDs))
done:
	for msg := range recvChan {
		switch accusationMsg := msg.(type) {
		case AccusationsMessage:
			if accusationMsg.id.IsEqual(myID) {
				continue
			}

			for _, accusedID := range accusationMsg.accusedIDs {
				justifyingMember.AddAccusationFromID(*accusationMsg.id, accusedID)
			}

			seenAccusations[*accusationMsg.id] = true
			if len(seenAccusations) == len(memberIDs) {
				break done
			}
		}
	}

	return nil
}

func sendJustifications(channel net.BroadcastChannel, justifyingMember *thresholdgroup.JustifyingMember) error {
	channel.Send(
		&JustificationsMessage{
			&justifyingMember.BlsID,
			justifyingMember.Justifications()})

	return nil
}

func waitForJustifications(myID *bls.ID, recvChan <-chan interface{}, justifyingMember *thresholdgroup.JustifyingMember) error {
	memberIDs := justifyingMember.OtherMemberIDs()
	seenJustifications := make(map[bls.ID]bool, len(memberIDs))
done:
	for msg := range recvChan {
		switch justificationsMsg := msg.(type) {
		case JustificationsMessage:
			if justificationsMsg.id.IsEqual(myID) {
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
				break done
			}
		}
	}

	return nil
}
