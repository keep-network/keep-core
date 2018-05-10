package libp2p

import (
	"context"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	testutils "github.com/libp2p/go-testutil"
)

func testServer() *proxy {
	return nil
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func GenNetworkIdentity(t *testing.T, ctx context.Context) pstore.Peerstore {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := peerIdentifier{id: p.ID, sk: p.PrivKey}
	// n, err := NewNetwork(ctx, []ma.Multiaddr{p.Addr}, p.ID, ps, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	ps, err := addIdentityToStore(pi)
	if err != nil {
		t.Fatal(err)
	}
	// ps.AddAddrs(p.ID, n.ListenAddresses(), pstore.PermanentAddrTTL)
	return ps
}

func TestConnect(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	testConfig := &net.Config{Port: 8080}
	_, err := Connect(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%+v", provider)
}
