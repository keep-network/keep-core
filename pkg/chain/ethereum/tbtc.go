package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*Chain
}

// NewTbtcChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func NewTbtcChain(
	ctx context.Context,
	config ethereum.Config,
	client *ethclient.Client,
) (*TbtcChain, error) {
	chain, err := NewChain(ctx, config, client)
	if err != nil {
		return nil, fmt.Errorf("cannot create base chain handle: [%v]", err)
	}

	return &TbtcChain{
		Chain: chain,
	}, nil
}
