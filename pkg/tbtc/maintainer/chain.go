package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// BitcoinDifficultyChain is an interface that provides the ability to
// communicate with the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyChain interface {
	// Retarget adds a new epoch to the Bitcoin difficulty relay by providing
	// a proof of the difficulty before and after the retarget.
	Retarget(headers []bitcoin.BlockHeader) error
}

// TODO: Description
type WalletChain interface {
	// TODO: Implement
}
