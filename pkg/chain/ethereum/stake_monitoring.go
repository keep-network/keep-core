package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumStakeMonitoring struct {
	config *ethereumChain
}

func (esm *ethereumStakeMonitoring) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return esm.config.HasMinimumStake(common.HexToAddress(address))
}

func (ec *ethereumChain) StakeMonitoring() (chain.StakeMonitoring, error) {
	stakeMonitoring := &ethereumStakeMonitoring{
		config: ec,
	}

	return stakeMonitoring, nil
}
