package relay

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/go/beacon/broadcast"
	"github.com/keep-network/keep-core/go/beacon/chain"
	"github.com/keep-network/keep-core/go/thresholdgroup"
)

// GroupSignatureShareMessage is a message payload that carries the sender's group signature share for the given message
type GroupSignatureShareMessage struct {
	Share []byte
}

// ExecuteGroupSignature triggers group signature process for the given message
func ExecuteGroupSignature(message string, blockCounter chain.BlockCounter, channel broadcast.Channel, member *thresholdgroup.Member) error {

	recvChan := channel.RecvChan()

	blockCounter.WaitForBlocks(15)

	sendGroupSignatureShare(message, channel, member)

	waiter := blockCounter.BlockWaiter(10)
	waitForGroupSignatureShares(&member.BlsID, recvChan, member)
	fmt.Printf("[member:%v] Waiting for other group signature share...\n", member.ID)
	<-waiter

	return nil
}

func sendGroupSignatureShare(message string, channel broadcast.Channel, member *thresholdgroup.Member) error {
	share := member.SignatureShare(message)
	fmt.Printf("[member:%v] Despatching group signature share!\n", member.ID)
	channel.Send(broadcast.NewBroadcastMessage(member.BlsID, GroupSignatureShareMessage{share}))
	fmt.Printf("[member:%v] Group signature share despatched!\n", member.ID)

	return nil
}

func waitForGroupSignatureShares(myID *bls.ID, recvChan <-chan broadcast.Message, member *thresholdgroup.Member) error {
done:
	for msg := range recvChan {
		switch shareMsg := msg.Data.(type) {
		case GroupSignatureShareMessage:

			if msg.Sender.IsEqual(myID) {
				continue
			}

			member.AddGroupSignatureShareFromID(msg.Sender, shareMsg.Share)

			if member.GroupSignatureSharesComplete() {
				fmt.Printf("[member:%v] Group signature process is complete\n", member.ID)
				break done
			}
		}
	}

	return nil
}
