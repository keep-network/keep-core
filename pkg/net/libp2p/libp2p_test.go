package libp2p

import (
	"context"
	"testing"
	"time"

	testutils "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func TestConnect(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	testProxy := testProxy(ctx, t)
	testProxy.Type()
	// test proxy is a type libp2p
	// test proxy returns a channel

}

func TestNetworkConnect(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	_ = testProxy(ctx, t)
	_ = testProxy(ctx, t)

}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func testProxy(ctx context.Context, t *testing.T) *Proxy {
	testConfig := generateDeterministicNetworkConfig(t)
	proxy, err := newProxy(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	return proxy
}

func connectNetworks(ctx context.Context, t *testing.T, proxies []*Proxy) {
	// var wg sync.WaitGroup

	// for i, proxy1 := range proxies {
	// 	for j, proxy2 := range proxies[i+1:] {
	// 		wg.Add(1)
	// 		// proxy1.host.Peerstore().AddAddr(proxy1.host.ID(), proxy2.host., ttl time.Duration)
	// 	}
	// }
}

func generateDeterministicNetworkConfig(t *testing.T) *Config {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := &identity{id: p.ID, privKey: p.PrivKey, pubKey: p.PubKey}
	return &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}, identity: pi}
}
