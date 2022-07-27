package local_v1

import (
	"encoding/hex"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"reflect"
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

// StakerFor returns a staker.Staker instance for the given operator public key.
func (lsm *StakeMonitor) stakerFor(
	operatorPublicKey *operator.PublicKey,
) (*localStaker, error) {
	if staker := lsm.findStakerByPublicKey(operatorPublicKey); staker != nil {
		return staker, nil
	}

	newStaker := &localStaker{
		publicKey: operatorPublicKey,
		stake:     big.NewInt(0),
	}
	lsm.stakers = append(lsm.stakers, newStaker)

	return newStaker, nil
}

func (lsm *StakeMonitor) findStakerByPublicKey(
	publicKey *operator.PublicKey,
) *localStaker {
	for _, staker := range lsm.stakers {
		if reflect.DeepEqual(staker.publicKey, publicKey) {
			return staker
		}
	}
	return nil
}

// HasMinimumStake checks if the provided public key staked enough to become
// a network operator. The minimum stake is an on-chain parameter.
func (lsm *StakeMonitor) HasMinimumStake(
	operatorPublicKey *operator.PublicKey,
) (bool, error) {
	staker, err := lsm.stakerFor(operatorPublicKey)
	if err != nil {
		return false, err
	}

	stake, err := staker.Stake()
	if err != nil {
		return false, err
	}

	return stake.Cmp(lsm.minimumStake) >= 0, nil
}

// StakeTokens stakes enough tokens for the provided public key to be a network
// operator. It stakes `5 * minimumStake` by default.
func (lsm *StakeMonitor) StakeTokens(operatorPublicKey *operator.PublicKey) error {
	staker, err := lsm.stakerFor(operatorPublicKey)
	if err != nil {
		return err
	}

	staker.stake = new(big.Int).Mul(big.NewInt(5), lsm.minimumStake)

	return nil
}

// UnstakeTokens unstakes all tokens from the provided public key so it can no
// longer be a network operator.
func (lsm *StakeMonitor) UnstakeTokens(operatorPublicKey *operator.PublicKey) error {
	staker, err := lsm.stakerFor(operatorPublicKey)
	if err != nil {
		return err
	}

	staker.stake = big.NewInt(0)

	return nil
}

type localStaker struct {
	publicKey *operator.PublicKey
	stake     *big.Int
}

func (ls *localStaker) Address() chain.Address {
	return chain.Address(hex.EncodeToString(operator.MarshalCompressed(ls.publicKey)))
}

func (ls *localStaker) Stake() (*big.Int, error) {
	return ls.stake, nil
}