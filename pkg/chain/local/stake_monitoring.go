package local

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type localStakeMonitoring struct {
	stakers map[string]bool
}

func newLocalStakeMonitoring() *localStakeMonitoring {
	return &localStakeMonitoring{
		stakers: make(map[string]bool),
	}
}

func (lsm *localStakeMonitoring) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return lsm.stakers[address], nil
}

func (lsm *localStakeMonitoring) stakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	lsm.stakers[address] = true
	return nil
}

func (lsm *localStakeMonitoring) unstakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	delete(lsm.stakers, address)
	return nil
}
