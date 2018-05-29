package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	floodsub "github.com/libp2p/go-floodsub"
)

type channel struct {
	name         string
	subscription *floodsub.Subscription

	messageBus chan net.Message

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler

	identifiersMutex            sync.Mutex
	transportToProtoIdentifiers map[net.TransportIdentifier]net.ProtocolIdentifier
	protoToTransportIdentifiers map[net.ProtocolIdentifier]net.TransportIdentifier
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Send(message net.TaggedMarshaler) error {
	return nil
}

func (c *channel) SendTo(
	recipientIdentifier interface{},
	message net.TaggedMarshaler,
) error {
	return nil
}

func (c *channel) Recv(h net.HandleMessageFunc) error {
	for message := range c.messageBus {
		if err := h(message); err != nil {
			return err
		}
	}
	return nil
}

func (c *channel) RegisterIdentifier(
	transportIdentifier net.TransportIdentifier,
	protocolIdentifier net.ProtocolIdentifier,
) error {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	if _, ok := transportIdentifier.(*identity); !ok {
		return fmt.Errorf(
			"incorrect type for transportIdentifier: [%v] in channel [%s]",
			transportIdentifier, c.name,
		)
	}

	if _, exists := c.transportToProtoIdentifiers[transportIdentifier]; exists {
		return fmt.Errorf(
			"protocol identifier in channel [%s] already associated with [%v]",
			c.name, transportIdentifier,
		)
	}
	if _, exists := c.protoToTransportIdentifiers[protocolIdentifier]; exists {
		return fmt.Errorf(
			"transport identifier in channel [%s] already associated with [%v]",
			c.name, protocolIdentifier,
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

func (c *channel) handleMessages() {
	defer c.subscription.Cancel()
	for {
		// TODO: thread in a context with cancel
		msg, err := c.subscription.Next(context.Background())
		if err != nil {
			// TODO: handle error - different error types
			// result in different outcomes
			fmt.Println(err)
			return
		}
		if err := c.processMessage(msg); err != nil {
			fmt.Println(err)
		}
		// TODO: handle message
		fmt.Println(msg)
	}
}

func (c *channel) processMessage(message *floodsub.Message) error {
	senderIdentifier := &identity{id: message.GetFrom()}

	c.identifiersMutex.Lock()
	protocolIdentifier := c.transportToProtoIdentifiers[senderIdentifier]
	c.identifiersMutex.Unlock()

	data := message.GetData()
	// 1. Unmarshall the message in to the Envelope
	// 2. Do the whole receiver thing -> senderIdentifier
	// 3. payload is a gossip message
	// 4. Unmarshall again to get the underlying gossip message
	// 5. Since the protocol type is on the gossip message, let's
	//    have an enum on the gossip message that keys us to the type of message that this is.
	// 6. Construct an internal.BasicMessage to fire back to the protocol
	// var payload interface{}

	// if err := proto.Unmarshal(data, payload); err != nil {
	// 	return err
	// }

	// unmarshaler, found := c.unmarshalersByType[payload.Type()]
	// if !found {
	// 	return fmt.Errorf("Couldn't find unmarshaler for type %s", payload.Type())
	// }

	// unmarshaled := unmarshaler()
	// if err := unmarshaled.Unmarshal(bytes); err != nil {
	// 	return err
	// }

	// protocolMessage :=
	// 	internal.BasicMessage(
	// 		senderIdentifier,
	// 		protocolIdentifier,
	// 		interface{},
	// 	)

	// c.messageBus <- data
	fmt.Println(protocolIdentifier, data)
	return nil
}
