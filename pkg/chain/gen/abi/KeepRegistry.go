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

// KeepRegistryABI is the input ABI used to generate the binding from.
const KeepRegistryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"defaultPanicButton\",\"type\":\"address\"}],\"name\":\"DefaultPanicButtonUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"governance\",\"type\":\"address\"}],\"name\":\"GovernanceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractPanicButtonDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"panicButton\",\"type\":\"address\"}],\"name\":\"OperatorContractPanicButtonUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"serviceContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"upgrader\",\"type\":\"address\"}],\"name\":\"OperatorContractUpgraderUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"registryKeeper\",\"type\":\"address\"}],\"name\":\"RegistryKeeperUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"keeper\",\"type\":\"address\"}],\"name\":\"ServiceContractUpgraderUpdated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"approveOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"defaultPanicButton\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"disableOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"disableOperatorContractPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"governance\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"isApprovedOperatorContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"isNewOperatorContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceContract\",\"type\":\"address\"}],\"name\":\"operatorContractUpgraderFor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorContractUpgraders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorContracts\",\"outputs\":[{\"internalType\":\"enumKeepRegistry.ContractStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"panicButtons\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"registryKeeper\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"serviceContractUpgraderFor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"serviceContractUpgraders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_panicButton\",\"type\":\"address\"}],\"name\":\"setDefaultPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_governance\",\"type\":\"address\"}],\"name\":\"setGovernance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_panicButton\",\"type\":\"address\"}],\"name\":\"setOperatorContractPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContractUpgrader\",\"type\":\"address\"}],\"name\":\"setOperatorContractUpgrader\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_registryKeeper\",\"type\":\"address\"}],\"name\":\"setRegistryKeeper\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_serviceContractUpgrader\",\"type\":\"address\"}],\"name\":\"setServiceContractUpgrader\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KeepRegistry is an auto generated Go binding around an Ethereum contract.
type KeepRegistry struct {
	KeepRegistryCaller     // Read-only binding to the contract
	KeepRegistryTransactor // Write-only binding to the contract
	KeepRegistryFilterer   // Log filterer for contract events
}

// KeepRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeepRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeepRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeepRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeepRegistrySession struct {
	Contract     *KeepRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeepRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeepRegistryCallerSession struct {
	Contract *KeepRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// KeepRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeepRegistryTransactorSession struct {
	Contract     *KeepRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// KeepRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeepRegistryRaw struct {
	Contract *KeepRegistry // Generic contract binding to access the raw methods on
}

// KeepRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeepRegistryCallerRaw struct {
	Contract *KeepRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// KeepRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeepRegistryTransactorRaw struct {
	Contract *KeepRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeepRegistry creates a new instance of KeepRegistry, bound to a specific deployed contract.
func NewKeepRegistry(address common.Address, backend bind.ContractBackend) (*KeepRegistry, error) {
	contract, err := bindKeepRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepRegistry{KeepRegistryCaller: KeepRegistryCaller{contract: contract}, KeepRegistryTransactor: KeepRegistryTransactor{contract: contract}, KeepRegistryFilterer: KeepRegistryFilterer{contract: contract}}, nil
}

// NewKeepRegistryCaller creates a new read-only instance of KeepRegistry, bound to a specific deployed contract.
func NewKeepRegistryCaller(address common.Address, caller bind.ContractCaller) (*KeepRegistryCaller, error) {
	contract, err := bindKeepRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRegistryCaller{contract: contract}, nil
}

// NewKeepRegistryTransactor creates a new write-only instance of KeepRegistry, bound to a specific deployed contract.
func NewKeepRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*KeepRegistryTransactor, error) {
	contract, err := bindKeepRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRegistryTransactor{contract: contract}, nil
}

// NewKeepRegistryFilterer creates a new log filterer instance of KeepRegistry, bound to a specific deployed contract.
func NewKeepRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*KeepRegistryFilterer, error) {
	contract, err := bindKeepRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepRegistryFilterer{contract: contract}, nil
}

// bindKeepRegistry binds a generic wrapper to an already deployed contract.
func bindKeepRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeepRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRegistry *KeepRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRegistry.Contract.KeepRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRegistry *KeepRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRegistry.Contract.KeepRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRegistry *KeepRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRegistry.Contract.KeepRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRegistry *KeepRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRegistry *KeepRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRegistry *KeepRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRegistry.Contract.contract.Transact(opts, method, params...)
}

// DefaultPanicButton is a free data retrieval call binding the contract method 0x2afb015e.
//
// Solidity: function defaultPanicButton() constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) DefaultPanicButton(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "defaultPanicButton")
	return *ret0, err
}

// DefaultPanicButton is a free data retrieval call binding the contract method 0x2afb015e.
//
// Solidity: function defaultPanicButton() constant returns(address)
func (_KeepRegistry *KeepRegistrySession) DefaultPanicButton() (common.Address, error) {
	return _KeepRegistry.Contract.DefaultPanicButton(&_KeepRegistry.CallOpts)
}

// DefaultPanicButton is a free data retrieval call binding the contract method 0x2afb015e.
//
// Solidity: function defaultPanicButton() constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) DefaultPanicButton() (common.Address, error) {
	return _KeepRegistry.Contract.DefaultPanicButton(&_KeepRegistry.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) Governance(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "governance")
	return *ret0, err
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() constant returns(address)
func (_KeepRegistry *KeepRegistrySession) Governance() (common.Address, error) {
	return _KeepRegistry.Contract.Governance(&_KeepRegistry.CallOpts)
}

// Governance is a free data retrieval call binding the contract method 0x5aa6e675.
//
// Solidity: function governance() constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) Governance() (common.Address, error) {
	return _KeepRegistry.Contract.Governance(&_KeepRegistry.CallOpts)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistryCaller) IsApprovedOperatorContract(opts *bind.CallOpts, operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "isApprovedOperatorContract", operatorContract)
	return *ret0, err
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistrySession) IsApprovedOperatorContract(operatorContract common.Address) (bool, error) {
	return _KeepRegistry.Contract.IsApprovedOperatorContract(&_KeepRegistry.CallOpts, operatorContract)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistryCallerSession) IsApprovedOperatorContract(operatorContract common.Address) (bool, error) {
	return _KeepRegistry.Contract.IsApprovedOperatorContract(&_KeepRegistry.CallOpts, operatorContract)
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistryCaller) IsNewOperatorContract(opts *bind.CallOpts, operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "isNewOperatorContract", operatorContract)
	return *ret0, err
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistrySession) IsNewOperatorContract(operatorContract common.Address) (bool, error) {
	return _KeepRegistry.Contract.IsNewOperatorContract(&_KeepRegistry.CallOpts, operatorContract)
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_KeepRegistry *KeepRegistryCallerSession) IsNewOperatorContract(operatorContract common.Address) (bool, error) {
	return _KeepRegistry.Contract.IsNewOperatorContract(&_KeepRegistry.CallOpts, operatorContract)
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) OperatorContractUpgraderFor(opts *bind.CallOpts, _serviceContract common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "operatorContractUpgraderFor", _serviceContract)
	return *ret0, err
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_KeepRegistry *KeepRegistrySession) OperatorContractUpgraderFor(_serviceContract common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.OperatorContractUpgraderFor(&_KeepRegistry.CallOpts, _serviceContract)
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) OperatorContractUpgraderFor(_serviceContract common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.OperatorContractUpgraderFor(&_KeepRegistry.CallOpts, _serviceContract)
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) OperatorContractUpgraders(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "operatorContractUpgraders", arg0)
	return *ret0, err
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistrySession) OperatorContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.OperatorContractUpgraders(&_KeepRegistry.CallOpts, arg0)
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) OperatorContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.OperatorContractUpgraders(&_KeepRegistry.CallOpts, arg0)
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_KeepRegistry *KeepRegistryCaller) OperatorContracts(opts *bind.CallOpts, arg0 common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "operatorContracts", arg0)
	return *ret0, err
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_KeepRegistry *KeepRegistrySession) OperatorContracts(arg0 common.Address) (uint8, error) {
	return _KeepRegistry.Contract.OperatorContracts(&_KeepRegistry.CallOpts, arg0)
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_KeepRegistry *KeepRegistryCallerSession) OperatorContracts(arg0 common.Address) (uint8, error) {
	return _KeepRegistry.Contract.OperatorContracts(&_KeepRegistry.CallOpts, arg0)
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) PanicButtons(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "panicButtons", arg0)
	return *ret0, err
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_KeepRegistry *KeepRegistrySession) PanicButtons(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.PanicButtons(&_KeepRegistry.CallOpts, arg0)
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) PanicButtons(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.PanicButtons(&_KeepRegistry.CallOpts, arg0)
}

// RegistryKeeper is a free data retrieval call binding the contract method 0x84c07bd2.
//
// Solidity: function registryKeeper() constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) RegistryKeeper(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "registryKeeper")
	return *ret0, err
}

// RegistryKeeper is a free data retrieval call binding the contract method 0x84c07bd2.
//
// Solidity: function registryKeeper() constant returns(address)
func (_KeepRegistry *KeepRegistrySession) RegistryKeeper() (common.Address, error) {
	return _KeepRegistry.Contract.RegistryKeeper(&_KeepRegistry.CallOpts)
}

// RegistryKeeper is a free data retrieval call binding the contract method 0x84c07bd2.
//
// Solidity: function registryKeeper() constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) RegistryKeeper() (common.Address, error) {
	return _KeepRegistry.Contract.RegistryKeeper(&_KeepRegistry.CallOpts)
}

// ServiceContractUpgraderFor is a free data retrieval call binding the contract method 0x6557eccf.
//
// Solidity: function serviceContractUpgraderFor(address _operatorContract) constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) ServiceContractUpgraderFor(opts *bind.CallOpts, _operatorContract common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "serviceContractUpgraderFor", _operatorContract)
	return *ret0, err
}

// ServiceContractUpgraderFor is a free data retrieval call binding the contract method 0x6557eccf.
//
// Solidity: function serviceContractUpgraderFor(address _operatorContract) constant returns(address)
func (_KeepRegistry *KeepRegistrySession) ServiceContractUpgraderFor(_operatorContract common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.ServiceContractUpgraderFor(&_KeepRegistry.CallOpts, _operatorContract)
}

// ServiceContractUpgraderFor is a free data retrieval call binding the contract method 0x6557eccf.
//
// Solidity: function serviceContractUpgraderFor(address _operatorContract) constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) ServiceContractUpgraderFor(_operatorContract common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.ServiceContractUpgraderFor(&_KeepRegistry.CallOpts, _operatorContract)
}

// ServiceContractUpgraders is a free data retrieval call binding the contract method 0xff111e3b.
//
// Solidity: function serviceContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCaller) ServiceContractUpgraders(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRegistry.contract.Call(opts, out, "serviceContractUpgraders", arg0)
	return *ret0, err
}

// ServiceContractUpgraders is a free data retrieval call binding the contract method 0xff111e3b.
//
// Solidity: function serviceContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistrySession) ServiceContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.ServiceContractUpgraders(&_KeepRegistry.CallOpts, arg0)
}

// ServiceContractUpgraders is a free data retrieval call binding the contract method 0xff111e3b.
//
// Solidity: function serviceContractUpgraders(address ) constant returns(address)
func (_KeepRegistry *KeepRegistryCallerSession) ServiceContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _KeepRegistry.Contract.ServiceContractUpgraders(&_KeepRegistry.CallOpts, arg0)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactor) ApproveOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "approveOperatorContract", operatorContract)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistrySession) ApproveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.ApproveOperatorContract(&_KeepRegistry.TransactOpts, operatorContract)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) ApproveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.ApproveOperatorContract(&_KeepRegistry.TransactOpts, operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactor) DisableOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "disableOperatorContract", operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistrySession) DisableOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.DisableOperatorContract(&_KeepRegistry.TransactOpts, operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) DisableOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.DisableOperatorContract(&_KeepRegistry.TransactOpts, operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactor) DisableOperatorContractPanicButton(opts *bind.TransactOpts, _operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "disableOperatorContractPanicButton", _operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_KeepRegistry *KeepRegistrySession) DisableOperatorContractPanicButton(_operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.DisableOperatorContractPanicButton(&_KeepRegistry.TransactOpts, _operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) DisableOperatorContractPanicButton(_operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.DisableOperatorContractPanicButton(&_KeepRegistry.TransactOpts, _operatorContract)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetDefaultPanicButton(opts *bind.TransactOpts, _panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setDefaultPanicButton", _panicButton)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_KeepRegistry *KeepRegistrySession) SetDefaultPanicButton(_panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetDefaultPanicButton(&_KeepRegistry.TransactOpts, _panicButton)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetDefaultPanicButton(_panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetDefaultPanicButton(&_KeepRegistry.TransactOpts, _panicButton)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetGovernance(opts *bind.TransactOpts, _governance common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setGovernance", _governance)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_KeepRegistry *KeepRegistrySession) SetGovernance(_governance common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetGovernance(&_KeepRegistry.TransactOpts, _governance)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetGovernance(_governance common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetGovernance(&_KeepRegistry.TransactOpts, _governance)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetOperatorContractPanicButton(opts *bind.TransactOpts, _operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setOperatorContractPanicButton", _operatorContract, _panicButton)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_KeepRegistry *KeepRegistrySession) SetOperatorContractPanicButton(_operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetOperatorContractPanicButton(&_KeepRegistry.TransactOpts, _operatorContract, _panicButton)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetOperatorContractPanicButton(_operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetOperatorContractPanicButton(&_KeepRegistry.TransactOpts, _operatorContract, _panicButton)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetOperatorContractUpgrader(opts *bind.TransactOpts, _serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setOperatorContractUpgrader", _serviceContract, _operatorContractUpgrader)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_KeepRegistry *KeepRegistrySession) SetOperatorContractUpgrader(_serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetOperatorContractUpgrader(&_KeepRegistry.TransactOpts, _serviceContract, _operatorContractUpgrader)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetOperatorContractUpgrader(_serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetOperatorContractUpgrader(&_KeepRegistry.TransactOpts, _serviceContract, _operatorContractUpgrader)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetRegistryKeeper(opts *bind.TransactOpts, _registryKeeper common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setRegistryKeeper", _registryKeeper)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_KeepRegistry *KeepRegistrySession) SetRegistryKeeper(_registryKeeper common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetRegistryKeeper(&_KeepRegistry.TransactOpts, _registryKeeper)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetRegistryKeeper(_registryKeeper common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetRegistryKeeper(&_KeepRegistry.TransactOpts, _registryKeeper)
}

// SetServiceContractUpgrader is a paid mutator transaction binding the contract method 0x3ae70f46.
//
// Solidity: function setServiceContractUpgrader(address _operatorContract, address _serviceContractUpgrader) returns()
func (_KeepRegistry *KeepRegistryTransactor) SetServiceContractUpgrader(opts *bind.TransactOpts, _operatorContract common.Address, _serviceContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.contract.Transact(opts, "setServiceContractUpgrader", _operatorContract, _serviceContractUpgrader)
}

// SetServiceContractUpgrader is a paid mutator transaction binding the contract method 0x3ae70f46.
//
// Solidity: function setServiceContractUpgrader(address _operatorContract, address _serviceContractUpgrader) returns()
func (_KeepRegistry *KeepRegistrySession) SetServiceContractUpgrader(_operatorContract common.Address, _serviceContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetServiceContractUpgrader(&_KeepRegistry.TransactOpts, _operatorContract, _serviceContractUpgrader)
}

// SetServiceContractUpgrader is a paid mutator transaction binding the contract method 0x3ae70f46.
//
// Solidity: function setServiceContractUpgrader(address _operatorContract, address _serviceContractUpgrader) returns()
func (_KeepRegistry *KeepRegistryTransactorSession) SetServiceContractUpgrader(_operatorContract common.Address, _serviceContractUpgrader common.Address) (*types.Transaction, error) {
	return _KeepRegistry.Contract.SetServiceContractUpgrader(&_KeepRegistry.TransactOpts, _operatorContract, _serviceContractUpgrader)
}

// KeepRegistryDefaultPanicButtonUpdatedIterator is returned from FilterDefaultPanicButtonUpdated and is used to iterate over the raw logs and unpacked data for DefaultPanicButtonUpdated events raised by the KeepRegistry contract.
type KeepRegistryDefaultPanicButtonUpdatedIterator struct {
	Event *KeepRegistryDefaultPanicButtonUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryDefaultPanicButtonUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryDefaultPanicButtonUpdated)
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
		it.Event = new(KeepRegistryDefaultPanicButtonUpdated)
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
func (it *KeepRegistryDefaultPanicButtonUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryDefaultPanicButtonUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryDefaultPanicButtonUpdated represents a DefaultPanicButtonUpdated event raised by the KeepRegistry contract.
type KeepRegistryDefaultPanicButtonUpdated struct {
	DefaultPanicButton common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDefaultPanicButtonUpdated is a free log retrieval operation binding the contract event 0x06f6f3a56d213caa17abfd21597ab6f8660771e05f9d430505ce755a2fdbb551.
//
// Solidity: event DefaultPanicButtonUpdated(address defaultPanicButton)
func (_KeepRegistry *KeepRegistryFilterer) FilterDefaultPanicButtonUpdated(opts *bind.FilterOpts) (*KeepRegistryDefaultPanicButtonUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "DefaultPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryDefaultPanicButtonUpdatedIterator{contract: _KeepRegistry.contract, event: "DefaultPanicButtonUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultPanicButtonUpdated is a free log subscription operation binding the contract event 0x06f6f3a56d213caa17abfd21597ab6f8660771e05f9d430505ce755a2fdbb551.
//
// Solidity: event DefaultPanicButtonUpdated(address defaultPanicButton)
func (_KeepRegistry *KeepRegistryFilterer) WatchDefaultPanicButtonUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryDefaultPanicButtonUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "DefaultPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryDefaultPanicButtonUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "DefaultPanicButtonUpdated", log); err != nil {
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

// ParseDefaultPanicButtonUpdated is a log parse operation binding the contract event 0x06f6f3a56d213caa17abfd21597ab6f8660771e05f9d430505ce755a2fdbb551.
//
// Solidity: event DefaultPanicButtonUpdated(address defaultPanicButton)
func (_KeepRegistry *KeepRegistryFilterer) ParseDefaultPanicButtonUpdated(log types.Log) (*KeepRegistryDefaultPanicButtonUpdated, error) {
	event := new(KeepRegistryDefaultPanicButtonUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "DefaultPanicButtonUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryGovernanceUpdatedIterator is returned from FilterGovernanceUpdated and is used to iterate over the raw logs and unpacked data for GovernanceUpdated events raised by the KeepRegistry contract.
type KeepRegistryGovernanceUpdatedIterator struct {
	Event *KeepRegistryGovernanceUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryGovernanceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryGovernanceUpdated)
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
		it.Event = new(KeepRegistryGovernanceUpdated)
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
func (it *KeepRegistryGovernanceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryGovernanceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryGovernanceUpdated represents a GovernanceUpdated event raised by the KeepRegistry contract.
type KeepRegistryGovernanceUpdated struct {
	Governance common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterGovernanceUpdated is a free log retrieval operation binding the contract event 0x9d3e522e1e47a2f6009739342b9cc7b252a1888154e843ab55ee1c81745795ab.
//
// Solidity: event GovernanceUpdated(address governance)
func (_KeepRegistry *KeepRegistryFilterer) FilterGovernanceUpdated(opts *bind.FilterOpts) (*KeepRegistryGovernanceUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "GovernanceUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryGovernanceUpdatedIterator{contract: _KeepRegistry.contract, event: "GovernanceUpdated", logs: logs, sub: sub}, nil
}

// WatchGovernanceUpdated is a free log subscription operation binding the contract event 0x9d3e522e1e47a2f6009739342b9cc7b252a1888154e843ab55ee1c81745795ab.
//
// Solidity: event GovernanceUpdated(address governance)
func (_KeepRegistry *KeepRegistryFilterer) WatchGovernanceUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryGovernanceUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "GovernanceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryGovernanceUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "GovernanceUpdated", log); err != nil {
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

// ParseGovernanceUpdated is a log parse operation binding the contract event 0x9d3e522e1e47a2f6009739342b9cc7b252a1888154e843ab55ee1c81745795ab.
//
// Solidity: event GovernanceUpdated(address governance)
func (_KeepRegistry *KeepRegistryFilterer) ParseGovernanceUpdated(log types.Log) (*KeepRegistryGovernanceUpdated, error) {
	event := new(KeepRegistryGovernanceUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "GovernanceUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryOperatorContractApprovedIterator is returned from FilterOperatorContractApproved and is used to iterate over the raw logs and unpacked data for OperatorContractApproved events raised by the KeepRegistry contract.
type KeepRegistryOperatorContractApprovedIterator struct {
	Event *KeepRegistryOperatorContractApproved // Event containing the contract specifics and raw log

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
func (it *KeepRegistryOperatorContractApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryOperatorContractApproved)
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
		it.Event = new(KeepRegistryOperatorContractApproved)
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
func (it *KeepRegistryOperatorContractApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryOperatorContractApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryOperatorContractApproved represents a OperatorContractApproved event raised by the KeepRegistry contract.
type KeepRegistryOperatorContractApproved struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractApproved is a free log retrieval operation binding the contract event 0xb5d7fbd1d67a09be2dd3ba169be0ebfb43f1d3fd31b12e9dc9e5bd98517bfafb.
//
// Solidity: event OperatorContractApproved(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) FilterOperatorContractApproved(opts *bind.FilterOpts) (*KeepRegistryOperatorContractApprovedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "OperatorContractApproved")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryOperatorContractApprovedIterator{contract: _KeepRegistry.contract, event: "OperatorContractApproved", logs: logs, sub: sub}, nil
}

// WatchOperatorContractApproved is a free log subscription operation binding the contract event 0xb5d7fbd1d67a09be2dd3ba169be0ebfb43f1d3fd31b12e9dc9e5bd98517bfafb.
//
// Solidity: event OperatorContractApproved(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) WatchOperatorContractApproved(opts *bind.WatchOpts, sink chan<- *KeepRegistryOperatorContractApproved) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "OperatorContractApproved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryOperatorContractApproved)
				if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractApproved", log); err != nil {
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

// ParseOperatorContractApproved is a log parse operation binding the contract event 0xb5d7fbd1d67a09be2dd3ba169be0ebfb43f1d3fd31b12e9dc9e5bd98517bfafb.
//
// Solidity: event OperatorContractApproved(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) ParseOperatorContractApproved(log types.Log) (*KeepRegistryOperatorContractApproved, error) {
	event := new(KeepRegistryOperatorContractApproved)
	if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractApproved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryOperatorContractDisabledIterator is returned from FilterOperatorContractDisabled and is used to iterate over the raw logs and unpacked data for OperatorContractDisabled events raised by the KeepRegistry contract.
type KeepRegistryOperatorContractDisabledIterator struct {
	Event *KeepRegistryOperatorContractDisabled // Event containing the contract specifics and raw log

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
func (it *KeepRegistryOperatorContractDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryOperatorContractDisabled)
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
		it.Event = new(KeepRegistryOperatorContractDisabled)
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
func (it *KeepRegistryOperatorContractDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryOperatorContractDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryOperatorContractDisabled represents a OperatorContractDisabled event raised by the KeepRegistry contract.
type KeepRegistryOperatorContractDisabled struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractDisabled is a free log retrieval operation binding the contract event 0x61e8cab67fc6f073d7e3bbd8d936b519daf95f5655e8b22563c63eaba75f8fa3.
//
// Solidity: event OperatorContractDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) FilterOperatorContractDisabled(opts *bind.FilterOpts) (*KeepRegistryOperatorContractDisabledIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "OperatorContractDisabled")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryOperatorContractDisabledIterator{contract: _KeepRegistry.contract, event: "OperatorContractDisabled", logs: logs, sub: sub}, nil
}

// WatchOperatorContractDisabled is a free log subscription operation binding the contract event 0x61e8cab67fc6f073d7e3bbd8d936b519daf95f5655e8b22563c63eaba75f8fa3.
//
// Solidity: event OperatorContractDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) WatchOperatorContractDisabled(opts *bind.WatchOpts, sink chan<- *KeepRegistryOperatorContractDisabled) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "OperatorContractDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryOperatorContractDisabled)
				if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractDisabled", log); err != nil {
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

// ParseOperatorContractDisabled is a log parse operation binding the contract event 0x61e8cab67fc6f073d7e3bbd8d936b519daf95f5655e8b22563c63eaba75f8fa3.
//
// Solidity: event OperatorContractDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) ParseOperatorContractDisabled(log types.Log) (*KeepRegistryOperatorContractDisabled, error) {
	event := new(KeepRegistryOperatorContractDisabled)
	if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractDisabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryOperatorContractPanicButtonDisabledIterator is returned from FilterOperatorContractPanicButtonDisabled and is used to iterate over the raw logs and unpacked data for OperatorContractPanicButtonDisabled events raised by the KeepRegistry contract.
type KeepRegistryOperatorContractPanicButtonDisabledIterator struct {
	Event *KeepRegistryOperatorContractPanicButtonDisabled // Event containing the contract specifics and raw log

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
func (it *KeepRegistryOperatorContractPanicButtonDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryOperatorContractPanicButtonDisabled)
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
		it.Event = new(KeepRegistryOperatorContractPanicButtonDisabled)
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
func (it *KeepRegistryOperatorContractPanicButtonDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryOperatorContractPanicButtonDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryOperatorContractPanicButtonDisabled represents a OperatorContractPanicButtonDisabled event raised by the KeepRegistry contract.
type KeepRegistryOperatorContractPanicButtonDisabled struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractPanicButtonDisabled is a free log retrieval operation binding the contract event 0x825c13d61db27b3c0dae8523f70310e87c2dfd99ea7023664547208322ff679c.
//
// Solidity: event OperatorContractPanicButtonDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) FilterOperatorContractPanicButtonDisabled(opts *bind.FilterOpts) (*KeepRegistryOperatorContractPanicButtonDisabledIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "OperatorContractPanicButtonDisabled")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryOperatorContractPanicButtonDisabledIterator{contract: _KeepRegistry.contract, event: "OperatorContractPanicButtonDisabled", logs: logs, sub: sub}, nil
}

// WatchOperatorContractPanicButtonDisabled is a free log subscription operation binding the contract event 0x825c13d61db27b3c0dae8523f70310e87c2dfd99ea7023664547208322ff679c.
//
// Solidity: event OperatorContractPanicButtonDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) WatchOperatorContractPanicButtonDisabled(opts *bind.WatchOpts, sink chan<- *KeepRegistryOperatorContractPanicButtonDisabled) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "OperatorContractPanicButtonDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryOperatorContractPanicButtonDisabled)
				if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractPanicButtonDisabled", log); err != nil {
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

// ParseOperatorContractPanicButtonDisabled is a log parse operation binding the contract event 0x825c13d61db27b3c0dae8523f70310e87c2dfd99ea7023664547208322ff679c.
//
// Solidity: event OperatorContractPanicButtonDisabled(address operatorContract)
func (_KeepRegistry *KeepRegistryFilterer) ParseOperatorContractPanicButtonDisabled(log types.Log) (*KeepRegistryOperatorContractPanicButtonDisabled, error) {
	event := new(KeepRegistryOperatorContractPanicButtonDisabled)
	if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractPanicButtonDisabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryOperatorContractPanicButtonUpdatedIterator is returned from FilterOperatorContractPanicButtonUpdated and is used to iterate over the raw logs and unpacked data for OperatorContractPanicButtonUpdated events raised by the KeepRegistry contract.
type KeepRegistryOperatorContractPanicButtonUpdatedIterator struct {
	Event *KeepRegistryOperatorContractPanicButtonUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryOperatorContractPanicButtonUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryOperatorContractPanicButtonUpdated)
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
		it.Event = new(KeepRegistryOperatorContractPanicButtonUpdated)
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
func (it *KeepRegistryOperatorContractPanicButtonUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryOperatorContractPanicButtonUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryOperatorContractPanicButtonUpdated represents a OperatorContractPanicButtonUpdated event raised by the KeepRegistry contract.
type KeepRegistryOperatorContractPanicButtonUpdated struct {
	OperatorContract common.Address
	PanicButton      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractPanicButtonUpdated is a free log retrieval operation binding the contract event 0x5f94d794a8a48b75e2ca89133a5b8748d14217648260c5247e57e3d314fd52a7.
//
// Solidity: event OperatorContractPanicButtonUpdated(address operatorContract, address panicButton)
func (_KeepRegistry *KeepRegistryFilterer) FilterOperatorContractPanicButtonUpdated(opts *bind.FilterOpts) (*KeepRegistryOperatorContractPanicButtonUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "OperatorContractPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryOperatorContractPanicButtonUpdatedIterator{contract: _KeepRegistry.contract, event: "OperatorContractPanicButtonUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorContractPanicButtonUpdated is a free log subscription operation binding the contract event 0x5f94d794a8a48b75e2ca89133a5b8748d14217648260c5247e57e3d314fd52a7.
//
// Solidity: event OperatorContractPanicButtonUpdated(address operatorContract, address panicButton)
func (_KeepRegistry *KeepRegistryFilterer) WatchOperatorContractPanicButtonUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryOperatorContractPanicButtonUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "OperatorContractPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryOperatorContractPanicButtonUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractPanicButtonUpdated", log); err != nil {
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

// ParseOperatorContractPanicButtonUpdated is a log parse operation binding the contract event 0x5f94d794a8a48b75e2ca89133a5b8748d14217648260c5247e57e3d314fd52a7.
//
// Solidity: event OperatorContractPanicButtonUpdated(address operatorContract, address panicButton)
func (_KeepRegistry *KeepRegistryFilterer) ParseOperatorContractPanicButtonUpdated(log types.Log) (*KeepRegistryOperatorContractPanicButtonUpdated, error) {
	event := new(KeepRegistryOperatorContractPanicButtonUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractPanicButtonUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryOperatorContractUpgraderUpdatedIterator is returned from FilterOperatorContractUpgraderUpdated and is used to iterate over the raw logs and unpacked data for OperatorContractUpgraderUpdated events raised by the KeepRegistry contract.
type KeepRegistryOperatorContractUpgraderUpdatedIterator struct {
	Event *KeepRegistryOperatorContractUpgraderUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryOperatorContractUpgraderUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryOperatorContractUpgraderUpdated)
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
		it.Event = new(KeepRegistryOperatorContractUpgraderUpdated)
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
func (it *KeepRegistryOperatorContractUpgraderUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryOperatorContractUpgraderUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryOperatorContractUpgraderUpdated represents a OperatorContractUpgraderUpdated event raised by the KeepRegistry contract.
type KeepRegistryOperatorContractUpgraderUpdated struct {
	ServiceContract common.Address
	Upgrader        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractUpgraderUpdated is a free log retrieval operation binding the contract event 0x6f5a686a46a4dd38461a6a5ce54c3b8dd7b989914b5a0748c8d6669c79477852.
//
// Solidity: event OperatorContractUpgraderUpdated(address serviceContract, address upgrader)
func (_KeepRegistry *KeepRegistryFilterer) FilterOperatorContractUpgraderUpdated(opts *bind.FilterOpts) (*KeepRegistryOperatorContractUpgraderUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "OperatorContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryOperatorContractUpgraderUpdatedIterator{contract: _KeepRegistry.contract, event: "OperatorContractUpgraderUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorContractUpgraderUpdated is a free log subscription operation binding the contract event 0x6f5a686a46a4dd38461a6a5ce54c3b8dd7b989914b5a0748c8d6669c79477852.
//
// Solidity: event OperatorContractUpgraderUpdated(address serviceContract, address upgrader)
func (_KeepRegistry *KeepRegistryFilterer) WatchOperatorContractUpgraderUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryOperatorContractUpgraderUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "OperatorContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryOperatorContractUpgraderUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractUpgraderUpdated", log); err != nil {
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

// ParseOperatorContractUpgraderUpdated is a log parse operation binding the contract event 0x6f5a686a46a4dd38461a6a5ce54c3b8dd7b989914b5a0748c8d6669c79477852.
//
// Solidity: event OperatorContractUpgraderUpdated(address serviceContract, address upgrader)
func (_KeepRegistry *KeepRegistryFilterer) ParseOperatorContractUpgraderUpdated(log types.Log) (*KeepRegistryOperatorContractUpgraderUpdated, error) {
	event := new(KeepRegistryOperatorContractUpgraderUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "OperatorContractUpgraderUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryRegistryKeeperUpdatedIterator is returned from FilterRegistryKeeperUpdated and is used to iterate over the raw logs and unpacked data for RegistryKeeperUpdated events raised by the KeepRegistry contract.
type KeepRegistryRegistryKeeperUpdatedIterator struct {
	Event *KeepRegistryRegistryKeeperUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryRegistryKeeperUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryRegistryKeeperUpdated)
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
		it.Event = new(KeepRegistryRegistryKeeperUpdated)
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
func (it *KeepRegistryRegistryKeeperUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryRegistryKeeperUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryRegistryKeeperUpdated represents a RegistryKeeperUpdated event raised by the KeepRegistry contract.
type KeepRegistryRegistryKeeperUpdated struct {
	RegistryKeeper common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRegistryKeeperUpdated is a free log retrieval operation binding the contract event 0xe5ecd8c922f03a416959a48c633d1bfdd1615ea46923a0a3e5d65ec90a38e038.
//
// Solidity: event RegistryKeeperUpdated(address registryKeeper)
func (_KeepRegistry *KeepRegistryFilterer) FilterRegistryKeeperUpdated(opts *bind.FilterOpts) (*KeepRegistryRegistryKeeperUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "RegistryKeeperUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryRegistryKeeperUpdatedIterator{contract: _KeepRegistry.contract, event: "RegistryKeeperUpdated", logs: logs, sub: sub}, nil
}

// WatchRegistryKeeperUpdated is a free log subscription operation binding the contract event 0xe5ecd8c922f03a416959a48c633d1bfdd1615ea46923a0a3e5d65ec90a38e038.
//
// Solidity: event RegistryKeeperUpdated(address registryKeeper)
func (_KeepRegistry *KeepRegistryFilterer) WatchRegistryKeeperUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryRegistryKeeperUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "RegistryKeeperUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryRegistryKeeperUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "RegistryKeeperUpdated", log); err != nil {
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

// ParseRegistryKeeperUpdated is a log parse operation binding the contract event 0xe5ecd8c922f03a416959a48c633d1bfdd1615ea46923a0a3e5d65ec90a38e038.
//
// Solidity: event RegistryKeeperUpdated(address registryKeeper)
func (_KeepRegistry *KeepRegistryFilterer) ParseRegistryKeeperUpdated(log types.Log) (*KeepRegistryRegistryKeeperUpdated, error) {
	event := new(KeepRegistryRegistryKeeperUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "RegistryKeeperUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRegistryServiceContractUpgraderUpdatedIterator is returned from FilterServiceContractUpgraderUpdated and is used to iterate over the raw logs and unpacked data for ServiceContractUpgraderUpdated events raised by the KeepRegistry contract.
type KeepRegistryServiceContractUpgraderUpdatedIterator struct {
	Event *KeepRegistryServiceContractUpgraderUpdated // Event containing the contract specifics and raw log

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
func (it *KeepRegistryServiceContractUpgraderUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRegistryServiceContractUpgraderUpdated)
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
		it.Event = new(KeepRegistryServiceContractUpgraderUpdated)
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
func (it *KeepRegistryServiceContractUpgraderUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRegistryServiceContractUpgraderUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRegistryServiceContractUpgraderUpdated represents a ServiceContractUpgraderUpdated event raised by the KeepRegistry contract.
type KeepRegistryServiceContractUpgraderUpdated struct {
	OperatorContract common.Address
	Keeper           common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterServiceContractUpgraderUpdated is a free log retrieval operation binding the contract event 0x067d7db72df8e5428d719a48ef14b072cdb3942706bfc6becc7198a316d73b78.
//
// Solidity: event ServiceContractUpgraderUpdated(address operatorContract, address keeper)
func (_KeepRegistry *KeepRegistryFilterer) FilterServiceContractUpgraderUpdated(opts *bind.FilterOpts) (*KeepRegistryServiceContractUpgraderUpdatedIterator, error) {

	logs, sub, err := _KeepRegistry.contract.FilterLogs(opts, "ServiceContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return &KeepRegistryServiceContractUpgraderUpdatedIterator{contract: _KeepRegistry.contract, event: "ServiceContractUpgraderUpdated", logs: logs, sub: sub}, nil
}

// WatchServiceContractUpgraderUpdated is a free log subscription operation binding the contract event 0x067d7db72df8e5428d719a48ef14b072cdb3942706bfc6becc7198a316d73b78.
//
// Solidity: event ServiceContractUpgraderUpdated(address operatorContract, address keeper)
func (_KeepRegistry *KeepRegistryFilterer) WatchServiceContractUpgraderUpdated(opts *bind.WatchOpts, sink chan<- *KeepRegistryServiceContractUpgraderUpdated) (event.Subscription, error) {

	logs, sub, err := _KeepRegistry.contract.WatchLogs(opts, "ServiceContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRegistryServiceContractUpgraderUpdated)
				if err := _KeepRegistry.contract.UnpackLog(event, "ServiceContractUpgraderUpdated", log); err != nil {
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

// ParseServiceContractUpgraderUpdated is a log parse operation binding the contract event 0x067d7db72df8e5428d719a48ef14b072cdb3942706bfc6becc7198a316d73b78.
//
// Solidity: event ServiceContractUpgraderUpdated(address operatorContract, address keeper)
func (_KeepRegistry *KeepRegistryFilterer) ParseServiceContractUpgraderUpdated(log types.Log) (*KeepRegistryServiceContractUpgraderUpdated, error) {
	event := new(KeepRegistryServiceContractUpgraderUpdated)
	if err := _KeepRegistry.contract.UnpackLog(event, "ServiceContractUpgraderUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
