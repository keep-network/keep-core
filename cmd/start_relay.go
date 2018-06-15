package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
	"time"
	"strconv"
)

const (
	broadcastChannelName = "test"
	peerStartPortNo = 6000
)

// StartFlags for group size and threshold settings
var (
	StartFlags     []cli.Flag
)

func init() {
	StartFlags = []cli.Flag{
		&cli.BoolFlag{
			Name: "disable-relay",
		},
		&cli.BoolFlag{
			Name: "disable-provider",
		},
		&cli.IntFlag{
			Name: "node-count",
		},
	}
}

// Start performs a simulated distributed key generation and verifyies that the members can do a threshold signature
func StartRelay(c *cli.Context) error {

	header(fmt.Sprintf("starting DKG - GroupSize (%d), Threshold (%d)", defaultGroupSize, defaultThreshold))

	disableRelay := c.Bool("disable-relay")
	if disableRelay {
		return errors.New("no clients were selected for startup, so the program is exiting.")
	}

	disableProvider := c.Bool("disable-provider")
	if disableProvider {
		return errors.New("Keep provider has not yet been implemented.  Try back later!")
	}

	nodeCount := c.Int("node-count")
	if nodeCount == 0 {
		nodeCount = 5
	}

	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return errors.New(fmt.Sprintf("error reading config file", err))
	}

	go func() {
		port, err := portFromMa(cfg.Bootstrap.URL)
		if err != nil {
			fmt.Println(fmt.Sprintf("error parsing port from URL (%s): %v\n", cfg.Bootstrap.URL, err))
			return
		}
		provider, err := libp2p.Connect(context.Background(), &libp2p.Config{
			Port: port,
			Seed: cfg.Bootstrap.Seed,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = provider.ChannelFor(broadcastChannelName)
		if err != nil {
			fmt.Println(err)
			return
		}

		select {}
	}()

	fmt.Println("sleep a bit")
	time.Sleep(3 * time.Second)

	for i := 1; i < nodeCount; i++ {
		portStr := fmt.Sprintf("%d", peerStartPortNo + i)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}

		var (
			ctx = context.Background()
		)

		provider, err := libp2p.Connect(ctx, &libp2p.Config{
			Port:  port,
			Peers: []string{"/ip4/127.0.0.1/tcp/8080/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"},
		})
		if err != nil {
			return err
		}
		broadcastChannel, err := provider.ChannelFor("test")
		if err != nil {
			return err
		}

		if err := broadcastChannel.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			return err
		}

		payload := fmt.Sprintf("some text from %d", port)
		if err := broadcastChannel.Send(
			&testMessage{Payload: payload},
		); err != nil {
			return err
		}

		recvChan := make(chan net.Message, 5)
		if err := broadcastChannel.Recv(func(msg net.Message) error {
			fmt.Printf("Got %v\n", msg)
			recvChan <- msg
			return nil
		}); err != nil {
			return err
		}

		go func(portStr string) {
			for {
				select {
				case msg := <-recvChan:
					testPayload := msg.Payload().(*testMessage)
					fmt.Printf("Message [%+v]\nRead by %s\n", testPayload, portStr)
				case <-ctx.Done():
					return
				}
			}
		}(portStr)
	}

	select {}
}