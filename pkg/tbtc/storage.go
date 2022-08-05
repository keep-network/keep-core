package tbtc

import (
	"github.com/keep-network/keep-common/pkg/persistence"
	"sync"
)

// walletStorage is the component that persists data of the wallets managed
// by the given node. All functions of the storage are safe for concurrent use.
type walletStorage struct {
	mutex       sync.Mutex
	persistence persistence.Handle
	wallets     map[string]*wallet
}

// newWalletStorage constructs a new instance of the walletStorage.
func newWalletStorage(persistence persistence.Handle) *walletStorage {
	return &walletStorage{
		persistence: persistence,
		wallets:     make(map[string]*wallet),
	}
}
