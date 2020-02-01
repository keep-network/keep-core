package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
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

// Staker returns an instance for the given address that allows insight into a
// staker's stake on Ethereum.
func (esm *ethereumStakeMonitor) StakerFor(address string) (chain.Staker, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("not a valid ethereum address: %v", address)
	}

	return &ethereumStaker{
		address:  address,
		ethereum: esm.config,
	}, nil
}

// StakeMonitor creates and returns a StakeMonitor instance operating on
// Ethereum chain.
func (ec *ethereumChain) StakeMonitor() (chain.StakeMonitor, error) {
	stakeMonitor := &ethereumStakeMonitor{
		config: ec,
	}

	return stakeMonitor, nil
}

func (ec *ethereumChain) BalanceOf(address string) (*big.Int, error) {
	ethereumAddress := common.HexToAddress(address)

	return ec.stakingContract.BalanceOf(ethereumAddress)
}

type ethereumStaker struct {
	address  string
	ethereum *ethereumChain
}

func (es *ethereumStaker) Address() relaychain.StakerAddress {
	return common.HexToAddress(es.address).Bytes()
}

func (es *ethereumStaker) Stake() (*big.Int, error) {
	return es.ethereum.BalanceOf(es.address)
}
