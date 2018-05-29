package libp2p

import (
	"context"
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
		t.Fatalf("expected: provider type [%s]\nactual:   provider type [%s]",
			provider.Type(), expectedType,
		)
	}
}

func TestProviderReturnsChannel(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	testName := "testname"

	provider, err := Connect(ctx, generateDeterministicNetworkConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	if _, err = provider.ChannelFor(testName); err != nil {
		t.Fatalf("expected: test to fail with [%v]\nactual:   failed with [%v]",
			nil, err,
		)
	}
}

func TestBroadcastChannel(t *testing.T) {
	t.Skip()

	ctx, cancel := newTestContext()
	defer cancel()

	tests := map[string]struct {
		name          string
		expectedError func(string) error
	}{
		"Send succeeds": {
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
			_, err := provider.ChannelFor(test.name)
			if !reflect.DeepEqual(test.expectedError(test.name), err) {
				t.Fatalf("expected test to fail with [%v] instead failed with [%v]",
					test.expectedError(test.name), err,
				)
			}
		})
	}
}

func TestNetworkConnect(t *testing.T) {
	t.Skip()

	ctx, cancel := newTestContext()
	defer cancel()

	proxies, err := buildTestProxies(ctx, t, 2)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: fix this
	connectNetworks(ctx, t, proxies)

	// TODO: have providers send messages to each other
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func generateDeterministicNetworkConfig(t *testing.T) *Config {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := &identity{id: p.ID, privKey: p.PrivKey, pubKey: p.PubKey}
	return &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}, identity: pi}
}

func testProvider(ctx context.Context, t *testing.T) (*provider, error) {
	testConfig := generateDeterministicNetworkConfig(t)

	host, err := discoverAndListen(ctx, testConfig)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, host)
	if err != nil {
		return nil, err
	}

	return &provider{cm: cm, host: host}, nil
}

func buildTestProxies(ctx context.Context, t *testing.T, num int) ([]*provider, error) {
	proxies := make([]*provider, num)
	for i := 0; i < num; i++ {
		proxy, err := testProvider(ctx, t)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

func connectNetworks(ctx context.Context, t *testing.T, proxies []*provider) {
	var waitGroup sync.WaitGroup

	for i, proxy := range proxies {
		// connect to all other peers, proxies after i+1, for good connectivity
		for _, peer := range proxies[i+1:] {
			waitGroup.Add(1)
			proxy.host.Peerstore().AddAddr(
				peer.host.ID(),
				peer.host.Network().ListenAddresses()[0],
				peerstore.PermanentAddrTTL,
			)
			_, err := proxy.host.Network().DialPeer(ctx, peer.host.ID())
			if err != nil {
				t.Fatal(err)
			}
			waitGroup.Done()
		}
	}
	waitGroup.Wait()
}
