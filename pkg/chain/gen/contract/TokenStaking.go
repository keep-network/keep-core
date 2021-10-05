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
var tsLogger = log.Logger("keep-contract-TokenStaking")

type TokenStaking struct {
	contract          *abi.TokenStaking
	contractAddress   common.Address
	contractABI       *hostchainabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *chainutil.ErrorResolver
	nonceManager      *ethlike.NonceManager
	miningWaiter      *chainutil.MiningWaiter
	blockCounter      *ethlike.BlockCounter

	transactionMutex *sync.Mutex
}

func NewTokenStaking(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*TokenStaking, error) {
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

	contract, err := abi.NewTokenStaking(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.TokenStakingABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TokenStaking{
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
func (ts *TokenStaking) AuthorizeOperatorContract(
	_operator common.Address,
	_operatorContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction authorizeOperatorContract",
		" params: ",
		fmt.Sprint(
			_operator,
			_operatorContract,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.AuthorizeOperatorContract(
		transactorOptions,
		_operator,
		_operatorContract,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"authorizeOperatorContract",
			_operator,
			_operatorContract,
		)
	}

	tsLogger.Infof(
		"submitted transaction authorizeOperatorContract with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.AuthorizeOperatorContract(
				newTransactorOptions,
				_operator,
				_operatorContract,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"authorizeOperatorContract",
					_operator,
					_operatorContract,
				)
			}

			tsLogger.Infof(
				"submitted transaction authorizeOperatorContract with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallAuthorizeOperatorContract(
	_operator common.Address,
	_operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"authorizeOperatorContract",
		&result,
		_operator,
		_operatorContract,
	)

	return err
}

func (ts *TokenStaking) AuthorizeOperatorContractGasEstimate(
	_operator common.Address,
	_operatorContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"authorizeOperatorContract",
		ts.contractABI,
		ts.transactor,
		_operator,
		_operatorContract,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) CancelStake(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction cancelStake",
		" params: ",
		fmt.Sprint(
			_operator,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.CancelStake(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"cancelStake",
			_operator,
		)
	}

	tsLogger.Infof(
		"submitted transaction cancelStake with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.CancelStake(
				newTransactorOptions,
				_operator,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"cancelStake",
					_operator,
				)
			}

			tsLogger.Infof(
				"submitted transaction cancelStake with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallCancelStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"cancelStake",
		&result,
		_operator,
	)

	return err
}

func (ts *TokenStaking) CancelStakeGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"cancelStake",
		ts.contractABI,
		ts.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) ClaimDelegatedAuthority(
	delegatedAuthoritySource common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction claimDelegatedAuthority",
		" params: ",
		fmt.Sprint(
			delegatedAuthoritySource,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.ClaimDelegatedAuthority(
		transactorOptions,
		delegatedAuthoritySource,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"claimDelegatedAuthority",
			delegatedAuthoritySource,
		)
	}

	tsLogger.Infof(
		"submitted transaction claimDelegatedAuthority with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.ClaimDelegatedAuthority(
				newTransactorOptions,
				delegatedAuthoritySource,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"claimDelegatedAuthority",
					delegatedAuthoritySource,
				)
			}

			tsLogger.Infof(
				"submitted transaction claimDelegatedAuthority with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallClaimDelegatedAuthority(
	delegatedAuthoritySource common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"claimDelegatedAuthority",
		&result,
		delegatedAuthoritySource,
	)

	return err
}

func (ts *TokenStaking) ClaimDelegatedAuthorityGasEstimate(
	delegatedAuthoritySource common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"claimDelegatedAuthority",
		ts.contractABI,
		ts.transactor,
		delegatedAuthoritySource,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) CommitTopUp(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction commitTopUp",
		" params: ",
		fmt.Sprint(
			_operator,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.CommitTopUp(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"commitTopUp",
			_operator,
		)
	}

	tsLogger.Infof(
		"submitted transaction commitTopUp with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.CommitTopUp(
				newTransactorOptions,
				_operator,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"commitTopUp",
					_operator,
				)
			}

			tsLogger.Infof(
				"submitted transaction commitTopUp with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallCommitTopUp(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"commitTopUp",
		&result,
		_operator,
	)

	return err
}

func (ts *TokenStaking) CommitTopUpGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"commitTopUp",
		ts.contractABI,
		ts.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) LockStake(
	operator common.Address,
	duration *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction lockStake",
		" params: ",
		fmt.Sprint(
			operator,
			duration,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.LockStake(
		transactorOptions,
		operator,
		duration,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"lockStake",
			operator,
			duration,
		)
	}

	tsLogger.Infof(
		"submitted transaction lockStake with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.LockStake(
				newTransactorOptions,
				operator,
				duration,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"lockStake",
					operator,
					duration,
				)
			}

			tsLogger.Infof(
				"submitted transaction lockStake with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallLockStake(
	operator common.Address,
	duration *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"lockStake",
		&result,
		operator,
		duration,
	)

	return err
}

func (ts *TokenStaking) LockStakeGasEstimate(
	operator common.Address,
	duration *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"lockStake",
		ts.contractABI,
		ts.transactor,
		operator,
		duration,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) ReceiveApproval(
	_from common.Address,
	_value *big.Int,
	_token common.Address,
	_extraData []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction receiveApproval",
		" params: ",
		fmt.Sprint(
			_from,
			_value,
			_token,
			_extraData,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.ReceiveApproval(
		transactorOptions,
		_from,
		_value,
		_token,
		_extraData,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"receiveApproval",
			_from,
			_value,
			_token,
			_extraData,
		)
	}

	tsLogger.Infof(
		"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.ReceiveApproval(
				newTransactorOptions,
				_from,
				_value,
				_token,
				_extraData,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"receiveApproval",
					_from,
					_value,
					_token,
					_extraData,
				)
			}

			tsLogger.Infof(
				"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallReceiveApproval(
	_from common.Address,
	_value *big.Int,
	_token common.Address,
	_extraData []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"receiveApproval",
		&result,
		_from,
		_value,
		_token,
		_extraData,
	)

	return err
}

func (ts *TokenStaking) ReceiveApprovalGasEstimate(
	_from common.Address,
	_value *big.Int,
	_token common.Address,
	_extraData []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"receiveApproval",
		ts.contractABI,
		ts.transactor,
		_from,
		_value,
		_token,
		_extraData,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) RecoverStake(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction recoverStake",
		" params: ",
		fmt.Sprint(
			_operator,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.RecoverStake(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"recoverStake",
			_operator,
		)
	}

	tsLogger.Infof(
		"submitted transaction recoverStake with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.RecoverStake(
				newTransactorOptions,
				_operator,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"recoverStake",
					_operator,
				)
			}

			tsLogger.Infof(
				"submitted transaction recoverStake with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallRecoverStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"recoverStake",
		&result,
		_operator,
	)

	return err
}

func (ts *TokenStaking) RecoverStakeGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"recoverStake",
		ts.contractABI,
		ts.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) ReleaseExpiredLock(
	operator common.Address,
	operatorContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction releaseExpiredLock",
		" params: ",
		fmt.Sprint(
			operator,
			operatorContract,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.ReleaseExpiredLock(
		transactorOptions,
		operator,
		operatorContract,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"releaseExpiredLock",
			operator,
			operatorContract,
		)
	}

	tsLogger.Infof(
		"submitted transaction releaseExpiredLock with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.ReleaseExpiredLock(
				newTransactorOptions,
				operator,
				operatorContract,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"releaseExpiredLock",
					operator,
					operatorContract,
				)
			}

			tsLogger.Infof(
				"submitted transaction releaseExpiredLock with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallReleaseExpiredLock(
	operator common.Address,
	operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"releaseExpiredLock",
		&result,
		operator,
		operatorContract,
	)

	return err
}

func (ts *TokenStaking) ReleaseExpiredLockGasEstimate(
	operator common.Address,
	operatorContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"releaseExpiredLock",
		ts.contractABI,
		ts.transactor,
		operator,
		operatorContract,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) Seize(
	amountToSeize *big.Int,
	rewardMultiplier *big.Int,
	tattletale common.Address,
	misbehavedOperators []common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction seize",
		" params: ",
		fmt.Sprint(
			amountToSeize,
			rewardMultiplier,
			tattletale,
			misbehavedOperators,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.Seize(
		transactorOptions,
		amountToSeize,
		rewardMultiplier,
		tattletale,
		misbehavedOperators,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"seize",
			amountToSeize,
			rewardMultiplier,
			tattletale,
			misbehavedOperators,
		)
	}

	tsLogger.Infof(
		"submitted transaction seize with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.Seize(
				newTransactorOptions,
				amountToSeize,
				rewardMultiplier,
				tattletale,
				misbehavedOperators,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"seize",
					amountToSeize,
					rewardMultiplier,
					tattletale,
					misbehavedOperators,
				)
			}

			tsLogger.Infof(
				"submitted transaction seize with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallSeize(
	amountToSeize *big.Int,
	rewardMultiplier *big.Int,
	tattletale common.Address,
	misbehavedOperators []common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"seize",
		&result,
		amountToSeize,
		rewardMultiplier,
		tattletale,
		misbehavedOperators,
	)

	return err
}

func (ts *TokenStaking) SeizeGasEstimate(
	amountToSeize *big.Int,
	rewardMultiplier *big.Int,
	tattletale common.Address,
	misbehavedOperators []common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"seize",
		ts.contractABI,
		ts.transactor,
		amountToSeize,
		rewardMultiplier,
		tattletale,
		misbehavedOperators,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) Slash(
	amountToSlash *big.Int,
	misbehavedOperators []common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction slash",
		" params: ",
		fmt.Sprint(
			amountToSlash,
			misbehavedOperators,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.Slash(
		transactorOptions,
		amountToSlash,
		misbehavedOperators,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"slash",
			amountToSlash,
			misbehavedOperators,
		)
	}

	tsLogger.Infof(
		"submitted transaction slash with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.Slash(
				newTransactorOptions,
				amountToSlash,
				misbehavedOperators,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"slash",
					amountToSlash,
					misbehavedOperators,
				)
			}

			tsLogger.Infof(
				"submitted transaction slash with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallSlash(
	amountToSlash *big.Int,
	misbehavedOperators []common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"slash",
		&result,
		amountToSlash,
		misbehavedOperators,
	)

	return err
}

func (ts *TokenStaking) SlashGasEstimate(
	amountToSlash *big.Int,
	misbehavedOperators []common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"slash",
		ts.contractABI,
		ts.transactor,
		amountToSlash,
		misbehavedOperators,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) TransferStakeOwnership(
	operator common.Address,
	newOwner common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction transferStakeOwnership",
		" params: ",
		fmt.Sprint(
			operator,
			newOwner,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.TransferStakeOwnership(
		transactorOptions,
		operator,
		newOwner,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"transferStakeOwnership",
			operator,
			newOwner,
		)
	}

	tsLogger.Infof(
		"submitted transaction transferStakeOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.TransferStakeOwnership(
				newTransactorOptions,
				operator,
				newOwner,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"transferStakeOwnership",
					operator,
					newOwner,
				)
			}

			tsLogger.Infof(
				"submitted transaction transferStakeOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallTransferStakeOwnership(
	operator common.Address,
	newOwner common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"transferStakeOwnership",
		&result,
		operator,
		newOwner,
	)

	return err
}

func (ts *TokenStaking) TransferStakeOwnershipGasEstimate(
	operator common.Address,
	newOwner common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"transferStakeOwnership",
		ts.contractABI,
		ts.transactor,
		operator,
		newOwner,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) Undelegate(
	_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction undelegate",
		" params: ",
		fmt.Sprint(
			_operator,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.Undelegate(
		transactorOptions,
		_operator,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"undelegate",
			_operator,
		)
	}

	tsLogger.Infof(
		"submitted transaction undelegate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.Undelegate(
				newTransactorOptions,
				_operator,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"undelegate",
					_operator,
				)
			}

			tsLogger.Infof(
				"submitted transaction undelegate with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUndelegate(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"undelegate",
		&result,
		_operator,
	)

	return err
}

func (ts *TokenStaking) UndelegateGasEstimate(
	_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"undelegate",
		ts.contractABI,
		ts.transactor,
		_operator,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) UndelegateAt(
	_operator common.Address,
	_undelegationTimestamp *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction undelegateAt",
		" params: ",
		fmt.Sprint(
			_operator,
			_undelegationTimestamp,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.UndelegateAt(
		transactorOptions,
		_operator,
		_undelegationTimestamp,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"undelegateAt",
			_operator,
			_undelegationTimestamp,
		)
	}

	tsLogger.Infof(
		"submitted transaction undelegateAt with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.UndelegateAt(
				newTransactorOptions,
				_operator,
				_undelegationTimestamp,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"undelegateAt",
					_operator,
					_undelegationTimestamp,
				)
			}

			tsLogger.Infof(
				"submitted transaction undelegateAt with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUndelegateAt(
	_operator common.Address,
	_undelegationTimestamp *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"undelegateAt",
		&result,
		_operator,
		_undelegationTimestamp,
	)

	return err
}

func (ts *TokenStaking) UndelegateAtGasEstimate(
	_operator common.Address,
	_undelegationTimestamp *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"undelegateAt",
		ts.contractABI,
		ts.transactor,
		_operator,
		_undelegationTimestamp,
	)

	return result, err
}

// Transaction submission.
func (ts *TokenStaking) UnlockStake(
	operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction unlockStake",
		" params: ",
		fmt.Sprint(
			operator,
		),
	)

	ts.transactionMutex.Lock()
	defer ts.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *ts.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := ts.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := ts.contract.UnlockStake(
		transactorOptions,
		operator,
	)
	if err != nil {
		return transaction, ts.errorResolver.ResolveError(
			err,
			ts.transactorOptions.From,
			nil,
			"unlockStake",
			operator,
		)
	}

	tsLogger.Infof(
		"submitted transaction unlockStake with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go ts.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := ts.contract.UnlockStake(
				newTransactorOptions,
				operator,
			)
			if err != nil {
				return nil, ts.errorResolver.ResolveError(
					err,
					ts.transactorOptions.From,
					nil,
					"unlockStake",
					operator,
				)
			}

			tsLogger.Infof(
				"submitted transaction unlockStake with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	ts.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUnlockStake(
	operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		ts.transactorOptions.From,
		blockNumber, nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"unlockStake",
		&result,
		operator,
	)

	return err
}

func (ts *TokenStaking) UnlockStakeGasEstimate(
	operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		ts.callerOptions.From,
		ts.contractAddress,
		"unlockStake",
		ts.contractABI,
		ts.transactor,
		operator,
	)

	return result, err
}

// ----- Const Methods ------

func (ts *TokenStaking) ActiveStake(
	_operator common.Address,
	_operatorContract common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.ActiveStake(
		ts.callerOptions,
		_operator,
		_operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"activeStake",
			_operator,
			_operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) ActiveStakeAtBlock(
	_operator common.Address,
	_operatorContract common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"activeStake",
		&result,
		_operator,
		_operatorContract,
	)

	return result, err
}

func (ts *TokenStaking) AuthorizerOf(
	_operator common.Address,
) (common.Address, error) {
	var result common.Address
	result, err := ts.contract.AuthorizerOf(
		ts.callerOptions,
		_operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"authorizerOf",
			_operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) AuthorizerOfAtBlock(
	_operator common.Address,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"authorizerOf",
		&result,
		_operator,
	)

	return result, err
}

func (ts *TokenStaking) BalanceOf(
	_address common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.BalanceOf(
		ts.callerOptions,
		_address,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"balanceOf",
			_address,
		)
	}

	return result, err
}

func (ts *TokenStaking) BalanceOfAtBlock(
	_address common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"balanceOf",
		&result,
		_address,
	)

	return result, err
}

func (ts *TokenStaking) BeneficiaryOf(
	_operator common.Address,
) (common.Address, error) {
	var result common.Address
	result, err := ts.contract.BeneficiaryOf(
		ts.callerOptions,
		_operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"beneficiaryOf",
			_operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) BeneficiaryOfAtBlock(
	_operator common.Address,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"beneficiaryOf",
		&result,
		_operator,
	)

	return result, err
}

func (ts *TokenStaking) DeployedAt() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.DeployedAt(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"deployedAt",
		)
	}

	return result, err
}

func (ts *TokenStaking) DeployedAtAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"deployedAt",
		&result,
	)

	return result, err
}

func (ts *TokenStaking) EligibleStake(
	_operator common.Address,
	_operatorContract common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.EligibleStake(
		ts.callerOptions,
		_operator,
		_operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"eligibleStake",
			_operator,
			_operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) EligibleStakeAtBlock(
	_operator common.Address,
	_operatorContract common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"eligibleStake",
		&result,
		_operator,
		_operatorContract,
	)

	return result, err
}

func (ts *TokenStaking) GetAuthoritySource(
	operatorContract common.Address,
) (common.Address, error) {
	var result common.Address
	result, err := ts.contract.GetAuthoritySource(
		ts.callerOptions,
		operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"getAuthoritySource",
			operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) GetAuthoritySourceAtBlock(
	operatorContract common.Address,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"getAuthoritySource",
		&result,
		operatorContract,
	)

	return result, err
}

type DelegationInfo struct {
	Amount        *big.Int
	CreatedAt     *big.Int
	UndelegatedAt *big.Int
}

func (ts *TokenStaking) GetDelegationInfo(
	_operator common.Address,
) (DelegationInfo, error) {
	var result DelegationInfo
	result, err := ts.contract.GetDelegationInfo(
		ts.callerOptions,
		_operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"getDelegationInfo",
			_operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) GetDelegationInfoAtBlock(
	_operator common.Address,
	blockNumber *big.Int,
) (DelegationInfo, error) {
	var result DelegationInfo

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"getDelegationInfo",
		&result,
		_operator,
	)

	return result, err
}

type Locks struct {
	Creators    []common.Address
	Expirations []*big.Int
}

func (ts *TokenStaking) GetLocks(
	operator common.Address,
) (Locks, error) {
	var result Locks
	result, err := ts.contract.GetLocks(
		ts.callerOptions,
		operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"getLocks",
			operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) GetLocksAtBlock(
	operator common.Address,
	blockNumber *big.Int,
) (Locks, error) {
	var result Locks

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"getLocks",
		&result,
		operator,
	)

	return result, err
}

func (ts *TokenStaking) HasMinimumStake(
	staker common.Address,
	operatorContract common.Address,
) (bool, error) {
	var result bool
	result, err := ts.contract.HasMinimumStake(
		ts.callerOptions,
		staker,
		operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"hasMinimumStake",
			staker,
			operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) HasMinimumStakeAtBlock(
	staker common.Address,
	operatorContract common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"hasMinimumStake",
		&result,
		staker,
		operatorContract,
	)

	return result, err
}

func (ts *TokenStaking) InitializationPeriod() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.InitializationPeriod(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"initializationPeriod",
		)
	}

	return result, err
}

func (ts *TokenStaking) InitializationPeriodAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"initializationPeriod",
		&result,
	)

	return result, err
}

func (ts *TokenStaking) IsApprovedOperatorContract(
	_operatorContract common.Address,
) (bool, error) {
	var result bool
	result, err := ts.contract.IsApprovedOperatorContract(
		ts.callerOptions,
		_operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"isApprovedOperatorContract",
			_operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) IsApprovedOperatorContractAtBlock(
	_operatorContract common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"isApprovedOperatorContract",
		&result,
		_operatorContract,
	)

	return result, err
}

func (ts *TokenStaking) IsAuthorizedForOperator(
	_operator common.Address,
	_operatorContract common.Address,
) (bool, error) {
	var result bool
	result, err := ts.contract.IsAuthorizedForOperator(
		ts.callerOptions,
		_operator,
		_operatorContract,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"isAuthorizedForOperator",
			_operator,
			_operatorContract,
		)
	}

	return result, err
}

func (ts *TokenStaking) IsAuthorizedForOperatorAtBlock(
	_operator common.Address,
	_operatorContract common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"isAuthorizedForOperator",
		&result,
		_operator,
		_operatorContract,
	)

	return result, err
}

func (ts *TokenStaking) IsStakeLocked(
	operator common.Address,
) (bool, error) {
	var result bool
	result, err := ts.contract.IsStakeLocked(
		ts.callerOptions,
		operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"isStakeLocked",
			operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) IsStakeLockedAtBlock(
	operator common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"isStakeLocked",
		&result,
		operator,
	)

	return result, err
}

func (ts *TokenStaking) MinimumStake() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MinimumStake(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"minimumStake",
		)
	}

	return result, err
}

func (ts *TokenStaking) MinimumStakeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"minimumStake",
		&result,
	)

	return result, err
}

func (ts *TokenStaking) OwnerOf(
	_operator common.Address,
) (common.Address, error) {
	var result common.Address
	result, err := ts.contract.OwnerOf(
		ts.callerOptions,
		_operator,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"ownerOf",
			_operator,
		)
	}

	return result, err
}

func (ts *TokenStaking) OwnerOfAtBlock(
	_operator common.Address,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"ownerOf",
		&result,
		_operator,
	)

	return result, err
}

func (ts *TokenStaking) UndelegationPeriod() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.UndelegationPeriod(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"undelegationPeriod",
		)
	}

	return result, err
}

func (ts *TokenStaking) UndelegationPeriodAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"undelegationPeriod",
		&result,
	)

	return result, err
}

// ------ Events -------

func (ts *TokenStaking) ExpiredLockReleased(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsExpiredLockReleasedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsExpiredLockReleasedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsExpiredLockReleasedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingExpiredLockReleasedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	blockNumber uint64,
)

func (elrs *TsExpiredLockReleasedSubscription) OnEvent(
	handler tokenStakingExpiredLockReleasedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingExpiredLockReleased)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.LockCreator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := elrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (elrs *TsExpiredLockReleasedSubscription) Pipe(
	sink chan *abi.TokenStakingExpiredLockReleased,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(elrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := elrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - elrs.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past ExpiredLockReleased events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := elrs.contract.PastExpiredLockReleasedEvents(
					fromBlock,
					nil,
					elrs.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past ExpiredLockReleased events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := elrs.contract.watchExpiredLockReleased(
		sink,
		elrs.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchExpiredLockReleased(
	sink chan *abi.TokenStakingExpiredLockReleased,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchExpiredLockReleased(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event ExpiredLockReleased had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event ExpiredLockReleased failed "+
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

func (ts *TokenStaking) PastExpiredLockReleasedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingExpiredLockReleased, error) {
	iterator, err := ts.contract.FilterExpiredLockReleased(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past ExpiredLockReleased events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingExpiredLockReleased, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) LockReleased(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsLockReleasedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsLockReleasedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsLockReleasedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingLockReleasedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	blockNumber uint64,
)

func (lrs *TsLockReleasedSubscription) OnEvent(
	handler tokenStakingLockReleasedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingLockReleased)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.LockCreator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := lrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (lrs *TsLockReleasedSubscription) Pipe(
	sink chan *abi.TokenStakingLockReleased,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(lrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := lrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - lrs.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past LockReleased events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := lrs.contract.PastLockReleasedEvents(
					fromBlock,
					nil,
					lrs.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past LockReleased events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := lrs.contract.watchLockReleased(
		sink,
		lrs.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchLockReleased(
	sink chan *abi.TokenStakingLockReleased,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchLockReleased(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event LockReleased had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event LockReleased failed "+
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

func (ts *TokenStaking) PastLockReleasedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingLockReleased, error) {
	iterator, err := ts.contract.FilterLockReleased(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past LockReleased events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingLockReleased, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) OperatorStaked(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
	beneficiaryFilter []common.Address,
	authorizerFilter []common.Address,
) *TsOperatorStakedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsOperatorStakedSubscription{
		ts,
		opts,
		operatorFilter,
		beneficiaryFilter,
		authorizerFilter,
	}
}

type TsOperatorStakedSubscription struct {
	contract          *TokenStaking
	opts              *ethlike.SubscribeOpts
	operatorFilter    []common.Address
	beneficiaryFilter []common.Address
	authorizerFilter  []common.Address
}

type tokenStakingOperatorStakedFunc func(
	Operator common.Address,
	Beneficiary common.Address,
	Authorizer common.Address,
	Value *big.Int,
	blockNumber uint64,
)

func (oss *TsOperatorStakedSubscription) OnEvent(
	handler tokenStakingOperatorStakedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingOperatorStaked)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.Beneficiary,
					event.Authorizer,
					event.Value,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := oss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (oss *TsOperatorStakedSubscription) Pipe(
	sink chan *abi.TokenStakingOperatorStaked,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(oss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := oss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - oss.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past OperatorStaked events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := oss.contract.PastOperatorStakedEvents(
					fromBlock,
					nil,
					oss.operatorFilter,
					oss.beneficiaryFilter,
					oss.authorizerFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past OperatorStaked events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := oss.contract.watchOperatorStaked(
		sink,
		oss.operatorFilter,
		oss.beneficiaryFilter,
		oss.authorizerFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchOperatorStaked(
	sink chan *abi.TokenStakingOperatorStaked,
	operatorFilter []common.Address,
	beneficiaryFilter []common.Address,
	authorizerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchOperatorStaked(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
			beneficiaryFilter,
			authorizerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event OperatorStaked had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event OperatorStaked failed "+
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

func (ts *TokenStaking) PastOperatorStakedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
	beneficiaryFilter []common.Address,
	authorizerFilter []common.Address,
) ([]*abi.TokenStakingOperatorStaked, error) {
	iterator, err := ts.contract.FilterOperatorStaked(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
		beneficiaryFilter,
		authorizerFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past OperatorStaked events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingOperatorStaked, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) RecoveredStake(
	opts *ethlike.SubscribeOpts,
) *TsRecoveredStakeSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsRecoveredStakeSubscription{
		ts,
		opts,
	}
}

type TsRecoveredStakeSubscription struct {
	contract *TokenStaking
	opts     *ethlike.SubscribeOpts
}

type tokenStakingRecoveredStakeFunc func(
	Operator common.Address,
	blockNumber uint64,
)

func (rss *TsRecoveredStakeSubscription) OnEvent(
	handler tokenStakingRecoveredStakeFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingRecoveredStake)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rss *TsRecoveredStakeSubscription) Pipe(
	sink chan *abi.TokenStakingRecoveredStake,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rss.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past RecoveredStake events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rss.contract.PastRecoveredStakeEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past RecoveredStake events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rss.contract.watchRecoveredStake(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchRecoveredStake(
	sink chan *abi.TokenStakingRecoveredStake,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchRecoveredStake(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event RecoveredStake had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event RecoveredStake failed "+
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

func (ts *TokenStaking) PastRecoveredStakeEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TokenStakingRecoveredStake, error) {
	iterator, err := ts.contract.FilterRecoveredStake(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RecoveredStake events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingRecoveredStake, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) StakeDelegated(
	opts *ethlike.SubscribeOpts,
	ownerFilter []common.Address,
	operatorFilter []common.Address,
) *TsStakeDelegatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsStakeDelegatedSubscription{
		ts,
		opts,
		ownerFilter,
		operatorFilter,
	}
}

type TsStakeDelegatedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	ownerFilter    []common.Address
	operatorFilter []common.Address
}

type tokenStakingStakeDelegatedFunc func(
	Owner common.Address,
	Operator common.Address,
	blockNumber uint64,
)

func (sds *TsStakeDelegatedSubscription) OnEvent(
	handler tokenStakingStakeDelegatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingStakeDelegated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Owner,
					event.Operator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sds.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sds *TsStakeDelegatedSubscription) Pipe(
	sink chan *abi.TokenStakingStakeDelegated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sds.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sds.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sds.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past StakeDelegated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sds.contract.PastStakeDelegatedEvents(
					fromBlock,
					nil,
					sds.ownerFilter,
					sds.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past StakeDelegated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sds.contract.watchStakeDelegated(
		sink,
		sds.ownerFilter,
		sds.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchStakeDelegated(
	sink chan *abi.TokenStakingStakeDelegated,
	ownerFilter []common.Address,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchStakeDelegated(
			&bind.WatchOpts{Context: ctx},
			sink,
			ownerFilter,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event StakeDelegated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event StakeDelegated failed "+
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

func (ts *TokenStaking) PastStakeDelegatedEvents(
	startBlock uint64,
	endBlock *uint64,
	ownerFilter []common.Address,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingStakeDelegated, error) {
	iterator, err := ts.contract.FilterStakeDelegated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		ownerFilter,
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past StakeDelegated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingStakeDelegated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) StakeLocked(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsStakeLockedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsStakeLockedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsStakeLockedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingStakeLockedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	Until *big.Int,
	blockNumber uint64,
)

func (sls *TsStakeLockedSubscription) OnEvent(
	handler tokenStakingStakeLockedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingStakeLocked)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.LockCreator,
					event.Until,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sls.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sls *TsStakeLockedSubscription) Pipe(
	sink chan *abi.TokenStakingStakeLocked,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sls.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sls.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sls.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past StakeLocked events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sls.contract.PastStakeLockedEvents(
					fromBlock,
					nil,
					sls.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past StakeLocked events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sls.contract.watchStakeLocked(
		sink,
		sls.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchStakeLocked(
	sink chan *abi.TokenStakingStakeLocked,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchStakeLocked(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event StakeLocked had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event StakeLocked failed "+
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

func (ts *TokenStaking) PastStakeLockedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingStakeLocked, error) {
	iterator, err := ts.contract.FilterStakeLocked(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past StakeLocked events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingStakeLocked, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) StakeOwnershipTransferred(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
	newOwnerFilter []common.Address,
) *TsStakeOwnershipTransferredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsStakeOwnershipTransferredSubscription{
		ts,
		opts,
		operatorFilter,
		newOwnerFilter,
	}
}

type TsStakeOwnershipTransferredSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
	newOwnerFilter []common.Address
}

type tokenStakingStakeOwnershipTransferredFunc func(
	Operator common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (sots *TsStakeOwnershipTransferredSubscription) OnEvent(
	handler tokenStakingStakeOwnershipTransferredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingStakeOwnershipTransferred)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.NewOwner,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sots.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sots *TsStakeOwnershipTransferredSubscription) Pipe(
	sink chan *abi.TokenStakingStakeOwnershipTransferred,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sots.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sots.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sots.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past StakeOwnershipTransferred events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sots.contract.PastStakeOwnershipTransferredEvents(
					fromBlock,
					nil,
					sots.operatorFilter,
					sots.newOwnerFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past StakeOwnershipTransferred events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sots.contract.watchStakeOwnershipTransferred(
		sink,
		sots.operatorFilter,
		sots.newOwnerFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchStakeOwnershipTransferred(
	sink chan *abi.TokenStakingStakeOwnershipTransferred,
	operatorFilter []common.Address,
	newOwnerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchStakeOwnershipTransferred(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
			newOwnerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event StakeOwnershipTransferred had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event StakeOwnershipTransferred failed "+
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

func (ts *TokenStaking) PastStakeOwnershipTransferredEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
	newOwnerFilter []common.Address,
) ([]*abi.TokenStakingStakeOwnershipTransferred, error) {
	iterator, err := ts.contract.FilterStakeOwnershipTransferred(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
		newOwnerFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past StakeOwnershipTransferred events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingStakeOwnershipTransferred, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) TokensSeized(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsTokensSeizedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsTokensSeizedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsTokensSeizedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingTokensSeizedFunc func(
	Operator common.Address,
	Amount *big.Int,
	blockNumber uint64,
)

func (tss *TsTokensSeizedSubscription) OnEvent(
	handler tokenStakingTokensSeizedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingTokensSeized)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.Amount,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tss *TsTokensSeizedSubscription) Pipe(
	sink chan *abi.TokenStakingTokensSeized,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tss.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past TokensSeized events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tss.contract.PastTokensSeizedEvents(
					fromBlock,
					nil,
					tss.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past TokensSeized events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tss.contract.watchTokensSeized(
		sink,
		tss.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchTokensSeized(
	sink chan *abi.TokenStakingTokensSeized,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchTokensSeized(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event TokensSeized had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event TokensSeized failed "+
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

func (ts *TokenStaking) PastTokensSeizedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingTokensSeized, error) {
	iterator, err := ts.contract.FilterTokensSeized(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokensSeized events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingTokensSeized, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) TokensSlashed(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsTokensSlashedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsTokensSlashedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsTokensSlashedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingTokensSlashedFunc func(
	Operator common.Address,
	Amount *big.Int,
	blockNumber uint64,
)

func (tss *TsTokensSlashedSubscription) OnEvent(
	handler tokenStakingTokensSlashedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingTokensSlashed)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.Amount,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tss *TsTokensSlashedSubscription) Pipe(
	sink chan *abi.TokenStakingTokensSlashed,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tss.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past TokensSlashed events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tss.contract.PastTokensSlashedEvents(
					fromBlock,
					nil,
					tss.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past TokensSlashed events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tss.contract.watchTokensSlashed(
		sink,
		tss.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchTokensSlashed(
	sink chan *abi.TokenStakingTokensSlashed,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchTokensSlashed(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event TokensSlashed had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event TokensSlashed failed "+
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

func (ts *TokenStaking) PastTokensSlashedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingTokensSlashed, error) {
	iterator, err := ts.contract.FilterTokensSlashed(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TokensSlashed events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingTokensSlashed, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) TopUpCompleted(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsTopUpCompletedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsTopUpCompletedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsTopUpCompletedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingTopUpCompletedFunc func(
	Operator common.Address,
	NewAmount *big.Int,
	blockNumber uint64,
)

func (tucs *TsTopUpCompletedSubscription) OnEvent(
	handler tokenStakingTopUpCompletedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingTopUpCompleted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.NewAmount,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tucs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tucs *TsTopUpCompletedSubscription) Pipe(
	sink chan *abi.TokenStakingTopUpCompleted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tucs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tucs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tucs.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past TopUpCompleted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tucs.contract.PastTopUpCompletedEvents(
					fromBlock,
					nil,
					tucs.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past TopUpCompleted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tucs.contract.watchTopUpCompleted(
		sink,
		tucs.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchTopUpCompleted(
	sink chan *abi.TokenStakingTopUpCompleted,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchTopUpCompleted(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event TopUpCompleted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event TopUpCompleted failed "+
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

func (ts *TokenStaking) PastTopUpCompletedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingTopUpCompleted, error) {
	iterator, err := ts.contract.FilterTopUpCompleted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TopUpCompleted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingTopUpCompleted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) TopUpInitiated(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsTopUpInitiatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsTopUpInitiatedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsTopUpInitiatedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingTopUpInitiatedFunc func(
	Operator common.Address,
	TopUp *big.Int,
	blockNumber uint64,
)

func (tuis *TsTopUpInitiatedSubscription) OnEvent(
	handler tokenStakingTopUpInitiatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingTopUpInitiated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.TopUp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := tuis.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tuis *TsTopUpInitiatedSubscription) Pipe(
	sink chan *abi.TokenStakingTopUpInitiated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(tuis.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := tuis.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - tuis.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past TopUpInitiated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := tuis.contract.PastTopUpInitiatedEvents(
					fromBlock,
					nil,
					tuis.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past TopUpInitiated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := tuis.contract.watchTopUpInitiated(
		sink,
		tuis.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchTopUpInitiated(
	sink chan *abi.TokenStakingTopUpInitiated,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchTopUpInitiated(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event TopUpInitiated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event TopUpInitiated failed "+
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

func (ts *TokenStaking) PastTopUpInitiatedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingTopUpInitiated, error) {
	iterator, err := ts.contract.FilterTopUpInitiated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past TopUpInitiated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingTopUpInitiated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (ts *TokenStaking) Undelegated(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
) *TsUndelegatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TsUndelegatedSubscription{
		ts,
		opts,
		operatorFilter,
	}
}

type TsUndelegatedSubscription struct {
	contract       *TokenStaking
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
}

type tokenStakingUndelegatedFunc func(
	Operator common.Address,
	UndelegatedAt *big.Int,
	blockNumber uint64,
)

func (us *TsUndelegatedSubscription) OnEvent(
	handler tokenStakingUndelegatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TokenStakingUndelegated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.UndelegatedAt,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := us.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (us *TsUndelegatedSubscription) Pipe(
	sink chan *abi.TokenStakingUndelegated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(us.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := us.contract.blockCounter.CurrentBlock()
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - us.opts.PastBlocks

				tsLogger.Infof(
					"subscription monitoring fetching past Undelegated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := us.contract.PastUndelegatedEvents(
					fromBlock,
					nil,
					us.operatorFilter,
				)
				if err != nil {
					tsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tsLogger.Infof(
					"subscription monitoring fetched [%v] past Undelegated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := us.contract.watchUndelegated(
		sink,
		us.operatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ts *TokenStaking) watchUndelegated(
	sink chan *abi.TokenStakingUndelegated,
	operatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return ts.contract.WatchUndelegated(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tsLogger.Errorf(
			"subscription to event Undelegated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tsLogger.Errorf(
			"subscription to event Undelegated failed "+
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

func (ts *TokenStaking) PastUndelegatedEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
) ([]*abi.TokenStakingUndelegated, error) {
	iterator, err := ts.contract.FilterUndelegated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Undelegated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TokenStakingUndelegated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
