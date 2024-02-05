package libp2p

import (
	"context"
	"strconv"
	"testing"

	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/libp2p/go-libp2p/core/peer"

	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/multiformats/go-multiaddr"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
)

func TestFuzzBroadcastChannelReceiveSide(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 7000, func(
		_ *identity,
		_ *identity,
		provider net.Provider,
		_ net.Provider,
	) {
		testChannel, err := provider.BroadcastChannelFor("test-channel")
		if err != nil {
			t.Fatal(err)
		}

		broadcastTestChannel, ok := testChannel.(*channel)
		if !ok {
			t.Fatal("could not cast to broadcast channel")
		}

		for i := 0; i < 100; i++ {
			var (
				pbMessage     pubsubpb.Message
				receivedFrom  peer.ID
				validatorData string
			)

			f := fuzz.New().NilChance(0.01).NumElements(0, 1024)
			f.Fuzz(&pbMessage)
			f.Fuzz(&receivedFrom)
			f.Fuzz(&validatorData)

			message := &pubsub.Message{
				Message:       &pbMessage,
				ReceivedFrom:  receivedFrom,
				ValidatorData: validatorData,
			}
			_ = broadcastTestChannel.processPubsubMessage(message)
		}
	})
}

func TestFuzzBroadcastChannelRoundtrip(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 8000, func(
		_ *identity,
		_ *identity,
		provider net.Provider,
		_ net.Provider,
	) {
		testChannel, err := provider.BroadcastChannelFor("test-channel")
		if err != nil {
			t.Fatal(err)
		}

		testChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
			return &fuzzingMessage{}
		})

		for i := 0; i < 100; i++ {
			var message fuzzingMessage

			f := fuzz.New().NilChance(0.01).NumElements(0, 1024)
			f.Fuzz(&message)

			_ = testChannel.Send(ctx, &message)
		}
	})
}

type fuzzingMessage struct {
	Bytes []byte
}

func (m *fuzzingMessage) Type() string {
	return "test/fuzzing"
}

func (m *fuzzingMessage) Marshal() ([]byte, error) {
	return m.Bytes, nil
}

func (m *fuzzingMessage) Unmarshal(bytes []byte) error {
	m.Bytes = bytes
	return nil
}

func withNetwork(
	ctx context.Context,
	t *testing.T,
	startingPort int,
	testFn func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	),
) {
	operatorPrivateKey1, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey1, _, err := operatorPrivateKeyToNetworkKeyPair(operatorPrivateKey1)
	if err != nil {
		t.Fatal(err)
	}

	identity1, err := createIdentity(networkPrivateKey1)
	if err != nil {
		t.Fatal(err)
	}

	provider1Port := startingPort

	multiaddr1, err := multiaddr.NewMultiaddr(
		"/ip4/127.0.0.1/tcp/" + strconv.Itoa(provider1Port),
	)
	if err != nil {
		t.Fatal(err)
	}

	operatorPrivateKey2, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey2, _, err := operatorPrivateKeyToNetworkKeyPair(operatorPrivateKey2)
	if err != nil {
		t.Fatal(err)
	}

	identity2, err := createIdentity(networkPrivateKey2)
	if err != nil {
		t.Fatal(err)
	}

	provider2Port := startingPort + 1

	multiaddr2, err := multiaddr.NewMultiaddr(
		"/ip4/127.0.0.1/tcp/" + strconv.Itoa(provider2Port),
	)
	if err != nil {
		t.Fatal(err)
	}

	provider1, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr2,
					identity2.id,
				),
			},
			Port: provider1Port,
		},
		operatorPrivateKey1,
		firewall.Disabled,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	provider2, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr1,
					identity1.id,
				),
			},
			Port: provider2Port,
		},
		operatorPrivateKey2,
		firewall.Disabled,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	testFn(
		identity1,
		identity2,
		provider1,
		provider2,
	)
}
