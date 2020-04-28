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

// KeepRandomBeaconServiceImplV1ABI is the input ABI used to generate the binding from.
const KeepRandomBeaconServiceImplV1ABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"entry\",\"type\":\"uint256\"}],\"name\":\"RelayEntryGenerated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"RelayEntryRequested\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"addOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"baseCallbackGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"entry\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"submitter\",\"type\":\"address\"}],\"name\":\"entryCreated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"entryFeeBreakdown\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"entryVerificationFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dkgContributionFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"groupProfitFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPriceCeiling\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"callbackGas\",\"type\":\"uint256\"}],\"name\":\"entryFeeEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"entry\",\"type\":\"uint256\"}],\"name\":\"executeCallback\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"surplusRecipient\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"fundDkgFeePool\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"fundRequestSubsidyFeePool\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dkgContributionMargin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"registry\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"previousEntry\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operatorContract\",\"type\":\"address\"}],\"name\":\"removeOperatorContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"requestRelayEntry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"callbackContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"callbackGas\",\"type\":\"uint256\"}],\"name\":\"requestRelayEntry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"seed\",\"type\":\"uint256\"}],\"name\":\"selectOperatorContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// KeepRandomBeaconServiceImplV1 is an auto generated Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1 struct {
	KeepRandomBeaconServiceImplV1Caller     // Read-only binding to the contract
	KeepRandomBeaconServiceImplV1Transactor // Write-only binding to the contract
	KeepRandomBeaconServiceImplV1Filterer   // Log filterer for contract events
}

// KeepRandomBeaconServiceImplV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceImplV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceImplV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeepRandomBeaconServiceImplV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeepRandomBeaconServiceImplV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeepRandomBeaconServiceImplV1Session struct {
	Contract     *KeepRandomBeaconServiceImplV1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// KeepRandomBeaconServiceImplV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeepRandomBeaconServiceImplV1CallerSession struct {
	Contract *KeepRandomBeaconServiceImplV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// KeepRandomBeaconServiceImplV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeepRandomBeaconServiceImplV1TransactorSession struct {
	Contract     *KeepRandomBeaconServiceImplV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// KeepRandomBeaconServiceImplV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1Raw struct {
	Contract *KeepRandomBeaconServiceImplV1 // Generic contract binding to access the raw methods on
}

// KeepRandomBeaconServiceImplV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1CallerRaw struct {
	Contract *KeepRandomBeaconServiceImplV1Caller // Generic read-only contract binding to access the raw methods on
}

// KeepRandomBeaconServiceImplV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeepRandomBeaconServiceImplV1TransactorRaw struct {
	Contract *KeepRandomBeaconServiceImplV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewKeepRandomBeaconServiceImplV1 creates a new instance of KeepRandomBeaconServiceImplV1, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceImplV1(address common.Address, backend bind.ContractBackend) (*KeepRandomBeaconServiceImplV1, error) {
	contract, err := bindKeepRandomBeaconServiceImplV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1{KeepRandomBeaconServiceImplV1Caller: KeepRandomBeaconServiceImplV1Caller{contract: contract}, KeepRandomBeaconServiceImplV1Transactor: KeepRandomBeaconServiceImplV1Transactor{contract: contract}, KeepRandomBeaconServiceImplV1Filterer: KeepRandomBeaconServiceImplV1Filterer{contract: contract}}, nil
}

// NewKeepRandomBeaconServiceImplV1Caller creates a new read-only instance of KeepRandomBeaconServiceImplV1, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceImplV1Caller(address common.Address, caller bind.ContractCaller) (*KeepRandomBeaconServiceImplV1Caller, error) {
	contract, err := bindKeepRandomBeaconServiceImplV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1Caller{contract: contract}, nil
}

// NewKeepRandomBeaconServiceImplV1Transactor creates a new write-only instance of KeepRandomBeaconServiceImplV1, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceImplV1Transactor(address common.Address, transactor bind.ContractTransactor) (*KeepRandomBeaconServiceImplV1Transactor, error) {
	contract, err := bindKeepRandomBeaconServiceImplV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1Transactor{contract: contract}, nil
}

// NewKeepRandomBeaconServiceImplV1Filterer creates a new log filterer instance of KeepRandomBeaconServiceImplV1, bound to a specific deployed contract.
func NewKeepRandomBeaconServiceImplV1Filterer(address common.Address, filterer bind.ContractFilterer) (*KeepRandomBeaconServiceImplV1Filterer, error) {
	contract, err := bindKeepRandomBeaconServiceImplV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1Filterer{contract: contract}, nil
}

// bindKeepRandomBeaconServiceImplV1 binds a generic wrapper to an already deployed contract.
func bindKeepRandomBeaconServiceImplV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeepRandomBeaconServiceImplV1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconServiceImplV1.Contract.KeepRandomBeaconServiceImplV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.KeepRandomBeaconServiceImplV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.KeepRandomBeaconServiceImplV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KeepRandomBeaconServiceImplV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.contract.Transact(opts, method, params...)
}

// BaseCallbackGas is a free data retrieval call binding the contract method 0xd09dd574.
//
// Solidity: function baseCallbackGas() constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) BaseCallbackGas(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "baseCallbackGas")
	return *ret0, err
}

// BaseCallbackGas is a free data retrieval call binding the contract method 0xd09dd574.
//
// Solidity: function baseCallbackGas() constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) BaseCallbackGas() (*big.Int, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.BaseCallbackGas(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// BaseCallbackGas is a free data retrieval call binding the contract method 0xd09dd574.
//
// Solidity: function baseCallbackGas() constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) BaseCallbackGas() (*big.Int, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.BaseCallbackGas(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// EntryFeeBreakdown is a free data retrieval call binding the contract method 0xdede8e95.
//
// Solidity: function entryFeeBreakdown() constant returns(uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee, uint256 gasPriceCeiling)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) EntryFeeBreakdown(opts *bind.CallOpts) (struct {
	EntryVerificationFee *big.Int
	DkgContributionFee   *big.Int
	GroupProfitFee       *big.Int
	GasPriceCeiling      *big.Int
}, error) {
	ret := new(struct {
		EntryVerificationFee *big.Int
		DkgContributionFee   *big.Int
		GroupProfitFee       *big.Int
		GasPriceCeiling      *big.Int
	})
	out := ret
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "entryFeeBreakdown")
	return *ret, err
}

// EntryFeeBreakdown is a free data retrieval call binding the contract method 0xdede8e95.
//
// Solidity: function entryFeeBreakdown() constant returns(uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee, uint256 gasPriceCeiling)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) EntryFeeBreakdown() (struct {
	EntryVerificationFee *big.Int
	DkgContributionFee   *big.Int
	GroupProfitFee       *big.Int
	GasPriceCeiling      *big.Int
}, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryFeeBreakdown(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// EntryFeeBreakdown is a free data retrieval call binding the contract method 0xdede8e95.
//
// Solidity: function entryFeeBreakdown() constant returns(uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee, uint256 gasPriceCeiling)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) EntryFeeBreakdown() (struct {
	EntryVerificationFee *big.Int
	DkgContributionFee   *big.Int
	GroupProfitFee       *big.Int
	GasPriceCeiling      *big.Int
}, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryFeeBreakdown(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) EntryFeeEstimate(opts *bind.CallOpts, callbackGas *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "entryFeeEstimate", callbackGas)
	return *ret0, err
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) EntryFeeEstimate(callbackGas *big.Int) (*big.Int, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryFeeEstimate(&_KeepRandomBeaconServiceImplV1.CallOpts, callbackGas)
}

// EntryFeeEstimate is a free data retrieval call binding the contract method 0xd13f1391.
//
// Solidity: function entryFeeEstimate(uint256 callbackGas) constant returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) EntryFeeEstimate(callbackGas *big.Int) (*big.Int, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryFeeEstimate(&_KeepRandomBeaconServiceImplV1.CallOpts, callbackGas)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() constant returns(bool)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) Initialized(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "initialized")
	return *ret0, err
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() constant returns(bool)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) Initialized() (bool, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Initialized(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() constant returns(bool)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) Initialized() (bool, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Initialized(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// PreviousEntry is a free data retrieval call binding the contract method 0xac6136db.
//
// Solidity: function previousEntry() constant returns(bytes)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) PreviousEntry(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "previousEntry")
	return *ret0, err
}

// PreviousEntry is a free data retrieval call binding the contract method 0xac6136db.
//
// Solidity: function previousEntry() constant returns(bytes)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) PreviousEntry() ([]byte, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.PreviousEntry(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// PreviousEntry is a free data retrieval call binding the contract method 0xac6136db.
//
// Solidity: function previousEntry() constant returns(bytes)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) PreviousEntry() ([]byte, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.PreviousEntry(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// SelectOperatorContract is a free data retrieval call binding the contract method 0xefc4971f.
//
// Solidity: function selectOperatorContract(uint256 seed) constant returns(address)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) SelectOperatorContract(opts *bind.CallOpts, seed *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "selectOperatorContract", seed)
	return *ret0, err
}

// SelectOperatorContract is a free data retrieval call binding the contract method 0xefc4971f.
//
// Solidity: function selectOperatorContract(uint256 seed) constant returns(address)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) SelectOperatorContract(seed *big.Int) (common.Address, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.SelectOperatorContract(&_KeepRandomBeaconServiceImplV1.CallOpts, seed)
}

// SelectOperatorContract is a free data retrieval call binding the contract method 0xefc4971f.
//
// Solidity: function selectOperatorContract(uint256 seed) constant returns(address)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) SelectOperatorContract(seed *big.Int) (common.Address, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.SelectOperatorContract(&_KeepRandomBeaconServiceImplV1.CallOpts, seed)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Caller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KeepRandomBeaconServiceImplV1.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) Version() (string, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Version(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1CallerSession) Version() (string, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Version(&_KeepRandomBeaconServiceImplV1.CallOpts)
}

// AddOperatorContract is a paid mutator transaction binding the contract method 0x60e07ffd.
//
// Solidity: function addOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) AddOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "addOperatorContract", operatorContract)
}

// AddOperatorContract is a paid mutator transaction binding the contract method 0x60e07ffd.
//
// Solidity: function addOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) AddOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.AddOperatorContract(&_KeepRandomBeaconServiceImplV1.TransactOpts, operatorContract)
}

// AddOperatorContract is a paid mutator transaction binding the contract method 0x60e07ffd.
//
// Solidity: function addOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) AddOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.AddOperatorContract(&_KeepRandomBeaconServiceImplV1.TransactOpts, operatorContract)
}

// EntryCreated is a paid mutator transaction binding the contract method 0xef7284e3.
//
// Solidity: function entryCreated(uint256 requestId, bytes entry, address submitter) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) EntryCreated(opts *bind.TransactOpts, requestId *big.Int, entry []byte, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "entryCreated", requestId, entry, submitter)
}

// EntryCreated is a paid mutator transaction binding the contract method 0xef7284e3.
//
// Solidity: function entryCreated(uint256 requestId, bytes entry, address submitter) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) EntryCreated(requestId *big.Int, entry []byte, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryCreated(&_KeepRandomBeaconServiceImplV1.TransactOpts, requestId, entry, submitter)
}

// EntryCreated is a paid mutator transaction binding the contract method 0xef7284e3.
//
// Solidity: function entryCreated(uint256 requestId, bytes entry, address submitter) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) EntryCreated(requestId *big.Int, entry []byte, submitter common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.EntryCreated(&_KeepRandomBeaconServiceImplV1.TransactOpts, requestId, entry, submitter)
}

// ExecuteCallback is a paid mutator transaction binding the contract method 0xfc3fcec7.
//
// Solidity: function executeCallback(uint256 requestId, uint256 entry) returns(address surplusRecipient)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) ExecuteCallback(opts *bind.TransactOpts, requestId *big.Int, entry *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "executeCallback", requestId, entry)
}

// ExecuteCallback is a paid mutator transaction binding the contract method 0xfc3fcec7.
//
// Solidity: function executeCallback(uint256 requestId, uint256 entry) returns(address surplusRecipient)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) ExecuteCallback(requestId *big.Int, entry *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.ExecuteCallback(&_KeepRandomBeaconServiceImplV1.TransactOpts, requestId, entry)
}

// ExecuteCallback is a paid mutator transaction binding the contract method 0xfc3fcec7.
//
// Solidity: function executeCallback(uint256 requestId, uint256 entry) returns(address surplusRecipient)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) ExecuteCallback(requestId *big.Int, entry *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.ExecuteCallback(&_KeepRandomBeaconServiceImplV1.TransactOpts, requestId, entry)
}

// FundDkgFeePool is a paid mutator transaction binding the contract method 0x4611b648.
//
// Solidity: function fundDkgFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) FundDkgFeePool(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "fundDkgFeePool")
}

// FundDkgFeePool is a paid mutator transaction binding the contract method 0x4611b648.
//
// Solidity: function fundDkgFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) FundDkgFeePool() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.FundDkgFeePool(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// FundDkgFeePool is a paid mutator transaction binding the contract method 0x4611b648.
//
// Solidity: function fundDkgFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) FundDkgFeePool() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.FundDkgFeePool(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// FundRequestSubsidyFeePool is a paid mutator transaction binding the contract method 0x11e816ee.
//
// Solidity: function fundRequestSubsidyFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) FundRequestSubsidyFeePool(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "fundRequestSubsidyFeePool")
}

// FundRequestSubsidyFeePool is a paid mutator transaction binding the contract method 0x11e816ee.
//
// Solidity: function fundRequestSubsidyFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) FundRequestSubsidyFeePool() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.FundRequestSubsidyFeePool(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// FundRequestSubsidyFeePool is a paid mutator transaction binding the contract method 0x11e816ee.
//
// Solidity: function fundRequestSubsidyFeePool() returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) FundRequestSubsidyFeePool() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.FundRequestSubsidyFeePool(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 dkgContributionMargin, address registry) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) Initialize(opts *bind.TransactOpts, dkgContributionMargin *big.Int, registry common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "initialize", dkgContributionMargin, registry)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 dkgContributionMargin, address registry) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) Initialize(dkgContributionMargin *big.Int, registry common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Initialize(&_KeepRandomBeaconServiceImplV1.TransactOpts, dkgContributionMargin, registry)
}

// Initialize is a paid mutator transaction binding the contract method 0xda35a26f.
//
// Solidity: function initialize(uint256 dkgContributionMargin, address registry) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) Initialize(dkgContributionMargin *big.Int, registry common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.Initialize(&_KeepRandomBeaconServiceImplV1.TransactOpts, dkgContributionMargin, registry)
}

// RemoveOperatorContract is a paid mutator transaction binding the contract method 0xe58990c5.
//
// Solidity: function removeOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) RemoveOperatorContract(opts *bind.TransactOpts, operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "removeOperatorContract", operatorContract)
}

// RemoveOperatorContract is a paid mutator transaction binding the contract method 0xe58990c5.
//
// Solidity: function removeOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) RemoveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RemoveOperatorContract(&_KeepRandomBeaconServiceImplV1.TransactOpts, operatorContract)
}

// RemoveOperatorContract is a paid mutator transaction binding the contract method 0xe58990c5.
//
// Solidity: function removeOperatorContract(address operatorContract) returns()
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) RemoveOperatorContract(operatorContract common.Address) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RemoveOperatorContract(&_KeepRandomBeaconServiceImplV1.TransactOpts, operatorContract)
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) RequestRelayEntry(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "requestRelayEntry")
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) RequestRelayEntry() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RequestRelayEntry(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// RequestRelayEntry is a paid mutator transaction binding the contract method 0x2d53fe8b.
//
// Solidity: function requestRelayEntry() returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) RequestRelayEntry() (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RequestRelayEntry(&_KeepRandomBeaconServiceImplV1.TransactOpts)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Transactor) RequestRelayEntry0(opts *bind.TransactOpts, callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.contract.Transact(opts, "requestRelayEntry0", callbackContract, callbackGas)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Session) RequestRelayEntry0(callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RequestRelayEntry0(&_KeepRandomBeaconServiceImplV1.TransactOpts, callbackContract, callbackGas)
}

// RequestRelayEntry0 is a paid mutator transaction binding the contract method 0xe1f95890.
//
// Solidity: function requestRelayEntry(address callbackContract, uint256 callbackGas) returns(uint256)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1TransactorSession) RequestRelayEntry0(callbackContract common.Address, callbackGas *big.Int) (*types.Transaction, error) {
	return _KeepRandomBeaconServiceImplV1.Contract.RequestRelayEntry0(&_KeepRandomBeaconServiceImplV1.TransactOpts, callbackContract, callbackGas)
}

// KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator is returned from FilterRelayEntryGenerated and is used to iterate over the raw logs and unpacked data for RelayEntryGenerated events raised by the KeepRandomBeaconServiceImplV1 contract.
type KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator struct {
	Event *KeepRandomBeaconServiceImplV1RelayEntryGenerated // Event containing the contract specifics and raw log

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
func (it *KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconServiceImplV1RelayEntryGenerated)
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
		it.Event = new(KeepRandomBeaconServiceImplV1RelayEntryGenerated)
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
func (it *KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconServiceImplV1RelayEntryGenerated represents a RelayEntryGenerated event raised by the KeepRandomBeaconServiceImplV1 contract.
type KeepRandomBeaconServiceImplV1RelayEntryGenerated struct {
	RequestId *big.Int
	Entry     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryGenerated is a free log retrieval operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) FilterRelayEntryGenerated(opts *bind.FilterOpts) (*KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator, error) {

	logs, sub, err := _KeepRandomBeaconServiceImplV1.contract.FilterLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1RelayEntryGeneratedIterator{contract: _KeepRandomBeaconServiceImplV1.contract, event: "RelayEntryGenerated", logs: logs, sub: sub}, nil
}

// WatchRelayEntryGenerated is a free log subscription operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) WatchRelayEntryGenerated(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconServiceImplV1RelayEntryGenerated) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconServiceImplV1.contract.WatchLogs(opts, "RelayEntryGenerated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconServiceImplV1RelayEntryGenerated)
				if err := _KeepRandomBeaconServiceImplV1.contract.UnpackLog(event, "RelayEntryGenerated", log); err != nil {
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

// ParseRelayEntryGenerated is a log parse operation binding the contract event 0xddb7473cb5eaf3aaa2cf09ea9d573d4876aac471614f99773f32a4b5994ea7e1.
//
// Solidity: event RelayEntryGenerated(uint256 requestId, uint256 entry)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) ParseRelayEntryGenerated(log types.Log) (*KeepRandomBeaconServiceImplV1RelayEntryGenerated, error) {
	event := new(KeepRandomBeaconServiceImplV1RelayEntryGenerated)
	if err := _KeepRandomBeaconServiceImplV1.contract.UnpackLog(event, "RelayEntryGenerated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator is returned from FilterRelayEntryRequested and is used to iterate over the raw logs and unpacked data for RelayEntryRequested events raised by the KeepRandomBeaconServiceImplV1 contract.
type KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator struct {
	Event *KeepRandomBeaconServiceImplV1RelayEntryRequested // Event containing the contract specifics and raw log

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
func (it *KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KeepRandomBeaconServiceImplV1RelayEntryRequested)
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
		it.Event = new(KeepRandomBeaconServiceImplV1RelayEntryRequested)
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
func (it *KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KeepRandomBeaconServiceImplV1RelayEntryRequested represents a RelayEntryRequested event raised by the KeepRandomBeaconServiceImplV1 contract.
type KeepRandomBeaconServiceImplV1RelayEntryRequested struct {
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRelayEntryRequested is a free log retrieval operation binding the contract event 0xc65bedbc1f2ce56540b8b4e1b7b41261bef9dc7aee32cb3b6806a70c5f4828fa.
//
// Solidity: event RelayEntryRequested(uint256 requestId)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) FilterRelayEntryRequested(opts *bind.FilterOpts) (*KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator, error) {

	logs, sub, err := _KeepRandomBeaconServiceImplV1.contract.FilterLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return &KeepRandomBeaconServiceImplV1RelayEntryRequestedIterator{contract: _KeepRandomBeaconServiceImplV1.contract, event: "RelayEntryRequested", logs: logs, sub: sub}, nil
}

// WatchRelayEntryRequested is a free log subscription operation binding the contract event 0xc65bedbc1f2ce56540b8b4e1b7b41261bef9dc7aee32cb3b6806a70c5f4828fa.
//
// Solidity: event RelayEntryRequested(uint256 requestId)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) WatchRelayEntryRequested(opts *bind.WatchOpts, sink chan<- *KeepRandomBeaconServiceImplV1RelayEntryRequested) (event.Subscription, error) {

	logs, sub, err := _KeepRandomBeaconServiceImplV1.contract.WatchLogs(opts, "RelayEntryRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KeepRandomBeaconServiceImplV1RelayEntryRequested)
				if err := _KeepRandomBeaconServiceImplV1.contract.UnpackLog(event, "RelayEntryRequested", log); err != nil {
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

// ParseRelayEntryRequested is a log parse operation binding the contract event 0xc65bedbc1f2ce56540b8b4e1b7b41261bef9dc7aee32cb3b6806a70c5f4828fa.
//
// Solidity: event RelayEntryRequested(uint256 requestId)
func (_KeepRandomBeaconServiceImplV1 *KeepRandomBeaconServiceImplV1Filterer) ParseRelayEntryRequested(log types.Log) (*KeepRandomBeaconServiceImplV1RelayEntryRequested, error) {
	event := new(KeepRandomBeaconServiceImplV1RelayEntryRequested)
	if err := _KeepRandomBeaconServiceImplV1.contract.UnpackLog(event, "RelayEntryRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}
