package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	floodsub "github.com/libp2p/go-floodsub"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

type channel struct {
	name string

	identity *identity
	store    pstore.Peerstore

	pubsubLock   sync.Mutex
	pubsub       *floodsub.PubSub
	subscription *floodsub.Subscription

	messagesLock sync.RWMutex
	messages     []net.Message

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
	return c.doSend(message, c.identity)
}

func envelopeProto(message net.TaggedMarshaler, sender *identity) ([]byte, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	identityBytes, err := sender.Marshal()
	if err != nil {
		return nil, err
	}

	return (&pb.Envelope{
		Payload: payloadBytes,
		Sender:  identityBytes,
		Type:    []byte(message.Type()),
	}).Marshal()
}

func (c *channel) doSend(message net.TaggedMarshaler, sender *identity) error {
	// Transform net.TaggedMarshaler to a protobuf message
	envelopeBytes, err := envelopeProto(message, sender)
	if err != nil {
		return err
	}

	c.pubsubLock.Lock()
	defer c.pubsubLock.Unlock()

	// Publish the proto to the network
	return c.pubsub.Publish(c.name, envelopeBytes)
}

func (c *channel) SendTo(
	recipientIdentifier interface{},
	message net.TaggedMarshaler,
) error {
	return nil
}

func (c *channel) Recv(h net.HandleMessageFunc) error {
	c.messagesLock.RLock()
	snapshot := make([]net.Message, len(c.messages))
	copy(snapshot, c.messages)
	// drain messages from buffer
	// FIXME: this will be a GC hotspot; use pools
	c.messages = make([]net.Message, 0)
	c.messagesLock.RUnlock()

	for _, message := range snapshot {
		if err := h(message); err != nil {
			return err
		}
	}

	snapshot = nil // release copy to the gc

	return nil
}

func (c *channel) RegisterIdentifier(
	transportIdentifier net.TransportIdentifier,
	protocolIdentifier net.ProtocolIdentifier,
) error {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	if _, ok := transportIdentifier.(networkIdentity); !ok {
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

func (c *channel) handleMessages(ctx context.Context) {
	defer c.subscription.Cancel()

	for {
		// TODO: thread in a context with cancel
		msg, err := c.subscription.Next(ctx)
		if err != nil {
			// TODO: handle error - different error types
			// result in different outcomes
			fmt.Println(err)
			return
		}

		if err := c.processMessage(msg); err != nil {
			// TODO: handle error - different error types
			// result in different outcomes
			fmt.Println(err)
			return
		}

		select {
		case <-ctx.Done():
			return
		}
	}
}

func (c *channel) processMessage(message *floodsub.Message) error {
	var envelope pb.Envelope
	if err := proto.Unmarshal(message.Data, &envelope); err != nil {
		return err
	}

	// TODO: handle receivers, authentication, etc

	// The protocol type is on the envelope; let's pull that type
	// from our map of unmarshallers.
	unmarshaled, err := c.getUnmarshalerByType(string(envelope.Type))
	if err != nil {
		return err
	}

	if err := unmarshaled.Unmarshal(envelope.GetPayload()); err != nil {
		return err
	}

	// Construct an identifier from the sender (on the message)
	senderIdentifier := &identity{}
	if err := senderIdentifier.Unmarshal(envelope.Sender); err != nil {
		return err
	}

	// Get the associated protocol identifier from an association map
	protocolIdentifier, err := c.getProtocolIdentifier(senderIdentifier)
	if err != nil {
		return err
	}

	// Fire a message back to the protocol
	protocolMessage := internal.BasicMessage(senderIdentifier.id,
		protocolIdentifier, unmarshaled,
	)

	// We'll drain the list of messages when called
	c.messagesLock.Lock()
	c.messages = append(c.messages, protocolMessage)
	c.messagesLock.Unlock()

	return nil
}

func (c *channel) getUnmarshalerByType(envelopeType string) (net.TaggedUnmarshaler, error) {
	c.unmarshalersMutex.Lock()
	defer c.unmarshalersMutex.Unlock()

	unmarshaler, found := c.unmarshalersByType[envelopeType]
	if !found {
		return nil, fmt.Errorf(
			"Couldn't find unmarshaler for type %s", envelopeType,
		)
	}

	return unmarshaler(), nil
}

func (c *channel) getProtocolIdentifier(senderIdentifier *identity) (net.ProtocolIdentifier, error) {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	protocolIdentifier, found := c.transportToProtoIdentifiers[senderIdentifier.id]
	if !found {
		return nil, fmt.Errorf(
			"Couldn't find protocol identifier for sender identifier %v",
			senderIdentifier,
		)
	}
	return protocolIdentifier, nil
}
