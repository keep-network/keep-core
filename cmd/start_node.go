package cmd

import (
	"context"
	"fmt"

	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

const (
	sampleText             = "sample text"
	broadcastChannelName   = "test"
	resetBroadcastTimerSec = 5
)

// StartFlags for bootstrap and port
var StartFlags []cli.Flag

func init() {
	StartFlags = []cli.Flag{
		&cli.BoolFlag{
			Name: "bootstrap",
		},
		&cli.IntFlag{
			Name: "port",
		},
	}
}

// StartNode starts a node; if it's not a bootstrap node it will get the
// Node.URLs from the config file
func StartNode(c *cli.Context) error {
	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var port int
	if c.Int("port") > 0 {
		port = c.Int("port")
	} else {
		port = cfg.Node.Port
	}

	var (
		seed          int
		nodeName      string
		bootstrapURLs []string
	)
	if c.Bool("bootstrap") {
		nodeName = " bootstrap"
		seed = cfg.Bootstrap.Seed
	} else {
		bootstrapURLs = cfg.Bootstrap.URLs
	}

	ctx := context.Background()
	provider, err := libp2p.Connect(ctx, &libp2p.Config{
		Port:  port,
		Peers: bootstrapURLs,
		Seed:  seed,
	})
	if err != nil {
		return err
	}

	header(fmt.Sprintf(
		"starting%s node, connnecting to network at port %d",
		nodeName,
		port,
	))

	broadcastChannel, err := provider.ChannelFor(broadcastChannelName)
	if err != nil {
		return err
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &testMessage{} },
	); err != nil {
		return err
	}

	recvChan := make(chan net.Message)
	if err := broadcastChannel.Recv(func(msg net.Message) error {
		fmt.Printf("Got %s\n", msg.Payload())
		recvChan <- msg
		return nil
	}); err != nil {
		return err
	}

	go broadcastMessages(ctx, broadcastChannel, port)
	go receiveMessages(ctx, recvChan, port)

	select {}
}

func broadcastMessages(
	ctx context.Context,
	broadcastChannel net.BroadcastChannel,
	port int,
) {
	t := time.NewTimer(1) // first tick is immediate
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := broadcastChannel.Send(
				&testMessage{
					Payload: fmt.Sprintf(
						"%s on port %d",
						sampleText,
						port,
					),
				},
			); err != nil {
				return
			}
			t.Reset(resetBroadcastTimerSec * time.Second)
		case <-ctx.Done():
			return
		}

	}
}

func receiveMessages(
	ctx context.Context,
	recvChan <-chan net.Message,
	port int,
) {
	for {
		select {
		case msg := <-recvChan:
			testPayload := msg.Payload().(*testMessage)
			fmt.Printf("%d read message: %+v\n", port, testPayload)
		case <-ctx.Done():
			return
		}
	}
}
