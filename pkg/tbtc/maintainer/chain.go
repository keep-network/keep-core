package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// BitcoinDifficultyChain is an interface that provides the ability to
// communicate with the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyChain interface {
	// Ready checks whether the relay is active (i.e. genesis has been performed).
	// Note that if the relay is used by querying the current and previous epoch
	// difficulty, at least one retarget needs to be provided after genesis;
	// otherwise the prevEpochDifficulty will be uninitialised and zero.
	Ready() (bool, error)

	// AuthorizationRequired checks whether the relay requires the address
	// submitting a retarget to be authorised in advance by governance.
	AuthorizationRequired() (bool, error)

	// IsAuthorized checks whether the given address has been authorised by
	// governance to submit a retarget.
	IsAuthorized(address chain.Address) (bool, error)

	// Signing returns the signing associated with the chain.
	Signing() chain.Signing

	// Retarget adds a new epoch to the relay by providing a proof
	// of the difficulty before and after the retarget.
	Retarget(headers []*bitcoin.BlockHeader) error

	// CurrentEpoch returns the number of the latest epoch whose difficulty is
	// proven to the relay.
	CurrentEpoch() (uint64, error)

	// ProofLength returns the number of blocks required for each side of a
	// retarget proof.
	ProofLength() (uint64, error)
}
