package entry

import (
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

import "github.com/ipfs/go-log"

var logger = log.Logger("keep-entry")

const (
	setupBlocks     = state.MessagingStateDelayBlocks
	signatureBlocks = state.MessagingStateActiveBlocks
)

func initializeChannel(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &SignatureShareMessage{} })
}

// SignAndSubmit triggers the threshold signature process for the combination of
// the previous relay entry and seed and publishes the signature to the chain as
// a new relay entry.
func SignAndSubmit(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	previousEntry *big.Int,
	seed *big.Int,
	threshold int,
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
) error {
	initializeChannel(channel)

	initialState := &signatureShareState{
		signingStateBase: signingStateBase{
			channel:       channel,
			relayChain:    relayChain,
			blockCounter:  blockCounter,
			signer:        signer,
			previousEntry: previousEntry,
			seed:          seed,
			threshold:     threshold,
		},
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)
	_, _, err := stateMachine.Execute(startBlockHeight)

	return err
}

// CombineToSign takes the previous relay entry value and the current
// requests's seed and combines it into a slice of bytes that is going to be
// signed by the selected group and as a result, will form a new relay entry
// value.
func CombineToSign(previousEntry *big.Int, seed *big.Int) []byte {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntry.Bytes()...)
	combinedEntryToSign = append(combinedEntryToSign, seed.Bytes()...)
	return combinedEntryToSign
}
