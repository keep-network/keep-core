package thresholdsignature

import (
	"fmt"
	"os"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg2"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/bls"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

const (
	setupBlocks     = 1
	signatureBlocks = 2
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
	threshold int,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	signer *dkg2.ThresholdSigner,
) ([]byte, error) {
	// Use an unbuffered channel to serialize message processing.
	recvChan := make(chan net.Message)
	handler := net.HandleMessageFunc{
		Type: fmt.Sprintf("relay/signature/%s", string(time.Now().UTC().UnixNano())),
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}

	// Initialize channel to perform threshold signing process.
	Init(channel)

	channel.Recv(handler)
	defer channel.UnregisterRecv(handler.Type)

	fmt.Printf(
		"[member:%v] Waiting for other group members to enter signing state...\n",
		signer.MemberID(),
	)

	err := blockCounter.WaitForBlocks(setupBlocks)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to wait %d blocks entering threshold setup: [%v]",
			setupBlocks,
			err,
		)
	}

	fmt.Printf("[member:%v] Sending signature share...\n", signer.MemberID())

	seenShares := make(map[group.MemberIndex]*bn256.G1)
	share := signer.CalculateSignatureShare(bytes)

	// Add local share to map rather than receiving from the network.
	seenShares[signer.MemberID()] = share

	err = sendSignatureShare(share.Marshal(), channel, signer.MemberID())
	if err != nil {
		return nil, err
	}

	blockWaiter, err := blockCounter.BlockWaiter(signatureBlocks)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[member:%v] Receiving other group signature share\n", signer.MemberID())

	for {
		select {
		case msg := <-recvChan:
			fmt.Printf(
				"[member:%v] Processing signing message\n",
				signer.MemberID(),
			)

			switch signatureShareMsg := msg.Payload().(type) {
			case *SignatureShareMessage:
				// Ignore our own share, we already have it.
				if signatureShareMsg.senderID == signer.MemberID() {
					continue
				}

				share := new(bn256.G1)
				_, err := share.Unmarshal(signatureShareMsg.ShareBytes)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"[member:%v] failed to unmarshal signature share: [%v]",
						signer.MemberID(),
						err,
					)
				} else {
					seenShares[signatureShareMsg.senderID] = share
				}
			}
		case <-blockWaiter:
			// put all seen shares into a slice and complete the signature
			seenSharesSlice := make([]*bls.SignatureShare, 0)
			for memberID, share := range seenShares {
				signatureShare := &bls.SignatureShare{I: int(memberID), V: share}
				seenSharesSlice = append(seenSharesSlice, signatureShare)
			}

			signature, err := signer.CompleteSignature(seenSharesSlice, threshold)
			if err != nil {
				return nil, err
			}

			return altbn128.G1Point{G1: signature}.Compress(), nil
		}
	}
}

func sendSignatureShare(
	share []byte,
	channel net.BroadcastChannel,
	memberID group.MemberIndex) error {
	return channel.Send(&SignatureShareMessage{memberID, share})
}
