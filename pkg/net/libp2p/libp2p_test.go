package libp2p

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	testutils "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func TestProviderReturnsType(t *testing.T) {
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	expectedType := "libp2p"
	provider, err := Connect(
		ctx, generateDeterministicNetworkConfig(t),
	)
	if err != nil {
		t.Fatal(err)
	}

	if provider.Type() != expectedType {
		t.Fatalf("%s %s", provider.Type(), expectedType)
	}
}

func TestProviderReturnsChannel(t *testing.T) {
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	tests := map[string]struct {
		name string
	}{
		"channel for name does not exist": {
			name: "",
		},
		"channel for name does exist": {
			name: "testchannel",
		},
	}

	provider, err := Connect(ctx, generateDeterministicNetworkConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			broadcastChannel := provider.ChannelFor(test.name)
			// if broadcastChannel == net.BroadcastChannel {
			// 	t.Fatalf("failed to return a valid broadcast channel")
			// }
			fmt.Println(broadcastChannel)
		})
	}
}

func TestNetworkConnect(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	proxies := buildTestProxies(ctx, t, 2)
	connectNetworks(ctx, t, proxies)
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func generateDeterministicNetworkConfig(t *testing.T) *Config {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := &identity{id: p.ID, privKey: p.PrivKey, pubKey: p.PubKey}
	return &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}, identity: pi}
}

func testProxy(ctx context.Context, t *testing.T) *proxy {
	testConfig := generateDeterministicNetworkConfig(t)
	proxy, err := newProxy(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	return proxy
}

func buildTestProxies(ctx context.Context, t *testing.T, num int) []*proxy {
	proxies := make([]*proxy, num)
	for i := 0; i < num; i++ {
		proxy := testProxy(ctx, t)
		proxies = append(proxies, proxy)
	}
	return proxies
}

func connectNetworks(ctx context.Context, t *testing.T, proxies []*proxy) {
	var wg sync.WaitGroup

	for i, proxy := range proxies {
		// connect to all other peers, proxies after i+1, for good connectivity
		for _, peer := range proxies[i+1:] {
			wg.Add(1)
			proxy.host.Peerstore().AddAddr(
				peer.host.ID(),
				peer.host.Network().ListenAddresses()[0],
				peerstore.PermanentAddrTTL,
			)
			_, err := proxy.host.Network().DialPeer(ctx, peer.host.ID())
			if err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}
	}
	wg.Wait()
}
