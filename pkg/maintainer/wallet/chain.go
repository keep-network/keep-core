package wallet

import (
	"time"

	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// Chain represents the interface that the wallet maintainer module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// coordinator.Chain is an interface required by the coordinator package.
	coordinator.Chain

	// GetDepositSweepMaxSize gets the maximum number of deposits that can
	// be part of a deposit sweep proposal.
	GetDepositSweepMaxSize() (uint16, error)

	// GetWalletLock gets the current wallet lock for the given wallet.
	// Returned values represent the expiration time and the cause of the lock.
	// The expiration time can be UNIX timestamp 0 which means there is no lock
	// on the wallet at the given moment.
	GetWalletLock(
		walletPublicKeyHash [20]byte,
	) (time.Time, tbtc.WalletActionType, error)
}
