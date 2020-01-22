package libp2p

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net/retransmission"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/multiformats/go-multiaddr"
)

func TestProviderCreatesUnicastChannel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	withNetwork(t, ctx, func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		_, err := provider1.ChannelWith(identity2.id.String())
		if err != nil {
			t.Fatal(err)
		}

		_, err = provider2.ChannelWith(identity1.id.String())
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSendUnicastMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	withNetwork(t, ctx, func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		channel1, err := provider1.ChannelWith(identity2.id.String())
		if err != nil {
			t.Fatal(err)
		}

		channel2, err := provider2.ChannelWith(identity1.id.String())
		if err != nil {
			t.Fatal(err)
		}

		if err := channel1.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			t.Fatal(err)
		}

		if err := channel2.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			t.Fatal(err)
		}

		go func() {
			if err := channel1.Send(
				ctx,
				&testMessage{Sender: identity1, Payload: "yolo"},
			); err != nil {
				t.Fatal(err)
			}
		}()

		time.Sleep(10 * time.Second)
	})
}

func withNetwork(
	t *testing.T,
	ctx context.Context,
	testFn func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	),
) {
	privKey1, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	identity1, err := createIdentity(privKey1)
	if err != nil {
		t.Fatal(err)
	}

	multiaddr1, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/8081")
	if err != nil {
		t.Fatal(err)
	}

	privKey2, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	identity2, err := createIdentity(privKey2)
	if err != nil {
		t.Fatal(err)
	}

	multiaddr2, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/8082")
	if err != nil {
		t.Fatal(err)
	}

	stakeMonitor := local.NewStakeMonitor(big.NewInt(0))

	provider1, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr2,
					identity2.id,
				)},
			Port: 8081,
		},
		privKey1,
		stakeMonitor,
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
				)},
			Port: 8082,
		},
		privKey2,
		stakeMonitor,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	testFn(identity1, identity2, provider1, provider2)
}
