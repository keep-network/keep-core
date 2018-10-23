package handshake

import peer "github.com/libp2p/go-libp2p-peer"

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
const ID = "/keep/handshake/1.0.0"

// Transport is stream security transport. It provides no
// confidentiality, but provides guarantees of authentication and integrity.
type Transport struct {
	id peer.ID
}

// New constructs a new insecure transport.
func New(id peer.ID) *Transport {
	return &Transport{
		id: id,
	}
}
