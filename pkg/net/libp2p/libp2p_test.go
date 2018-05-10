package libp2p

import (
	"context"
	"testing"
	"time"

	pstore "github.com/libp2p/go-libp2p-peerstore"
	testutils "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func testServer() *proxy {
	return nil
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func genNetworkConfig(t *testing.T, ctx context.Context) (pstore.Peerstore, *Config) {
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
	testConfig := &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}}
	return ps, testConfig
}

func TestConnect(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	_, testConfig := genNetworkConfig(t, ctx)
	_, err := Connect(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	// ps.AddAddrs(p.ID, n.ListenAddresses(), pstore.PermanentAddrTTL)
	// fmt.Printf("%+v", provider)
}
