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
var krboLogger = log.Logger("keep-contract-KeepRandomBeaconOperator")

type KeepRandomBeaconOperator struct {
	contract          *abi.KeepRandomBeaconOperator
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

func NewKeepRandomBeaconOperator(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*KeepRandomBeaconOperator, error) {
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

	contract, err := abi.NewKeepRandomBeaconOperator(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.KeepRandomBeaconOperatorABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &KeepRandomBeaconOperator{
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
func (krbo *KeepRandomBeaconOperator) AddServiceContract(
	serviceContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction addServiceContract",
		" params: ",
		fmt.Sprint(
			serviceContract,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.AddServiceContract(
		transactorOptions,
		serviceContract,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"addServiceContract",
			serviceContract,
		)
	}

	krboLogger.Infof(
		"submitted transaction addServiceContract with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.AddServiceContract(
				newTransactorOptions,
				serviceContract,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"addServiceContract",
					serviceContract,
				)
			}

			krboLogger.Infof(
				"submitted transaction addServiceContract with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallAddServiceContract(
	serviceContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"addServiceContract",
		&result,
		serviceContract,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) AddServiceContractGasEstimate(
	serviceContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"addServiceContract",
		krbo.contractABI,
		krbo.transactor,
		serviceContract,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) CreateGroup(
	_newEntry *big.Int,
	submitter common.Address,
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction createGroup",
		" params: ",
		fmt.Sprint(
			_newEntry,
			submitter,
		),
		" value: ", value,
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.CreateGroup(
		transactorOptions,
		_newEntry,
		submitter,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			value,
			"createGroup",
			_newEntry,
			submitter,
		)
	}

	krboLogger.Infof(
		"submitted transaction createGroup with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.CreateGroup(
				newTransactorOptions,
				_newEntry,
				submitter,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					value,
					"createGroup",
					_newEntry,
					submitter,
				)
			}

			krboLogger.Infof(
				"submitted transaction createGroup with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallCreateGroup(
	_newEntry *big.Int,
	submitter common.Address,
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, value,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"createGroup",
		&result,
		_newEntry,
		submitter,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) CreateGroupGasEstimate(
	_newEntry *big.Int,
	submitter common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"createGroup",
		krbo.contractABI,
		krbo.transactor,
		_newEntry,
		submitter,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) Genesis(
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction genesis",
		" value: ", value,
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.Genesis(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			value,
			"genesis",
		)
	}

	krboLogger.Infof(
		"submitted transaction genesis with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.Genesis(
				newTransactorOptions,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					value,
					"genesis",
				)
			}

			krboLogger.Infof(
				"submitted transaction genesis with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallGenesis(
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, value,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"genesis",
		&result,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) GenesisGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"genesis",
		krbo.contractABI,
		krbo.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) RefreshGasPrice(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction refreshGasPrice",
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.RefreshGasPrice(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"refreshGasPrice",
		)
	}

	krboLogger.Infof(
		"submitted transaction refreshGasPrice with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.RefreshGasPrice(
				newTransactorOptions,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"refreshGasPrice",
				)
			}

			krboLogger.Infof(
				"submitted transaction refreshGasPrice with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallRefreshGasPrice(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"refreshGasPrice",
		&result,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) RefreshGasPriceGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"refreshGasPrice",
		krbo.contractABI,
		krbo.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) RelayEntry(
	_groupSignature []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction relayEntry",
		" params: ",
		fmt.Sprint(
			_groupSignature,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.RelayEntry(
		transactorOptions,
		_groupSignature,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"relayEntry",
			_groupSignature,
		)
	}

	krboLogger.Infof(
		"submitted transaction relayEntry with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.RelayEntry(
				newTransactorOptions,
				_groupSignature,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"relayEntry",
					_groupSignature,
				)
			}

			krboLogger.Infof(
				"submitted transaction relayEntry with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallRelayEntry(
	_groupSignature []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"relayEntry",
		&result,
		_groupSignature,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) RelayEntryGasEstimate(
	_groupSignature []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"relayEntry",
		krbo.contractABI,
		krbo.transactor,
		_groupSignature,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) ReportRelayEntryTimeout(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction reportRelayEntryTimeout",
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.ReportRelayEntryTimeout(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"reportRelayEntryTimeout",
		)
	}

	krboLogger.Infof(
		"submitted transaction reportRelayEntryTimeout with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.ReportRelayEntryTimeout(
				newTransactorOptions,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"reportRelayEntryTimeout",
				)
			}

			krboLogger.Infof(
				"submitted transaction reportRelayEntryTimeout with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallReportRelayEntryTimeout(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"reportRelayEntryTimeout",
		&result,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) ReportRelayEntryTimeoutGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"reportRelayEntryTimeout",
		krbo.contractABI,
		krbo.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) ReportUnauthorizedSigning(
	groupIndex *big.Int,
	signedMsgSender []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction reportUnauthorizedSigning",
		" params: ",
		fmt.Sprint(
			groupIndex,
			signedMsgSender,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.ReportUnauthorizedSigning(
		transactorOptions,
		groupIndex,
		signedMsgSender,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"reportUnauthorizedSigning",
			groupIndex,
			signedMsgSender,
		)
	}

	krboLogger.Infof(
		"submitted transaction reportUnauthorizedSigning with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.ReportUnauthorizedSigning(
				newTransactorOptions,
				groupIndex,
				signedMsgSender,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"reportUnauthorizedSigning",
					groupIndex,
					signedMsgSender,
				)
			}

			krboLogger.Infof(
				"submitted transaction reportUnauthorizedSigning with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallReportUnauthorizedSigning(
	groupIndex *big.Int,
	signedMsgSender []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"reportUnauthorizedSigning",
		&result,
		groupIndex,
		signedMsgSender,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) ReportUnauthorizedSigningGasEstimate(
	groupIndex *big.Int,
	signedMsgSender []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"reportUnauthorizedSigning",
		krbo.contractABI,
		krbo.transactor,
		groupIndex,
		signedMsgSender,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) Sign(
	requestId *big.Int,
	previousEntry []uint8,
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction sign",
		" params: ",
		fmt.Sprint(
			requestId,
			previousEntry,
		),
		" value: ", value,
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.Sign(
		transactorOptions,
		requestId,
		previousEntry,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			value,
			"sign",
			requestId,
			previousEntry,
		)
	}

	krboLogger.Infof(
		"submitted transaction sign with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.Sign(
				newTransactorOptions,
				requestId,
				previousEntry,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					value,
					"sign",
					requestId,
					previousEntry,
				)
			}

			krboLogger.Infof(
				"submitted transaction sign with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallSign(
	requestId *big.Int,
	previousEntry []uint8,
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, value,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"sign",
		&result,
		requestId,
		previousEntry,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) SignGasEstimate(
	requestId *big.Int,
	previousEntry []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"sign",
		krbo.contractABI,
		krbo.transactor,
		requestId,
		previousEntry,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) SubmitDkgResult(
	submitterMemberIndex *big.Int,
	groupPubKey []uint8,
	misbehaved []uint8,
	signatures []uint8,
	signingMembersIndexes []*big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction submitDkgResult",
		" params: ",
		fmt.Sprint(
			submitterMemberIndex,
			groupPubKey,
			misbehaved,
			signatures,
			signingMembersIndexes,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.SubmitDkgResult(
		transactorOptions,
		submitterMemberIndex,
		groupPubKey,
		misbehaved,
		signatures,
		signingMembersIndexes,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"submitDkgResult",
			submitterMemberIndex,
			groupPubKey,
			misbehaved,
			signatures,
			signingMembersIndexes,
		)
	}

	krboLogger.Infof(
		"submitted transaction submitDkgResult with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.SubmitDkgResult(
				newTransactorOptions,
				submitterMemberIndex,
				groupPubKey,
				misbehaved,
				signatures,
				signingMembersIndexes,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"submitDkgResult",
					submitterMemberIndex,
					groupPubKey,
					misbehaved,
					signatures,
					signingMembersIndexes,
				)
			}

			krboLogger.Infof(
				"submitted transaction submitDkgResult with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallSubmitDkgResult(
	submitterMemberIndex *big.Int,
	groupPubKey []uint8,
	misbehaved []uint8,
	signatures []uint8,
	signingMembersIndexes []*big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"submitDkgResult",
		&result,
		submitterMemberIndex,
		groupPubKey,
		misbehaved,
		signatures,
		signingMembersIndexes,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) SubmitDkgResultGasEstimate(
	submitterMemberIndex *big.Int,
	groupPubKey []uint8,
	misbehaved []uint8,
	signatures []uint8,
	signingMembersIndexes []*big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"submitDkgResult",
		krbo.contractABI,
		krbo.transactor,
		submitterMemberIndex,
		groupPubKey,
		misbehaved,
		signatures,
		signingMembersIndexes,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) SubmitTicket(
	ticket [32]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction submitTicket",
		" params: ",
		fmt.Sprint(
			ticket,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.SubmitTicket(
		transactorOptions,
		ticket,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"submitTicket",
			ticket,
		)
	}

	krboLogger.Infof(
		"submitted transaction submitTicket with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.SubmitTicket(
				newTransactorOptions,
				ticket,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"submitTicket",
					ticket,
				)
			}

			krboLogger.Infof(
				"submitted transaction submitTicket with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallSubmitTicket(
	ticket [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"submitTicket",
		&result,
		ticket,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) SubmitTicketGasEstimate(
	ticket [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"submitTicket",
		krbo.contractABI,
		krbo.transactor,
		ticket,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) WithdrawGroupMemberRewards(
	operator common.Address,
	groupIndex *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction withdrawGroupMemberRewards",
		" params: ",
		fmt.Sprint(
			operator,
			groupIndex,
		),
	)

	krbo.transactionMutex.Lock()
	defer krbo.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbo.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbo.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbo.contract.WithdrawGroupMemberRewards(
		transactorOptions,
		operator,
		groupIndex,
	)
	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"withdrawGroupMemberRewards",
			operator,
			groupIndex,
		)
	}

	krboLogger.Infof(
		"submitted transaction withdrawGroupMemberRewards with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbo.miningWaiter.ForceMining(
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

			transaction, err := krbo.contract.WithdrawGroupMemberRewards(
				newTransactorOptions,
				operator,
				groupIndex,
			)
			if err != nil {
				return nil, krbo.errorResolver.ResolveError(
					err,
					krbo.transactorOptions.From,
					nil,
					"withdrawGroupMemberRewards",
					operator,
					groupIndex,
				)
			}

			krboLogger.Infof(
				"submitted transaction withdrawGroupMemberRewards with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	krbo.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallWithdrawGroupMemberRewards(
	operator common.Address,
	groupIndex *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"withdrawGroupMemberRewards",
		&result,
		operator,
		groupIndex,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) WithdrawGroupMemberRewardsGasEstimate(
	operator common.Address,
	groupIndex *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"withdrawGroupMemberRewards",
		krbo.contractABI,
		krbo.transactor,
		operator,
		groupIndex,
	)

	return result, err
}

// ----- Const Methods ------

func (krbo *KeepRandomBeaconOperator) CurrentRequestGroupIndex() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.CurrentRequestGroupIndex(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"currentRequestGroupIndex",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentRequestGroupIndexAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"currentRequestGroupIndex",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentRequestPreviousEntry() ([]uint8, error) {
	var result []uint8
	result, err := krbo.contract.CurrentRequestPreviousEntry(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"currentRequestPreviousEntry",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentRequestPreviousEntryAtBlock(
	blockNumber *big.Int,
) ([]uint8, error) {
	var result []uint8

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"currentRequestPreviousEntry",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentRequestStartBlock() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.CurrentRequestStartBlock(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"currentRequestStartBlock",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentRequestStartBlockAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"currentRequestStartBlock",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) DkgGasEstimate() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.DkgGasEstimate(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"dkgGasEstimate",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) DkgGasEstimateAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"dkgGasEstimate",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) DkgSubmitterReimbursementFee() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.DkgSubmitterReimbursementFee(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"dkgSubmitterReimbursementFee",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) DkgSubmitterReimbursementFeeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"dkgSubmitterReimbursementFee",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) EntryVerificationFee() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.EntryVerificationFee(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"entryVerificationFee",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) EntryVerificationFeeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"entryVerificationFee",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) EntryVerificationGasEstimate() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.EntryVerificationGasEstimate(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"entryVerificationGasEstimate",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) EntryVerificationGasEstimateAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"entryVerificationGasEstimate",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GasPriceCeiling() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GasPriceCeiling(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"gasPriceCeiling",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GasPriceCeilingAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"gasPriceCeiling",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetFirstActiveGroupIndex() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GetFirstActiveGroupIndex(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getFirstActiveGroupIndex",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetFirstActiveGroupIndexAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getFirstActiveGroupIndex",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupMemberRewards(
	groupPubKey []uint8,
) (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GetGroupMemberRewards(
		krbo.callerOptions,
		groupPubKey,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getGroupMemberRewards",
			groupPubKey,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupMemberRewardsAtBlock(
	groupPubKey []uint8,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getGroupMemberRewards",
		&result,
		groupPubKey,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupMembers(
	groupPubKey []uint8,
) ([]common.Address, error) {
	var result []common.Address
	result, err := krbo.contract.GetGroupMembers(
		krbo.callerOptions,
		groupPubKey,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getGroupMembers",
			groupPubKey,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupMembersAtBlock(
	groupPubKey []uint8,
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getGroupMembers",
		&result,
		groupPubKey,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupPublicKey(
	groupIndex *big.Int,
) ([]uint8, error) {
	var result []uint8
	result, err := krbo.contract.GetGroupPublicKey(
		krbo.callerOptions,
		groupIndex,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getGroupPublicKey",
			groupIndex,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupPublicKeyAtBlock(
	groupIndex *big.Int,
	blockNumber *big.Int,
) ([]uint8, error) {
	var result []uint8

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getGroupPublicKey",
		&result,
		groupIndex,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupRegistrationTime(
	groupIndex *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GetGroupRegistrationTime(
		krbo.callerOptions,
		groupIndex,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getGroupRegistrationTime",
			groupIndex,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetGroupRegistrationTimeAtBlock(
	groupIndex *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getGroupRegistrationTime",
		&result,
		groupIndex,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetNumberOfCreatedGroups() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GetNumberOfCreatedGroups(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"getNumberOfCreatedGroups",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GetNumberOfCreatedGroupsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"getNumberOfCreatedGroups",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupCreationFee() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupCreationFee(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupCreationFee",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupCreationFeeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupCreationFee",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupMemberBaseReward() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupMemberBaseReward(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupMemberBaseReward",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupMemberBaseRewardAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupMemberBaseReward",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupProfitFee() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupProfitFee(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupProfitFee",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupProfitFeeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupProfitFee",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupSelectionGasEstimate() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupSelectionGasEstimate(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupSelectionGasEstimate",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupSelectionGasEstimateAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupSelectionGasEstimate",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupSize() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupSize(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupSize",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupSizeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupSize",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupThreshold() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.GroupThreshold(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"groupThreshold",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) GroupThresholdAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"groupThreshold",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) HasMinimumStake(
	staker common.Address,
) (bool, error) {
	var result bool
	result, err := krbo.contract.HasMinimumStake(
		krbo.callerOptions,
		staker,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"hasMinimumStake",
			staker,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) HasMinimumStakeAtBlock(
	staker common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"hasMinimumStake",
		&result,
		staker,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) HasWithdrawnRewards(
	operator common.Address,
	groupIndex *big.Int,
) (bool, error) {
	var result bool
	result, err := krbo.contract.HasWithdrawnRewards(
		krbo.callerOptions,
		operator,
		groupIndex,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"hasWithdrawnRewards",
			operator,
			groupIndex,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) HasWithdrawnRewardsAtBlock(
	operator common.Address,
	groupIndex *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"hasWithdrawnRewards",
		&result,
		operator,
		groupIndex,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsEntryInProgress() (bool, error) {
	var result bool
	result, err := krbo.contract.IsEntryInProgress(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"isEntryInProgress",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsEntryInProgressAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"isEntryInProgress",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupRegistered(
	groupPubKey []uint8,
) (bool, error) {
	var result bool
	result, err := krbo.contract.IsGroupRegistered(
		krbo.callerOptions,
		groupPubKey,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"isGroupRegistered",
			groupPubKey,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupRegisteredAtBlock(
	groupPubKey []uint8,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"isGroupRegistered",
		&result,
		groupPubKey,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupSelectionPossible() (bool, error) {
	var result bool
	result, err := krbo.contract.IsGroupSelectionPossible(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"isGroupSelectionPossible",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupSelectionPossibleAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"isGroupSelectionPossible",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupTerminated(
	groupIndex *big.Int,
) (bool, error) {
	var result bool
	result, err := krbo.contract.IsGroupTerminated(
		krbo.callerOptions,
		groupIndex,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"isGroupTerminated",
			groupIndex,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsGroupTerminatedAtBlock(
	groupIndex *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"isGroupTerminated",
		&result,
		groupIndex,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsStaleGroup(
	groupPubKey []uint8,
) (bool, error) {
	var result bool
	result, err := krbo.contract.IsStaleGroup(
		krbo.callerOptions,
		groupPubKey,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"isStaleGroup",
			groupPubKey,
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) IsStaleGroupAtBlock(
	groupPubKey []uint8,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"isStaleGroup",
		&result,
		groupPubKey,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) NumberOfGroups() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.NumberOfGroups(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"numberOfGroups",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) NumberOfGroupsAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"numberOfGroups",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) RelayEntryTimeout() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.RelayEntryTimeout(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"relayEntryTimeout",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) RelayEntryTimeoutAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"relayEntryTimeout",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) ResultPublicationBlockStep() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.ResultPublicationBlockStep(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"resultPublicationBlockStep",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) ResultPublicationBlockStepAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"resultPublicationBlockStep",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) SelectedParticipants() ([]common.Address, error) {
	var result []common.Address
	result, err := krbo.contract.SelectedParticipants(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"selectedParticipants",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) SelectedParticipantsAtBlock(
	blockNumber *big.Int,
) ([]common.Address, error) {
	var result []common.Address

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"selectedParticipants",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) SubmittedTickets() ([]uint64, error) {
	var result []uint64
	result, err := krbo.contract.SubmittedTickets(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"submittedTickets",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) SubmittedTicketsAtBlock(
	blockNumber *big.Int,
) ([]uint64, error) {
	var result []uint64

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"submittedTickets",
		&result,
	)

	return result, err
}

func (krbo *KeepRandomBeaconOperator) TicketSubmissionTimeout() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.TicketSubmissionTimeout(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"ticketSubmissionTimeout",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) TicketSubmissionTimeoutAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"ticketSubmissionTimeout",
		&result,
	)

	return result, err
}

// ------ Events -------

func (krbo *KeepRandomBeaconOperator) DkgResultSubmittedEvent(
	opts *ethlike.SubscribeOpts,
) *KrboDkgResultSubmittedEventSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboDkgResultSubmittedEventSubscription{
		krbo,
		opts,
	}
}

type KrboDkgResultSubmittedEventSubscription struct {
	contract *KeepRandomBeaconOperator
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconOperatorDkgResultSubmittedEventFunc func(
	MemberIndex *big.Int,
	GroupPubKey []uint8,
	Misbehaved []uint8,
	blockNumber uint64,
)

func (drses *KrboDkgResultSubmittedEventSubscription) OnEvent(
	handler keepRandomBeaconOperatorDkgResultSubmittedEventFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.MemberIndex,
					event.GroupPubKey,
					event.Misbehaved,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := drses.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (drses *KrboDkgResultSubmittedEventSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(drses.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := drses.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - drses.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past DkgResultSubmittedEvent events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := drses.contract.PastDkgResultSubmittedEventEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past DkgResultSubmittedEvent events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := drses.contract.watchDkgResultSubmittedEvent(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchDkgResultSubmittedEvent(
	sink chan *abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchDkgResultSubmittedEvent(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event DkgResultSubmittedEvent had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event DkgResultSubmittedEvent failed "+
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

func (krbo *KeepRandomBeaconOperator) PastDkgResultSubmittedEventEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent, error) {
	iterator, err := krbo.contract.FilterDkgResultSubmittedEvent(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DkgResultSubmittedEvent events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) GroupMemberRewardsWithdrawn(
	opts *ethlike.SubscribeOpts,
	beneficiaryFilter []common.Address,
) *KrboGroupMemberRewardsWithdrawnSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboGroupMemberRewardsWithdrawnSubscription{
		krbo,
		opts,
		beneficiaryFilter,
	}
}

type KrboGroupMemberRewardsWithdrawnSubscription struct {
	contract          *KeepRandomBeaconOperator
	opts              *ethlike.SubscribeOpts
	beneficiaryFilter []common.Address
}

type keepRandomBeaconOperatorGroupMemberRewardsWithdrawnFunc func(
	Beneficiary common.Address,
	Operator common.Address,
	Amount *big.Int,
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (gmrws *KrboGroupMemberRewardsWithdrawnSubscription) OnEvent(
	handler keepRandomBeaconOperatorGroupMemberRewardsWithdrawnFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Beneficiary,
					event.Operator,
					event.Amount,
					event.GroupIndex,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := gmrws.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (gmrws *KrboGroupMemberRewardsWithdrawnSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(gmrws.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := gmrws.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - gmrws.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past GroupMemberRewardsWithdrawn events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := gmrws.contract.PastGroupMemberRewardsWithdrawnEvents(
					fromBlock,
					nil,
					gmrws.beneficiaryFilter,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past GroupMemberRewardsWithdrawn events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := gmrws.contract.watchGroupMemberRewardsWithdrawn(
		sink,
		gmrws.beneficiaryFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchGroupMemberRewardsWithdrawn(
	sink chan *abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn,
	beneficiaryFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchGroupMemberRewardsWithdrawn(
			&bind.WatchOpts{Context: ctx},
			sink,
			beneficiaryFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event GroupMemberRewardsWithdrawn had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event GroupMemberRewardsWithdrawn failed "+
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

func (krbo *KeepRandomBeaconOperator) PastGroupMemberRewardsWithdrawnEvents(
	startBlock uint64,
	endBlock *uint64,
	beneficiaryFilter []common.Address,
) ([]*abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn, error) {
	iterator, err := krbo.contract.FilterGroupMemberRewardsWithdrawn(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		beneficiaryFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past GroupMemberRewardsWithdrawn events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) GroupSelectionStarted(
	opts *ethlike.SubscribeOpts,
) *KrboGroupSelectionStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboGroupSelectionStartedSubscription{
		krbo,
		opts,
	}
}

type KrboGroupSelectionStartedSubscription struct {
	contract *KeepRandomBeaconOperator
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconOperatorGroupSelectionStartedFunc func(
	NewEntry *big.Int,
	blockNumber uint64,
)

func (gsss *KrboGroupSelectionStartedSubscription) OnEvent(
	handler keepRandomBeaconOperatorGroupSelectionStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorGroupSelectionStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.NewEntry,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := gsss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (gsss *KrboGroupSelectionStartedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorGroupSelectionStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(gsss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := gsss.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - gsss.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past GroupSelectionStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := gsss.contract.PastGroupSelectionStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past GroupSelectionStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := gsss.contract.watchGroupSelectionStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchGroupSelectionStarted(
	sink chan *abi.KeepRandomBeaconOperatorGroupSelectionStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchGroupSelectionStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event GroupSelectionStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event GroupSelectionStarted failed "+
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

func (krbo *KeepRandomBeaconOperator) PastGroupSelectionStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconOperatorGroupSelectionStarted, error) {
	iterator, err := krbo.contract.FilterGroupSelectionStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past GroupSelectionStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorGroupSelectionStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) OnGroupRegistered(
	opts *ethlike.SubscribeOpts,
) *KrboOnGroupRegisteredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboOnGroupRegisteredSubscription{
		krbo,
		opts,
	}
}

type KrboOnGroupRegisteredSubscription struct {
	contract *KeepRandomBeaconOperator
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconOperatorOnGroupRegisteredFunc func(
	GroupPubKey []uint8,
	blockNumber uint64,
)

func (ogrs *KrboOnGroupRegisteredSubscription) OnEvent(
	handler keepRandomBeaconOperatorOnGroupRegisteredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorOnGroupRegistered)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GroupPubKey,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ogrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ogrs *KrboOnGroupRegisteredSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorOnGroupRegistered,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ogrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ogrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ogrs.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past OnGroupRegistered events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ogrs.contract.PastOnGroupRegisteredEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past OnGroupRegistered events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ogrs.contract.watchOnGroupRegistered(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchOnGroupRegistered(
	sink chan *abi.KeepRandomBeaconOperatorOnGroupRegistered,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchOnGroupRegistered(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event OnGroupRegistered had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event OnGroupRegistered failed "+
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

func (krbo *KeepRandomBeaconOperator) PastOnGroupRegisteredEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconOperatorOnGroupRegistered, error) {
	iterator, err := krbo.contract.FilterOnGroupRegistered(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past OnGroupRegistered events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorOnGroupRegistered, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) RelayEntryRequested(
	opts *ethlike.SubscribeOpts,
) *KrboRelayEntryRequestedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboRelayEntryRequestedSubscription{
		krbo,
		opts,
	}
}

type KrboRelayEntryRequestedSubscription struct {
	contract *KeepRandomBeaconOperator
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconOperatorRelayEntryRequestedFunc func(
	PreviousEntry []uint8,
	GroupPublicKey []uint8,
	blockNumber uint64,
)

func (rers *KrboRelayEntryRequestedSubscription) OnEvent(
	handler keepRandomBeaconOperatorRelayEntryRequestedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntryRequested)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.PreviousEntry,
					event.GroupPublicKey,
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

func (rers *KrboRelayEntryRequestedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntryRequested,
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
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rers.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past RelayEntryRequested events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rers.contract.PastRelayEntryRequestedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past RelayEntryRequested events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rers.contract.watchRelayEntryRequested(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchRelayEntryRequested(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntryRequested,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchRelayEntryRequested(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event RelayEntryRequested had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event RelayEntryRequested failed "+
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

func (krbo *KeepRandomBeaconOperator) PastRelayEntryRequestedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconOperatorRelayEntryRequested, error) {
	iterator, err := krbo.contract.FilterRelayEntryRequested(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RelayEntryRequested events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorRelayEntryRequested, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) RelayEntrySubmitted(
	opts *ethlike.SubscribeOpts,
) *KrboRelayEntrySubmittedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboRelayEntrySubmittedSubscription{
		krbo,
		opts,
	}
}

type KrboRelayEntrySubmittedSubscription struct {
	contract *KeepRandomBeaconOperator
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconOperatorRelayEntrySubmittedFunc func(
	blockNumber uint64,
)

func (ress *KrboRelayEntrySubmittedSubscription) OnEvent(
	handler keepRandomBeaconOperatorRelayEntrySubmittedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntrySubmitted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ress.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ress *KrboRelayEntrySubmittedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntrySubmitted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ress.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ress.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ress.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past RelayEntrySubmitted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ress.contract.PastRelayEntrySubmittedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past RelayEntrySubmitted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ress.contract.watchRelayEntrySubmitted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchRelayEntrySubmitted(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntrySubmitted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchRelayEntrySubmitted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event RelayEntrySubmitted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event RelayEntrySubmitted failed "+
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

func (krbo *KeepRandomBeaconOperator) PastRelayEntrySubmittedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconOperatorRelayEntrySubmitted, error) {
	iterator, err := krbo.contract.FilterRelayEntrySubmitted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RelayEntrySubmitted events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorRelayEntrySubmitted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) RelayEntryTimeoutReported(
	opts *ethlike.SubscribeOpts,
	groupIndexFilter []*big.Int,
) *KrboRelayEntryTimeoutReportedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboRelayEntryTimeoutReportedSubscription{
		krbo,
		opts,
		groupIndexFilter,
	}
}

type KrboRelayEntryTimeoutReportedSubscription struct {
	contract         *KeepRandomBeaconOperator
	opts             *ethlike.SubscribeOpts
	groupIndexFilter []*big.Int
}

type keepRandomBeaconOperatorRelayEntryTimeoutReportedFunc func(
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (retrs *KrboRelayEntryTimeoutReportedSubscription) OnEvent(
	handler keepRandomBeaconOperatorRelayEntryTimeoutReportedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GroupIndex,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := retrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (retrs *KrboRelayEntryTimeoutReportedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(retrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := retrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - retrs.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past RelayEntryTimeoutReported events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := retrs.contract.PastRelayEntryTimeoutReportedEvents(
					fromBlock,
					nil,
					retrs.groupIndexFilter,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past RelayEntryTimeoutReported events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := retrs.contract.watchRelayEntryTimeoutReported(
		sink,
		retrs.groupIndexFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchRelayEntryTimeoutReported(
	sink chan *abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported,
	groupIndexFilter []*big.Int,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchRelayEntryTimeoutReported(
			&bind.WatchOpts{Context: ctx},
			sink,
			groupIndexFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event RelayEntryTimeoutReported had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event RelayEntryTimeoutReported failed "+
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

func (krbo *KeepRandomBeaconOperator) PastRelayEntryTimeoutReportedEvents(
	startBlock uint64,
	endBlock *uint64,
	groupIndexFilter []*big.Int,
) ([]*abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported, error) {
	iterator, err := krbo.contract.FilterRelayEntryTimeoutReported(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		groupIndexFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RelayEntryTimeoutReported events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbo *KeepRandomBeaconOperator) UnauthorizedSigningReported(
	opts *ethlike.SubscribeOpts,
	groupIndexFilter []*big.Int,
) *KrboUnauthorizedSigningReportedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrboUnauthorizedSigningReportedSubscription{
		krbo,
		opts,
		groupIndexFilter,
	}
}

type KrboUnauthorizedSigningReportedSubscription struct {
	contract         *KeepRandomBeaconOperator
	opts             *ethlike.SubscribeOpts
	groupIndexFilter []*big.Int
}

type keepRandomBeaconOperatorUnauthorizedSigningReportedFunc func(
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (usrs *KrboUnauthorizedSigningReportedSubscription) OnEvent(
	handler keepRandomBeaconOperatorUnauthorizedSigningReportedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorUnauthorizedSigningReported)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.GroupIndex,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := usrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (usrs *KrboUnauthorizedSigningReportedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconOperatorUnauthorizedSigningReported,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(usrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := usrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - usrs.opts.PastBlocks

				krboLogger.Infof(
					"subscription monitoring fetching past UnauthorizedSigningReported events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := usrs.contract.PastUnauthorizedSigningReportedEvents(
					fromBlock,
					nil,
					usrs.groupIndexFilter,
				)
				if err != nil {
					krboLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krboLogger.Infof(
					"subscription monitoring fetched [%v] past UnauthorizedSigningReported events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := usrs.contract.watchUnauthorizedSigningReported(
		sink,
		usrs.groupIndexFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbo *KeepRandomBeaconOperator) watchUnauthorizedSigningReported(
	sink chan *abi.KeepRandomBeaconOperatorUnauthorizedSigningReported,
	groupIndexFilter []*big.Int,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbo.contract.WatchUnauthorizedSigningReported(
			&bind.WatchOpts{Context: ctx},
			sink,
			groupIndexFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krboLogger.Errorf(
			"subscription to event UnauthorizedSigningReported had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krboLogger.Errorf(
			"subscription to event UnauthorizedSigningReported failed "+
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

func (krbo *KeepRandomBeaconOperator) PastUnauthorizedSigningReportedEvents(
	startBlock uint64,
	endBlock *uint64,
	groupIndexFilter []*big.Int,
) ([]*abi.KeepRandomBeaconOperatorUnauthorizedSigningReported, error) {
	iterator, err := krbo.contract.FilterUnauthorizedSigningReported(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		groupIndexFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past UnauthorizedSigningReported events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconOperatorUnauthorizedSigningReported, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
