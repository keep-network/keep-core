package thresholdsignature

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

const (
	setupBlocks     = 15
	signatureBlocks = 10
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &SignatureShareMessage{} })
}

// Execute triggers the threshold signature process for the given bytes. After
// the process has completed, it returns either the threshold signature's final
// bytes, or an error.
func Execute(
	bytes []byte,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	member *thresholdgroup.Member,
) ([]byte, error) {
	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	channel.Recv(func(msg net.Message) error {
		recvChan <- msg
		return nil
	})

	seenShares := make(map[bls.ID][]byte)

	fmt.Printf(
		"[member:%v, state:signing] Waiting for other group members to enter signing state...\n",
		member.MemberID(),
	)
	blockCounter.WaitForBlocks(15)

	fmt.Printf(
		"[member:%v] Sending signature share...\n",
		member.MemberID(),
	)
	err := sendSignatureShare(bytes, channel, member)
	if err != nil {
		return nil, err
	}

	blockWaiter, err := blockCounter.BlockWaiter(10)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[member:%v] Receiving other group signature share.\n", member.ID)

	for {
		select {
		case msg := <-recvChan:
			fmt.Printf(
				"[member:%v, state:signing] Processing message.\n",
				member.MemberID(),
			)

			switch signatureShareMsg := msg.Payload().(type) {
			case *SignatureShareMessage:
				if senderID, ok := msg.ProtocolSenderID().(*bls.ID); ok {
					if senderID.IsEqual(&member.BlsID) {
						continue
					}

					seenShares[*senderID] = signatureShareMsg.ShareBytes
				}
			}

		case <-blockWaiter:
			signature, err := member.CompleteSignature(seenShares)
			if err != nil {
				return signature.Serialize(), nil
			}

			return nil, fmt.Errorf(
				"[member:%v] failed to complete signature inside active period [%v]",
				member.MemberID(),
				signatureBlocks,
			)
		}
	}
}

func sendSignatureShare(
	bytes []byte,
	channel net.BroadcastChannel,
	member *thresholdgroup.Member,
) error {
	share := member.SignatureShare(string(bytes))
	return channel.Send(&SignatureShareMessage{&member.BlsID, share})
}
