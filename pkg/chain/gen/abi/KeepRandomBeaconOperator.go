// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// KeepRandomBeaconOperatorABI is the input ABI used to generate the binding from.
const KeepRandomBeaconOperatorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_registryContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"memberIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"misbehaved\",\"type\":\"bytes\"}],\"name\":\"DkgResultSubmittedEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"GroupMemberRewardsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newEntry\",\"type\":\"uint256\"}],\"name\":\"GroupSelectionStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"}],\"name\":\"OnGroupRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousEntry\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"groupPublicKey\",\"type\":\"bytes\"}],\"name\":\"RelayEntryRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RelayEntrySubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"RelayEntryTimeoutReported\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"UnauthorizedSigningReported\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"serviceContract\",\"type\":\"address\"}],\"name\":\"addServiceContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newEntry\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"submitter\",\"type\":\"address\"}],\"name\":\"createGroup\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEntryStartBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentRequestGroupIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentRequestPreviousEntry\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dkgGasEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dkgSubmitterReimbursementFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"entryVerificationFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"entryVerificationGasEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasPriceCeiling\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"genesis\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getFirstActiveGroupIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"}],\"name\":\"getGroupMemberRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"}],\"name\":\"getGroupMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"getGroupPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupCreationFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupMemberBaseReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupProfitFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupSelectionGasEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"groupThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"}],\"name\":\"hasMinimumStake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"hasWithdrawnRewards\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isEntryInProgress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"}],\"name\":\"isGroupRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isGroupSelectionPossible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"}],\"name\":\"isStaleGroup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numberOfGroups\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_groupSignature\",\"type\":\"bytes\"}],\"name\":\"relayEntry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"relayEntryTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"serviceContract\",\"type\":\"address\"}],\"name\":\"removeServiceContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"reportRelayEntryTimeout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signedMsgSender\",\"type\":\"bytes\"}],\"name\":\"reportUnauthorizedSigning\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"resultPublicationBlockStep\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"selectedParticipants\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"previousEntry\",\"type\":\"bytes\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"submitterMemberIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"groupPubKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"misbehaved\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"signingMembersIndexes\",\"type\":\"uint256[]\"}],\"name\":\"submitDkgResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ticket\",\"type\":\"bytes32\"}],\"name\":\"submitTicket\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"submittedTickets\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ticketSubmissionTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"groupIndex\",\"type\":\"uint256\"}],\"name\":\"withdrawGroupMemberRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KeepRandomBeaconOperator is an auto generated Go binding around an Ethereum contract.
type KeepRandomBeaconOperator struct {
	KeepRandomBeaconOperatorCaller     // Read-only binding to the contract
	KeepRandomBeaconOperatorTransactor // Write-only binding to the contract
	KeepRandomBeaconOperatorFilterer   // Log filterer for contract events
}

// KeepRandomBeaconOperatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeepRandomBeaconOperatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconOperatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeepRandomBeaconOperatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconOperatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeepRandomBeaconOperatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconOperatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeepRandomBeaconOperatorSession struct {
	Contract     *KeepRandomBeaconOperator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// KeepRandomBeaconOperatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeepRandomBeaconOperatorCallerSession struct {
	Contract *KeepRandomBeaconOperatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// KeepRandomBeaconOperatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeepRandomBeaconOperatorTransactorSession struct {
	Contract     *KeepRandomBeaconOperatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// KeepRandomBeaconOperatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeepRandomBeaconOperatorRaw struct {
	Contract *KeepRandomBeaconOperator // Generic contract binding to access the raw methods on
}

// KeepRandomBeaconOperatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeepRandomBeaconOperatorCallerRaw struct {
	Contract *KeepRandomBeaconOperatorCaller // Generic read-only contract binding to access the raw methods on
}

// KeepRandomBeaconOperatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeepRandomBeaconOperatorTransactorRaw struct {
	Contract *KeepRandomBeaconOperatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeepRandomBeaconOperator creates a new instance of KeepRandomBeaconOperator, bound to a specific deployed contract.
func NewKeepRandomBeaconOperator(address common.Address, backend bind.ContractBackend) (*KeepRandomBeaconOperator, error) {
	contract, err := bindKeepRandomBeaconOperator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperator{KeepRandomBeaconOperatorCaller: KeepRandomBeaconOperatorCaller{contract: contract}, KeepRandomBeaconOperatorTransactor: KeepRandomBeaconOperatorTransactor{contract: contract}, KeepRandomBeaconOperatorFilterer: KeepRandomBeaconOperatorFilterer{contract: contract}}, nil
}

// NewKeepRandomBeaconOperatorCaller creates a new read-only instance of KeepRandomBeaconOperator, bound to a specific deployed contract.
func NewKeepRandomBeaconOperatorCaller(address common.Address, caller bind.ContractCaller) (*KeepRandomBeaconOperatorCaller, error) {
	contract, err := bindKeepRandomBeaconOperator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorCaller{contract: contract}, nil
}

// NewKeepRandomBeaconOperatorTransactor creates a new write-only instance of KeepRandomBeaconOperator, bound to a specific deployed contract.
func NewKeepRandomBeaconOperatorTransactor(address common.Address, transactor bind.ContractTransactor) (*KeepRandomBeaconOperatorTransactor, error) {
	contract, err := bindKeepRandomBeaconOperator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorTransactor{contract: contract}, nil
}

// NewKeepRandomBeaconOperatorFilterer creates a new log filterer instance of KeepRandomBeaconOperator, bound to a specific deployed contract.
func NewKeepRandomBeaconOperatorFilterer(address common.Address, filterer bind.ContractFilterer) (*KeepRandomBeaconOperatorFilterer, error) {
	contract, err := bindKeepRandomBeaconOperator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorFilterer{contract: contract}, nil
}

// bindKeepRandomBeaconOperator binds a generic wrapper to an already deployed contract.
func bindKeepRandomBeaconOperator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeepRandomBeaconOperatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconOperator.Contract.KeepRandomBeaconOperatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.KeepRandomBeaconOperatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.KeepRandomBeaconOperatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconOperator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.contract.Transact(opts, method, params...)
}

// CurrentEntryStartBlock is a free data retrieval call binding the contract method 0xb0dfe104.
//
// Solidity: function currentEntryStartBlock() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) CurrentEntryStartBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "currentEntryStartBlock")
	return *ret0, err
}

// CurrentEntryStartBlock is a free data retrieval call binding the contract method 0xb0dfe104.
//
// Solidity: function currentEntryStartBlock() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) CurrentEntryStartBlock() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentEntryStartBlock(&_KeepRandomBeaconOperator.CallOpts)
}

// CurrentEntryStartBlock is a free data retrieval call binding the contract method 0xb0dfe104.
//
// Solidity: function currentEntryStartBlock() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) CurrentEntryStartBlock() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentEntryStartBlock(&_KeepRandomBeaconOperator.CallOpts)
}

// CurrentRequestGroupIndex is a free data retrieval call binding the contract method 0x7031b7ff.
//
// Solidity: function currentRequestGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) CurrentRequestGroupIndex(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "currentRequestGroupIndex")
	return *ret0, err
}

// CurrentRequestGroupIndex is a free data retrieval call binding the contract method 0x7031b7ff.
//
// Solidity: function currentRequestGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) CurrentRequestGroupIndex() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentRequestGroupIndex(&_KeepRandomBeaconOperator.CallOpts)
}

// CurrentRequestGroupIndex is a free data retrieval call binding the contract method 0x7031b7ff.
//
// Solidity: function currentRequestGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) CurrentRequestGroupIndex() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentRequestGroupIndex(&_KeepRandomBeaconOperator.CallOpts)
}

// CurrentRequestPreviousEntry is a free data retrieval call binding the contract method 0x618c2656.
//
// Solidity: function currentRequestPreviousEntry() constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) CurrentRequestPreviousEntry(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "currentRequestPreviousEntry")
	return *ret0, err
}

// CurrentRequestPreviousEntry is a free data retrieval call binding the contract method 0x618c2656.
//
// Solidity: function currentRequestPreviousEntry() constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) CurrentRequestPreviousEntry() ([]byte, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentRequestPreviousEntry(&_KeepRandomBeaconOperator.CallOpts)
}

// CurrentRequestPreviousEntry is a free data retrieval call binding the contract method 0x618c2656.
//
// Solidity: function currentRequestPreviousEntry() constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) CurrentRequestPreviousEntry() ([]byte, error) {
	return _KeepRandomBeaconOperator.Contract.CurrentRequestPreviousEntry(&_KeepRandomBeaconOperator.CallOpts)
}

// DkgGasEstimate is a free data retrieval call binding the contract method 0x003bf87e.
//
// Solidity: function dkgGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) DkgGasEstimate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "dkgGasEstimate")
	return *ret0, err
}

// DkgGasEstimate is a free data retrieval call binding the contract method 0x003bf87e.
//
// Solidity: function dkgGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) DkgGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.DkgGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// DkgGasEstimate is a free data retrieval call binding the contract method 0x003bf87e.
//
// Solidity: function dkgGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) DkgGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.DkgGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// DkgSubmitterReimbursementFee is a free data retrieval call binding the contract method 0xb1c77c8f.
//
// Solidity: function dkgSubmitterReimbursementFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) DkgSubmitterReimbursementFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "dkgSubmitterReimbursementFee")
	return *ret0, err
}

// DkgSubmitterReimbursementFee is a free data retrieval call binding the contract method 0xb1c77c8f.
//
// Solidity: function dkgSubmitterReimbursementFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) DkgSubmitterReimbursementFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.DkgSubmitterReimbursementFee(&_KeepRandomBeaconOperator.CallOpts)
}

// DkgSubmitterReimbursementFee is a free data retrieval call binding the contract method 0xb1c77c8f.
//
// Solidity: function dkgSubmitterReimbursementFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) DkgSubmitterReimbursementFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.DkgSubmitterReimbursementFee(&_KeepRandomBeaconOperator.CallOpts)
}

// EntryVerificationFee is a free data retrieval call binding the contract method 0x517471a9.
//
// Solidity: function entryVerificationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) EntryVerificationFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "entryVerificationFee")
	return *ret0, err
}

// EntryVerificationFee is a free data retrieval call binding the contract method 0x517471a9.
//
// Solidity: function entryVerificationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) EntryVerificationFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.EntryVerificationFee(&_KeepRandomBeaconOperator.CallOpts)
}

// EntryVerificationFee is a free data retrieval call binding the contract method 0x517471a9.
//
// Solidity: function entryVerificationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) EntryVerificationFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.EntryVerificationFee(&_KeepRandomBeaconOperator.CallOpts)
}

// EntryVerificationGasEstimate is a free data retrieval call binding the contract method 0x79f9fb7e.
//
// Solidity: function entryVerificationGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) EntryVerificationGasEstimate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "entryVerificationGasEstimate")
	return *ret0, err
}

// EntryVerificationGasEstimate is a free data retrieval call binding the contract method 0x79f9fb7e.
//
// Solidity: function entryVerificationGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) EntryVerificationGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.EntryVerificationGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// EntryVerificationGasEstimate is a free data retrieval call binding the contract method 0x79f9fb7e.
//
// Solidity: function entryVerificationGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) EntryVerificationGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.EntryVerificationGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// GasPriceCeiling is a free data retrieval call binding the contract method 0xe1f4d632.
//
// Solidity: function gasPriceCeiling() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GasPriceCeiling(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "gasPriceCeiling")
	return *ret0, err
}

// GasPriceCeiling is a free data retrieval call binding the contract method 0xe1f4d632.
//
// Solidity: function gasPriceCeiling() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GasPriceCeiling() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GasPriceCeiling(&_KeepRandomBeaconOperator.CallOpts)
}

// GasPriceCeiling is a free data retrieval call binding the contract method 0xe1f4d632.
//
// Solidity: function gasPriceCeiling() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GasPriceCeiling() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GasPriceCeiling(&_KeepRandomBeaconOperator.CallOpts)
}

// GetFirstActiveGroupIndex is a free data retrieval call binding the contract method 0xeb9488d3.
//
// Solidity: function getFirstActiveGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GetFirstActiveGroupIndex(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "getFirstActiveGroupIndex")
	return *ret0, err
}

// GetFirstActiveGroupIndex is a free data retrieval call binding the contract method 0xeb9488d3.
//
// Solidity: function getFirstActiveGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GetFirstActiveGroupIndex() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GetFirstActiveGroupIndex(&_KeepRandomBeaconOperator.CallOpts)
}

// GetFirstActiveGroupIndex is a free data retrieval call binding the contract method 0xeb9488d3.
//
// Solidity: function getFirstActiveGroupIndex() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GetFirstActiveGroupIndex() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GetFirstActiveGroupIndex(&_KeepRandomBeaconOperator.CallOpts)
}

// GetGroupMemberRewards is a free data retrieval call binding the contract method 0x9dabee44.
//
// Solidity: function getGroupMemberRewards(bytes groupPubKey) constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GetGroupMemberRewards(opts *bind.CallOpts, groupPubKey []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "getGroupMemberRewards", groupPubKey)
	return *ret0, err
}

// GetGroupMemberRewards is a free data retrieval call binding the contract method 0x9dabee44.
//
// Solidity: function getGroupMemberRewards(bytes groupPubKey) constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GetGroupMemberRewards(groupPubKey []byte) (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupMemberRewards(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// GetGroupMemberRewards is a free data retrieval call binding the contract method 0x9dabee44.
//
// Solidity: function getGroupMemberRewards(bytes groupPubKey) constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GetGroupMemberRewards(groupPubKey []byte) (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupMemberRewards(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// GetGroupMembers is a free data retrieval call binding the contract method 0xd12f5e69.
//
// Solidity: function getGroupMembers(bytes groupPubKey) constant returns(address[] members)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GetGroupMembers(opts *bind.CallOpts, groupPubKey []byte) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "getGroupMembers", groupPubKey)
	return *ret0, err
}

// GetGroupMembers is a free data retrieval call binding the contract method 0xd12f5e69.
//
// Solidity: function getGroupMembers(bytes groupPubKey) constant returns(address[] members)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GetGroupMembers(groupPubKey []byte) ([]common.Address, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupMembers(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// GetGroupMembers is a free data retrieval call binding the contract method 0xd12f5e69.
//
// Solidity: function getGroupMembers(bytes groupPubKey) constant returns(address[] members)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GetGroupMembers(groupPubKey []byte) ([]common.Address, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupMembers(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// GetGroupPublicKey is a free data retrieval call binding the contract method 0xef7c8f9c.
//
// Solidity: function getGroupPublicKey(uint256 groupIndex) constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GetGroupPublicKey(opts *bind.CallOpts, groupIndex *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "getGroupPublicKey", groupIndex)
	return *ret0, err
}

// GetGroupPublicKey is a free data retrieval call binding the contract method 0xef7c8f9c.
//
// Solidity: function getGroupPublicKey(uint256 groupIndex) constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GetGroupPublicKey(groupIndex *big.Int) ([]byte, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupPublicKey(&_KeepRandomBeaconOperator.CallOpts, groupIndex)
}

// GetGroupPublicKey is a free data retrieval call binding the contract method 0xef7c8f9c.
//
// Solidity: function getGroupPublicKey(uint256 groupIndex) constant returns(bytes)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GetGroupPublicKey(groupIndex *big.Int) ([]byte, error) {
	return _KeepRandomBeaconOperator.Contract.GetGroupPublicKey(&_KeepRandomBeaconOperator.CallOpts, groupIndex)
}

// GroupCreationFee is a free data retrieval call binding the contract method 0xc300d058.
//
// Solidity: function groupCreationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupCreationFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupCreationFee")
	return *ret0, err
}

// GroupCreationFee is a free data retrieval call binding the contract method 0xc300d058.
//
// Solidity: function groupCreationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupCreationFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupCreationFee(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupCreationFee is a free data retrieval call binding the contract method 0xc300d058.
//
// Solidity: function groupCreationFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupCreationFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupCreationFee(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupMemberBaseReward is a free data retrieval call binding the contract method 0x7d7d7dd9.
//
// Solidity: function groupMemberBaseReward() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupMemberBaseReward(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupMemberBaseReward")
	return *ret0, err
}

// GroupMemberBaseReward is a free data retrieval call binding the contract method 0x7d7d7dd9.
//
// Solidity: function groupMemberBaseReward() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupMemberBaseReward() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupMemberBaseReward(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupMemberBaseReward is a free data retrieval call binding the contract method 0x7d7d7dd9.
//
// Solidity: function groupMemberBaseReward() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupMemberBaseReward() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupMemberBaseReward(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupProfitFee is a free data retrieval call binding the contract method 0xc4438946.
//
// Solidity: function groupProfitFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupProfitFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupProfitFee")
	return *ret0, err
}

// GroupProfitFee is a free data retrieval call binding the contract method 0xc4438946.
//
// Solidity: function groupProfitFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupProfitFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupProfitFee(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupProfitFee is a free data retrieval call binding the contract method 0xc4438946.
//
// Solidity: function groupProfitFee() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupProfitFee() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupProfitFee(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupSelectionGasEstimate is a free data retrieval call binding the contract method 0x24f17313.
//
// Solidity: function groupSelectionGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupSelectionGasEstimate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupSelectionGasEstimate")
	return *ret0, err
}

// GroupSelectionGasEstimate is a free data retrieval call binding the contract method 0x24f17313.
//
// Solidity: function groupSelectionGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupSelectionGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupSelectionGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupSelectionGasEstimate is a free data retrieval call binding the contract method 0x24f17313.
//
// Solidity: function groupSelectionGasEstimate() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupSelectionGasEstimate() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupSelectionGasEstimate(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupSize is a free data retrieval call binding the contract method 0x63b635ea.
//
// Solidity: function groupSize() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupSize")
	return *ret0, err
}

// GroupSize is a free data retrieval call binding the contract method 0x63b635ea.
//
// Solidity: function groupSize() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupSize() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupSize(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupSize is a free data retrieval call binding the contract method 0x63b635ea.
//
// Solidity: function groupSize() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupSize() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupSize(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupThreshold is a free data retrieval call binding the contract method 0x6dcc64f8.
//
// Solidity: function groupThreshold() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) GroupThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "groupThreshold")
	return *ret0, err
}

// GroupThreshold is a free data retrieval call binding the contract method 0x6dcc64f8.
//
// Solidity: function groupThreshold() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) GroupThreshold() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupThreshold(&_KeepRandomBeaconOperator.CallOpts)
}

// GroupThreshold is a free data retrieval call binding the contract method 0x6dcc64f8.
//
// Solidity: function groupThreshold() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) GroupThreshold() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.GroupThreshold(&_KeepRandomBeaconOperator.CallOpts)
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x5c1c0710.
//
// Solidity: function hasMinimumStake(address staker) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) HasMinimumStake(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "hasMinimumStake", staker)
	return *ret0, err
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x5c1c0710.
//
// Solidity: function hasMinimumStake(address staker) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) HasMinimumStake(staker common.Address) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.HasMinimumStake(&_KeepRandomBeaconOperator.CallOpts, staker)
}

// HasMinimumStake is a free data retrieval call binding the contract method 0x5c1c0710.
//
// Solidity: function hasMinimumStake(address staker) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) HasMinimumStake(staker common.Address) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.HasMinimumStake(&_KeepRandomBeaconOperator.CallOpts, staker)
}

// HasWithdrawnRewards is a free data retrieval call binding the contract method 0x376f7a11.
//
// Solidity: function hasWithdrawnRewards(address operator, uint256 groupIndex) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) HasWithdrawnRewards(opts *bind.CallOpts, operator common.Address, groupIndex *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "hasWithdrawnRewards", operator, groupIndex)
	return *ret0, err
}

// HasWithdrawnRewards is a free data retrieval call binding the contract method 0x376f7a11.
//
// Solidity: function hasWithdrawnRewards(address operator, uint256 groupIndex) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) HasWithdrawnRewards(operator common.Address, groupIndex *big.Int) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.HasWithdrawnRewards(&_KeepRandomBeaconOperator.CallOpts, operator, groupIndex)
}

// HasWithdrawnRewards is a free data retrieval call binding the contract method 0x376f7a11.
//
// Solidity: function hasWithdrawnRewards(address operator, uint256 groupIndex) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) HasWithdrawnRewards(operator common.Address, groupIndex *big.Int) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.HasWithdrawnRewards(&_KeepRandomBeaconOperator.CallOpts, operator, groupIndex)
}

// IsEntryInProgress is a free data retrieval call binding the contract method 0x6e5636e4.
//
// Solidity: function isEntryInProgress() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) IsEntryInProgress(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "isEntryInProgress")
	return *ret0, err
}

// IsEntryInProgress is a free data retrieval call binding the contract method 0x6e5636e4.
//
// Solidity: function isEntryInProgress() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) IsEntryInProgress() (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsEntryInProgress(&_KeepRandomBeaconOperator.CallOpts)
}

// IsEntryInProgress is a free data retrieval call binding the contract method 0x6e5636e4.
//
// Solidity: function isEntryInProgress() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) IsEntryInProgress() (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsEntryInProgress(&_KeepRandomBeaconOperator.CallOpts)
}

// IsGroupRegistered is a free data retrieval call binding the contract method 0x1c524ac2.
//
// Solidity: function isGroupRegistered(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) IsGroupRegistered(opts *bind.CallOpts, groupPubKey []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "isGroupRegistered", groupPubKey)
	return *ret0, err
}

// IsGroupRegistered is a free data retrieval call binding the contract method 0x1c524ac2.
//
// Solidity: function isGroupRegistered(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) IsGroupRegistered(groupPubKey []byte) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsGroupRegistered(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// IsGroupRegistered is a free data retrieval call binding the contract method 0x1c524ac2.
//
// Solidity: function isGroupRegistered(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) IsGroupRegistered(groupPubKey []byte) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsGroupRegistered(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// IsGroupSelectionPossible is a free data retrieval call binding the contract method 0x21a8f86c.
//
// Solidity: function isGroupSelectionPossible() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) IsGroupSelectionPossible(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "isGroupSelectionPossible")
	return *ret0, err
}

// IsGroupSelectionPossible is a free data retrieval call binding the contract method 0x21a8f86c.
//
// Solidity: function isGroupSelectionPossible() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) IsGroupSelectionPossible() (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsGroupSelectionPossible(&_KeepRandomBeaconOperator.CallOpts)
}

// IsGroupSelectionPossible is a free data retrieval call binding the contract method 0x21a8f86c.
//
// Solidity: function isGroupSelectionPossible() constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) IsGroupSelectionPossible() (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsGroupSelectionPossible(&_KeepRandomBeaconOperator.CallOpts)
}

// IsStaleGroup is a free data retrieval call binding the contract method 0x2d6f8f31.
//
// Solidity: function isStaleGroup(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) IsStaleGroup(opts *bind.CallOpts, groupPubKey []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "isStaleGroup", groupPubKey)
	return *ret0, err
}

// IsStaleGroup is a free data retrieval call binding the contract method 0x2d6f8f31.
//
// Solidity: function isStaleGroup(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) IsStaleGroup(groupPubKey []byte) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsStaleGroup(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// IsStaleGroup is a free data retrieval call binding the contract method 0x2d6f8f31.
//
// Solidity: function isStaleGroup(bytes groupPubKey) constant returns(bool)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) IsStaleGroup(groupPubKey []byte) (bool, error) {
	return _KeepRandomBeaconOperator.Contract.IsStaleGroup(&_KeepRandomBeaconOperator.CallOpts, groupPubKey)
}

// NumberOfGroups is a free data retrieval call binding the contract method 0xbf952496.
//
// Solidity: function numberOfGroups() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) NumberOfGroups(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "numberOfGroups")
	return *ret0, err
}

// NumberOfGroups is a free data retrieval call binding the contract method 0xbf952496.
//
// Solidity: function numberOfGroups() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) NumberOfGroups() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.NumberOfGroups(&_KeepRandomBeaconOperator.CallOpts)
}

// NumberOfGroups is a free data retrieval call binding the contract method 0xbf952496.
//
// Solidity: function numberOfGroups() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) NumberOfGroups() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.NumberOfGroups(&_KeepRandomBeaconOperator.CallOpts)
}

// RelayEntryTimeout is a free data retrieval call binding the contract method 0xb99f0c43.
//
// Solidity: function relayEntryTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) RelayEntryTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "relayEntryTimeout")
	return *ret0, err
}

// RelayEntryTimeout is a free data retrieval call binding the contract method 0xb99f0c43.
//
// Solidity: function relayEntryTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) RelayEntryTimeout() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.RelayEntryTimeout(&_KeepRandomBeaconOperator.CallOpts)
}

// RelayEntryTimeout is a free data retrieval call binding the contract method 0xb99f0c43.
//
// Solidity: function relayEntryTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) RelayEntryTimeout() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.RelayEntryTimeout(&_KeepRandomBeaconOperator.CallOpts)
}

// ResultPublicationBlockStep is a free data retrieval call binding the contract method 0x36c85717.
//
// Solidity: function resultPublicationBlockStep() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) ResultPublicationBlockStep(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "resultPublicationBlockStep")
	return *ret0, err
}

// ResultPublicationBlockStep is a free data retrieval call binding the contract method 0x36c85717.
//
// Solidity: function resultPublicationBlockStep() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) ResultPublicationBlockStep() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.ResultPublicationBlockStep(&_KeepRandomBeaconOperator.CallOpts)
}

// ResultPublicationBlockStep is a free data retrieval call binding the contract method 0x36c85717.
//
// Solidity: function resultPublicationBlockStep() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) ResultPublicationBlockStep() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.ResultPublicationBlockStep(&_KeepRandomBeaconOperator.CallOpts)
}

// SelectedParticipants is a free data retrieval call binding the contract method 0x0b19991f.
//
// Solidity: function selectedParticipants() constant returns(address[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) SelectedParticipants(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "selectedParticipants")
	return *ret0, err
}

// SelectedParticipants is a free data retrieval call binding the contract method 0x0b19991f.
//
// Solidity: function selectedParticipants() constant returns(address[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) SelectedParticipants() ([]common.Address, error) {
	return _KeepRandomBeaconOperator.Contract.SelectedParticipants(&_KeepRandomBeaconOperator.CallOpts)
}

// SelectedParticipants is a free data retrieval call binding the contract method 0x0b19991f.
//
// Solidity: function selectedParticipants() constant returns(address[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) SelectedParticipants() ([]common.Address, error) {
	return _KeepRandomBeaconOperator.Contract.SelectedParticipants(&_KeepRandomBeaconOperator.CallOpts)
}

// SubmittedTickets is a free data retrieval call binding the contract method 0x6262d54e.
//
// Solidity: function submittedTickets() constant returns(uint64[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) SubmittedTickets(opts *bind.CallOpts) ([]uint64, error) {
	var (
		ret0 = new([]uint64)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "submittedTickets")
	return *ret0, err
}

// SubmittedTickets is a free data retrieval call binding the contract method 0x6262d54e.
//
// Solidity: function submittedTickets() constant returns(uint64[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) SubmittedTickets() ([]uint64, error) {
	return _KeepRandomBeaconOperator.Contract.SubmittedTickets(&_KeepRandomBeaconOperator.CallOpts)
}

// SubmittedTickets is a free data retrieval call binding the contract method 0x6262d54e.
//
// Solidity: function submittedTickets() constant returns(uint64[])
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) SubmittedTickets() ([]uint64, error) {
	return _KeepRandomBeaconOperator.Contract.SubmittedTickets(&_KeepRandomBeaconOperator.CallOpts)
}

// TicketSubmissionTimeout is a free data retrieval call binding the contract method 0xc98622fb.
//
// Solidity: function ticketSubmissionTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCaller) TicketSubmissionTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconOperator.contract.Call(opts, out, "ticketSubmissionTimeout")
	return *ret0, err
}

// TicketSubmissionTimeout is a free data retrieval call binding the contract method 0xc98622fb.
//
// Solidity: function ticketSubmissionTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) TicketSubmissionTimeout() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.TicketSubmissionTimeout(&_KeepRandomBeaconOperator.CallOpts)
}

// TicketSubmissionTimeout is a free data retrieval call binding the contract method 0xc98622fb.
//
// Solidity: function ticketSubmissionTimeout() constant returns(uint256)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorCallerSession) TicketSubmissionTimeout() (*big.Int, error) {
	return _KeepRandomBeaconOperator.Contract.TicketSubmissionTimeout(&_KeepRandomBeaconOperator.CallOpts)
}

// AddServiceContract is a paid mutator transaction binding the contract method 0x7760c6c7.
//
// Solidity: function addServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) AddServiceContract(opts *bind.TransactOpts, serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "addServiceContract", serviceContract)
}

// AddServiceContract is a paid mutator transaction binding the contract method 0x7760c6c7.
//
// Solidity: function addServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) AddServiceContract(serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.AddServiceContract(&_KeepRandomBeaconOperator.TransactOpts, serviceContract)
}

// AddServiceContract is a paid mutator transaction binding the contract method 0x7760c6c7.
//
// Solidity: function addServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) AddServiceContract(serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.AddServiceContract(&_KeepRandomBeaconOperator.TransactOpts, serviceContract)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xc96e71fb.
//
// Solidity: function createGroup(uint256 _newEntry, address submitter) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) CreateGroup(opts *bind.TransactOpts, _newEntry *big.Int, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "createGroup", _newEntry, submitter)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xc96e71fb.
//
// Solidity: function createGroup(uint256 _newEntry, address submitter) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) CreateGroup(_newEntry *big.Int, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.CreateGroup(&_KeepRandomBeaconOperator.TransactOpts, _newEntry, submitter)
}

// CreateGroup is a paid mutator transaction binding the contract method 0xc96e71fb.
//
// Solidity: function createGroup(uint256 _newEntry, address submitter) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) CreateGroup(_newEntry *big.Int, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.CreateGroup(&_KeepRandomBeaconOperator.TransactOpts, _newEntry, submitter)
}

// Genesis is a paid mutator transaction binding the contract method 0xa7f0b3de.
//
// Solidity: function genesis() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) Genesis(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "genesis")
}

// Genesis is a paid mutator transaction binding the contract method 0xa7f0b3de.
//
// Solidity: function genesis() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) Genesis() (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.Genesis(&_KeepRandomBeaconOperator.TransactOpts)
}

// Genesis is a paid mutator transaction binding the contract method 0xa7f0b3de.
//
// Solidity: function genesis() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) Genesis() (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.Genesis(&_KeepRandomBeaconOperator.TransactOpts)
}

// RelayEntry is a paid mutator transaction binding the contract method 0xac374f4b.
//
// Solidity: function relayEntry(bytes _groupSignature) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) RelayEntry(opts *bind.TransactOpts, _groupSignature []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "relayEntry", _groupSignature)
}

// RelayEntry is a paid mutator transaction binding the contract method 0xac374f4b.
//
// Solidity: function relayEntry(bytes _groupSignature) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) RelayEntry(_groupSignature []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.RelayEntry(&_KeepRandomBeaconOperator.TransactOpts, _groupSignature)
}

// RelayEntry is a paid mutator transaction binding the contract method 0xac374f4b.
//
// Solidity: function relayEntry(bytes _groupSignature) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) RelayEntry(_groupSignature []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.RelayEntry(&_KeepRandomBeaconOperator.TransactOpts, _groupSignature)
}

// RemoveServiceContract is a paid mutator transaction binding the contract method 0xb64ac8e4.
//
// Solidity: function removeServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) RemoveServiceContract(opts *bind.TransactOpts, serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "removeServiceContract", serviceContract)
}

// RemoveServiceContract is a paid mutator transaction binding the contract method 0xb64ac8e4.
//
// Solidity: function removeServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) RemoveServiceContract(serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.RemoveServiceContract(&_KeepRandomBeaconOperator.TransactOpts, serviceContract)
}

// RemoveServiceContract is a paid mutator transaction binding the contract method 0xb64ac8e4.
//
// Solidity: function removeServiceContract(address serviceContract) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) RemoveServiceContract(serviceContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.RemoveServiceContract(&_KeepRandomBeaconOperator.TransactOpts, serviceContract)
}

// ReportRelayEntryTimeout is a paid mutator transaction binding the contract method 0x1ed74070.
//
// Solidity: function reportRelayEntryTimeout() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) ReportRelayEntryTimeout(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "reportRelayEntryTimeout")
}

// ReportRelayEntryTimeout is a paid mutator transaction binding the contract method 0x1ed74070.
//
// Solidity: function reportRelayEntryTimeout() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) ReportRelayEntryTimeout() (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.ReportRelayEntryTimeout(&_KeepRandomBeaconOperator.TransactOpts)
}

// ReportRelayEntryTimeout is a paid mutator transaction binding the contract method 0x1ed74070.
//
// Solidity: function reportRelayEntryTimeout() returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) ReportRelayEntryTimeout() (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.ReportRelayEntryTimeout(&_KeepRandomBeaconOperator.TransactOpts)
}

// ReportUnauthorizedSigning is a paid mutator transaction binding the contract method 0xe581ff74.
//
// Solidity: function reportUnauthorizedSigning(uint256 groupIndex, bytes signedMsgSender) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) ReportUnauthorizedSigning(opts *bind.TransactOpts, groupIndex *big.Int, signedMsgSender []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "reportUnauthorizedSigning", groupIndex, signedMsgSender)
}

// ReportUnauthorizedSigning is a paid mutator transaction binding the contract method 0xe581ff74.
//
// Solidity: function reportUnauthorizedSigning(uint256 groupIndex, bytes signedMsgSender) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) ReportUnauthorizedSigning(groupIndex *big.Int, signedMsgSender []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.ReportUnauthorizedSigning(&_KeepRandomBeaconOperator.TransactOpts, groupIndex, signedMsgSender)
}

// ReportUnauthorizedSigning is a paid mutator transaction binding the contract method 0xe581ff74.
//
// Solidity: function reportUnauthorizedSigning(uint256 groupIndex, bytes signedMsgSender) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) ReportUnauthorizedSigning(groupIndex *big.Int, signedMsgSender []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.ReportUnauthorizedSigning(&_KeepRandomBeaconOperator.TransactOpts, groupIndex, signedMsgSender)
}

// Sign is a paid mutator transaction binding the contract method 0x9b3d270a.
//
// Solidity: function sign(uint256 requestId, bytes previousEntry) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) Sign(opts *bind.TransactOpts, requestId *big.Int, previousEntry []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "sign", requestId, previousEntry)
}

// Sign is a paid mutator transaction binding the contract method 0x9b3d270a.
//
// Solidity: function sign(uint256 requestId, bytes previousEntry) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) Sign(requestId *big.Int, previousEntry []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.Sign(&_KeepRandomBeaconOperator.TransactOpts, requestId, previousEntry)
}

// Sign is a paid mutator transaction binding the contract method 0x9b3d270a.
//
// Solidity: function sign(uint256 requestId, bytes previousEntry) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) Sign(requestId *big.Int, previousEntry []byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.Sign(&_KeepRandomBeaconOperator.TransactOpts, requestId, previousEntry)
}

// SubmitDkgResult is a paid mutator transaction binding the contract method 0x73f1daab.
//
// Solidity: function submitDkgResult(uint256 submitterMemberIndex, bytes groupPubKey, bytes misbehaved, bytes signatures, uint256[] signingMembersIndexes) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) SubmitDkgResult(opts *bind.TransactOpts, submitterMemberIndex *big.Int, groupPubKey []byte, misbehaved []byte, signatures []byte, signingMembersIndexes []*big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "submitDkgResult", submitterMemberIndex, groupPubKey, misbehaved, signatures, signingMembersIndexes)
}

// SubmitDkgResult is a paid mutator transaction binding the contract method 0x73f1daab.
//
// Solidity: function submitDkgResult(uint256 submitterMemberIndex, bytes groupPubKey, bytes misbehaved, bytes signatures, uint256[] signingMembersIndexes) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) SubmitDkgResult(submitterMemberIndex *big.Int, groupPubKey []byte, misbehaved []byte, signatures []byte, signingMembersIndexes []*big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.SubmitDkgResult(&_KeepRandomBeaconOperator.TransactOpts, submitterMemberIndex, groupPubKey, misbehaved, signatures, signingMembersIndexes)
}

// SubmitDkgResult is a paid mutator transaction binding the contract method 0x73f1daab.
//
// Solidity: function submitDkgResult(uint256 submitterMemberIndex, bytes groupPubKey, bytes misbehaved, bytes signatures, uint256[] signingMembersIndexes) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) SubmitDkgResult(submitterMemberIndex *big.Int, groupPubKey []byte, misbehaved []byte, signatures []byte, signingMembersIndexes []*big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.SubmitDkgResult(&_KeepRandomBeaconOperator.TransactOpts, submitterMemberIndex, groupPubKey, misbehaved, signatures, signingMembersIndexes)
}

// SubmitTicket is a paid mutator transaction binding the contract method 0x8a3a3da8.
//
// Solidity: function submitTicket(bytes32 ticket) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) SubmitTicket(opts *bind.TransactOpts, ticket [32]byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "submitTicket", ticket)
}

// SubmitTicket is a paid mutator transaction binding the contract method 0x8a3a3da8.
//
// Solidity: function submitTicket(bytes32 ticket) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) SubmitTicket(ticket [32]byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.SubmitTicket(&_KeepRandomBeaconOperator.TransactOpts, ticket)
}

// SubmitTicket is a paid mutator transaction binding the contract method 0x8a3a3da8.
//
// Solidity: function submitTicket(bytes32 ticket) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) SubmitTicket(ticket [32]byte) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.SubmitTicket(&_KeepRandomBeaconOperator.TransactOpts, ticket)
}

// WithdrawGroupMemberRewards is a paid mutator transaction binding the contract method 0x3926c28e.
//
// Solidity: function withdrawGroupMemberRewards(address operator, uint256 groupIndex) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactor) WithdrawGroupMemberRewards(opts *bind.TransactOpts, operator common.Address, groupIndex *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.contract.Transact(opts, "withdrawGroupMemberRewards", operator, groupIndex)
}

// WithdrawGroupMemberRewards is a paid mutator transaction binding the contract method 0x3926c28e.
//
// Solidity: function withdrawGroupMemberRewards(address operator, uint256 groupIndex) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorSession) WithdrawGroupMemberRewards(operator common.Address, groupIndex *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.WithdrawGroupMemberRewards(&_KeepRandomBeaconOperator.TransactOpts, operator, groupIndex)
}

// WithdrawGroupMemberRewards is a paid mutator transaction binding the contract method 0x3926c28e.
//
// Solidity: function withdrawGroupMemberRewards(address operator, uint256 groupIndex) returns()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorTransactorSession) WithdrawGroupMemberRewards(operator common.Address, groupIndex *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconOperator.Contract.WithdrawGroupMemberRewards(&_KeepRandomBeaconOperator.TransactOpts, operator, groupIndex)
}

// KeepRandomBeaconOperatorDkgResultSubmittedEventIterator is returned from FilterDkgResultSubmittedEvent and is used to iterate over the raw logs and unpacked data for DkgResultSubmittedEvent events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorDkgResultSubmittedEventIterator struct {
	Event *KeepRandomBeaconOperatorDkgResultSubmittedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorDkgResultSubmittedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorDkgResultSubmittedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorDkgResultSubmittedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorDkgResultSubmittedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorDkgResultSubmittedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorDkgResultSubmittedEvent represents a DkgResultSubmittedEvent event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorDkgResultSubmittedEvent struct {
	MemberIndex *big.Int
	GroupPubKey []byte
	Misbehaved  []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDkgResultSubmittedEvent is a free log retrieval operation binding the contract event 0xd1d71346ed0f1479c55b14e7c48b084207a7e1bcd9abe4f22d425ec23a518a1f.
//
// Solidity: event DkgResultSubmittedEvent(uint256 memberIndex, bytes groupPubKey, bytes misbehaved)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterDkgResultSubmittedEvent(opts *bind.FilterOpts) (*KeepRandomBeaconOperatorDkgResultSubmittedEventIterator, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "DkgResultSubmittedEvent")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorDkgResultSubmittedEventIterator{contract: _KeepRandomBeaconOperator.contract, event: "DkgResultSubmittedEvent", logs: logs, sub: sub}, nil
}

// WatchDkgResultSubmittedEvent is a free log subscription operation binding the contract event 0xd1d71346ed0f1479c55b14e7c48b084207a7e1bcd9abe4f22d425ec23a518a1f.
//
// Solidity: event DkgResultSubmittedEvent(uint256 memberIndex, bytes groupPubKey, bytes misbehaved)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchDkgResultSubmittedEvent(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorDkgResultSubmittedEvent) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "DkgResultSubmittedEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorDkgResultSubmittedEvent)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "DkgResultSubmittedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDkgResultSubmittedEvent is a log parse operation binding the contract event 0xd1d71346ed0f1479c55b14e7c48b084207a7e1bcd9abe4f22d425ec23a518a1f.
//
// Solidity: event DkgResultSubmittedEvent(uint256 memberIndex, bytes groupPubKey, bytes misbehaved)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseDkgResultSubmittedEvent(log types.Log) (*KeepRandomBeaconOperatorDkgResultSubmittedEvent, error) {
	event := new(KeepRandomBeaconOperatorDkgResultSubmittedEvent)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "DkgResultSubmittedEvent", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator is returned from FilterGroupMemberRewardsWithdrawn and is used to iterate over the raw logs and unpacked data for GroupMemberRewardsWithdrawn events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator struct {
	Event *KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn represents a GroupMemberRewardsWithdrawn event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn struct {
	Beneficiary common.Address
	Operator    common.Address
	Amount      *big.Int
	GroupIndex  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGroupMemberRewardsWithdrawn is a free log retrieval operation binding the contract event 0xd2d1d8bb9db82c3480418ddcddf25a021102ad139edec2a62b274595d408a88d.
//
// Solidity: event GroupMemberRewardsWithdrawn(address indexed beneficiary, address operator, uint256 amount, uint256 groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterGroupMemberRewardsWithdrawn(opts *bind.FilterOpts, beneficiary []common.Address) (*KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "GroupMemberRewardsWithdrawn", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorGroupMemberRewardsWithdrawnIterator{contract: _KeepRandomBeaconOperator.contract, event: "GroupMemberRewardsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchGroupMemberRewardsWithdrawn is a free log subscription operation binding the contract event 0xd2d1d8bb9db82c3480418ddcddf25a021102ad139edec2a62b274595d408a88d.
//
// Solidity: event GroupMemberRewardsWithdrawn(address indexed beneficiary, address operator, uint256 amount, uint256 groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchGroupMemberRewardsWithdrawn(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "GroupMemberRewardsWithdrawn", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "GroupMemberRewardsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGroupMemberRewardsWithdrawn is a log parse operation binding the contract event 0xd2d1d8bb9db82c3480418ddcddf25a021102ad139edec2a62b274595d408a88d.
//
// Solidity: event GroupMemberRewardsWithdrawn(address indexed beneficiary, address operator, uint256 amount, uint256 groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseGroupMemberRewardsWithdrawn(log types.Log) (*KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn, error) {
	event := new(KeepRandomBeaconOperatorGroupMemberRewardsWithdrawn)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "GroupMemberRewardsWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorGroupSelectionStartedIterator is returned from FilterGroupSelectionStarted and is used to iterate over the raw logs and unpacked data for GroupSelectionStarted events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorGroupSelectionStartedIterator struct {
	Event *KeepRandomBeaconOperatorGroupSelectionStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorGroupSelectionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorGroupSelectionStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorGroupSelectionStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorGroupSelectionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorGroupSelectionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorGroupSelectionStarted represents a GroupSelectionStarted event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorGroupSelectionStarted struct {
	NewEntry *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterGroupSelectionStarted is a free log retrieval operation binding the contract event 0x0769b89b6dbd96af3cdebccc7b68ce1e4ae748abc3e6b19a73b8b58460c57a94.
//
// Solidity: event GroupSelectionStarted(uint256 newEntry)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterGroupSelectionStarted(opts *bind.FilterOpts) (*KeepRandomBeaconOperatorGroupSelectionStartedIterator, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "GroupSelectionStarted")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorGroupSelectionStartedIterator{contract: _KeepRandomBeaconOperator.contract, event: "GroupSelectionStarted", logs: logs, sub: sub}, nil
}

// WatchGroupSelectionStarted is a free log subscription operation binding the contract event 0x0769b89b6dbd96af3cdebccc7b68ce1e4ae748abc3e6b19a73b8b58460c57a94.
//
// Solidity: event GroupSelectionStarted(uint256 newEntry)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchGroupSelectionStarted(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorGroupSelectionStarted) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "GroupSelectionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorGroupSelectionStarted)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "GroupSelectionStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGroupSelectionStarted is a log parse operation binding the contract event 0x0769b89b6dbd96af3cdebccc7b68ce1e4ae748abc3e6b19a73b8b58460c57a94.
//
// Solidity: event GroupSelectionStarted(uint256 newEntry)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseGroupSelectionStarted(log types.Log) (*KeepRandomBeaconOperatorGroupSelectionStarted, error) {
	event := new(KeepRandomBeaconOperatorGroupSelectionStarted)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "GroupSelectionStarted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorOnGroupRegisteredIterator is returned from FilterOnGroupRegistered and is used to iterate over the raw logs and unpacked data for OnGroupRegistered events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorOnGroupRegisteredIterator struct {
	Event *KeepRandomBeaconOperatorOnGroupRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorOnGroupRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorOnGroupRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorOnGroupRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorOnGroupRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorOnGroupRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorOnGroupRegistered represents a OnGroupRegistered event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorOnGroupRegistered struct {
	GroupPubKey []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOnGroupRegistered is a free log retrieval operation binding the contract event 0x3a32cfbcd60aee14d75b837115f97fc141e9f9fa0d1fcd310b6306abbf329c15.
//
// Solidity: event OnGroupRegistered(bytes groupPubKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterOnGroupRegistered(opts *bind.FilterOpts) (*KeepRandomBeaconOperatorOnGroupRegisteredIterator, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "OnGroupRegistered")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorOnGroupRegisteredIterator{contract: _KeepRandomBeaconOperator.contract, event: "OnGroupRegistered", logs: logs, sub: sub}, nil
}

// WatchOnGroupRegistered is a free log subscription operation binding the contract event 0x3a32cfbcd60aee14d75b837115f97fc141e9f9fa0d1fcd310b6306abbf329c15.
//
// Solidity: event OnGroupRegistered(bytes groupPubKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchOnGroupRegistered(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorOnGroupRegistered) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "OnGroupRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorOnGroupRegistered)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "OnGroupRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOnGroupRegistered is a log parse operation binding the contract event 0x3a32cfbcd60aee14d75b837115f97fc141e9f9fa0d1fcd310b6306abbf329c15.
//
// Solidity: event OnGroupRegistered(bytes groupPubKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseOnGroupRegistered(log types.Log) (*KeepRandomBeaconOperatorOnGroupRegistered, error) {
	event := new(KeepRandomBeaconOperatorOnGroupRegistered)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "OnGroupRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorRelayEntryRequestedIterator is returned from FilterRelayEntryRequested and is used to iterate over the raw logs and unpacked data for RelayEntryRequested events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntryRequestedIterator struct {
	Event *KeepRandomBeaconOperatorRelayEntryRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorRelayEntryRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorRelayEntryRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorRelayEntryRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorRelayEntryRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorRelayEntryRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorRelayEntryRequested represents a RelayEntryRequested event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntryRequested struct {
	PreviousEntry  []byte
	GroupPublicKey []byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryRequested is a free log retrieval operation binding the contract event 0xf3a8bf09e4f9146a48f9b91226985ac8d83d971beb4fc9ffdc569790e85a97e4.
//
// Solidity: event RelayEntryRequested(bytes previousEntry, bytes groupPublicKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterRelayEntryRequested(opts *bind.FilterOpts) (*KeepRandomBeaconOperatorRelayEntryRequestedIterator, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorRelayEntryRequestedIterator{contract: _KeepRandomBeaconOperator.contract, event: "RelayEntryRequested", logs: logs, sub: sub}, nil
}

// WatchRelayEntryRequested is a free log subscription operation binding the contract event 0xf3a8bf09e4f9146a48f9b91226985ac8d83d971beb4fc9ffdc569790e85a97e4.
//
// Solidity: event RelayEntryRequested(bytes previousEntry, bytes groupPublicKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchRelayEntryRequested(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorRelayEntryRequested) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorRelayEntryRequested)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntryRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayEntryRequested is a log parse operation binding the contract event 0xf3a8bf09e4f9146a48f9b91226985ac8d83d971beb4fc9ffdc569790e85a97e4.
//
// Solidity: event RelayEntryRequested(bytes previousEntry, bytes groupPublicKey)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseRelayEntryRequested(log types.Log) (*KeepRandomBeaconOperatorRelayEntryRequested, error) {
	event := new(KeepRandomBeaconOperatorRelayEntryRequested)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntryRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorRelayEntrySubmittedIterator is returned from FilterRelayEntrySubmitted and is used to iterate over the raw logs and unpacked data for RelayEntrySubmitted events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntrySubmittedIterator struct {
	Event *KeepRandomBeaconOperatorRelayEntrySubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorRelayEntrySubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorRelayEntrySubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorRelayEntrySubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorRelayEntrySubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorRelayEntrySubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorRelayEntrySubmitted represents a RelayEntrySubmitted event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntrySubmitted struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRelayEntrySubmitted is a free log retrieval operation binding the contract event 0x8711cae111460cf9bde0d890f0dc09abcb8851e39bf020f406e53e86394cdbd7.
//
// Solidity: event RelayEntrySubmitted()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterRelayEntrySubmitted(opts *bind.FilterOpts) (*KeepRandomBeaconOperatorRelayEntrySubmittedIterator, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "RelayEntrySubmitted")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorRelayEntrySubmittedIterator{contract: _KeepRandomBeaconOperator.contract, event: "RelayEntrySubmitted", logs: logs, sub: sub}, nil
}

// WatchRelayEntrySubmitted is a free log subscription operation binding the contract event 0x8711cae111460cf9bde0d890f0dc09abcb8851e39bf020f406e53e86394cdbd7.
//
// Solidity: event RelayEntrySubmitted()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchRelayEntrySubmitted(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorRelayEntrySubmitted) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "RelayEntrySubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorRelayEntrySubmitted)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntrySubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayEntrySubmitted is a log parse operation binding the contract event 0x8711cae111460cf9bde0d890f0dc09abcb8851e39bf020f406e53e86394cdbd7.
//
// Solidity: event RelayEntrySubmitted()
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseRelayEntrySubmitted(log types.Log) (*KeepRandomBeaconOperatorRelayEntrySubmitted, error) {
	event := new(KeepRandomBeaconOperatorRelayEntrySubmitted)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntrySubmitted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator is returned from FilterRelayEntryTimeoutReported and is used to iterate over the raw logs and unpacked data for RelayEntryTimeoutReported events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator struct {
	Event *KeepRandomBeaconOperatorRelayEntryTimeoutReported // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorRelayEntryTimeoutReported)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorRelayEntryTimeoutReported)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorRelayEntryTimeoutReported represents a RelayEntryTimeoutReported event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorRelayEntryTimeoutReported struct {
	GroupIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryTimeoutReported is a free log retrieval operation binding the contract event 0x6675fe3ae219641aa4ec9e58867bd7af88bf03caf819d8858f3ddf4cc635eed2.
//
// Solidity: event RelayEntryTimeoutReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterRelayEntryTimeoutReported(opts *bind.FilterOpts, groupIndex []*big.Int) (*KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator, error) {

	var groupIndexRule []interface{}
	for _, groupIndexItem := range groupIndex {
		groupIndexRule = append(groupIndexRule, groupIndexItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "RelayEntryTimeoutReported", groupIndexRule)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorRelayEntryTimeoutReportedIterator{contract: _KeepRandomBeaconOperator.contract, event: "RelayEntryTimeoutReported", logs: logs, sub: sub}, nil
}

// WatchRelayEntryTimeoutReported is a free log subscription operation binding the contract event 0x6675fe3ae219641aa4ec9e58867bd7af88bf03caf819d8858f3ddf4cc635eed2.
//
// Solidity: event RelayEntryTimeoutReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchRelayEntryTimeoutReported(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorRelayEntryTimeoutReported, groupIndex []*big.Int) (event.Subscription, error) {

	var groupIndexRule []interface{}
	for _, groupIndexItem := range groupIndex {
		groupIndexRule = append(groupIndexRule, groupIndexItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "RelayEntryTimeoutReported", groupIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorRelayEntryTimeoutReported)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntryTimeoutReported", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRelayEntryTimeoutReported is a log parse operation binding the contract event 0x6675fe3ae219641aa4ec9e58867bd7af88bf03caf819d8858f3ddf4cc635eed2.
//
// Solidity: event RelayEntryTimeoutReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseRelayEntryTimeoutReported(log types.Log) (*KeepRandomBeaconOperatorRelayEntryTimeoutReported, error) {
	event := new(KeepRandomBeaconOperatorRelayEntryTimeoutReported)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "RelayEntryTimeoutReported", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator is returned from FilterUnauthorizedSigningReported and is used to iterate over the raw logs and unpacked data for UnauthorizedSigningReported events raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator struct {
	Event *KeepRandomBeaconOperatorUnauthorizedSigningReported // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconOperatorUnauthorizedSigningReported)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(KeepRandomBeaconOperatorUnauthorizedSigningReported)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconOperatorUnauthorizedSigningReported represents a UnauthorizedSigningReported event raised by the KeepRandomBeaconOperator contract.
type KeepRandomBeaconOperatorUnauthorizedSigningReported struct {
	GroupIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUnauthorizedSigningReported is a free log retrieval operation binding the contract event 0x6124f28ae7240a98a8ad3410bcd1f3bb0a113fc9834d5bd16426f9e1bd698fde.
//
// Solidity: event UnauthorizedSigningReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) FilterUnauthorizedSigningReported(opts *bind.FilterOpts, groupIndex []*big.Int) (*KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator, error) {

	var groupIndexRule []interface{}
	for _, groupIndexItem := range groupIndex {
		groupIndexRule = append(groupIndexRule, groupIndexItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.FilterLogs(opts, "UnauthorizedSigningReported", groupIndexRule)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconOperatorUnauthorizedSigningReportedIterator{contract: _KeepRandomBeaconOperator.contract, event: "UnauthorizedSigningReported", logs: logs, sub: sub}, nil
}

// WatchUnauthorizedSigningReported is a free log subscription operation binding the contract event 0x6124f28ae7240a98a8ad3410bcd1f3bb0a113fc9834d5bd16426f9e1bd698fde.
//
// Solidity: event UnauthorizedSigningReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) WatchUnauthorizedSigningReported(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconOperatorUnauthorizedSigningReported, groupIndex []*big.Int) (event.Subscription, error) {

	var groupIndexRule []interface{}
	for _, groupIndexItem := range groupIndex {
		groupIndexRule = append(groupIndexRule, groupIndexItem)
	}

	logs, sub, err := _KeepRandomBeaconOperator.contract.WatchLogs(opts, "UnauthorizedSigningReported", groupIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconOperatorUnauthorizedSigningReported)
				if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "UnauthorizedSigningReported", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnauthorizedSigningReported is a log parse operation binding the contract event 0x6124f28ae7240a98a8ad3410bcd1f3bb0a113fc9834d5bd16426f9e1bd698fde.
//
// Solidity: event UnauthorizedSigningReported(uint256 indexed groupIndex)
func (_KeepRandomBeaconOperator *KeepRandomBeaconOperatorFilterer) ParseUnauthorizedSigningReported(log types.Log) (*KeepRandomBeaconOperatorUnauthorizedSigningReported, error) {
	event := new(KeepRandomBeaconOperatorUnauthorizedSigningReported)
	if err := _KeepRandomBeaconOperator.contract.UnpackLog(event, "UnauthorizedSigningReported", log); err != nil {
		return nil, err
	}
	return event, nil
}
