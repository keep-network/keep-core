package thresholdsignature

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

const (
	setupBlocks     = state.MessagingStateDelayBlocks
	signatureBlocks = state.MessagingStateActiveBlocks
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func initializeChannel(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &SignatureShareMessage{} })
}

// Execute triggers the threshold signature process for the given bytes.
func Execute(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	relayChain relaychain.RelayEntryInterface,
	requestID *big.Int,
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
			requestID:     requestID,
			previousEntry: previousEntry,
			seed:          seed,
			threshold:     threshold,
		},
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)
	_, _, err := stateMachine.Execute(startBlockHeight)

	return err
}

// CombineEntryToSign takes the previous relay entry value and the current
// requests's seed and combines it into a slice of bytes that is going to be
// signed by the selected group and as a result, will form a new relay entry
// value.
func CombineEntryToSign(previousEntry *big.Int, seed *big.Int) []byte {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntry.Bytes()...)
	combinedEntryToSign = append(combinedEntryToSign, seed.Bytes()...)
	return combinedEntryToSign
}
