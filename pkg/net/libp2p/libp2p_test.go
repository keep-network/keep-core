package libp2p

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/watchtower"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

func TestProviderReturnsType(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	privKey, _, err := key.GenerateStaticNetworkKey()
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

	privKey, _, err := key.GenerateStaticNetworkKey()
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

	privKey, _, err := key.GenerateStaticNetworkKey()
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

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func newTestIdentity() (*identity, error) {
	privKey, _, err := key.GenerateStaticNetworkKey()
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

// disconnect a peer that drops below min stake (unstake?)
// test that you are no longer connected
func TestDisconnectPeerUnderMinStake(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// initialize the first peer
	bootstrapPeerPrivKey, bootstrapPeerPubKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	minStake := big.NewInt(200)
	stakeMonitor := local.NewStakeMonitor(minStake)

	bootstrapPeerAddress := key.NetworkPubKeyToEthAddress(bootstrapPeerPubKey)
	bootstrapPeerStaker, err := stakeMonitor.StakerFor(bootstrapPeerAddress)
	if err != nil {
		t.Fatal(err)
	}
	err = stakeMonitor.StakeTokens(bootstrapPeerAddress)
	if err != nil {
		t.Fatal(err)
	}

	_, err = bootstrapPeerStaker.Stake()
	if err != nil {
		t.Fatal(err)
	}

	bootstrapIdentity, err := createIdentity(bootstrapPeerPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	// kick off the network
	bsHost, err := discoverAndListen(
		ctx,
		bootstrapIdentity,
		8080,
		stakeMonitor,
	)
	if err != nil {
		t.Fatal(err)
	}

	bootstrapRouter := dht.NewDHT(ctx, bsHost, dssync.MutexWrap(dstore.NewMapDatastore()))
	bootstrapHost := rhost.Wrap(bsHost, bootstrapRouter)

	// set watchtower for bootstrap peer
	_ = watchtower.NewGuard(ctx, stakeMonitor, bootstrapHost)

	// initialize second peer
	peerPrivKey, peerPubKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	peerIdentity, err := createIdentity(peerPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	peerAddress := key.NetworkPubKeyToEthAddress(peerPubKey)
	peerStaker, err := stakeMonitor.StakerFor(peerAddress)
	if err != nil {
		t.Fatal(err)
	}
	err = stakeMonitor.StakeTokens(peerAddress)
	if err != nil {
		t.Fatal(err)
	}

	_, err = peerStaker.Stake()
	if err != nil {
		t.Fatal(err)
	}

	peerHost, err := discoverAndListen(
		ctx,
		peerIdentity,
		8081,
		stakeMonitor,
	)
	if err != nil {
		t.Fatal(err)
	}

	peerRouter := dht.NewDHT(ctx, peerHost, dssync.MutexWrap(dstore.NewMapDatastore()))
	routedPeerHost := rhost.Wrap(peerHost, peerRouter)

	// set watchtower for our second peer
	_ = watchtower.NewGuard(ctx, stakeMonitor, peerHost)

	// connect the two peers
	pinfo := bootstrapHost.Peerstore().PeerInfo(bootstrapHost.ID())
	fmt.Printf("bootstrap peer info %+v\n", pinfo)
	if err := routedPeerHost.Connect(context.Background(), pinfo); err != nil {
		// TODO: this fails as we can't get addresses for the first
		// peer back from the peerstore.
		fmt.Printf("peer info of peer: %+v\n", peerHost.Peerstore().Addrs(bootstrapHost.ID()))
		t.Fatal(err)
	}
	// TODO: remove or decrease this time.
	time.Sleep(4 * time.Second)

	// make sure we have a valid connection
	if routedPeerHost.Network().Connectedness(bootstrapIdentity.id) != inet.Connected {
		t.Fatal("Failed to connect bootstrap peer to other peer")
	}

	// drop our second peer below the min stake
	// make sure the connection that the first peer has to the second peer has been untethered.
}
