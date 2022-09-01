package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/keep-network/keep-core/pkg/operator"
)

// PingCommand contains the definition of the ping command-line subcommand.
var PingCommand = &cobra.Command{
	Use:                   "ping [multiaddr]...",
	Short:                 `bidirectional send between two peers to test the network`,
	Long:                  pingDescription,
	DisableFlagsInUseLine: true,
	RunE:                  pingRequest,
}

const (
	ping         = "PING"
	pong         = "PONG"
	messageCount = 64 * 64
)

const pingDescription = `The ping command conducts a simple peer-to-peer test
   between a bootstrap node and another peer: can known peers communicate over
   a peer-to-peer network. Both peers send a "PING" and expect to receive a
   corresponding "PONG". Notably, this does not exercise peer discovery.`

func isBootstrapNode(args []string) (bool, []string) {
	var bootstrapPeers []string

	// Not a bootstrap node
	if len(args) > 0 {
		bootstrapPeers = append(bootstrapPeers, args[0])
	}

	return len(bootstrapPeers) == 0, bootstrapPeers
}

// pingRequest tests the functionality and availability of Keep's libp2p
// network layer.
func pingRequest(cmd *cobra.Command, args []string) error {
	isBootstrapNode, bootstrapPeers := isBootstrapNode(args)
	var (
		libp2pConfig = libp2p.Config{Peers: bootstrapPeers}
		ctx          = context.Background()
		privKey      *operator.PrivateKey
	)

	bootstrapPeerPrivKey, _ := getBootstrapPeerOperatorKey()
	standardPeerPrivKey, _ := getStandardPeerOperatorKey()

	if isBootstrapNode {
		privKey = bootstrapPeerPrivKey
	} else {
		privKey = standardPeerPrivKey
	}

	netProvider, err := libp2p.Connect(
		ctx,
		libp2pConfig,
		privKey,
		firewall.Disabled,
		retransmission.NewTimeTicker(ctx, 50*time.Millisecond),
	)
	if err != nil {
		return err
	}

	if isBootstrapNode {
		var bootstrapAddr string
		for _, addr := range netProvider.ConnectionManager().AddrStrings() {
			if strings.Contains(addr, "ip4") && !strings.Contains(addr, "127.0.0.1") {
				bootstrapAddr = addr
				break
			}
		}

		fmt.Printf("You can ping this node using:\n"+
			"    %s %s\n\n",
			cmd.CommandPath(),
			bootstrapAddr,
		)
	}

	// When we call ChannelFor, we create a coordination point for peers
	broadcastChannel, err := netProvider.BroadcastChannelFor(ping)
	if err != nil {
		return err
	}

	// PingMessage and PongMessage conform to the net.Message interface
	// (Type, Unmarshal, Marshal); ensure our network knows how to serialize
	// them when sending over the wire
	broadcastChannel.SetUnmarshaler(
		func() net.TaggedUnmarshaler { return &PingMessage{} },
	)
	broadcastChannel.SetUnmarshaler(
		func() net.TaggedUnmarshaler { return &PongMessage{} },
	)

	var (
		pingChan = make(chan net.Message)
		pongChan = make(chan net.Message)
	)

	broadcastChannel.Recv(ctx, func(msg net.Message) {
		// Do some message routing
		if msg.Type() == pong {
			pongChan <- msg
		}
	})

	broadcastChannel.Recv(ctx, func(msg net.Message) {
		// Do some message routing
		if msg.Type() == ping {
			pingChan <- msg
		}
	})

	// Give ourselves a moment to form a mesh with the other peer
	for {
		time.Sleep(3 * time.Second)
		peers := netProvider.ConnectionManager().ConnectedPeers()
		if len(peers) < 1 {
			fmt.Println("waiting for peer...")
			continue
		}
		break
	}

	start := make(chan struct{})
	receivedMessages := make(map[string]bool)

	for i := 1; i <= messageCount; i++ {
		message := &PingMessage{
			Sender:  netProvider.ID().String(),
			Payload: ping + " number " + strconv.Itoa(i),
		}

		go func(msg *PingMessage) {
			<-start
			err := broadcastChannel.Send(ctx, message)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Error while sending PING with payload [%v]: [%v]\n",
					message.Payload,
					err,
				)
			} else {
				fmt.Printf("Sent PING with payload [%v]\n", message.Payload)
			}
		}(message)
	}

	close(start)

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

			message := &PongMessage{
				Sender:  netProvider.ID().String(),
				Payload: pong + " corresponding to " + pingPayload.Payload,
			}
			err := broadcastChannel.Send(ctx, message)
			if err != nil {
				return err
			}
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

			fmt.Printf(
				"Received PONG from [%s] with payload [%v]\n",
				msg.TransportSenderID().String(),
				pongPayload.Payload,
			)

			receivedMessages[pongPayload.Payload] = true

			if len(receivedMessages) == messageCount {
				fmt.Println("All expected messages received")
			}
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

// PingMessage is a network message sent between bootstrap peer and
// non-bootstrap peer in order to test the connection.
type PingMessage struct {
	Sender  string
	Payload string
}

// Type returns a string type of the `PingMessage` so that it conforms to
// `net.Message` interface.
func (pm *PingMessage) Type() string {
	return ping
}

// Marshal converts this PingMessage to a byte array suitable for network
// communication.
func (pm *PingMessage) Marshal() ([]byte, error) {
	return json.Marshal(pm)
}

// Unmarshal converts a byte array produced by Marshal to a PingMessage.
func (pm *PingMessage) Unmarshal(bytes []byte) error {
	var message PingMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}
	pm.Sender = message.Sender
	pm.Payload = message.Payload

	return nil
}

// PongMessage is a network message sent between bootstrap peer and
// non-bootstrap peer in order to test the connection.
type PongMessage struct {
	Sender  string
	Payload string
}

// Type returns a string type of the `PongMessage` so that it conforms to
// `net.Message` interface.
func (pm *PongMessage) Type() string {
	return pong
}

// Marshal converts this PongMessage to a byte array suitable for network
// communication.
func (pm *PongMessage) Marshal() ([]byte, error) {
	return json.Marshal(pm)
}

// Unmarshal converts a byte array produced by Marshal to a PongMessage.
func (pm *PongMessage) Unmarshal(bytes []byte) error {
	var message PongMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}
	pm.Sender = message.Sender
	pm.Payload = message.Payload

	return nil
}

// getBootstrapPeerOperatorKey returns hardcoded public and private operator key
// of the bootstrap peer.
func getBootstrapPeerOperatorKey() (
	*operator.PrivateKey,
	*operator.PublicKey,
) {
	return getPeerOperatorKey(big.NewInt(128838122312))
}

// getStandardPeerOperatorKey returns hardcoded public and private operator key
// of the standard peer.
func getStandardPeerOperatorKey() (
	*operator.PrivateKey,
	*operator.PublicKey,
) {
	return getPeerOperatorKey(big.NewInt(6743262236222))
}

func getPeerOperatorKey(privateEcdsaKey *big.Int) (
	*operator.PrivateKey,
	*operator.PublicKey,
) {
	x, y := secp256k1.S256().ScalarBaseMult(privateEcdsaKey.Bytes())

	operatorPublicKey := operator.PublicKey{
		Curve: operator.Secp256k1,
		X:     x,
		Y:     y,
	}

	operatorPrivateKey := operator.PrivateKey{
		PublicKey: operatorPublicKey,
		D:         privateEcdsaKey,
	}

	return &operatorPrivateKey, &operatorPublicKey
}
