package ethereum

import "github.com/keep-network/keep-core/pkg/chain"

func (bc *baseChain) BlockCounter() (chain.BlockCounter, error) {
	return bc.blockCounter, nil
}
