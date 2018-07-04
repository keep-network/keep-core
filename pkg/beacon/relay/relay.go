package relay

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

// NewNode returns an empty Node with no group, zero group count, and a nil last
// seen entry, tied to the given net.Provider.
func NewNode(
	netProvider net.Provider,
	blockCounter chain.BlockCounter,
	chainConfig config.Chain,
) Node {
	return Node{
		netProvider:  netProvider,
		blockCounter: blockCounter,
		chainConfig:  chainConfig,
	}
}
