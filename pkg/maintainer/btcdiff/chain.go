package btcdiff

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// Chain is an interface that provides the ability to
// communicate with the Bitcoin difficulty on-chain contract.
type Chain interface {
	// Ready checks whether the relay is active (i.e. genesis has been performed).
	// Note that if the relay is used by querying the current and previous epoch
	// difficulty, at least one retarget needs to be provided after genesis;
	// otherwise the prevEpochDifficulty will be uninitialized and zero.
	Ready() (bool, error)

	// IsAuthorized checks whether the given address has been authorized to
	// submit a retarget directly to LightRelay. This function should be used
	// when retargetting via LightRelayMaintainerProxy is disabled.
	IsAuthorized(address chain.Address) (bool, error)

	// IsAuthorizedForRefund checks whether the given address has been
	// authorized to submit a retarget via LightRelayMaintainerProxy. This
	// function should be used when retargetting via LightRelayMaintainerProxy
	// is not disabled.
	IsAuthorizedForRefund(address chain.Address) (bool, error)

	// Signing returns the signing associated with the chain.
	Signing() chain.Signing

	// Retarget adds a new epoch to the relay by providing a proof of the
	// difficulty before and after the retarget. The cost of calling this
	// function is not refunded to the caller.
	Retarget(headers []*bitcoin.BlockHeader) error

	// RetargetWithRefund adds a new epoch to the relay by providing a proof of
	// the difficulty before and after the retarget. The cost of calling this
	// function is refunded to the caller.
	RetargetWithRefund(headers []*bitcoin.BlockHeader) error

	// CurrentEpoch returns the number of the latest difficulty epoch which is
	// proven to the relay. If the genesis epoch's number is set correctly, and
	// retargets along the way have been legitimate, this equals the height of
	// the block starting the most recent epoch, divided by 2016.
	CurrentEpoch() (uint64, error)

	// ProofLength returns the number of blocks required for each side of a
	// retarget proof.
	ProofLength() (uint64, error)
}
