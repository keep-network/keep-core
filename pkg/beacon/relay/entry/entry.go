package entry

import (
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

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform relay entry signing protocol interactions by registering all the
// required protocol message unmarshallers.
// The channel has to be initialized before the SignAndSubmit is called.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
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
	previousEntry []byte,
	honestThreshold int,
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
) error {
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
