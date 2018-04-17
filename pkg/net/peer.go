package net

import (
	"context"
	"fmt"
	"strings"

	"github.com/keep-network/keep-core/pkg/net/identity"
	addrutil "github.com/libp2p/go-addr-util"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

type Authenticator interface {
	Sign(data []byte) ([]byte, error)
	Sign(data interface{}) ([]byte, error)
	Verify(data []byte, sig []byte, peerID peer.ID, pubKey []byte) bool
}

type Peer struct {
	ID    identity.Identity
	Store pstore.Peerstore
	ph    host.Host
}

type Connector interface {
	Connect(ctx context.Context, port int) error
	Bootstrap(ctx context.Context) error
	Close() error
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

func (p *Peer) Connect(ctx context.Context, port int) error {
	// Convert available network ifaces to listen on into multiaddrs
	addrs, err := getListenAdresses(port)
	if err != nil {
		return err
	}

	// TODO: should we limit which interfaces we attempt to listen to?
	nonLocalAddrs := make([]ma.Multiaddr, 0)
	for _, addr := range addrs {
		stringily := fmt.Sprintf("%v", addr)
		// TODO: limited to this cidr in testing
		if strings.Contains(stringily, "10.240") {
			nonLocalAddrs = append(nonLocalAddrs, addr)
		}
	}

	p.ph, err = p.buildPeerHost(nonLocalAddrs)
	if err != nil {
		return err
	}

	// Ok, now we're ready to listen
	if err := p.ph.Network().Listen(addrs...); err != nil {
		return err
	}

	return nil
}

// TODO: Allow for user-scoped listeners to either override this or union with this.
func getListenAdresses(port int) ([]ma.Multiaddr, error) {
	// TODO: figure out go-libp2p-interface-pnet.Protector and go-libp2p-pnet.NewProtector - later
	ia, err := addrutil.InterfaceAddresses()
	if err != nil {
		return nil, err
	}
	addrs := make([]ma.Multiaddr, len(ia), len(ia))
	for _, addr := range ia {
		portAddr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", port))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(portAddr))
	}
	return addrs, nil
}

func (p *Peer) buildPeerHost(listenAddrs []ma.Multiaddr) (host.Host, error)

func (p *Peer) Sign(data []byte) ([]byte, error)
func (p *Peer) Sign(data interface{}) ([]byte, error)
func (p *Peer) Verify(data []byte, sig []byte, peerID peer.ID, pubkey []byte) bool
