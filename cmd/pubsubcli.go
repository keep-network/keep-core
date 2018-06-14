package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

// PubSubTestFlags for group size and threshold settings
var PubSubTestFlags []cli.Flag

func init() {
	PubSubTestFlags = []cli.Flag{
		&cli.IntFlag{},
	}
}

func PubSubTest(c *cli.Context) {
	go func() {
		provider, err := libp2p.Connect(context.Background(), &libp2p.Config{
			Port: 8080,
			Seed: 2,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = provider.ChannelFor("test")
		if err != nil {
			fmt.Println(err)
			return
		}

		select {}
	}()

	fmt.Println("sleep a bit")
	time.Sleep(3 * time.Second)

	for i := 1; i < 5; i++ {
		port := fmt.Sprintf("27%d%d", i, i)
		p, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println(err)
			return
		}

		var (
			ctx = context.Background()
		)

		provider, err := libp2p.Connect(ctx, &libp2p.Config{
			Port:  p,
			Peers: []string{"/ip4/127.0.0.1/tcp/8080/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		broadcastChannel, err := provider.ChannelFor("test")
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := broadcastChannel.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			fmt.Println(err)
			return
		}

		payload := fmt.Sprintf("some text from %s", port)
		if err := broadcastChannel.Send(
			&testMessage{Payload: payload},
		); err != nil {
			fmt.Println(err)
			return
		}

		recvChan := make(chan net.Message, 5)
		if err := broadcastChannel.Recv(func(msg net.Message) error {
			recvChan <- msg
			return nil
		}); err != nil {
			fmt.Println(err)
			return
		}

		go func(port string) {
			for {
				select {
				case msg := <-recvChan:
					testPayload := msg.Payload().(*testMessage)
					fmt.Printf("Message [%+v]\nRead by %s\n", testPayload, port)
				case <-ctx.Done():
					return
				}
			}
		}(port)
	}

	select {}
}

type protocolIdentifier struct {
	id string
}

type testMessage struct {
	Payload string
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
		fmt.Println("hit this error")
		return err
	}
	m.Payload = message.Payload

	return nil
}
