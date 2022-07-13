package ethereum

import "github.com/keep-network/keep-core/pkg/chain"

func (c *Chain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}
