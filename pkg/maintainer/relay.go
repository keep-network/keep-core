package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"

)

func NewRelay(btcChain bitcoin.Chain, tbtcChain tbtc.Chain) *Relay {
	return &Relay{
		btcChain:  btcChain,
		tbtcChain: tbtcChain,
	}
}

type Relay struct {
	btcChain  bitcoin.Chain
	tbtcChain tbtc.Chain
}
