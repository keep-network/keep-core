package net

import (
	"fmt"

	host "github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	peer "github.com/rargulati/go-libp2p-peer"
)

type Authenticator interface {
	Sign(data []byte) ([]byte, error)
	Sign(data interface{}) ([]byte, error)
	Verify(data []byte, sig []byte, peerID peer.ID, pubKey []byte) bool
}

type Peer struct {
	ID    Identity
	Store pstore.Peerstore
	ph    host.Host
}

func NewPeer(randseed int64, filepath string) *Peer {
	pi, err := LoadOrGenerateIdentity(randseed, filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate Identity with error %s", err))
	}
	ps, err := pi.AddIdentityToStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to add Identity to PeerStore with error %s", err))
	}

	return &Peer{ID: pi, Store: ps}
}

func (p *Peer) Sign(data []byte) ([]byte, error)
func (p *Peer) Sign(data interface{}) ([]byte, error)
func (p *Peer) Verify(data []byte, sig []byte, peerID peer.ID, pubkey []byte) bool
