package key

import (
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

// NetworkKey represents static peer's key equal to the delegate key associated
// with an on-chain stake. It is used to authenticate the peer and for message
// attributability - each message leaving the peer is signed with its static
// key.
type NetworkKey interface {
	PrivateKey() libp2pcrypto.PrivKey
}
