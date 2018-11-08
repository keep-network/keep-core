package local

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// StakeMonitoring implements `chain.StakeMonitoring` interface and works
// as a local stub for testing.
type StakeMonitoring struct {
	stakers map[string]bool
}

// NewStakeMonitoring creates a new instance of `StakeMonitoring` test stub.
func NewStakeMonitoring() *StakeMonitoring {
	return &StakeMonitoring{
		stakers: make(map[string]bool),
	}
}

// HasMinimumStake checks if the provided address staked the number of
// ERC20 KEEP tokens above the required minimum to become a network operator.
// The minimum number of KEEP tokens required to be staked is an on-chain
// parameter.
func (lsm *StakeMonitoring) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return lsm.stakers[address], nil
}

// StakeTokens stakes for the provided address enough KEEP tokens allowing to
// become a network operator.
func (lsm *StakeMonitoring) StakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	lsm.stakers[address] = true
	return nil
}

// UnstakeTokens unstakes from the provided address all KEEP tokens so it can
// no longer be a network operator.
func (lsm *StakeMonitoring) UnstakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	delete(lsm.stakers, address)
	return nil
}
