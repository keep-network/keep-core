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

// TokenGrantStakeABI is the input ABI used to generate the binding from.
const TokenGrantStakeABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_grantId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_tokenStaking\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancelStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_grantId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_tokenStaking\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGrantId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getStakingContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"recoverStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenGrantStake is an auto generated Go binding around an Ethereum contract.
type TokenGrantStake struct {
	TokenGrantStakeCaller     // Read-only binding to the contract
	TokenGrantStakeTransactor // Write-only binding to the contract
	TokenGrantStakeFilterer   // Log filterer for contract events
}

// TokenGrantStakeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenGrantStakeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantStakeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenGrantStakeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantStakeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenGrantStakeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantStakeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenGrantStakeSession struct {
	Contract     *TokenGrantStake  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenGrantStakeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenGrantStakeCallerSession struct {
	Contract *TokenGrantStakeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// TokenGrantStakeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenGrantStakeTransactorSession struct {
	Contract     *TokenGrantStakeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// TokenGrantStakeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenGrantStakeRaw struct {
	Contract *TokenGrantStake // Generic contract binding to access the raw methods on
}

// TokenGrantStakeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenGrantStakeCallerRaw struct {
	Contract *TokenGrantStakeCaller // Generic read-only contract binding to access the raw methods on
}

// TokenGrantStakeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenGrantStakeTransactorRaw struct {
	Contract *TokenGrantStakeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenGrantStake creates a new instance of TokenGrantStake, bound to a specific deployed contract.
func NewTokenGrantStake(address common.Address, backend bind.ContractBackend) (*TokenGrantStake, error) {
	contract, err := bindTokenGrantStake(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenGrantStake{TokenGrantStakeCaller: TokenGrantStakeCaller{contract: contract}, TokenGrantStakeTransactor: TokenGrantStakeTransactor{contract: contract}, TokenGrantStakeFilterer: TokenGrantStakeFilterer{contract: contract}}, nil
}

// NewTokenGrantStakeCaller creates a new read-only instance of TokenGrantStake, bound to a specific deployed contract.
func NewTokenGrantStakeCaller(address common.Address, caller bind.ContractCaller) (*TokenGrantStakeCaller, error) {
	contract, err := bindTokenGrantStake(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGrantStakeCaller{contract: contract}, nil
}

// NewTokenGrantStakeTransactor creates a new write-only instance of TokenGrantStake, bound to a specific deployed contract.
func NewTokenGrantStakeTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenGrantStakeTransactor, error) {
	contract, err := bindTokenGrantStake(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGrantStakeTransactor{contract: contract}, nil
}

// NewTokenGrantStakeFilterer creates a new log filterer instance of TokenGrantStake, bound to a specific deployed contract.
func NewTokenGrantStakeFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenGrantStakeFilterer, error) {
	contract, err := bindTokenGrantStake(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenGrantStakeFilterer{contract: contract}, nil
}

// bindTokenGrantStake binds a generic wrapper to an already deployed contract.
func bindTokenGrantStake(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenGrantStakeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGrantStake *TokenGrantStakeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenGrantStake.Contract.TokenGrantStakeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGrantStake *TokenGrantStakeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.TokenGrantStakeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGrantStake *TokenGrantStakeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.TokenGrantStakeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGrantStake *TokenGrantStakeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenGrantStake.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGrantStake *TokenGrantStakeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGrantStake *TokenGrantStakeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.contract.Transact(opts, method, params...)
}

// GetAmount is a free data retrieval call binding the contract method 0xd321fe29.
//
// Solidity: function getAmount() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeCaller) GetAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrantStake.contract.Call(opts, out, "getAmount")
	return *ret0, err
}

// GetAmount is a free data retrieval call binding the contract method 0xd321fe29.
//
// Solidity: function getAmount() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeSession) GetAmount() (*big.Int, error) {
	return _TokenGrantStake.Contract.GetAmount(&_TokenGrantStake.CallOpts)
}

// GetAmount is a free data retrieval call binding the contract method 0xd321fe29.
//
// Solidity: function getAmount() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeCallerSession) GetAmount() (*big.Int, error) {
	return _TokenGrantStake.Contract.GetAmount(&_TokenGrantStake.CallOpts)
}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() constant returns(uint256 _grantId, uint256 _amount, address _tokenStaking)
func (_TokenGrantStake *TokenGrantStakeCaller) GetDetails(opts *bind.CallOpts) (struct {
	GrantId      *big.Int
	Amount       *big.Int
	TokenStaking common.Address
}, error) {
	ret := new(struct {
		GrantId      *big.Int
		Amount       *big.Int
		TokenStaking common.Address
	})
	out := ret
	err := _TokenGrantStake.contract.Call(opts, out, "getDetails")
	return *ret, err
}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() constant returns(uint256 _grantId, uint256 _amount, address _tokenStaking)
func (_TokenGrantStake *TokenGrantStakeSession) GetDetails() (struct {
	GrantId      *big.Int
	Amount       *big.Int
	TokenStaking common.Address
}, error) {
	return _TokenGrantStake.Contract.GetDetails(&_TokenGrantStake.CallOpts)
}

// GetDetails is a free data retrieval call binding the contract method 0xfbbf93a0.
//
// Solidity: function getDetails() constant returns(uint256 _grantId, uint256 _amount, address _tokenStaking)
func (_TokenGrantStake *TokenGrantStakeCallerSession) GetDetails() (struct {
	GrantId      *big.Int
	Amount       *big.Int
	TokenStaking common.Address
}, error) {
	return _TokenGrantStake.Contract.GetDetails(&_TokenGrantStake.CallOpts)
}

// GetGrantId is a free data retrieval call binding the contract method 0x5c8fb3cf.
//
// Solidity: function getGrantId() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeCaller) GetGrantId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrantStake.contract.Call(opts, out, "getGrantId")
	return *ret0, err
}

// GetGrantId is a free data retrieval call binding the contract method 0x5c8fb3cf.
//
// Solidity: function getGrantId() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeSession) GetGrantId() (*big.Int, error) {
	return _TokenGrantStake.Contract.GetGrantId(&_TokenGrantStake.CallOpts)
}

// GetGrantId is a free data retrieval call binding the contract method 0x5c8fb3cf.
//
// Solidity: function getGrantId() constant returns(uint256)
func (_TokenGrantStake *TokenGrantStakeCallerSession) GetGrantId() (*big.Int, error) {
	return _TokenGrantStake.Contract.GetGrantId(&_TokenGrantStake.CallOpts)
}

// GetStakingContract is a free data retrieval call binding the contract method 0x8e68dce4.
//
// Solidity: function getStakingContract() constant returns(address)
func (_TokenGrantStake *TokenGrantStakeCaller) GetStakingContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenGrantStake.contract.Call(opts, out, "getStakingContract")
	return *ret0, err
}

// GetStakingContract is a free data retrieval call binding the contract method 0x8e68dce4.
//
// Solidity: function getStakingContract() constant returns(address)
func (_TokenGrantStake *TokenGrantStakeSession) GetStakingContract() (common.Address, error) {
	return _TokenGrantStake.Contract.GetStakingContract(&_TokenGrantStake.CallOpts)
}

// GetStakingContract is a free data retrieval call binding the contract method 0x8e68dce4.
//
// Solidity: function getStakingContract() constant returns(address)
func (_TokenGrantStake *TokenGrantStakeCallerSession) GetStakingContract() (common.Address, error) {
	return _TokenGrantStake.Contract.GetStakingContract(&_TokenGrantStake.CallOpts)
}

// CancelStake is a paid mutator transaction binding the contract method 0xca9ce291.
//
// Solidity: function cancelStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeTransactor) CancelStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrantStake.contract.Transact(opts, "cancelStake")
}

// CancelStake is a paid mutator transaction binding the contract method 0xca9ce291.
//
// Solidity: function cancelStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeSession) CancelStake() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.CancelStake(&_TokenGrantStake.TransactOpts)
}

// CancelStake is a paid mutator transaction binding the contract method 0xca9ce291.
//
// Solidity: function cancelStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeTransactorSession) CancelStake() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.CancelStake(&_TokenGrantStake.TransactOpts)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x2a9a347a.
//
// Solidity: function recoverStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeTransactor) RecoverStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrantStake.contract.Transact(opts, "recoverStake")
}

// RecoverStake is a paid mutator transaction binding the contract method 0x2a9a347a.
//
// Solidity: function recoverStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeSession) RecoverStake() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.RecoverStake(&_TokenGrantStake.TransactOpts)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x2a9a347a.
//
// Solidity: function recoverStake() returns(uint256)
func (_TokenGrantStake *TokenGrantStakeTransactorSession) RecoverStake() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.RecoverStake(&_TokenGrantStake.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x0e89439b.
//
// Solidity: function stake(uint256 _amount, bytes _extraData) returns()
func (_TokenGrantStake *TokenGrantStakeTransactor) Stake(opts *bind.TransactOpts, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrantStake.contract.Transact(opts, "stake", _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x0e89439b.
//
// Solidity: function stake(uint256 _amount, bytes _extraData) returns()
func (_TokenGrantStake *TokenGrantStakeSession) Stake(_amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.Stake(&_TokenGrantStake.TransactOpts, _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x0e89439b.
//
// Solidity: function stake(uint256 _amount, bytes _extraData) returns()
func (_TokenGrantStake *TokenGrantStakeTransactorSession) Stake(_amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrantStake.Contract.Stake(&_TokenGrantStake.TransactOpts, _amount, _extraData)
}

// Undelegate is a paid mutator transaction binding the contract method 0x92ab89bb.
//
// Solidity: function undelegate() returns()
func (_TokenGrantStake *TokenGrantStakeTransactor) Undelegate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrantStake.contract.Transact(opts, "undelegate")
}

// Undelegate is a paid mutator transaction binding the contract method 0x92ab89bb.
//
// Solidity: function undelegate() returns()
func (_TokenGrantStake *TokenGrantStakeSession) Undelegate() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.Undelegate(&_TokenGrantStake.TransactOpts)
}

// Undelegate is a paid mutator transaction binding the contract method 0x92ab89bb.
//
// Solidity: function undelegate() returns()
func (_TokenGrantStake *TokenGrantStakeTransactorSession) Undelegate() (*types.Transaction, error) {
	return _TokenGrantStake.Contract.Undelegate(&_TokenGrantStake.TransactOpts)
}
