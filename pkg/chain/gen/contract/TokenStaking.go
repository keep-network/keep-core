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
var tsLogger = log.Logger("keep-contract-TokenStaking")

type TokenStaking struct {
	contract          *abi.TokenStaking
	contractAddress   common.Address
	contractABI       *ethereumabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *ethutil.ErrorResolver

	transactionMutex *sync.Mutex
}

func NewTokenStaking(
	contractAddress common.Address,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	transactionMutex *sync.Mutex,
) (*TokenStaking, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	transactorOptions := bind.NewKeyedTransactor(
		accountKey.PrivateKey,
	)

	randomBeaconContract, err := abi.NewTokenStaking(
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

	contractABI, err := ethereumabi.JSON(strings.NewReader(abi.TokenStakingABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TokenStaking{
		contract:          randomBeaconContract,
		contractAddress:   contractAddress,
		contractABI:       &contractABI,
		caller:            backend,
		transactor:        backend,
		callerOptions:     callerOptions,
		transactorOptions: transactorOptions,
		errorResolver:     ethutil.NewErrorResolver(backend, &contractABI, &contractAddress),
		transactionMutex:  transactionMutex,
	}, nil
}

// ----- Non-const Methods ------

// Transaction submission.
func (ts *TokenStaking) LockStake(
	operator common.Address,
	duration *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction lockStake",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction lockStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallLockStake(
	operator common.Address,
	duration *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction receiveApproval",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction receiveApproval with id: [%v]",
		transaction.Hash().Hex(),
	)

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

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (ts *TokenStaking) Undelegate(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction undelegate",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction undelegate with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUndelegate(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (ts *TokenStaking) AuthorizeOperatorContract(
	_operator common.Address,
	_operatorContract common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction authorizeOperatorContract",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction authorizeOperatorContract with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallAuthorizeOperatorContract(
	_operator common.Address,
	_operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction cancelStake",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction cancelStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallCancelStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction claimDelegatedAuthority",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction claimDelegatedAuthority with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallClaimDelegatedAuthority(
	delegatedAuthoritySource common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (ts *TokenStaking) RecoverStake(
	_operator common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction recoverStake",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction recoverStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallRecoverStake(
	_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction releaseExpiredLock",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction releaseExpiredLock with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallReleaseExpiredLock(
	operator common.Address,
	operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction seize",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction seize with id: [%v]",
		transaction.Hash().Hex(),
	)

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

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction slash",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction slash with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallSlash(
	amountToSlash *big.Int,
	misbehavedOperators []common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (ts *TokenStaking) UndelegateAt(
	_operator common.Address,
	_undelegationTimestamp *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction undelegateAt",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction undelegateAt with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUndelegateAt(
	_operator common.Address,
	_undelegationTimestamp *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tsLogger.Debug(
		"submitting transaction unlockStake",
		"params: ",
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

	tsLogger.Debugf(
		"submitted transaction unlockStake with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (ts *TokenStaking) CallUnlockStake(
	operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

func (ts *TokenStaking) MinimumStakeSchedule() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MinimumStakeSchedule(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"minimumStakeSchedule",
		)
	}

	return result, err
}

func (ts *TokenStaking) MinimumStakeScheduleAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"minimumStakeSchedule",
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

func (ts *TokenStaking) MinimumStakeSteps() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MinimumStakeSteps(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"minimumStakeSteps",
		)
	}

	return result, err
}

func (ts *TokenStaking) MinimumStakeStepsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"minimumStakeSteps",
		&result,
	)

	return result, err
}

func (ts *TokenStaking) OwnerOperators(
	arg0 common.Address,
	arg1 *big.Int,
) (common.Address, error) {
	var result common.Address
	result, err := ts.contract.OwnerOperators(
		ts.callerOptions,
		arg0,
		arg1,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"ownerOperators",
			arg0,
			arg1,
		)
	}

	return result, err
}

func (ts *TokenStaking) OwnerOperatorsAtBlock(
	arg0 common.Address,
	arg1 *big.Int,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"ownerOperators",
		&result,
		arg0,
		arg1,
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

func (ts *TokenStaking) MinimumStakeBase() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MinimumStakeBase(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"minimumStakeBase",
		)
	}

	return result, err
}

func (ts *TokenStaking) MinimumStakeBaseAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"minimumStakeBase",
		&result,
	)

	return result, err
}

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

	err := ethutil.CallAtBlock(
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

func (ts *TokenStaking) MaximumLockDuration() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MaximumLockDuration(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"maximumLockDuration",
		)
	}

	return result, err
}

func (ts *TokenStaking) MaximumLockDurationAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"maximumLockDuration",
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

func (ts *TokenStaking) Registry() (common.Address, error) {
	var result common.Address
	result, err := ts.contract.Registry(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"registry",
		)
	}

	return result, err
}

func (ts *TokenStaking) RegistryAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"registry",
		&result,
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

type operators struct {
	PackedParams *big.Int
	Owner        common.Address
	Beneficiary  common.Address
	Authorizer   common.Address
}

func (ts *TokenStaking) Operators(
	arg0 common.Address,
) (operators, error) {
	var result operators
	result, err := ts.contract.Operators(
		ts.callerOptions,
		arg0,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"operators",
			arg0,
		)
	}

	return result, err
}

func (ts *TokenStaking) OperatorsAtBlock(
	arg0 common.Address,
	blockNumber *big.Int,
) (operators, error) {
	var result operators

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"operators",
		&result,
		arg0,
	)

	return result, err
}

func (ts *TokenStaking) OperatorsOf(
	_address common.Address,
) ([]common.Address, error) {
	var result []common.Address
	result, err := ts.contract.OperatorsOf(
		ts.callerOptions,
		_address,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"operatorsOf",
			_address,
		)
	}

	return result, err
}

func (ts *TokenStaking) OperatorsOfAtBlock(
	_address common.Address,
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"operatorsOf",
		&result,
		_address,
	)

	return result, err
}

func (ts *TokenStaking) Token() (common.Address, error) {
	var result common.Address
	result, err := ts.contract.Token(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"token",
		)
	}

	return result, err
}

func (ts *TokenStaking) TokenAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"token",
		&result,
	)

	return result, err
}

func (ts *TokenStaking) MinimumStakeScheduleStart() (*big.Int, error) {
	var result *big.Int
	result, err := ts.contract.MinimumStakeScheduleStart(
		ts.callerOptions,
	)

	if err != nil {
		return result, ts.errorResolver.ResolveError(
			err,
			ts.callerOptions.From,
			nil,
			"minimumStakeScheduleStart",
		)
	}

	return result, err
}

func (ts *TokenStaking) MinimumStakeScheduleStartAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		ts.callerOptions.From,
		blockNumber,
		nil,
		ts.contractABI,
		ts.caller,
		ts.errorResolver,
		ts.contractAddress,
		"minimumStakeScheduleStart",
		&result,
	)

	return result, err
}

// ------ Events -------

type tokenStakingStakedFunc func(
	From common.Address,
	Value *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchStaked(
	success tokenStakingStakedFunc,
	fail func(err error) error,
	fromFilter []common.Address,
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

		subscription, err := ts.subscribeStaked(
			success,
			failCallback,
			fromFilter,
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
				tsLogger.Warning(
					"subscription to event Staked terminated with error; " +
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

func (ts *TokenStaking) subscribeStaked(
	success tokenStakingStakedFunc,
	fail func(err error) error,
	fromFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingStaked)
	eventSubscription, err := ts.contract.WatchStaked(
		nil,
		eventChan,
		fromFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Staked events: [%v]",
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
					event.From,
					event.Value,
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

type tokenStakingTokensSeizedFunc func(
	Operator common.Address,
	Amount *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchTokensSeized(
	success tokenStakingTokensSeizedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeTokensSeized(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event TokensSeized terminated with error; " +
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

func (ts *TokenStaking) subscribeTokensSeized(
	success tokenStakingTokensSeizedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingTokensSeized)
	eventSubscription, err := ts.contract.WatchTokensSeized(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokensSeized events: [%v]",
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
					event.Operator,
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

type tokenStakingTokensSlashedFunc func(
	Operator common.Address,
	Amount *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchTokensSlashed(
	success tokenStakingTokensSlashedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeTokensSlashed(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event TokensSlashed terminated with error; " +
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

func (ts *TokenStaking) subscribeTokensSlashed(
	success tokenStakingTokensSlashedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingTokensSlashed)
	eventSubscription, err := ts.contract.WatchTokensSlashed(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for TokensSlashed events: [%v]",
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
					event.Operator,
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

type tokenStakingUndelegatedFunc func(
	Operator common.Address,
	UndelegatedAt *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchUndelegated(
	success tokenStakingUndelegatedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeUndelegated(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event Undelegated terminated with error; " +
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

func (ts *TokenStaking) subscribeUndelegated(
	success tokenStakingUndelegatedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingUndelegated)
	eventSubscription, err := ts.contract.WatchUndelegated(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Undelegated events: [%v]",
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
					event.Operator,
					event.UndelegatedAt,
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

type tokenStakingExpiredLockReleasedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchExpiredLockReleased(
	success tokenStakingExpiredLockReleasedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeExpiredLockReleased(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event ExpiredLockReleased terminated with error; " +
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

func (ts *TokenStaking) subscribeExpiredLockReleased(
	success tokenStakingExpiredLockReleasedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingExpiredLockReleased)
	eventSubscription, err := ts.contract.WatchExpiredLockReleased(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for ExpiredLockReleased events: [%v]",
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
					event.Operator,
					event.LockCreator,
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

type tokenStakingLockReleasedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchLockReleased(
	success tokenStakingLockReleasedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeLockReleased(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event LockReleased terminated with error; " +
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

func (ts *TokenStaking) subscribeLockReleased(
	success tokenStakingLockReleasedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingLockReleased)
	eventSubscription, err := ts.contract.WatchLockReleased(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for LockReleased events: [%v]",
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
					event.Operator,
					event.LockCreator,
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

type tokenStakingRecoveredStakeFunc func(
	Operator common.Address,
	RecoveredAt *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchRecoveredStake(
	success tokenStakingRecoveredStakeFunc,
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

		subscription, err := ts.subscribeRecoveredStake(
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
				tsLogger.Warning(
					"subscription to event RecoveredStake terminated with error; " +
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

func (ts *TokenStaking) subscribeRecoveredStake(
	success tokenStakingRecoveredStakeFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingRecoveredStake)
	eventSubscription, err := ts.contract.WatchRecoveredStake(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RecoveredStake events: [%v]",
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
					event.Operator,
					event.RecoveredAt,
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

type tokenStakingStakeLockedFunc func(
	Operator common.Address,
	LockCreator common.Address,
	Until *big.Int,
	blockNumber uint64,
)

func (ts *TokenStaking) WatchStakeLocked(
	success tokenStakingStakeLockedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
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

		subscription, err := ts.subscribeStakeLocked(
			success,
			failCallback,
			operatorFilter,
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
				tsLogger.Warning(
					"subscription to event StakeLocked terminated with error; " +
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

func (ts *TokenStaking) subscribeStakeLocked(
	success tokenStakingStakeLockedFunc,
	fail func(err error) error,
	operatorFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TokenStakingStakeLocked)
	eventSubscription, err := ts.contract.WatchStakeLocked(
		nil,
		eventChan,
		operatorFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for StakeLocked events: [%v]",
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
					event.Operator,
					event.LockCreator,
					event.Until,
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
