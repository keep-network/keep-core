package ethereum

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumStakeMonitor struct {
	ethereum *ethereumChain
}

func (esm *ethereumStakeMonitor) HasMinimumStake(
	operatorPublicKey *operator.PublicKey,
) (bool, error) {
	address, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return false, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	return esm.ethereum.HasMinimumStake(address)
}

func (esm *ethereumStakeMonitor) StakerFor(
	operatorPublicKey *operator.PublicKey,
) (chain.Staker, error) {
	address, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	return &ethereumStaker{
		address:  address.Hex(),
		ethereum: esm.ethereum,
	}, nil
}

func (ec *ethereumChain) StakeMonitor() (chain.StakeMonitor, error) {
	stakeMonitor := &ethereumStakeMonitor{
		ethereum: ec,
	}

	return stakeMonitor, nil
}

type ethereumStaker struct {
	address  string
	ethereum *ethereumChain
}

func (es *ethereumStaker) Address() relaychain.StakerAddress {
	return common.HexToAddress(es.address).Bytes()
}

func (es *ethereumStaker) Stake() (*big.Int, error) {
	return es.ethereum.stakingContract.BalanceOf(common.HexToAddress(es.address))
}
