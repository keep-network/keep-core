// Package ethutil provides utilities used for dealing with Ethereum concerns in
// the context of implementing cross-chain interfaces defined in pkg/chain.
package ethutil

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// AddressFromHex converts the passed string to a common.Address and returns it,
// unless it is not a valid address, in which case it returns an error. Compare
// to common.HexToAddress, which assumes the address is valid and does not
// provide for an error return.
func AddressFromHex(hex string) (common.Address, error) {
	if common.IsHexAddress(hex) {
		return common.HexToAddress(hex), nil
	}

	return common.Address{}, fmt.Errorf("[%v] is not a valid Ethereum address", hex)
}

// DecryptKeyFile reads in a key file and uses the password to decrypt it.
func DecryptKeyFile(keyFile, password string) (*keystore.Key, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read KeyFile %s [%v]", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt %s [%v]", keyFile, err)
	}
	return key, nil
}
