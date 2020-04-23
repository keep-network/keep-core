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

// TokenGrantABI is the input ABI used to generate the binding from.
const TokenGrantABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"TokenGrantCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"TokenGrantRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"grantId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"TokenGrantStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"grantId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenGrantWithdrawn\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"}],\"name\":\"authorizeStakingContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_grantId\",\"type\":\"uint256\"}],\"name\":\"availableToStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"cancelRevokedStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"cancelStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getGrant\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"staked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revokedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revokedAt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"getGrantStakeDetails\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"grantId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getGrantUnlockingSchedule\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"grantManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"policy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"}],\"name\":\"getGranteeOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_granteeOrGrantManager\",\"type\":\"address\"}],\"name\":\"getGrants\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"grantIndices\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"grantStakes\",\"outputs\":[{\"internalType\":\"contractTokenGrantStake\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"granteesToOperators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"grants\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"grantManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"grantee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"revokedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revokedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"revokedWithdrawn\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"revocable\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cliff\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"staked\",\"type\":\"uint256\"},{\"internalType\":\"contractGrantStakingPolicy\",\"name\":\"stakingPolicy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numGrants\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"receiveApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"recoverStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"revoke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakingContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"stakeBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractERC20Burnable\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"undelegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"undelegateRevoked\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"unlockedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"withdrawRevoked\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"withdrawable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenGrant is an auto generated Go binding around an Ethereum contract.
type TokenGrant struct {
	TokenGrantCaller     // Read-only binding to the contract
	TokenGrantTransactor // Write-only binding to the contract
	TokenGrantFilterer   // Log filterer for contract events
}

// TokenGrantCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenGrantCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenGrantTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenGrantFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGrantSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenGrantSession struct {
	Contract     *TokenGrant       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenGrantCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenGrantCallerSession struct {
	Contract *TokenGrantCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TokenGrantTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenGrantTransactorSession struct {
	Contract     *TokenGrantTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TokenGrantRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenGrantRaw struct {
	Contract *TokenGrant // Generic contract binding to access the raw methods on
}

// TokenGrantCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenGrantCallerRaw struct {
	Contract *TokenGrantCaller // Generic read-only contract binding to access the raw methods on
}

// TokenGrantTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenGrantTransactorRaw struct {
	Contract *TokenGrantTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenGrant creates a new instance of TokenGrant, bound to a specific deployed contract.
func NewTokenGrant(address common.Address, backend bind.ContractBackend) (*TokenGrant, error) {
	contract, err := bindTokenGrant(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenGrant{TokenGrantCaller: TokenGrantCaller{contract: contract}, TokenGrantTransactor: TokenGrantTransactor{contract: contract}, TokenGrantFilterer: TokenGrantFilterer{contract: contract}}, nil
}

// NewTokenGrantCaller creates a new read-only instance of TokenGrant, bound to a specific deployed contract.
func NewTokenGrantCaller(address common.Address, caller bind.ContractCaller) (*TokenGrantCaller, error) {
	contract, err := bindTokenGrant(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGrantCaller{contract: contract}, nil
}

// NewTokenGrantTransactor creates a new write-only instance of TokenGrant, bound to a specific deployed contract.
func NewTokenGrantTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenGrantTransactor, error) {
	contract, err := bindTokenGrant(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGrantTransactor{contract: contract}, nil
}

// NewTokenGrantFilterer creates a new log filterer instance of TokenGrant, bound to a specific deployed contract.
func NewTokenGrantFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenGrantFilterer, error) {
	contract, err := bindTokenGrant(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenGrantFilterer{contract: contract}, nil
}

// bindTokenGrant binds a generic wrapper to an already deployed contract.
func bindTokenGrant(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenGrantABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGrant *TokenGrantRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenGrant.Contract.TokenGrantCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGrant *TokenGrantRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrant.Contract.TokenGrantTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGrant *TokenGrantRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGrant.Contract.TokenGrantTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGrant *TokenGrantCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TokenGrant.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGrant *TokenGrantTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGrant.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGrant *TokenGrantTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGrant.Contract.contract.Transact(opts, method, params...)
}

// AvailableToStake is a free data retrieval call binding the contract method 0x910bcf61.
//
// Solidity: function availableToStake(uint256 _grantId) constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) AvailableToStake(opts *bind.CallOpts, _grantId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "availableToStake", _grantId)
	return *ret0, err
}

// AvailableToStake is a free data retrieval call binding the contract method 0x910bcf61.
//
// Solidity: function availableToStake(uint256 _grantId) constant returns(uint256)
func (_TokenGrant *TokenGrantSession) AvailableToStake(_grantId *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.AvailableToStake(&_TokenGrant.CallOpts, _grantId)
}

// AvailableToStake is a free data retrieval call binding the contract method 0x910bcf61.
//
// Solidity: function availableToStake(uint256 _grantId) constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) AvailableToStake(_grantId *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.AvailableToStake(&_TokenGrant.CallOpts, _grantId)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.BalanceOf(&_TokenGrant.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.BalanceOf(&_TokenGrant.CallOpts, _owner)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TokenGrant *TokenGrantSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.Balances(&_TokenGrant.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.Balances(&_TokenGrant.CallOpts, arg0)
}

// GetGrant is a free data retrieval call binding the contract method 0xb8cc6c93.
//
// Solidity: function getGrant(uint256 _id) constant returns(uint256 amount, uint256 withdrawn, uint256 staked, uint256 revokedAmount, uint256 revokedAt, address grantee)
func (_TokenGrant *TokenGrantCaller) GetGrant(opts *bind.CallOpts, _id *big.Int) (struct {
	Amount        *big.Int
	Withdrawn     *big.Int
	Staked        *big.Int
	RevokedAmount *big.Int
	RevokedAt     *big.Int
	Grantee       common.Address
}, error) {
	ret := new(struct {
		Amount        *big.Int
		Withdrawn     *big.Int
		Staked        *big.Int
		RevokedAmount *big.Int
		RevokedAt     *big.Int
		Grantee       common.Address
	})
	out := ret
	err := _TokenGrant.contract.Call(opts, out, "getGrant", _id)
	return *ret, err
}

// GetGrant is a free data retrieval call binding the contract method 0xb8cc6c93.
//
// Solidity: function getGrant(uint256 _id) constant returns(uint256 amount, uint256 withdrawn, uint256 staked, uint256 revokedAmount, uint256 revokedAt, address grantee)
func (_TokenGrant *TokenGrantSession) GetGrant(_id *big.Int) (struct {
	Amount        *big.Int
	Withdrawn     *big.Int
	Staked        *big.Int
	RevokedAmount *big.Int
	RevokedAt     *big.Int
	Grantee       common.Address
}, error) {
	return _TokenGrant.Contract.GetGrant(&_TokenGrant.CallOpts, _id)
}

// GetGrant is a free data retrieval call binding the contract method 0xb8cc6c93.
//
// Solidity: function getGrant(uint256 _id) constant returns(uint256 amount, uint256 withdrawn, uint256 staked, uint256 revokedAmount, uint256 revokedAt, address grantee)
func (_TokenGrant *TokenGrantCallerSession) GetGrant(_id *big.Int) (struct {
	Amount        *big.Int
	Withdrawn     *big.Int
	Staked        *big.Int
	RevokedAmount *big.Int
	RevokedAt     *big.Int
	Grantee       common.Address
}, error) {
	return _TokenGrant.Contract.GetGrant(&_TokenGrant.CallOpts, _id)
}

// GetGrantStakeDetails is a free data retrieval call binding the contract method 0xf29e2737.
//
// Solidity: function getGrantStakeDetails(address operator) constant returns(uint256 grantId, uint256 amount, address stakingContract)
func (_TokenGrant *TokenGrantCaller) GetGrantStakeDetails(opts *bind.CallOpts, operator common.Address) (struct {
	GrantId         *big.Int
	Amount          *big.Int
	StakingContract common.Address
}, error) {
	ret := new(struct {
		GrantId         *big.Int
		Amount          *big.Int
		StakingContract common.Address
	})
	out := ret
	err := _TokenGrant.contract.Call(opts, out, "getGrantStakeDetails", operator)
	return *ret, err
}

// GetGrantStakeDetails is a free data retrieval call binding the contract method 0xf29e2737.
//
// Solidity: function getGrantStakeDetails(address operator) constant returns(uint256 grantId, uint256 amount, address stakingContract)
func (_TokenGrant *TokenGrantSession) GetGrantStakeDetails(operator common.Address) (struct {
	GrantId         *big.Int
	Amount          *big.Int
	StakingContract common.Address
}, error) {
	return _TokenGrant.Contract.GetGrantStakeDetails(&_TokenGrant.CallOpts, operator)
}

// GetGrantStakeDetails is a free data retrieval call binding the contract method 0xf29e2737.
//
// Solidity: function getGrantStakeDetails(address operator) constant returns(uint256 grantId, uint256 amount, address stakingContract)
func (_TokenGrant *TokenGrantCallerSession) GetGrantStakeDetails(operator common.Address) (struct {
	GrantId         *big.Int
	Amount          *big.Int
	StakingContract common.Address
}, error) {
	return _TokenGrant.Contract.GetGrantStakeDetails(&_TokenGrant.CallOpts, operator)
}

// GetGrantUnlockingSchedule is a free data retrieval call binding the contract method 0xc6076b7a.
//
// Solidity: function getGrantUnlockingSchedule(uint256 _id) constant returns(address grantManager, uint256 duration, uint256 start, uint256 cliff, address policy)
func (_TokenGrant *TokenGrantCaller) GetGrantUnlockingSchedule(opts *bind.CallOpts, _id *big.Int) (struct {
	GrantManager common.Address
	Duration     *big.Int
	Start        *big.Int
	Cliff        *big.Int
	Policy       common.Address
}, error) {
	ret := new(struct {
		GrantManager common.Address
		Duration     *big.Int
		Start        *big.Int
		Cliff        *big.Int
		Policy       common.Address
	})
	out := ret
	err := _TokenGrant.contract.Call(opts, out, "getGrantUnlockingSchedule", _id)
	return *ret, err
}

// GetGrantUnlockingSchedule is a free data retrieval call binding the contract method 0xc6076b7a.
//
// Solidity: function getGrantUnlockingSchedule(uint256 _id) constant returns(address grantManager, uint256 duration, uint256 start, uint256 cliff, address policy)
func (_TokenGrant *TokenGrantSession) GetGrantUnlockingSchedule(_id *big.Int) (struct {
	GrantManager common.Address
	Duration     *big.Int
	Start        *big.Int
	Cliff        *big.Int
	Policy       common.Address
}, error) {
	return _TokenGrant.Contract.GetGrantUnlockingSchedule(&_TokenGrant.CallOpts, _id)
}

// GetGrantUnlockingSchedule is a free data retrieval call binding the contract method 0xc6076b7a.
//
// Solidity: function getGrantUnlockingSchedule(uint256 _id) constant returns(address grantManager, uint256 duration, uint256 start, uint256 cliff, address policy)
func (_TokenGrant *TokenGrantCallerSession) GetGrantUnlockingSchedule(_id *big.Int) (struct {
	GrantManager common.Address
	Duration     *big.Int
	Start        *big.Int
	Cliff        *big.Int
	Policy       common.Address
}, error) {
	return _TokenGrant.Contract.GetGrantUnlockingSchedule(&_TokenGrant.CallOpts, _id)
}

// GetGranteeOperators is a free data retrieval call binding the contract method 0xca9a33ce.
//
// Solidity: function getGranteeOperators(address grantee) constant returns(address[])
func (_TokenGrant *TokenGrantCaller) GetGranteeOperators(opts *bind.CallOpts, grantee common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "getGranteeOperators", grantee)
	return *ret0, err
}

// GetGranteeOperators is a free data retrieval call binding the contract method 0xca9a33ce.
//
// Solidity: function getGranteeOperators(address grantee) constant returns(address[])
func (_TokenGrant *TokenGrantSession) GetGranteeOperators(grantee common.Address) ([]common.Address, error) {
	return _TokenGrant.Contract.GetGranteeOperators(&_TokenGrant.CallOpts, grantee)
}

// GetGranteeOperators is a free data retrieval call binding the contract method 0xca9a33ce.
//
// Solidity: function getGranteeOperators(address grantee) constant returns(address[])
func (_TokenGrant *TokenGrantCallerSession) GetGranteeOperators(grantee common.Address) ([]common.Address, error) {
	return _TokenGrant.Contract.GetGranteeOperators(&_TokenGrant.CallOpts, grantee)
}

// GetGrants is a free data retrieval call binding the contract method 0xf3ad1c1a.
//
// Solidity: function getGrants(address _granteeOrGrantManager) constant returns(uint256[])
func (_TokenGrant *TokenGrantCaller) GetGrants(opts *bind.CallOpts, _granteeOrGrantManager common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "getGrants", _granteeOrGrantManager)
	return *ret0, err
}

// GetGrants is a free data retrieval call binding the contract method 0xf3ad1c1a.
//
// Solidity: function getGrants(address _granteeOrGrantManager) constant returns(uint256[])
func (_TokenGrant *TokenGrantSession) GetGrants(_granteeOrGrantManager common.Address) ([]*big.Int, error) {
	return _TokenGrant.Contract.GetGrants(&_TokenGrant.CallOpts, _granteeOrGrantManager)
}

// GetGrants is a free data retrieval call binding the contract method 0xf3ad1c1a.
//
// Solidity: function getGrants(address _granteeOrGrantManager) constant returns(uint256[])
func (_TokenGrant *TokenGrantCallerSession) GetGrants(_granteeOrGrantManager common.Address) ([]*big.Int, error) {
	return _TokenGrant.Contract.GetGrants(&_TokenGrant.CallOpts, _granteeOrGrantManager)
}

// GrantIndices is a free data retrieval call binding the contract method 0x73297585.
//
// Solidity: function grantIndices(address , uint256 ) constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) GrantIndices(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "grantIndices", arg0, arg1)
	return *ret0, err
}

// GrantIndices is a free data retrieval call binding the contract method 0x73297585.
//
// Solidity: function grantIndices(address , uint256 ) constant returns(uint256)
func (_TokenGrant *TokenGrantSession) GrantIndices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.GrantIndices(&_TokenGrant.CallOpts, arg0, arg1)
}

// GrantIndices is a free data retrieval call binding the contract method 0x73297585.
//
// Solidity: function grantIndices(address , uint256 ) constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) GrantIndices(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.GrantIndices(&_TokenGrant.CallOpts, arg0, arg1)
}

// GrantStakes is a free data retrieval call binding the contract method 0x2496bfb3.
//
// Solidity: function grantStakes(address ) constant returns(address)
func (_TokenGrant *TokenGrantCaller) GrantStakes(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "grantStakes", arg0)
	return *ret0, err
}

// GrantStakes is a free data retrieval call binding the contract method 0x2496bfb3.
//
// Solidity: function grantStakes(address ) constant returns(address)
func (_TokenGrant *TokenGrantSession) GrantStakes(arg0 common.Address) (common.Address, error) {
	return _TokenGrant.Contract.GrantStakes(&_TokenGrant.CallOpts, arg0)
}

// GrantStakes is a free data retrieval call binding the contract method 0x2496bfb3.
//
// Solidity: function grantStakes(address ) constant returns(address)
func (_TokenGrant *TokenGrantCallerSession) GrantStakes(arg0 common.Address) (common.Address, error) {
	return _TokenGrant.Contract.GrantStakes(&_TokenGrant.CallOpts, arg0)
}

// GranteesToOperators is a free data retrieval call binding the contract method 0x9ae6e9c4.
//
// Solidity: function granteesToOperators(address , uint256 ) constant returns(address)
func (_TokenGrant *TokenGrantCaller) GranteesToOperators(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "granteesToOperators", arg0, arg1)
	return *ret0, err
}

// GranteesToOperators is a free data retrieval call binding the contract method 0x9ae6e9c4.
//
// Solidity: function granteesToOperators(address , uint256 ) constant returns(address)
func (_TokenGrant *TokenGrantSession) GranteesToOperators(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _TokenGrant.Contract.GranteesToOperators(&_TokenGrant.CallOpts, arg0, arg1)
}

// GranteesToOperators is a free data retrieval call binding the contract method 0x9ae6e9c4.
//
// Solidity: function granteesToOperators(address , uint256 ) constant returns(address)
func (_TokenGrant *TokenGrantCallerSession) GranteesToOperators(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _TokenGrant.Contract.GranteesToOperators(&_TokenGrant.CallOpts, arg0, arg1)
}

// Grants is a free data retrieval call binding the contract method 0x0c0debea.
//
// Solidity: function grants(uint256 ) constant returns(address grantManager, address grantee, uint256 revokedAt, uint256 revokedAmount, uint256 revokedWithdrawn, bool revocable, uint256 amount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn, uint256 staked, address stakingPolicy)
func (_TokenGrant *TokenGrantCaller) Grants(opts *bind.CallOpts, arg0 *big.Int) (struct {
	GrantManager     common.Address
	Grantee          common.Address
	RevokedAt        *big.Int
	RevokedAmount    *big.Int
	RevokedWithdrawn *big.Int
	Revocable        bool
	Amount           *big.Int
	Duration         *big.Int
	Start            *big.Int
	Cliff            *big.Int
	Withdrawn        *big.Int
	Staked           *big.Int
	StakingPolicy    common.Address
}, error) {
	ret := new(struct {
		GrantManager     common.Address
		Grantee          common.Address
		RevokedAt        *big.Int
		RevokedAmount    *big.Int
		RevokedWithdrawn *big.Int
		Revocable        bool
		Amount           *big.Int
		Duration         *big.Int
		Start            *big.Int
		Cliff            *big.Int
		Withdrawn        *big.Int
		Staked           *big.Int
		StakingPolicy    common.Address
	})
	out := ret
	err := _TokenGrant.contract.Call(opts, out, "grants", arg0)
	return *ret, err
}

// Grants is a free data retrieval call binding the contract method 0x0c0debea.
//
// Solidity: function grants(uint256 ) constant returns(address grantManager, address grantee, uint256 revokedAt, uint256 revokedAmount, uint256 revokedWithdrawn, bool revocable, uint256 amount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn, uint256 staked, address stakingPolicy)
func (_TokenGrant *TokenGrantSession) Grants(arg0 *big.Int) (struct {
	GrantManager     common.Address
	Grantee          common.Address
	RevokedAt        *big.Int
	RevokedAmount    *big.Int
	RevokedWithdrawn *big.Int
	Revocable        bool
	Amount           *big.Int
	Duration         *big.Int
	Start            *big.Int
	Cliff            *big.Int
	Withdrawn        *big.Int
	Staked           *big.Int
	StakingPolicy    common.Address
}, error) {
	return _TokenGrant.Contract.Grants(&_TokenGrant.CallOpts, arg0)
}

// Grants is a free data retrieval call binding the contract method 0x0c0debea.
//
// Solidity: function grants(uint256 ) constant returns(address grantManager, address grantee, uint256 revokedAt, uint256 revokedAmount, uint256 revokedWithdrawn, bool revocable, uint256 amount, uint256 duration, uint256 start, uint256 cliff, uint256 withdrawn, uint256 staked, address stakingPolicy)
func (_TokenGrant *TokenGrantCallerSession) Grants(arg0 *big.Int) (struct {
	GrantManager     common.Address
	Grantee          common.Address
	RevokedAt        *big.Int
	RevokedAmount    *big.Int
	RevokedWithdrawn *big.Int
	Revocable        bool
	Amount           *big.Int
	Duration         *big.Int
	Start            *big.Int
	Cliff            *big.Int
	Withdrawn        *big.Int
	Staked           *big.Int
	StakingPolicy    common.Address
}, error) {
	return _TokenGrant.Contract.Grants(&_TokenGrant.CallOpts, arg0)
}

// NumGrants is a free data retrieval call binding the contract method 0xf81a6463.
//
// Solidity: function numGrants() constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) NumGrants(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "numGrants")
	return *ret0, err
}

// NumGrants is a free data retrieval call binding the contract method 0xf81a6463.
//
// Solidity: function numGrants() constant returns(uint256)
func (_TokenGrant *TokenGrantSession) NumGrants() (*big.Int, error) {
	return _TokenGrant.Contract.NumGrants(&_TokenGrant.CallOpts)
}

// NumGrants is a free data retrieval call binding the contract method 0xf81a6463.
//
// Solidity: function numGrants() constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) NumGrants() (*big.Int, error) {
	return _TokenGrant.Contract.NumGrants(&_TokenGrant.CallOpts)
}

// StakeBalanceOf is a free data retrieval call binding the contract method 0xaf46aa08.
//
// Solidity: function stakeBalanceOf(address _address) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantCaller) StakeBalanceOf(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "stakeBalanceOf", _address)
	return *ret0, err
}

// StakeBalanceOf is a free data retrieval call binding the contract method 0xaf46aa08.
//
// Solidity: function stakeBalanceOf(address _address) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantSession) StakeBalanceOf(_address common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.StakeBalanceOf(&_TokenGrant.CallOpts, _address)
}

// StakeBalanceOf is a free data retrieval call binding the contract method 0xaf46aa08.
//
// Solidity: function stakeBalanceOf(address _address) constant returns(uint256 balance)
func (_TokenGrant *TokenGrantCallerSession) StakeBalanceOf(_address common.Address) (*big.Int, error) {
	return _TokenGrant.Contract.StakeBalanceOf(&_TokenGrant.CallOpts, _address)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenGrant *TokenGrantCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenGrant *TokenGrantSession) Token() (common.Address, error) {
	return _TokenGrant.Contract.Token(&_TokenGrant.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_TokenGrant *TokenGrantCallerSession) Token() (common.Address, error) {
	return _TokenGrant.Contract.Token(&_TokenGrant.CallOpts)
}

// UnlockedAmount is a free data retrieval call binding the contract method 0xc90376e4.
//
// Solidity: function unlockedAmount(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) UnlockedAmount(opts *bind.CallOpts, _id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "unlockedAmount", _id)
	return *ret0, err
}

// UnlockedAmount is a free data retrieval call binding the contract method 0xc90376e4.
//
// Solidity: function unlockedAmount(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantSession) UnlockedAmount(_id *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.UnlockedAmount(&_TokenGrant.CallOpts, _id)
}

// UnlockedAmount is a free data retrieval call binding the contract method 0xc90376e4.
//
// Solidity: function unlockedAmount(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) UnlockedAmount(_id *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.UnlockedAmount(&_TokenGrant.CallOpts, _id)
}

// Withdrawable is a free data retrieval call binding the contract method 0xf11988e0.
//
// Solidity: function withdrawable(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantCaller) Withdrawable(opts *bind.CallOpts, _id *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TokenGrant.contract.Call(opts, out, "withdrawable", _id)
	return *ret0, err
}

// Withdrawable is a free data retrieval call binding the contract method 0xf11988e0.
//
// Solidity: function withdrawable(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantSession) Withdrawable(_id *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.Withdrawable(&_TokenGrant.CallOpts, _id)
}

// Withdrawable is a free data retrieval call binding the contract method 0xf11988e0.
//
// Solidity: function withdrawable(uint256 _id) constant returns(uint256)
func (_TokenGrant *TokenGrantCallerSession) Withdrawable(_id *big.Int) (*big.Int, error) {
	return _TokenGrant.Contract.Withdrawable(&_TokenGrant.CallOpts, _id)
}

// AuthorizeStakingContract is a paid mutator transaction binding the contract method 0xdba6d924.
//
// Solidity: function authorizeStakingContract(address _stakingContract) returns()
func (_TokenGrant *TokenGrantTransactor) AuthorizeStakingContract(opts *bind.TransactOpts, _stakingContract common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "authorizeStakingContract", _stakingContract)
}

// AuthorizeStakingContract is a paid mutator transaction binding the contract method 0xdba6d924.
//
// Solidity: function authorizeStakingContract(address _stakingContract) returns()
func (_TokenGrant *TokenGrantSession) AuthorizeStakingContract(_stakingContract common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.AuthorizeStakingContract(&_TokenGrant.TransactOpts, _stakingContract)
}

// AuthorizeStakingContract is a paid mutator transaction binding the contract method 0xdba6d924.
//
// Solidity: function authorizeStakingContract(address _stakingContract) returns()
func (_TokenGrant *TokenGrantTransactorSession) AuthorizeStakingContract(_stakingContract common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.AuthorizeStakingContract(&_TokenGrant.TransactOpts, _stakingContract)
}

// CancelRevokedStake is a paid mutator transaction binding the contract method 0x3934af5c.
//
// Solidity: function cancelRevokedStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactor) CancelRevokedStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "cancelRevokedStake", _operator)
}

// CancelRevokedStake is a paid mutator transaction binding the contract method 0x3934af5c.
//
// Solidity: function cancelRevokedStake(address _operator) returns()
func (_TokenGrant *TokenGrantSession) CancelRevokedStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.CancelRevokedStake(&_TokenGrant.TransactOpts, _operator)
}

// CancelRevokedStake is a paid mutator transaction binding the contract method 0x3934af5c.
//
// Solidity: function cancelRevokedStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactorSession) CancelRevokedStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.CancelRevokedStake(&_TokenGrant.TransactOpts, _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactor) CancelStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "cancelStake", _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenGrant *TokenGrantSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.CancelStake(&_TokenGrant.TransactOpts, _operator)
}

// CancelStake is a paid mutator transaction binding the contract method 0xe064172e.
//
// Solidity: function cancelStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactorSession) CancelStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.CancelStake(&_TokenGrant.TransactOpts, _operator)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_TokenGrant *TokenGrantTransactor) ReceiveApproval(opts *bind.TransactOpts, _from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "receiveApproval", _from, _amount, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_TokenGrant *TokenGrantSession) ReceiveApproval(_from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.Contract.ReceiveApproval(&_TokenGrant.TransactOpts, _from, _amount, _token, _extraData)
}

// ReceiveApproval is a paid mutator transaction binding the contract method 0x8f4ffcb1.
//
// Solidity: function receiveApproval(address _from, uint256 _amount, address _token, bytes _extraData) returns()
func (_TokenGrant *TokenGrantTransactorSession) ReceiveApproval(_from common.Address, _amount *big.Int, _token common.Address, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.Contract.ReceiveApproval(&_TokenGrant.TransactOpts, _from, _amount, _token, _extraData)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactor) RecoverStake(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "recoverStake", _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenGrant *TokenGrantSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.RecoverStake(&_TokenGrant.TransactOpts, _operator)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x525835f9.
//
// Solidity: function recoverStake(address _operator) returns()
func (_TokenGrant *TokenGrantTransactorSession) RecoverStake(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.RecoverStake(&_TokenGrant.TransactOpts, _operator)
}

// Revoke is a paid mutator transaction binding the contract method 0x20c5429b.
//
// Solidity: function revoke(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactor) Revoke(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "revoke", _id)
}

// Revoke is a paid mutator transaction binding the contract method 0x20c5429b.
//
// Solidity: function revoke(uint256 _id) returns()
func (_TokenGrant *TokenGrantSession) Revoke(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.Revoke(&_TokenGrant.TransactOpts, _id)
}

// Revoke is a paid mutator transaction binding the contract method 0x20c5429b.
//
// Solidity: function revoke(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactorSession) Revoke(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.Revoke(&_TokenGrant.TransactOpts, _id)
}

// Stake is a paid mutator transaction binding the contract method 0x5860280c.
//
// Solidity: function stake(uint256 _id, address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_TokenGrant *TokenGrantTransactor) Stake(opts *bind.TransactOpts, _id *big.Int, _stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "stake", _id, _stakingContract, _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x5860280c.
//
// Solidity: function stake(uint256 _id, address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_TokenGrant *TokenGrantSession) Stake(_id *big.Int, _stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.Contract.Stake(&_TokenGrant.TransactOpts, _id, _stakingContract, _amount, _extraData)
}

// Stake is a paid mutator transaction binding the contract method 0x5860280c.
//
// Solidity: function stake(uint256 _id, address _stakingContract, uint256 _amount, bytes _extraData) returns()
func (_TokenGrant *TokenGrantTransactorSession) Stake(_id *big.Int, _stakingContract common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _TokenGrant.Contract.Stake(&_TokenGrant.TransactOpts, _id, _stakingContract, _amount, _extraData)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenGrant *TokenGrantTransactor) Undelegate(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "undelegate", _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenGrant *TokenGrantSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.Undelegate(&_TokenGrant.TransactOpts, _operator)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address _operator) returns()
func (_TokenGrant *TokenGrantTransactorSession) Undelegate(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.Undelegate(&_TokenGrant.TransactOpts, _operator)
}

// UndelegateRevoked is a paid mutator transaction binding the contract method 0x9bf85078.
//
// Solidity: function undelegateRevoked(address _operator) returns()
func (_TokenGrant *TokenGrantTransactor) UndelegateRevoked(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "undelegateRevoked", _operator)
}

// UndelegateRevoked is a paid mutator transaction binding the contract method 0x9bf85078.
//
// Solidity: function undelegateRevoked(address _operator) returns()
func (_TokenGrant *TokenGrantSession) UndelegateRevoked(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.UndelegateRevoked(&_TokenGrant.TransactOpts, _operator)
}

// UndelegateRevoked is a paid mutator transaction binding the contract method 0x9bf85078.
//
// Solidity: function undelegateRevoked(address _operator) returns()
func (_TokenGrant *TokenGrantTransactorSession) UndelegateRevoked(_operator common.Address) (*types.Transaction, error) {
	return _TokenGrant.Contract.UndelegateRevoked(&_TokenGrant.TransactOpts, _operator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactor) Withdraw(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "withdraw", _id)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _id) returns()
func (_TokenGrant *TokenGrantSession) Withdraw(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.Withdraw(&_TokenGrant.TransactOpts, _id)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactorSession) Withdraw(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.Withdraw(&_TokenGrant.TransactOpts, _id)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0x8e34e561.
//
// Solidity: function withdrawRevoked(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactor) WithdrawRevoked(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.contract.Transact(opts, "withdrawRevoked", _id)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0x8e34e561.
//
// Solidity: function withdrawRevoked(uint256 _id) returns()
func (_TokenGrant *TokenGrantSession) WithdrawRevoked(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.WithdrawRevoked(&_TokenGrant.TransactOpts, _id)
}

// WithdrawRevoked is a paid mutator transaction binding the contract method 0x8e34e561.
//
// Solidity: function withdrawRevoked(uint256 _id) returns()
func (_TokenGrant *TokenGrantTransactorSession) WithdrawRevoked(_id *big.Int) (*types.Transaction, error) {
	return _TokenGrant.Contract.WithdrawRevoked(&_TokenGrant.TransactOpts, _id)
}

// TokenGrantTokenGrantCreatedIterator is returned from FilterTokenGrantCreated and is used to iterate over the raw logs and unpacked data for TokenGrantCreated events raised by the TokenGrant contract.
type TokenGrantTokenGrantCreatedIterator struct {
	Event *TokenGrantTokenGrantCreated // Event containing the contract specifics and raw log

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
func (it *TokenGrantTokenGrantCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGrantTokenGrantCreated)
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
		it.Event = new(TokenGrantTokenGrantCreated)
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
func (it *TokenGrantTokenGrantCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGrantTokenGrantCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGrantTokenGrantCreated represents a TokenGrantCreated event raised by the TokenGrant contract.
type TokenGrantTokenGrantCreated struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTokenGrantCreated is a free log retrieval operation binding the contract event 0x74b1b578efe34f4348abe981c705578eb9c7f0e25fb6e5bb0eeeb33c67854f83.
//
// Solidity: event TokenGrantCreated(uint256 id)
func (_TokenGrant *TokenGrantFilterer) FilterTokenGrantCreated(opts *bind.FilterOpts) (*TokenGrantTokenGrantCreatedIterator, error) {

	logs, sub, err := _TokenGrant.contract.FilterLogs(opts, "TokenGrantCreated")
	if err != nil {
		return nil, err
	}
	return &TokenGrantTokenGrantCreatedIterator{contract: _TokenGrant.contract, event: "TokenGrantCreated", logs: logs, sub: sub}, nil
}

// WatchTokenGrantCreated is a free log subscription operation binding the contract event 0x74b1b578efe34f4348abe981c705578eb9c7f0e25fb6e5bb0eeeb33c67854f83.
//
// Solidity: event TokenGrantCreated(uint256 id)
func (_TokenGrant *TokenGrantFilterer) WatchTokenGrantCreated(opts *bind.WatchOpts, sink chan<- *TokenGrantTokenGrantCreated) (event.Subscription, error) {

	logs, sub, err := _TokenGrant.contract.WatchLogs(opts, "TokenGrantCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGrantTokenGrantCreated)
				if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantCreated", log); err != nil {
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

// ParseTokenGrantCreated is a log parse operation binding the contract event 0x74b1b578efe34f4348abe981c705578eb9c7f0e25fb6e5bb0eeeb33c67854f83.
//
// Solidity: event TokenGrantCreated(uint256 id)
func (_TokenGrant *TokenGrantFilterer) ParseTokenGrantCreated(log types.Log) (*TokenGrantTokenGrantCreated, error) {
	event := new(TokenGrantTokenGrantCreated)
	if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenGrantTokenGrantRevokedIterator is returned from FilterTokenGrantRevoked and is used to iterate over the raw logs and unpacked data for TokenGrantRevoked events raised by the TokenGrant contract.
type TokenGrantTokenGrantRevokedIterator struct {
	Event *TokenGrantTokenGrantRevoked // Event containing the contract specifics and raw log

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
func (it *TokenGrantTokenGrantRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGrantTokenGrantRevoked)
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
		it.Event = new(TokenGrantTokenGrantRevoked)
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
func (it *TokenGrantTokenGrantRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGrantTokenGrantRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGrantTokenGrantRevoked represents a TokenGrantRevoked event raised by the TokenGrant contract.
type TokenGrantTokenGrantRevoked struct {
	Id  *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTokenGrantRevoked is a free log retrieval operation binding the contract event 0x19791073e5f2b48fe653897875e9f5893ac441ce6df89f2053538191e5c47631.
//
// Solidity: event TokenGrantRevoked(uint256 id)
func (_TokenGrant *TokenGrantFilterer) FilterTokenGrantRevoked(opts *bind.FilterOpts) (*TokenGrantTokenGrantRevokedIterator, error) {

	logs, sub, err := _TokenGrant.contract.FilterLogs(opts, "TokenGrantRevoked")
	if err != nil {
		return nil, err
	}
	return &TokenGrantTokenGrantRevokedIterator{contract: _TokenGrant.contract, event: "TokenGrantRevoked", logs: logs, sub: sub}, nil
}

// WatchTokenGrantRevoked is a free log subscription operation binding the contract event 0x19791073e5f2b48fe653897875e9f5893ac441ce6df89f2053538191e5c47631.
//
// Solidity: event TokenGrantRevoked(uint256 id)
func (_TokenGrant *TokenGrantFilterer) WatchTokenGrantRevoked(opts *bind.WatchOpts, sink chan<- *TokenGrantTokenGrantRevoked) (event.Subscription, error) {

	logs, sub, err := _TokenGrant.contract.WatchLogs(opts, "TokenGrantRevoked")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGrantTokenGrantRevoked)
				if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantRevoked", log); err != nil {
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

// ParseTokenGrantRevoked is a log parse operation binding the contract event 0x19791073e5f2b48fe653897875e9f5893ac441ce6df89f2053538191e5c47631.
//
// Solidity: event TokenGrantRevoked(uint256 id)
func (_TokenGrant *TokenGrantFilterer) ParseTokenGrantRevoked(log types.Log) (*TokenGrantTokenGrantRevoked, error) {
	event := new(TokenGrantTokenGrantRevoked)
	if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantRevoked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenGrantTokenGrantStakedIterator is returned from FilterTokenGrantStaked and is used to iterate over the raw logs and unpacked data for TokenGrantStaked events raised by the TokenGrant contract.
type TokenGrantTokenGrantStakedIterator struct {
	Event *TokenGrantTokenGrantStaked // Event containing the contract specifics and raw log

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
func (it *TokenGrantTokenGrantStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGrantTokenGrantStaked)
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
		it.Event = new(TokenGrantTokenGrantStaked)
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
func (it *TokenGrantTokenGrantStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGrantTokenGrantStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGrantTokenGrantStaked represents a TokenGrantStaked event raised by the TokenGrant contract.
type TokenGrantTokenGrantStaked struct {
	GrantId  *big.Int
	Amount   *big.Int
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTokenGrantStaked is a free log retrieval operation binding the contract event 0xf05c07b89b3e4ff57b17aa6883ec35e90d7710c57a038c49c3ec3911a80c2445.
//
// Solidity: event TokenGrantStaked(uint256 indexed grantId, uint256 amount, address operator)
func (_TokenGrant *TokenGrantFilterer) FilterTokenGrantStaked(opts *bind.FilterOpts, grantId []*big.Int) (*TokenGrantTokenGrantStakedIterator, error) {

	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenGrant.contract.FilterLogs(opts, "TokenGrantStaked", grantIdRule)
	if err != nil {
		return nil, err
	}
	return &TokenGrantTokenGrantStakedIterator{contract: _TokenGrant.contract, event: "TokenGrantStaked", logs: logs, sub: sub}, nil
}

// WatchTokenGrantStaked is a free log subscription operation binding the contract event 0xf05c07b89b3e4ff57b17aa6883ec35e90d7710c57a038c49c3ec3911a80c2445.
//
// Solidity: event TokenGrantStaked(uint256 indexed grantId, uint256 amount, address operator)
func (_TokenGrant *TokenGrantFilterer) WatchTokenGrantStaked(opts *bind.WatchOpts, sink chan<- *TokenGrantTokenGrantStaked, grantId []*big.Int) (event.Subscription, error) {

	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenGrant.contract.WatchLogs(opts, "TokenGrantStaked", grantIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGrantTokenGrantStaked)
				if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantStaked", log); err != nil {
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

// ParseTokenGrantStaked is a log parse operation binding the contract event 0xf05c07b89b3e4ff57b17aa6883ec35e90d7710c57a038c49c3ec3911a80c2445.
//
// Solidity: event TokenGrantStaked(uint256 indexed grantId, uint256 amount, address operator)
func (_TokenGrant *TokenGrantFilterer) ParseTokenGrantStaked(log types.Log) (*TokenGrantTokenGrantStaked, error) {
	event := new(TokenGrantTokenGrantStaked)
	if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantStaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenGrantTokenGrantWithdrawnIterator is returned from FilterTokenGrantWithdrawn and is used to iterate over the raw logs and unpacked data for TokenGrantWithdrawn events raised by the TokenGrant contract.
type TokenGrantTokenGrantWithdrawnIterator struct {
	Event *TokenGrantTokenGrantWithdrawn // Event containing the contract specifics and raw log

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
func (it *TokenGrantTokenGrantWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGrantTokenGrantWithdrawn)
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
		it.Event = new(TokenGrantTokenGrantWithdrawn)
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
func (it *TokenGrantTokenGrantWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGrantTokenGrantWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGrantTokenGrantWithdrawn represents a TokenGrantWithdrawn event raised by the TokenGrant contract.
type TokenGrantTokenGrantWithdrawn struct {
	GrantId *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenGrantWithdrawn is a free log retrieval operation binding the contract event 0x16b38bde3a1f1ffcf6135152793028f56443302494c5c811e63cf5f460a1e093.
//
// Solidity: event TokenGrantWithdrawn(uint256 indexed grantId, uint256 amount)
func (_TokenGrant *TokenGrantFilterer) FilterTokenGrantWithdrawn(opts *bind.FilterOpts, grantId []*big.Int) (*TokenGrantTokenGrantWithdrawnIterator, error) {

	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenGrant.contract.FilterLogs(opts, "TokenGrantWithdrawn", grantIdRule)
	if err != nil {
		return nil, err
	}
	return &TokenGrantTokenGrantWithdrawnIterator{contract: _TokenGrant.contract, event: "TokenGrantWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTokenGrantWithdrawn is a free log subscription operation binding the contract event 0x16b38bde3a1f1ffcf6135152793028f56443302494c5c811e63cf5f460a1e093.
//
// Solidity: event TokenGrantWithdrawn(uint256 indexed grantId, uint256 amount)
func (_TokenGrant *TokenGrantFilterer) WatchTokenGrantWithdrawn(opts *bind.WatchOpts, sink chan<- *TokenGrantTokenGrantWithdrawn, grantId []*big.Int) (event.Subscription, error) {

	var grantIdRule []interface{}
	for _, grantIdItem := range grantId {
		grantIdRule = append(grantIdRule, grantIdItem)
	}

	logs, sub, err := _TokenGrant.contract.WatchLogs(opts, "TokenGrantWithdrawn", grantIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGrantTokenGrantWithdrawn)
				if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantWithdrawn", log); err != nil {
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

// ParseTokenGrantWithdrawn is a log parse operation binding the contract event 0x16b38bde3a1f1ffcf6135152793028f56443302494c5c811e63cf5f460a1e093.
//
// Solidity: event TokenGrantWithdrawn(uint256 indexed grantId, uint256 amount)
func (_TokenGrant *TokenGrantFilterer) ParseTokenGrantWithdrawn(log types.Log) (*TokenGrantTokenGrantWithdrawn, error) {
	event := new(TokenGrantTokenGrantWithdrawn)
	if err := _TokenGrant.contract.UnpackLog(event, "TokenGrantWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}
