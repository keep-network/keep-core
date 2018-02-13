// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package GenRequestID

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// GenRequestIDABI is the input ABI used to generate the binding from.
const GenRequestIDABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"RequestIDSequence\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"GenerateNextRequestID\",\"outputs\":[{\"name\":\"RequestID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"}]"

// GenRequestID is an auto generated Go binding around an Ethereum contract.
type GenRequestID struct {
	GenRequestIDCaller     // Read-only binding to the contract
	GenRequestIDTransactor // Write-only binding to the contract
	GenRequestIDFilterer   // Log filterer for contract events
}

// GenRequestIDCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenRequestIDCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenRequestIDTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenRequestIDTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenRequestIDFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenRequestIDFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenRequestIDSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenRequestIDSession struct {
	Contract     *GenRequestID     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GenRequestIDCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenRequestIDCallerSession struct {
	Contract *GenRequestIDCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GenRequestIDTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenRequestIDTransactorSession struct {
	Contract     *GenRequestIDTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GenRequestIDRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenRequestIDRaw struct {
	Contract *GenRequestID // Generic contract binding to access the raw methods on
}

// GenRequestIDCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenRequestIDCallerRaw struct {
	Contract *GenRequestIDCaller // Generic read-only contract binding to access the raw methods on
}

// GenRequestIDTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenRequestIDTransactorRaw struct {
	Contract *GenRequestIDTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenRequestID creates a new instance of GenRequestID, bound to a specific deployed contract.
func NewGenRequestID(address common.Address, backend bind.ContractBackend) (*GenRequestID, error) {
	contract, err := bindGenRequestID(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenRequestID{GenRequestIDCaller: GenRequestIDCaller{contract: contract}, GenRequestIDTransactor: GenRequestIDTransactor{contract: contract}, GenRequestIDFilterer: GenRequestIDFilterer{contract: contract}}, nil
}

// NewGenRequestIDCaller creates a new read-only instance of GenRequestID, bound to a specific deployed contract.
func NewGenRequestIDCaller(address common.Address, caller bind.ContractCaller) (*GenRequestIDCaller, error) {
	contract, err := bindGenRequestID(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenRequestIDCaller{contract: contract}, nil
}

// NewGenRequestIDTransactor creates a new write-only instance of GenRequestID, bound to a specific deployed contract.
func NewGenRequestIDTransactor(address common.Address, transactor bind.ContractTransactor) (*GenRequestIDTransactor, error) {
	contract, err := bindGenRequestID(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenRequestIDTransactor{contract: contract}, nil
}

// NewGenRequestIDFilterer creates a new log filterer instance of GenRequestID, bound to a specific deployed contract.
func NewGenRequestIDFilterer(address common.Address, filterer bind.ContractFilterer) (*GenRequestIDFilterer, error) {
	contract, err := bindGenRequestID(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenRequestIDFilterer{contract: contract}, nil
}

// bindGenRequestID binds a generic wrapper to an already deployed contract.
func bindGenRequestID(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GenRequestIDABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenRequestID *GenRequestIDRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GenRequestID.Contract.GenRequestIDCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenRequestID *GenRequestIDRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenRequestID.Contract.GenRequestIDTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenRequestID *GenRequestIDRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenRequestID.Contract.GenRequestIDTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenRequestID *GenRequestIDCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GenRequestID.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenRequestID *GenRequestIDTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenRequestID.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenRequestID *GenRequestIDTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenRequestID.Contract.contract.Transact(opts, method, params...)
}

// RequestIDSequence is a free data retrieval call binding the contract method 0x2d7c9e3e.
//
// Solidity: function RequestIDSequence() constant returns(uint256)
func (_GenRequestID *GenRequestIDCaller) RequestIDSequence(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GenRequestID.contract.Call(opts, out, "RequestIDSequence")
	return *ret0, err
}

// RequestIDSequence is a free data retrieval call binding the contract method 0x2d7c9e3e.
//
// Solidity: function RequestIDSequence() constant returns(uint256)
func (_GenRequestID *GenRequestIDSession) RequestIDSequence() (*big.Int, error) {
	return _GenRequestID.Contract.RequestIDSequence(&_GenRequestID.CallOpts)
}

// RequestIDSequence is a free data retrieval call binding the contract method 0x2d7c9e3e.
//
// Solidity: function RequestIDSequence() constant returns(uint256)
func (_GenRequestID *GenRequestIDCallerSession) RequestIDSequence() (*big.Int, error) {
	return _GenRequestID.Contract.RequestIDSequence(&_GenRequestID.CallOpts)
}

// GenerateNextRequestID is a paid mutator transaction binding the contract method 0x68e706c3.
//
// Solidity: function GenerateNextRequestID() returns(RequestID uint256)
func (_GenRequestID *GenRequestIDTransactor) GenerateNextRequestID(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenRequestID.contract.Transact(opts, "GenerateNextRequestID")
}

// GenerateNextRequestID is a paid mutator transaction binding the contract method 0x68e706c3.
//
// Solidity: function GenerateNextRequestID() returns(RequestID uint256)
func (_GenRequestID *GenRequestIDSession) GenerateNextRequestID() (*types.Transaction, error) {
	return _GenRequestID.Contract.GenerateNextRequestID(&_GenRequestID.TransactOpts)
}

// GenerateNextRequestID is a paid mutator transaction binding the contract method 0x68e706c3.
//
// Solidity: function GenerateNextRequestID() returns(RequestID uint256)
func (_GenRequestID *GenRequestIDTransactorSession) GenerateNextRequestID() (*types.Transaction, error) {
	return _GenRequestID.Contract.GenerateNextRequestID(&_GenRequestID.TransactOpts)
}
