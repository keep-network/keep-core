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
	"github.com/keep-network/keep-core/pkg/chain/random-beacon/gen/abi"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var spLogger = log.Logger("keep-contract-SortitionPool")

type SortitionPool struct {
	contract          *abi.SortitionPool
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

func NewSortitionPool(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*SortitionPool, error) {
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

	contract, err := abi.NewSortitionPool(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.SortitionPoolABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &SortitionPool{
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
func (sp *SortitionPool) InsertOperator(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction insertOperator",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_authorizedStake,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.InsertOperator(
		transactorOptions,
		arg_operator,
		arg_authorizedStake,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"insertOperator",
			arg_operator,
			arg_authorizedStake,
		)
	}

	spLogger.Infof(
		"submitted transaction insertOperator with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.InsertOperator(
				newTransactorOptions,
				arg_operator,
				arg_authorizedStake,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"insertOperator",
					arg_operator,
					arg_authorizedStake,
				)
			}

			spLogger.Infof(
				"submitted transaction insertOperator with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallInsertOperator(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"insertOperator",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return err
}

func (sp *SortitionPool) InsertOperatorGasEstimate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"insertOperator",
		sp.contractABI,
		sp.transactor,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) Lock(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction lock",
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.Lock(
		transactorOptions,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"lock",
		)
	}

	spLogger.Infof(
		"submitted transaction lock with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.Lock(
				newTransactorOptions,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"lock",
				)
			}

			spLogger.Infof(
				"submitted transaction lock with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallLock(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"lock",
		&result,
	)

	return err
}

func (sp *SortitionPool) LockGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"lock",
		sp.contractABI,
		sp.transactor,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) ReceiveApproval(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction receiveApproval",
		" params: ",
		fmt.Sprint(
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.ReceiveApproval(
		transactorOptions,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"receiveApproval",
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
		)
	}

	spLogger.Infof(
		"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.ReceiveApproval(
				newTransactorOptions,
				arg_sender,
				arg_amount,
				arg_token,
				arg3,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"receiveApproval",
					arg_sender,
					arg_amount,
					arg_token,
					arg3,
				)
			}

			spLogger.Infof(
				"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallReceiveApproval(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"receiveApproval",
		&result,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)

	return err
}

func (sp *SortitionPool) ReceiveApprovalGasEstimate(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"receiveApproval",
		sp.contractABI,
		sp.transactor,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) RenounceOwnership(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction renounceOwnership",
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.RenounceOwnership(
		transactorOptions,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"renounceOwnership",
		)
	}

	spLogger.Infof(
		"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.RenounceOwnership(
				newTransactorOptions,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"renounceOwnership",
				)
			}

			spLogger.Infof(
				"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallRenounceOwnership(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"renounceOwnership",
		&result,
	)

	return err
}

func (sp *SortitionPool) RenounceOwnershipGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"renounceOwnership",
		sp.contractABI,
		sp.transactor,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) RestoreRewardEligibility(
	arg_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction restoreRewardEligibility",
		" params: ",
		fmt.Sprint(
			arg_operator,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.RestoreRewardEligibility(
		transactorOptions,
		arg_operator,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"restoreRewardEligibility",
			arg_operator,
		)
	}

	spLogger.Infof(
		"submitted transaction restoreRewardEligibility with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.RestoreRewardEligibility(
				newTransactorOptions,
				arg_operator,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"restoreRewardEligibility",
					arg_operator,
				)
			}

			spLogger.Infof(
				"submitted transaction restoreRewardEligibility with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallRestoreRewardEligibility(
	arg_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"restoreRewardEligibility",
		&result,
		arg_operator,
	)

	return err
}

func (sp *SortitionPool) RestoreRewardEligibilityGasEstimate(
	arg_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"restoreRewardEligibility",
		sp.contractABI,
		sp.transactor,
		arg_operator,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) SetRewardIneligibility(
	arg_operators []uint32,
	arg_until *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction setRewardIneligibility",
		" params: ",
		fmt.Sprint(
			arg_operators,
			arg_until,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.SetRewardIneligibility(
		transactorOptions,
		arg_operators,
		arg_until,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"setRewardIneligibility",
			arg_operators,
			arg_until,
		)
	}

	spLogger.Infof(
		"submitted transaction setRewardIneligibility with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.SetRewardIneligibility(
				newTransactorOptions,
				arg_operators,
				arg_until,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"setRewardIneligibility",
					arg_operators,
					arg_until,
				)
			}

			spLogger.Infof(
				"submitted transaction setRewardIneligibility with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallSetRewardIneligibility(
	arg_operators []uint32,
	arg_until *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"setRewardIneligibility",
		&result,
		arg_operators,
		arg_until,
	)

	return err
}

func (sp *SortitionPool) SetRewardIneligibilityGasEstimate(
	arg_operators []uint32,
	arg_until *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"setRewardIneligibility",
		sp.contractABI,
		sp.transactor,
		arg_operators,
		arg_until,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) TransferOwnership(
	arg_newOwner common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction transferOwnership",
		" params: ",
		fmt.Sprint(
			arg_newOwner,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.TransferOwnership(
		transactorOptions,
		arg_newOwner,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"transferOwnership",
			arg_newOwner,
		)
	}

	spLogger.Infof(
		"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.TransferOwnership(
				newTransactorOptions,
				arg_newOwner,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"transferOwnership",
					arg_newOwner,
				)
			}

			spLogger.Infof(
				"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallTransferOwnership(
	arg_newOwner common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"transferOwnership",
		&result,
		arg_newOwner,
	)

	return err
}

func (sp *SortitionPool) TransferOwnershipGasEstimate(
	arg_newOwner common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"transferOwnership",
		sp.contractABI,
		sp.transactor,
		arg_newOwner,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) Unlock(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction unlock",
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.Unlock(
		transactorOptions,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"unlock",
		)
	}

	spLogger.Infof(
		"submitted transaction unlock with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.Unlock(
				newTransactorOptions,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"unlock",
				)
			}

			spLogger.Infof(
				"submitted transaction unlock with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallUnlock(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"unlock",
		&result,
	)

	return err
}

func (sp *SortitionPool) UnlockGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"unlock",
		sp.contractABI,
		sp.transactor,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) UpdateOperatorStatus(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction updateOperatorStatus",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_authorizedStake,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.UpdateOperatorStatus(
		transactorOptions,
		arg_operator,
		arg_authorizedStake,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"updateOperatorStatus",
			arg_operator,
			arg_authorizedStake,
		)
	}

	spLogger.Infof(
		"submitted transaction updateOperatorStatus with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.UpdateOperatorStatus(
				newTransactorOptions,
				arg_operator,
				arg_authorizedStake,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"updateOperatorStatus",
					arg_operator,
					arg_authorizedStake,
				)
			}

			spLogger.Infof(
				"submitted transaction updateOperatorStatus with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallUpdateOperatorStatus(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"updateOperatorStatus",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return err
}

func (sp *SortitionPool) UpdateOperatorStatusGasEstimate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"updateOperatorStatus",
		sp.contractABI,
		sp.transactor,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) WithdrawIneligible(
	arg_recipient common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction withdrawIneligible",
		" params: ",
		fmt.Sprint(
			arg_recipient,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.WithdrawIneligible(
		transactorOptions,
		arg_recipient,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"withdrawIneligible",
			arg_recipient,
		)
	}

	spLogger.Infof(
		"submitted transaction withdrawIneligible with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.WithdrawIneligible(
				newTransactorOptions,
				arg_recipient,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"withdrawIneligible",
					arg_recipient,
				)
			}

			spLogger.Infof(
				"submitted transaction withdrawIneligible with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallWithdrawIneligible(
	arg_recipient common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"withdrawIneligible",
		&result,
		arg_recipient,
	)

	return err
}

func (sp *SortitionPool) WithdrawIneligibleGasEstimate(
	arg_recipient common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"withdrawIneligible",
		sp.contractABI,
		sp.transactor,
		arg_recipient,
	)

	return result, err
}

// Transaction submission.
func (sp *SortitionPool) WithdrawRewards(
	arg_operator common.Address,
	arg_beneficiary common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	spLogger.Debug(
		"submitting transaction withdrawRewards",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_beneficiary,
		),
	)

	sp.transactionMutex.Lock()
	defer sp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *sp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := sp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := sp.contract.WithdrawRewards(
		transactorOptions,
		arg_operator,
		arg_beneficiary,
	)
	if err != nil {
		return transaction, sp.errorResolver.ResolveError(
			err,
			sp.transactorOptions.From,
			nil,
			"withdrawRewards",
			arg_operator,
			arg_beneficiary,
		)
	}

	spLogger.Infof(
		"submitted transaction withdrawRewards with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go sp.miningWaiter.ForceMining(
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

			transaction, err := sp.contract.WithdrawRewards(
				newTransactorOptions,
				arg_operator,
				arg_beneficiary,
			)
			if err != nil {
				return nil, sp.errorResolver.ResolveError(
					err,
					sp.transactorOptions.From,
					nil,
					"withdrawRewards",
					arg_operator,
					arg_beneficiary,
				)
			}

			spLogger.Infof(
				"submitted transaction withdrawRewards with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	sp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (sp *SortitionPool) CallWithdrawRewards(
	arg_operator common.Address,
	arg_beneficiary common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.transactorOptions.From,
		blockNumber, nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"withdrawRewards",
		&result,
		arg_operator,
		arg_beneficiary,
	)

	return result, err
}

func (sp *SortitionPool) WithdrawRewardsGasEstimate(
	arg_operator common.Address,
	arg_beneficiary common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		sp.callerOptions.From,
		sp.contractAddress,
		"withdrawRewards",
		sp.contractABI,
		sp.transactor,
		arg_operator,
		arg_beneficiary,
	)

	return result, err
}

// ----- Const Methods ------

func (sp *SortitionPool) CanRestoreRewardEligibility(
	arg_operator uint32,
) (bool, error) {
	result, err := sp.contract.CanRestoreRewardEligibility(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"canRestoreRewardEligibility",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) CanRestoreRewardEligibilityAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"canRestoreRewardEligibility",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) GetAvailableRewards(
	arg_operator common.Address,
) (*big.Int, error) {
	result, err := sp.contract.GetAvailableRewards(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"getAvailableRewards",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) GetAvailableRewardsAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"getAvailableRewards",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) GetIDOperator(
	arg_id uint32,
) (common.Address, error) {
	result, err := sp.contract.GetIDOperator(
		sp.callerOptions,
		arg_id,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"getIDOperator",
			arg_id,
		)
	}

	return result, err
}

func (sp *SortitionPool) GetIDOperatorAtBlock(
	arg_id uint32,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"getIDOperator",
		&result,
		arg_id,
	)

	return result, err
}

func (sp *SortitionPool) GetIDOperators(
	arg_ids []uint32,
) ([]common.Address, error) {
	result, err := sp.contract.GetIDOperators(
		sp.callerOptions,
		arg_ids,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"getIDOperators",
			arg_ids,
		)
	}

	return result, err
}

func (sp *SortitionPool) GetIDOperatorsAtBlock(
	arg_ids []uint32,
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"getIDOperators",
		&result,
		arg_ids,
	)

	return result, err
}

func (sp *SortitionPool) GetOperatorID(
	arg_operator common.Address,
) (uint32, error) {
	result, err := sp.contract.GetOperatorID(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"getOperatorID",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) GetOperatorIDAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"getOperatorID",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) GetPoolWeight(
	arg_operator common.Address,
) (*big.Int, error) {
	result, err := sp.contract.GetPoolWeight(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"getPoolWeight",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) GetPoolWeightAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"getPoolWeight",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) IneligibleEarnedRewards() (*big.Int, error) {
	result, err := sp.contract.IneligibleEarnedRewards(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"ineligibleEarnedRewards",
		)
	}

	return result, err
}

func (sp *SortitionPool) IneligibleEarnedRewardsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"ineligibleEarnedRewards",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) IsEligibleForRewards(
	arg_operator uint32,
) (bool, error) {
	result, err := sp.contract.IsEligibleForRewards(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"isEligibleForRewards",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) IsEligibleForRewardsAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"isEligibleForRewards",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) IsLocked() (bool, error) {
	result, err := sp.contract.IsLocked(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"isLocked",
		)
	}

	return result, err
}

func (sp *SortitionPool) IsLockedAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"isLocked",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) IsOperatorInPool(
	arg_operator common.Address,
) (bool, error) {
	result, err := sp.contract.IsOperatorInPool(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"isOperatorInPool",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) IsOperatorInPoolAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"isOperatorInPool",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) IsOperatorRegistered(
	arg_operator common.Address,
) (bool, error) {
	result, err := sp.contract.IsOperatorRegistered(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"isOperatorRegistered",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) IsOperatorRegisteredAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"isOperatorRegistered",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) IsOperatorUpToDate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (bool, error) {
	result, err := sp.contract.IsOperatorUpToDate(
		sp.callerOptions,
		arg_operator,
		arg_authorizedStake,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"isOperatorUpToDate",
			arg_operator,
			arg_authorizedStake,
		)
	}

	return result, err
}

func (sp *SortitionPool) IsOperatorUpToDateAtBlock(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"isOperatorUpToDate",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

func (sp *SortitionPool) OperatorsInPool() (*big.Int, error) {
	result, err := sp.contract.OperatorsInPool(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"operatorsInPool",
		)
	}

	return result, err
}

func (sp *SortitionPool) OperatorsInPoolAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"operatorsInPool",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) Owner() (common.Address, error) {
	result, err := sp.contract.Owner(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"owner",
		)
	}

	return result, err
}

func (sp *SortitionPool) OwnerAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"owner",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) PoolWeightDivisor() (*big.Int, error) {
	result, err := sp.contract.PoolWeightDivisor(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"poolWeightDivisor",
		)
	}

	return result, err
}

func (sp *SortitionPool) PoolWeightDivisorAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"poolWeightDivisor",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) RewardToken() (common.Address, error) {
	result, err := sp.contract.RewardToken(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"rewardToken",
		)
	}

	return result, err
}

func (sp *SortitionPool) RewardTokenAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"rewardToken",
		&result,
	)

	return result, err
}

func (sp *SortitionPool) RewardsEligibilityRestorableAt(
	arg_operator uint32,
) (*big.Int, error) {
	result, err := sp.contract.RewardsEligibilityRestorableAt(
		sp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"rewardsEligibilityRestorableAt",
			arg_operator,
		)
	}

	return result, err
}

func (sp *SortitionPool) RewardsEligibilityRestorableAtAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"rewardsEligibilityRestorableAt",
		&result,
		arg_operator,
	)

	return result, err
}

func (sp *SortitionPool) SelectGroup(
	arg_groupSize *big.Int,
	arg_seed [32]byte,
) ([]uint32, error) {
	result, err := sp.contract.SelectGroup(
		sp.callerOptions,
		arg_groupSize,
		arg_seed,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"selectGroup",
			arg_groupSize,
			arg_seed,
		)
	}

	return result, err
}

func (sp *SortitionPool) SelectGroupAtBlock(
	arg_groupSize *big.Int,
	arg_seed [32]byte,
	blockNumber *big.Int,
) ([]uint32, error) {
	var result []uint32

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"selectGroup",
		&result,
		arg_groupSize,
		arg_seed,
	)

	return result, err
}

func (sp *SortitionPool) TotalWeight() (*big.Int, error) {
	result, err := sp.contract.TotalWeight(
		sp.callerOptions,
	)

	if err != nil {
		return result, sp.errorResolver.ResolveError(
			err,
			sp.callerOptions.From,
			nil,
			"totalWeight",
		)
	}

	return result, err
}

func (sp *SortitionPool) TotalWeightAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		sp.callerOptions.From,
		blockNumber,
		nil,
		sp.contractABI,
		sp.caller,
		sp.errorResolver,
		sp.contractAddress,
		"totalWeight",
		&result,
	)

	return result, err
}

// ------ Events -------

func (sp *SortitionPool) IneligibleForRewardsEvent(
	opts *ethlike.SubscribeOpts,
) *SpIneligibleForRewardsSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &SpIneligibleForRewardsSubscription{
		sp,
		opts,
	}
}

type SpIneligibleForRewardsSubscription struct {
	contract *SortitionPool
	opts     *ethlike.SubscribeOpts
}

type sortitionPoolIneligibleForRewardsFunc func(
	Ids []uint32,
	Until *big.Int,
	blockNumber uint64,
)

func (ifrs *SpIneligibleForRewardsSubscription) OnEvent(
	handler sortitionPoolIneligibleForRewardsFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.SortitionPoolIneligibleForRewards)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Ids,
					event.Until,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ifrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ifrs *SpIneligibleForRewardsSubscription) Pipe(
	sink chan *abi.SortitionPoolIneligibleForRewards,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ifrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ifrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ifrs.opts.PastBlocks

				spLogger.Infof(
					"subscription monitoring fetching past IneligibleForRewards events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ifrs.contract.PastIneligibleForRewardsEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				spLogger.Infof(
					"subscription monitoring fetched [%v] past IneligibleForRewards events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ifrs.contract.watchIneligibleForRewards(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sp *SortitionPool) watchIneligibleForRewards(
	sink chan *abi.SortitionPoolIneligibleForRewards,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return sp.contract.WatchIneligibleForRewards(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		spLogger.Errorf(
			"subscription to event IneligibleForRewards had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		spLogger.Errorf(
			"subscription to event IneligibleForRewards failed "+
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

func (sp *SortitionPool) PastIneligibleForRewardsEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.SortitionPoolIneligibleForRewards, error) {
	iterator, err := sp.contract.FilterIneligibleForRewards(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past IneligibleForRewards events: [%v]",
			err,
		)
	}

	events := make([]*abi.SortitionPoolIneligibleForRewards, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (sp *SortitionPool) OwnershipTransferredEvent(
	opts *ethlike.SubscribeOpts,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) *SpOwnershipTransferredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &SpOwnershipTransferredSubscription{
		sp,
		opts,
		previousOwnerFilter,
		newOwnerFilter,
	}
}

type SpOwnershipTransferredSubscription struct {
	contract            *SortitionPool
	opts                *ethlike.SubscribeOpts
	previousOwnerFilter []common.Address
	newOwnerFilter      []common.Address
}

type sortitionPoolOwnershipTransferredFunc func(
	PreviousOwner common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (ots *SpOwnershipTransferredSubscription) OnEvent(
	handler sortitionPoolOwnershipTransferredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.SortitionPoolOwnershipTransferred)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.PreviousOwner,
					event.NewOwner,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ots.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ots *SpOwnershipTransferredSubscription) Pipe(
	sink chan *abi.SortitionPoolOwnershipTransferred,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ots.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ots.contract.blockCounter.CurrentBlock()
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ots.opts.PastBlocks

				spLogger.Infof(
					"subscription monitoring fetching past OwnershipTransferred events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ots.contract.PastOwnershipTransferredEvents(
					fromBlock,
					nil,
					ots.previousOwnerFilter,
					ots.newOwnerFilter,
				)
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				spLogger.Infof(
					"subscription monitoring fetched [%v] past OwnershipTransferred events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ots.contract.watchOwnershipTransferred(
		sink,
		ots.previousOwnerFilter,
		ots.newOwnerFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sp *SortitionPool) watchOwnershipTransferred(
	sink chan *abi.SortitionPoolOwnershipTransferred,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return sp.contract.WatchOwnershipTransferred(
			&bind.WatchOpts{Context: ctx},
			sink,
			previousOwnerFilter,
			newOwnerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		spLogger.Errorf(
			"subscription to event OwnershipTransferred had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		spLogger.Errorf(
			"subscription to event OwnershipTransferred failed "+
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

func (sp *SortitionPool) PastOwnershipTransferredEvents(
	startBlock uint64,
	endBlock *uint64,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) ([]*abi.SortitionPoolOwnershipTransferred, error) {
	iterator, err := sp.contract.FilterOwnershipTransferred(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		previousOwnerFilter,
		newOwnerFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past OwnershipTransferred events: [%v]",
			err,
		)
	}

	events := make([]*abi.SortitionPoolOwnershipTransferred, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (sp *SortitionPool) RewardEligibilityRestoredEvent(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
	idFilter []uint32,
) *SpRewardEligibilityRestoredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &SpRewardEligibilityRestoredSubscription{
		sp,
		opts,
		operatorFilter,
		idFilter,
	}
}

type SpRewardEligibilityRestoredSubscription struct {
	contract       *SortitionPool
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
	idFilter       []uint32
}

type sortitionPoolRewardEligibilityRestoredFunc func(
	Operator common.Address,
	Id uint32,
	blockNumber uint64,
)

func (rers *SpRewardEligibilityRestoredSubscription) OnEvent(
	handler sortitionPoolRewardEligibilityRestoredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.SortitionPoolRewardEligibilityRestored)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Operator,
					event.Id,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rers.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rers *SpRewardEligibilityRestoredSubscription) Pipe(
	sink chan *abi.SortitionPoolRewardEligibilityRestored,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rers.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rers.contract.blockCounter.CurrentBlock()
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rers.opts.PastBlocks

				spLogger.Infof(
					"subscription monitoring fetching past RewardEligibilityRestored events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rers.contract.PastRewardEligibilityRestoredEvents(
					fromBlock,
					nil,
					rers.operatorFilter,
					rers.idFilter,
				)
				if err != nil {
					spLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				spLogger.Infof(
					"subscription monitoring fetched [%v] past RewardEligibilityRestored events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rers.contract.watchRewardEligibilityRestored(
		sink,
		rers.operatorFilter,
		rers.idFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sp *SortitionPool) watchRewardEligibilityRestored(
	sink chan *abi.SortitionPoolRewardEligibilityRestored,
	operatorFilter []common.Address,
	idFilter []uint32,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return sp.contract.WatchRewardEligibilityRestored(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
			idFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		spLogger.Errorf(
			"subscription to event RewardEligibilityRestored had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		spLogger.Errorf(
			"subscription to event RewardEligibilityRestored failed "+
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

func (sp *SortitionPool) PastRewardEligibilityRestoredEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
	idFilter []uint32,
) ([]*abi.SortitionPoolRewardEligibilityRestored, error) {
	iterator, err := sp.contract.FilterRewardEligibilityRestored(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		operatorFilter,
		idFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RewardEligibilityRestored events: [%v]",
			err,
		)
	}

	events := make([]*abi.SortitionPoolRewardEligibilityRestored, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
