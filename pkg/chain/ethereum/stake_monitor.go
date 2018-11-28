package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumStakeMonitor struct {
	config *ethereumChain
}

// HasMinimumStake checks if the provided address staked enough to become
// a network operator. The minimum stake is an on-chain parameter.
func (esm *ethereumStakeMonitor) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return esm.config.HasMinimumStake(common.HexToAddress(address))
}

// StakeMonitor creates and returns a StakeMonitor instance operating on
// Ethereum chain.
func (ec *ethereumChain) StakeMonitor() (chain.StakeMonitor, error) {
	stakeMonitor := &ethereumStakeMonitor{
		config: ec,
	}

	return stakeMonitor, nil
}
