// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BitcoinTxInfo2 is an auto generated low-level Go binding around an user-defined struct.
type BitcoinTxInfo2 struct {
	Version      [4]byte
	InputVector  []byte
	OutputVector []byte
	Locktime     [4]byte
}

// WalletCoordinatorDepositExtraInfo is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositExtraInfo struct {
	FundingTx        BitcoinTxInfo2
	BlindingFactor   [8]byte
	WalletPubKeyHash [20]byte
	RefundPubKeyHash [20]byte
	RefundLocktime   [4]byte
}

// WalletCoordinatorDepositKey is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositKey struct {
	FundingTxHash      [32]byte
	FundingOutputIndex uint32
}

// WalletCoordinatorDepositSweepProposal is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorDepositSweepProposal struct {
	WalletPubKeyHash [20]byte
	DepositsKeys     []WalletCoordinatorDepositKey
	SweepTxFee       *big.Int
}

// WalletCoordinatorWalletMemberContext is an auto generated low-level Go binding around an user-defined struct.
type WalletCoordinatorWalletMemberContext struct {
	WalletMembersIDs  []uint32
	WalletMemberIndex *big.Int
}

// WalletCoordinatorMetaData contains all meta data concerning the WalletCoordinator contract.
var WalletCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositMinAge\",\"type\":\"uint32\"}],\"name\":\"DepositMinAgeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositRefundSafetyMargin\",\"type\":\"uint32\"}],\"name\":\"DepositRefundSafetyMarginUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"depositSweepMaxSize\",\"type\":\"uint16\"}],\"name\":\"DepositSweepMaxSizeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositSweepProposalSubmissionGasOffset\",\"type\":\"uint32\"}],\"name\":\"DepositSweepProposalSubmissionGasOffsetUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposalSubmitter\",\"type\":\"address\"}],\"name\":\"DepositSweepProposalSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"depositSweepProposalValidity\",\"type\":\"uint32\"}],\"name\":\"DepositSweepProposalValidityUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposalSubmitter\",\"type\":\"address\"}],\"name\":\"ProposalSubmitterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposalSubmitter\",\"type\":\"address\"}],\"name\":\"ProposalSubmitterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newReimbursementPool\",\"type\":\"address\"}],\"name\":\"ReimbursementPoolUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"}],\"name\":\"WalletManuallyUnlocked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposalSubmitter\",\"type\":\"address\"}],\"name\":\"addProposalSubmitter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositMinAge\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositRefundSafetyMargin\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepMaxSize\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepProposalSubmissionGasOffset\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositSweepProposalValidity\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractBridge\",\"name\":\"_bridge\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isProposalSubmitter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reimbursementPool\",\"outputs\":[{\"internalType\":\"contractReimbursementPool\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposalSubmitter\",\"type\":\"address\"}],\"name\":\"removeProposalSubmitter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32[]\",\"name\":\"walletMembersIDs\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"walletMemberIndex\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.WalletMemberContext\",\"name\":\"walletMemberContext\",\"type\":\"tuple\"}],\"name\":\"submitDepositSweepProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32[]\",\"name\":\"walletMembersIDs\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"walletMemberIndex\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.WalletMemberContext\",\"name\":\"walletMemberContext\",\"type\":\"tuple\"}],\"name\":\"submitDepositSweepProposalWithReimbursement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"}],\"name\":\"unlockWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_depositMinAge\",\"type\":\"uint32\"}],\"name\":\"updateDepositMinAge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_depositRefundSafetyMargin\",\"type\":\"uint32\"}],\"name\":\"updateDepositRefundSafetyMargin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_depositSweepMaxSize\",\"type\":\"uint16\"}],\"name\":\"updateDepositSweepMaxSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_depositSweepProposalSubmissionGasOffset\",\"type\":\"uint32\"}],\"name\":\"updateDepositSweepProposalSubmissionGasOffset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"_depositSweepProposalValidity\",\"type\":\"uint32\"}],\"name\":\"updateDepositSweepProposalValidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractReimbursementPool\",\"name\":\"_reimbursementPool\",\"type\":\"address\"}],\"name\":\"updateReimbursementPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"fundingTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"fundingOutputIndex\",\"type\":\"uint32\"}],\"internalType\":\"structWalletCoordinator.DepositKey[]\",\"name\":\"depositsKeys\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sweepTxFee\",\"type\":\"uint256\"}],\"internalType\":\"structWalletCoordinator.DepositSweepProposal\",\"name\":\"proposal\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes4\",\"name\":\"version\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"inputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"outputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"locktime\",\"type\":\"bytes4\"}],\"internalType\":\"structBitcoinTx.Info\",\"name\":\"fundingTx\",\"type\":\"tuple\"},{\"internalType\":\"bytes8\",\"name\":\"blindingFactor\",\"type\":\"bytes8\"},{\"internalType\":\"bytes20\",\"name\":\"walletPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes20\",\"name\":\"refundPubKeyHash\",\"type\":\"bytes20\"},{\"internalType\":\"bytes4\",\"name\":\"refundLocktime\",\"type\":\"bytes4\"}],\"internalType\":\"structWalletCoordinator.DepositExtraInfo[]\",\"name\":\"depositsExtraInfo\",\"type\":\"tuple[]\"}],\"name\":\"validateDepositSweepProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes20\",\"name\":\"\",\"type\":\"bytes20\"}],\"name\":\"walletLock\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"expiresAt\",\"type\":\"uint32\"},{\"internalType\":\"enumWalletCoordinator.WalletAction\",\"name\":\"cause\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"walletRegistry\",\"outputs\":[{\"internalType\":\"contractIWalletRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// WalletCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use WalletCoordinatorMetaData.ABI instead.
var WalletCoordinatorABI = WalletCoordinatorMetaData.ABI

// WalletCoordinator is an auto generated Go binding around an Ethereum contract.
type WalletCoordinator struct {
	WalletCoordinatorCaller     // Read-only binding to the contract
	WalletCoordinatorTransactor // Write-only binding to the contract
	WalletCoordinatorFilterer   // Log filterer for contract events
}

// WalletCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type WalletCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WalletCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WalletCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WalletCoordinatorSession struct {
	Contract     *WalletCoordinator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// WalletCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WalletCoordinatorCallerSession struct {
	Contract *WalletCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// WalletCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WalletCoordinatorTransactorSession struct {
	Contract     *WalletCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// WalletCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type WalletCoordinatorRaw struct {
	Contract *WalletCoordinator // Generic contract binding to access the raw methods on
}

// WalletCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WalletCoordinatorCallerRaw struct {
	Contract *WalletCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// WalletCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WalletCoordinatorTransactorRaw struct {
	Contract *WalletCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWalletCoordinator creates a new instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinator(address common.Address, backend bind.ContractBackend) (*WalletCoordinator, error) {
	contract, err := bindWalletCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinator{WalletCoordinatorCaller: WalletCoordinatorCaller{contract: contract}, WalletCoordinatorTransactor: WalletCoordinatorTransactor{contract: contract}, WalletCoordinatorFilterer: WalletCoordinatorFilterer{contract: contract}}, nil
}

// NewWalletCoordinatorCaller creates a new read-only instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*WalletCoordinatorCaller, error) {
	contract, err := bindWalletCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorCaller{contract: contract}, nil
}

// NewWalletCoordinatorTransactor creates a new write-only instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*WalletCoordinatorTransactor, error) {
	contract, err := bindWalletCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorTransactor{contract: contract}, nil
}

// NewWalletCoordinatorFilterer creates a new log filterer instance of WalletCoordinator, bound to a specific deployed contract.
func NewWalletCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*WalletCoordinatorFilterer, error) {
	contract, err := bindWalletCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorFilterer{contract: contract}, nil
}

// bindWalletCoordinator binds a generic wrapper to an already deployed contract.
func bindWalletCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WalletCoordinatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WalletCoordinator *WalletCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WalletCoordinator.Contract.WalletCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WalletCoordinator *WalletCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.WalletCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WalletCoordinator *WalletCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.WalletCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WalletCoordinator *WalletCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WalletCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WalletCoordinator *WalletCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WalletCoordinator *WalletCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.contract.Transact(opts, method, params...)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) Bridge() (common.Address, error) {
	return _WalletCoordinator.Contract.Bridge(&_WalletCoordinator.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) Bridge() (common.Address, error) {
	return _WalletCoordinator.Contract.Bridge(&_WalletCoordinator.CallOpts)
}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositMinAge(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositMinAge")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositMinAge() (uint32, error) {
	return _WalletCoordinator.Contract.DepositMinAge(&_WalletCoordinator.CallOpts)
}

// DepositMinAge is a free data retrieval call binding the contract method 0xbc20a080.
//
// Solidity: function depositMinAge() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositMinAge() (uint32, error) {
	return _WalletCoordinator.Contract.DepositMinAge(&_WalletCoordinator.CallOpts)
}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositRefundSafetyMargin(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositRefundSafetyMargin")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositRefundSafetyMargin() (uint32, error) {
	return _WalletCoordinator.Contract.DepositRefundSafetyMargin(&_WalletCoordinator.CallOpts)
}

// DepositRefundSafetyMargin is a free data retrieval call binding the contract method 0x31d8c622.
//
// Solidity: function depositRefundSafetyMargin() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositRefundSafetyMargin() (uint32, error) {
	return _WalletCoordinator.Contract.DepositRefundSafetyMargin(&_WalletCoordinator.CallOpts)
}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepMaxSize(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepMaxSize")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepMaxSize() (uint16, error) {
	return _WalletCoordinator.Contract.DepositSweepMaxSize(&_WalletCoordinator.CallOpts)
}

// DepositSweepMaxSize is a free data retrieval call binding the contract method 0xacaac705.
//
// Solidity: function depositSweepMaxSize() view returns(uint16)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepMaxSize() (uint16, error) {
	return _WalletCoordinator.Contract.DepositSweepMaxSize(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepProposalSubmissionGasOffset(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepProposalSubmissionGasOffset")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepProposalSubmissionGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalSubmissionGasOffset is a free data retrieval call binding the contract method 0x08553292.
//
// Solidity: function depositSweepProposalSubmissionGasOffset() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepProposalSubmissionGasOffset() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCaller) DepositSweepProposalValidity(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "depositSweepProposalValidity")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorSession) DepositSweepProposalValidity() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalValidity(&_WalletCoordinator.CallOpts)
}

// DepositSweepProposalValidity is a free data retrieval call binding the contract method 0xf74b37bd.
//
// Solidity: function depositSweepProposalValidity() view returns(uint32)
func (_WalletCoordinator *WalletCoordinatorCallerSession) DepositSweepProposalValidity() (uint32, error) {
	return _WalletCoordinator.Contract.DepositSweepProposalValidity(&_WalletCoordinator.CallOpts)
}

// IsProposalSubmitter is a free data retrieval call binding the contract method 0x213f7df7.
//
// Solidity: function isProposalSubmitter(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCaller) IsProposalSubmitter(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "isProposalSubmitter", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsProposalSubmitter is a free data retrieval call binding the contract method 0x213f7df7.
//
// Solidity: function isProposalSubmitter(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorSession) IsProposalSubmitter(arg0 common.Address) (bool, error) {
	return _WalletCoordinator.Contract.IsProposalSubmitter(&_WalletCoordinator.CallOpts, arg0)
}

// IsProposalSubmitter is a free data retrieval call binding the contract method 0x213f7df7.
//
// Solidity: function isProposalSubmitter(address ) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCallerSession) IsProposalSubmitter(arg0 common.Address) (bool, error) {
	return _WalletCoordinator.Contract.IsProposalSubmitter(&_WalletCoordinator.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) Owner() (common.Address, error) {
	return _WalletCoordinator.Contract.Owner(&_WalletCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) Owner() (common.Address, error) {
	return _WalletCoordinator.Contract.Owner(&_WalletCoordinator.CallOpts)
}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) ReimbursementPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "reimbursementPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) ReimbursementPool() (common.Address, error) {
	return _WalletCoordinator.Contract.ReimbursementPool(&_WalletCoordinator.CallOpts)
}

// ReimbursementPool is a free data retrieval call binding the contract method 0xc09975cd.
//
// Solidity: function reimbursementPool() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) ReimbursementPool() (common.Address, error) {
	return _WalletCoordinator.Contract.ReimbursementPool(&_WalletCoordinator.CallOpts)
}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCaller) ValidateDepositSweepProposal(opts *bind.CallOpts, proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "validateDepositSweepProposal", proposal, depositsExtraInfo)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorSession) ValidateDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	return _WalletCoordinator.Contract.ValidateDepositSweepProposal(&_WalletCoordinator.CallOpts, proposal, depositsExtraInfo)
}

// ValidateDepositSweepProposal is a free data retrieval call binding the contract method 0x7698fa97.
//
// Solidity: function validateDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, ((bytes4,bytes,bytes,bytes4),bytes8,bytes20,bytes20,bytes4)[] depositsExtraInfo) view returns(bool)
func (_WalletCoordinator *WalletCoordinatorCallerSession) ValidateDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, depositsExtraInfo []WalletCoordinatorDepositExtraInfo) (bool, error) {
	return _WalletCoordinator.Contract.ValidateDepositSweepProposal(&_WalletCoordinator.CallOpts, proposal, depositsExtraInfo)
}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorCaller) WalletLock(opts *bind.CallOpts, arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "walletLock", arg0)

	outstruct := new(struct {
		ExpiresAt uint32
		Cause     uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ExpiresAt = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Cause = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorSession) WalletLock(arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	return _WalletCoordinator.Contract.WalletLock(&_WalletCoordinator.CallOpts, arg0)
}

// WalletLock is a free data retrieval call binding the contract method 0x2c259d2b.
//
// Solidity: function walletLock(bytes20 ) view returns(uint32 expiresAt, uint8 cause)
func (_WalletCoordinator *WalletCoordinatorCallerSession) WalletLock(arg0 [20]byte) (struct {
	ExpiresAt uint32
	Cause     uint8
}, error) {
	return _WalletCoordinator.Contract.WalletLock(&_WalletCoordinator.CallOpts, arg0)
}

// WalletRegistry is a free data retrieval call binding the contract method 0xab7aa6ad.
//
// Solidity: function walletRegistry() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCaller) WalletRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WalletCoordinator.contract.Call(opts, &out, "walletRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WalletRegistry is a free data retrieval call binding the contract method 0xab7aa6ad.
//
// Solidity: function walletRegistry() view returns(address)
func (_WalletCoordinator *WalletCoordinatorSession) WalletRegistry() (common.Address, error) {
	return _WalletCoordinator.Contract.WalletRegistry(&_WalletCoordinator.CallOpts)
}

// WalletRegistry is a free data retrieval call binding the contract method 0xab7aa6ad.
//
// Solidity: function walletRegistry() view returns(address)
func (_WalletCoordinator *WalletCoordinatorCallerSession) WalletRegistry() (common.Address, error) {
	return _WalletCoordinator.Contract.WalletRegistry(&_WalletCoordinator.CallOpts)
}

// AddProposalSubmitter is a paid mutator transaction binding the contract method 0x8e323a99.
//
// Solidity: function addProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) AddProposalSubmitter(opts *bind.TransactOpts, proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "addProposalSubmitter", proposalSubmitter)
}

// AddProposalSubmitter is a paid mutator transaction binding the contract method 0x8e323a99.
//
// Solidity: function addProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorSession) AddProposalSubmitter(proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.AddProposalSubmitter(&_WalletCoordinator.TransactOpts, proposalSubmitter)
}

// AddProposalSubmitter is a paid mutator transaction binding the contract method 0x8e323a99.
//
// Solidity: function addProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) AddProposalSubmitter(proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.AddProposalSubmitter(&_WalletCoordinator.TransactOpts, proposalSubmitter)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) Initialize(opts *bind.TransactOpts, _bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "initialize", _bridge)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorSession) Initialize(_bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.Initialize(&_WalletCoordinator.TransactOpts, _bridge)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _bridge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) Initialize(_bridge common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.Initialize(&_WalletCoordinator.TransactOpts, _bridge)
}

// RemoveProposalSubmitter is a paid mutator transaction binding the contract method 0x3a4159c8.
//
// Solidity: function removeProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RemoveProposalSubmitter(opts *bind.TransactOpts, proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "removeProposalSubmitter", proposalSubmitter)
}

// RemoveProposalSubmitter is a paid mutator transaction binding the contract method 0x3a4159c8.
//
// Solidity: function removeProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorSession) RemoveProposalSubmitter(proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RemoveProposalSubmitter(&_WalletCoordinator.TransactOpts, proposalSubmitter)
}

// RemoveProposalSubmitter is a paid mutator transaction binding the contract method 0x3a4159c8.
//
// Solidity: function removeProposalSubmitter(address proposalSubmitter) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RemoveProposalSubmitter(proposalSubmitter common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RemoveProposalSubmitter(&_WalletCoordinator.TransactOpts, proposalSubmitter)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RenounceOwnership(&_WalletCoordinator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WalletCoordinator.Contract.RenounceOwnership(&_WalletCoordinator.TransactOpts)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x1c6710dd.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) SubmitDepositSweepProposal(opts *bind.TransactOpts, proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "submitDepositSweepProposal", proposal, walletMemberContext)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x1c6710dd.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorSession) SubmitDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposal(&_WalletCoordinator.TransactOpts, proposal, walletMemberContext)
}

// SubmitDepositSweepProposal is a paid mutator transaction binding the contract method 0x1c6710dd.
//
// Solidity: function submitDepositSweepProposal((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) SubmitDepositSweepProposal(proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposal(&_WalletCoordinator.TransactOpts, proposal, walletMemberContext)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0xc1790d77.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) SubmitDepositSweepProposalWithReimbursement(opts *bind.TransactOpts, proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "submitDepositSweepProposalWithReimbursement", proposal, walletMemberContext)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0xc1790d77.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorSession) SubmitDepositSweepProposalWithReimbursement(proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposalWithReimbursement(&_WalletCoordinator.TransactOpts, proposal, walletMemberContext)
}

// SubmitDepositSweepProposalWithReimbursement is a paid mutator transaction binding the contract method 0xc1790d77.
//
// Solidity: function submitDepositSweepProposalWithReimbursement((bytes20,(bytes32,uint32)[],uint256) proposal, (uint32[],uint256) walletMemberContext) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) SubmitDepositSweepProposalWithReimbursement(proposal WalletCoordinatorDepositSweepProposal, walletMemberContext WalletCoordinatorWalletMemberContext) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.SubmitDepositSweepProposalWithReimbursement(&_WalletCoordinator.TransactOpts, proposal, walletMemberContext)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.TransferOwnership(&_WalletCoordinator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.TransferOwnership(&_WalletCoordinator.TransactOpts, newOwner)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UnlockWallet(opts *bind.TransactOpts, walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "unlockWallet", walletPubKeyHash)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UnlockWallet(walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UnlockWallet(&_WalletCoordinator.TransactOpts, walletPubKeyHash)
}

// UnlockWallet is a paid mutator transaction binding the contract method 0x9039bc1c.
//
// Solidity: function unlockWallet(bytes20 walletPubKeyHash) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UnlockWallet(walletPubKeyHash [20]byte) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UnlockWallet(&_WalletCoordinator.TransactOpts, walletPubKeyHash)
}

// UpdateDepositMinAge is a paid mutator transaction binding the contract method 0x0dc1b8df.
//
// Solidity: function updateDepositMinAge(uint32 _depositMinAge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositMinAge(opts *bind.TransactOpts, _depositMinAge uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositMinAge", _depositMinAge)
}

// UpdateDepositMinAge is a paid mutator transaction binding the contract method 0x0dc1b8df.
//
// Solidity: function updateDepositMinAge(uint32 _depositMinAge) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositMinAge(_depositMinAge uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositMinAge(&_WalletCoordinator.TransactOpts, _depositMinAge)
}

// UpdateDepositMinAge is a paid mutator transaction binding the contract method 0x0dc1b8df.
//
// Solidity: function updateDepositMinAge(uint32 _depositMinAge) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositMinAge(_depositMinAge uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositMinAge(&_WalletCoordinator.TransactOpts, _depositMinAge)
}

// UpdateDepositRefundSafetyMargin is a paid mutator transaction binding the contract method 0x6ff913c3.
//
// Solidity: function updateDepositRefundSafetyMargin(uint32 _depositRefundSafetyMargin) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositRefundSafetyMargin(opts *bind.TransactOpts, _depositRefundSafetyMargin uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositRefundSafetyMargin", _depositRefundSafetyMargin)
}

// UpdateDepositRefundSafetyMargin is a paid mutator transaction binding the contract method 0x6ff913c3.
//
// Solidity: function updateDepositRefundSafetyMargin(uint32 _depositRefundSafetyMargin) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositRefundSafetyMargin(_depositRefundSafetyMargin uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositRefundSafetyMargin(&_WalletCoordinator.TransactOpts, _depositRefundSafetyMargin)
}

// UpdateDepositRefundSafetyMargin is a paid mutator transaction binding the contract method 0x6ff913c3.
//
// Solidity: function updateDepositRefundSafetyMargin(uint32 _depositRefundSafetyMargin) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositRefundSafetyMargin(_depositRefundSafetyMargin uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositRefundSafetyMargin(&_WalletCoordinator.TransactOpts, _depositRefundSafetyMargin)
}

// UpdateDepositSweepMaxSize is a paid mutator transaction binding the contract method 0x9f41e13b.
//
// Solidity: function updateDepositSweepMaxSize(uint16 _depositSweepMaxSize) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositSweepMaxSize(opts *bind.TransactOpts, _depositSweepMaxSize uint16) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositSweepMaxSize", _depositSweepMaxSize)
}

// UpdateDepositSweepMaxSize is a paid mutator transaction binding the contract method 0x9f41e13b.
//
// Solidity: function updateDepositSweepMaxSize(uint16 _depositSweepMaxSize) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositSweepMaxSize(_depositSweepMaxSize uint16) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepMaxSize(&_WalletCoordinator.TransactOpts, _depositSweepMaxSize)
}

// UpdateDepositSweepMaxSize is a paid mutator transaction binding the contract method 0x9f41e13b.
//
// Solidity: function updateDepositSweepMaxSize(uint16 _depositSweepMaxSize) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositSweepMaxSize(_depositSweepMaxSize uint16) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepMaxSize(&_WalletCoordinator.TransactOpts, _depositSweepMaxSize)
}

// UpdateDepositSweepProposalSubmissionGasOffset is a paid mutator transaction binding the contract method 0xb82bc399.
//
// Solidity: function updateDepositSweepProposalSubmissionGasOffset(uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositSweepProposalSubmissionGasOffset(opts *bind.TransactOpts, _depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositSweepProposalSubmissionGasOffset", _depositSweepProposalSubmissionGasOffset)
}

// UpdateDepositSweepProposalSubmissionGasOffset is a paid mutator transaction binding the contract method 0xb82bc399.
//
// Solidity: function updateDepositSweepProposalSubmissionGasOffset(uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositSweepProposalSubmissionGasOffset(_depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.TransactOpts, _depositSweepProposalSubmissionGasOffset)
}

// UpdateDepositSweepProposalSubmissionGasOffset is a paid mutator transaction binding the contract method 0xb82bc399.
//
// Solidity: function updateDepositSweepProposalSubmissionGasOffset(uint32 _depositSweepProposalSubmissionGasOffset) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositSweepProposalSubmissionGasOffset(_depositSweepProposalSubmissionGasOffset uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalSubmissionGasOffset(&_WalletCoordinator.TransactOpts, _depositSweepProposalSubmissionGasOffset)
}

// UpdateDepositSweepProposalValidity is a paid mutator transaction binding the contract method 0x23a1c307.
//
// Solidity: function updateDepositSweepProposalValidity(uint32 _depositSweepProposalValidity) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateDepositSweepProposalValidity(opts *bind.TransactOpts, _depositSweepProposalValidity uint32) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateDepositSweepProposalValidity", _depositSweepProposalValidity)
}

// UpdateDepositSweepProposalValidity is a paid mutator transaction binding the contract method 0x23a1c307.
//
// Solidity: function updateDepositSweepProposalValidity(uint32 _depositSweepProposalValidity) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateDepositSweepProposalValidity(_depositSweepProposalValidity uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalValidity(&_WalletCoordinator.TransactOpts, _depositSweepProposalValidity)
}

// UpdateDepositSweepProposalValidity is a paid mutator transaction binding the contract method 0x23a1c307.
//
// Solidity: function updateDepositSweepProposalValidity(uint32 _depositSweepProposalValidity) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateDepositSweepProposalValidity(_depositSweepProposalValidity uint32) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateDepositSweepProposalValidity(&_WalletCoordinator.TransactOpts, _depositSweepProposalValidity)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorTransactor) UpdateReimbursementPool(opts *bind.TransactOpts, _reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.contract.Transact(opts, "updateReimbursementPool", _reimbursementPool)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorSession) UpdateReimbursementPool(_reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateReimbursementPool(&_WalletCoordinator.TransactOpts, _reimbursementPool)
}

// UpdateReimbursementPool is a paid mutator transaction binding the contract method 0x7b35b4e6.
//
// Solidity: function updateReimbursementPool(address _reimbursementPool) returns()
func (_WalletCoordinator *WalletCoordinatorTransactorSession) UpdateReimbursementPool(_reimbursementPool common.Address) (*types.Transaction, error) {
	return _WalletCoordinator.Contract.UpdateReimbursementPool(&_WalletCoordinator.TransactOpts, _reimbursementPool)
}

// WalletCoordinatorDepositMinAgeUpdatedIterator is returned from FilterDepositMinAgeUpdated and is used to iterate over the raw logs and unpacked data for DepositMinAgeUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositMinAgeUpdatedIterator struct {
	Event *WalletCoordinatorDepositMinAgeUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositMinAgeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositMinAgeUpdated)
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
		it.Event = new(WalletCoordinatorDepositMinAgeUpdated)
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
func (it *WalletCoordinatorDepositMinAgeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositMinAgeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositMinAgeUpdated represents a DepositMinAgeUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositMinAgeUpdated struct {
	DepositMinAge uint32
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDepositMinAgeUpdated is a free log retrieval operation binding the contract event 0xeda69d960fc800ae7f67255cd8b42b338c24e1a655053127fff8bf88a2480d3f.
//
// Solidity: event DepositMinAgeUpdated(uint32 depositMinAge)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositMinAgeUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositMinAgeUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositMinAgeUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositMinAgeUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositMinAgeUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositMinAgeUpdated is a free log subscription operation binding the contract event 0xeda69d960fc800ae7f67255cd8b42b338c24e1a655053127fff8bf88a2480d3f.
//
// Solidity: event DepositMinAgeUpdated(uint32 depositMinAge)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositMinAgeUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositMinAgeUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositMinAgeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositMinAgeUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositMinAgeUpdated", log); err != nil {
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

// ParseDepositMinAgeUpdated is a log parse operation binding the contract event 0xeda69d960fc800ae7f67255cd8b42b338c24e1a655053127fff8bf88a2480d3f.
//
// Solidity: event DepositMinAgeUpdated(uint32 depositMinAge)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositMinAgeUpdated(log types.Log) (*WalletCoordinatorDepositMinAgeUpdated, error) {
	event := new(WalletCoordinatorDepositMinAgeUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositMinAgeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator is returned from FilterDepositRefundSafetyMarginUpdated and is used to iterate over the raw logs and unpacked data for DepositRefundSafetyMarginUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator struct {
	Event *WalletCoordinatorDepositRefundSafetyMarginUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositRefundSafetyMarginUpdated)
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
		it.Event = new(WalletCoordinatorDepositRefundSafetyMarginUpdated)
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
func (it *WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositRefundSafetyMarginUpdated represents a DepositRefundSafetyMarginUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositRefundSafetyMarginUpdated struct {
	DepositRefundSafetyMargin uint32
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterDepositRefundSafetyMarginUpdated is a free log retrieval operation binding the contract event 0xcb0c11e4086618176ceab0afd46b7fa25688cf2a6566ed4845f1fb530571faee.
//
// Solidity: event DepositRefundSafetyMarginUpdated(uint32 depositRefundSafetyMargin)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositRefundSafetyMarginUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositRefundSafetyMarginUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositRefundSafetyMarginUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositRefundSafetyMarginUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositRefundSafetyMarginUpdated is a free log subscription operation binding the contract event 0xcb0c11e4086618176ceab0afd46b7fa25688cf2a6566ed4845f1fb530571faee.
//
// Solidity: event DepositRefundSafetyMarginUpdated(uint32 depositRefundSafetyMargin)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositRefundSafetyMarginUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositRefundSafetyMarginUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositRefundSafetyMarginUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositRefundSafetyMarginUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositRefundSafetyMarginUpdated", log); err != nil {
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

// ParseDepositRefundSafetyMarginUpdated is a log parse operation binding the contract event 0xcb0c11e4086618176ceab0afd46b7fa25688cf2a6566ed4845f1fb530571faee.
//
// Solidity: event DepositRefundSafetyMarginUpdated(uint32 depositRefundSafetyMargin)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositRefundSafetyMarginUpdated(log types.Log) (*WalletCoordinatorDepositRefundSafetyMarginUpdated, error) {
	event := new(WalletCoordinatorDepositRefundSafetyMarginUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositRefundSafetyMarginUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepMaxSizeUpdatedIterator is returned from FilterDepositSweepMaxSizeUpdated and is used to iterate over the raw logs and unpacked data for DepositSweepMaxSizeUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepMaxSizeUpdatedIterator struct {
	Event *WalletCoordinatorDepositSweepMaxSizeUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepMaxSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepMaxSizeUpdated)
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
		it.Event = new(WalletCoordinatorDepositSweepMaxSizeUpdated)
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
func (it *WalletCoordinatorDepositSweepMaxSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepMaxSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepMaxSizeUpdated represents a DepositSweepMaxSizeUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepMaxSizeUpdated struct {
	DepositSweepMaxSize uint16
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepMaxSizeUpdated is a free log retrieval operation binding the contract event 0xcc7b63f885d59808f386320218caf8e6bdaa095b61927dea4e59f9cf6954beab.
//
// Solidity: event DepositSweepMaxSizeUpdated(uint16 depositSweepMaxSize)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepMaxSizeUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositSweepMaxSizeUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepMaxSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepMaxSizeUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepMaxSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositSweepMaxSizeUpdated is a free log subscription operation binding the contract event 0xcc7b63f885d59808f386320218caf8e6bdaa095b61927dea4e59f9cf6954beab.
//
// Solidity: event DepositSweepMaxSizeUpdated(uint16 depositSweepMaxSize)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepMaxSizeUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepMaxSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepMaxSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepMaxSizeUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepMaxSizeUpdated", log); err != nil {
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

// ParseDepositSweepMaxSizeUpdated is a log parse operation binding the contract event 0xcc7b63f885d59808f386320218caf8e6bdaa095b61927dea4e59f9cf6954beab.
//
// Solidity: event DepositSweepMaxSizeUpdated(uint16 depositSweepMaxSize)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepMaxSizeUpdated(log types.Log) (*WalletCoordinatorDepositSweepMaxSizeUpdated, error) {
	event := new(WalletCoordinatorDepositSweepMaxSizeUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepMaxSizeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator is returned from FilterDepositSweepProposalSubmissionGasOffsetUpdated and is used to iterate over the raw logs and unpacked data for DepositSweepProposalSubmissionGasOffsetUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator struct {
	Event *WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated)
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
		it.Event = new(WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated)
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
func (it *WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated represents a DepositSweepProposalSubmissionGasOffsetUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated struct {
	DepositSweepProposalSubmissionGasOffset uint32
	Raw                                     types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepProposalSubmissionGasOffsetUpdated is a free log retrieval operation binding the contract event 0xafb79fa073ff248ba2f1cd9798891f4ee776a03c9f860eccdebcb1837305a197.
//
// Solidity: event DepositSweepProposalSubmissionGasOffsetUpdated(uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepProposalSubmissionGasOffsetUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepProposalSubmissionGasOffsetUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepProposalSubmissionGasOffsetUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositSweepProposalSubmissionGasOffsetUpdated is a free log subscription operation binding the contract event 0xafb79fa073ff248ba2f1cd9798891f4ee776a03c9f860eccdebcb1837305a197.
//
// Solidity: event DepositSweepProposalSubmissionGasOffsetUpdated(uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepProposalSubmissionGasOffsetUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepProposalSubmissionGasOffsetUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmissionGasOffsetUpdated", log); err != nil {
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

// ParseDepositSweepProposalSubmissionGasOffsetUpdated is a log parse operation binding the contract event 0xafb79fa073ff248ba2f1cd9798891f4ee776a03c9f860eccdebcb1837305a197.
//
// Solidity: event DepositSweepProposalSubmissionGasOffsetUpdated(uint32 depositSweepProposalSubmissionGasOffset)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepProposalSubmissionGasOffsetUpdated(log types.Log) (*WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated, error) {
	event := new(WalletCoordinatorDepositSweepProposalSubmissionGasOffsetUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmissionGasOffsetUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepProposalSubmittedIterator is returned from FilterDepositSweepProposalSubmitted and is used to iterate over the raw logs and unpacked data for DepositSweepProposalSubmitted events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmittedIterator struct {
	Event *WalletCoordinatorDepositSweepProposalSubmitted // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepProposalSubmitted)
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
		it.Event = new(WalletCoordinatorDepositSweepProposalSubmitted)
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
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepProposalSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepProposalSubmitted represents a DepositSweepProposalSubmitted event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalSubmitted struct {
	Proposal          WalletCoordinatorDepositSweepProposal
	ProposalSubmitter common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepProposalSubmitted is a free log retrieval operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepProposalSubmitted(opts *bind.FilterOpts, proposalSubmitter []common.Address) (*WalletCoordinatorDepositSweepProposalSubmittedIterator, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepProposalSubmitted", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepProposalSubmittedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepProposalSubmitted", logs: logs, sub: sub}, nil
}

// WatchDepositSweepProposalSubmitted is a free log subscription operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepProposalSubmitted(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepProposalSubmitted, proposalSubmitter []common.Address) (event.Subscription, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepProposalSubmitted", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepProposalSubmitted)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmitted", log); err != nil {
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

// ParseDepositSweepProposalSubmitted is a log parse operation binding the contract event 0xe17b4ed4828b9f018ca549cfd5f662e5551408093b1db16d94b9cd242e6df171.
//
// Solidity: event DepositSweepProposalSubmitted((bytes20,(bytes32,uint32)[],uint256) proposal, address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepProposalSubmitted(log types.Log) (*WalletCoordinatorDepositSweepProposalSubmitted, error) {
	event := new(WalletCoordinatorDepositSweepProposalSubmitted)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorDepositSweepProposalValidityUpdatedIterator is returned from FilterDepositSweepProposalValidityUpdated and is used to iterate over the raw logs and unpacked data for DepositSweepProposalValidityUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalValidityUpdatedIterator struct {
	Event *WalletCoordinatorDepositSweepProposalValidityUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorDepositSweepProposalValidityUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorDepositSweepProposalValidityUpdated)
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
		it.Event = new(WalletCoordinatorDepositSweepProposalValidityUpdated)
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
func (it *WalletCoordinatorDepositSweepProposalValidityUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorDepositSweepProposalValidityUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorDepositSweepProposalValidityUpdated represents a DepositSweepProposalValidityUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorDepositSweepProposalValidityUpdated struct {
	DepositSweepProposalValidity uint32
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterDepositSweepProposalValidityUpdated is a free log retrieval operation binding the contract event 0x085a81d11ac3cf6033d185eccdfc414b8a20da29a555406c068a0a1fdac8afea.
//
// Solidity: event DepositSweepProposalValidityUpdated(uint32 depositSweepProposalValidity)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterDepositSweepProposalValidityUpdated(opts *bind.FilterOpts) (*WalletCoordinatorDepositSweepProposalValidityUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "DepositSweepProposalValidityUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorDepositSweepProposalValidityUpdatedIterator{contract: _WalletCoordinator.contract, event: "DepositSweepProposalValidityUpdated", logs: logs, sub: sub}, nil
}

// WatchDepositSweepProposalValidityUpdated is a free log subscription operation binding the contract event 0x085a81d11ac3cf6033d185eccdfc414b8a20da29a555406c068a0a1fdac8afea.
//
// Solidity: event DepositSweepProposalValidityUpdated(uint32 depositSweepProposalValidity)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchDepositSweepProposalValidityUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorDepositSweepProposalValidityUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "DepositSweepProposalValidityUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorDepositSweepProposalValidityUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalValidityUpdated", log); err != nil {
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

// ParseDepositSweepProposalValidityUpdated is a log parse operation binding the contract event 0x085a81d11ac3cf6033d185eccdfc414b8a20da29a555406c068a0a1fdac8afea.
//
// Solidity: event DepositSweepProposalValidityUpdated(uint32 depositSweepProposalValidity)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseDepositSweepProposalValidityUpdated(log types.Log) (*WalletCoordinatorDepositSweepProposalValidityUpdated, error) {
	event := new(WalletCoordinatorDepositSweepProposalValidityUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "DepositSweepProposalValidityUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the WalletCoordinator contract.
type WalletCoordinatorInitializedIterator struct {
	Event *WalletCoordinatorInitialized // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorInitialized)
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
		it.Event = new(WalletCoordinatorInitialized)
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
func (it *WalletCoordinatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorInitialized represents a Initialized event raised by the WalletCoordinator contract.
type WalletCoordinatorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*WalletCoordinatorInitializedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorInitializedIterator{contract: _WalletCoordinator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorInitialized) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorInitialized)
				if err := _WalletCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseInitialized(log types.Log) (*WalletCoordinatorInitialized, error) {
	event := new(WalletCoordinatorInitialized)
	if err := _WalletCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WalletCoordinator contract.
type WalletCoordinatorOwnershipTransferredIterator struct {
	Event *WalletCoordinatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorOwnershipTransferred)
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
		it.Event = new(WalletCoordinatorOwnershipTransferred)
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
func (it *WalletCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorOwnershipTransferred represents a OwnershipTransferred event raised by the WalletCoordinator contract.
type WalletCoordinatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WalletCoordinatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorOwnershipTransferredIterator{contract: _WalletCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorOwnershipTransferred)
				if err := _WalletCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*WalletCoordinatorOwnershipTransferred, error) {
	event := new(WalletCoordinatorOwnershipTransferred)
	if err := _WalletCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorProposalSubmitterAddedIterator is returned from FilterProposalSubmitterAdded and is used to iterate over the raw logs and unpacked data for ProposalSubmitterAdded events raised by the WalletCoordinator contract.
type WalletCoordinatorProposalSubmitterAddedIterator struct {
	Event *WalletCoordinatorProposalSubmitterAdded // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorProposalSubmitterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorProposalSubmitterAdded)
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
		it.Event = new(WalletCoordinatorProposalSubmitterAdded)
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
func (it *WalletCoordinatorProposalSubmitterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorProposalSubmitterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorProposalSubmitterAdded represents a ProposalSubmitterAdded event raised by the WalletCoordinator contract.
type WalletCoordinatorProposalSubmitterAdded struct {
	ProposalSubmitter common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterProposalSubmitterAdded is a free log retrieval operation binding the contract event 0xe7005265f76a2d6482aaa3a0e969edd45868e3210ff126216a0425f83af1ef20.
//
// Solidity: event ProposalSubmitterAdded(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterProposalSubmitterAdded(opts *bind.FilterOpts, proposalSubmitter []common.Address) (*WalletCoordinatorProposalSubmitterAddedIterator, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "ProposalSubmitterAdded", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorProposalSubmitterAddedIterator{contract: _WalletCoordinator.contract, event: "ProposalSubmitterAdded", logs: logs, sub: sub}, nil
}

// WatchProposalSubmitterAdded is a free log subscription operation binding the contract event 0xe7005265f76a2d6482aaa3a0e969edd45868e3210ff126216a0425f83af1ef20.
//
// Solidity: event ProposalSubmitterAdded(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchProposalSubmitterAdded(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorProposalSubmitterAdded, proposalSubmitter []common.Address) (event.Subscription, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "ProposalSubmitterAdded", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorProposalSubmitterAdded)
				if err := _WalletCoordinator.contract.UnpackLog(event, "ProposalSubmitterAdded", log); err != nil {
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

// ParseProposalSubmitterAdded is a log parse operation binding the contract event 0xe7005265f76a2d6482aaa3a0e969edd45868e3210ff126216a0425f83af1ef20.
//
// Solidity: event ProposalSubmitterAdded(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseProposalSubmitterAdded(log types.Log) (*WalletCoordinatorProposalSubmitterAdded, error) {
	event := new(WalletCoordinatorProposalSubmitterAdded)
	if err := _WalletCoordinator.contract.UnpackLog(event, "ProposalSubmitterAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorProposalSubmitterRemovedIterator is returned from FilterProposalSubmitterRemoved and is used to iterate over the raw logs and unpacked data for ProposalSubmitterRemoved events raised by the WalletCoordinator contract.
type WalletCoordinatorProposalSubmitterRemovedIterator struct {
	Event *WalletCoordinatorProposalSubmitterRemoved // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorProposalSubmitterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorProposalSubmitterRemoved)
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
		it.Event = new(WalletCoordinatorProposalSubmitterRemoved)
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
func (it *WalletCoordinatorProposalSubmitterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorProposalSubmitterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorProposalSubmitterRemoved represents a ProposalSubmitterRemoved event raised by the WalletCoordinator contract.
type WalletCoordinatorProposalSubmitterRemoved struct {
	ProposalSubmitter common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterProposalSubmitterRemoved is a free log retrieval operation binding the contract event 0x4017c4f3f844d24f56ade93af205a7ed01bf7f858b90f4b62206be460af2d226.
//
// Solidity: event ProposalSubmitterRemoved(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterProposalSubmitterRemoved(opts *bind.FilterOpts, proposalSubmitter []common.Address) (*WalletCoordinatorProposalSubmitterRemovedIterator, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "ProposalSubmitterRemoved", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorProposalSubmitterRemovedIterator{contract: _WalletCoordinator.contract, event: "ProposalSubmitterRemoved", logs: logs, sub: sub}, nil
}

// WatchProposalSubmitterRemoved is a free log subscription operation binding the contract event 0x4017c4f3f844d24f56ade93af205a7ed01bf7f858b90f4b62206be460af2d226.
//
// Solidity: event ProposalSubmitterRemoved(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchProposalSubmitterRemoved(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorProposalSubmitterRemoved, proposalSubmitter []common.Address) (event.Subscription, error) {

	var proposalSubmitterRule []interface{}
	for _, proposalSubmitterItem := range proposalSubmitter {
		proposalSubmitterRule = append(proposalSubmitterRule, proposalSubmitterItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "ProposalSubmitterRemoved", proposalSubmitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorProposalSubmitterRemoved)
				if err := _WalletCoordinator.contract.UnpackLog(event, "ProposalSubmitterRemoved", log); err != nil {
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

// ParseProposalSubmitterRemoved is a log parse operation binding the contract event 0x4017c4f3f844d24f56ade93af205a7ed01bf7f858b90f4b62206be460af2d226.
//
// Solidity: event ProposalSubmitterRemoved(address indexed proposalSubmitter)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseProposalSubmitterRemoved(log types.Log) (*WalletCoordinatorProposalSubmitterRemoved, error) {
	event := new(WalletCoordinatorProposalSubmitterRemoved)
	if err := _WalletCoordinator.contract.UnpackLog(event, "ProposalSubmitterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorReimbursementPoolUpdatedIterator is returned from FilterReimbursementPoolUpdated and is used to iterate over the raw logs and unpacked data for ReimbursementPoolUpdated events raised by the WalletCoordinator contract.
type WalletCoordinatorReimbursementPoolUpdatedIterator struct {
	Event *WalletCoordinatorReimbursementPoolUpdated // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorReimbursementPoolUpdated)
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
		it.Event = new(WalletCoordinatorReimbursementPoolUpdated)
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
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorReimbursementPoolUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorReimbursementPoolUpdated represents a ReimbursementPoolUpdated event raised by the WalletCoordinator contract.
type WalletCoordinatorReimbursementPoolUpdated struct {
	NewReimbursementPool common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterReimbursementPoolUpdated is a free log retrieval operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterReimbursementPoolUpdated(opts *bind.FilterOpts) (*WalletCoordinatorReimbursementPoolUpdatedIterator, error) {

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "ReimbursementPoolUpdated")
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorReimbursementPoolUpdatedIterator{contract: _WalletCoordinator.contract, event: "ReimbursementPoolUpdated", logs: logs, sub: sub}, nil
}

// WatchReimbursementPoolUpdated is a free log subscription operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchReimbursementPoolUpdated(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorReimbursementPoolUpdated) (event.Subscription, error) {

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "ReimbursementPoolUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorReimbursementPoolUpdated)
				if err := _WalletCoordinator.contract.UnpackLog(event, "ReimbursementPoolUpdated", log); err != nil {
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

// ParseReimbursementPoolUpdated is a log parse operation binding the contract event 0x0e2d2343d31b085b7c4e56d1c8a6ec79f7ab07460386f1c9a1756239fe2533ac.
//
// Solidity: event ReimbursementPoolUpdated(address newReimbursementPool)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseReimbursementPoolUpdated(log types.Log) (*WalletCoordinatorReimbursementPoolUpdated, error) {
	event := new(WalletCoordinatorReimbursementPoolUpdated)
	if err := _WalletCoordinator.contract.UnpackLog(event, "ReimbursementPoolUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WalletCoordinatorWalletManuallyUnlockedIterator is returned from FilterWalletManuallyUnlocked and is used to iterate over the raw logs and unpacked data for WalletManuallyUnlocked events raised by the WalletCoordinator contract.
type WalletCoordinatorWalletManuallyUnlockedIterator struct {
	Event *WalletCoordinatorWalletManuallyUnlocked // Event containing the contract specifics and raw log

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
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCoordinatorWalletManuallyUnlocked)
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
		it.Event = new(WalletCoordinatorWalletManuallyUnlocked)
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
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCoordinatorWalletManuallyUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCoordinatorWalletManuallyUnlocked represents a WalletManuallyUnlocked event raised by the WalletCoordinator contract.
type WalletCoordinatorWalletManuallyUnlocked struct {
	WalletPubKeyHash [20]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWalletManuallyUnlocked is a free log retrieval operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) FilterWalletManuallyUnlocked(opts *bind.FilterOpts, walletPubKeyHash [][20]byte) (*WalletCoordinatorWalletManuallyUnlockedIterator, error) {

	var walletPubKeyHashRule []interface{}
	for _, walletPubKeyHashItem := range walletPubKeyHash {
		walletPubKeyHashRule = append(walletPubKeyHashRule, walletPubKeyHashItem)
	}

	logs, sub, err := _WalletCoordinator.contract.FilterLogs(opts, "WalletManuallyUnlocked", walletPubKeyHashRule)
	if err != nil {
		return nil, err
	}
	return &WalletCoordinatorWalletManuallyUnlockedIterator{contract: _WalletCoordinator.contract, event: "WalletManuallyUnlocked", logs: logs, sub: sub}, nil
}

// WatchWalletManuallyUnlocked is a free log subscription operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) WatchWalletManuallyUnlocked(opts *bind.WatchOpts, sink chan<- *WalletCoordinatorWalletManuallyUnlocked, walletPubKeyHash [][20]byte) (event.Subscription, error) {

	var walletPubKeyHashRule []interface{}
	for _, walletPubKeyHashItem := range walletPubKeyHash {
		walletPubKeyHashRule = append(walletPubKeyHashRule, walletPubKeyHashItem)
	}

	logs, sub, err := _WalletCoordinator.contract.WatchLogs(opts, "WalletManuallyUnlocked", walletPubKeyHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCoordinatorWalletManuallyUnlocked)
				if err := _WalletCoordinator.contract.UnpackLog(event, "WalletManuallyUnlocked", log); err != nil {
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

// ParseWalletManuallyUnlocked is a log parse operation binding the contract event 0x5ea2f0397ee52fe8e2e6460e92b334c996f3580dc76c022d41b290cb985f41e6.
//
// Solidity: event WalletManuallyUnlocked(bytes20 indexed walletPubKeyHash)
func (_WalletCoordinator *WalletCoordinatorFilterer) ParseWalletManuallyUnlocked(log types.Log) (*WalletCoordinatorWalletManuallyUnlocked, error) {
	event := new(WalletCoordinatorWalletManuallyUnlocked)
	if err := _WalletCoordinator.contract.UnpackLog(event, "WalletManuallyUnlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
