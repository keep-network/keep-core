package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type RelayChain interface {
	Retarget(headers []bitcoin.BlockHeader) error
}

func NewRelay(btcChain bitcoin.Chain, relayChain RelayChain) *Relay {
	return &Relay{
		btcChain:   btcChain,
		relayChain: relayChain,
	}
}

type Relay struct {
	btcChain   bitcoin.Chain
	relayChain RelayChain
}
