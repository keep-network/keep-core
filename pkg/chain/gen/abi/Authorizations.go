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

// AuthorizationsABI is the input ABI used to generate the binding from.
const AuthorizationsABI = "[{\"inputs\":[{\"internalType\":\"contractKeepRegistry\",\"name\":\"_registry\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"authorizeOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"authorizerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatedAuthoritySource\",\"type\":\"address\"}],\"name\":\"claimDelegatedAuthority\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"getAuthoritySource\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"isApprovedOperatorContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"isAuthorizedForOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Authorizations is an auto generated Go binding around an Ethereum contract.
type Authorizations struct {
	AuthorizationsCaller     // Read-only binding to the contract
	AuthorizationsTransactor // Write-only binding to the contract
	AuthorizationsFilterer   // Log filterer for contract events
}

// AuthorizationsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AuthorizationsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthorizationsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AuthorizationsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthorizationsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AuthorizationsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuthorizationsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AuthorizationsSession struct {
	Contract     *Authorizations   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuthorizationsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AuthorizationsCallerSession struct {
	Contract *AuthorizationsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AuthorizationsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AuthorizationsTransactorSession struct {
	Contract     *AuthorizationsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AuthorizationsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AuthorizationsRaw struct {
	Contract *Authorizations // Generic contract binding to access the raw methods on
}

// AuthorizationsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AuthorizationsCallerRaw struct {
	Contract *AuthorizationsCaller // Generic read-only contract binding to access the raw methods on
}

// AuthorizationsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AuthorizationsTransactorRaw struct {
	Contract *AuthorizationsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAuthorizations creates a new instance of Authorizations, bound to a specific deployed contract.
func NewAuthorizations(address common.Address, backend bind.ContractBackend) (*Authorizations, error) {
	contract, err := bindAuthorizations(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Authorizations{AuthorizationsCaller: AuthorizationsCaller{contract: contract}, AuthorizationsTransactor: AuthorizationsTransactor{contract: contract}, AuthorizationsFilterer: AuthorizationsFilterer{contract: contract}}, nil
}

// NewAuthorizationsCaller creates a new read-only instance of Authorizations, bound to a specific deployed contract.
func NewAuthorizationsCaller(address common.Address, caller bind.ContractCaller) (*AuthorizationsCaller, error) {
	contract, err := bindAuthorizations(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizationsCaller{contract: contract}, nil
}

// NewAuthorizationsTransactor creates a new write-only instance of Authorizations, bound to a specific deployed contract.
func NewAuthorizationsTransactor(address common.Address, transactor bind.ContractTransactor) (*AuthorizationsTransactor, error) {
	contract, err := bindAuthorizations(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizationsTransactor{contract: contract}, nil
}

// NewAuthorizationsFilterer creates a new log filterer instance of Authorizations, bound to a specific deployed contract.
func NewAuthorizationsFilterer(address common.Address, filterer bind.ContractFilterer) (*AuthorizationsFilterer, error) {
	contract, err := bindAuthorizations(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuthorizationsFilterer{contract: contract}, nil
}

// bindAuthorizations binds a generic wrapper to an already deployed contract.
func bindAuthorizations(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuthorizationsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Authorizations *AuthorizationsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Authorizations.Contract.AuthorizationsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Authorizations *AuthorizationsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Authorizations.Contract.AuthorizationsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Authorizations *AuthorizationsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Authorizations.Contract.AuthorizationsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Authorizations *AuthorizationsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Authorizations.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Authorizations *AuthorizationsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Authorizations.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Authorizations *AuthorizationsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Authorizations.Contract.contract.Transact(opts, method, params...)
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_Authorizations *AuthorizationsCaller) AuthorizerOf(opts *bind.CallOpts, _operator common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Authorizations.contract.Call(opts, out, "authorizerOf", _operator)
	return *ret0, err
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_Authorizations *AuthorizationsSession) AuthorizerOf(_operator common.Address) (common.Address, error) {
	return _Authorizations.Contract.AuthorizerOf(&_Authorizations.CallOpts, _operator)
}

// AuthorizerOf is a free data retrieval call binding the contract method 0xfb1677b1.
//
// Solidity: function authorizerOf(address _operator) constant returns(address)
func (_Authorizations *AuthorizationsCallerSession) AuthorizerOf(_operator common.Address) (common.Address, error) {
	return _Authorizations.Contract.AuthorizerOf(&_Authorizations.CallOpts, _operator)
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_Authorizations *AuthorizationsCaller) GetAuthoritySource(opts *bind.CallOpts, operatorContract common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Authorizations.contract.Call(opts, out, "getAuthoritySource", operatorContract)
	return *ret0, err
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_Authorizations *AuthorizationsSession) GetAuthoritySource(operatorContract common.Address) (common.Address, error) {
	return _Authorizations.Contract.GetAuthoritySource(&_Authorizations.CallOpts, operatorContract)
}

// GetAuthoritySource is a free data retrieval call binding the contract method 0xcbe945dc.
//
// Solidity: function getAuthoritySource(address operatorContract) constant returns(address)
func (_Authorizations *AuthorizationsCallerSession) GetAuthoritySource(operatorContract common.Address) (common.Address, error) {
	return _Authorizations.Contract.GetAuthoritySource(&_Authorizations.CallOpts, operatorContract)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsCaller) IsApprovedOperatorContract(opts *bind.CallOpts, _operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Authorizations.contract.Call(opts, out, "isApprovedOperatorContract", _operatorContract)
	return *ret0, err
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsSession) IsApprovedOperatorContract(_operatorContract common.Address) (bool, error) {
	return _Authorizations.Contract.IsApprovedOperatorContract(&_Authorizations.CallOpts, _operatorContract)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsCallerSession) IsApprovedOperatorContract(_operatorContract common.Address) (bool, error) {
	return _Authorizations.Contract.IsApprovedOperatorContract(&_Authorizations.CallOpts, _operatorContract)
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsCaller) IsAuthorizedForOperator(opts *bind.CallOpts, _operator common.Address, _operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Authorizations.contract.Call(opts, out, "isAuthorizedForOperator", _operator, _operatorContract)
	return *ret0, err
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsSession) IsAuthorizedForOperator(_operator common.Address, _operatorContract common.Address) (bool, error) {
	return _Authorizations.Contract.IsAuthorizedForOperator(&_Authorizations.CallOpts, _operator, _operatorContract)
}

// IsAuthorizedForOperator is a free data retrieval call binding the contract method 0xef1f9661.
//
// Solidity: function isAuthorizedForOperator(address _operator, address _operatorContract) constant returns(bool)
func (_Authorizations *AuthorizationsCallerSession) IsAuthorizedForOperator(_operator common.Address, _operatorContract common.Address) (bool, error) {
	return _Authorizations.Contract.IsAuthorizedForOperator(&_Authorizations.CallOpts, _operator, _operatorContract)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_Authorizations *AuthorizationsTransactor) AuthorizeOperatorContract(opts *bind.TransactOpts, _operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _Authorizations.contract.Transact(opts, "authorizeOperatorContract", _operator, _operatorContract)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_Authorizations *AuthorizationsSession) AuthorizeOperatorContract(_operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _Authorizations.Contract.AuthorizeOperatorContract(&_Authorizations.TransactOpts, _operator, _operatorContract)
}

// AuthorizeOperatorContract is a paid mutator transaction binding the contract method 0xf1654783.
//
// Solidity: function authorizeOperatorContract(address _operator, address _operatorContract) returns()
func (_Authorizations *AuthorizationsTransactorSession) AuthorizeOperatorContract(_operator common.Address, _operatorContract common.Address) (*types.Transaction, error) {
	return _Authorizations.Contract.AuthorizeOperatorContract(&_Authorizations.TransactOpts, _operator, _operatorContract)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_Authorizations *AuthorizationsTransactor) ClaimDelegatedAuthority(opts *bind.TransactOpts, delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _Authorizations.contract.Transact(opts, "claimDelegatedAuthority", delegatedAuthoritySource)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_Authorizations *AuthorizationsSession) ClaimDelegatedAuthority(delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _Authorizations.Contract.ClaimDelegatedAuthority(&_Authorizations.TransactOpts, delegatedAuthoritySource)
}

// ClaimDelegatedAuthority is a paid mutator transaction binding the contract method 0xa590ae36.
//
// Solidity: function claimDelegatedAuthority(address delegatedAuthoritySource) returns()
func (_Authorizations *AuthorizationsTransactorSession) ClaimDelegatedAuthority(delegatedAuthoritySource common.Address) (*types.Transaction, error) {
	return _Authorizations.Contract.ClaimDelegatedAuthority(&_Authorizations.TransactOpts, delegatedAuthoritySource)
}
