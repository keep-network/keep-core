package node

import (
	"context"

	cid "github.com/ipfs/go-cid"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	routing "github.com/libp2p/go-libp2p-routing"
)

// Implements the routing interface. Quick and dirty.
// See https://github.com/libp2p/go-libp2p-routing/blob/master/routing.go
type nilRouter struct{}

// satisfy routing.ContentRouting interface
func (nr *nilRouter) Provide(_ context.Context, _ *cid.Cid, _ bool) error {
	return nil
}

func (nr *nilRouter) FindProvidersAsync(_ context.Context, _ *cid.Cid, _ int) <-chan ps.PeerInfo {
	return nil
}

// satisfy routing.PeerRouting interface
func (nr *nilRouter) FindPeer(_ context.Context, _ peer.ID) (ps.PeerInfo, error) {
	return ps.PeerInfo{}, nil
}

// satisfy routing.ValueStore interface
func (nr *nilRouter) PutValue(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (nr *nilRouter) GetValue(_ context.Context, _ string) ([]byte, error) {
	return nil, nil
}

func (nr *nilRouter) GetValues(_ context.Context, _ string, _ int) ([]routing.RecvdVal, error) {
	return nil, nil
}

// hack to satisfy the interface concerns of routing.IpfsRouting
// TODO: gross, can we untangle this? Necessary for now...
var _ routing.IpfsRouting = &nilRouter{}

func (nr *nilRouter) Bootstrap(_ context.Context) error {
	return nil
}

// satisfy routing.PubKeyFetcher interface
func (nr *nilRouter) GetPublicKey(context.Context, peer.ID) (ci.PubKey, error) {
	return nil, nil
}

func NewNilRouter() *nilRouter {
	return &nilRouter{}
}
