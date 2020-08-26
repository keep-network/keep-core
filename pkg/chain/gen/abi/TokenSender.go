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

// TokenSenderABI is the input ABI used to generate the binding from.
const TokenSenderABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"approveAndCall\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenSender is an auto generated Go binding around an Ethereum contract.
type TokenSender struct {
	TokenSenderCaller     // Read-only binding to the contract
	TokenSenderTransactor // Write-only binding to the contract
	TokenSenderFilterer   // Log filterer for contract events
}

// TokenSenderCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenSenderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSenderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenSenderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSenderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenSenderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSenderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSenderSession struct {
	Contract     *TokenSender      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenSenderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenSenderCallerSession struct {
	Contract *TokenSenderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TokenSenderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenSenderTransactorSession struct {
	Contract     *TokenSenderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TokenSenderRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenSenderRaw struct {
	Contract *TokenSender // Generic contract binding to access the raw methods on
}

// TokenSenderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenSenderCallerRaw struct {
	Contract *TokenSenderCaller // Generic read-only contract binding to access the raw methods on
}

// TokenSenderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenSenderTransactorRaw struct {
	Contract *TokenSenderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenSender creates a new instance of TokenSender, bound to a specific deployed contract.
func NewTokenSender(address common.Address, backend bind.ContractBackend) (*TokenSender, error) {
	contract, err := bindTokenSender(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenSender{TokenSenderCaller: TokenSenderCaller{contract: contract}, TokenSenderTransactor: TokenSenderTransactor{contract: contract}, TokenSenderFilterer: TokenSenderFilterer{contract: contract}}, nil
}

// NewTokenSenderCaller creates a new read-only instance of TokenSender, bound to a specific deployed contract.
func NewTokenSenderCaller(address common.Address, caller bind.ContractCaller) (*TokenSenderCaller, error) {
	contract, err := bindTokenSender(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenSenderCaller{contract: contract}, nil
}

// NewTokenSenderTransactor creates a new write-only instance of TokenSender, bound to a specific deployed contract.
func NewTokenSenderTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenSenderTransactor, error) {
	contract, err := bindTokenSender(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenSenderTransactor{contract: contract}, nil
}

// NewTokenSenderFilterer creates a new log filterer instance of TokenSender, bound to a specific deployed contract.
func NewTokenSenderFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenSenderFilterer, error) {
	contract, err := bindTokenSender(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenSenderFilterer{contract: contract}, nil
}

// bindTokenSender binds a generic wrapper to an already deployed contract.
func bindTokenSender(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenSenderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenSender *TokenSenderRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenSender.Contract.TokenSenderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenSender *TokenSenderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenSender.Contract.TokenSenderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenSender *TokenSenderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenSender.Contract.TokenSenderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenSender *TokenSenderCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenSender.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenSender *TokenSenderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenSender.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenSender *TokenSenderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenSender.Contract.contract.Transact(opts, method, params...)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address _spender, uint256 _value, bytes _extraData) returns()
func (_TokenSender *TokenSenderTransactor) ApproveAndCall(opts *bind.TransactOpts, _spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenSender.contract.Transact(opts, "approveAndCall", _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address _spender, uint256 _value, bytes _extraData) returns()
func (_TokenSender *TokenSenderSession) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenSender.Contract.ApproveAndCall(&_TokenSender.TransactOpts, _spender, _value, _extraData)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address _spender, uint256 _value, bytes _extraData) returns()
func (_TokenSender *TokenSenderTransactorSession) ApproveAndCall(_spender common.Address, _value *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenSender.Contract.ApproveAndCall(&_TokenSender.TransactOpts, _spender, _value, _extraData)
}
