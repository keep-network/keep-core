package cmd

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

const (
	broadcastChannelName = "test"
	sampleText           = "sample text"
)

// StartFlags for bootstrap, port and disable-provider
var	StartFlags []cli.Flag

func init() {
	StartFlags = []cli.Flag{
		&cli.BoolFlag{
			Name: "bootstrap",
		},
		&cli.IntFlag{
			Name: "port",
		},
		&cli.StringFlag{
			Name: "preferred-ip-address",
		},
		&cli.BoolFlag{
			Name: "disable-provider",
		},
	}
}

// StartNode starts a node; if it's not a bootstrap node it will get the Node.URLs from the config file
func StartNode(c *cli.Context) error {
	disableProvider := c.Bool("disable-provider")
	if disableProvider {
		return errors.New("Keep provider has not yet been implemented.  Try back later!")
	}

	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	preferredIPAddress := c.String("preferred-ip-address")
	if len(preferredIPAddress) == 0 {
		preferredIPAddress = cfg.Node.MyPreferredOutboundIP
	}

	myIPAddress := GetMyIPv4Address(preferredIPAddress)
	var port int
	if c.Int("port") > 0 {
		port = c.Int("port")
	} else {
		port = cfg.Node.Port
	}

	var (
		seed int
		nodeName string
		bootstrapURLs []string
	)
	if c.Bool("bootstrap") {
		nodeName = " bootstrap"
		seed = cfg.Bootstrap.Seed
	} else {
		bootstrapURLs = cfg.Bootstrap.URLs
	}

	header(fmt.Sprintf("starting%s node, connnecting to network and listening at %s port %d", nodeName, myIPAddress, port))

	ctx := context.Background()
	provider, err := libp2p.Connect(ctx, &libp2p.Config{
		Port:  port,
		Peers: bootstrapURLs,
		Seed:  seed,
	})
	if err != nil {
		return err
	}
	broadcastChannel, err := provider.ChannelFor(broadcastChannelName)
	if err != nil {
		return err
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &testMessage{} },
	); err != nil {
		return err
	}

	go func() {
		t := time.NewTimer(1) // first tick is immediate
		defer t.Stop()
		for {
			select {
			case <-t.C:
				payload := fmt.Sprintf("%s from %s on port %d", sampleText, myIPAddress, port)
				if err := broadcastChannel.Send(
					&testMessage{Payload: payload},
				); err != nil {
					return
				}
				t.Reset(5 * time.Second)
			case <-ctx.Done():
				return
			}

		}
	}()

	recvChan := make(chan net.Message)

	if err := broadcastChannel.Recv(func(msg net.Message) error {
		fmt.Printf("Got %s\n", msg.Payload())
		recvChan <- msg
		return nil
	}); err != nil {
		return err
	}

	go func(port int) {
		for {
			select {
			case msg := <-recvChan:
				testPayload := msg.Payload().(*testMessage)
				fmt.Printf("%s:%d read message: %+v\n", myIPAddress, port, testPayload)
			case <-ctx.Done():
				return
			}
		}
	}(port)

	select {}
}
