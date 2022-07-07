package ethereum

import (
	"context"
	"fmt"
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
) (*BeaconChain, error) {
	chain, err := NewChain(ctx, config, client)
	if err != nil {
		return nil, fmt.Errorf("cannot create base chain handle: [%v]", err)
	}

	return &BeaconChain{
		Chain: chain,
	}, nil
}
