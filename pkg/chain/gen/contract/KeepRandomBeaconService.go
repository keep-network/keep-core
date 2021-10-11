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
var krbsLogger = log.Logger("keep-contract-KeepRandomBeaconService")

type KeepRandomBeaconService struct {
	contract          *abi.KeepRandomBeaconServiceImplV1
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

func NewKeepRandomBeaconService(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *ethlike.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*KeepRandomBeaconService, error) {
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

	contract, err := abi.NewKeepRandomBeaconServiceImplV1(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.KeepRandomBeaconServiceImplV1ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &KeepRandomBeaconService{
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
func (krbs *KeepRandomBeaconService) AddOperatorContract(
	operatorContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction addOperatorContract",
		" params: ",
		fmt.Sprint(
			operatorContract,
		),
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.AddOperatorContract(
		transactorOptions,
		operatorContract,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			nil,
			"addOperatorContract",
			operatorContract,
		)
	}

	krbsLogger.Infof(
		"submitted transaction addOperatorContract with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.AddOperatorContract(
				transactorOptions,
				operatorContract,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					nil,
					"addOperatorContract",
					operatorContract,
				)
			}

			krbsLogger.Infof(
				"submitted transaction addOperatorContract with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallAddOperatorContract(
	operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"addOperatorContract",
		&result,
		operatorContract,
	)

	return err
}

func (krbs *KeepRandomBeaconService) AddOperatorContractGasEstimate(
	operatorContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"addOperatorContract",
		krbs.contractABI,
		krbs.transactor,
		operatorContract,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) EntryCreated(
	requestId *big.Int,
	entry []uint8,
	submitter common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction entryCreated",
		" params: ",
		fmt.Sprint(
			requestId,
			entry,
			submitter,
		),
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.EntryCreated(
		transactorOptions,
		requestId,
		entry,
		submitter,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			nil,
			"entryCreated",
			requestId,
			entry,
			submitter,
		)
	}

	krbsLogger.Infof(
		"submitted transaction entryCreated with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.EntryCreated(
				transactorOptions,
				requestId,
				entry,
				submitter,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					nil,
					"entryCreated",
					requestId,
					entry,
					submitter,
				)
			}

			krbsLogger.Infof(
				"submitted transaction entryCreated with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallEntryCreated(
	requestId *big.Int,
	entry []uint8,
	submitter common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"entryCreated",
		&result,
		requestId,
		entry,
		submitter,
	)

	return err
}

func (krbs *KeepRandomBeaconService) EntryCreatedGasEstimate(
	requestId *big.Int,
	entry []uint8,
	submitter common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"entryCreated",
		krbs.contractABI,
		krbs.transactor,
		requestId,
		entry,
		submitter,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) ExecuteCallback(
	requestId *big.Int,
	entry *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction executeCallback",
		" params: ",
		fmt.Sprint(
			requestId,
			entry,
		),
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.ExecuteCallback(
		transactorOptions,
		requestId,
		entry,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			nil,
			"executeCallback",
			requestId,
			entry,
		)
	}

	krbsLogger.Infof(
		"submitted transaction executeCallback with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.ExecuteCallback(
				transactorOptions,
				requestId,
				entry,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					nil,
					"executeCallback",
					requestId,
					entry,
				)
			}

			krbsLogger.Infof(
				"submitted transaction executeCallback with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallExecuteCallback(
	requestId *big.Int,
	entry *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"executeCallback",
		&result,
		requestId,
		entry,
	)

	return err
}

func (krbs *KeepRandomBeaconService) ExecuteCallbackGasEstimate(
	requestId *big.Int,
	entry *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"executeCallback",
		krbs.contractABI,
		krbs.transactor,
		requestId,
		entry,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) FundDkgFeePool(
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction fundDkgFeePool",
		" value: ", value,
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.FundDkgFeePool(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			value,
			"fundDkgFeePool",
		)
	}

	krbsLogger.Infof(
		"submitted transaction fundDkgFeePool with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.FundDkgFeePool(
				transactorOptions,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					value,
					"fundDkgFeePool",
				)
			}

			krbsLogger.Infof(
				"submitted transaction fundDkgFeePool with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallFundDkgFeePool(
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, value,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"fundDkgFeePool",
		&result,
	)

	return err
}

func (krbs *KeepRandomBeaconService) FundDkgFeePoolGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"fundDkgFeePool",
		krbs.contractABI,
		krbs.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) FundRequestSubsidyFeePool(
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction fundRequestSubsidyFeePool",
		" value: ", value,
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.FundRequestSubsidyFeePool(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			value,
			"fundRequestSubsidyFeePool",
		)
	}

	krbsLogger.Infof(
		"submitted transaction fundRequestSubsidyFeePool with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.FundRequestSubsidyFeePool(
				transactorOptions,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					value,
					"fundRequestSubsidyFeePool",
				)
			}

			krbsLogger.Infof(
				"submitted transaction fundRequestSubsidyFeePool with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallFundRequestSubsidyFeePool(
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, value,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"fundRequestSubsidyFeePool",
		&result,
	)

	return err
}

func (krbs *KeepRandomBeaconService) FundRequestSubsidyFeePoolGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"fundRequestSubsidyFeePool",
		krbs.contractABI,
		krbs.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) Initialize(
	dkgContributionMargin *big.Int,
	registry common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction initialize",
		" params: ",
		fmt.Sprint(
			dkgContributionMargin,
			registry,
		),
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.Initialize(
		transactorOptions,
		dkgContributionMargin,
		registry,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			nil,
			"initialize",
			dkgContributionMargin,
			registry,
		)
	}

	krbsLogger.Infof(
		"submitted transaction initialize with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.Initialize(
				transactorOptions,
				dkgContributionMargin,
				registry,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					nil,
					"initialize",
					dkgContributionMargin,
					registry,
				)
			}

			krbsLogger.Infof(
				"submitted transaction initialize with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallInitialize(
	dkgContributionMargin *big.Int,
	registry common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"initialize",
		&result,
		dkgContributionMargin,
		registry,
	)

	return err
}

func (krbs *KeepRandomBeaconService) InitializeGasEstimate(
	dkgContributionMargin *big.Int,
	registry common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"initialize",
		krbs.contractABI,
		krbs.transactor,
		dkgContributionMargin,
		registry,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) RemoveOperatorContract(
	operatorContract common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction removeOperatorContract",
		" params: ",
		fmt.Sprint(
			operatorContract,
		),
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.RemoveOperatorContract(
		transactorOptions,
		operatorContract,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			nil,
			"removeOperatorContract",
			operatorContract,
		)
	}

	krbsLogger.Infof(
		"submitted transaction removeOperatorContract with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.RemoveOperatorContract(
				transactorOptions,
				operatorContract,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					nil,
					"removeOperatorContract",
					operatorContract,
				)
			}

			krbsLogger.Infof(
				"submitted transaction removeOperatorContract with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallRemoveOperatorContract(
	operatorContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"removeOperatorContract",
		&result,
		operatorContract,
	)

	return err
}

func (krbs *KeepRandomBeaconService) RemoveOperatorContractGasEstimate(
	operatorContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"removeOperatorContract",
		krbs.contractABI,
		krbs.transactor,
		operatorContract,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) RequestRelayEntry(
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction requestRelayEntry",
		" value: ", value,
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.RequestRelayEntry(
		transactorOptions,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			value,
			"requestRelayEntry",
		)
	}

	krbsLogger.Infof(
		"submitted transaction requestRelayEntry with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.RequestRelayEntry(
				transactorOptions,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					value,
					"requestRelayEntry",
				)
			}

			krbsLogger.Infof(
				"submitted transaction requestRelayEntry with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallRequestRelayEntry(
	value *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, value,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"requestRelayEntry",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) RequestRelayEntryGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"requestRelayEntry",
		krbs.contractABI,
		krbs.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbs *KeepRandomBeaconService) RequestRelayEntry0(
	callbackContract common.Address,
	callbackGas *big.Int,
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	krbsLogger.Debug(
		"submitting transaction requestRelayEntry0",
		" params: ",
		fmt.Sprint(
			callbackContract,
			callbackGas,
		),
		" value: ", value,
	)

	krbs.transactionMutex.Lock()
	defer krbs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *krbs.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := krbs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := krbs.contract.RequestRelayEntry0(
		transactorOptions,
		callbackContract,
		callbackGas,
	)
	if err != nil {
		return transaction, krbs.errorResolver.ResolveError(
			err,
			krbs.transactorOptions.From,
			value,
			"requestRelayEntry0",
			callbackContract,
			callbackGas,
		)
	}

	krbsLogger.Infof(
		"submitted transaction requestRelayEntry0 with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go krbs.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := krbs.contract.RequestRelayEntry0(
				transactorOptions,
				callbackContract,
				callbackGas,
			)
			if err != nil {
				return nil, krbs.errorResolver.ResolveError(
					err,
					krbs.transactorOptions.From,
					value,
					"requestRelayEntry0",
					callbackContract,
					callbackGas,
				)
			}

			krbsLogger.Infof(
				"submitted transaction requestRelayEntry0 with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	krbs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbs *KeepRandomBeaconService) CallRequestRelayEntry0(
	callbackContract common.Address,
	callbackGas *big.Int,
	value *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.transactorOptions.From,
		blockNumber, value,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"requestRelayEntry0",
		&result,
		callbackContract,
		callbackGas,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) RequestRelayEntry0GasEstimate(
	callbackContract common.Address,
	callbackGas *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		krbs.callerOptions.From,
		krbs.contractAddress,
		"requestRelayEntry0",
		krbs.contractABI,
		krbs.transactor,
		callbackContract,
		callbackGas,
	)

	return result, err
}

// ----- Const Methods ------

func (krbs *KeepRandomBeaconService) BaseCallbackGas() (*big.Int, error) {
	var result *big.Int
	result, err := krbs.contract.BaseCallbackGas(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"baseCallbackGas",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) BaseCallbackGasAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"baseCallbackGas",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) CallbackSurplusRecipient(
	requestId *big.Int,
) (common.Address, error) {
	var result common.Address
	result, err := krbs.contract.CallbackSurplusRecipient(
		krbs.callerOptions,
		requestId,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"callbackSurplusRecipient",
			requestId,
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) CallbackSurplusRecipientAtBlock(
	requestId *big.Int,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"callbackSurplusRecipient",
		&result,
		requestId,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) DkgContributionMargin() (*big.Int, error) {
	var result *big.Int
	result, err := krbs.contract.DkgContributionMargin(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"dkgContributionMargin",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) DkgContributionMarginAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"dkgContributionMargin",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) DkgFeePool() (*big.Int, error) {
	var result *big.Int
	result, err := krbs.contract.DkgFeePool(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"dkgFeePool",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) DkgFeePoolAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"dkgFeePool",
		&result,
	)

	return result, err
}

type entryFeeBreakdown struct {
	EntryVerificationFee *big.Int
	DkgContributionFee   *big.Int
	GroupProfitFee       *big.Int
	GasPriceCeiling      *big.Int
}

func (krbs *KeepRandomBeaconService) EntryFeeBreakdown() (entryFeeBreakdown, error) {
	var result entryFeeBreakdown
	result, err := krbs.contract.EntryFeeBreakdown(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"entryFeeBreakdown",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) EntryFeeBreakdownAtBlock(
	blockNumber *big.Int,
) (entryFeeBreakdown, error) {
	var result entryFeeBreakdown

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"entryFeeBreakdown",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) EntryFeeEstimate(
	callbackGas *big.Int,
) (*big.Int, error) {
	var result *big.Int
	result, err := krbs.contract.EntryFeeEstimate(
		krbs.callerOptions,
		callbackGas,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"entryFeeEstimate",
			callbackGas,
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) EntryFeeEstimateAtBlock(
	callbackGas *big.Int,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"entryFeeEstimate",
		&result,
		callbackGas,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) Initialized() (bool, error) {
	var result bool
	result, err := krbs.contract.Initialized(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"initialized",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) InitializedAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"initialized",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) RequestSubsidyFeePool() (*big.Int, error) {
	var result *big.Int
	result, err := krbs.contract.RequestSubsidyFeePool(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"requestSubsidyFeePool",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) RequestSubsidyFeePoolAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"requestSubsidyFeePool",
		&result,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) SelectOperatorContract(
	seed *big.Int,
) (common.Address, error) {
	var result common.Address
	result, err := krbs.contract.SelectOperatorContract(
		krbs.callerOptions,
		seed,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"selectOperatorContract",
			seed,
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) SelectOperatorContractAtBlock(
	seed *big.Int,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"selectOperatorContract",
		&result,
		seed,
	)

	return result, err
}

func (krbs *KeepRandomBeaconService) Version() (string, error) {
	var result string
	result, err := krbs.contract.Version(
		krbs.callerOptions,
	)

	if err != nil {
		return result, krbs.errorResolver.ResolveError(
			err,
			krbs.callerOptions.From,
			nil,
			"version",
		)
	}

	return result, err
}

func (krbs *KeepRandomBeaconService) VersionAtBlock(
	blockNumber *big.Int,
) (string, error) {
	var result string

	err := chainutil.CallAtBlock(
		krbs.callerOptions.From,
		blockNumber,
		nil,
		krbs.contractABI,
		krbs.caller,
		krbs.errorResolver,
		krbs.contractAddress,
		"version",
		&result,
	)

	return result, err
}

// ------ Events -------

func (krbs *KeepRandomBeaconService) RelayEntryGenerated(
	opts *ethlike.SubscribeOpts,
) *KrbsRelayEntryGeneratedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrbsRelayEntryGeneratedSubscription{
		krbs,
		opts,
	}
}

type KrbsRelayEntryGeneratedSubscription struct {
	contract *KeepRandomBeaconService
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconServiceRelayEntryGeneratedFunc func(
	RequestId *big.Int,
	Entry *big.Int,
	blockNumber uint64,
)

func (regs *KrbsRelayEntryGeneratedSubscription) OnEvent(
	handler keepRandomBeaconServiceRelayEntryGeneratedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconServiceImplV1RelayEntryGenerated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.RequestId,
					event.Entry,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := regs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (regs *KrbsRelayEntryGeneratedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconServiceImplV1RelayEntryGenerated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(regs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := regs.contract.blockCounter.CurrentBlock()
				if err != nil {
					krbsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - regs.opts.PastBlocks

				krbsLogger.Infof(
					"subscription monitoring fetching past RelayEntryGenerated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := regs.contract.PastRelayEntryGeneratedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krbsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krbsLogger.Infof(
					"subscription monitoring fetched [%v] past RelayEntryGenerated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := regs.contract.watchRelayEntryGenerated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (krbs *KeepRandomBeaconService) watchRelayEntryGenerated(
	sink chan *abi.KeepRandomBeaconServiceImplV1RelayEntryGenerated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbs.contract.WatchRelayEntryGenerated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krbsLogger.Errorf(
			"subscription to event RelayEntryGenerated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krbsLogger.Errorf(
			"subscription to event RelayEntryGenerated failed "+
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

func (krbs *KeepRandomBeaconService) PastRelayEntryGeneratedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconServiceImplV1RelayEntryGenerated, error) {
	iterator, err := krbs.contract.FilterRelayEntryGenerated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RelayEntryGenerated events: [%v]",
			err,
		)
	}

	events := make([]*abi.KeepRandomBeaconServiceImplV1RelayEntryGenerated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (krbs *KeepRandomBeaconService) RelayEntryRequested(
	opts *ethlike.SubscribeOpts,
) *KrbsRelayEntryRequestedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &KrbsRelayEntryRequestedSubscription{
		krbs,
		opts,
	}
}

type KrbsRelayEntryRequestedSubscription struct {
	contract *KeepRandomBeaconService
	opts     *ethlike.SubscribeOpts
}

type keepRandomBeaconServiceRelayEntryRequestedFunc func(
	RequestId *big.Int,
	blockNumber uint64,
)

func (rers *KrbsRelayEntryRequestedSubscription) OnEvent(
	handler keepRandomBeaconServiceRelayEntryRequestedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.KeepRandomBeaconServiceImplV1RelayEntryRequested)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.RequestId,
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

func (rers *KrbsRelayEntryRequestedSubscription) Pipe(
	sink chan *abi.KeepRandomBeaconServiceImplV1RelayEntryRequested,
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
					krbsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rers.opts.PastBlocks

				krbsLogger.Infof(
					"subscription monitoring fetching past RelayEntryRequested events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rers.contract.PastRelayEntryRequestedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					krbsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				krbsLogger.Infof(
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

func (krbs *KeepRandomBeaconService) watchRelayEntryRequested(
	sink chan *abi.KeepRandomBeaconServiceImplV1RelayEntryRequested,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return krbs.contract.WatchRelayEntryRequested(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		krbsLogger.Errorf(
			"subscription to event RelayEntryRequested had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		krbsLogger.Errorf(
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

func (krbs *KeepRandomBeaconService) PastRelayEntryRequestedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.KeepRandomBeaconServiceImplV1RelayEntryRequested, error) {
	iterator, err := krbs.contract.FilterRelayEntryRequested(
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

	events := make([]*abi.KeepRandomBeaconServiceImplV1RelayEntryRequested, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
