package tbtc

import (
	"sync"
)

// heartbeatFailureCounter holds counters keeping track of consecutive
// heartbeat failures. Each wallet has a separate counter. The key used in
// the map is the uncompressed public key (with 04 prefix) of the wallet.
type heartbeatFailureCounter struct {
	mutex    sync.Mutex
	counters map[string]uint
}

func newHeartbeatFailureCounter() *heartbeatFailureCounter {
	return &heartbeatFailureCounter{
		counters: make(map[string]uint),
	}
}

func (hfc *heartbeatFailureCounter) increment(walletPublicKey string) {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	hfc.counters[walletPublicKey]++

}

func (hfc *heartbeatFailureCounter) reset(walletPublicKey string) {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	hfc.counters[walletPublicKey] = 0
}

func (hfc *heartbeatFailureCounter) get(walletPublicKey string) uint {
	hfc.mutex.Lock()
	defer hfc.mutex.Unlock()

	return hfc.counters[walletPublicKey]
}
