package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// RelayChain is an interface that provides the ability to communicate with the
// relay on-chain contract.
type RelayChain interface {
	// Retarget adds a new epoch to the relay by providing a proof
	// of the difficulty before and after the retarget.
	Retarget(headers []bitcoin.BlockHeader) error
}

func NewRelay(btcChain bitcoin.Chain, relayChain RelayChain) *Relay {
	return &Relay{
		btcChain:   btcChain,
		relayChain: relayChain,
	}
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain   bitcoin.Chain
	relayChain RelayChain
}
