package libp2p

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	testutils "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func TestProviderReturnsType(t *testing.T) {
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
	ctx, cancel := newTestContext()
	defer cancel()

	tests := map[string]struct {
		name          string
		expectedError func(string) error
	}{
		"channel for name does not exist": {
			name: "",
			expectedError: func(name string) error {
				return fmt.Errorf("invalid channel name")
			},
		},
		"channel for name does exist": {
			name: "testchannel",
			expectedError: func(name string) error {
				return nil
			},
		},
	}

	provider, err := Connect(ctx, generateDeterministicNetworkConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			broadcastChannel, err := provider.ChannelFor(test.name)
			if !reflect.DeepEqual(test.expectedError(test.name), err) {
				t.Fatalf("expected test to fail with [%v] instead failed with [%v]",
					test.expectedError(test.name), err,
				)
			}
			if err == nil {
				// TODO: Test that broadcastChannel does things
			}
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
	// TODO: fix this
	connectNetworks(ctx, t, proxies)

	// TODO: have proxies send messages to each other
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
