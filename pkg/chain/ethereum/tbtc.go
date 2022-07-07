package ethereum

import (
	"context"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*Chain
}

// NewTbtcChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func newTbtcChain(
	ctx context.Context,
	baseChain *Chain,
) (*TbtcChain, error) {

	return &TbtcChain{
		Chain: baseChain,
	}, nil
}
