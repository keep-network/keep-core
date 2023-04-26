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

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/abi"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var wcLogger = log.Logger("keep-contract-WalletCoordinator")

type WalletCoordinator struct {
	contract          *abi.WalletCoordinator
	contractAddress   common.Address
	contractABI       *hostchainabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *chainutil.ErrorResolver
	nonceManager      *ethereum.NonceManager
	miningWaiter      *chainutil.MiningWaiter
	blockCounter      *ethereum.BlockCounter

	transactionMutex *sync.Mutex
}

func NewWalletCoordinator(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethereum.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethereum.BlockCounter,
	transactionMutex *sync.Mutex,
) (*WalletCoordinator, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	transactorOptions, err := bind.NewKeyedTransactorWithChainID(
		accountKey.PrivateKey,
		chainId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate transactor: [%v]", err)
	}

	contract, err := abi.NewWalletCoordinator(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.WalletCoordinatorABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &WalletCoordinator{
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
func (wc *WalletCoordinator) AddCoordinator(
	arg_coordinator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction addCoordinator",
		" params: ",
		fmt.Sprint(
			arg_coordinator,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.AddCoordinator(
		transactorOptions,
		arg_coordinator,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"addCoordinator",
			arg_coordinator,
		)
	}

	wcLogger.Infof(
		"submitted transaction addCoordinator with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.AddCoordinator(
				newTransactorOptions,
				arg_coordinator,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"addCoordinator",
					arg_coordinator,
				)
			}

			wcLogger.Infof(
				"submitted transaction addCoordinator with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallAddCoordinator(
	arg_coordinator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"addCoordinator",
		&result,
		arg_coordinator,
	)

	return err
}

func (wc *WalletCoordinator) AddCoordinatorGasEstimate(
	arg_coordinator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"addCoordinator",
		wc.contractABI,
		wc.transactor,
		arg_coordinator,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) Initialize(
	arg__bridge common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction initialize",
		" params: ",
		fmt.Sprint(
			arg__bridge,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.Initialize(
		transactorOptions,
		arg__bridge,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"initialize",
			arg__bridge,
		)
	}

	wcLogger.Infof(
		"submitted transaction initialize with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.Initialize(
				newTransactorOptions,
				arg__bridge,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"initialize",
					arg__bridge,
				)
			}

			wcLogger.Infof(
				"submitted transaction initialize with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallInitialize(
	arg__bridge common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"initialize",
		&result,
		arg__bridge,
	)

	return err
}

func (wc *WalletCoordinator) InitializeGasEstimate(
	arg__bridge common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"initialize",
		wc.contractABI,
		wc.transactor,
		arg__bridge,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) RemoveCoordinator(
	arg_coordinator common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction removeCoordinator",
		" params: ",
		fmt.Sprint(
			arg_coordinator,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.RemoveCoordinator(
		transactorOptions,
		arg_coordinator,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"removeCoordinator",
			arg_coordinator,
		)
	}

	wcLogger.Infof(
		"submitted transaction removeCoordinator with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.RemoveCoordinator(
				newTransactorOptions,
				arg_coordinator,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"removeCoordinator",
					arg_coordinator,
				)
			}

			wcLogger.Infof(
				"submitted transaction removeCoordinator with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallRemoveCoordinator(
	arg_coordinator common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"removeCoordinator",
		&result,
		arg_coordinator,
	)

	return err
}

func (wc *WalletCoordinator) RemoveCoordinatorGasEstimate(
	arg_coordinator common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"removeCoordinator",
		wc.contractABI,
		wc.transactor,
		arg_coordinator,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) RenounceOwnership(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction renounceOwnership",
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.RenounceOwnership(
		transactorOptions,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"renounceOwnership",
		)
	}

	wcLogger.Infof(
		"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.RenounceOwnership(
				newTransactorOptions,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"renounceOwnership",
				)
			}

			wcLogger.Infof(
				"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallRenounceOwnership(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"renounceOwnership",
		&result,
	)

	return err
}

func (wc *WalletCoordinator) RenounceOwnershipGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"renounceOwnership",
		wc.contractABI,
		wc.transactor,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) RequestHeartbeat(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction requestHeartbeat",
		" params: ",
		fmt.Sprint(
			arg_walletPubKeyHash,
			arg_message,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.RequestHeartbeat(
		transactorOptions,
		arg_walletPubKeyHash,
		arg_message,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"requestHeartbeat",
			arg_walletPubKeyHash,
			arg_message,
		)
	}

	wcLogger.Infof(
		"submitted transaction requestHeartbeat with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.RequestHeartbeat(
				newTransactorOptions,
				arg_walletPubKeyHash,
				arg_message,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"requestHeartbeat",
					arg_walletPubKeyHash,
					arg_message,
				)
			}

			wcLogger.Infof(
				"submitted transaction requestHeartbeat with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallRequestHeartbeat(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"requestHeartbeat",
		&result,
		arg_walletPubKeyHash,
		arg_message,
	)

	return err
}

func (wc *WalletCoordinator) RequestHeartbeatGasEstimate(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"requestHeartbeat",
		wc.contractABI,
		wc.transactor,
		arg_walletPubKeyHash,
		arg_message,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) RequestHeartbeatWithReimbursement(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction requestHeartbeatWithReimbursement",
		" params: ",
		fmt.Sprint(
			arg_walletPubKeyHash,
			arg_message,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.RequestHeartbeatWithReimbursement(
		transactorOptions,
		arg_walletPubKeyHash,
		arg_message,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"requestHeartbeatWithReimbursement",
			arg_walletPubKeyHash,
			arg_message,
		)
	}

	wcLogger.Infof(
		"submitted transaction requestHeartbeatWithReimbursement with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.RequestHeartbeatWithReimbursement(
				newTransactorOptions,
				arg_walletPubKeyHash,
				arg_message,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"requestHeartbeatWithReimbursement",
					arg_walletPubKeyHash,
					arg_message,
				)
			}

			wcLogger.Infof(
				"submitted transaction requestHeartbeatWithReimbursement with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallRequestHeartbeatWithReimbursement(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"requestHeartbeatWithReimbursement",
		&result,
		arg_walletPubKeyHash,
		arg_message,
	)

	return err
}

func (wc *WalletCoordinator) RequestHeartbeatWithReimbursementGasEstimate(
	arg_walletPubKeyHash [20]byte,
	arg_message []byte,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"requestHeartbeatWithReimbursement",
		wc.contractABI,
		wc.transactor,
		arg_walletPubKeyHash,
		arg_message,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) SubmitDepositSweepProposal(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction submitDepositSweepProposal",
		" params: ",
		fmt.Sprint(
			arg_proposal,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.SubmitDepositSweepProposal(
		transactorOptions,
		arg_proposal,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"submitDepositSweepProposal",
			arg_proposal,
		)
	}

	wcLogger.Infof(
		"submitted transaction submitDepositSweepProposal with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.SubmitDepositSweepProposal(
				newTransactorOptions,
				arg_proposal,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"submitDepositSweepProposal",
					arg_proposal,
				)
			}

			wcLogger.Infof(
				"submitted transaction submitDepositSweepProposal with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallSubmitDepositSweepProposal(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"submitDepositSweepProposal",
		&result,
		arg_proposal,
	)

	return err
}

func (wc *WalletCoordinator) SubmitDepositSweepProposalGasEstimate(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"submitDepositSweepProposal",
		wc.contractABI,
		wc.transactor,
		arg_proposal,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) SubmitDepositSweepProposalWithReimbursement(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction submitDepositSweepProposalWithReimbursement",
		" params: ",
		fmt.Sprint(
			arg_proposal,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.SubmitDepositSweepProposalWithReimbursement(
		transactorOptions,
		arg_proposal,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"submitDepositSweepProposalWithReimbursement",
			arg_proposal,
		)
	}

	wcLogger.Infof(
		"submitted transaction submitDepositSweepProposalWithReimbursement with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.SubmitDepositSweepProposalWithReimbursement(
				newTransactorOptions,
				arg_proposal,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"submitDepositSweepProposalWithReimbursement",
					arg_proposal,
				)
			}

			wcLogger.Infof(
				"submitted transaction submitDepositSweepProposalWithReimbursement with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallSubmitDepositSweepProposalWithReimbursement(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"submitDepositSweepProposalWithReimbursement",
		&result,
		arg_proposal,
	)

	return err
}

func (wc *WalletCoordinator) SubmitDepositSweepProposalWithReimbursementGasEstimate(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"submitDepositSweepProposalWithReimbursement",
		wc.contractABI,
		wc.transactor,
		arg_proposal,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) TransferOwnership(
	arg_newOwner common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction transferOwnership",
		" params: ",
		fmt.Sprint(
			arg_newOwner,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.TransferOwnership(
		transactorOptions,
		arg_newOwner,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"transferOwnership",
			arg_newOwner,
		)
	}

	wcLogger.Infof(
		"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.TransferOwnership(
				newTransactorOptions,
				arg_newOwner,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"transferOwnership",
					arg_newOwner,
				)
			}

			wcLogger.Infof(
				"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallTransferOwnership(
	arg_newOwner common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"transferOwnership",
		&result,
		arg_newOwner,
	)

	return err
}

func (wc *WalletCoordinator) TransferOwnershipGasEstimate(
	arg_newOwner common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"transferOwnership",
		wc.contractABI,
		wc.transactor,
		arg_newOwner,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UnlockWallet(
	arg_walletPubKeyHash [20]byte,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction unlockWallet",
		" params: ",
		fmt.Sprint(
			arg_walletPubKeyHash,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.UnlockWallet(
		transactorOptions,
		arg_walletPubKeyHash,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"unlockWallet",
			arg_walletPubKeyHash,
		)
	}

	wcLogger.Infof(
		"submitted transaction unlockWallet with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.UnlockWallet(
				newTransactorOptions,
				arg_walletPubKeyHash,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"unlockWallet",
					arg_walletPubKeyHash,
				)
			}

			wcLogger.Infof(
				"submitted transaction unlockWallet with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallUnlockWallet(
	arg_walletPubKeyHash [20]byte,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"unlockWallet",
		&result,
		arg_walletPubKeyHash,
	)

	return err
}

func (wc *WalletCoordinator) UnlockWalletGasEstimate(
	arg_walletPubKeyHash [20]byte,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"unlockWallet",
		wc.contractABI,
		wc.transactor,
		arg_walletPubKeyHash,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateDepositSweepProposalParameters(
	arg__depositSweepProposalValidity uint32,
	arg__depositMinAge uint32,
	arg__depositRefundSafetyMargin uint32,
	arg__depositSweepMaxSize uint16,
	arg__depositSweepProposalSubmissionGasOffset uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositSweepProposalParameters",
		" params: ",
		fmt.Sprint(
			arg__depositSweepProposalValidity,
			arg__depositMinAge,
			arg__depositRefundSafetyMargin,
			arg__depositSweepMaxSize,
			arg__depositSweepProposalSubmissionGasOffset,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.UpdateDepositSweepProposalParameters(
		transactorOptions,
		arg__depositSweepProposalValidity,
		arg__depositMinAge,
		arg__depositRefundSafetyMargin,
		arg__depositSweepMaxSize,
		arg__depositSweepProposalSubmissionGasOffset,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositSweepProposalParameters",
			arg__depositSweepProposalValidity,
			arg__depositMinAge,
			arg__depositRefundSafetyMargin,
			arg__depositSweepMaxSize,
			arg__depositSweepProposalSubmissionGasOffset,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositSweepProposalParameters with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.UpdateDepositSweepProposalParameters(
				newTransactorOptions,
				arg__depositSweepProposalValidity,
				arg__depositMinAge,
				arg__depositRefundSafetyMargin,
				arg__depositSweepMaxSize,
				arg__depositSweepProposalSubmissionGasOffset,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositSweepProposalParameters",
					arg__depositSweepProposalValidity,
					arg__depositMinAge,
					arg__depositRefundSafetyMargin,
					arg__depositSweepMaxSize,
					arg__depositSweepProposalSubmissionGasOffset,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositSweepProposalParameters with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallUpdateDepositSweepProposalParameters(
	arg__depositSweepProposalValidity uint32,
	arg__depositMinAge uint32,
	arg__depositRefundSafetyMargin uint32,
	arg__depositSweepMaxSize uint16,
	arg__depositSweepProposalSubmissionGasOffset uint32,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"updateDepositSweepProposalParameters",
		&result,
		arg__depositSweepProposalValidity,
		arg__depositMinAge,
		arg__depositRefundSafetyMargin,
		arg__depositSweepMaxSize,
		arg__depositSweepProposalSubmissionGasOffset,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositSweepProposalParametersGasEstimate(
	arg__depositSweepProposalValidity uint32,
	arg__depositMinAge uint32,
	arg__depositRefundSafetyMargin uint32,
	arg__depositSweepMaxSize uint16,
	arg__depositSweepProposalSubmissionGasOffset uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositSweepProposalParameters",
		wc.contractABI,
		wc.transactor,
		arg__depositSweepProposalValidity,
		arg__depositMinAge,
		arg__depositRefundSafetyMargin,
		arg__depositSweepMaxSize,
		arg__depositSweepProposalSubmissionGasOffset,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateHeartbeatRequestParameters(
	arg__heartbeatRequestValidity uint32,
	arg__heartbeatRequestGasOffset uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateHeartbeatRequestParameters",
		" params: ",
		fmt.Sprint(
			arg__heartbeatRequestValidity,
			arg__heartbeatRequestGasOffset,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.UpdateHeartbeatRequestParameters(
		transactorOptions,
		arg__heartbeatRequestValidity,
		arg__heartbeatRequestGasOffset,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateHeartbeatRequestParameters",
			arg__heartbeatRequestValidity,
			arg__heartbeatRequestGasOffset,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateHeartbeatRequestParameters with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.UpdateHeartbeatRequestParameters(
				newTransactorOptions,
				arg__heartbeatRequestValidity,
				arg__heartbeatRequestGasOffset,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateHeartbeatRequestParameters",
					arg__heartbeatRequestValidity,
					arg__heartbeatRequestGasOffset,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateHeartbeatRequestParameters with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallUpdateHeartbeatRequestParameters(
	arg__heartbeatRequestValidity uint32,
	arg__heartbeatRequestGasOffset uint32,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"updateHeartbeatRequestParameters",
		&result,
		arg__heartbeatRequestValidity,
		arg__heartbeatRequestGasOffset,
	)

	return err
}

func (wc *WalletCoordinator) UpdateHeartbeatRequestParametersGasEstimate(
	arg__heartbeatRequestValidity uint32,
	arg__heartbeatRequestGasOffset uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateHeartbeatRequestParameters",
		wc.contractABI,
		wc.transactor,
		arg__heartbeatRequestValidity,
		arg__heartbeatRequestGasOffset,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateReimbursementPool(
	arg__reimbursementPool common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateReimbursementPool",
		" params: ",
		fmt.Sprint(
			arg__reimbursementPool,
		),
	)

	wc.transactionMutex.Lock()
	defer wc.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *wc.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := wc.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := wc.contract.UpdateReimbursementPool(
		transactorOptions,
		arg__reimbursementPool,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateReimbursementPool",
			arg__reimbursementPool,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateReimbursementPool with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go wc.miningWaiter.ForceMining(
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

			transaction, err := wc.contract.UpdateReimbursementPool(
				newTransactorOptions,
				arg__reimbursementPool,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateReimbursementPool",
					arg__reimbursementPool,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateReimbursementPool with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	wc.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (wc *WalletCoordinator) CallUpdateReimbursementPool(
	arg__reimbursementPool common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		wc.transactorOptions.From,
		blockNumber, nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"updateReimbursementPool",
		&result,
		arg__reimbursementPool,
	)

	return err
}

func (wc *WalletCoordinator) UpdateReimbursementPoolGasEstimate(
	arg__reimbursementPool common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateReimbursementPool",
		wc.contractABI,
		wc.transactor,
		arg__reimbursementPool,
	)

	return result, err
}

// ----- Const Methods ------

func (wc *WalletCoordinator) Bridge() (common.Address, error) {
	result, err := wc.contract.Bridge(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"bridge",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) BridgeAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"bridge",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) DepositMinAge() (uint32, error) {
	result, err := wc.contract.DepositMinAge(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"depositMinAge",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) DepositMinAgeAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"depositMinAge",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) DepositRefundSafetyMargin() (uint32, error) {
	result, err := wc.contract.DepositRefundSafetyMargin(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"depositRefundSafetyMargin",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) DepositRefundSafetyMarginAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"depositRefundSafetyMargin",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) DepositSweepMaxSize() (uint16, error) {
	result, err := wc.contract.DepositSweepMaxSize(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"depositSweepMaxSize",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) DepositSweepMaxSizeAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"depositSweepMaxSize",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) DepositSweepProposalSubmissionGasOffset() (uint32, error) {
	result, err := wc.contract.DepositSweepProposalSubmissionGasOffset(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"depositSweepProposalSubmissionGasOffset",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) DepositSweepProposalSubmissionGasOffsetAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"depositSweepProposalSubmissionGasOffset",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) DepositSweepProposalValidity() (uint32, error) {
	result, err := wc.contract.DepositSweepProposalValidity(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"depositSweepProposalValidity",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) DepositSweepProposalValidityAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"depositSweepProposalValidity",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) HeartbeatRequestGasOffset() (uint32, error) {
	result, err := wc.contract.HeartbeatRequestGasOffset(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"heartbeatRequestGasOffset",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) HeartbeatRequestGasOffsetAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"heartbeatRequestGasOffset",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) HeartbeatRequestValidity() (uint32, error) {
	result, err := wc.contract.HeartbeatRequestValidity(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"heartbeatRequestValidity",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) HeartbeatRequestValidityAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"heartbeatRequestValidity",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) IsCoordinator(
	arg0 common.Address,
) (bool, error) {
	result, err := wc.contract.IsCoordinator(
		wc.callerOptions,
		arg0,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"isCoordinator",
			arg0,
		)
	}

	return result, err
}

func (wc *WalletCoordinator) IsCoordinatorAtBlock(
	arg0 common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"isCoordinator",
		&result,
		arg0,
	)

	return result, err
}

func (wc *WalletCoordinator) Owner() (common.Address, error) {
	result, err := wc.contract.Owner(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"owner",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) OwnerAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"owner",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) ReimbursementPool() (common.Address, error) {
	result, err := wc.contract.ReimbursementPool(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"reimbursementPool",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) ReimbursementPoolAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"reimbursementPool",
		&result,
	)

	return result, err
}

func (wc *WalletCoordinator) ValidateDepositSweepProposal(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_depositsExtraInfo []abi.WalletCoordinatorDepositExtraInfo,
) (bool, error) {
	result, err := wc.contract.ValidateDepositSweepProposal(
		wc.callerOptions,
		arg_proposal,
		arg_depositsExtraInfo,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"validateDepositSweepProposal",
			arg_proposal,
			arg_depositsExtraInfo,
		)
	}

	return result, err
}

func (wc *WalletCoordinator) ValidateDepositSweepProposalAtBlock(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_depositsExtraInfo []abi.WalletCoordinatorDepositExtraInfo,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"validateDepositSweepProposal",
		&result,
		arg_proposal,
		arg_depositsExtraInfo,
	)

	return result, err
}

type walletLock struct {
	ExpiresAt uint32
	Cause     uint8
}

func (wc *WalletCoordinator) WalletLock(
	arg0 [20]byte,
) (walletLock, error) {
	result, err := wc.contract.WalletLock(
		wc.callerOptions,
		arg0,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"walletLock",
			arg0,
		)
	}

	return result, err
}

func (wc *WalletCoordinator) WalletLockAtBlock(
	arg0 [20]byte,
	blockNumber *big.Int,
) (walletLock, error) {
	var result walletLock

	err := chainutil.CallAtBlock(
		wc.callerOptions.From,
		blockNumber,
		nil,
		wc.contractABI,
		wc.caller,
		wc.errorResolver,
		wc.contractAddress,
		"walletLock",
		&result,
		arg0,
	)

	return result, err
}

// ------ Events -------

func (wc *WalletCoordinator) CoordinatorAddedEvent(
	opts *ethereum.SubscribeOpts,
	coordinatorFilter []common.Address,
) *WcCoordinatorAddedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcCoordinatorAddedSubscription{
		wc,
		opts,
		coordinatorFilter,
	}
}

type WcCoordinatorAddedSubscription struct {
	contract          *WalletCoordinator
	opts              *ethereum.SubscribeOpts
	coordinatorFilter []common.Address
}

type walletCoordinatorCoordinatorAddedFunc func(
	Coordinator common.Address,
	blockNumber uint64,
)

func (cas *WcCoordinatorAddedSubscription) OnEvent(
	handler walletCoordinatorCoordinatorAddedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorCoordinatorAdded)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Coordinator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := cas.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (cas *WcCoordinatorAddedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorCoordinatorAdded,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(cas.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := cas.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - cas.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past CoordinatorAdded events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := cas.contract.PastCoordinatorAddedEvents(
					fromBlock,
					nil,
					cas.coordinatorFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past CoordinatorAdded events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := cas.contract.watchCoordinatorAdded(
		sink,
		cas.coordinatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchCoordinatorAdded(
	sink chan *abi.WalletCoordinatorCoordinatorAdded,
	coordinatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchCoordinatorAdded(
			&bind.WatchOpts{Context: ctx},
			sink,
			coordinatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event CoordinatorAdded had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event CoordinatorAdded failed "+
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

func (wc *WalletCoordinator) PastCoordinatorAddedEvents(
	startBlock uint64,
	endBlock *uint64,
	coordinatorFilter []common.Address,
) ([]*abi.WalletCoordinatorCoordinatorAdded, error) {
	iterator, err := wc.contract.FilterCoordinatorAdded(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		coordinatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past CoordinatorAdded events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorCoordinatorAdded, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) CoordinatorRemovedEvent(
	opts *ethereum.SubscribeOpts,
	coordinatorFilter []common.Address,
) *WcCoordinatorRemovedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcCoordinatorRemovedSubscription{
		wc,
		opts,
		coordinatorFilter,
	}
}

type WcCoordinatorRemovedSubscription struct {
	contract          *WalletCoordinator
	opts              *ethereum.SubscribeOpts
	coordinatorFilter []common.Address
}

type walletCoordinatorCoordinatorRemovedFunc func(
	Coordinator common.Address,
	blockNumber uint64,
)

func (crs *WcCoordinatorRemovedSubscription) OnEvent(
	handler walletCoordinatorCoordinatorRemovedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorCoordinatorRemoved)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Coordinator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := crs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (crs *WcCoordinatorRemovedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorCoordinatorRemoved,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(crs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := crs.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - crs.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past CoordinatorRemoved events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := crs.contract.PastCoordinatorRemovedEvents(
					fromBlock,
					nil,
					crs.coordinatorFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past CoordinatorRemoved events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := crs.contract.watchCoordinatorRemoved(
		sink,
		crs.coordinatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchCoordinatorRemoved(
	sink chan *abi.WalletCoordinatorCoordinatorRemoved,
	coordinatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchCoordinatorRemoved(
			&bind.WatchOpts{Context: ctx},
			sink,
			coordinatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event CoordinatorRemoved had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event CoordinatorRemoved failed "+
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

func (wc *WalletCoordinator) PastCoordinatorRemovedEvents(
	startBlock uint64,
	endBlock *uint64,
	coordinatorFilter []common.Address,
) ([]*abi.WalletCoordinatorCoordinatorRemoved, error) {
	iterator, err := wc.contract.FilterCoordinatorRemoved(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		coordinatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past CoordinatorRemoved events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorCoordinatorRemoved, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositSweepProposalParametersUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositSweepProposalParametersUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositSweepProposalParametersUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositSweepProposalParametersUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositSweepProposalParametersUpdatedFunc func(
	DepositSweepProposalValidity uint32,
	DepositMinAge uint32,
	DepositRefundSafetyMargin uint32,
	DepositSweepMaxSize uint16,
	DepositSweepProposalSubmissionGasOffset uint32,
	blockNumber uint64,
)

func (dsppus *WcDepositSweepProposalParametersUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositSweepProposalParametersUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositSweepProposalParametersUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositSweepProposalValidity,
					event.DepositMinAge,
					event.DepositRefundSafetyMargin,
					event.DepositSweepMaxSize,
					event.DepositSweepProposalSubmissionGasOffset,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dsppus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dsppus *WcDepositSweepProposalParametersUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositSweepProposalParametersUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dsppus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dsppus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dsppus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositSweepProposalParametersUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dsppus.contract.PastDepositSweepProposalParametersUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past DepositSweepProposalParametersUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dsppus.contract.watchDepositSweepProposalParametersUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepProposalParametersUpdated(
	sink chan *abi.WalletCoordinatorDepositSweepProposalParametersUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepProposalParametersUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositSweepProposalParametersUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositSweepProposalParametersUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositSweepProposalParametersUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositSweepProposalParametersUpdated, error) {
	iterator, err := wc.contract.FilterDepositSweepProposalParametersUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositSweepProposalParametersUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositSweepProposalParametersUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositSweepProposalSubmittedEvent(
	opts *ethereum.SubscribeOpts,
	coordinatorFilter []common.Address,
) *WcDepositSweepProposalSubmittedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositSweepProposalSubmittedSubscription{
		wc,
		opts,
		coordinatorFilter,
	}
}

type WcDepositSweepProposalSubmittedSubscription struct {
	contract          *WalletCoordinator
	opts              *ethereum.SubscribeOpts
	coordinatorFilter []common.Address
}

type walletCoordinatorDepositSweepProposalSubmittedFunc func(
	Proposal abi.WalletCoordinatorDepositSweepProposal,
	Coordinator common.Address,
	blockNumber uint64,
)

func (dspss *WcDepositSweepProposalSubmittedSubscription) OnEvent(
	handler walletCoordinatorDepositSweepProposalSubmittedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositSweepProposalSubmitted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Proposal,
					event.Coordinator,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dspss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dspss *WcDepositSweepProposalSubmittedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositSweepProposalSubmitted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dspss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dspss.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dspss.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositSweepProposalSubmitted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dspss.contract.PastDepositSweepProposalSubmittedEvents(
					fromBlock,
					nil,
					dspss.coordinatorFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past DepositSweepProposalSubmitted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dspss.contract.watchDepositSweepProposalSubmitted(
		sink,
		dspss.coordinatorFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepProposalSubmitted(
	sink chan *abi.WalletCoordinatorDepositSweepProposalSubmitted,
	coordinatorFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepProposalSubmitted(
			&bind.WatchOpts{Context: ctx},
			sink,
			coordinatorFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositSweepProposalSubmitted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositSweepProposalSubmitted failed "+
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

func (wc *WalletCoordinator) PastDepositSweepProposalSubmittedEvents(
	startBlock uint64,
	endBlock *uint64,
	coordinatorFilter []common.Address,
) ([]*abi.WalletCoordinatorDepositSweepProposalSubmitted, error) {
	iterator, err := wc.contract.FilterDepositSweepProposalSubmitted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		coordinatorFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositSweepProposalSubmitted events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositSweepProposalSubmitted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) HeartbeatRequestParametersUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcHeartbeatRequestParametersUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcHeartbeatRequestParametersUpdatedSubscription{
		wc,
		opts,
	}
}

type WcHeartbeatRequestParametersUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorHeartbeatRequestParametersUpdatedFunc func(
	HeartbeatRequestValidity uint32,
	HeartbeatRequestGasOffset uint32,
	blockNumber uint64,
)

func (hrpus *WcHeartbeatRequestParametersUpdatedSubscription) OnEvent(
	handler walletCoordinatorHeartbeatRequestParametersUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorHeartbeatRequestParametersUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.HeartbeatRequestValidity,
					event.HeartbeatRequestGasOffset,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := hrpus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (hrpus *WcHeartbeatRequestParametersUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorHeartbeatRequestParametersUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(hrpus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := hrpus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - hrpus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past HeartbeatRequestParametersUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := hrpus.contract.PastHeartbeatRequestParametersUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past HeartbeatRequestParametersUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := hrpus.contract.watchHeartbeatRequestParametersUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchHeartbeatRequestParametersUpdated(
	sink chan *abi.WalletCoordinatorHeartbeatRequestParametersUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchHeartbeatRequestParametersUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event HeartbeatRequestParametersUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event HeartbeatRequestParametersUpdated failed "+
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

func (wc *WalletCoordinator) PastHeartbeatRequestParametersUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorHeartbeatRequestParametersUpdated, error) {
	iterator, err := wc.contract.FilterHeartbeatRequestParametersUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past HeartbeatRequestParametersUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorHeartbeatRequestParametersUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) HeartbeatRequestSubmittedEvent(
	opts *ethereum.SubscribeOpts,
) *WcHeartbeatRequestSubmittedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcHeartbeatRequestSubmittedSubscription{
		wc,
		opts,
	}
}

type WcHeartbeatRequestSubmittedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorHeartbeatRequestSubmittedFunc func(
	WalletPubKeyHash [20]byte,
	Message []byte,
	blockNumber uint64,
)

func (hrss *WcHeartbeatRequestSubmittedSubscription) OnEvent(
	handler walletCoordinatorHeartbeatRequestSubmittedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorHeartbeatRequestSubmitted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.WalletPubKeyHash,
					event.Message,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := hrss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (hrss *WcHeartbeatRequestSubmittedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorHeartbeatRequestSubmitted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(hrss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := hrss.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - hrss.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past HeartbeatRequestSubmitted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := hrss.contract.PastHeartbeatRequestSubmittedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past HeartbeatRequestSubmitted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := hrss.contract.watchHeartbeatRequestSubmitted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchHeartbeatRequestSubmitted(
	sink chan *abi.WalletCoordinatorHeartbeatRequestSubmitted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchHeartbeatRequestSubmitted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event HeartbeatRequestSubmitted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event HeartbeatRequestSubmitted failed "+
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

func (wc *WalletCoordinator) PastHeartbeatRequestSubmittedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorHeartbeatRequestSubmitted, error) {
	iterator, err := wc.contract.FilterHeartbeatRequestSubmitted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past HeartbeatRequestSubmitted events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorHeartbeatRequestSubmitted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) InitializedEvent(
	opts *ethereum.SubscribeOpts,
) *WcInitializedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcInitializedSubscription{
		wc,
		opts,
	}
}

type WcInitializedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorInitializedFunc func(
	Version uint8,
	blockNumber uint64,
)

func (is *WcInitializedSubscription) OnEvent(
	handler walletCoordinatorInitializedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorInitialized)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.Version,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := is.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (is *WcInitializedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorInitialized,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(is.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := is.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - is.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past Initialized events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := is.contract.PastInitializedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past Initialized events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := is.contract.watchInitialized(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchInitialized(
	sink chan *abi.WalletCoordinatorInitialized,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchInitialized(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event Initialized had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event Initialized failed "+
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

func (wc *WalletCoordinator) PastInitializedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorInitialized, error) {
	iterator, err := wc.contract.FilterInitialized(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Initialized events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorInitialized, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) OwnershipTransferredEvent(
	opts *ethereum.SubscribeOpts,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) *WcOwnershipTransferredSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcOwnershipTransferredSubscription{
		wc,
		opts,
		previousOwnerFilter,
		newOwnerFilter,
	}
}

type WcOwnershipTransferredSubscription struct {
	contract            *WalletCoordinator
	opts                *ethereum.SubscribeOpts
	previousOwnerFilter []common.Address
	newOwnerFilter      []common.Address
}

type walletCoordinatorOwnershipTransferredFunc func(
	PreviousOwner common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (ots *WcOwnershipTransferredSubscription) OnEvent(
	handler walletCoordinatorOwnershipTransferredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorOwnershipTransferred)
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

func (ots *WcOwnershipTransferredSubscription) Pipe(
	sink chan *abi.WalletCoordinatorOwnershipTransferred,
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
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ots.opts.PastBlocks

				wcLogger.Infof(
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
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
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

func (wc *WalletCoordinator) watchOwnershipTransferred(
	sink chan *abi.WalletCoordinatorOwnershipTransferred,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchOwnershipTransferred(
			&bind.WatchOpts{Context: ctx},
			sink,
			previousOwnerFilter,
			newOwnerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event OwnershipTransferred had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
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

func (wc *WalletCoordinator) PastOwnershipTransferredEvents(
	startBlock uint64,
	endBlock *uint64,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) ([]*abi.WalletCoordinatorOwnershipTransferred, error) {
	iterator, err := wc.contract.FilterOwnershipTransferred(
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

	events := make([]*abi.WalletCoordinatorOwnershipTransferred, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) ReimbursementPoolUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcReimbursementPoolUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcReimbursementPoolUpdatedSubscription{
		wc,
		opts,
	}
}

type WcReimbursementPoolUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorReimbursementPoolUpdatedFunc func(
	NewReimbursementPool common.Address,
	blockNumber uint64,
)

func (rpus *WcReimbursementPoolUpdatedSubscription) OnEvent(
	handler walletCoordinatorReimbursementPoolUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorReimbursementPoolUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.NewReimbursementPool,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rpus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rpus *WcReimbursementPoolUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorReimbursementPoolUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rpus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rpus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rpus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past ReimbursementPoolUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rpus.contract.PastReimbursementPoolUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past ReimbursementPoolUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rpus.contract.watchReimbursementPoolUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchReimbursementPoolUpdated(
	sink chan *abi.WalletCoordinatorReimbursementPoolUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchReimbursementPoolUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event ReimbursementPoolUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event ReimbursementPoolUpdated failed "+
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

func (wc *WalletCoordinator) PastReimbursementPoolUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorReimbursementPoolUpdated, error) {
	iterator, err := wc.contract.FilterReimbursementPoolUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past ReimbursementPoolUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorReimbursementPoolUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) WalletManuallyUnlockedEvent(
	opts *ethereum.SubscribeOpts,
	walletPubKeyHashFilter [][20]byte,
) *WcWalletManuallyUnlockedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcWalletManuallyUnlockedSubscription{
		wc,
		opts,
		walletPubKeyHashFilter,
	}
}

type WcWalletManuallyUnlockedSubscription struct {
	contract               *WalletCoordinator
	opts                   *ethereum.SubscribeOpts
	walletPubKeyHashFilter [][20]byte
}

type walletCoordinatorWalletManuallyUnlockedFunc func(
	WalletPubKeyHash [20]byte,
	blockNumber uint64,
)

func (wmus *WcWalletManuallyUnlockedSubscription) OnEvent(
	handler walletCoordinatorWalletManuallyUnlockedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorWalletManuallyUnlocked)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.WalletPubKeyHash,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := wmus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wmus *WcWalletManuallyUnlockedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorWalletManuallyUnlocked,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(wmus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := wmus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - wmus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past WalletManuallyUnlocked events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := wmus.contract.PastWalletManuallyUnlockedEvents(
					fromBlock,
					nil,
					wmus.walletPubKeyHashFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past WalletManuallyUnlocked events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := wmus.contract.watchWalletManuallyUnlocked(
		sink,
		wmus.walletPubKeyHashFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchWalletManuallyUnlocked(
	sink chan *abi.WalletCoordinatorWalletManuallyUnlocked,
	walletPubKeyHashFilter [][20]byte,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchWalletManuallyUnlocked(
			&bind.WatchOpts{Context: ctx},
			sink,
			walletPubKeyHashFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event WalletManuallyUnlocked had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event WalletManuallyUnlocked failed "+
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

func (wc *WalletCoordinator) PastWalletManuallyUnlockedEvents(
	startBlock uint64,
	endBlock *uint64,
	walletPubKeyHashFilter [][20]byte,
) ([]*abi.WalletCoordinatorWalletManuallyUnlocked, error) {
	iterator, err := wc.contract.FilterWalletManuallyUnlocked(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		walletPubKeyHashFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past WalletManuallyUnlocked events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorWalletManuallyUnlocked, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
