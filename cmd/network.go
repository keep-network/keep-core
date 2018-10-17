package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

// PingCommand contains the definition of the ping command-line subcommand.
var PingCommand cli.Command

const (
	ping = "PING"
	pong = "PONG"

	bootstrapPeerFlag = "bootstrap-peer"
	bootstrapShort    = "b"
)

const pingDescription = `The ping command allows a peer to construct an adhoc
   peer-to-peer network with other known peers. Known peers can then ping this
   client (or vis versa) to ensure that they are able to properly build up and
   tear down all of the connections. At least one peer must take the
   responsibility of being the bootstrap node.`

func init() {
	PingCommand =
		cli.Command{
			Name:        "ping",
			Usage:       ``,
			Description: pingDescription,
			Action:      pingRequest,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: bootstrapPeerFlag + "," + bootstrapShort,
				},
			},
		}
}

// pingRequest tests our network layer code
//
// requests a new entry from the threshold relay and prints the
// request id. By default, it also waits until the associated relay entry is
// generated and prints out the entry.
func pingRequest(c *cli.Context) error {
	var bootstrapPeers []string

	if c.String(bootstrapPeerFlag) != "" {
		bootstrapPeers = append(bootstrapPeers, c.String(bootstrapPeerFlag))
	}

	libp2pConfig := libp2p.Config{Peers: bootstrapPeers}

	ctx := context.Background()
	netProvider, err := libp2p.Connect(ctx, libp2pConfig)
	if err != nil {
		return err
	}

	isBootstrapNode := len(libp2pConfig.Peers) == 0

	if isBootstrapNode {
		var bootstrapAddr string
		for _, addr := range netProvider.AddrStrings() {
			if strings.Contains(addr, "ip4") && !strings.Contains(addr, "127.0.0.1") {
				bootstrapAddr = addr
				break
			}
		}

		fmt.Printf("Enable other peer with:\n"+
			"   > ./keep-core ping -bootstrap-peer %s\n"+
			"modifications to the above may be necessary\n",
			bootstrapAddr,
		)
	}

	// When we call ChannelFor, we create a coordination point for peers
	broadcastChannel, err := netProvider.ChannelFor(ping)
	if err != nil {
		return err
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &PingMessage{} },
	); err != nil {
		return err
	}
	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &PongMessage{} },
	); err != nil {
		return err
	}

	if err := broadcastChannel.RegisterIdentifier(
		netProvider.ID(),
		netProvider.ID().String(),
	); err != nil {
		return err
	}

	var (
		pingChan = make(chan net.Message)
		pongChan = make(chan net.Message)
	)

	if err := broadcastChannel.Recv(net.HandleMessageFunc{
		Type: pong,
		Handler: func(msg net.Message) error {
			// Do some message routing
			if msg.Type() == pong {
				pongChan <- msg
			}
			return nil
		},
	}); err != nil {
		return err
	}
	defer broadcastChannel.UnregisterRecv(pong)

	if err := broadcastChannel.Recv(net.HandleMessageFunc{
		Type: ping,
		Handler: func(msg net.Message) error {
			// Do some message routing
			if msg.Type() == ping {
				pingChan <- msg
			}
			return nil
		},
	}); err != nil {
		return err
	}
	defer broadcastChannel.UnregisterRecv(ping)

	// Give ourselves a moment to form a mesh with potential peers
	for {
		time.Sleep(3 * time.Second)
		peers := netProvider.Peers()
		if len(peers) < 1 {
			fmt.Println("waiting for peer...\n")
			continue
		}
		break
	}

	err = broadcastChannel.Send(
		&PingMessage{
			Sender:  netProvider.ID().String(),
			Payload: ping})
	if err != nil {
		return err
	}
	fmt.Println("Sent PING")

	for {
		select {
		case msg := <-pingChan:
			pingPayload, ok := msg.Payload().(*PingMessage)
			if !ok {
				return fmt.Errorf(
					"expected: payload type PingMessage\nactual:   payload type [%v]",
					pingPayload,
				)
			}

			err := broadcastChannel.Send(
				&PongMessage{
					Sender:  netProvider.ID().String(),
					Payload: pong})
			if err != nil {
				return err
			}
		case msg := <-pongChan:
			// if you read a pong message, go ahead and ack and close out
			pongPayload, ok := msg.Payload().(*PongMessage)
			if !ok {
				return fmt.Errorf(
					"expected: payload type PongMessage\nactual:   payload type [%v]",
					pongPayload,
				)
			}

			fmt.Println("Received PONG")
			return nil
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Request errored out: [%v].\n",
					err,
				)
			} else {
				fmt.Fprintf(os.Stderr, "Request errored for unknown reason.\n")
			}

			os.Exit(1)
		}
	}
}

type PingMessage struct {
	Sender  string
	Payload string
}

func (p *PingMessage) Type() string {
	return ping
}

// Marshal converts this PingMessage to a byte array suitable for network
// communication.
func (p *PingMessage) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// Unmarshal converts a byte array produced by Marshal to a PingMessage.
func (m *PingMessage) Unmarshal(bytes []byte) error {
	var message PingMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}
	m.Sender = message.Sender
	m.Payload = message.Payload

	return nil
}

type PongMessage struct {
	Sender  string
	Payload string
}

func (p *PongMessage) Type() string {
	return pong
}

// Marshal converts this PongMessage to a byte array suitable for network
// communication.
func (p *PongMessage) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// Unmarshal converts a byte array produced by Marshal to a PongMessage.
func (m *PongMessage) Unmarshal(bytes []byte) error {
	var message PongMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}
	m.Sender = message.Sender
	m.Payload = message.Payload

	return nil
}
