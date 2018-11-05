package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pschlump/godebug"
)

func (kg *keepGroup) SetGroupSizeThreshold(groupSize, threshold int) (*types.Transaction, error) {
	bigGroupSize := big.NewInt(int64(groupSize))
	bigThreshold := big.NewInt(int64(threshold))
	kg.transactorOpts.GasLimit = 4712388
	kg.transactorOpts.GasPrice = big.NewInt(int64(100000000000))
	// func (_KeepGroupImplV1 *KeepGroupImplV1Transactor) SetGroupSizeThreshold(opts *bind.TransactOpts, _groupSize *big.Int, _groupThreshold *big.Int) (*types.Transaction, error) {
	fmt.Printf("tOpts=%s\n", godebug.SVarI(kg.transactorOpts))
	return kg.transactor.SetGroupSizeThreshold(kg.transactorOpts, bigGroupSize, bigThreshold)
}

func (kg *keepGroup) NumberOfGroups() (int, error) {
	// func (_KeepGroupImplV1 *KeepGroupImplV1Caller) NumberOfGroups(opts *bind.CallOpts) (*big.Int, error) {
	nog, err := kg.caller.NumberOfGroups(kg.callerOpts)
	if err != nil {
		return 0, err
	}
	return int(nog.Int64()), nil
}

func (kg *keepGroup) GetGroupIndex(groupPubKey []byte) (int, error) {
	groupPubKeyArray, err := toByte32(groupPubKey)
	if err != nil {
		return -1, err
	}
	// func (_KeepGroupImplV1 *KeepGroupImplV1Caller) GetGroupIndex(opts *bind.CallOpts, groupPubKey [32]byte) (*big.Int, error) {
	idx, err := kg.caller.GetGroupIndex(kg.callerOpts, groupPubKeyArray)
	if err != nil {
		return -1, err
	}
	return int(idx.Int64()), nil
}

func (kg *keepGroup) GetGroupPubKey(idx int) ([]byte, error) {
	// func (_KeepGroupImplV1 *KeepGroupImplV1Session) GetGroupPubKey(groupIndex *big.Int) ([32]byte, error) {
	bigIdx := big.NewInt(int64(idx))
	pkArray, err := kg.caller.GetGroupPubKey(kg.callerOpts, bigIdx)
	if err != nil {
		return []byte{}, err
	}
	return pkArray[:], nil
}
