package ethereum

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
)

// TBTCChain represents a TBTC-specific chain handle.
type TBTCChain struct {
	*Chain
}

// NewTBTCChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func NewTBTCChain(
	ctx context.Context,
	config *ethereum.Config,
	client ethutil.EthereumClient,
) (*TBTCChain, error) {
	chain, err := NewChain(ctx, config, client)
	if err != nil {
		return nil, fmt.Errorf("cannot create base chain handle: [%v]", err)
	}

	return &TBTCChain{
		Chain: chain,
	}, nil
}