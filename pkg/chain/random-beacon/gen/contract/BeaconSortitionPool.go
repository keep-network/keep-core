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
var bspLogger = log.Logger("keep-contract-BeaconSortitionPool")

type BeaconSortitionPool struct {
	contract          *abi.BeaconSortitionPool
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

func NewBeaconSortitionPool(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*BeaconSortitionPool, error) {
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

	contract, err := abi.NewBeaconSortitionPool(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.BeaconSortitionPoolABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &BeaconSortitionPool{
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
func (bsp *BeaconSortitionPool) InsertOperator(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction insertOperator",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_authorizedStake,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.InsertOperator(
		transactorOptions,
		arg_operator,
		arg_authorizedStake,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"insertOperator",
			arg_operator,
			arg_authorizedStake,
		)
	}

	bspLogger.Infof(
		"submitted transaction insertOperator with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.InsertOperator(
				newTransactorOptions,
				arg_operator,
				arg_authorizedStake,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"insertOperator",
					arg_operator,
					arg_authorizedStake,
				)
			}

			bspLogger.Infof(
				"submitted transaction insertOperator with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallInsertOperator(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"insertOperator",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return err
}

func (bsp *BeaconSortitionPool) InsertOperatorGasEstimate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"insertOperator",
		bsp.contractABI,
		bsp.transactor,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) Lock(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction lock",
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.Lock(
		transactorOptions,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"lock",
		)
	}

	bspLogger.Infof(
		"submitted transaction lock with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.Lock(
				newTransactorOptions,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"lock",
				)
			}

			bspLogger.Infof(
				"submitted transaction lock with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallLock(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"lock",
		&result,
	)

	return err
}

func (bsp *BeaconSortitionPool) LockGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"lock",
		bsp.contractABI,
		bsp.transactor,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) ReceiveApproval(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction receiveApproval",
		" params: ",
		fmt.Sprint(
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.ReceiveApproval(
		transactorOptions,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"receiveApproval",
			arg_sender,
			arg_amount,
			arg_token,
			arg3,
		)
	}

	bspLogger.Infof(
		"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.ReceiveApproval(
				newTransactorOptions,
				arg_sender,
				arg_amount,
				arg_token,
				arg3,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"receiveApproval",
					arg_sender,
					arg_amount,
					arg_token,
					arg3,
				)
			}

			bspLogger.Infof(
				"submitted transaction receiveApproval with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallReceiveApproval(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"receiveApproval",
		&result,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)

	return err
}

func (bsp *BeaconSortitionPool) ReceiveApprovalGasEstimate(
	arg_sender common.Address,
	arg_amount *big.Int,
	arg_token common.Address,
	arg3 []byte,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"receiveApproval",
		bsp.contractABI,
		bsp.transactor,
		arg_sender,
		arg_amount,
		arg_token,
		arg3,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) RenounceOwnership(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction renounceOwnership",
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.RenounceOwnership(
		transactorOptions,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"renounceOwnership",
		)
	}

	bspLogger.Infof(
		"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.RenounceOwnership(
				newTransactorOptions,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"renounceOwnership",
				)
			}

			bspLogger.Infof(
				"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallRenounceOwnership(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"renounceOwnership",
		&result,
	)

	return err
}

func (bsp *BeaconSortitionPool) RenounceOwnershipGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"renounceOwnership",
		bsp.contractABI,
		bsp.transactor,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) RestoreRewardEligibility(
	arg_operator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction restoreRewardEligibility",
		" params: ",
		fmt.Sprint(
			arg_operator,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.RestoreRewardEligibility(
		transactorOptions,
		arg_operator,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"restoreRewardEligibility",
			arg_operator,
		)
	}

	bspLogger.Infof(
		"submitted transaction restoreRewardEligibility with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.RestoreRewardEligibility(
				newTransactorOptions,
				arg_operator,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"restoreRewardEligibility",
					arg_operator,
				)
			}

			bspLogger.Infof(
				"submitted transaction restoreRewardEligibility with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallRestoreRewardEligibility(
	arg_operator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"restoreRewardEligibility",
		&result,
		arg_operator,
	)

	return err
}

func (bsp *BeaconSortitionPool) RestoreRewardEligibilityGasEstimate(
	arg_operator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"restoreRewardEligibility",
		bsp.contractABI,
		bsp.transactor,
		arg_operator,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) SetRewardIneligibility(
	arg_operators []uint32,
	arg_until *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction setRewardIneligibility",
		" params: ",
		fmt.Sprint(
			arg_operators,
			arg_until,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.SetRewardIneligibility(
		transactorOptions,
		arg_operators,
		arg_until,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"setRewardIneligibility",
			arg_operators,
			arg_until,
		)
	}

	bspLogger.Infof(
		"submitted transaction setRewardIneligibility with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.SetRewardIneligibility(
				newTransactorOptions,
				arg_operators,
				arg_until,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"setRewardIneligibility",
					arg_operators,
					arg_until,
				)
			}

			bspLogger.Infof(
				"submitted transaction setRewardIneligibility with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallSetRewardIneligibility(
	arg_operators []uint32,
	arg_until *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"setRewardIneligibility",
		&result,
		arg_operators,
		arg_until,
	)

	return err
}

func (bsp *BeaconSortitionPool) SetRewardIneligibilityGasEstimate(
	arg_operators []uint32,
	arg_until *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"setRewardIneligibility",
		bsp.contractABI,
		bsp.transactor,
		arg_operators,
		arg_until,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) TransferOwnership(
	arg_newOwner common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction transferOwnership",
		" params: ",
		fmt.Sprint(
			arg_newOwner,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.TransferOwnership(
		transactorOptions,
		arg_newOwner,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"transferOwnership",
			arg_newOwner,
		)
	}

	bspLogger.Infof(
		"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.TransferOwnership(
				newTransactorOptions,
				arg_newOwner,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"transferOwnership",
					arg_newOwner,
				)
			}

			bspLogger.Infof(
				"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallTransferOwnership(
	arg_newOwner common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"transferOwnership",
		&result,
		arg_newOwner,
	)

	return err
}

func (bsp *BeaconSortitionPool) TransferOwnershipGasEstimate(
	arg_newOwner common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"transferOwnership",
		bsp.contractABI,
		bsp.transactor,
		arg_newOwner,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) Unlock(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction unlock",
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.Unlock(
		transactorOptions,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"unlock",
		)
	}

	bspLogger.Infof(
		"submitted transaction unlock with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.Unlock(
				newTransactorOptions,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"unlock",
				)
			}

			bspLogger.Infof(
				"submitted transaction unlock with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallUnlock(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"unlock",
		&result,
	)

	return err
}

func (bsp *BeaconSortitionPool) UnlockGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"unlock",
		bsp.contractABI,
		bsp.transactor,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) UpdateOperatorStatus(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction updateOperatorStatus",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_authorizedStake,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.UpdateOperatorStatus(
		transactorOptions,
		arg_operator,
		arg_authorizedStake,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"updateOperatorStatus",
			arg_operator,
			arg_authorizedStake,
		)
	}

	bspLogger.Infof(
		"submitted transaction updateOperatorStatus with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.UpdateOperatorStatus(
				newTransactorOptions,
				arg_operator,
				arg_authorizedStake,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"updateOperatorStatus",
					arg_operator,
					arg_authorizedStake,
				)
			}

			bspLogger.Infof(
				"submitted transaction updateOperatorStatus with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallUpdateOperatorStatus(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"updateOperatorStatus",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return err
}

func (bsp *BeaconSortitionPool) UpdateOperatorStatusGasEstimate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"updateOperatorStatus",
		bsp.contractABI,
		bsp.transactor,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) WithdrawIneligible(
	arg_recipient common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction withdrawIneligible",
		" params: ",
		fmt.Sprint(
			arg_recipient,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.WithdrawIneligible(
		transactorOptions,
		arg_recipient,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"withdrawIneligible",
			arg_recipient,
		)
	}

	bspLogger.Infof(
		"submitted transaction withdrawIneligible with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.WithdrawIneligible(
				newTransactorOptions,
				arg_recipient,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"withdrawIneligible",
					arg_recipient,
				)
			}

			bspLogger.Infof(
				"submitted transaction withdrawIneligible with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallWithdrawIneligible(
	arg_recipient common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"withdrawIneligible",
		&result,
		arg_recipient,
	)

	return err
}

func (bsp *BeaconSortitionPool) WithdrawIneligibleGasEstimate(
	arg_recipient common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"withdrawIneligible",
		bsp.contractABI,
		bsp.transactor,
		arg_recipient,
	)

	return result, err
}

// Transaction submission.
func (bsp *BeaconSortitionPool) WithdrawRewards(
	arg_operator common.Address,
	arg_beneficiary common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	bspLogger.Debug(
		"submitting transaction withdrawRewards",
		" params: ",
		fmt.Sprint(
			arg_operator,
			arg_beneficiary,
		),
	)

	bsp.transactionMutex.Lock()
	defer bsp.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *bsp.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := bsp.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := bsp.contract.WithdrawRewards(
		transactorOptions,
		arg_operator,
		arg_beneficiary,
	)
	if err != nil {
		return transaction, bsp.errorResolver.ResolveError(
			err,
			bsp.transactorOptions.From,
			nil,
			"withdrawRewards",
			arg_operator,
			arg_beneficiary,
		)
	}

	bspLogger.Infof(
		"submitted transaction withdrawRewards with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go bsp.miningWaiter.ForceMining(
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

			transaction, err := bsp.contract.WithdrawRewards(
				newTransactorOptions,
				arg_operator,
				arg_beneficiary,
			)
			if err != nil {
				return nil, bsp.errorResolver.ResolveError(
					err,
					bsp.transactorOptions.From,
					nil,
					"withdrawRewards",
					arg_operator,
					arg_beneficiary,
				)
			}

			bspLogger.Infof(
				"submitted transaction withdrawRewards with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	bsp.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (bsp *BeaconSortitionPool) CallWithdrawRewards(
	arg_operator common.Address,
	arg_beneficiary common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.transactorOptions.From,
		blockNumber, nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"withdrawRewards",
		&result,
		arg_operator,
		arg_beneficiary,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) WithdrawRewardsGasEstimate(
	arg_operator common.Address,
	arg_beneficiary common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		bsp.callerOptions.From,
		bsp.contractAddress,
		"withdrawRewards",
		bsp.contractABI,
		bsp.transactor,
		arg_operator,
		arg_beneficiary,
	)

	return result, err
}

// ----- Const Methods ------

func (bsp *BeaconSortitionPool) CanRestoreRewardEligibility(
	arg_operator uint32,
) (bool, error) {
	result, err := bsp.contract.CanRestoreRewardEligibility(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"canRestoreRewardEligibility",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) CanRestoreRewardEligibilityAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"canRestoreRewardEligibility",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) GetAvailableRewards(
	arg_operator common.Address,
) (*big.Int, error) {
	result, err := bsp.contract.GetAvailableRewards(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"getAvailableRewards",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) GetAvailableRewardsAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"getAvailableRewards",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) GetIDOperator(
	arg_id uint32,
) (common.Address, error) {
	result, err := bsp.contract.GetIDOperator(
		bsp.callerOptions,
		arg_id,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"getIDOperator",
			arg_id,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) GetIDOperatorAtBlock(
	arg_id uint32,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"getIDOperator",
		&result,
		arg_id,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) GetIDOperators(
	arg_ids []uint32,
) ([]common.Address, error) {
	result, err := bsp.contract.GetIDOperators(
		bsp.callerOptions,
		arg_ids,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"getIDOperators",
			arg_ids,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) GetIDOperatorsAtBlock(
	arg_ids []uint32,
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"getIDOperators",
		&result,
		arg_ids,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) GetOperatorID(
	arg_operator common.Address,
) (uint32, error) {
	result, err := bsp.contract.GetOperatorID(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"getOperatorID",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) GetOperatorIDAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"getOperatorID",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) GetPoolWeight(
	arg_operator common.Address,
) (*big.Int, error) {
	result, err := bsp.contract.GetPoolWeight(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"getPoolWeight",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) GetPoolWeightAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"getPoolWeight",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IneligibleEarnedRewards() (*big.Int, error) {
	result, err := bsp.contract.IneligibleEarnedRewards(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"ineligibleEarnedRewards",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IneligibleEarnedRewardsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"ineligibleEarnedRewards",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IsEligibleForRewards(
	arg_operator uint32,
) (bool, error) {
	result, err := bsp.contract.IsEligibleForRewards(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"isEligibleForRewards",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IsEligibleForRewardsAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"isEligibleForRewards",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IsLocked() (bool, error) {
	result, err := bsp.contract.IsLocked(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"isLocked",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IsLockedAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"isLocked",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorInPool(
	arg_operator common.Address,
) (bool, error) {
	result, err := bsp.contract.IsOperatorInPool(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"isOperatorInPool",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorInPoolAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"isOperatorInPool",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorRegistered(
	arg_operator common.Address,
) (bool, error) {
	result, err := bsp.contract.IsOperatorRegistered(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"isOperatorRegistered",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorRegisteredAtBlock(
	arg_operator common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"isOperatorRegistered",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorUpToDate(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
) (bool, error) {
	result, err := bsp.contract.IsOperatorUpToDate(
		bsp.callerOptions,
		arg_operator,
		arg_authorizedStake,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"isOperatorUpToDate",
			arg_operator,
			arg_authorizedStake,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) IsOperatorUpToDateAtBlock(
	arg_operator common.Address,
	arg_authorizedStake *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"isOperatorUpToDate",
		&result,
		arg_operator,
		arg_authorizedStake,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) OperatorsInPool() (*big.Int, error) {
	result, err := bsp.contract.OperatorsInPool(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"operatorsInPool",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) OperatorsInPoolAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"operatorsInPool",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) Owner() (common.Address, error) {
	result, err := bsp.contract.Owner(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"owner",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) OwnerAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"owner",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) PoolWeightDivisor() (*big.Int, error) {
	result, err := bsp.contract.PoolWeightDivisor(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"poolWeightDivisor",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) PoolWeightDivisorAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"poolWeightDivisor",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) RewardToken() (common.Address, error) {
	result, err := bsp.contract.RewardToken(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"rewardToken",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) RewardTokenAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"rewardToken",
		&result,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) RewardsEligibilityRestorableAt(
	arg_operator uint32,
) (*big.Int, error) {
	result, err := bsp.contract.RewardsEligibilityRestorableAt(
		bsp.callerOptions,
		arg_operator,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"rewardsEligibilityRestorableAt",
			arg_operator,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) RewardsEligibilityRestorableAtAtBlock(
	arg_operator uint32,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"rewardsEligibilityRestorableAt",
		&result,
		arg_operator,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) SelectGroup(
	arg_groupSize *big.Int,
	arg_seed [32]byte,
) ([]uint32, error) {
	result, err := bsp.contract.SelectGroup(
		bsp.callerOptions,
		arg_groupSize,
		arg_seed,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"selectGroup",
			arg_groupSize,
			arg_seed,
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) SelectGroupAtBlock(
	arg_groupSize *big.Int,
	arg_seed [32]byte,
	blockNumber *big.Int,
) ([]uint32, error) {
	var result []uint32

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"selectGroup",
		&result,
		arg_groupSize,
		arg_seed,
	)

	return result, err
}

func (bsp *BeaconSortitionPool) TotalWeight() (*big.Int, error) {
	result, err := bsp.contract.TotalWeight(
		bsp.callerOptions,
	)

	if err != nil {
		return result, bsp.errorResolver.ResolveError(
			err,
			bsp.callerOptions.From,
			nil,
			"totalWeight",
		)
	}

	return result, err
}

func (bsp *BeaconSortitionPool) TotalWeightAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		bsp.callerOptions.From,
		blockNumber,
		nil,
		bsp.contractABI,
		bsp.caller,
		bsp.errorResolver,
		bsp.contractAddress,
		"totalWeight",
		&result,
	)

	return result, err
}

// ------ Events -------

func (bsp *BeaconSortitionPool) IneligibleForRewardsEvent(
	opts *ethlike.SubscribeOpts,
) *BspIneligibleForRewardsSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &BspIneligibleForRewardsSubscription{
		bsp,
		opts,
	}
}

type BspIneligibleForRewardsSubscription struct {
	contract *BeaconSortitionPool
	opts     *ethlike.SubscribeOpts
}

type beaconSortitionPoolIneligibleForRewardsFunc func(
	Ids []uint32,
	Until *big.Int,
	blockNumber uint64,
)

func (ifrs *BspIneligibleForRewardsSubscription) OnEvent(
	handler beaconSortitionPoolIneligibleForRewardsFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.BeaconSortitionPoolIneligibleForRewards)
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

func (ifrs *BspIneligibleForRewardsSubscription) Pipe(
	sink chan *abi.BeaconSortitionPoolIneligibleForRewards,
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
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ifrs.opts.PastBlocks

				bspLogger.Infof(
					"subscription monitoring fetching past IneligibleForRewards events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ifrs.contract.PastIneligibleForRewardsEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				bspLogger.Infof(
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

func (bsp *BeaconSortitionPool) watchIneligibleForRewards(
	sink chan *abi.BeaconSortitionPoolIneligibleForRewards,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return bsp.contract.WatchIneligibleForRewards(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		bspLogger.Errorf(
			"subscription to event IneligibleForRewards had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		bspLogger.Errorf(
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

func (bsp *BeaconSortitionPool) PastIneligibleForRewardsEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.BeaconSortitionPoolIneligibleForRewards, error) {
	iterator, err := bsp.contract.FilterIneligibleForRewards(
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

	events := make([]*abi.BeaconSortitionPoolIneligibleForRewards, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (bsp *BeaconSortitionPool) OwnershipTransferredEvent(
	opts *ethlike.SubscribeOpts,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) *BspOwnershipTransferredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &BspOwnershipTransferredSubscription{
		bsp,
		opts,
		previousOwnerFilter,
		newOwnerFilter,
	}
}

type BspOwnershipTransferredSubscription struct {
	contract            *BeaconSortitionPool
	opts                *ethlike.SubscribeOpts
	previousOwnerFilter []common.Address
	newOwnerFilter      []common.Address
}

type beaconSortitionPoolOwnershipTransferredFunc func(
	PreviousOwner common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (ots *BspOwnershipTransferredSubscription) OnEvent(
	handler beaconSortitionPoolOwnershipTransferredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.BeaconSortitionPoolOwnershipTransferred)
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

func (ots *BspOwnershipTransferredSubscription) Pipe(
	sink chan *abi.BeaconSortitionPoolOwnershipTransferred,
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
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ots.opts.PastBlocks

				bspLogger.Infof(
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
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				bspLogger.Infof(
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

func (bsp *BeaconSortitionPool) watchOwnershipTransferred(
	sink chan *abi.BeaconSortitionPoolOwnershipTransferred,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return bsp.contract.WatchOwnershipTransferred(
			&bind.WatchOpts{Context: ctx},
			sink,
			previousOwnerFilter,
			newOwnerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		bspLogger.Errorf(
			"subscription to event OwnershipTransferred had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		bspLogger.Errorf(
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

func (bsp *BeaconSortitionPool) PastOwnershipTransferredEvents(
	startBlock uint64,
	endBlock *uint64,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) ([]*abi.BeaconSortitionPoolOwnershipTransferred, error) {
	iterator, err := bsp.contract.FilterOwnershipTransferred(
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

	events := make([]*abi.BeaconSortitionPoolOwnershipTransferred, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (bsp *BeaconSortitionPool) RewardEligibilityRestoredEvent(
	opts *ethlike.SubscribeOpts,
	operatorFilter []common.Address,
	idFilter []uint32,
) *BspRewardEligibilityRestoredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &BspRewardEligibilityRestoredSubscription{
		bsp,
		opts,
		operatorFilter,
		idFilter,
	}
}

type BspRewardEligibilityRestoredSubscription struct {
	contract       *BeaconSortitionPool
	opts           *ethlike.SubscribeOpts
	operatorFilter []common.Address
	idFilter       []uint32
}

type beaconSortitionPoolRewardEligibilityRestoredFunc func(
	Operator common.Address,
	Id uint32,
	blockNumber uint64,
)

func (rers *BspRewardEligibilityRestoredSubscription) OnEvent(
	handler beaconSortitionPoolRewardEligibilityRestoredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.BeaconSortitionPoolRewardEligibilityRestored)
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

func (rers *BspRewardEligibilityRestoredSubscription) Pipe(
	sink chan *abi.BeaconSortitionPoolRewardEligibilityRestored,
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
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rers.opts.PastBlocks

				bspLogger.Infof(
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
					bspLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				bspLogger.Infof(
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

func (bsp *BeaconSortitionPool) watchRewardEligibilityRestored(
	sink chan *abi.BeaconSortitionPoolRewardEligibilityRestored,
	operatorFilter []common.Address,
	idFilter []uint32,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return bsp.contract.WatchRewardEligibilityRestored(
			&bind.WatchOpts{Context: ctx},
			sink,
			operatorFilter,
			idFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		bspLogger.Errorf(
			"subscription to event RewardEligibilityRestored had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		bspLogger.Errorf(
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

func (bsp *BeaconSortitionPool) PastRewardEligibilityRestoredEvents(
	startBlock uint64,
	endBlock *uint64,
	operatorFilter []common.Address,
	idFilter []uint32,
) ([]*abi.BeaconSortitionPoolRewardEligibilityRestored, error) {
	iterator, err := bsp.contract.FilterRewardEligibilityRestored(
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

	events := make([]*abi.BeaconSortitionPoolRewardEligibilityRestored, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
