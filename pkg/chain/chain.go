package chain

import (
	"context"
	"github.com/keep-network/keep-core/pkg/operator"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks or watch as they are mined.
// Block height is expected to increase monotonically, though
// the time between blocks depends on the underlying implementation.
type BlockCounter interface {
	// WaitForBlockHeight blocks at the caller until the given block height is
	// reached. If the number of blocks is zero or negative or if the given
	// block height has been already reached, it returns immediately.
	WaitForBlockHeight(blockNumber uint64) error

	// BlockHeightWaiter returns a channel that will emit the block number after
	// the given block height is reached and then immediately close.
	// Reading from the returned channel immediately will effectively behave the
	// same way as calling WaitForBlockHeight.
	BlockHeightWaiter(blockNumber uint64) (<-chan uint64, error)

	// CurrentBlock returns the current block height.
	CurrentBlock() (uint64, error)

	// WatchBlocks returns a channel that will emit new block numbers as they
	// are mined. When the context provided as the parameter ends, new blocks
	// are no longer pushed to the channel and the channel is closed. If there
	// is no reader for the channel or reader is too slow, block updates can be
	// dropped.
	WatchBlocks(ctx context.Context) <-chan uint64
}

// StakeMonitor is an interface that provides ability to check and monitor
// the stake for the provided address.
type StakeMonitor interface {
	// HasMinimumStake checks if the specified account has enough active stake
	// to become network operator and that the operator contract the client is
	// working with has been authorized for potential slashing.
	//
	// Having the required minimum of active stake makes the operator eligible
	// to join the network. If the active stake is not currently undelegating,
	// operator is also eligible for work selection.
	HasMinimumStake(operatorPublicKey *operator.PublicKey) (bool, error)

	// StakerFor returns a Staker for the given operator public key.
	StakerFor(operatorPublicKey *operator.PublicKey) (Staker, error)
}

// Signing is an interface that provides ability to sign and verify
// signatures using operator's key associated with the chain.
type Signing interface {
	// PublicKey returns operator's public key in a serialized format.
	// The returned public key is used to Sign messages and can be later used
	// for verification.
	PublicKey() []byte

	// Sign the provided message with operator's private key. Returns the
	// signature or error in case signing failed.
	Sign(message []byte) ([]byte, error)

	// Verify the provided message against the signature using operator's
	// public key. Returns true if signature is valid and false otherwise.
	// If signature verification failed for some reason, an error is returned.
	Verify(message []byte, signature []byte) (bool, error)

	// VerifyWithPublicKey verifies the provided message against the signature
	// using the provided operator's public key. Returns true if signature is
	// valid and false otherwise. If signature verification failed for some
	// reason, an error is returned.
	VerifyWithPublicKey(
		message []byte,
		signature []byte,
		publicKey []byte,
	) (bool, error)

	// PublicKeyToAddress converts operator's public key to an address
	// associated with the chain.
	PublicKeyToAddress(publicKey *operator.PublicKey) ([]byte, error)

	// PublicKeyBytesToAddress converts operator's public key bytes to an address
	// associated with the chain.
	PublicKeyBytesToAddress(publicKey []byte) []byte
}

// Handle represents a handle to a blockchain that provides access to the core
// operator functionality needed for Keep network interactions.
type Handle interface {
	BlockCounter() (BlockCounter, error)
	StakeMonitor() (StakeMonitor, error)
	ThresholdRelay() relaychain.Interface
	Signing() Signing
}
