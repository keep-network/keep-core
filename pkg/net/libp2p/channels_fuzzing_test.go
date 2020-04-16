package libp2p

import (
	"context"
	"testing"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/libp2p/go-libp2p-core/peer"

	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/net"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
)

func TestFuzzUnicastChannelReceiveSide(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 5000, func(
		_ *identity,
		identity2 *identity,
		provider1 net.Provider,
		_ net.Provider,
	) {
		testChannel, err := provider1.UnicastChannelWith(identity2.id)
		if err != nil {
			t.Fatal(err)
		}

		unicastTestChannel, ok := testChannel.(*unicastChannel)
		if !ok {
			t.Fatal("could not cast to unicast channel")
		}

		for i := 0; i < 100; i++ {
			var message pb.UnicastNetworkMessage

			f := fuzz.New().NilChance(0.01).NumElements(0, 1024)
			f.Fuzz(&message)

			_ = unicastTestChannel.processMessage(&message)
		}
	})
}

func TestFuzzUnicastChannelRoundtrip(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 6000, func(
		_ *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		testChannel, err := provider1.UnicastChannelWith(identity2.id)
		if err != nil {
			t.Fatal(err)
		}

		provider2.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
			channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
				return &fuzzingMessage{}
			})
		})

		for i := 0; i < 100; i++ {
			var message fuzzingMessage

			f := fuzz.New().NilChance(0.01).NumElements(0, 1024)
			f.Fuzz(&message)

			_ = testChannel.Send(&message)
		}
	})
}

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

		err = testChannel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &fuzzingMessage{}
		})
		if err != nil {
			t.Fatal(err)
		}

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
