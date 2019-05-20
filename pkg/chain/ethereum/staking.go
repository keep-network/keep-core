package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
)

// staking contains connection information for interface to the staking proxy
// contract.
type staking struct {
	caller          *abi.StakingProxyCaller
	callerOpts      *bind.CallOpts
	transactor      *abi.StakingProxyTransactor
	transactorOpts  *bind.TransactOpts
	contract        *abi.StakingProxy
	contractAddress common.Address
}

// NewTokenStaking creates the necessary connections and configurations for
// accessing the contract.
func newStaking(chainConfig *ethereumChain) (*staking, error) {
	contractAddressHex, exists := chainConfig.config.ContractAddresses["Staking"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'Staking' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	stakingTransactor, err := abi.NewStakingProxyTransactor(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a StakingProxyTransactor contract: [%v]",
			err,
		)
	}

	if chainConfig.accountKey == nil {
		key, err := ethutil.DecryptKeyFile(
			chainConfig.config.Account.KeyFile,
			chainConfig.config.Account.KeyFilePassword,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read KeyFile: %s: [%v]",
				chainConfig.config.Account.KeyFile,
				err,
			)
		}
		chainConfig.accountKey = key
	}

	optsTransactor := bind.NewKeyedTransactor(
		chainConfig.accountKey.PrivateKey,
	)

	stakingCaller, err := abi.NewStakingProxyCaller(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a StakingProxyCaller contract: [%v]",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		From: contractAddress,
	}

	stakingProxyContract, err := abi.NewStakingProxy(
		contractAddress,
		chainConfig.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract at address: %s [%v]",
			contractAddressHex,
			err,
		)
	}

	return &staking{
		transactor:      stakingTransactor,
		transactorOpts:  optsTransactor,
		caller:          stakingCaller,
		callerOpts:      optsCaller,
		contract:        stakingProxyContract,
		contractAddress: contractAddress,
	}, nil
}

// BalanceOf returns a big.Int containing the currently staked balance of
// the given Ethereum address in this token staking contract.
func (s *staking) BalanceOf(address common.Address) (*big.Int, error) {
	return s.caller.BalanceOf(s.callerOpts, address)
}

// WatchUnstakedFor
func (s *staking) WatchUnstakedFor(
	address common.Address,
	success func(common.Address, *big.Int),
	fail errorCallback,
) error {
	eventChan := make(chan *abi.StakingProxyUnstaked)
	eventSubscription, err := s.contract.WatchUnstaked(nil, eventChan, []common.Address{address})
	if err != nil {
		close(eventChan)
		return fmt.Errorf(
			"error creating watch for Unstaked event: [%v]",
			err,
		)
	}
	go func() {
		defer close(eventChan)
		defer eventSubscription.Unsubscribe()
		for {
			select {
			case event := <-eventChan:
				success(
					event.Staker,
					event.Amount,
				)

			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()
	return nil
}
