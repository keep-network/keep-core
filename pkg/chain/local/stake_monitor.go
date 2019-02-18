package local

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// StakeMonitor implements `chain.StakeMonitor` interface and works
// as a local stub for testing.
type StakeMonitor struct {
	minimumStake *big.Int
	stakers      []*localStaker
}

// NewStakeMonitor creates a new instance of `StakeMonitor` test stub.
func NewStakeMonitor(minimumStake *big.Int) *StakeMonitor {
	return &StakeMonitor{
		minimumStake: minimumStake,
		stakers:      make([]*localStaker, 0),
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

	newStaker := &localStaker{
		address: address,
		stake:   big.NewInt(0),
	}
	lsm.stakers = append(lsm.stakers, newStaker)

	return newStaker, nil
}

func (lsm *StakeMonitor) findStakerByAddress(address string) *localStaker {
	for _, staker := range lsm.stakers {
		if staker.address == address {
			return staker
		}
	}
	return nil
}

// HasMinimumStake checks if the provided address staked enough to become
// a network operator. The minimum stake is an on-chain parameter.
func (lsm *StakeMonitor) HasMinimumStake(address string) (bool, error) {
	staker, err := lsm.StakerFor(address)
	if err != nil {
		return false, err
	}

	stake, err := staker.Stake()
	if err != nil {
		return false, err
	}

	return stake.Cmp(lsm.minimumStake) >= 0, nil
}

// StakeTokens stakes enough tokens for the provided address to be a network
// operator. It stakes `5 * minimumStake` by default.
func (lsm *StakeMonitor) StakeTokens(address string) error {
	staker, err := lsm.StakerFor(address)
	if err != nil {
		return err
	}

	stakerLocal, ok := staker.(*localStaker)
	if !ok {
		return fmt.Errorf("invalid type of staker")
	}

	stakerLocal.stake = new(big.Int).Mul(big.NewInt(5), lsm.minimumStake)

	return nil
}

// UnstakeTokens unstakes all tokens from the provided address so it can no
// longer be a network operator.
func (lsm *StakeMonitor) UnstakeTokens(address string) error {
	staker, err := lsm.StakerFor(address)
	if err != nil {
		return err
	}

	stakerLocal, ok := staker.(*localStaker)
	if !ok {
		return fmt.Errorf("invalid type of staker")
	}

	stakerLocal.stake = big.NewInt(0)

	return nil
}

type localStaker struct {
	address string
	stake   *big.Int
}

func (ls *localStaker) ID() relaychain.StakerAddress {
	return []byte(ls.address)
}

func (ls *localStaker) Stake() (*big.Int, error) {
	return ls.stake, nil
}

func (ls *localStaker) OnStakeChanged(func(newStake *big.Int)) {
	// Do nothing for now.
}
