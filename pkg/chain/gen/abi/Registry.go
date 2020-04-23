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

// RegistryABI is the input ABI used to generate the binding from.
const RegistryABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DefaultPanicButtonUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"GovernanceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"OperatorContractPanicButtonDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"panicButton\",\"type\":\"address\"}],\"name\":\"OperatorContractPanicButtonUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"serviceContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"upgrader\",\"type\":\"address\"}],\"name\":\"OperatorContractUpgraderUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RegistryKeeperUpdated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"approveOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"disableOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"}],\"name\":\"disableOperatorContractPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"isApprovedOperatorContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"isNewOperatorContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceContract\",\"type\":\"address\"}],\"name\":\"operatorContractUpgraderFor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorContractUpgraders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorContracts\",\"outputs\":[{\"internalType\":\"enumRegistry.ContractStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"panicButtons\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_panicButton\",\"type\":\"address\"}],\"name\":\"setDefaultPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_governance\",\"type\":\"address\"}],\"name\":\"setGovernance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operatorContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_panicButton\",\"type\":\"address\"}],\"name\":\"setOperatorContractPanicButton\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operatorContractUpgrader\",\"type\":\"address\"}],\"name\":\"setOperatorContractUpgrader\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_registryKeeper\",\"type\":\"address\"}],\"name\":\"setRegistryKeeper\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistryCaller) IsApprovedOperatorContract(opts *bind.CallOpts, operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "isApprovedOperatorContract", operatorContract)
	return *ret0, err
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistrySession) IsApprovedOperatorContract(operatorContract common.Address) (bool, error) {
	return _Registry.Contract.IsApprovedOperatorContract(&_Registry.CallOpts, operatorContract)
}

// IsApprovedOperatorContract is a free data retrieval call binding the contract method 0x84d57689.
//
// Solidity: function isApprovedOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistryCallerSession) IsApprovedOperatorContract(operatorContract common.Address) (bool, error) {
	return _Registry.Contract.IsApprovedOperatorContract(&_Registry.CallOpts, operatorContract)
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistryCaller) IsNewOperatorContract(opts *bind.CallOpts, operatorContract common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "isNewOperatorContract", operatorContract)
	return *ret0, err
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistrySession) IsNewOperatorContract(operatorContract common.Address) (bool, error) {
	return _Registry.Contract.IsNewOperatorContract(&_Registry.CallOpts, operatorContract)
}

// IsNewOperatorContract is a free data retrieval call binding the contract method 0x5a9671f5.
//
// Solidity: function isNewOperatorContract(address operatorContract) constant returns(bool)
func (_Registry *RegistryCallerSession) IsNewOperatorContract(operatorContract common.Address) (bool, error) {
	return _Registry.Contract.IsNewOperatorContract(&_Registry.CallOpts, operatorContract)
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_Registry *RegistryCaller) OperatorContractUpgraderFor(opts *bind.CallOpts, _serviceContract common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "operatorContractUpgraderFor", _serviceContract)
	return *ret0, err
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_Registry *RegistrySession) OperatorContractUpgraderFor(_serviceContract common.Address) (common.Address, error) {
	return _Registry.Contract.OperatorContractUpgraderFor(&_Registry.CallOpts, _serviceContract)
}

// OperatorContractUpgraderFor is a free data retrieval call binding the contract method 0x63dc5fcf.
//
// Solidity: function operatorContractUpgraderFor(address _serviceContract) constant returns(address)
func (_Registry *RegistryCallerSession) OperatorContractUpgraderFor(_serviceContract common.Address) (common.Address, error) {
	return _Registry.Contract.OperatorContractUpgraderFor(&_Registry.CallOpts, _serviceContract)
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_Registry *RegistryCaller) OperatorContractUpgraders(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "operatorContractUpgraders", arg0)
	return *ret0, err
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_Registry *RegistrySession) OperatorContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.OperatorContractUpgraders(&_Registry.CallOpts, arg0)
}

// OperatorContractUpgraders is a free data retrieval call binding the contract method 0xd858008a.
//
// Solidity: function operatorContractUpgraders(address ) constant returns(address)
func (_Registry *RegistryCallerSession) OperatorContractUpgraders(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.OperatorContractUpgraders(&_Registry.CallOpts, arg0)
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_Registry *RegistryCaller) OperatorContracts(opts *bind.CallOpts, arg0 common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "operatorContracts", arg0)
	return *ret0, err
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_Registry *RegistrySession) OperatorContracts(arg0 common.Address) (uint8, error) {
	return _Registry.Contract.OperatorContracts(&_Registry.CallOpts, arg0)
}

// OperatorContracts is a free data retrieval call binding the contract method 0x3dbcab45.
//
// Solidity: function operatorContracts(address ) constant returns(uint8)
func (_Registry *RegistryCallerSession) OperatorContracts(arg0 common.Address) (uint8, error) {
	return _Registry.Contract.OperatorContracts(&_Registry.CallOpts, arg0)
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_Registry *RegistryCaller) PanicButtons(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Registry.contract.Call(opts, out, "panicButtons", arg0)
	return *ret0, err
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_Registry *RegistrySession) PanicButtons(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.PanicButtons(&_Registry.CallOpts, arg0)
}

// PanicButtons is a free data retrieval call binding the contract method 0x7b4df938.
//
// Solidity: function panicButtons(address ) constant returns(address)
func (_Registry *RegistryCallerSession) PanicButtons(arg0 common.Address) (common.Address, error) {
	return _Registry.Contract.PanicButtons(&_Registry.CallOpts, arg0)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_Registry *RegistryTransactor) ApproveOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "approveOperatorContract", operatorContract)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_Registry *RegistrySession) ApproveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ApproveOperatorContract(&_Registry.TransactOpts, operatorContract)
}

// ApproveOperatorContract is a paid mutator transaction binding the contract method 0x02db580b.
//
// Solidity: function approveOperatorContract(address operatorContract) returns()
func (_Registry *RegistryTransactorSession) ApproveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.ApproveOperatorContract(&_Registry.TransactOpts, operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_Registry *RegistryTransactor) DisableOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "disableOperatorContract", operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_Registry *RegistrySession) DisableOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.DisableOperatorContract(&_Registry.TransactOpts, operatorContract)
}

// DisableOperatorContract is a paid mutator transaction binding the contract method 0x09a6f73f.
//
// Solidity: function disableOperatorContract(address operatorContract) returns()
func (_Registry *RegistryTransactorSession) DisableOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.DisableOperatorContract(&_Registry.TransactOpts, operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_Registry *RegistryTransactor) DisableOperatorContractPanicButton(opts *bind.TransactOpts, _operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "disableOperatorContractPanicButton", _operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_Registry *RegistrySession) DisableOperatorContractPanicButton(_operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.DisableOperatorContractPanicButton(&_Registry.TransactOpts, _operatorContract)
}

// DisableOperatorContractPanicButton is a paid mutator transaction binding the contract method 0x637e9dbe.
//
// Solidity: function disableOperatorContractPanicButton(address _operatorContract) returns()
func (_Registry *RegistryTransactorSession) DisableOperatorContractPanicButton(_operatorContract common.Address) (*types.Transaction, error) {
	return _Registry.Contract.DisableOperatorContractPanicButton(&_Registry.TransactOpts, _operatorContract)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_Registry *RegistryTransactor) SetDefaultPanicButton(opts *bind.TransactOpts, _panicButton common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setDefaultPanicButton", _panicButton)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_Registry *RegistrySession) SetDefaultPanicButton(_panicButton common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetDefaultPanicButton(&_Registry.TransactOpts, _panicButton)
}

// SetDefaultPanicButton is a paid mutator transaction binding the contract method 0x618dae20.
//
// Solidity: function setDefaultPanicButton(address _panicButton) returns()
func (_Registry *RegistryTransactorSession) SetDefaultPanicButton(_panicButton common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetDefaultPanicButton(&_Registry.TransactOpts, _panicButton)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_Registry *RegistryTransactor) SetGovernance(opts *bind.TransactOpts, _governance common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setGovernance", _governance)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_Registry *RegistrySession) SetGovernance(_governance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetGovernance(&_Registry.TransactOpts, _governance)
}

// SetGovernance is a paid mutator transaction binding the contract method 0xab033ea9.
//
// Solidity: function setGovernance(address _governance) returns()
func (_Registry *RegistryTransactorSession) SetGovernance(_governance common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetGovernance(&_Registry.TransactOpts, _governance)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_Registry *RegistryTransactor) SetOperatorContractPanicButton(opts *bind.TransactOpts, _operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setOperatorContractPanicButton", _operatorContract, _panicButton)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_Registry *RegistrySession) SetOperatorContractPanicButton(_operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetOperatorContractPanicButton(&_Registry.TransactOpts, _operatorContract, _panicButton)
}

// SetOperatorContractPanicButton is a paid mutator transaction binding the contract method 0xd82437f4.
//
// Solidity: function setOperatorContractPanicButton(address _operatorContract, address _panicButton) returns()
func (_Registry *RegistryTransactorSession) SetOperatorContractPanicButton(_operatorContract common.Address, _panicButton common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetOperatorContractPanicButton(&_Registry.TransactOpts, _operatorContract, _panicButton)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_Registry *RegistryTransactor) SetOperatorContractUpgrader(opts *bind.TransactOpts, _serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setOperatorContractUpgrader", _serviceContract, _operatorContractUpgrader)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_Registry *RegistrySession) SetOperatorContractUpgrader(_serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetOperatorContractUpgrader(&_Registry.TransactOpts, _serviceContract, _operatorContractUpgrader)
}

// SetOperatorContractUpgrader is a paid mutator transaction binding the contract method 0x716d6061.
//
// Solidity: function setOperatorContractUpgrader(address _serviceContract, address _operatorContractUpgrader) returns()
func (_Registry *RegistryTransactorSession) SetOperatorContractUpgrader(_serviceContract common.Address, _operatorContractUpgrader common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetOperatorContractUpgrader(&_Registry.TransactOpts, _serviceContract, _operatorContractUpgrader)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_Registry *RegistryTransactor) SetRegistryKeeper(opts *bind.TransactOpts, _registryKeeper common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setRegistryKeeper", _registryKeeper)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_Registry *RegistrySession) SetRegistryKeeper(_registryKeeper common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetRegistryKeeper(&_Registry.TransactOpts, _registryKeeper)
}

// SetRegistryKeeper is a paid mutator transaction binding the contract method 0xdb816b62.
//
// Solidity: function setRegistryKeeper(address _registryKeeper) returns()
func (_Registry *RegistryTransactorSession) SetRegistryKeeper(_registryKeeper common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetRegistryKeeper(&_Registry.TransactOpts, _registryKeeper)
}

// RegistryDefaultPanicButtonUpdatedIterator is returned from FilterDefaultPanicButtonUpdated and is used to iterate over the raw logs and unpacked data for DefaultPanicButtonUpdated events raised by the Registry contract.
type RegistryDefaultPanicButtonUpdatedIterator struct {
	Event *RegistryDefaultPanicButtonUpdated // Event containing the contract specifics and raw log

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
func (it *RegistryDefaultPanicButtonUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryDefaultPanicButtonUpdated)
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
		it.Event = new(RegistryDefaultPanicButtonUpdated)
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
func (it *RegistryDefaultPanicButtonUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryDefaultPanicButtonUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryDefaultPanicButtonUpdated represents a DefaultPanicButtonUpdated event raised by the Registry contract.
type RegistryDefaultPanicButtonUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultPanicButtonUpdated is a free log retrieval operation binding the contract event 0xd2d642ac7b8a5ebfc5dcbe35662d1c1ba8c961ca8092caddf1ab66bb08940af4.
//
// Solidity: event DefaultPanicButtonUpdated()
func (_Registry *RegistryFilterer) FilterDefaultPanicButtonUpdated(opts *bind.FilterOpts) (*RegistryDefaultPanicButtonUpdatedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "DefaultPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return &RegistryDefaultPanicButtonUpdatedIterator{contract: _Registry.contract, event: "DefaultPanicButtonUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultPanicButtonUpdated is a free log subscription operation binding the contract event 0xd2d642ac7b8a5ebfc5dcbe35662d1c1ba8c961ca8092caddf1ab66bb08940af4.
//
// Solidity: event DefaultPanicButtonUpdated()
func (_Registry *RegistryFilterer) WatchDefaultPanicButtonUpdated(opts *bind.WatchOpts, sink chan<- *RegistryDefaultPanicButtonUpdated) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "DefaultPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryDefaultPanicButtonUpdated)
				if err := _Registry.contract.UnpackLog(event, "DefaultPanicButtonUpdated", log); err != nil {
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

// ParseDefaultPanicButtonUpdated is a log parse operation binding the contract event 0xd2d642ac7b8a5ebfc5dcbe35662d1c1ba8c961ca8092caddf1ab66bb08940af4.
//
// Solidity: event DefaultPanicButtonUpdated()
func (_Registry *RegistryFilterer) ParseDefaultPanicButtonUpdated(log types.Log) (*RegistryDefaultPanicButtonUpdated, error) {
	event := new(RegistryDefaultPanicButtonUpdated)
	if err := _Registry.contract.UnpackLog(event, "DefaultPanicButtonUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryGovernanceUpdatedIterator is returned from FilterGovernanceUpdated and is used to iterate over the raw logs and unpacked data for GovernanceUpdated events raised by the Registry contract.
type RegistryGovernanceUpdatedIterator struct {
	Event *RegistryGovernanceUpdated // Event containing the contract specifics and raw log

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
func (it *RegistryGovernanceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryGovernanceUpdated)
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
		it.Event = new(RegistryGovernanceUpdated)
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
func (it *RegistryGovernanceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryGovernanceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryGovernanceUpdated represents a GovernanceUpdated event raised by the Registry contract.
type RegistryGovernanceUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterGovernanceUpdated is a free log retrieval operation binding the contract event 0xfae4c5ac7a8b7a33c68caa47e5984dd4be94cb1a1a353742693624832df2bbd5.
//
// Solidity: event GovernanceUpdated()
func (_Registry *RegistryFilterer) FilterGovernanceUpdated(opts *bind.FilterOpts) (*RegistryGovernanceUpdatedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "GovernanceUpdated")
	if err != nil {
		return nil, err
	}
	return &RegistryGovernanceUpdatedIterator{contract: _Registry.contract, event: "GovernanceUpdated", logs: logs, sub: sub}, nil
}

// WatchGovernanceUpdated is a free log subscription operation binding the contract event 0xfae4c5ac7a8b7a33c68caa47e5984dd4be94cb1a1a353742693624832df2bbd5.
//
// Solidity: event GovernanceUpdated()
func (_Registry *RegistryFilterer) WatchGovernanceUpdated(opts *bind.WatchOpts, sink chan<- *RegistryGovernanceUpdated) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "GovernanceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryGovernanceUpdated)
				if err := _Registry.contract.UnpackLog(event, "GovernanceUpdated", log); err != nil {
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

// ParseGovernanceUpdated is a log parse operation binding the contract event 0xfae4c5ac7a8b7a33c68caa47e5984dd4be94cb1a1a353742693624832df2bbd5.
//
// Solidity: event GovernanceUpdated()
func (_Registry *RegistryFilterer) ParseGovernanceUpdated(log types.Log) (*RegistryGovernanceUpdated, error) {
	event := new(RegistryGovernanceUpdated)
	if err := _Registry.contract.UnpackLog(event, "GovernanceUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryOperatorContractApprovedIterator is returned from FilterOperatorContractApproved and is used to iterate over the raw logs and unpacked data for OperatorContractApproved events raised by the Registry contract.
type RegistryOperatorContractApprovedIterator struct {
	Event *RegistryOperatorContractApproved // Event containing the contract specifics and raw log

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
func (it *RegistryOperatorContractApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOperatorContractApproved)
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
		it.Event = new(RegistryOperatorContractApproved)
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
func (it *RegistryOperatorContractApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOperatorContractApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOperatorContractApproved represents a OperatorContractApproved event raised by the Registry contract.
type RegistryOperatorContractApproved struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractApproved is a free log retrieval operation binding the contract event 0xb5d7fbd1d67a09be2dd3ba169be0ebfb43f1d3fd31b12e9dc9e5bd98517bfafb.
//
// Solidity: event OperatorContractApproved(address operatorContract)
func (_Registry *RegistryFilterer) FilterOperatorContractApproved(opts *bind.FilterOpts) (*RegistryOperatorContractApprovedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OperatorContractApproved")
	if err != nil {
		return nil, err
	}
	return &RegistryOperatorContractApprovedIterator{contract: _Registry.contract, event: "OperatorContractApproved", logs: logs, sub: sub}, nil
}

// WatchOperatorContractApproved is a free log subscription operation binding the contract event 0xb5d7fbd1d67a09be2dd3ba169be0ebfb43f1d3fd31b12e9dc9e5bd98517bfafb.
//
// Solidity: event OperatorContractApproved(address operatorContract)
func (_Registry *RegistryFilterer) WatchOperatorContractApproved(opts *bind.WatchOpts, sink chan<- *RegistryOperatorContractApproved) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OperatorContractApproved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOperatorContractApproved)
				if err := _Registry.contract.UnpackLog(event, "OperatorContractApproved", log); err != nil {
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
func (_Registry *RegistryFilterer) ParseOperatorContractApproved(log types.Log) (*RegistryOperatorContractApproved, error) {
	event := new(RegistryOperatorContractApproved)
	if err := _Registry.contract.UnpackLog(event, "OperatorContractApproved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryOperatorContractDisabledIterator is returned from FilterOperatorContractDisabled and is used to iterate over the raw logs and unpacked data for OperatorContractDisabled events raised by the Registry contract.
type RegistryOperatorContractDisabledIterator struct {
	Event *RegistryOperatorContractDisabled // Event containing the contract specifics and raw log

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
func (it *RegistryOperatorContractDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOperatorContractDisabled)
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
		it.Event = new(RegistryOperatorContractDisabled)
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
func (it *RegistryOperatorContractDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOperatorContractDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOperatorContractDisabled represents a OperatorContractDisabled event raised by the Registry contract.
type RegistryOperatorContractDisabled struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractDisabled is a free log retrieval operation binding the contract event 0x61e8cab67fc6f073d7e3bbd8d936b519daf95f5655e8b22563c63eaba75f8fa3.
//
// Solidity: event OperatorContractDisabled(address operatorContract)
func (_Registry *RegistryFilterer) FilterOperatorContractDisabled(opts *bind.FilterOpts) (*RegistryOperatorContractDisabledIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OperatorContractDisabled")
	if err != nil {
		return nil, err
	}
	return &RegistryOperatorContractDisabledIterator{contract: _Registry.contract, event: "OperatorContractDisabled", logs: logs, sub: sub}, nil
}

// WatchOperatorContractDisabled is a free log subscription operation binding the contract event 0x61e8cab67fc6f073d7e3bbd8d936b519daf95f5655e8b22563c63eaba75f8fa3.
//
// Solidity: event OperatorContractDisabled(address operatorContract)
func (_Registry *RegistryFilterer) WatchOperatorContractDisabled(opts *bind.WatchOpts, sink chan<- *RegistryOperatorContractDisabled) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OperatorContractDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOperatorContractDisabled)
				if err := _Registry.contract.UnpackLog(event, "OperatorContractDisabled", log); err != nil {
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
func (_Registry *RegistryFilterer) ParseOperatorContractDisabled(log types.Log) (*RegistryOperatorContractDisabled, error) {
	event := new(RegistryOperatorContractDisabled)
	if err := _Registry.contract.UnpackLog(event, "OperatorContractDisabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryOperatorContractPanicButtonDisabledIterator is returned from FilterOperatorContractPanicButtonDisabled and is used to iterate over the raw logs and unpacked data for OperatorContractPanicButtonDisabled events raised by the Registry contract.
type RegistryOperatorContractPanicButtonDisabledIterator struct {
	Event *RegistryOperatorContractPanicButtonDisabled // Event containing the contract specifics and raw log

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
func (it *RegistryOperatorContractPanicButtonDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOperatorContractPanicButtonDisabled)
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
		it.Event = new(RegistryOperatorContractPanicButtonDisabled)
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
func (it *RegistryOperatorContractPanicButtonDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOperatorContractPanicButtonDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOperatorContractPanicButtonDisabled represents a OperatorContractPanicButtonDisabled event raised by the Registry contract.
type RegistryOperatorContractPanicButtonDisabled struct {
	OperatorContract common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractPanicButtonDisabled is a free log retrieval operation binding the contract event 0x825c13d61db27b3c0dae8523f70310e87c2dfd99ea7023664547208322ff679c.
//
// Solidity: event OperatorContractPanicButtonDisabled(address operatorContract)
func (_Registry *RegistryFilterer) FilterOperatorContractPanicButtonDisabled(opts *bind.FilterOpts) (*RegistryOperatorContractPanicButtonDisabledIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OperatorContractPanicButtonDisabled")
	if err != nil {
		return nil, err
	}
	return &RegistryOperatorContractPanicButtonDisabledIterator{contract: _Registry.contract, event: "OperatorContractPanicButtonDisabled", logs: logs, sub: sub}, nil
}

// WatchOperatorContractPanicButtonDisabled is a free log subscription operation binding the contract event 0x825c13d61db27b3c0dae8523f70310e87c2dfd99ea7023664547208322ff679c.
//
// Solidity: event OperatorContractPanicButtonDisabled(address operatorContract)
func (_Registry *RegistryFilterer) WatchOperatorContractPanicButtonDisabled(opts *bind.WatchOpts, sink chan<- *RegistryOperatorContractPanicButtonDisabled) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OperatorContractPanicButtonDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOperatorContractPanicButtonDisabled)
				if err := _Registry.contract.UnpackLog(event, "OperatorContractPanicButtonDisabled", log); err != nil {
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
func (_Registry *RegistryFilterer) ParseOperatorContractPanicButtonDisabled(log types.Log) (*RegistryOperatorContractPanicButtonDisabled, error) {
	event := new(RegistryOperatorContractPanicButtonDisabled)
	if err := _Registry.contract.UnpackLog(event, "OperatorContractPanicButtonDisabled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryOperatorContractPanicButtonUpdatedIterator is returned from FilterOperatorContractPanicButtonUpdated and is used to iterate over the raw logs and unpacked data for OperatorContractPanicButtonUpdated events raised by the Registry contract.
type RegistryOperatorContractPanicButtonUpdatedIterator struct {
	Event *RegistryOperatorContractPanicButtonUpdated // Event containing the contract specifics and raw log

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
func (it *RegistryOperatorContractPanicButtonUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOperatorContractPanicButtonUpdated)
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
		it.Event = new(RegistryOperatorContractPanicButtonUpdated)
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
func (it *RegistryOperatorContractPanicButtonUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOperatorContractPanicButtonUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOperatorContractPanicButtonUpdated represents a OperatorContractPanicButtonUpdated event raised by the Registry contract.
type RegistryOperatorContractPanicButtonUpdated struct {
	OperatorContract common.Address
	PanicButton      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractPanicButtonUpdated is a free log retrieval operation binding the contract event 0x5f94d794a8a48b75e2ca89133a5b8748d14217648260c5247e57e3d314fd52a7.
//
// Solidity: event OperatorContractPanicButtonUpdated(address operatorContract, address panicButton)
func (_Registry *RegistryFilterer) FilterOperatorContractPanicButtonUpdated(opts *bind.FilterOpts) (*RegistryOperatorContractPanicButtonUpdatedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OperatorContractPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return &RegistryOperatorContractPanicButtonUpdatedIterator{contract: _Registry.contract, event: "OperatorContractPanicButtonUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorContractPanicButtonUpdated is a free log subscription operation binding the contract event 0x5f94d794a8a48b75e2ca89133a5b8748d14217648260c5247e57e3d314fd52a7.
//
// Solidity: event OperatorContractPanicButtonUpdated(address operatorContract, address panicButton)
func (_Registry *RegistryFilterer) WatchOperatorContractPanicButtonUpdated(opts *bind.WatchOpts, sink chan<- *RegistryOperatorContractPanicButtonUpdated) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OperatorContractPanicButtonUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOperatorContractPanicButtonUpdated)
				if err := _Registry.contract.UnpackLog(event, "OperatorContractPanicButtonUpdated", log); err != nil {
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
func (_Registry *RegistryFilterer) ParseOperatorContractPanicButtonUpdated(log types.Log) (*RegistryOperatorContractPanicButtonUpdated, error) {
	event := new(RegistryOperatorContractPanicButtonUpdated)
	if err := _Registry.contract.UnpackLog(event, "OperatorContractPanicButtonUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryOperatorContractUpgraderUpdatedIterator is returned from FilterOperatorContractUpgraderUpdated and is used to iterate over the raw logs and unpacked data for OperatorContractUpgraderUpdated events raised by the Registry contract.
type RegistryOperatorContractUpgraderUpdatedIterator struct {
	Event *RegistryOperatorContractUpgraderUpdated // Event containing the contract specifics and raw log

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
func (it *RegistryOperatorContractUpgraderUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryOperatorContractUpgraderUpdated)
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
		it.Event = new(RegistryOperatorContractUpgraderUpdated)
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
func (it *RegistryOperatorContractUpgraderUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryOperatorContractUpgraderUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryOperatorContractUpgraderUpdated represents a OperatorContractUpgraderUpdated event raised by the Registry contract.
type RegistryOperatorContractUpgraderUpdated struct {
	ServiceContract common.Address
	Upgrader        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOperatorContractUpgraderUpdated is a free log retrieval operation binding the contract event 0x6f5a686a46a4dd38461a6a5ce54c3b8dd7b989914b5a0748c8d6669c79477852.
//
// Solidity: event OperatorContractUpgraderUpdated(address serviceContract, address upgrader)
func (_Registry *RegistryFilterer) FilterOperatorContractUpgraderUpdated(opts *bind.FilterOpts) (*RegistryOperatorContractUpgraderUpdatedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "OperatorContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return &RegistryOperatorContractUpgraderUpdatedIterator{contract: _Registry.contract, event: "OperatorContractUpgraderUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorContractUpgraderUpdated is a free log subscription operation binding the contract event 0x6f5a686a46a4dd38461a6a5ce54c3b8dd7b989914b5a0748c8d6669c79477852.
//
// Solidity: event OperatorContractUpgraderUpdated(address serviceContract, address upgrader)
func (_Registry *RegistryFilterer) WatchOperatorContractUpgraderUpdated(opts *bind.WatchOpts, sink chan<- *RegistryOperatorContractUpgraderUpdated) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "OperatorContractUpgraderUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryOperatorContractUpgraderUpdated)
				if err := _Registry.contract.UnpackLog(event, "OperatorContractUpgraderUpdated", log); err != nil {
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
func (_Registry *RegistryFilterer) ParseOperatorContractUpgraderUpdated(log types.Log) (*RegistryOperatorContractUpgraderUpdated, error) {
	event := new(RegistryOperatorContractUpgraderUpdated)
	if err := _Registry.contract.UnpackLog(event, "OperatorContractUpgraderUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RegistryRegistryKeeperUpdatedIterator is returned from FilterRegistryKeeperUpdated and is used to iterate over the raw logs and unpacked data for RegistryKeeperUpdated events raised by the Registry contract.
type RegistryRegistryKeeperUpdatedIterator struct {
	Event *RegistryRegistryKeeperUpdated // Event containing the contract specifics and raw log

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
func (it *RegistryRegistryKeeperUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryRegistryKeeperUpdated)
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
		it.Event = new(RegistryRegistryKeeperUpdated)
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
func (it *RegistryRegistryKeeperUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryRegistryKeeperUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryRegistryKeeperUpdated represents a RegistryKeeperUpdated event raised by the Registry contract.
type RegistryRegistryKeeperUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRegistryKeeperUpdated is a free log retrieval operation binding the contract event 0x5c579a7c46def7fd67b8114bb3ebbf46ef84b8d7939b2177abd34565ab7dd971.
//
// Solidity: event RegistryKeeperUpdated()
func (_Registry *RegistryFilterer) FilterRegistryKeeperUpdated(opts *bind.FilterOpts) (*RegistryRegistryKeeperUpdatedIterator, error) {

	logs, sub, err := _Registry.contract.FilterLogs(opts, "RegistryKeeperUpdated")
	if err != nil {
		return nil, err
	}
	return &RegistryRegistryKeeperUpdatedIterator{contract: _Registry.contract, event: "RegistryKeeperUpdated", logs: logs, sub: sub}, nil
}

// WatchRegistryKeeperUpdated is a free log subscription operation binding the contract event 0x5c579a7c46def7fd67b8114bb3ebbf46ef84b8d7939b2177abd34565ab7dd971.
//
// Solidity: event RegistryKeeperUpdated()
func (_Registry *RegistryFilterer) WatchRegistryKeeperUpdated(opts *bind.WatchOpts, sink chan<- *RegistryRegistryKeeperUpdated) (event.Subscription, error) {

	logs, sub, err := _Registry.contract.WatchLogs(opts, "RegistryKeeperUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryRegistryKeeperUpdated)
				if err := _Registry.contract.UnpackLog(event, "RegistryKeeperUpdated", log); err != nil {
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

// ParseRegistryKeeperUpdated is a log parse operation binding the contract event 0x5c579a7c46def7fd67b8114bb3ebbf46ef84b8d7939b2177abd34565ab7dd971.
//
// Solidity: event RegistryKeeperUpdated()
func (_Registry *RegistryFilterer) ParseRegistryKeeperUpdated(log types.Log) (*RegistryRegistryKeeperUpdated, error) {
	event := new(RegistryRegistryKeeperUpdated)
	if err := _Registry.contract.UnpackLog(event, "RegistryKeeperUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
