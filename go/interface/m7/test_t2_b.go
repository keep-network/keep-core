// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TokenRecipientTypeABI is the input ABI used to generate the binding from.
const TokenRecipientTypeABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenRecipientType is an auto generated Go binding around an Ethereum contract.
type TokenRecipientType struct {
	TokenRecipientTypeCaller     // Read-only binding to the contract
	TokenRecipientTypeTransactor // Write-only binding to the contract
	TokenRecipientTypeFilterer   // Log filterer for contract events
}

// TokenRecipientTypeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenRecipientTypeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRecipientTypeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenRecipientTypeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRecipientTypeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenRecipientTypeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRecipientTypeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenRecipientTypeSession struct {
	Contract     *TokenRecipientType // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenRecipientTypeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenRecipientTypeCallerSession struct {
	Contract *TokenRecipientTypeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// TokenRecipientTypeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenRecipientTypeTransactorSession struct {
	Contract     *TokenRecipientTypeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// TokenRecipientTypeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRecipientTypeRaw struct {
	Contract *TokenRecipientType // Generic contract binding to access the raw methods on
}

// TokenRecipientTypeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenRecipientTypeCallerRaw struct {
	Contract *TokenRecipientTypeCaller // Generic read-only contract binding to access the raw methods on
}

// TokenRecipientTypeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenRecipientTypeTransactorRaw struct {
	Contract *TokenRecipientTypeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenRecipientType creates a new instance of TokenRecipientType, bound to a specific deployed contract.
func NewTokenRecipientType(address common.Address, backend bind.ContractBackend) (*TokenRecipientType, error) {
	contract, err := bindTokenRecipientType(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenRecipientType{TokenRecipientTypeCaller: TokenRecipientTypeCaller{contract: contract}, TokenRecipientTypeTransactor: TokenRecipientTypeTransactor{contract: contract}, TokenRecipientTypeFilterer: TokenRecipientTypeFilterer{contract: contract}}, nil
}

// NewTokenRecipientTypeCaller creates a new read-only instance of TokenRecipientType, bound to a specific deployed contract.
func NewTokenRecipientTypeCaller(address common.Address, caller bind.ContractCaller) (*TokenRecipientTypeCaller, error) {
	contract, err := bindTokenRecipientType(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenRecipientTypeCaller{contract: contract}, nil
}

// NewTokenRecipientTypeTransactor creates a new write-only instance of TokenRecipientType, bound to a specific deployed contract.
func NewTokenRecipientTypeTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenRecipientTypeTransactor, error) {
	contract, err := bindTokenRecipientType(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenRecipientTypeTransactor{contract: contract}, nil
}

// NewTokenRecipientTypeFilterer creates a new log filterer instance of TokenRecipientType, bound to a specific deployed contract.
func NewTokenRecipientTypeFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenRecipientTypeFilterer, error) {
	contract, err := bindTokenRecipientType(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenRecipientTypeFilterer{contract: contract}, nil
}

// bindTokenRecipientType binds a generic wrapper to an already deployed contract.
func bindTokenRecipientType(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenRecipientTypeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenRecipientType *TokenRecipientTypeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenRecipientType.Contract.TokenRecipientTypeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenRecipientType *TokenRecipientTypeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.TokenRecipientTypeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenRecipientType *TokenRecipientTypeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.TokenRecipientTypeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenRecipientType *TokenRecipientTypeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenRecipientType.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenRecipientType *TokenRecipientTypeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenRecipientType *TokenRecipientTypeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.contract.Transact(opts, method, params...)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(_from address, _value uint256, _token address, _extraData bytes) returns()
func (_TokenRecipientType *TokenRecipientTypeTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenRecipientType.contract.Transact(opts, "receiveApproval", _from, _value, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(_from address, _value uint256, _token address, _extraData bytes) returns()
func (_TokenRecipientType *TokenRecipientTypeSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.ReceiveApproval(&_TokenRecipientType.TransactOpts, _from, _value, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(_from address, _value uint256, _token address, _extraData bytes) returns()
func (_TokenRecipientType *TokenRecipientTypeTransactorSession) ReceiveApproval(_from common.Address, _value *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenRecipientType.Contract.ReceiveApproval(&_TokenRecipientType.TransactOpts, _from, _value, _token, _extraData)
}
