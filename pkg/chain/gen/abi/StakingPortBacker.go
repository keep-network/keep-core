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

// StakingPortBackerABI is the input ABI used to generate the binding from.
const StakingPortBackerABI = "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"_keepToken\",\"type\":\"address\"},{\"internalType\":\"contractTokenGrant\",\"name\":\"_tokenGrant\",\"type\":\"address\"},{\"internalType\":\"contractStakeDelegatable\",\"name\":\"_oldStakingContract\",\"type\":\"address\"},{\"internalType\":\"contractTokenStaking\",\"name\":\"_newStakingContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"StakeCopied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"StakePaidBack\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensWithdrawn\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"allowOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"operators\",\"type\":\"address[]\"}],\"name\":\"allowOperators\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedOperators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"copiedStakes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"paidBack\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"copyStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"forceUndelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keepToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxAllowedBackingDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newStakingContract\",\"outputs\":[{\"internalType\":\"contractTokenStaking\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"oldStakingContract\",\"outputs\":[{\"internalType\":\"contractStakeDelegatable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"recoverStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StakingPortBacker is an auto generated Go binding around an Ethereum contract.
type StakingPortBacker struct {
	StakingPortBackerCaller     // Read-only binding to the contract
	StakingPortBackerTransactor // Write-only binding to the contract
	StakingPortBackerFilterer   // Log filterer for contract events
}

// StakingPortBackerCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingPortBackerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPortBackerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingPortBackerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPortBackerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingPortBackerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPortBackerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingPortBackerSession struct {
	Contract     *StakingPortBacker // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingPortBackerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingPortBackerCallerSession struct {
	Contract *StakingPortBackerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// StakingPortBackerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingPortBackerTransactorSession struct {
	Contract     *StakingPortBackerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// StakingPortBackerRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingPortBackerRaw struct {
	Contract *StakingPortBacker // Generic contract binding to access the raw methods on
}

// StakingPortBackerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingPortBackerCallerRaw struct {
	Contract *StakingPortBackerCaller // Generic read-only contract binding to access the raw methods on
}

// StakingPortBackerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingPortBackerTransactorRaw struct {
	Contract *StakingPortBackerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingPortBacker creates a new instance of StakingPortBacker, bound to a specific deployed contract.
func NewStakingPortBacker(address common.Address, backend bind.ContractBackend) (*StakingPortBacker, error) {
	contract, err := bindStakingPortBacker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingPortBacker{StakingPortBackerCaller: StakingPortBackerCaller{contract: contract}, StakingPortBackerTransactor: StakingPortBackerTransactor{contract: contract}, StakingPortBackerFilterer: StakingPortBackerFilterer{contract: contract}}, nil
}

// NewStakingPortBackerCaller creates a new read-only instance of StakingPortBacker, bound to a specific deployed contract.
func NewStakingPortBackerCaller(address common.Address, caller bind.ContractCaller) (*StakingPortBackerCaller, error) {
	contract, err := bindStakingPortBacker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerCaller{contract: contract}, nil
}

// NewStakingPortBackerTransactor creates a new write-only instance of StakingPortBacker, bound to a specific deployed contract.
func NewStakingPortBackerTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingPortBackerTransactor, error) {
	contract, err := bindStakingPortBacker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerTransactor{contract: contract}, nil
}

// NewStakingPortBackerFilterer creates a new log filterer instance of StakingPortBacker, bound to a specific deployed contract.
func NewStakingPortBackerFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingPortBackerFilterer, error) {
	contract, err := bindStakingPortBacker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerFilterer{contract: contract}, nil
}

// bindStakingPortBacker binds a generic wrapper to an already deployed contract.
func bindStakingPortBacker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingPortBackerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPortBacker *StakingPortBackerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakingPortBacker.Contract.StakingPortBackerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPortBacker *StakingPortBackerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.StakingPortBackerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPortBacker *StakingPortBackerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.StakingPortBackerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPortBacker *StakingPortBackerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakingPortBacker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPortBacker *StakingPortBackerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPortBacker *StakingPortBackerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.contract.Transact(opts, method, params...)
}

// AllowedOperators is a free data retrieval call binding the contract method 0xa6674f88.
//
// Solidity: function allowedOperators(address ) constant returns(bool)
func (_StakingPortBacker *StakingPortBackerCaller) AllowedOperators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "allowedOperators", arg0)
	return *ret0, err
}

// AllowedOperators is a free data retrieval call binding the contract method 0xa6674f88.
//
// Solidity: function allowedOperators(address ) constant returns(bool)
func (_StakingPortBacker *StakingPortBackerSession) AllowedOperators(arg0 common.Address) (bool, error) {
	return _StakingPortBacker.Contract.AllowedOperators(&_StakingPortBacker.CallOpts, arg0)
}

// AllowedOperators is a free data retrieval call binding the contract method 0xa6674f88.
//
// Solidity: function allowedOperators(address ) constant returns(bool)
func (_StakingPortBacker *StakingPortBackerCallerSession) AllowedOperators(arg0 common.Address) (bool, error) {
	return _StakingPortBacker.Contract.AllowedOperators(&_StakingPortBacker.CallOpts, arg0)
}

// CopiedStakes is a free data retrieval call binding the contract method 0x4d85d9cb.
//
// Solidity: function copiedStakes(address ) constant returns(address owner, uint256 amount, uint256 timestamp, bool paidBack)
func (_StakingPortBacker *StakingPortBackerCaller) CopiedStakes(opts *bind.CallOpts, arg0 common.Address) (struct {
	Owner     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	PaidBack  bool
}, error) {
	ret := new(struct {
		Owner     common.Address
		Amount    *big.Int
		Timestamp *big.Int
		PaidBack  bool
	})
	out := ret
	err := _StakingPortBacker.contract.Call(opts, out, "copiedStakes", arg0)
	return *ret, err
}

// CopiedStakes is a free data retrieval call binding the contract method 0x4d85d9cb.
//
// Solidity: function copiedStakes(address ) constant returns(address owner, uint256 amount, uint256 timestamp, bool paidBack)
func (_StakingPortBacker *StakingPortBackerSession) CopiedStakes(arg0 common.Address) (struct {
	Owner     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	PaidBack  bool
}, error) {
	return _StakingPortBacker.Contract.CopiedStakes(&_StakingPortBacker.CallOpts, arg0)
}

// CopiedStakes is a free data retrieval call binding the contract method 0x4d85d9cb.
//
// Solidity: function copiedStakes(address ) constant returns(address owner, uint256 amount, uint256 timestamp, bool paidBack)
func (_StakingPortBacker *StakingPortBackerCallerSession) CopiedStakes(arg0 common.Address) (struct {
	Owner     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	PaidBack  bool
}, error) {
	return _StakingPortBacker.Contract.CopiedStakes(&_StakingPortBacker.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_StakingPortBacker *StakingPortBackerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_StakingPortBacker *StakingPortBackerSession) IsOwner() (bool, error) {
	return _StakingPortBacker.Contract.IsOwner(&_StakingPortBacker.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_StakingPortBacker *StakingPortBackerCallerSession) IsOwner() (bool, error) {
	return _StakingPortBacker.Contract.IsOwner(&_StakingPortBacker.CallOpts)
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCaller) KeepToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "keepToken")
	return *ret0, err
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_StakingPortBacker *StakingPortBackerSession) KeepToken() (common.Address, error) {
	return _StakingPortBacker.Contract.KeepToken(&_StakingPortBacker.CallOpts)
}

// KeepToken is a free data retrieval call binding the contract method 0x467c5adf.
//
// Solidity: function keepToken() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCallerSession) KeepToken() (common.Address, error) {
	return _StakingPortBacker.Contract.KeepToken(&_StakingPortBacker.CallOpts)
}

// MaxAllowedBackingDuration is a free data retrieval call binding the contract method 0xd97d5953.
//
// Solidity: function maxAllowedBackingDuration() constant returns(uint256)
func (_StakingPortBacker *StakingPortBackerCaller) MaxAllowedBackingDuration(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "maxAllowedBackingDuration")
	return *ret0, err
}

// MaxAllowedBackingDuration is a free data retrieval call binding the contract method 0xd97d5953.
//
// Solidity: function maxAllowedBackingDuration() constant returns(uint256)
func (_StakingPortBacker *StakingPortBackerSession) MaxAllowedBackingDuration() (*big.Int, error) {
	return _StakingPortBacker.Contract.MaxAllowedBackingDuration(&_StakingPortBacker.CallOpts)
}

// MaxAllowedBackingDuration is a free data retrieval call binding the contract method 0xd97d5953.
//
// Solidity: function maxAllowedBackingDuration() constant returns(uint256)
func (_StakingPortBacker *StakingPortBackerCallerSession) MaxAllowedBackingDuration() (*big.Int, error) {
	return _StakingPortBacker.Contract.MaxAllowedBackingDuration(&_StakingPortBacker.CallOpts)
}

// NewStakingContract is a free data retrieval call binding the contract method 0xae81dfe4.
//
// Solidity: function newStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCaller) NewStakingContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "newStakingContract")
	return *ret0, err
}

// NewStakingContract is a free data retrieval call binding the contract method 0xae81dfe4.
//
// Solidity: function newStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerSession) NewStakingContract() (common.Address, error) {
	return _StakingPortBacker.Contract.NewStakingContract(&_StakingPortBacker.CallOpts)
}

// NewStakingContract is a free data retrieval call binding the contract method 0xae81dfe4.
//
// Solidity: function newStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCallerSession) NewStakingContract() (common.Address, error) {
	return _StakingPortBacker.Contract.NewStakingContract(&_StakingPortBacker.CallOpts)
}

// OldStakingContract is a free data retrieval call binding the contract method 0xe59dd90d.
//
// Solidity: function oldStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCaller) OldStakingContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "oldStakingContract")
	return *ret0, err
}

// OldStakingContract is a free data retrieval call binding the contract method 0xe59dd90d.
//
// Solidity: function oldStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerSession) OldStakingContract() (common.Address, error) {
	return _StakingPortBacker.Contract.OldStakingContract(&_StakingPortBacker.CallOpts)
}

// OldStakingContract is a free data retrieval call binding the contract method 0xe59dd90d.
//
// Solidity: function oldStakingContract() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCallerSession) OldStakingContract() (common.Address, error) {
	return _StakingPortBacker.Contract.OldStakingContract(&_StakingPortBacker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _StakingPortBacker.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StakingPortBacker *StakingPortBackerSession) Owner() (common.Address, error) {
	return _StakingPortBacker.Contract.Owner(&_StakingPortBacker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_StakingPortBacker *StakingPortBackerCallerSession) Owner() (common.Address, error) {
	return _StakingPortBacker.Contract.Owner(&_StakingPortBacker.CallOpts)
}

// AllowOperator is a paid mutator transaction binding the contract method 0x1ba9a458.
//
// Solidity: function allowOperator(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) AllowOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "allowOperator", operator)
}

// AllowOperator is a paid mutator transaction binding the contract method 0x1ba9a458.
//
// Solidity: function allowOperator(address operator) returns()
func (_StakingPortBacker *StakingPortBackerSession) AllowOperator(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.AllowOperator(&_StakingPortBacker.TransactOpts, operator)
}

// AllowOperator is a paid mutator transaction binding the contract method 0x1ba9a458.
//
// Solidity: function allowOperator(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) AllowOperator(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.AllowOperator(&_StakingPortBacker.TransactOpts, operator)
}

// AllowOperators is a paid mutator transaction binding the contract method 0x432de9c8.
//
// Solidity: function allowOperators(address[] operators) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) AllowOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "allowOperators", operators)
}

// AllowOperators is a paid mutator transaction binding the contract method 0x432de9c8.
//
// Solidity: function allowOperators(address[] operators) returns()
func (_StakingPortBacker *StakingPortBackerSession) AllowOperators(operators []common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.AllowOperators(&_StakingPortBacker.TransactOpts, operators)
}

// AllowOperators is a paid mutator transaction binding the contract method 0x432de9c8.
//
// Solidity: function allowOperators(address[] operators) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) AllowOperators(operators []common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.AllowOperators(&_StakingPortBacker.TransactOpts, operators)
}

// CopyStake is a paid mutator transaction binding the contract method 0x67d847e6.
//
// Solidity: function copyStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) CopyStake(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "copyStake", operator)
}

// CopyStake is a paid mutator transaction binding the contract method 0x67d847e6.
//
// Solidity: function copyStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerSession) CopyStake(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.CopyStake(&_StakingPortBacker.TransactOpts, operator)
}

// CopyStake is a paid mutator transaction binding the contract method 0x67d847e6.
//
// Solidity: function copyStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) CopyStake(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.CopyStake(&_StakingPortBacker.TransactOpts, operator)
}

// ForceUndelegate is a paid mutator transaction binding the contract method 0x391cc9f6.
//
// Solidity: function forceUndelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) ForceUndelegate(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "forceUndelegate", operator)
}

// ForceUndelegate is a paid mutator transaction binding the contract method 0x391cc9f6.
//
// Solidity: function forceUndelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerSession) ForceUndelegate(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.ForceUndelegate(&_StakingPortBacker.TransactOpts, operator)
}

// ForceUndelegate is a paid mutator transaction binding the contract method 0x391cc9f6.
//
// Solidity: function forceUndelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) ForceUndelegate(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.ForceUndelegate(&_StakingPortBacker.TransactOpts, operator)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) ReceiveApproval(opts *bind.TransactOpts, from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "receiveApproval", from, value, token, extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_StakingPortBacker *StakingPortBackerSession) ReceiveApproval(from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.ReceiveApproval(&_StakingPortBacker.TransactOpts, from, value, token, extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address from, uint256 value, address token, bytes extraData) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) ReceiveApproval(from common.Address, value *big.Int, token common.Address, extraData []byte) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.ReceiveApproval(&_StakingPortBacker.TransactOpts, from, value, token, extraData)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) RecoverStake(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "recoverStake", operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerSession) RecoverStake(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.RecoverStake(&_StakingPortBacker.TransactOpts, operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) RecoverStake(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.RecoverStake(&_StakingPortBacker.TransactOpts, operator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPortBacker *StakingPortBackerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPortBacker *StakingPortBackerSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingPortBacker.Contract.RenounceOwnership(&_StakingPortBacker.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingPortBacker.Contract.RenounceOwnership(&_StakingPortBacker.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPortBacker *StakingPortBackerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.TransferOwnership(&_StakingPortBacker.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.TransferOwnership(&_StakingPortBacker.TransactOpts, newOwner)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) Undelegate(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "undelegate", operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerSession) Undelegate(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.Undelegate(&_StakingPortBacker.TransactOpts, operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address operator) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) Undelegate(operator common.Address) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.Undelegate(&_StakingPortBacker.TransactOpts, operator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_StakingPortBacker *StakingPortBackerTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _StakingPortBacker.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_StakingPortBacker *StakingPortBackerSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.Withdraw(&_StakingPortBacker.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_StakingPortBacker *StakingPortBackerTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _StakingPortBacker.Contract.Withdraw(&_StakingPortBacker.TransactOpts, amount)
}

// StakingPortBackerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the StakingPortBacker contract.
type StakingPortBackerOwnershipTransferredIterator struct {
	Event *StakingPortBackerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakingPortBackerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPortBackerOwnershipTransferred)
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
		it.Event = new(StakingPortBackerOwnershipTransferred)
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
func (it *StakingPortBackerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPortBackerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPortBackerOwnershipTransferred represents a OwnershipTransferred event raised by the StakingPortBacker contract.
type StakingPortBackerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingPortBacker *StakingPortBackerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakingPortBackerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingPortBacker.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerOwnershipTransferredIterator{contract: _StakingPortBacker.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingPortBacker *StakingPortBackerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingPortBackerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingPortBacker.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPortBackerOwnershipTransferred)
				if err := _StakingPortBacker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_StakingPortBacker *StakingPortBackerFilterer) ParseOwnershipTransferred(log types.Log) (*StakingPortBackerOwnershipTransferred, error) {
	event := new(StakingPortBackerOwnershipTransferred)
	if err := _StakingPortBacker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingPortBackerStakeCopiedIterator is returned from FilterStakeCopied and is used to iterate over the raw logs and unpacked data for StakeCopied events raised by the StakingPortBacker contract.
type StakingPortBackerStakeCopiedIterator struct {
	Event *StakingPortBackerStakeCopied // Event containing the contract specifics and raw log

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
func (it *StakingPortBackerStakeCopiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPortBackerStakeCopied)
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
		it.Event = new(StakingPortBackerStakeCopied)
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
func (it *StakingPortBackerStakeCopiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPortBackerStakeCopiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPortBackerStakeCopied represents a StakeCopied event raised by the StakingPortBacker contract.
type StakingPortBackerStakeCopied struct {
	Owner    common.Address
	Operator common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakeCopied is a free log retrieval operation binding the contract event 0x7801a97571720a1fb25f80e759ed0db445c94c7d267e1158af2574de0b0b5ef5.
//
// Solidity: event StakeCopied(address indexed owner, address indexed operator, uint256 value)
func (_StakingPortBacker *StakingPortBackerFilterer) FilterStakeCopied(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*StakingPortBackerStakeCopiedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakingPortBacker.contract.FilterLogs(opts, "StakeCopied", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerStakeCopiedIterator{contract: _StakingPortBacker.contract, event: "StakeCopied", logs: logs, sub: sub}, nil
}

// WatchStakeCopied is a free log subscription operation binding the contract event 0x7801a97571720a1fb25f80e759ed0db445c94c7d267e1158af2574de0b0b5ef5.
//
// Solidity: event StakeCopied(address indexed owner, address indexed operator, uint256 value)
func (_StakingPortBacker *StakingPortBackerFilterer) WatchStakeCopied(opts *bind.WatchOpts, sink chan<- *StakingPortBackerStakeCopied, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakingPortBacker.contract.WatchLogs(opts, "StakeCopied", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPortBackerStakeCopied)
				if err := _StakingPortBacker.contract.UnpackLog(event, "StakeCopied", log); err != nil {
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

// ParseStakeCopied is a log parse operation binding the contract event 0x7801a97571720a1fb25f80e759ed0db445c94c7d267e1158af2574de0b0b5ef5.
//
// Solidity: event StakeCopied(address indexed owner, address indexed operator, uint256 value)
func (_StakingPortBacker *StakingPortBackerFilterer) ParseStakeCopied(log types.Log) (*StakingPortBackerStakeCopied, error) {
	event := new(StakingPortBackerStakeCopied)
	if err := _StakingPortBacker.contract.UnpackLog(event, "StakeCopied", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingPortBackerStakePaidBackIterator is returned from FilterStakePaidBack and is used to iterate over the raw logs and unpacked data for StakePaidBack events raised by the StakingPortBacker contract.
type StakingPortBackerStakePaidBackIterator struct {
	Event *StakingPortBackerStakePaidBack // Event containing the contract specifics and raw log

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
func (it *StakingPortBackerStakePaidBackIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPortBackerStakePaidBack)
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
		it.Event = new(StakingPortBackerStakePaidBack)
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
func (it *StakingPortBackerStakePaidBackIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPortBackerStakePaidBackIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPortBackerStakePaidBack represents a StakePaidBack event raised by the StakingPortBacker contract.
type StakingPortBackerStakePaidBack struct {
	Owner    common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakePaidBack is a free log retrieval operation binding the contract event 0x76934ffc7c84991135ccda0b52fddf17b1fe2db343a49ca2863baded3deb328d.
//
// Solidity: event StakePaidBack(address indexed owner, address indexed operator)
func (_StakingPortBacker *StakingPortBackerFilterer) FilterStakePaidBack(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*StakingPortBackerStakePaidBackIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakingPortBacker.contract.FilterLogs(opts, "StakePaidBack", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerStakePaidBackIterator{contract: _StakingPortBacker.contract, event: "StakePaidBack", logs: logs, sub: sub}, nil
}

// WatchStakePaidBack is a free log subscription operation binding the contract event 0x76934ffc7c84991135ccda0b52fddf17b1fe2db343a49ca2863baded3deb328d.
//
// Solidity: event StakePaidBack(address indexed owner, address indexed operator)
func (_StakingPortBacker *StakingPortBackerFilterer) WatchStakePaidBack(opts *bind.WatchOpts, sink chan<- *StakingPortBackerStakePaidBack, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakingPortBacker.contract.WatchLogs(opts, "StakePaidBack", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPortBackerStakePaidBack)
				if err := _StakingPortBacker.contract.UnpackLog(event, "StakePaidBack", log); err != nil {
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

// ParseStakePaidBack is a log parse operation binding the contract event 0x76934ffc7c84991135ccda0b52fddf17b1fe2db343a49ca2863baded3deb328d.
//
// Solidity: event StakePaidBack(address indexed owner, address indexed operator)
func (_StakingPortBacker *StakingPortBackerFilterer) ParseStakePaidBack(log types.Log) (*StakingPortBackerStakePaidBack, error) {
	event := new(StakingPortBackerStakePaidBack)
	if err := _StakingPortBacker.contract.UnpackLog(event, "StakePaidBack", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingPortBackerTokensWithdrawnIterator is returned from FilterTokensWithdrawn and is used to iterate over the raw logs and unpacked data for TokensWithdrawn events raised by the StakingPortBacker contract.
type StakingPortBackerTokensWithdrawnIterator struct {
	Event *StakingPortBackerTokensWithdrawn // Event containing the contract specifics and raw log

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
func (it *StakingPortBackerTokensWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPortBackerTokensWithdrawn)
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
		it.Event = new(StakingPortBackerTokensWithdrawn)
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
func (it *StakingPortBackerTokensWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPortBackerTokensWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPortBackerTokensWithdrawn represents a TokensWithdrawn event raised by the StakingPortBacker contract.
type StakingPortBackerTokensWithdrawn struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokensWithdrawn is a free log retrieval operation binding the contract event 0x9c6393f251205f9e03559951cab4c9ae71767b6174f77944a5b0c2fa51fbda9f.
//
// Solidity: event TokensWithdrawn(uint256 amount)
func (_StakingPortBacker *StakingPortBackerFilterer) FilterTokensWithdrawn(opts *bind.FilterOpts) (*StakingPortBackerTokensWithdrawnIterator, error) {

	logs, sub, err := _StakingPortBacker.contract.FilterLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return &StakingPortBackerTokensWithdrawnIterator{contract: _StakingPortBacker.contract, event: "TokensWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTokensWithdrawn is a free log subscription operation binding the contract event 0x9c6393f251205f9e03559951cab4c9ae71767b6174f77944a5b0c2fa51fbda9f.
//
// Solidity: event TokensWithdrawn(uint256 amount)
func (_StakingPortBacker *StakingPortBackerFilterer) WatchTokensWithdrawn(opts *bind.WatchOpts, sink chan<- *StakingPortBackerTokensWithdrawn) (event.Subscription, error) {

	logs, sub, err := _StakingPortBacker.contract.WatchLogs(opts, "TokensWithdrawn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPortBackerTokensWithdrawn)
				if err := _StakingPortBacker.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
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

// ParseTokensWithdrawn is a log parse operation binding the contract event 0x9c6393f251205f9e03559951cab4c9ae71767b6174f77944a5b0c2fa51fbda9f.
//
// Solidity: event TokensWithdrawn(uint256 amount)
func (_StakingPortBacker *StakingPortBackerFilterer) ParseTokensWithdrawn(log types.Log) (*StakingPortBackerTokensWithdrawn, error) {
	event := new(StakingPortBackerTokensWithdrawn)
	if err := _StakingPortBacker.contract.UnpackLog(event, "TokensWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}
