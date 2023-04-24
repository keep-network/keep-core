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
func (wc *WalletCoordinator) AddProposalSubmitter(
	arg_proposalSubmitter common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction addProposalSubmitter",
		" params: ",
		fmt.Sprint(
			arg_proposalSubmitter,
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

	transaction, err := wc.contract.AddProposalSubmitter(
		transactorOptions,
		arg_proposalSubmitter,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"addProposalSubmitter",
			arg_proposalSubmitter,
		)
	}

	wcLogger.Infof(
		"submitted transaction addProposalSubmitter with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.AddProposalSubmitter(
				newTransactorOptions,
				arg_proposalSubmitter,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"addProposalSubmitter",
					arg_proposalSubmitter,
				)
			}

			wcLogger.Infof(
				"submitted transaction addProposalSubmitter with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallAddProposalSubmitter(
	arg_proposalSubmitter common.Address,
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
		"addProposalSubmitter",
		&result,
		arg_proposalSubmitter,
	)

	return err
}

func (wc *WalletCoordinator) AddProposalSubmitterGasEstimate(
	arg_proposalSubmitter common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"addProposalSubmitter",
		wc.contractABI,
		wc.transactor,
		arg_proposalSubmitter,
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
func (wc *WalletCoordinator) RemoveProposalSubmitter(
	arg_proposalSubmitter common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction removeProposalSubmitter",
		" params: ",
		fmt.Sprint(
			arg_proposalSubmitter,
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

	transaction, err := wc.contract.RemoveProposalSubmitter(
		transactorOptions,
		arg_proposalSubmitter,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"removeProposalSubmitter",
			arg_proposalSubmitter,
		)
	}

	wcLogger.Infof(
		"submitted transaction removeProposalSubmitter with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.RemoveProposalSubmitter(
				newTransactorOptions,
				arg_proposalSubmitter,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"removeProposalSubmitter",
					arg_proposalSubmitter,
				)
			}

			wcLogger.Infof(
				"submitted transaction removeProposalSubmitter with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallRemoveProposalSubmitter(
	arg_proposalSubmitter common.Address,
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
		"removeProposalSubmitter",
		&result,
		arg_proposalSubmitter,
	)

	return err
}

func (wc *WalletCoordinator) RemoveProposalSubmitterGasEstimate(
	arg_proposalSubmitter common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"removeProposalSubmitter",
		wc.contractABI,
		wc.transactor,
		arg_proposalSubmitter,
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
func (wc *WalletCoordinator) SubmitDepositSweepProposal(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction submitDepositSweepProposal",
		" params: ",
		fmt.Sprint(
			arg_proposal,
			arg_walletMemberContext,
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
		arg_walletMemberContext,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"submitDepositSweepProposal",
			arg_proposal,
			arg_walletMemberContext,
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
				arg_walletMemberContext,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"submitDepositSweepProposal",
					arg_proposal,
					arg_walletMemberContext,
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
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,
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
		arg_walletMemberContext,
	)

	return err
}

func (wc *WalletCoordinator) SubmitDepositSweepProposalGasEstimate(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"submitDepositSweepProposal",
		wc.contractABI,
		wc.transactor,
		arg_proposal,
		arg_walletMemberContext,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) SubmitDepositSweepProposalWithReimbursement(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction submitDepositSweepProposalWithReimbursement",
		" params: ",
		fmt.Sprint(
			arg_proposal,
			arg_walletMemberContext,
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
		arg_walletMemberContext,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"submitDepositSweepProposalWithReimbursement",
			arg_proposal,
			arg_walletMemberContext,
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
				arg_walletMemberContext,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"submitDepositSweepProposalWithReimbursement",
					arg_proposal,
					arg_walletMemberContext,
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
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,
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
		arg_walletMemberContext,
	)

	return err
}

func (wc *WalletCoordinator) SubmitDepositSweepProposalWithReimbursementGasEstimate(
	arg_proposal abi.WalletCoordinatorDepositSweepProposal,
	arg_walletMemberContext abi.WalletCoordinatorWalletMemberContext,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"submitDepositSweepProposalWithReimbursement",
		wc.contractABI,
		wc.transactor,
		arg_proposal,
		arg_walletMemberContext,
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
func (wc *WalletCoordinator) UpdateDepositMinAge(
	arg__depositMinAge uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositMinAge",
		" params: ",
		fmt.Sprint(
			arg__depositMinAge,
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

	transaction, err := wc.contract.UpdateDepositMinAge(
		transactorOptions,
		arg__depositMinAge,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositMinAge",
			arg__depositMinAge,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositMinAge with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.UpdateDepositMinAge(
				newTransactorOptions,
				arg__depositMinAge,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositMinAge",
					arg__depositMinAge,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositMinAge with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallUpdateDepositMinAge(
	arg__depositMinAge uint32,
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
		"updateDepositMinAge",
		&result,
		arg__depositMinAge,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositMinAgeGasEstimate(
	arg__depositMinAge uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositMinAge",
		wc.contractABI,
		wc.transactor,
		arg__depositMinAge,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateDepositRefundSafetyMargin(
	arg__depositRefundSafetyMargin uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositRefundSafetyMargin",
		" params: ",
		fmt.Sprint(
			arg__depositRefundSafetyMargin,
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

	transaction, err := wc.contract.UpdateDepositRefundSafetyMargin(
		transactorOptions,
		arg__depositRefundSafetyMargin,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositRefundSafetyMargin",
			arg__depositRefundSafetyMargin,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositRefundSafetyMargin with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.UpdateDepositRefundSafetyMargin(
				newTransactorOptions,
				arg__depositRefundSafetyMargin,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositRefundSafetyMargin",
					arg__depositRefundSafetyMargin,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositRefundSafetyMargin with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallUpdateDepositRefundSafetyMargin(
	arg__depositRefundSafetyMargin uint32,
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
		"updateDepositRefundSafetyMargin",
		&result,
		arg__depositRefundSafetyMargin,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositRefundSafetyMarginGasEstimate(
	arg__depositRefundSafetyMargin uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositRefundSafetyMargin",
		wc.contractABI,
		wc.transactor,
		arg__depositRefundSafetyMargin,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateDepositSweepMaxSize(
	arg__depositSweepMaxSize uint16,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositSweepMaxSize",
		" params: ",
		fmt.Sprint(
			arg__depositSweepMaxSize,
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

	transaction, err := wc.contract.UpdateDepositSweepMaxSize(
		transactorOptions,
		arg__depositSweepMaxSize,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositSweepMaxSize",
			arg__depositSweepMaxSize,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositSweepMaxSize with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.UpdateDepositSweepMaxSize(
				newTransactorOptions,
				arg__depositSweepMaxSize,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositSweepMaxSize",
					arg__depositSweepMaxSize,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositSweepMaxSize with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallUpdateDepositSweepMaxSize(
	arg__depositSweepMaxSize uint16,
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
		"updateDepositSweepMaxSize",
		&result,
		arg__depositSweepMaxSize,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositSweepMaxSizeGasEstimate(
	arg__depositSweepMaxSize uint16,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositSweepMaxSize",
		wc.contractABI,
		wc.transactor,
		arg__depositSweepMaxSize,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateDepositSweepProposalSubmissionGasOffset(
	arg__depositSweepProposalSubmissionGasOffset uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositSweepProposalSubmissionGasOffset",
		" params: ",
		fmt.Sprint(
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

	transaction, err := wc.contract.UpdateDepositSweepProposalSubmissionGasOffset(
		transactorOptions,
		arg__depositSweepProposalSubmissionGasOffset,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositSweepProposalSubmissionGasOffset",
			arg__depositSweepProposalSubmissionGasOffset,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositSweepProposalSubmissionGasOffset with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.UpdateDepositSweepProposalSubmissionGasOffset(
				newTransactorOptions,
				arg__depositSweepProposalSubmissionGasOffset,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositSweepProposalSubmissionGasOffset",
					arg__depositSweepProposalSubmissionGasOffset,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositSweepProposalSubmissionGasOffset with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallUpdateDepositSweepProposalSubmissionGasOffset(
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
		"updateDepositSweepProposalSubmissionGasOffset",
		&result,
		arg__depositSweepProposalSubmissionGasOffset,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositSweepProposalSubmissionGasOffsetGasEstimate(
	arg__depositSweepProposalSubmissionGasOffset uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositSweepProposalSubmissionGasOffset",
		wc.contractABI,
		wc.transactor,
		arg__depositSweepProposalSubmissionGasOffset,
	)

	return result, err
}

// Transaction submission.
func (wc *WalletCoordinator) UpdateDepositSweepProposalValidity(
	arg__depositSweepProposalValidity uint32,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	wcLogger.Debug(
		"submitting transaction updateDepositSweepProposalValidity",
		" params: ",
		fmt.Sprint(
			arg__depositSweepProposalValidity,
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

	transaction, err := wc.contract.UpdateDepositSweepProposalValidity(
		transactorOptions,
		arg__depositSweepProposalValidity,
	)
	if err != nil {
		return transaction, wc.errorResolver.ResolveError(
			err,
			wc.transactorOptions.From,
			nil,
			"updateDepositSweepProposalValidity",
			arg__depositSweepProposalValidity,
		)
	}

	wcLogger.Infof(
		"submitted transaction updateDepositSweepProposalValidity with id: [%s] and nonce [%v]",
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

			transaction, err := wc.contract.UpdateDepositSweepProposalValidity(
				newTransactorOptions,
				arg__depositSweepProposalValidity,
			)
			if err != nil {
				return nil, wc.errorResolver.ResolveError(
					err,
					wc.transactorOptions.From,
					nil,
					"updateDepositSweepProposalValidity",
					arg__depositSweepProposalValidity,
				)
			}

			wcLogger.Infof(
				"submitted transaction updateDepositSweepProposalValidity with id: [%s] and nonce [%v]",
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
func (wc *WalletCoordinator) CallUpdateDepositSweepProposalValidity(
	arg__depositSweepProposalValidity uint32,
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
		"updateDepositSweepProposalValidity",
		&result,
		arg__depositSweepProposalValidity,
	)

	return err
}

func (wc *WalletCoordinator) UpdateDepositSweepProposalValidityGasEstimate(
	arg__depositSweepProposalValidity uint32,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		wc.callerOptions.From,
		wc.contractAddress,
		"updateDepositSweepProposalValidity",
		wc.contractABI,
		wc.transactor,
		arg__depositSweepProposalValidity,
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

func (wc *WalletCoordinator) IsProposalSubmitter(
	arg0 common.Address,
) (bool, error) {
	result, err := wc.contract.IsProposalSubmitter(
		wc.callerOptions,
		arg0,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"isProposalSubmitter",
			arg0,
		)
	}

	return result, err
}

func (wc *WalletCoordinator) IsProposalSubmitterAtBlock(
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
		"isProposalSubmitter",
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

func (wc *WalletCoordinator) WalletRegistry() (common.Address, error) {
	result, err := wc.contract.WalletRegistry(
		wc.callerOptions,
	)

	if err != nil {
		return result, wc.errorResolver.ResolveError(
			err,
			wc.callerOptions.From,
			nil,
			"walletRegistry",
		)
	}

	return result, err
}

func (wc *WalletCoordinator) WalletRegistryAtBlock(
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
		"walletRegistry",
		&result,
	)

	return result, err
}

// ------ Events -------

func (wc *WalletCoordinator) DepositMinAgeUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositMinAgeUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositMinAgeUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositMinAgeUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositMinAgeUpdatedFunc func(
	DepositMinAge uint32,
	blockNumber uint64,
)

func (dmaus *WcDepositMinAgeUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositMinAgeUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositMinAgeUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositMinAge,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dmaus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dmaus *WcDepositMinAgeUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositMinAgeUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dmaus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dmaus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dmaus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositMinAgeUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dmaus.contract.PastDepositMinAgeUpdatedEvents(
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
					"subscription monitoring fetched [%v] past DepositMinAgeUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dmaus.contract.watchDepositMinAgeUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositMinAgeUpdated(
	sink chan *abi.WalletCoordinatorDepositMinAgeUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositMinAgeUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositMinAgeUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositMinAgeUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositMinAgeUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositMinAgeUpdated, error) {
	iterator, err := wc.contract.FilterDepositMinAgeUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositMinAgeUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositMinAgeUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositRefundSafetyMarginUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositRefundSafetyMarginUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositRefundSafetyMarginUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositRefundSafetyMarginUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositRefundSafetyMarginUpdatedFunc func(
	DepositRefundSafetyMargin uint32,
	blockNumber uint64,
)

func (drsmus *WcDepositRefundSafetyMarginUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositRefundSafetyMarginUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositRefundSafetyMarginUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositRefundSafetyMargin,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := drsmus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (drsmus *WcDepositRefundSafetyMarginUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositRefundSafetyMarginUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(drsmus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := drsmus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - drsmus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositRefundSafetyMarginUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := drsmus.contract.PastDepositRefundSafetyMarginUpdatedEvents(
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
					"subscription monitoring fetched [%v] past DepositRefundSafetyMarginUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := drsmus.contract.watchDepositRefundSafetyMarginUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositRefundSafetyMarginUpdated(
	sink chan *abi.WalletCoordinatorDepositRefundSafetyMarginUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositRefundSafetyMarginUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositRefundSafetyMarginUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositRefundSafetyMarginUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositRefundSafetyMarginUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositRefundSafetyMarginUpdated, error) {
	iterator, err := wc.contract.FilterDepositRefundSafetyMarginUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositRefundSafetyMarginUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositRefundSafetyMarginUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositSweepMaxSizeUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositSweepMaxSizeUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositSweepMaxSizeUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositSweepMaxSizeUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositSweepMaxSizeUpdatedFunc func(
	DepositSweepMaxSize uint16,
	blockNumber uint64,
)

func (dsmsus *WcDepositSweepMaxSizeUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositSweepMaxSizeUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositSweepMaxSizeUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositSweepMaxSize,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dsmsus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dsmsus *WcDepositSweepMaxSizeUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositSweepMaxSizeUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dsmsus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dsmsus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dsmsus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositSweepMaxSizeUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dsmsus.contract.PastDepositSweepMaxSizeUpdatedEvents(
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
					"subscription monitoring fetched [%v] past DepositSweepMaxSizeUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dsmsus.contract.watchDepositSweepMaxSizeUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepMaxSizeUpdated(
	sink chan *abi.WalletCoordinatorDepositSweepMaxSizeUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepMaxSizeUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositSweepMaxSizeUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositSweepMaxSizeUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositSweepMaxSizeUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositSweepMaxSizeUpdated, error) {
	iterator, err := wc.contract.FilterDepositSweepMaxSizeUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositSweepMaxSizeUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositSweepMaxSizeUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositSweepProposalSubmissionGasOffsetUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositSweepProposalSubmissionGasOffsetUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositSweepProposalSubmissionGasOffsetUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositSweepProposalSubmissionGasOffsetUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedFunc func(
	DepositSweepProposalSubmissionGasOffset uint32,
	blockNumber uint64,
)

func (dspsgous *WcDepositSweepProposalSubmissionGasOffsetUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositSweepProposalSubmissionGasOffset,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dspsgous.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dspsgous *WcDepositSweepProposalSubmissionGasOffsetUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dspsgous.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dspsgous.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dspsgous.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositSweepProposalSubmissionGasOffsetUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dspsgous.contract.PastDepositSweepProposalSubmissionGasOffsetUpdatedEvents(
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
					"subscription monitoring fetched [%v] past DepositSweepProposalSubmissionGasOffsetUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dspsgous.contract.watchDepositSweepProposalSubmissionGasOffsetUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepProposalSubmissionGasOffsetUpdated(
	sink chan *abi.WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepProposalSubmissionGasOffsetUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositSweepProposalSubmissionGasOffsetUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositSweepProposalSubmissionGasOffsetUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositSweepProposalSubmissionGasOffsetUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated, error) {
	iterator, err := wc.contract.FilterDepositSweepProposalSubmissionGasOffsetUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositSweepProposalSubmissionGasOffsetUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) DepositSweepProposalSubmittedEvent(
	opts *ethereum.SubscribeOpts,
	proposalSubmitterFilter []common.Address,
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
		proposalSubmitterFilter,
	}
}

type WcDepositSweepProposalSubmittedSubscription struct {
	contract                *WalletCoordinator
	opts                    *ethereum.SubscribeOpts
	proposalSubmitterFilter []common.Address
}

type walletCoordinatorDepositSweepProposalSubmittedFunc func(
	Proposal abi.WalletCoordinatorDepositSweepProposal,
	ProposalSubmitter common.Address,
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
					event.ProposalSubmitter,
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
					dspss.proposalSubmitterFilter,
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
		dspss.proposalSubmitterFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepProposalSubmitted(
	sink chan *abi.WalletCoordinatorDepositSweepProposalSubmitted,
	proposalSubmitterFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepProposalSubmitted(
			&bind.WatchOpts{Context: ctx},
			sink,
			proposalSubmitterFilter,
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
	proposalSubmitterFilter []common.Address,
) ([]*abi.WalletCoordinatorDepositSweepProposalSubmitted, error) {
	iterator, err := wc.contract.FilterDepositSweepProposalSubmitted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		proposalSubmitterFilter,
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

func (wc *WalletCoordinator) DepositSweepProposalValidityUpdatedEvent(
	opts *ethereum.SubscribeOpts,
) *WcDepositSweepProposalValidityUpdatedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcDepositSweepProposalValidityUpdatedSubscription{
		wc,
		opts,
	}
}

type WcDepositSweepProposalValidityUpdatedSubscription struct {
	contract *WalletCoordinator
	opts     *ethereum.SubscribeOpts
}

type walletCoordinatorDepositSweepProposalValidityUpdatedFunc func(
	DepositSweepProposalValidity uint32,
	blockNumber uint64,
)

func (dspvus *WcDepositSweepProposalValidityUpdatedSubscription) OnEvent(
	handler walletCoordinatorDepositSweepProposalValidityUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorDepositSweepProposalValidityUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositSweepProposalValidity,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := dspvus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (dspvus *WcDepositSweepProposalValidityUpdatedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorDepositSweepProposalValidityUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(dspvus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := dspvus.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - dspvus.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past DepositSweepProposalValidityUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := dspvus.contract.PastDepositSweepProposalValidityUpdatedEvents(
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
					"subscription monitoring fetched [%v] past DepositSweepProposalValidityUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := dspvus.contract.watchDepositSweepProposalValidityUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchDepositSweepProposalValidityUpdated(
	sink chan *abi.WalletCoordinatorDepositSweepProposalValidityUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchDepositSweepProposalValidityUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event DepositSweepProposalValidityUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event DepositSweepProposalValidityUpdated failed "+
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

func (wc *WalletCoordinator) PastDepositSweepProposalValidityUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.WalletCoordinatorDepositSweepProposalValidityUpdated, error) {
	iterator, err := wc.contract.FilterDepositSweepProposalValidityUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past DepositSweepProposalValidityUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorDepositSweepProposalValidityUpdated, 0)

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

func (wc *WalletCoordinator) ProposalSubmitterAddedEvent(
	opts *ethereum.SubscribeOpts,
	proposalSubmitterFilter []common.Address,
) *WcProposalSubmitterAddedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcProposalSubmitterAddedSubscription{
		wc,
		opts,
		proposalSubmitterFilter,
	}
}

type WcProposalSubmitterAddedSubscription struct {
	contract                *WalletCoordinator
	opts                    *ethereum.SubscribeOpts
	proposalSubmitterFilter []common.Address
}

type walletCoordinatorProposalSubmitterAddedFunc func(
	ProposalSubmitter common.Address,
	blockNumber uint64,
)

func (psas *WcProposalSubmitterAddedSubscription) OnEvent(
	handler walletCoordinatorProposalSubmitterAddedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorProposalSubmitterAdded)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.ProposalSubmitter,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := psas.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (psas *WcProposalSubmitterAddedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorProposalSubmitterAdded,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(psas.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := psas.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - psas.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past ProposalSubmitterAdded events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := psas.contract.PastProposalSubmitterAddedEvents(
					fromBlock,
					nil,
					psas.proposalSubmitterFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past ProposalSubmitterAdded events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := psas.contract.watchProposalSubmitterAdded(
		sink,
		psas.proposalSubmitterFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchProposalSubmitterAdded(
	sink chan *abi.WalletCoordinatorProposalSubmitterAdded,
	proposalSubmitterFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchProposalSubmitterAdded(
			&bind.WatchOpts{Context: ctx},
			sink,
			proposalSubmitterFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event ProposalSubmitterAdded had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event ProposalSubmitterAdded failed "+
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

func (wc *WalletCoordinator) PastProposalSubmitterAddedEvents(
	startBlock uint64,
	endBlock *uint64,
	proposalSubmitterFilter []common.Address,
) ([]*abi.WalletCoordinatorProposalSubmitterAdded, error) {
	iterator, err := wc.contract.FilterProposalSubmitterAdded(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		proposalSubmitterFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past ProposalSubmitterAdded events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorProposalSubmitterAdded, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (wc *WalletCoordinator) ProposalSubmitterRemovedEvent(
	opts *ethereum.SubscribeOpts,
	proposalSubmitterFilter []common.Address,
) *WcProposalSubmitterRemovedSubscription {
	if opts == nil {
		opts = new(ethereum.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &WcProposalSubmitterRemovedSubscription{
		wc,
		opts,
		proposalSubmitterFilter,
	}
}

type WcProposalSubmitterRemovedSubscription struct {
	contract                *WalletCoordinator
	opts                    *ethereum.SubscribeOpts
	proposalSubmitterFilter []common.Address
}

type walletCoordinatorProposalSubmitterRemovedFunc func(
	ProposalSubmitter common.Address,
	blockNumber uint64,
)

func (psrs *WcProposalSubmitterRemovedSubscription) OnEvent(
	handler walletCoordinatorProposalSubmitterRemovedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.WalletCoordinatorProposalSubmitterRemoved)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.ProposalSubmitter,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := psrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (psrs *WcProposalSubmitterRemovedSubscription) Pipe(
	sink chan *abi.WalletCoordinatorProposalSubmitterRemoved,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(psrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := psrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - psrs.opts.PastBlocks

				wcLogger.Infof(
					"subscription monitoring fetching past ProposalSubmitterRemoved events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := psrs.contract.PastProposalSubmitterRemovedEvents(
					fromBlock,
					nil,
					psrs.proposalSubmitterFilter,
				)
				if err != nil {
					wcLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				wcLogger.Infof(
					"subscription monitoring fetched [%v] past ProposalSubmitterRemoved events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := psrs.contract.watchProposalSubmitterRemoved(
		sink,
		psrs.proposalSubmitterFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (wc *WalletCoordinator) watchProposalSubmitterRemoved(
	sink chan *abi.WalletCoordinatorProposalSubmitterRemoved,
	proposalSubmitterFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return wc.contract.WatchProposalSubmitterRemoved(
			&bind.WatchOpts{Context: ctx},
			sink,
			proposalSubmitterFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		wcLogger.Warnf(
			"subscription to event ProposalSubmitterRemoved had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		wcLogger.Errorf(
			"subscription to event ProposalSubmitterRemoved failed "+
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

func (wc *WalletCoordinator) PastProposalSubmitterRemovedEvents(
	startBlock uint64,
	endBlock *uint64,
	proposalSubmitterFilter []common.Address,
) ([]*abi.WalletCoordinatorProposalSubmitterRemoved, error) {
	iterator, err := wc.contract.FilterProposalSubmitterRemoved(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		proposalSubmitterFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past ProposalSubmitterRemoved events: [%v]",
			err,
		)
	}

	events := make([]*abi.WalletCoordinatorProposalSubmitterRemoved, 0)

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
