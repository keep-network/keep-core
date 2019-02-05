package libp2p

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

func TestProviderReturnsType(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	privKey, _, err := key.GenerateStaticNetworkKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	expectedType := "libp2p"
	provider, err := Connect(
		ctx,
		generateDeterministicNetworkConfig(t),
		privKey,
		local.NewStakeMonitor(big.NewInt(200)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if provider.Type() != expectedType {
		t.Fatalf("expected: provider type [%s]\nactual:   provider type [%s]",
			provider.Type(),
			expectedType,
		)
	}
}

func TestProviderReturnsChannel(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	testName := "testname"

	privKey, _, err := key.GenerateStaticNetworkKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	provider, err := Connect(
		ctx,
		generateDeterministicNetworkConfig(t),
		privKey,
		local.NewStakeMonitor(big.NewInt(200)),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = provider.ChannelFor(testName); err != nil {
		t.Fatalf("expected: test to fail with [%v]\nactual:   failed with [%v]",
			nil,
			err,
		)
	}
}

func TestSendReceive(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	var (
		config          = generateDeterministicNetworkConfig(t)
		name            = "testchannel"
		expectedPayload = "some text"
	)

	privKey, _, err := key.GenerateStaticNetworkKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	identity, err := createIdentity(privKey)
	if err != nil {
		t.Fatal(err)
	}

	provider, err := Connect(
		ctx,
		config,
		privKey,
		local.NewStakeMonitor(big.NewInt(200)),
	)
	if err != nil {
		t.Fatal(err)
	}
	broadcastChannel, err := provider.ChannelFor(name)
	if err != nil {
		t.Fatal(err)
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &testMessage{} },
	); err != nil {
		t.Fatal(err)
	}

	if err := broadcastChannel.Send(
		&testMessage{Sender: identity, Payload: expectedPayload},
	); err != nil {
		t.Fatal(err)
	}

	recvChan := make(chan net.Message)
	if err := broadcastChannel.Recv(net.HandleMessageFunc{
		Type: "test",
		Handler: func(msg net.Message) error {
			recvChan <- msg
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case msg := <-recvChan:
			testPayload, ok := msg.Payload().(*testMessage)
			if !ok {
				t.Fatalf(
					"expected: payload type string\nactual:   payload type [%v]",
					testPayload,
				)
			}

			if expectedPayload != testPayload.Payload {
				t.Fatalf(
					"expected: message payload [%s]\ngot:   payload [%s]",
					expectedPayload,
					testPayload.Payload,
				)
			}
			return
		case <-ctx.Done():
			t.Fatal(err)
		}
	}
}

type protocolIdentifier struct {
	id string
}

type testMessage struct {
	Sender    *identity
	Recipient *identity
	Payload   string
}

func (m *testMessage) Type() string {
	return "test/unmarshaler"
}

func (m *testMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *testMessage) Unmarshal(bytes []byte) error {
	var message testMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}
	m.Sender = message.Sender
	m.Recipient = message.Recipient
	m.Payload = message.Payload

	return nil
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

func newTestIdentity() (*identity, error) {
	privKey, _, err := key.GenerateStaticNetworkKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return createIdentity(privKey)
}

func generateDeterministicNetworkConfig(t *testing.T) Config {
	return Config{Port: 8080}
}

func testProvider(ctx context.Context, t *testing.T) (*provider, error) {
	identity, err := newTestIdentity()
	if err != nil {
		return nil, err
	}

	host, err := discoverAndListen(
		ctx,
		identity,
		8080,
		local.NewStakeMonitor(big.NewInt(200)),
	)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, identity, host)
	if err != nil {
		return nil, err
	}

	return &provider{channelManagr: cm, host: host}, nil
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

func connectNetworks(ctx context.Context, t *testing.T, providers []*provider) {
	var waitGroup sync.WaitGroup

	for i, provider := range providers {
		// connect to all other peers, proxies after i+1, for good connectivity
		for _, peer := range providers[i+1:] {
			waitGroup.Add(1)
			provider.host.Peerstore().AddAddr(
				peer.host.ID(),
				peer.host.Network().ListenAddresses()[0],
				peerstore.PermanentAddrTTL,
			)
			_, err := provider.host.Network().DialPeer(ctx, peer.host.ID())
			if err != nil {
				t.Fatal(err)
			}
			waitGroup.Done()
		}
	}
	waitGroup.Wait()
}
