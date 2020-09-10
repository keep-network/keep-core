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

// GasPriceOracleABI is the input ABI used to generate the binding from.
const GasPriceOracleABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"}],\"name\":\"GasPriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"consumerContract\",\"type\":\"address\"}],\"name\":\"addConsumerContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newGasPrice\",\"type\":\"uint256\"}],\"name\":\"beginGasPriceUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"consumerContracts\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeGasPriceUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gasPriceChangeInitiated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getConsumerContracts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"governanceDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"removeConsumerContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// GasPriceOracle is an auto generated Go binding around an Ethereum contract.
type GasPriceOracle struct {
	GasPriceOracleCaller     // Read-only binding to the contract
	GasPriceOracleTransactor // Write-only binding to the contract
	GasPriceOracleFilterer   // Log filterer for contract events
}

// GasPriceOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type GasPriceOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GasPriceOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GasPriceOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GasPriceOracleSession struct {
	Contract     *GasPriceOracle   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GasPriceOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GasPriceOracleCallerSession struct {
	Contract *GasPriceOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// GasPriceOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GasPriceOracleTransactorSession struct {
	Contract     *GasPriceOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// GasPriceOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type GasPriceOracleRaw struct {
	Contract *GasPriceOracle // Generic contract binding to access the raw methods on
}

// GasPriceOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GasPriceOracleCallerRaw struct {
	Contract *GasPriceOracleCaller // Generic read-only contract binding to access the raw methods on
}

// GasPriceOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GasPriceOracleTransactorRaw struct {
	Contract *GasPriceOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGasPriceOracle creates a new instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracle(address common.Address, backend bind.ContractBackend) (*GasPriceOracle, error) {
	contract, err := bindGasPriceOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracle{GasPriceOracleCaller: GasPriceOracleCaller{contract: contract}, GasPriceOracleTransactor: GasPriceOracleTransactor{contract: contract}, GasPriceOracleFilterer: GasPriceOracleFilterer{contract: contract}}, nil
}

// NewGasPriceOracleCaller creates a new read-only instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleCaller(address common.Address, caller bind.ContractCaller) (*GasPriceOracleCaller, error) {
	contract, err := bindGasPriceOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleCaller{contract: contract}, nil
}

// NewGasPriceOracleTransactor creates a new write-only instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*GasPriceOracleTransactor, error) {
	contract, err := bindGasPriceOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleTransactor{contract: contract}, nil
}

// NewGasPriceOracleFilterer creates a new log filterer instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*GasPriceOracleFilterer, error) {
	contract, err := bindGasPriceOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleFilterer{contract: contract}, nil
}

// bindGasPriceOracle binds a generic wrapper to an already deployed contract.
func bindGasPriceOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GasPriceOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasPriceOracle *GasPriceOracleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GasPriceOracle.Contract.GasPriceOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasPriceOracle *GasPriceOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.GasPriceOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasPriceOracle *GasPriceOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.GasPriceOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasPriceOracle *GasPriceOracleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _GasPriceOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasPriceOracle *GasPriceOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasPriceOracle *GasPriceOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.contract.Transact(opts, method, params...)
}

// ConsumerContracts is a free data retrieval call binding the contract method 0x9955efd0.
//
// Solidity: function consumerContracts(uint256 ) constant returns(address)
func (_GasPriceOracle *GasPriceOracleCaller) ConsumerContracts(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "consumerContracts", arg0)
	return *ret0, err
}

// ConsumerContracts is a free data retrieval call binding the contract method 0x9955efd0.
//
// Solidity: function consumerContracts(uint256 ) constant returns(address)
func (_GasPriceOracle *GasPriceOracleSession) ConsumerContracts(arg0 *big.Int) (common.Address, error) {
	return _GasPriceOracle.Contract.ConsumerContracts(&_GasPriceOracle.CallOpts, arg0)
}

// ConsumerContracts is a free data retrieval call binding the contract method 0x9955efd0.
//
// Solidity: function consumerContracts(uint256 ) constant returns(address)
func (_GasPriceOracle *GasPriceOracleCallerSession) ConsumerContracts(arg0 *big.Int) (common.Address, error) {
	return _GasPriceOracle.Contract.ConsumerContracts(&_GasPriceOracle.CallOpts, arg0)
}

// GasPrice is a free data retrieval call binding the contract method 0xfe173b97.
//
// Solidity: function gasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCaller) GasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "gasPrice")
	return *ret0, err
}

// GasPrice is a free data retrieval call binding the contract method 0xfe173b97.
//
// Solidity: function gasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleSession) GasPrice() (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPrice(&_GasPriceOracle.CallOpts)
}

// GasPrice is a free data retrieval call binding the contract method 0xfe173b97.
//
// Solidity: function gasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCallerSession) GasPrice() (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPrice(&_GasPriceOracle.CallOpts)
}

// GasPriceChangeInitiated is a free data retrieval call binding the contract method 0xe02e9783.
//
// Solidity: function gasPriceChangeInitiated() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCaller) GasPriceChangeInitiated(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "gasPriceChangeInitiated")
	return *ret0, err
}

// GasPriceChangeInitiated is a free data retrieval call binding the contract method 0xe02e9783.
//
// Solidity: function gasPriceChangeInitiated() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleSession) GasPriceChangeInitiated() (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPriceChangeInitiated(&_GasPriceOracle.CallOpts)
}

// GasPriceChangeInitiated is a free data retrieval call binding the contract method 0xe02e9783.
//
// Solidity: function gasPriceChangeInitiated() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCallerSession) GasPriceChangeInitiated() (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPriceChangeInitiated(&_GasPriceOracle.CallOpts)
}

// GetConsumerContracts is a free data retrieval call binding the contract method 0x5b4fb912.
//
// Solidity: function getConsumerContracts() constant returns(address[])
func (_GasPriceOracle *GasPriceOracleCaller) GetConsumerContracts(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "getConsumerContracts")
	return *ret0, err
}

// GetConsumerContracts is a free data retrieval call binding the contract method 0x5b4fb912.
//
// Solidity: function getConsumerContracts() constant returns(address[])
func (_GasPriceOracle *GasPriceOracleSession) GetConsumerContracts() ([]common.Address, error) {
	return _GasPriceOracle.Contract.GetConsumerContracts(&_GasPriceOracle.CallOpts)
}

// GetConsumerContracts is a free data retrieval call binding the contract method 0x5b4fb912.
//
// Solidity: function getConsumerContracts() constant returns(address[])
func (_GasPriceOracle *GasPriceOracleCallerSession) GetConsumerContracts() ([]common.Address, error) {
	return _GasPriceOracle.Contract.GetConsumerContracts(&_GasPriceOracle.CallOpts)
}

// GovernanceDelay is a free data retrieval call binding the contract method 0xbba32939.
//
// Solidity: function governanceDelay() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCaller) GovernanceDelay(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "governanceDelay")
	return *ret0, err
}

// GovernanceDelay is a free data retrieval call binding the contract method 0xbba32939.
//
// Solidity: function governanceDelay() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleSession) GovernanceDelay() (*big.Int, error) {
	return _GasPriceOracle.Contract.GovernanceDelay(&_GasPriceOracle.CallOpts)
}

// GovernanceDelay is a free data retrieval call binding the contract method 0xbba32939.
//
// Solidity: function governanceDelay() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCallerSession) GovernanceDelay() (*big.Int, error) {
	return _GasPriceOracle.Contract.GovernanceDelay(&_GasPriceOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GasPriceOracle *GasPriceOracleCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GasPriceOracle *GasPriceOracleSession) IsOwner() (bool, error) {
	return _GasPriceOracle.Contract.IsOwner(&_GasPriceOracle.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_GasPriceOracle *GasPriceOracleCallerSession) IsOwner() (bool, error) {
	return _GasPriceOracle.Contract.IsOwner(&_GasPriceOracle.CallOpts)
}

// NewGasPrice is a free data retrieval call binding the contract method 0xfad2c8fd.
//
// Solidity: function newGasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCaller) NewGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "newGasPrice")
	return *ret0, err
}

// NewGasPrice is a free data retrieval call binding the contract method 0xfad2c8fd.
//
// Solidity: function newGasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleSession) NewGasPrice() (*big.Int, error) {
	return _GasPriceOracle.Contract.NewGasPrice(&_GasPriceOracle.CallOpts)
}

// NewGasPrice is a free data retrieval call binding the contract method 0xfad2c8fd.
//
// Solidity: function newGasPrice() constant returns(uint256)
func (_GasPriceOracle *GasPriceOracleCallerSession) NewGasPrice() (*big.Int, error) {
	return _GasPriceOracle.Contract.NewGasPrice(&_GasPriceOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GasPriceOracle *GasPriceOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _GasPriceOracle.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GasPriceOracle *GasPriceOracleSession) Owner() (common.Address, error) {
	return _GasPriceOracle.Contract.Owner(&_GasPriceOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_GasPriceOracle *GasPriceOracleCallerSession) Owner() (common.Address, error) {
	return _GasPriceOracle.Contract.Owner(&_GasPriceOracle.CallOpts)
}

// AddConsumerContract is a paid mutator transaction binding the contract method 0xbcb22e4c.
//
// Solidity: function addConsumerContract(address consumerContract) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) AddConsumerContract(opts *bind.TransactOpts, consumerContract common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "addConsumerContract", consumerContract)
}

// AddConsumerContract is a paid mutator transaction binding the contract method 0xbcb22e4c.
//
// Solidity: function addConsumerContract(address consumerContract) returns()
func (_GasPriceOracle *GasPriceOracleSession) AddConsumerContract(consumerContract common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.AddConsumerContract(&_GasPriceOracle.TransactOpts, consumerContract)
}

// AddConsumerContract is a paid mutator transaction binding the contract method 0xbcb22e4c.
//
// Solidity: function addConsumerContract(address consumerContract) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) AddConsumerContract(consumerContract common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.AddConsumerContract(&_GasPriceOracle.TransactOpts, consumerContract)
}

// BeginGasPriceUpdate is a paid mutator transaction binding the contract method 0x4ce8b3bc.
//
// Solidity: function beginGasPriceUpdate(uint256 _newGasPrice) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) BeginGasPriceUpdate(opts *bind.TransactOpts, _newGasPrice *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "beginGasPriceUpdate", _newGasPrice)
}

// BeginGasPriceUpdate is a paid mutator transaction binding the contract method 0x4ce8b3bc.
//
// Solidity: function beginGasPriceUpdate(uint256 _newGasPrice) returns()
func (_GasPriceOracle *GasPriceOracleSession) BeginGasPriceUpdate(_newGasPrice *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.BeginGasPriceUpdate(&_GasPriceOracle.TransactOpts, _newGasPrice)
}

// BeginGasPriceUpdate is a paid mutator transaction binding the contract method 0x4ce8b3bc.
//
// Solidity: function beginGasPriceUpdate(uint256 _newGasPrice) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) BeginGasPriceUpdate(_newGasPrice *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.BeginGasPriceUpdate(&_GasPriceOracle.TransactOpts, _newGasPrice)
}

// FinalizeGasPriceUpdate is a paid mutator transaction binding the contract method 0xd615e4f4.
//
// Solidity: function finalizeGasPriceUpdate() returns()
func (_GasPriceOracle *GasPriceOracleTransactor) FinalizeGasPriceUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "finalizeGasPriceUpdate")
}

// FinalizeGasPriceUpdate is a paid mutator transaction binding the contract method 0xd615e4f4.
//
// Solidity: function finalizeGasPriceUpdate() returns()
func (_GasPriceOracle *GasPriceOracleSession) FinalizeGasPriceUpdate() (*types.Transaction, error) {
	return _GasPriceOracle.Contract.FinalizeGasPriceUpdate(&_GasPriceOracle.TransactOpts)
}

// FinalizeGasPriceUpdate is a paid mutator transaction binding the contract method 0xd615e4f4.
//
// Solidity: function finalizeGasPriceUpdate() returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) FinalizeGasPriceUpdate() (*types.Transaction, error) {
	return _GasPriceOracle.Contract.FinalizeGasPriceUpdate(&_GasPriceOracle.TransactOpts)
}

// RemoveConsumerContract is a paid mutator transaction binding the contract method 0xfb656c8c.
//
// Solidity: function removeConsumerContract(uint256 index) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) RemoveConsumerContract(opts *bind.TransactOpts, index *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "removeConsumerContract", index)
}

// RemoveConsumerContract is a paid mutator transaction binding the contract method 0xfb656c8c.
//
// Solidity: function removeConsumerContract(uint256 index) returns()
func (_GasPriceOracle *GasPriceOracleSession) RemoveConsumerContract(index *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.RemoveConsumerContract(&_GasPriceOracle.TransactOpts, index)
}

// RemoveConsumerContract is a paid mutator transaction binding the contract method 0xfb656c8c.
//
// Solidity: function removeConsumerContract(uint256 index) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) RemoveConsumerContract(index *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.RemoveConsumerContract(&_GasPriceOracle.TransactOpts, index)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GasPriceOracle *GasPriceOracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GasPriceOracle *GasPriceOracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _GasPriceOracle.Contract.RenounceOwnership(&_GasPriceOracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GasPriceOracle.Contract.RenounceOwnership(&_GasPriceOracle.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GasPriceOracle *GasPriceOracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.TransferOwnership(&_GasPriceOracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.TransferOwnership(&_GasPriceOracle.TransactOpts, newOwner)
}

// GasPriceOracleGasPriceUpdatedIterator is returned from FilterGasPriceUpdated and is used to iterate over the raw logs and unpacked data for GasPriceUpdated events raised by the GasPriceOracle contract.
type GasPriceOracleGasPriceUpdatedIterator struct {
	Event *GasPriceOracleGasPriceUpdated // Event containing the contract specifics and raw log

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
func (it *GasPriceOracleGasPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GasPriceOracleGasPriceUpdated)
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
		it.Event = new(GasPriceOracleGasPriceUpdated)
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
func (it *GasPriceOracleGasPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GasPriceOracleGasPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GasPriceOracleGasPriceUpdated represents a GasPriceUpdated event raised by the GasPriceOracle contract.
type GasPriceOracleGasPriceUpdated struct {
	NewValue *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterGasPriceUpdated is a free log retrieval operation binding the contract event 0xfcdccc6074c6c42e4bd578aa9870c697dc976a270968452d2b8c8dc369fae396.
//
// Solidity: event GasPriceUpdated(uint256 newValue)
func (_GasPriceOracle *GasPriceOracleFilterer) FilterGasPriceUpdated(opts *bind.FilterOpts) (*GasPriceOracleGasPriceUpdatedIterator, error) {

	logs, sub, err := _GasPriceOracle.contract.FilterLogs(opts, "GasPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleGasPriceUpdatedIterator{contract: _GasPriceOracle.contract, event: "GasPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchGasPriceUpdated is a free log subscription operation binding the contract event 0xfcdccc6074c6c42e4bd578aa9870c697dc976a270968452d2b8c8dc369fae396.
//
// Solidity: event GasPriceUpdated(uint256 newValue)
func (_GasPriceOracle *GasPriceOracleFilterer) WatchGasPriceUpdated(opts *bind.WatchOpts, sink chan<- *GasPriceOracleGasPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _GasPriceOracle.contract.WatchLogs(opts, "GasPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GasPriceOracleGasPriceUpdated)
				if err := _GasPriceOracle.contract.UnpackLog(event, "GasPriceUpdated", log); err != nil {
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

// ParseGasPriceUpdated is a log parse operation binding the contract event 0xfcdccc6074c6c42e4bd578aa9870c697dc976a270968452d2b8c8dc369fae396.
//
// Solidity: event GasPriceUpdated(uint256 newValue)
func (_GasPriceOracle *GasPriceOracleFilterer) ParseGasPriceUpdated(log types.Log) (*GasPriceOracleGasPriceUpdated, error) {
	event := new(GasPriceOracleGasPriceUpdated)
	if err := _GasPriceOracle.contract.UnpackLog(event, "GasPriceUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// GasPriceOracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GasPriceOracle contract.
type GasPriceOracleOwnershipTransferredIterator struct {
	Event *GasPriceOracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GasPriceOracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GasPriceOracleOwnershipTransferred)
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
		it.Event = new(GasPriceOracleOwnershipTransferred)
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
func (it *GasPriceOracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GasPriceOracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GasPriceOracleOwnershipTransferred represents a OwnershipTransferred event raised by the GasPriceOracle contract.
type GasPriceOracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GasPriceOracle *GasPriceOracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GasPriceOracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GasPriceOracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleOwnershipTransferredIterator{contract: _GasPriceOracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GasPriceOracle *GasPriceOracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GasPriceOracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GasPriceOracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GasPriceOracleOwnershipTransferred)
				if err := _GasPriceOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GasPriceOracle *GasPriceOracleFilterer) ParseOwnershipTransferred(log types.Log) (*GasPriceOracleOwnershipTransferred, error) {
	event := new(GasPriceOracleOwnershipTransferred)
	if err := _GasPriceOracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
