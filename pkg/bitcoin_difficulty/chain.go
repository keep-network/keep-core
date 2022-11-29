package btcdiff

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// Chain defines the interface that pertains specifically to the process of
// updating Bitcoin difficulty.
type Chain interface {
	// Ready checks whether the relay is active (i.e. genesis has been
	// performed).
	Ready() (bool, error)

	// AuthorizationRequired checks whether the relay requires the address
	// submitting a retarget to be authorised in advance by governance.
	AuthorizationRequired() (bool, error)

	// IsAuthorized checks whether the given address has been authorised to
	// submit a retarget by governance.
	IsAuthorized(address chain.Address) (bool, error)

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
