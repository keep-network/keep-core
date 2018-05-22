package libp2p

import (
	"context"
	"testing"
	"time"

	testutils "github.com/libp2p/go-testutil"
	ma "github.com/multiformats/go-multiaddr"
)

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func generateDeterministicNetworkConfig(t *testing.T) *Config {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := &identity{id: p.ID, privKey: p.PrivKey, pubKey: p.PubKey}
	return &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}, identity: pi}
}

func TestConnect(t *testing.T) {
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	testConfig := generateDeterministicNetworkConfig(t)
	provider, err := Connect(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	// provider
	// ps.AddAddrs(p.ID, n.ListenAddresses(), pstore.PermanentAddrTTL)
	// fmt.Printf("%+v", provider)
}
