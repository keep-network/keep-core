// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	ethereumabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/chain/gen/abi"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var tgLogger = log.Logger("keep-contract-TokenGrant")

type TokenGrant struct {
	contract          *abi.TokenGrant
	contractAddress   common.Address
	contractABI       *ethereumabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *ethutil.ErrorResolver
	nonceManager      *ethutil.NonceManager
	miningWaiter      *ethutil.MiningWaiter

	transactionMutex *sync.Mutex
}

func NewTokenGrant(
	contractAddress common.Address,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethutil.NonceManager,
	miningWaiter *ethutil.MiningWaiter,
	transactionMutex *sync.Mutex,
) (*TokenGrant, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	transactorOptions := bind.NewKeyedTransactor(
		accountKey.PrivateKey,
	)

	randomBeaconContract, err := abi.NewTokenGrant(
		contractAddress,
		backend,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract at address: %s [%v]",
			contractAddress.String(),
			err,
		)
	}

	contractABI, err := ethereumabi.JSON(strings.NewReader(abi.TokenGrantABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TokenGrant{
		contract:          randomBeaconContract,
		contractAddress:   contractAddress,
		contractABI:       &contractABI,
		caller:            backend,
		transactor:        backend,
		callerOptions:     callerOptions,
		transactorOptions: transactorOptions,
		errorResolver:     ethutil.NewErrorResolver(backend, &contractABI, &contractAddress),
		nonceManager:      nonceManager,
		miningWaiter:      miningWaiter,
		transactionMutex:  transactionMutex,
	}, nil
}

// ----- Non-const Methods ------

// Transaction submission.
func (tg *TokenGrant) Stake(
	_id *big.Int,
	_stakingContract common.Address,
	_amount *big.Int,
	_extraData []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction stake",
		"params: ",
		fmt.Sprint(
			_id,
			_stakingContract,
			_amount,
			_extraData,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.Stake(
		transactorOptions,
		_id,
		_stakingContract,
		_amount,
		_extraData,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"stake",
			_id,
			_stakingContract,
			_amount,
			_extraData,
		)
	}

	tgLogger.Debugf(
		"submitted transaction stake with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Stake(
				transactorOptions,
				_id,
				_stakingContract,
				_amount,
				_extraData,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"stake",
					_id,
					_stakingContract,
					_amount,
					_extraData,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallStake(
	_id *big.Int,
	_stakingContract common.Address,
	_amount *big.Int,
	_extraData []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"stake",
		&result,
		_id,
		_stakingContract,
		_amount,
		_extraData,
	)

	return err
}

func (tg *TokenGrant) StakeGasEstimate(
	_id *big.Int,
	_stakingContract common.Address,
	_amount *big.Int,
	_extraData []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"stake",
		tg.contractABI,
		tg.transactor,
		_id,
		_stakingContract,
		_amount,
		_extraData,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) UndelegateRevoked(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction undelegateRevoked",
		"params: ",
		fmt.Sprint(
			_operator,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.UndelegateRevoked(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"undelegateRevoked",
			_operator,
		)
	}

	tgLogger.Debugf(
		"submitted transaction undelegateRevoked with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.UndelegateRevoked(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"undelegateRevoked",
					_operator,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallUndelegateRevoked(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"undelegateRevoked",
		&result,
		_operator,
	)

	return err
}

func (tg *TokenGrant) UndelegateRevokedGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"undelegateRevoked",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) WithdrawRevoked(
	_id *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction withdrawRevoked",
		"params: ",
		fmt.Sprint(
			_id,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.WithdrawRevoked(
		transactorOptions,
		_id,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"withdrawRevoked",
			_id,
		)
	}

	tgLogger.Debugf(
		"submitted transaction withdrawRevoked with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.WithdrawRevoked(
				transactorOptions,
				_id,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"withdrawRevoked",
					_id,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallWithdrawRevoked(
	_id *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"withdrawRevoked",
		&result,
		_id,
	)

	return err
}

func (tg *TokenGrant) WithdrawRevokedGasEstimate(
	_id *big.Int,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"withdrawRevoked",
		tg.contractABI,
		tg.transactor,
		_id,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) Revoke(
	_id *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction revoke",
		"params: ",
		fmt.Sprint(
			_id,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.Revoke(
		transactorOptions,
		_id,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"revoke",
			_id,
		)
	}

	tgLogger.Debugf(
		"submitted transaction revoke with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Revoke(
				transactorOptions,
				_id,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"revoke",
					_id,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallRevoke(
	_id *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"revoke",
		&result,
		_id,
	)

	return err
}

func (tg *TokenGrant) RevokeGasEstimate(
	_id *big.Int,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"revoke",
		tg.contractABI,
		tg.transactor,
		_id,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) CancelStake(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction cancelStake",
		"params: ",
		fmt.Sprint(
			_operator,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.CancelStake(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"cancelStake",
			_operator,
		)
	}

	tgLogger.Debugf(
		"submitted transaction cancelStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.CancelStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"cancelStake",
					_operator,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallCancelStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"cancelStake",
		&result,
		_operator,
	)

	return err
}

func (tg *TokenGrant) CancelStakeGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"cancelStake",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) Withdraw(
	_id *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction withdraw",
		"params: ",
		fmt.Sprint(
			_id,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.Withdraw(
		transactorOptions,
		_id,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"withdraw",
			_id,
		)
	}

	tgLogger.Debugf(
		"submitted transaction withdraw with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Withdraw(
				transactorOptions,
				_id,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"withdraw",
					_id,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallWithdraw(
	_id *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"withdraw",
		&result,
		_id,
	)

	return err
}

func (tg *TokenGrant) WithdrawGasEstimate(
	_id *big.Int,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"withdraw",
		tg.contractABI,
		tg.transactor,
		_id,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) AuthorizeStakingContract(
	_stakingContract common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction authorizeStakingContract",
		"params: ",
		fmt.Sprint(
			_stakingContract,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.AuthorizeStakingContract(
		transactorOptions,
		_stakingContract,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"authorizeStakingContract",
			_stakingContract,
		)
	}

	tgLogger.Debugf(
		"submitted transaction authorizeStakingContract with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.AuthorizeStakingContract(
				transactorOptions,
				_stakingContract,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"authorizeStakingContract",
					_stakingContract,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallAuthorizeStakingContract(
	_stakingContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"authorizeStakingContract",
		&result,
		_stakingContract,
	)

	return err
}

func (tg *TokenGrant) AuthorizeStakingContractGasEstimate(
	_stakingContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"authorizeStakingContract",
		tg.contractABI,
		tg.transactor,
		_stakingContract,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) RecoverStake(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction recoverStake",
		"params: ",
		fmt.Sprint(
			_operator,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.RecoverStake(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"recoverStake",
			_operator,
		)
	}

	tgLogger.Debugf(
		"submitted transaction recoverStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.RecoverStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"recoverStake",
					_operator,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallRecoverStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"recoverStake",
		&result,
		_operator,
	)

	return err
}

func (tg *TokenGrant) RecoverStakeGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"recoverStake",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) ReceiveApproval(
	_from common.Address,
	_amount *big.Int,
	_token common.Address,
	_extraData []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction receiveApproval",
		"params: ",
		fmt.Sprint(
			_from,
			_amount,
			_token,
			_extraData,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.ReceiveApproval(
		transactorOptions,
		_from,
		_amount,
		_token,
		_extraData,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"receiveApproval",
			_from,
			_amount,
			_token,
			_extraData,
		)
	}

	tgLogger.Debugf(
		"submitted transaction receiveApproval with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.ReceiveApproval(
				transactorOptions,
				_from,
				_amount,
				_token,
				_extraData,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"receiveApproval",
					_from,
					_amount,
					_token,
					_extraData,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallReceiveApproval(
	_from common.Address,
	_amount *big.Int,
	_token common.Address,
	_extraData []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"receiveApproval",
		&result,
		_from,
		_amount,
		_token,
		_extraData,
	)

	return err
}

func (tg *TokenGrant) ReceiveApprovalGasEstimate(
	_from common.Address,
	_amount *big.Int,
	_token common.Address,
	_extraData []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"receiveApproval",
		tg.contractABI,
		tg.transactor,
		_from,
		_amount,
		_token,
		_extraData,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) Undelegate(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction undelegate",
		"params: ",
		fmt.Sprint(
			_operator,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.Undelegate(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"undelegate",
			_operator,
		)
	}

	tgLogger.Debugf(
		"submitted transaction undelegate with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Undelegate(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"undelegate",
					_operator,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallUndelegate(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"undelegate",
		&result,
		_operator,
	)

	return err
}

func (tg *TokenGrant) UndelegateGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"undelegate",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) CancelRevokedStake(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tgLogger.Debug(
		"submitting transaction cancelRevokedStake",
		"params: ",
		fmt.Sprint(
			_operator,
		),
	)

	tg.transactionMutex.Lock()
	defer tg.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tg.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tg.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tg.contract.CancelRevokedStake(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, tg.errorResolver.ResolveError(
			err,
			tg.transactorOptions.From,
			nil,
			"cancelRevokedStake",
			_operator,
		)
	}

	tgLogger.Debugf(
		"submitted transaction cancelRevokedStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	go tg.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.CancelRevokedStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return transaction, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"cancelRevokedStake",
					_operator,
				)
			}

			return transaction, nil
		},
	)

	tg.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tg *TokenGrant) CallCancelRevokedStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tg.transactorOptions.From,
		blockNumber, nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"cancelRevokedStake",
		&result,
		_operator,
	)

	return err
}

func (tg *TokenGrant) CancelRevokedStakeGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"cancelRevokedStake",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// ----- Const Methods ------

type GrantStakeDetails struct {
	GrantId         *big.Int
	Amount          *big.Int
	StakingContract common.Address
}

func (tg *TokenGrant) GetGrantStakeDetails(
	operator common.Address,
) (GrantStakeDetails, error) {
	var result GrantStakeDetails
	result, err := tg.contract.GetGrantStakeDetails(
		tg.callerOptions,
		operator,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"getGrantStakeDetails",
			operator,
		)
	}

	return result, err
}

func (tg *TokenGrant) GetGrantStakeDetailsAtBlock(
	operator common.Address,
	blockNumber *big.Int,
) (GrantStakeDetails, error) {
	var result GrantStakeDetails

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"getGrantStakeDetails",
		&result,
		operator,
	)

	return result, err
}

type grants struct {
	GrantManager     common.Address
	Grantee          common.Address
	RevokedAt        *big.Int
	RevokedAmount    *big.Int
	RevokedWithdrawn *big.Int
	Revocable        bool
	Amount           *big.Int
	Duration         *big.Int
	Start            *big.Int
	Cliff            *big.Int
	Withdrawn        *big.Int
	Staked           *big.Int
	StakingPolicy    common.Address
}

func (tg *TokenGrant) Grants(
	arg0 *big.Int,
) (grants, error) {
	var result grants
	result, err := tg.contract.Grants(
		tg.callerOptions,
		arg0,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"grants",
			arg0,
		)
	}

	return result, err
}

func (tg *TokenGrant) GrantsAtBlock(
	arg0 *big.Int,
	blockNumber *big.Int,
) (grants, error) {
	var result grants

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"grants",
		&result,
		arg0,
	)

	return result, err
}

func (tg *TokenGrant) NumGrants() (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.NumGrants(
		tg.callerOptions,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"numGrants",
		)
	}

	return result, err
}

func (tg *TokenGrant) NumGrantsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"numGrants",
		&result,
	)

	return result, err
}

func (tg *TokenGrant) GetGranteeOperators(
	grantee common.Address,
) ([]common.Address, error) {
	var result []common.Address
	result, err := tg.contract.GetGranteeOperators(
		tg.callerOptions,
		grantee,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"getGranteeOperators",
			grantee,
		)
	}

	return result, err
}

func (tg *TokenGrant) GetGranteeOperatorsAtBlock(
	grantee common.Address,
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"getGranteeOperators",
		&result,
		grantee,
	)

	return result, err
}

func (tg *TokenGrant) AvailableToStake(
	_grantId *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.AvailableToStake(
		tg.callerOptions,
		_grantId,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"availableToStake",
			_grantId,
		)
	}

	return result, err
}

func (tg *TokenGrant) AvailableToStakeAtBlock(
	_grantId *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"availableToStake",
		&result,
		_grantId,
	)

	return result, err
}

func (tg *TokenGrant) BalanceOf(
	_owner common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.BalanceOf(
		tg.callerOptions,
		_owner,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"balanceOf",
			_owner,
		)
	}

	return result, err
}

func (tg *TokenGrant) BalanceOfAtBlock(
	_owner common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"balanceOf",
		&result,
		_owner,
	)

	return result, err
}

func (tg *TokenGrant) Token() (common.Address, error) {
	var result common.Address
	result, err := tg.contract.Token(
		tg.callerOptions,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"token",
		)
	}

	return result, err
}

func (tg *TokenGrant) TokenAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"token",
		&result,
	)

	return result, err
}

func (tg *TokenGrant) Withdrawable(
	_id *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.Withdrawable(
		tg.callerOptions,
		_id,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"withdrawable",
			_id,
		)
	}

	return result, err
}

func (tg *TokenGrant) WithdrawableAtBlock(
	_id *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"withdrawable",
		&result,
		_id,
	)

	return result, err
}

func (tg *TokenGrant) GetGrants(
	_granteeOrGrantManager common.Address,
) ([]*big.Int, error) {
	var result []*big.Int
	result, err := tg.contract.GetGrants(
		tg.callerOptions,
		_granteeOrGrantManager,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"getGrants",
			_granteeOrGrantManager,
		)
	}

	return result, err
}

func (tg *TokenGrant) GetGrantsAtBlock(
	_granteeOrGrantManager common.Address,
	blockNumber *big.Int,
) ([]*big.Int, error) {
	var result []*big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"getGrants",
		&result,
		_granteeOrGrantManager,
	)

	return result, err
}

func (tg *TokenGrant) GrantIndices(
	arg0 common.Address,
	arg1 *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.GrantIndices(
		tg.callerOptions,
		arg0,
		arg1,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"grantIndices",
			arg0,
			arg1,
		)
	}

	return result, err
}

func (tg *TokenGrant) GrantIndicesAtBlock(
	arg0 common.Address,
	arg1 *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"grantIndices",
		&result,
		arg0,
		arg1,
	)

	return result, err
}

func (tg *TokenGrant) GranteesToOperators(
	arg0 common.Address,
	arg1 *big.Int,
) (common.Address, error) {
	var result common.Address
	result, err := tg.contract.GranteesToOperators(
		tg.callerOptions,
		arg0,
		arg1,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"granteesToOperators",
			arg0,
			arg1,
		)
	}

	return result, err
}

func (tg *TokenGrant) GranteesToOperatorsAtBlock(
	arg0 common.Address,
	arg1 *big.Int,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"granteesToOperators",
		&result,
		arg0,
		arg1,
	)

	return result, err
}

type Grant struct {
	Amount        *big.Int
	Withdrawn     *big.Int
	Staked        *big.Int
	RevokedAmount *big.Int
	RevokedAt     *big.Int
	Grantee       common.Address
}

func (tg *TokenGrant) GetGrant(
	_id *big.Int,
) (Grant, error) {
	var result Grant
	result, err := tg.contract.GetGrant(
		tg.callerOptions,
		_id,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"getGrant",
			_id,
		)
	}

	return result, err
}

func (tg *TokenGrant) GetGrantAtBlock(
	_id *big.Int,
	blockNumber *big.Int,
) (Grant, error) {
	var result Grant

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"getGrant",
		&result,
		_id,
	)

	return result, err
}

func (tg *TokenGrant) StakeBalanceOf(
	_address common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.StakeBalanceOf(
		tg.callerOptions,
		_address,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"stakeBalanceOf",
			_address,
		)
	}

	return result, err
}

func (tg *TokenGrant) StakeBalanceOfAtBlock(
	_address common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"stakeBalanceOf",
		&result,
		_address,
	)

	return result, err
}

type GrantUnlockingSchedule struct {
	GrantManager common.Address
	Duration     *big.Int
	Start        *big.Int
	Cliff        *big.Int
	Policy       common.Address
}

func (tg *TokenGrant) GetGrantUnlockingSchedule(
	_id *big.Int,
) (GrantUnlockingSchedule, error) {
	var result GrantUnlockingSchedule
	result, err := tg.contract.GetGrantUnlockingSchedule(
		tg.callerOptions,
		_id,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"getGrantUnlockingSchedule",
			_id,
		)
	}

	return result, err
}

func (tg *TokenGrant) GetGrantUnlockingScheduleAtBlock(
	_id *big.Int,
	blockNumber *big.Int,
) (GrantUnlockingSchedule, error) {
	var result GrantUnlockingSchedule

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"getGrantUnlockingSchedule",
		&result,
		_id,
	)

	return result, err
}

func (tg *TokenGrant) UnlockedAmount(
	_id *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.UnlockedAmount(
		tg.callerOptions,
		_id,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"unlockedAmount",
			_id,
		)
	}

	return result, err
}

func (tg *TokenGrant) UnlockedAmountAtBlock(
	_id *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"unlockedAmount",
		&result,
		_id,
	)

	return result, err
}

func (tg *TokenGrant) Balances(
	arg0 common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := tg.contract.Balances(
		tg.callerOptions,
		arg0,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"balances",
			arg0,
		)
	}

	return result, err
}

func (tg *TokenGrant) BalancesAtBlock(
	arg0 common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"balances",
		&result,
		arg0,
	)

	return result, err
}

func (tg *TokenGrant) GrantStakes(
	arg0 common.Address,
) (common.Address, error) {
	var result common.Address
	result, err := tg.contract.GrantStakes(
		tg.callerOptions,
		arg0,
	)

	if err != nil {
		return result, tg.errorResolver.ResolveError(
			err,
			tg.callerOptions.From,
			nil,
			"grantStakes",
			arg0,
		)
	}

	return result, err
}

func (tg *TokenGrant) GrantStakesAtBlock(
	arg0 common.Address,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tg.callerOptions.From,
		blockNumber,
		nil,
		tg.contractABI,
		tg.caller,
		tg.errorResolver,
		tg.contractAddress,
		"grantStakes",
		&result,
		arg0,
	)

	return result, err
}

// ------ Events -------

type tokenGrantStakingContractAuthorizedFunc func(
	GrantManager common.Address,
	StakingContract common.Address,
	blockNumber uint64,
)

func (tg *TokenGrant) WatchStakingContractAuthorized(
	success tokenGrantStakingContractAuthorizedFunc,
	fail func(err error) error,
	grantManagerFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tg.subscribeStakingContractAuthorized(
			success,
			failCallback,
			grantManagerFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tgLogger.Warning(
					"subscription to event StakingContractAuthorized terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tg *TokenGrant) subscribeStakingContractAuthorized(
	success tokenGrantStakingContractAuthorizedFunc,
	fail func(err error) error,
	grantManagerFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenGrantStakingContractAuthorized)
	eventSubscription, err := tg.contract.WatchStakingContractAuthorized(
		nil,
		eventChan,
		grantManagerFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for StakingContractAuthorized events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.GrantManager,
					event.StakingContract,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tokenGrantTokenGrantCreatedFunc func(
	Id *big.Int,
	blockNumber uint64,
)

func (tg *TokenGrant) WatchTokenGrantCreated(
	success tokenGrantTokenGrantCreatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tg.subscribeTokenGrantCreated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tgLogger.Warning(
					"subscription to event TokenGrantCreated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tg *TokenGrant) subscribeTokenGrantCreated(
	success tokenGrantTokenGrantCreatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenGrantTokenGrantCreated)
	eventSubscription, err := tg.contract.WatchTokenGrantCreated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokenGrantCreated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.Id,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tokenGrantTokenGrantRevokedFunc func(
	Id *big.Int,
	blockNumber uint64,
)

func (tg *TokenGrant) WatchTokenGrantRevoked(
	success tokenGrantTokenGrantRevokedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tg.subscribeTokenGrantRevoked(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tgLogger.Warning(
					"subscription to event TokenGrantRevoked terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tg *TokenGrant) subscribeTokenGrantRevoked(
	success tokenGrantTokenGrantRevokedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenGrantTokenGrantRevoked)
	eventSubscription, err := tg.contract.WatchTokenGrantRevoked(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokenGrantRevoked events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.Id,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tokenGrantTokenGrantStakedFunc func(
	GrantId *big.Int,
	Amount *big.Int,
	Operator common.Address,
	blockNumber uint64,
)

func (tg *TokenGrant) WatchTokenGrantStaked(
	success tokenGrantTokenGrantStakedFunc,
	fail func(err error) error,
	grantIdFilter []*big.Int,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tg.subscribeTokenGrantStaked(
			success,
			failCallback,
			grantIdFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tgLogger.Warning(
					"subscription to event TokenGrantStaked terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tg *TokenGrant) subscribeTokenGrantStaked(
	success tokenGrantTokenGrantStakedFunc,
	fail func(err error) error,
	grantIdFilter []*big.Int,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenGrantTokenGrantStaked)
	eventSubscription, err := tg.contract.WatchTokenGrantStaked(
		nil,
		eventChan,
		grantIdFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokenGrantStaked events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.GrantId,
					event.Amount,
					event.Operator,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tokenGrantTokenGrantWithdrawnFunc func(
	GrantId *big.Int,
	Amount *big.Int,
	blockNumber uint64,
)

func (tg *TokenGrant) WatchTokenGrantWithdrawn(
	success tokenGrantTokenGrantWithdrawnFunc,
	fail func(err error) error,
	grantIdFilter []*big.Int,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tg.subscribeTokenGrantWithdrawn(
			success,
			failCallback,
			grantIdFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tgLogger.Warning(
					"subscription to event TokenGrantWithdrawn terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tg *TokenGrant) subscribeTokenGrantWithdrawn(
	success tokenGrantTokenGrantWithdrawnFunc,
	fail func(err error) error,
	grantIdFilter []*big.Int,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenGrantTokenGrantWithdrawn)
	eventSubscription, err := tg.contract.WatchTokenGrantWithdrawn(
		nil,
		eventChan,
		grantIdFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokenGrantWithdrawn events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.GrantId,
					event.Amount,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}
