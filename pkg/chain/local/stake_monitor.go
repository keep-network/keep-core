package local

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// StakeMonitor implements `chain.StakeMonitor` interface and works
// as a local stub for testing.
type StakeMonitor struct {
	stakers map[string]bool
}

// NewStakeMonitor creates a new instance of `StakeMonitor` test stub.
func NewStakeMonitor() *StakeMonitor {
	return &StakeMonitor{
		stakers: make(map[string]bool),
	}
}

// HasMinimumStake checks if the provided address staked enough to become
// a network operator. The minimum stake is an on-chain parameter.
func (lsm *StakeMonitor) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return lsm.stakers[address], nil
}

// StakeTokens stakes for the provided address enough number of tokens allowing
// to become a network operator.
func (lsm *StakeMonitor) StakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	lsm.stakers[address] = true
	return nil
}

// UnstakeTokens unstakes from the provided address all tokens so it can
// no longer be a network operator.
func (lsm *StakeMonitor) UnstakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	delete(lsm.stakers, address)
	return nil
}
