package libp2p

import (
	"context"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestUnicastChannelFuzzing(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 7000, func(
		_ *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		channel, err := provider1.UnicastChannelWith(identity2.id)
		if err != nil {
			t.Fatal(err)
		}

		provider2.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
			channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
				return &fuzzingMessage{}
			})
		})

		for i := 0; i < 100000; i++ {
			var message fuzzingMessage

			f := fuzz.New().NilChance(0.1).NumElements(0, 1024)
			f.Fuzz(&message)

			_ = channel.Send(&message)
		}
	})
}

func TestBroadcastChannelFuzzing(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 8000, func(
		_ *identity,
		_ *identity,
		provider net.Provider,
		_ net.Provider,
	) {
		channel, err := provider.BroadcastChannelFor("test-channel")
		if err != nil {
			t.Fatal(err)
		}

		err = channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &fuzzingMessage{}
		})
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 100000; i++ {
			var message fuzzingMessage

			f := fuzz.New().NilChance(0.1).NumElements(0, 1024)
			f.Fuzz(&message)

			_ = channel.Send(ctx, &message)
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
