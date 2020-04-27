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
var krboLogger = log.Logger("keep-contract-KeepRandomBeaconOperator")

type KeepRandomBeaconOperator struct {
	contract          *abi.KeepRandomBeaconOperator
	contractAddress   common.Address
	contractABI       *ethereumabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *ethutil.ErrorResolver

	transactionMutex *sync.Mutex
}

func NewKeepRandomBeaconOperator(
	contractAddress common.Address,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	transactionMutex *sync.Mutex,
) (*KeepRandomBeaconOperator, error) {
	callerOptions := &bind.CallOpts{
		From: contractAddress,
	}
	transactorOptions := bind.NewKeyedTransactor(
		accountKey.PrivateKey,
	)

	randomBeaconContract, err := abi.NewKeepRandomBeaconOperator(
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

	contractABI, err := ethereumabi.JSON(strings.NewReader(abi.KeepRandomBeaconOperatorABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &KeepRandomBeaconOperator{
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
func (krbo *KeepRandomBeaconOperator) SubmitDkgResult(
	submitterMemberIndex *big.Int,
	groupPubKey []uint8,
	misbehaved []uint8,
	signatures []uint8,
	signingMembersIndexes []*big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction submitDkgResult",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction submitDkgResult with id: [%v]",
		transaction.Hash().Hex(),
	)

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

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) ReportRelayEntryTimeout(

	transactionOptions ...ethutil.TransactionOptions,
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

	krboLogger.Debugf(
		"submitted transaction reportRelayEntryTimeout with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallReportRelayEntryTimeout(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction reportUnauthorizedSigning",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction reportUnauthorizedSigning with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallReportUnauthorizedSigning(
	groupIndex *big.Int,
	signedMsgSender []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) RelayEntry(
	_groupSignature []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction relayEntry",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction relayEntry with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallRelayEntry(
	_groupSignature []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) RemoveServiceContract(
	serviceContract common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction removeServiceContract",
		"params: ",
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

	transaction, err := krbo.contract.RemoveServiceContract(
		transactorOptions,
		serviceContract,
	)

	if err != nil {
		return transaction, krbo.errorResolver.ResolveError(
			err,
			krbo.transactorOptions.From,
			nil,
			"removeServiceContract",
			serviceContract,
		)
	}

	krboLogger.Debugf(
		"submitted transaction removeServiceContract with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallRemoveServiceContract(
	serviceContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		krbo.transactorOptions.From,
		blockNumber, nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"removeServiceContract",
		&result,
		serviceContract,
	)

	return err
}

func (krbo *KeepRandomBeaconOperator) RemoveServiceContractGasEstimate(
	serviceContract common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"removeServiceContract",
		krbo.contractABI,
		krbo.transactor,
		serviceContract,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) AddServiceContract(
	serviceContract common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction addServiceContract",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction addServiceContract with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallAddServiceContract(
	serviceContract common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) Genesis(
	value *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction genesis",
		"value: ", value,
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

	krboLogger.Debugf(
		"submitted transaction genesis with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallGenesis(
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
		krbo.callerOptions.From,
		krbo.contractAddress,
		"genesis",
		krbo.contractABI,
		krbo.transactor,
	)

	return result, err
}

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) WithdrawGroupMemberRewards(
	operator common.Address,
	groupIndex *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction withdrawGroupMemberRewards",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction withdrawGroupMemberRewards with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallWithdrawGroupMemberRewards(
	operator common.Address,
	groupIndex *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

// Transaction submission.
func (krbo *KeepRandomBeaconOperator) SubmitTicket(
	ticket [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction submitTicket",
		"params: ",
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

	krboLogger.Debugf(
		"submitted transaction submitTicket with id: [%v]",
		transaction.Hash().Hex(),
	)

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (krbo *KeepRandomBeaconOperator) CallSubmitTicket(
	ticket [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) CreateGroup(
	_newEntry *big.Int,
	submitter common.Address,
	value *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction createGroup",
		"params: ",
		fmt.Sprint(
			_newEntry,
			submitter,
		),
		"value: ", value,
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

	krboLogger.Debugf(
		"submitted transaction createGroup with id: [%v]",
		transaction.Hash().Hex(),
	)

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

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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
func (krbo *KeepRandomBeaconOperator) Sign(
	requestId *big.Int,
	previousEntry []uint8,
	value *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	krboLogger.Debug(
		"submitting transaction sign",
		"params: ",
		fmt.Sprint(
			requestId,
			previousEntry,
		),
		"value: ", value,
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

	krboLogger.Debugf(
		"submitted transaction sign with id: [%v]",
		transaction.Hash().Hex(),
	)

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

	err := ethutil.CallAtBlock(
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

	result, err := ethutil.EstimateGas(
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

// ----- Const Methods ------

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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

func (krbo *KeepRandomBeaconOperator) CurrentEntryStartBlock() (*big.Int, error) {
	var result *big.Int
	result, err := krbo.contract.CurrentEntryStartBlock(
		krbo.callerOptions,
	)

	if err != nil {
		return result, krbo.errorResolver.ResolveError(
			err,
			krbo.callerOptions.From,
			nil,
			"currentEntryStartBlock",
		)
	}

	return result, err
}

func (krbo *KeepRandomBeaconOperator) CurrentEntryStartBlockAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		krbo.callerOptions.From,
		blockNumber,
		nil,
		krbo.contractABI,
		krbo.caller,
		krbo.errorResolver,
		krbo.contractAddress,
		"currentEntryStartBlock",
		&result,
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

	err := ethutil.CallAtBlock(
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

// ------ Events -------

type keepRandomBeaconOperatorDkgResultSubmittedEventFunc func(
	MemberIndex *big.Int,
	GroupPubKey []uint8,
	Misbehaved []uint8,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchDkgResultSubmittedEvent(
	success keepRandomBeaconOperatorDkgResultSubmittedEventFunc,
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

		subscription, err := krbo.subscribeDkgResultSubmittedEvent(
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
				krboLogger.Warning(
					"subscription to event DkgResultSubmittedEvent terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeDkgResultSubmittedEvent(
	success keepRandomBeaconOperatorDkgResultSubmittedEventFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorDkgResultSubmittedEvent)
	eventSubscription, err := krbo.contract.WatchDkgResultSubmittedEvent(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for DkgResultSubmittedEvent events: [%v]",
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
					event.MemberIndex,
					event.GroupPubKey,
					event.Misbehaved,
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

type keepRandomBeaconOperatorGroupMemberRewardsWithdrawnFunc func(
	Beneficiary common.Address,
	Operator common.Address,
	Amount *big.Int,
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchGroupMemberRewardsWithdrawn(
	success keepRandomBeaconOperatorGroupMemberRewardsWithdrawnFunc,
	fail func(err error) error,
	beneficiaryFilter []common.Address,
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

		subscription, err := krbo.subscribeGroupMemberRewardsWithdrawn(
			success,
			failCallback,
			beneficiaryFilter,
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
				krboLogger.Warning(
					"subscription to event GroupMemberRewardsWithdrawn terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeGroupMemberRewardsWithdrawn(
	success keepRandomBeaconOperatorGroupMemberRewardsWithdrawnFunc,
	fail func(err error) error,
	beneficiaryFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
	eventSubscription, err := krbo.contract.WatchGroupMemberRewardsWithdrawn(
		nil,
		eventChan,
		beneficiaryFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for GroupMemberRewardsWithdrawn events: [%v]",
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
					event.Beneficiary,
					event.Operator,
					event.Amount,
					event.GroupIndex,
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

type keepRandomBeaconOperatorGroupSelectionStartedFunc func(
	NewEntry *big.Int,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchGroupSelectionStarted(
	success keepRandomBeaconOperatorGroupSelectionStartedFunc,
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

		subscription, err := krbo.subscribeGroupSelectionStarted(
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
				krboLogger.Warning(
					"subscription to event GroupSelectionStarted terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeGroupSelectionStarted(
	success keepRandomBeaconOperatorGroupSelectionStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorGroupSelectionStarted)
	eventSubscription, err := krbo.contract.WatchGroupSelectionStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for GroupSelectionStarted events: [%v]",
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
					event.NewEntry,
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

type keepRandomBeaconOperatorOnGroupRegisteredFunc func(
	GroupPubKey []uint8,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchOnGroupRegistered(
	success keepRandomBeaconOperatorOnGroupRegisteredFunc,
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

		subscription, err := krbo.subscribeOnGroupRegistered(
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
				krboLogger.Warning(
					"subscription to event OnGroupRegistered terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeOnGroupRegistered(
	success keepRandomBeaconOperatorOnGroupRegisteredFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorOnGroupRegistered)
	eventSubscription, err := krbo.contract.WatchOnGroupRegistered(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for OnGroupRegistered events: [%v]",
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
					event.GroupPubKey,
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

type keepRandomBeaconOperatorRelayEntryRequestedFunc func(
	PreviousEntry []uint8,
	GroupPublicKey []uint8,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchRelayEntryRequested(
	success keepRandomBeaconOperatorRelayEntryRequestedFunc,
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

		subscription, err := krbo.subscribeRelayEntryRequested(
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
				krboLogger.Warning(
					"subscription to event RelayEntryRequested terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeRelayEntryRequested(
	success keepRandomBeaconOperatorRelayEntryRequestedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntryRequested)
	eventSubscription, err := krbo.contract.WatchRelayEntryRequested(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RelayEntryRequested events: [%v]",
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
					event.PreviousEntry,
					event.GroupPublicKey,
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

type keepRandomBeaconOperatorRelayEntrySubmittedFunc func(
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchRelayEntrySubmitted(
	success keepRandomBeaconOperatorRelayEntrySubmittedFunc,
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

		subscription, err := krbo.subscribeRelayEntrySubmitted(
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
				krboLogger.Warning(
					"subscription to event RelayEntrySubmitted terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeRelayEntrySubmitted(
	success keepRandomBeaconOperatorRelayEntrySubmittedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntrySubmitted)
	eventSubscription, err := krbo.contract.WatchRelayEntrySubmitted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RelayEntrySubmitted events: [%v]",
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

type keepRandomBeaconOperatorRelayEntryTimeoutReportedFunc func(
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchRelayEntryTimeoutReported(
	success keepRandomBeaconOperatorRelayEntryTimeoutReportedFunc,
	fail func(err error) error,
	groupIndexFilter []*big.Int,
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

		subscription, err := krbo.subscribeRelayEntryTimeoutReported(
			success,
			failCallback,
			groupIndexFilter,
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
				krboLogger.Warning(
					"subscription to event RelayEntryTimeoutReported terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeRelayEntryTimeoutReported(
	success keepRandomBeaconOperatorRelayEntryTimeoutReportedFunc,
	fail func(err error) error,
	groupIndexFilter []*big.Int,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorRelayEntryTimeoutReported)
	eventSubscription, err := krbo.contract.WatchRelayEntryTimeoutReported(
		nil,
		eventChan,
		groupIndexFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RelayEntryTimeoutReported events: [%v]",
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
					event.GroupIndex,
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

type keepRandomBeaconOperatorUnauthorizedSigningReportedFunc func(
	GroupIndex *big.Int,
	blockNumber uint64,
)

func (krbo *KeepRandomBeaconOperator) WatchUnauthorizedSigningReported(
	success keepRandomBeaconOperatorUnauthorizedSigningReportedFunc,
	fail func(err error) error,
	groupIndexFilter []*big.Int,
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

		subscription, err := krbo.subscribeUnauthorizedSigningReported(
			success,
			failCallback,
			groupIndexFilter,
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
				krboLogger.Warning(
					"subscription to event UnauthorizedSigningReported terminated with error; " +
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

func (krbo *KeepRandomBeaconOperator) subscribeUnauthorizedSigningReported(
	success keepRandomBeaconOperatorUnauthorizedSigningReportedFunc,
	fail func(err error) error,
	groupIndexFilter []*big.Int,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.KeepRandomBeaconOperatorUnauthorizedSigningReported)
	eventSubscription, err := krbo.contract.WatchUnauthorizedSigningReported(
		nil,
		eventChan,
		groupIndexFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for UnauthorizedSigningReported events: [%v]",
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
					event.GroupIndex,
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
