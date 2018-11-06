package cmd

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

// PingCommand contains the definition of the ping command-line subcommand.
var PingCommand cli.Command

const (
	ping = "PING"
	pong = "PONG"
)

const pingDescription = `The ping command conducts a simple peer-to-peer test
   between a bootstrap node and another peer: can known peers communicate over
   a peer-to-peer network. Both peers send a "PING" and expect to receive a
   corresponding "PONG". Notably, this does not exercise peer discovery.`

func init() {
	PingCommand =
		cli.Command{
			Name:        "ping",
			Usage:       `bidirectional send between two peers to test the network`,
			ArgsUsage:   "[multiaddr]",
			Description: pingDescription,
			Action:      pingRequest,
		}
}

func isBootstrapNode(args cli.Args) (bool, []string) {
	var bootstrapPeers []string

	// Not a bootstrap node
	if len(args) > 0 {
		bootstrapPeers = append(bootstrapPeers, args.Get(0))
	}

	return len(bootstrapPeers) == 0, bootstrapPeers
}

// pingRequest tests the functionality and availability of Keep's libp2p
// network layer.
func pingRequest(c *cli.Context) error {
	isBootstrapNode, bootstrapPeers := isBootstrapNode(c.Args())
	var (
		libp2pConfig = libp2p.Config{Peers: bootstrapPeers}
		ctx          = context.Background()
	)

	privKey, pubKey, err := key.GenerateStaticNetworkKey(rand.Reader)
	if err != nil {
		return err
	}

	stakeMonitoring := local.NewStakeMonitoring()
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(pubKey))

	netProvider, err := libp2p.Connect(
		ctx,
		libp2pConfig,
		privKey,
		stakeMonitoring,
	)
	if err != nil {
		return err
	}

	if isBootstrapNode {
		var bootstrapAddr string
		for _, addr := range netProvider.AddrStrings() {
			if strings.Contains(addr, "ip4") && !strings.Contains(addr, "127.0.0.1") {
				bootstrapAddr = addr
				break
			}
		}

		fmt.Printf("You can ping this node using:\n"+
			"    %s ping %s\n\n",
			c.App.Name,
			bootstrapAddr,
		)
	}

	// When we call ChannelFor, we create a coordination point for peers
	broadcastChannel, err := netProvider.ChannelFor(ping)
	if err != nil {
		return err
	}

	// PingMessage and PongMessage conform to the net.Message interface
	// (Type, Unmarshal, Marshal); ensure our network knows how to serialize
	// them when sending over the wire
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

	err = broadcastChannel.Recv(
		net.HandleMessageFunc{
			Type: pong,
			Handler: func(msg net.Message) error {
				// Do some message routing
				if msg.Type() == pong {
					pongChan <- msg
				}
				return nil
			},
		},
	)
	if err != nil {
		return err
	}
	defer broadcastChannel.UnregisterRecv(pong)

	err = broadcastChannel.Recv(
		net.HandleMessageFunc{
			Type: ping,
			Handler: func(msg net.Message) error {
				// Do some message routing
				if msg.Type() == ping {
					pingChan <- msg
				}
				return nil
			},
		},
	)
	if err != nil {
		return err
	}
	defer broadcastChannel.UnregisterRecv(ping)

	// Give ourselves a moment to form a mesh with the other peer
	for {
		time.Sleep(3 * time.Second)
		peers := netProvider.Peers()
		if len(peers) < 1 {
			fmt.Println("waiting for peer...")
			continue
		}
		break
	}

	if err := broadcastChannel.Send(
		&PingMessage{Sender: netProvider.ID().String(), Payload: ping},
	); err != nil {
		return err
	}
	fmt.Println("Sent PING")

	for {
		select {
		case msg := <-pingChan:
			// don't read our own ping
			if msg.TransportSenderID().String() == netProvider.ID().String() {
				continue
			}
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

			// Help with synchronization between slow clients.
			// Occasionally the client exits before successfully
			// writing the ping to the wire.
			time.Sleep(1 * time.Second)
		case msg := <-pongChan:
			// don't read our own pong
			if msg.TransportSenderID().String() == netProvider.ID().String() {
				continue
			}
			// if you read a pong message, go ahead and ack and close out
			pongPayload, ok := msg.Payload().(*PongMessage)
			if !ok {
				return fmt.Errorf(
					"expected: payload type PongMessage\nactual:   payload type [%v]",
					pongPayload,
				)
			}

			fmt.Printf("Received PONG from %s", msg.TransportSenderID().String())
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
