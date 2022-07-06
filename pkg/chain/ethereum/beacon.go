package ethereum

import (
	"context"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
)

// BeaconChain represents a beacon-specific chain handle.
type BeaconChain struct {
	*Chain
}

// NewBeaconChain construct a new instance of the beacon-specific Ethereum
// chain handle.
func NewBeaconChain(
	ctx context.Context,
	config *ethereum.Config,
	client ethutil.EthereumClient,
) *BeaconChain {
	return &BeaconChain{
		Chain: NewChain(ctx, config, client),
	}
}
