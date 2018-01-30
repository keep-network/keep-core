// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"www.2c-why.com/job/ethrpc/ethABI/bind"
)

// TokenTypeABI is the input ABI used to generate the binding from.
const TokenTypeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"approveAndCall\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"spentAllowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"decimalUnits\",\"type\":\"uint8\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// TokenType is an auto generated Go binding around an Ethereum contract.
type TokenType struct {
	TokenTypeCaller     // Read-only binding to the contract
	TokenTypeTransactor // Write-only binding to the contract
}

// TokenTypeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenTypeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTypeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTypeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTypeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenTypeSession struct {
	Contract     *TokenType        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenTypeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenTypeCallerSession struct {
	Contract *TokenTypeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TokenTypeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTypeTransactorSession struct {
	Contract     *TokenTypeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TokenTypeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenTypeRaw struct {
	Contract *TokenType // Generic contract binding to access the raw methods on
}

// TokenTypeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenTypeCallerRaw struct {
	Contract *TokenTypeCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTypeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTypeTransactorRaw struct {
	Contract *TokenTypeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenType creates a new instance of TokenType, bound to a specific deployed contract.
func NewTokenType(address common.Address, backend bind.ContractBackend) (*TokenType, error) {
	contract, err := bindTokenType(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenType{TokenTypeCaller: TokenTypeCaller{contract: contract}, TokenTypeTransactor: TokenTypeTransactor{contract: contract}}, nil
}

// NewTokenTypeCaller creates a new read-only instance of TokenType, bound to a specific deployed contract.
func NewTokenTypeCaller(address common.Address, caller bind.ContractCaller) (*TokenTypeCaller, error) {
	contract, err := bindTokenType(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTypeCaller{contract: contract}, nil
}

// NewTokenTypeTransactor creates a new write-only instance of TokenType, bound to a specific deployed contract.
func NewTokenTypeTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTypeTransactor, error) {
	contract, err := bindTokenType(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &TokenTypeTransactor{contract: contract}, nil
}

// bindTokenType binds a generic wrapper to an already deployed contract.
func bindTokenType(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenTypeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenType *TokenTypeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenType.Contract.TokenTypeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenType *TokenTypeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenType.Contract.TokenTypeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenType *TokenTypeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenType.Contract.TokenTypeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenType *TokenTypeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenType.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenType *TokenTypeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenType.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenType *TokenTypeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenType.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "allowance", arg0, arg1)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TokenType.Contract.Allowance(&_TokenType.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TokenType.Contract.Allowance(&_TokenType.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_TokenType *TokenTypeCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "balanceOf", arg0)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_TokenType *TokenTypeSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _TokenType.Contract.BalanceOf(&_TokenType.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf( address) constant returns(uint256)
func (_TokenType *TokenTypeCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _TokenType.Contract.BalanceOf(&_TokenType.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TokenType *TokenTypeCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TokenType *TokenTypeSession) Decimals() (uint8, error) {
	return _TokenType.Contract.Decimals(&_TokenType.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TokenType *TokenTypeCallerSession) Decimals() (uint8, error) {
	return _TokenType.Contract.Decimals(&_TokenType.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TokenType *TokenTypeCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TokenType *TokenTypeSession) Name() (string, error) {
	return _TokenType.Contract.Name(&_TokenType.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TokenType *TokenTypeCallerSession) Name() (string, error) {
	return _TokenType.Contract.Name(&_TokenType.CallOpts)
}

// SpentAllowance is a free data retrieval call binding the contract method 0xdc3080f2.
//
// Solidity: function spentAllowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeCaller) SpentAllowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "spentAllowance", arg0, arg1)
	return *ret0, err
}

// SpentAllowance is a free data retrieval call binding the contract method 0xdc3080f2.
//
// Solidity: function spentAllowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeSession) SpentAllowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TokenType.Contract.SpentAllowance(&_TokenType.CallOpts, arg0, arg1)
}

// SpentAllowance is a free data retrieval call binding the contract method 0xdc3080f2.
//
// Solidity: function spentAllowance( address,  address) constant returns(uint256)
func (_TokenType *TokenTypeCallerSession) SpentAllowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _TokenType.Contract.SpentAllowance(&_TokenType.CallOpts, arg0, arg1)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TokenType *TokenTypeCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TokenType.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TokenType *TokenTypeSession) Symbol() (string, error) {
	return _TokenType.Contract.Symbol(&_TokenType.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TokenType *TokenTypeCallerSession) Symbol() (string, error) {
	return _TokenType.Contract.Symbol(&_TokenType.CallOpts)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(success bool)
func (_TokenType *TokenTypeTransactor) ApproveAndCall(opts *bind.TransactOpts, _spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenType.contract.Transact(opts, "approveAndCall", _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(success bool)
func (_TokenType *TokenTypeSession) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenType.Contract.ApproveAndCall(&_TokenType.TransactOpts, _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(_spender address, _value uint256, _extraData bytes) returns(success bool)
func (_TokenType *TokenTypeTransactorSession) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenType.Contract.ApproveAndCall(&_TokenType.TransactOpts, _spender, _value, _extraData)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns()
func (_TokenType *TokenTypeTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns()
func (_TokenType *TokenTypeSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.Contract.Transfer(&_TokenType.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns()
func (_TokenType *TokenTypeTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.Contract.Transfer(&_TokenType.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_TokenType *TokenTypeTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_TokenType *TokenTypeSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.Contract.TransferFrom(&_TokenType.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(_from address, _to address, _value uint256) returns(success bool)
func (_TokenType *TokenTypeTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _TokenType.Contract.TransferFrom(&_TokenType.TransactOpts, _from, _to, _value)
}
