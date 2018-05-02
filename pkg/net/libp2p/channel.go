package libp2p

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core-dkg-branch/go/beacon/broadcast"
	"github.com/keep-network/keep-core/pkg/net"
	floodsub "github.com/libp2p/go-floodsub"
	peer "github.com/libp2p/go-libp2p-peer"
)

// TODO: if it's absolutely necessary to have a pk/sk, make this a struct
type peerIdentifier peer.ID

func (n peerIdentifier) ProviderName() string {
	return "libp2p"
}

type channel struct {
	broadcast.Channel
	name string
	sub  *floodsub.Subscription

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler

	identifiersMutex            sync.Mutex
	transportToProtoIdentifiers map[net.TransportIdentifier]net.ProtocolIdentifier
	protoToTransportIdentifiers map[net.ProtocolIdentifier]net.TransportIdentifier
}

func (c *channel) Name() string {
	// TODO: lock this if being updated by an external party
	return c.name
}

func (c *channel) Send(message net.TaggedMarshaler) error {
	return nil
}

func (c *channel) SendTo(
	recipientIdentifier interface{},
	message net.TaggedMarshaler,
) error {
	// TODO: do a check between the registered protocol identifier and the network ident
	return nil
}

func (c *channel) Recv(h net.HandleMessageFunc) error {
	return nil
}

func (c *channel) RegisterIdentifier(
	transportIdentifier net.TransportIdentifier,
	protocolIdentifier net.ProtocolIdentifier,
) error {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	if _, ok := transportIdentifier.(peerIdentifier); !ok {
		return fmt.Errorf(
			"incorrect type for transportIdentifier: [%v]",
			transportIdentifier,
		)
	}

	if _, exists := c.transportToProtoIdentifiers[transportIdentifier]; exists {
		return fmt.Errorf(
			"already have a protocol identifier in channel [%v] associated with [%v]",
			c, transportIdentifier,
		)
	}
	if _, exists := c.protoToTransportIdentifiers[protocolIdentifier]; exists {
		return fmt.Errorf(
			"already have a transport identifier in channel [%v] associated with [%v]",
			c, protocolIdentifier,
		)
	}

	c.transportToProtoIdentifiers[transportIdentifier] = protocolIdentifier
	c.protoToTransportIdentifiers[protocolIdentifier] = transportIdentifier

	return nil
}

func (c *channel) RegisterUnmarshaler(unmarshaler func() net.TaggedUnmarshaler) error {
	tpe := unmarshaler().Type()

	c.unmarshalersMutex.Lock()
	defer c.unmarshalersMutex.Unlock()

	if _, exists := c.unmarshalersByType[tpe]; exists {
		return fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	}

	c.unmarshalersByType[tpe] = unmarshaler
	return nil
}
