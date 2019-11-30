package entry

import (
	"math/big"

	"github.com/ipfs/go-log"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-entry")

const (
	setupBlocks     = state.MessagingStateDelayBlocks
	signatureBlocks = state.MessagingStateActiveBlocks
)

func initializeChannel(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &SignatureShareMessage{} })
}

// SignAndSubmit triggers the threshold signature process for the
// previous relay entry and publishes the signature to the chain as
// a new relay entry.
func SignAndSubmit(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	previousEntry *big.Int,
	honestThreshold int,
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
) error {
	initializeChannel(channel)

	initialState := &signatureShareState{
		signingStateBase: signingStateBase{
			channel:         channel,
			relayChain:      relayChain,
			blockCounter:    blockCounter,
			signer:          signer,
			previousEntry:   previousEntry,
			honestThreshold: honestThreshold,
		},
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)
	_, _, err := stateMachine.Execute(startBlockHeight)

	return err
}
