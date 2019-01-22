package local

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain"
)

var minimumStake = big.NewInt(20000000)

// StakeMonitor implements `chain.StakeMonitor` interface and works
// as a local stub for testing.
type StakeMonitor struct {
	stakers []*localStaker
}

// NewStakeMonitor creates a new instance of `StakeMonitor` test stub.
func NewStakeMonitor() *StakeMonitor {
	return &StakeMonitor{
		stakers: make([]*localStaker, 0),
	}
}

// StakerFor returns a staker.Staker instance for the given address. Returns an
// error if the address is invalid.
func (lsm *StakeMonitor) StakerFor(address string) (chain.Staker, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	if staker := lsm.findStakerByAddress(address); staker != nil {
		return staker, nil
	}

	newStaker := &localStaker{address: address, stake: big.NewInt(0)}
	lsm.stakers = append(lsm.stakers, newStaker)

	return newStaker, nil
}

// HasMinimumStake checks if the provided address staked enough to become
// a network operator. The minimum stake is an on-chain parameter.
func (lsm *StakeMonitor) HasMinimumStake(address string) (bool, error) {
	if !common.IsHexAddress(address) {
		return false, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	staker := lsm.findStakerByAddress(address)

	return staker.stake.Cmp(minimumStake) >= 0, nil
}

// StakeTokens stakes enough tokens for the provided address to be a network
// operator.
func (lsm *StakeMonitor) StakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	lsm.findStakerByAddress(address).stake = minimumStake

	return nil
}

// UnstakeTokens unstakes all tokens from the provided address so it can no
// longer be a network operator.
func (lsm *StakeMonitor) UnstakeTokens(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("not a valid ethereum address: %v", address)
	}

	lsm.findStakerByAddress(address).stake = big.NewInt(0)

	return nil
}

func (lsm *StakeMonitor) findStakerByAddress(address string) *localStaker {
	for _, staker := range lsm.stakers {
		if staker.address == address {
			return staker
		}
	}
	return nil
}

type localStaker struct {
	address string
	stake   *big.Int
}

func (ls *localStaker) ID() string {
	return ls.address
}

func (ls *localStaker) Stake() (*big.Int, error) {
	return ls.stake, nil
}

func (ls *localStaker) OnStakeChanged(func(newStake *big.Int)) {
	// Do nothing for now.
}
