package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// BitcoinDifficultyChain is an interface that provides the ability to
// communicate with the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyChain interface {
	// Ready checks whether the Bitcoin difficulty chain is active (i.e. genesis
	// has been performed).
	Ready() (bool, error)

	// IsAuthorizationRequired checks whether the Bitcoin difficulty chain
	// requires the address submitting a retarget to be authorised in advance by
	// governance.
	IsAuthorizationRequired() (bool, error)

	// IsAuthorized checks whether the given address has been authorised to
	// submit a retarget by governance.
	IsAuthorized(address chain.Address) (bool, error)

	// Signing returns the signing associated with the chain.
	Signing() chain.Signing

	// Retarget adds a new epoch to the Bitcoin difficulty chain by providing
	// a proof of the difficulty before and after the retarget.
	Retarget(headers []*bitcoin.BlockHeader) error

	// CurrentEpoch returns the number of the latest epoch whose difficulty is
	// proven in the Bitcoin difficulty chain. If the genesis epoch's number is
	// set correctly, and retargets along the way have been legitimate,
	// the current epoch equals the height of the block starting the most recent
	// epoch, divided by 2016.
	CurrentEpoch() (uint64, error)

	// ProofLength returns the number of blocks required for each side of a
	// retarget proof: a retarget must provide `proofLength` blocks before
	// the retarget and `proofLength` blocks after it.
	ProofLength() (uint64, error)
}
