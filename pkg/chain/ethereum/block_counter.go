package ethereum

import "github.com/keep-network/keep-core/pkg/chain"

// BlockCounter creates a BlockCounter that uses the block number in Ethereum.
func (c *Chain) BlockCounter() (chain.BlockCounter, error) {
	return c.blockCounter, nil
}