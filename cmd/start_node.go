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

type recvParams struct {
	port     int
	ipaddr   string
	recvChan chan net.Message
}

type broadcastParams struct {
	port      int
	ipaddr    string
	bcastChan net.BroadcastChannel
}

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
		bootstrapURLs []string
	)
	if c.Bool("bootstrap") {
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

	myIPv4Address := GetIPv4Address(provider.Addrs())

	nodeHeader(c.Bool("bootstrap"), myIPv4Address, port)

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
		fmt.Printf("got %s\n", msg.Payload())
		recvChan <- msg
		return nil
	}); err != nil {
		return err
	}

	go broadcastMessages(
		ctx,
		broadcastParams{
			port:      port,
			ipaddr:    myIPv4Address,
			bcastChan: broadcastChannel,
		},
	)
	go receiveMessage(
		ctx,
		recvParams{
			port:     port,
			ipaddr:   myIPv4Address,
			recvChan: recvChan,
		},
	)

	select {}
}

func broadcastMessages(ctx context.Context, params broadcastParams) {
	t := time.NewTimer(1) // first tick is immediate
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := params.bcastChan.Send(
				&testMessage{
					Payload: fmt.Sprintf(
						"%s from %s on port %d",
						sampleText,
						params.ipaddr,
						params.port,
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

func receiveMessage(ctx context.Context, params recvParams) {
	for {
		select {
		case msg := <-params.recvChan:
			testPayload := msg.Payload().(*testMessage)
			fmt.Printf(
				"%s:%d read message: %+v\n",
				params.ipaddr,
				params.port,
				testPayload,
			)
		case <-ctx.Done():
			return
		}
	}
}
