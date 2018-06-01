package libp2p

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
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
	ctx, cancel := newTestContext()
	defer cancel()
	// ctx := context.Background()

	config := generateDeterministicNetworkConfig(t)

	tests := map[string]struct {
		name                    string
		testIdentity            *identity
		protocolIdentifier      *protocolIdentifier
		expectedChannelForError func(string) error
	}{
		"Send succeeds": {
			name: "testchannel",
			expectedChannelForError: func(name string) error {
				return nil
			},
			testIdentity:       config.identity,
			protocolIdentifier: &protocolIdentifier{id: "testProtocolIdentifier"},
		},
	}

	provider, err := Connect(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			broadcastChannel, err := provider.ChannelFor(test.name)
			if !reflect.DeepEqual(test.expectedChannelForError(test.name), err) {
				t.Fatalf("expected test to fail with [%v] instead failed with [%v]",
					test.expectedChannelForError(test.name), err,
				)
			}

			if err := broadcastChannel.RegisterUnmarshaler(
				func() net.TaggedUnmarshaler { return &TestMessage{} },
			); err != nil {
				t.Fatal(err)
			}

			if err := broadcastChannel.RegisterIdentifier(
				test.testIdentity.id,
				test.protocolIdentifier,
			); err != nil {
				t.Fatal(err)
			}

			if err := broadcastChannel.Send(
				&TestMessage{ID: test.testIdentity, Payload: "some text"},
			); err != nil {
				t.Fatal(err)
			}

			recvChan := make(chan net.Message, 1)
			if err := broadcastChannel.Recv(func(msg net.Message) error {
				// slap something onto a channel and move on?
				recvChan <- msg

				// if msg.Payload() != test.ExpectedPayload {
				// 	t.Fatal("expected message payload %s, got payload %s", msg.Payload(), test.ExpectedPayload)
				// }
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			select {
			case msg := <-recvChan:
				fmt.Printf("Message: %+v\n", msg)
				return
			case <-ctx.Done():
				return
			}
		})
	}
}

type protocolIdentifier struct {
	id string
}

type TestMessage struct {
	ID      *identity
	Payload string
}

// Type returns a string describing a TestMessage's type.
func (m *TestMessage) Type() string {
	return "test/unmarshaler"
}

// Marshal converts this TestMessage to a byte array suitable for network
// communication.
func (m *TestMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal converts a byte array produced by Marshal to a JoinMessage.
func (m *TestMessage) Unmarshal(bytes []byte) error {
	var message TestMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println("hit this error")
		return err
	}
	m.ID = message.ID
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

func generateDeterministicNetworkConfig(t *testing.T) *Config {
	p := testutils.RandPeerNetParamsOrFatal(t)
	pi := &identity{id: networkIdentity{p.ID}, privKey: p.PrivKey, pubKey: p.PubKey}
	return &Config{port: 8080, listenAddrs: []ma.Multiaddr{p.Addr}, identity: pi}
}

func testProvider(ctx context.Context, t *testing.T) (*provider, error) {
	testConfig := generateDeterministicNetworkConfig(t)

	host, identity, err := discoverAndListen(ctx, testConfig)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, identity, host)
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
