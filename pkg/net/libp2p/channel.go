package libp2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	floodsub "github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p-peerstore"
)

type channel struct {
	name string

	clientIdentity *identity
	peerStore      peerstore.Peerstore

	pubsubMutex  sync.Mutex
	pubsub       *floodsub.PubSub
	subscription *floodsub.Subscription

	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc

	tempBufferLock sync.Mutex
	tempBuffer     []net.Message

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
	return c.doSend(message, c.clientIdentity)
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

	c.pubsubMutex.Lock()
	defer c.pubsubMutex.Unlock()

	// Publish the proto to the network
	return c.pubsub.Publish(c.name, envelopeBytes)
}

func (c *channel) SendTo(
	recipientIdentifier interface{},
	message net.TaggedMarshaler,
) error {
	return nil
}

func (c *channel) Recv(handler net.HandleMessageFunc) error {
	c.messageHandlersMutex.Lock()
	c.messageHandlers = append(c.messageHandlers, handler)
	c.messageHandlersMutex.Unlock()

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

	t := time.NewTimer(1) // first tick is immediate
	defer t.Stop()

	for {
		select {
		case <-t.C:
			msg, err := c.subscription.Next(ctx)
			if err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				fmt.Println(err)
				continue
			}

			if err := c.processMessage(msg); err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				fmt.Println(err)
				continue
			}

			t.Reset(1 * time.Second)
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
	unmarshaled, err := c.getUnmarshalingContainerByType(string(envelope.Type))
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

	return c.deliver(protocolMessage)
}

func (c *channel) getUnmarshalingContainerByType(envelopeType string) (net.TaggedUnmarshaler, error) {
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

	return c.transportToProtoIdentifiers[senderIdentifier.id], nil
}

func (c *channel) deliver(message net.Message) error {
	c.messageHandlersMutex.Lock()
	defer c.messageHandlersMutex.Unlock()

	// If we haven't registered a callback, buffer the message
	if len(c.messageHandlers) == 0 {
		c.tempBufferLock.Lock()
		c.tempBuffer = append(c.tempBuffer, message)
		c.tempBufferLock.Unlock()
		return nil
	}

	handlerSnapshot := make([]net.HandleMessageFunc, len(c.messageHandlers))
	copy(handlerSnapshot, c.messageHandlers)

	// Once we've registered a callback, drain the buffer
	if c.tempBuffer != nil {
		c.tempBufferLock.Lock()
		bufferSnapshot := make([]net.Message, len(c.tempBuffer))
		copy(bufferSnapshot, c.tempBuffer)
		c.tempBufferLock.Unlock()

		// Block so that we can clear the temporary buffer
		c.executeHandler(bufferSnapshot, handlerSnapshot)
		c.tempBuffer = nil

		return nil
	}

	// The usual case: for each message, execute against registered handlers
	go c.executeHandler([]net.Message{message}, handlerSnapshot)

	return nil
}

func (c *channel) executeHandler(messages []net.Message, snapshot []net.HandleMessageFunc) {
	for _, message := range messages {
		for _, handler := range snapshot {
			if err := handler(message); err != nil {
				fmt.Println(err)
			}
		}
	}

	// release copy to the gc
	snapshot = nil
}
