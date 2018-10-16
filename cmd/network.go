package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

// PingCommand contains the definition of the ping command-line subcommand.
var PingCommand cli.Command

const (
	ping = "PING"
	pong = "PONG"
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
		}
}

// pingRequest tests our network layer code
//
// requests a new entry from the threshold relay and prints the
// request id. By default, it also waits until the associated relay entry is
// generated and prints out the entry.
func pingRequest(c *cli.Context) error {
	config, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: [%v]", err)
	}

	ctx := context.Background()
	netProvider, err := libp2p.Connect(ctx, config.LibP2P)
	if err != nil {
		return err
	}

	isBootstrapNode := config.LibP2P.Seed != 0
	// TODO: make this custom output
	nodeHeader(isBootstrapNode, netProvider.AddrStrings(), config.LibP2P.Port)

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
		fmt.Printf("Got peer %s\n", peers[0])
		break
	}

	if isBootstrapNode {
		time.Sleep(1 * time.Second)
	}

	if isBootstrapNode {
		err = broadcastChannel.Send(
			&PingMessage{
				Sender:  netProvider.ID().String(),
				Payload: ping})
		if err != nil {
			return err
		}
	}

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

			// // Ensure we don't send the response to the peer that
			// // sent the originating ping
			// if netProvider.ID().String() == pingPayload.Sender {
			// 	continue
			// }
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
