package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (ec *ethereumChain) SetGroupSizeThreshold(groupSize, threshold int) (*types.Transaction, error) {
	// func (kg *keepGroup) SetGroupSizeThreshold(groupSize, threshold int) (*types.Transaction, error) {
	return ec.keepGroupContract.SetGroupSizeThreshold(groupSize, threshold)
}

/*
func (kg *keepGroup) CreateGroup(
	groupPubKey []byte,
) (*types.Transaction, error) {
*/
func (ec *ethereumChain) CreateGroup(groupPubKey []byte) (*types.Transaction, error) {
	// func (kg *keepGroup) SetGroupSizeThreshold(groupSize, threshold int) (*types.Transaction, error) {
	return ec.keepGroupContract.CreateGroup(groupPubKey)
}

func (ec *ethereumChain) NumberOfGroups() (int, error) {
	// func (kg *keepGroup) NumberOfGroups() (int, error) {
	n, err := ec.keepGroupContract.NumberOfGroups()
	if err != nil {
		return 0, fmt.Errorf("error calling NumberOfGroups: [%v]", err)
	}
	return n, nil
}

func (ec *ethereumChain) GetGroupIndex(groupPubKey []byte) (int, error) {
	// func (kg *keepGroup) GetGroupIndex(pubKey []byte) (int, error) {
	idx, err := ec.keepGroupContract.GetGroupIndex(groupPubKey)
	if err != nil {
		return -1, fmt.Errorf("error calling GetGroupIndex: [%v]", err)
	}
	return idx, nil
}

func (ec *ethereumChain) GetGroupPubKey(idx int) ([]byte, error) {
	// func (kg *keepGroup) GetGroupPubKey(idx int) ([]byte, error) {
	pk, err := ec.keepGroupContract.GetGroupPubKey(idx)
	if err != nil {
		return []byte{}, fmt.Errorf("error calling GetGroupPubKey: [%v]", err)
	}
	return pk, nil
}

type TestConfigType struct {
	URL             string
	URLRPC          string
	ContractAddress map[string]string // set of addresses for contracts
	ProxyFor        map[string]string // Specify contract names that are proxies for other contracts
	Address         string
	KeyFile         string
	KeyFilePassword string
}

var TestConfig TestConfigType
var EthConn *ethereumChain
