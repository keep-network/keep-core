// Package ethutil provides utilities used for dealing with Ethereum concerns in
// the context of implementing cross-chain interfaces defined in pkg/chain.
package ethutil

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

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
