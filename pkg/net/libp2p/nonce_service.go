package libp2p

import (
	"time"
)

// TODO: Value from on-chain
const nonceServiceTimeout = time.Second * 60

type nonceService struct {
	identity *identity

	// nonce = hash(n_1 || n_2).toInt
	initial uint64
	latest  uint64
	max     uint64

	used map[uint64]bool
}

func NewNonceService(identity *identity) *nonceService {
	ns := &nonceService{
		identity: identity,
		used:     make(map[uint64]bool),
	}
	return ns
}
