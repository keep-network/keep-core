package libp2p

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

// TODO: Value from on-chain
const nonceServiceTimeout = time.Second * 60

type nonceService struct {
	identity *identity

	// nonce = hash(n_1 || n_2).toInt
	initial uint64
	latest  uint64

	used map[uint64]bool
}

func NewNonceService(identity *identity) *nonceService {
	ns := &nonceService{
		identity: identity,
		used:     make(map[uint64]bool),
	}
	return ns
}

func nonceFromBytes(proposedNonceBytes []byte) uint64 {
	return binary.LittleEndian.Uint64(proposedNonceBytes)
}

func (ns *nonceService) consumeNonce(nonce uint64) error {
	if _, ok := ns.used[nonce]; ok {
		return fmt.Errorf("bad times")
	}
	ns.used[nonce] = true
	return nil
}

func (ns *nonceService) incrememntNonce() uint64 {
	counter := atomic.AddUint64(&ns.latest, 1)
	return counter
}
