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
			seed:            seed,
			honestThreshold: honestThreshold,
		},
		signingStartBlockHeight: startBlockHeight,
	}

	stateMachine := state.NewMachine(channel, blockCounter, initialState)
	_, _, err := stateMachine.Execute(startBlockHeight)

	return err
}

// CombineToSign takes the previous relay entry value and the current
// requests's seed and:
// - pad them with zeros if their byte length is less than 32 bytes. These
//   values are used later on-chain as `uint256` values and are combined using
//   `abi.encodePacked` function during signature verification. This function
//   pads `uint256` type values with zeros, if they byte length is less than 32.
//   If such values are not also padding off-chain, the on-chain verification
//   will fail because of the padding difference.
// - combines it into a slice of bytes that is going to be signed by the
//   selected group and as a result, will form a new relay entry value.
func CombineToSign(previousEntry *big.Int, seed *big.Int) []byte {
	previousEntryBytes := toPaddedBytes(previousEntry, 32)
	seedBytes := toPaddedBytes(seed, 32)

	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntryBytes...)
	combinedEntryToSign = append(combinedEntryToSign, seedBytes...)
	return combinedEntryToSign
}

func toPaddedBytes(value *big.Int, minimumByteLength int) []byte {
	valueBytes := value.Bytes()
	valuePaddingBytes := minimumByteLength - len(valueBytes)
	for i := 0; i < valuePaddingBytes; i++ {
		valueBytes = append([]byte{0x00}, valueBytes...)
	}
	return valueBytes
}
