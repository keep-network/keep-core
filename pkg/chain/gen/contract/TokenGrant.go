// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	hostchainabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"github.com/ipfs/go-log"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/chain/ethlike"
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
	contractABI       *hostchainabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *chainutil.ErrorResolver
	nonceManager      *ethlike.NonceManager
	miningWaiter      *ethlike.MiningWaiter
	blockCounter      *ethlike.BlockCounter

	transactionMutex *sync.Mutex
}

func NewTokenGrant(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *ethlike.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*TokenGrant, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	// FIXME Switch to bind.NewKeyedTransactorWithChainID when
	// FIXME celo-org/celo-blockchain merges in changes from upstream
	// FIXME ethereum/go-ethereum beyond v1.9.25.
	transactorOptions, err := chainutil.NewKeyedTransactorWithChainID(
		accountKey.PrivateKey,
		chainId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate transactor: [%v]", err)
	}

	contract, err := abi.NewTokenGrant(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.TokenGrantABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TokenGrant{
		contract:          contract,
		contractAddress:   contractAddress,
		contractABI:       &contractABI,
		caller:            backend,
		transactor:        backend,
		callerOptions:     callerOptions,
		transactorOptions: transactorOptions,
		errorResolver:     chainutil.NewErrorResolver(backend, &contractABI, &contractAddress),
		nonceManager:      nonceManager,
		miningWaiter:      miningWaiter,
		blockCounter:      blockCounter,
		transactionMutex:  transactionMutex,
	}, nil
}

// ----- Non-const Methods ------

// Transaction submission.
func (tg *TokenGrant) AuthorizeStakingContract(
	_stakingContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction authorizeStakingContract with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.AuthorizeStakingContract(
				transactorOptions,
				_stakingContract,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"authorizeStakingContract",
					_stakingContract,
				)
			}

			tgLogger.Infof(
				"submitted transaction authorizeStakingContract with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) CancelRevokedStake(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction cancelRevokedStake with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.CancelRevokedStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"cancelRevokedStake",
					_operator,
				)
			}

			tgLogger.Infof(
				"submitted transaction cancelRevokedStake with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"cancelRevokedStake",
		tg.contractABI,
		tg.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (tg *TokenGrant) CancelStake(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction cancelStake with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.CancelStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"cancelStake",
					_operator,
				)
			}

			tgLogger.Infof(
				"submitted transaction cancelStake with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) ReceiveApproval(
	_from common.Address,
	_amount *big.Int,
	_token common.Address,
	_extraData []uint8,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction receiveApproval with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
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
				return nil, tg.errorResolver.ResolveError(
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

			tgLogger.Infof(
				"submitted transaction receiveApproval with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) RecoverStake(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction recoverStake with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.RecoverStake(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"recoverStake",
					_operator,
				)
			}

			tgLogger.Infof(
				"submitted transaction recoverStake with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) Revoke(
	_id *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction revoke with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Revoke(
				transactorOptions,
				_id,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"revoke",
					_id,
				)
			}

			tgLogger.Infof(
				"submitted transaction revoke with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) Stake(
	_id *big.Int,
	_stakingContract common.Address,
	_amount *big.Int,
	_extraData []uint8,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction stake with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
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
				return nil, tg.errorResolver.ResolveError(
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

			tgLogger.Infof(
				"submitted transaction stake with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) Undelegate(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction undelegate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Undelegate(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"undelegate",
					_operator,
				)
			}

			tgLogger.Infof(
				"submitted transaction undelegate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) UndelegateRevoked(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction undelegateRevoked with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.UndelegateRevoked(
				transactorOptions,
				_operator,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"undelegateRevoked",
					_operator,
				)
			}

			tgLogger.Infof(
				"submitted transaction undelegateRevoked with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) Withdraw(
	_id *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction withdraw with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.Withdraw(
				transactorOptions,
				_id,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"withdraw",
					_id,
				)
			}

			tgLogger.Infof(
				"submitted transaction withdraw with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tg *TokenGrant) WithdrawRevoked(
	_id *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
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

	tgLogger.Infof(
		"submitted transaction withdrawRevoked with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tg.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tg.contract.WithdrawRevoked(
				transactorOptions,
				_id,
			)
			if err != nil {
				return nil, tg.errorResolver.ResolveError(
					err,
					tg.transactorOptions.From,
					nil,
					"withdrawRevoked",
					_id,
				)
			}

			tgLogger.Infof(
				"submitted transaction withdrawRevoked with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tg.callerOptions.From,
		tg.contractAddress,
		"withdrawRevoked",
		tg.contractABI,
		tg.transactor,
		_id,
	)

	return result, err
}

// ----- Const Methods ------

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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

// ------ Events -------

func (tg *TokenGrant) StakingContractAuthorized(
	opts *ethlike.SubscribeOpts,
	grantManagerFilter []common.Address,
) *TgStakingContractAuthorizedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TgStakingContractAuthorizedSubscription{
		tg,
		opts,
		grantManagerFilter,
	}
}

type TgStakingContractAuthorizedSubscription struct {
	contract           *TokenGrant
	opts               *ethlike.SubscribeOpts
	grantManagerFilter []common.Address
}

type tokenGrantStakingContractAuthorizedFunc func(
	GrantManager common.Address,
	StakingContract common.Address,
	blockNumber uint64,
)

func (scas *TgStakingContractAuthorizedSubscription) OnEvent(
	handler tokenGrantStakingContractAuthorizedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenGrantStakingContractAuthorized)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GrantManager,
					event.StakingContract,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := scas.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (scas *TgStakingContractAuthorizedSubscription) Pipe(
	sink chan *abi.TokenGrantStakingContractAuthorized,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(scas.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := scas.contract.blockCounter.CurrentBlock()
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - scas.opts.PastBlocks

				tgLogger.Infof(
					"subscription monitoring fetching past StakingContractAuthorized events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := scas.contract.PastStakingContractAuthorizedEvents(
					fromBlock,
					nil,
					scas.grantManagerFilter,
				)
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tgLogger.Infof(
					"subscription monitoring fetched [%v] past StakingContractAuthorized events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := scas.contract.watchStakingContractAuthorized(
		sink,
		scas.grantManagerFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tg *TokenGrant) watchStakingContractAuthorized(
	sink chan *abi.TokenGrantStakingContractAuthorized,
	grantManagerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tg.contract.WatchStakingContractAuthorized(
			&bind.WatchOpts{Context: ctx},
			sink,
			grantManagerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tgLogger.Errorf(
			"subscription to event StakingContractAuthorized had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tgLogger.Errorf(
			"subscription to event StakingContractAuthorized failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

func (tg *TokenGrant) PastStakingContractAuthorizedEvents(
	startBlock uint64,
	endBlock *uint64,
	grantManagerFilter []common.Address,
) ([]*abi.TokenGrantStakingContractAuthorized, error) {
	iterator, err := tg.contract.FilterStakingContractAuthorized(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		grantManagerFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past StakingContractAuthorized events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenGrantStakingContractAuthorized, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tg *TokenGrant) TokenGrantCreated(
	opts *ethlike.SubscribeOpts,
) *TgTokenGrantCreatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TgTokenGrantCreatedSubscription{
		tg,
		opts,
	}
}

type TgTokenGrantCreatedSubscription struct {
	contract *TokenGrant
	opts     *ethlike.SubscribeOpts
}

type tokenGrantTokenGrantCreatedFunc func(
	Id *big.Int,
	blockNumber uint64,
)

func (tgcs *TgTokenGrantCreatedSubscription) OnEvent(
	handler tokenGrantTokenGrantCreatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenGrantTokenGrantCreated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Id,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tgcs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tgcs *TgTokenGrantCreatedSubscription) Pipe(
	sink chan *abi.TokenGrantTokenGrantCreated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tgcs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tgcs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tgcs.opts.PastBlocks

				tgLogger.Infof(
					"subscription monitoring fetching past TokenGrantCreated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tgcs.contract.PastTokenGrantCreatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tgLogger.Infof(
					"subscription monitoring fetched [%v] past TokenGrantCreated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tgcs.contract.watchTokenGrantCreated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tg *TokenGrant) watchTokenGrantCreated(
	sink chan *abi.TokenGrantTokenGrantCreated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tg.contract.WatchTokenGrantCreated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tgLogger.Errorf(
			"subscription to event TokenGrantCreated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tgLogger.Errorf(
			"subscription to event TokenGrantCreated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

func (tg *TokenGrant) PastTokenGrantCreatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TokenGrantTokenGrantCreated, error) {
	iterator, err := tg.contract.FilterTokenGrantCreated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokenGrantCreated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenGrantTokenGrantCreated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tg *TokenGrant) TokenGrantRevoked(
	opts *ethlike.SubscribeOpts,
) *TgTokenGrantRevokedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TgTokenGrantRevokedSubscription{
		tg,
		opts,
	}
}

type TgTokenGrantRevokedSubscription struct {
	contract *TokenGrant
	opts     *ethlike.SubscribeOpts
}

type tokenGrantTokenGrantRevokedFunc func(
	Id *big.Int,
	blockNumber uint64,
)

func (tgrs *TgTokenGrantRevokedSubscription) OnEvent(
	handler tokenGrantTokenGrantRevokedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenGrantTokenGrantRevoked)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Id,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tgrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tgrs *TgTokenGrantRevokedSubscription) Pipe(
	sink chan *abi.TokenGrantTokenGrantRevoked,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tgrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tgrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tgrs.opts.PastBlocks

				tgLogger.Infof(
					"subscription monitoring fetching past TokenGrantRevoked events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tgrs.contract.PastTokenGrantRevokedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tgLogger.Infof(
					"subscription monitoring fetched [%v] past TokenGrantRevoked events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tgrs.contract.watchTokenGrantRevoked(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tg *TokenGrant) watchTokenGrantRevoked(
	sink chan *abi.TokenGrantTokenGrantRevoked,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tg.contract.WatchTokenGrantRevoked(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tgLogger.Errorf(
			"subscription to event TokenGrantRevoked had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tgLogger.Errorf(
			"subscription to event TokenGrantRevoked failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

func (tg *TokenGrant) PastTokenGrantRevokedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TokenGrantTokenGrantRevoked, error) {
	iterator, err := tg.contract.FilterTokenGrantRevoked(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokenGrantRevoked events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenGrantTokenGrantRevoked, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tg *TokenGrant) TokenGrantStaked(
	opts *ethlike.SubscribeOpts,
	grantIdFilter []*big.Int,
) *TgTokenGrantStakedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TgTokenGrantStakedSubscription{
		tg,
		opts,
		grantIdFilter,
	}
}

type TgTokenGrantStakedSubscription struct {
	contract      *TokenGrant
	opts          *ethlike.SubscribeOpts
	grantIdFilter []*big.Int
}

type tokenGrantTokenGrantStakedFunc func(
	GrantId *big.Int,
	Amount *big.Int,
	Operator common.Address,
	blockNumber uint64,
)

func (tgss *TgTokenGrantStakedSubscription) OnEvent(
	handler tokenGrantTokenGrantStakedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenGrantTokenGrantStaked)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GrantId,
					event.Amount,
					event.Operator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tgss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tgss *TgTokenGrantStakedSubscription) Pipe(
	sink chan *abi.TokenGrantTokenGrantStaked,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tgss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tgss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tgss.opts.PastBlocks

				tgLogger.Infof(
					"subscription monitoring fetching past TokenGrantStaked events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tgss.contract.PastTokenGrantStakedEvents(
					fromBlock,
					nil,
					tgss.grantIdFilter,
				)
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tgLogger.Infof(
					"subscription monitoring fetched [%v] past TokenGrantStaked events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tgss.contract.watchTokenGrantStaked(
		sink,
		tgss.grantIdFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tg *TokenGrant) watchTokenGrantStaked(
	sink chan *abi.TokenGrantTokenGrantStaked,
	grantIdFilter []*big.Int,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tg.contract.WatchTokenGrantStaked(
			&bind.WatchOpts{Context: ctx},
			sink,
			grantIdFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tgLogger.Errorf(
			"subscription to event TokenGrantStaked had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tgLogger.Errorf(
			"subscription to event TokenGrantStaked failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

func (tg *TokenGrant) PastTokenGrantStakedEvents(
	startBlock uint64,
	endBlock *uint64,
	grantIdFilter []*big.Int,
) ([]*abi.TokenGrantTokenGrantStaked, error) {
	iterator, err := tg.contract.FilterTokenGrantStaked(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		grantIdFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokenGrantStaked events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenGrantTokenGrantStaked, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tg *TokenGrant) TokenGrantWithdrawn(
	opts *ethlike.SubscribeOpts,
	grantIdFilter []*big.Int,
) *TgTokenGrantWithdrawnSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TgTokenGrantWithdrawnSubscription{
		tg,
		opts,
		grantIdFilter,
	}
}

type TgTokenGrantWithdrawnSubscription struct {
	contract      *TokenGrant
	opts          *ethlike.SubscribeOpts
	grantIdFilter []*big.Int
}

type tokenGrantTokenGrantWithdrawnFunc func(
	GrantId *big.Int,
	Amount *big.Int,
	blockNumber uint64,
)

func (tgws *TgTokenGrantWithdrawnSubscription) OnEvent(
	handler tokenGrantTokenGrantWithdrawnFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenGrantTokenGrantWithdrawn)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GrantId,
					event.Amount,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tgws.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tgws *TgTokenGrantWithdrawnSubscription) Pipe(
	sink chan *abi.TokenGrantTokenGrantWithdrawn,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tgws.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tgws.contract.blockCounter.CurrentBlock()
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tgws.opts.PastBlocks

				tgLogger.Infof(
					"subscription monitoring fetching past TokenGrantWithdrawn events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tgws.contract.PastTokenGrantWithdrawnEvents(
					fromBlock,
					nil,
					tgws.grantIdFilter,
				)
				if err != nil {
					tgLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tgLogger.Infof(
					"subscription monitoring fetched [%v] past TokenGrantWithdrawn events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tgws.contract.watchTokenGrantWithdrawn(
		sink,
		tgws.grantIdFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tg *TokenGrant) watchTokenGrantWithdrawn(
	sink chan *abi.TokenGrantTokenGrantWithdrawn,
	grantIdFilter []*big.Int,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tg.contract.WatchTokenGrantWithdrawn(
			&bind.WatchOpts{Context: ctx},
			sink,
			grantIdFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tgLogger.Errorf(
			"subscription to event TokenGrantWithdrawn had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tgLogger.Errorf(
			"subscription to event TokenGrantWithdrawn failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

func (tg *TokenGrant) PastTokenGrantWithdrawnEvents(
	startBlock uint64,
	endBlock *uint64,
	grantIdFilter []*big.Int,
) ([]*abi.TokenGrantTokenGrantWithdrawn, error) {
	iterator, err := tg.contract.FilterTokenGrantWithdrawn(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		grantIdFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokenGrantWithdrawn events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenGrantTokenGrantWithdrawn, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
