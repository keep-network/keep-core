package net

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/identity"
	host "github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

type Peer struct {
	ID    identity.Identity
	Store pstore.Peerstore
	ph    host.Host
}

func NewPeer(randseed int64, filepath string) *Peer {
	pi, err := identity.LoadOrGenerateIdentity(randseed, filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate Identity with error %s", err))
	}
	ps, err := pi.AddIdentityToStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to add Identity to PeerStore with error %s", err))
	}

	return &Peer{ID: pi, Store: ps}
}
