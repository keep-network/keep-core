package ethereum

import (
	"context"
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
) *TBTCChain {
	return &TBTCChain{
		Chain: NewChain(ctx, config, client),
	}
}