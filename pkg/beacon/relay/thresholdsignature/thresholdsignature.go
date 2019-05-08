package thresholdsignature

import (
	"fmt"
	"os"
	"sync"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/bls"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

const (
	setupBlocks     = state.MessagingStateDelayBlocks
	signatureBlocks = state.MessagingStateActiveBlocks
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
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
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

	setupDelay := startBlockHeight + setupBlocks
	fmt.Printf(
		"[member:%v] Waiting for block [%v] to start threshold signing...\n",
		signer.MemberID(),
		setupDelay,
	)

	err := blockCounter.WaitForBlockHeight(setupDelay)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to wait %d blocks entering threshold setup: [%v]",
			setupBlocks,
			err,
		)
	}

	fmt.Printf("[member:%v] Sending signature share...\n", signer.MemberID())

	seenSharesMutex := sync.Mutex{}
	seenShares := make(map[group.MemberIndex]*bn256.G1)
	share := signer.CalculateSignatureShare(bytes)

	// Add local share to map rather than receiving from the network.
	seenShares[signer.MemberID()] = share

	err = sendSignatureShare(share.Marshal(), channel, signer.MemberID())
	if err != nil {
		return nil, err
	}

	blockWaiter, err := blockCounter.BlockHeightWaiter(setupDelay + signatureBlocks)
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
					fmt.Printf("[member:%v] Ignoring my own message (senderID = [%v])\n", signer.MemberID(), signatureShareMsg.senderID)
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
					seenSharesMutex.Lock()
					fmt.Printf("[member:%v] Adding signature share from sender [%v] = [%v]\n", signer.MemberID(), signatureShareMsg.senderID, share)
					seenShares[signatureShareMsg.senderID] = share
					seenSharesMutex.Unlock()
				}
			}
		case endBlockHeight := <-blockWaiter:
			fmt.Printf(
				"[member:%v] Stopped receiving signature shares at block [%v]\n",
				signer.MemberID(),
				endBlockHeight,
			)

			// put all seen shares into a slice and complete the signature
			seenSharesSlice := make([]*bls.SignatureShare, 0)
			for memberID, share := range seenShares {
				signatureShare := &bls.SignatureShare{I: int(memberID), V: share}
				seenSharesSlice = append(seenSharesSlice, signatureShare)
			}

			fmt.Printf("[member:%v] All seen shares:\n", signer.MemberID())
			for _, share := range seenSharesSlice {
				fmt.Printf("[member:%v] [I = %v, V = %v]\n", signer.MemberID(), share.I, share.V)
			}

			signature, err := signer.CompleteSignature(seenSharesSlice, threshold)
			if err != nil {
				return nil, err
			}

			fmt.Printf("[member:%v] Evaluated signature = [%v]\n", signer.MemberID(), signature)

			return altbn128.G1Point{G1: signature}.Compress(), nil
		}
	}
}

func sendSignatureShare(
	share []byte,
	channel net.BroadcastChannel,
	memberID group.MemberIndex,
) error {
	fmt.Printf("[member:%v] Sending my signature share: [%v]\n", memberID, share)
	return channel.Send(&SignatureShareMessage{memberID, share})
}
